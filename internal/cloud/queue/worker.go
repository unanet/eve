package queue

import (
	"context"
	"sync"
	"time"

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
}

func NewWorker(name string, q *Q, timeout time.Duration) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		name:    name,
		q:       q,
		log:     log.Logger.With(zap.Uint64("internal_queue_id", q.id), zap.String("worker", name)),
		timeout: timeout,
		ctx:     ctx,
		cancel:  cancel,
		done:    make(chan bool),
	}
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
				worker.log.Error("Error receiving message from queue", zap.Error(err))
				continue
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

func (worker *Worker) run(h Handler, messages []*M) {
	numMessages := len(messages)
	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *M) {
			ctx, cancel := context.WithTimeout(context.WithValue(context.Background(), RequestIDKey, m.ReqID), time.Duration(120)*time.Second)
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
