package queue

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TaskBroadcastSend = "broadcast.send"
)

type Client struct {
	client *asynq.Client
}

func NewClient(redisURL string) (*Client, error) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	return &Client{client: asynq.NewClient(opt)}, nil
}

func (c *Client) Enqueue(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return c.client.EnqueueContext(ctx, task, opts...)
}

func (c *Client) Close() error {
	return c.client.Close()
}

func NewBroadcastSendTask(broadcastID string) (*asynq.Task, error) {
	payload := fmt.Sprintf(`{"broadcast_id":"%s"}`, broadcastID)
	return asynq.NewTask(TaskBroadcastSend, []byte(payload)), nil
}

type Worker struct {
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewWorker(redisURL string, concurrency int) (*Worker, error) {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}
	srv := asynq.NewServer(opt, asynq.Config{
		Concurrency: concurrency,
		Queues: map[string]int{
			"default": 10,
		},
	})
	return &Worker{server: srv, mux: asynq.NewServeMux()}, nil
}

func (w *Worker) Handle(pattern string, handler asynq.HandlerFunc) {
	w.mux.HandleFunc(pattern, handler)
}

func (w *Worker) Run() error {
	return w.server.Run(w.mux)
}

func (w *Worker) Shutdown() {
	w.server.Shutdown()
}

func PingRedis(redisURL string) error {
	opt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return err
	}
	inspector := asynq.NewInspector(opt)
	defer inspector.Close()
	_, err = inspector.Queues()
	return err
}
