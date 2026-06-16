package producer

import (
	"context"
	"github.com/twmb/franz-go/pkg/kgo"
)

type AsyncProducer struct {
	ctx context.Context
	cfg *config
	cl  *kgo.Client

	errs chan error
	recv chan *kgo.Record
	done chan struct{}
}

func NewAsyncProducer(ctx context.Context, opts ...Option) (*AsyncProducer, error) {
	cfg := newCfg()
	for _, o := range opts {
		cfg = o.apply(cfg)
	}
	return newAsyncProducer(ctx, cfg)
}

func newAsyncProducer(ctx context.Context, cfg config) (*AsyncProducer, error) {
	c := &AsyncProducer{ctx: ctx, cfg: &cfg, done: make(chan struct{}), recv: make(chan *kgo.Record, 100)}
	//defaults := []kgo.Opt{
	//kgo.SASL(cfg.sasl),
	//}
	//defaults = append(defaults, cfg.opts...)
	client, err := kgo.NewClient(cfg.opts...)
	if err != nil {
		return nil, err
	}
	c.cl = client
	if err = c.cl.Ping(c.ctx); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *AsyncProducer) stop() {
	close(c.done)
}

func (c *AsyncProducer) listen() {
	i := 0
	for {
		select {
		case record := <-c.recv:
			i++
			c.cl.Produce(record.Context, record, func(record *kgo.Record, err error) {
				if err != nil {
					c.errs <- err
				}
			})
		case <-c.ctx.Done():
			return
		case <-c.done:
			return
		}
	}
}
func (c *AsyncProducer) ListenErrs() <-chan error {
	return c.errs
}

func (c *AsyncProducer) Publish(record *kgo.Record) {
	c.recv <- record
}

func (c *AsyncProducer) Start() {
	go c.listen()
}

func (c *AsyncProducer) Stop() {
	c.stop()
	if err := c.cl.Flush(c.ctx); err != nil {
		c.errs <- err
	}
}

func (c *AsyncProducer) Flush() {
	if err := c.cl.Flush(c.ctx); err != nil {
		c.errs <- err
	}
}
