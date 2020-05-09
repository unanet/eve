package queue

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"

	"gitlab.unanet.io/devops/eve/pkg/log"
)

type ctxKeyRequestID int

const (
	RequestIDKey ctxKeyRequestID = 0
)

// HandlerFunc is used to define the Handler that is run on for each message
type HandlerFunc func(ctx context.Context, msg *M) error

// HandleMessage wraps a function for handling sqs messages
func (f HandlerFunc) HandleMessage(ctx context.Context, msg *M) error {
	return f(ctx, msg)
}

type Handler interface {
	HandleMessage(ctx context.Context, msg *M) error
}

type Worker struct {
	q       *Q
	log     *zap.Logger
	name    string
	timeout time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
	done    chan bool
	wqs     map[string]*Q
	mutex   sync.Mutex
	sess    *session.Session
}

func NewWorker(name string, q *Q, timeout time.Duration) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	w := Worker{
		name:    name,
		q:       q,
		log:     log.Logger.With(zap.Uint64("internal_queue_id", q.id), zap.String("worker", name)),
		timeout: timeout,
		ctx:     ctx,
		cancel:  cancel,
		sess:    q.sess,
		done:    make(chan bool),
		wqs:     make(map[string]*Q),
	}

	return &w
}

func (worker *Worker) Start(h Handler) {
	worker.log.Info("Queue worker started")
	for {
		select {
		case <-worker.ctx.Done():
			worker.log.Info("Queue worker stopped")
			close(worker.done)
			return
		default:
			m, err := worker.q.Receive(worker.ctx)
			if err != nil {
				worker.log.Panic("Error receiving message from queue", zap.Error(err))
			}
			if len(m) == 0 {
				continue
			}
			worker.run(h, m)
		}
	}
}

func (worker *Worker) Stop() {
	worker.cancel()
	<-worker.done
}

func (worker *Worker) DeleteMessage(ctx context.Context, m *M) error {
	return worker.q.Delete(ctx, m)
}

func (worker *Worker) getQueue(qUrl string) *Q {
	worker.mutex.Lock()
	defer worker.mutex.Unlock()
	if val, ok := worker.wqs[qUrl]; ok {
		return val
	}
	q := NewQ(worker.sess, Config{
		QueueURL: qUrl,
	})
	worker.wqs[qUrl] = q
	return q
}

func (worker *Worker) Message(ctx context.Context, qUrl string, m *M) error {
	q := worker.getQueue(qUrl)
	return q.Message(ctx, m)
}

func (worker *Worker) run(h Handler, messages []*M) {
	numMessages := len(messages)
	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *M) {
			ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), RequestIDKey, m.ReqID), worker.timeout)
			defer cancel()
			defer wg.Done()
			if err := h.HandleMessage(ctx, m); err != nil {
				worker.log.Error("Error handling message", zap.Error(err))
			}
		}(messages[i])
	}
	wg.Wait()
}

func GetReqID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		return reqID
	}
	return ""
}

func GetLogger(ctx context.Context) *zap.Logger {
	reqID := GetReqID(ctx)
	if len(reqID) > 0 {
		return log.Logger.With(zap.String("req_id", reqID))
	} else {
		return log.Logger
	}
}
