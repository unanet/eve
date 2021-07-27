package queue

import (
	"context"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"

	"github.com/unanet/go/pkg/log"
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
			worker.log.Info("queue worker stopped")
			close(worker.done)
			return
		default:
			ctx := context.Background()
			m, err := worker.q.Receive(ctx)
			if err != nil {
				worker.log.Panic("error receiving message from queue", zap.Error(err))
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

func (worker *Worker) run(h Handler, mCtx []*mContext) {
	numMessages := len(mCtx)
	var wg sync.WaitGroup
	wg.Add(numMessages)
	for _, mc := range mCtx {
		go func(m *mContext) {
			ctx, cancel := context.WithTimeout(m.ctx, worker.timeout)
			defer cancel()
			defer wg.Done()
			if err := h.HandleMessage(ctx, &m.M); err != nil {
				worker.log.Error("error handling message", zap.Error(err))
			}
		}(mc)
	}
	wg.Wait()
}

func GetLogger(ctx context.Context) *zap.Logger {
	reqID := log.GetReqID(ctx)
	if len(reqID) > 0 {
		return log.Logger.With(zap.String("req_id", reqID))
	} else {
		return log.Logger
	}
}
