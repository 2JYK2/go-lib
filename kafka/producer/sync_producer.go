package producer

import (
	"context"
	"github.com/twmb/franz-go/pkg/kgo"
)

type SyncProducer struct {
	ctx context.Context
	cfg *config
	cl  *kgo.Client
}

func NewSyncProducer(ctx context.Context, opts ...Option) (*SyncProducer, error) {
	cfg := newCfg()
	for _, o := range opts {
		cfg = o.apply(cfg)
	}
	return newSyncProducer(ctx, cfg)
}

func newSyncProducer(ctx context.Context, cfg config) (*SyncProducer, error) {
	c := &SyncProducer{ctx: ctx, cfg: &cfg}
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

func (c *SyncProducer) Publish(records ...*kgo.Record) error {
	res := c.cl.ProduceSync(c.ctx, records...)
	return res.FirstErr()
}

func (c *SyncProducer) Start() {
	//
}

func (c *SyncProducer) Stop() {
	//
}
