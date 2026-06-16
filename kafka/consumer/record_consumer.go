package consumer

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type RecordConsumer struct {
	ctx context.Context
	cfg *config
	cl  *kgo.Client

	done chan struct{}
	errs chan error
}

func (c *RecordConsumer) revoked(ctx context.Context, client *kgo.Client, m map[string][]int32) {
	if err := client.CommitMarkedOffsets(ctx); err != nil {
		c.errs <- fmt.Errorf("commitMarkedOffsets failed in OnPartitionsRevoked, errMsg: %s", err)
		c.cfg.logger.Log(kgo.LogLevelError, "commitMarkedOffsets failed in OnPartitionsRevoked", "because", err)
	}
}

func NewRecordConsumer(ctx context.Context, opts ...Option) (*RecordConsumer, error) {
	cfg := newCfg()
	for _, o := range opts {
		cfg = o.apply(cfg)
	}
	return newRecordConsumer(ctx, cfg)
}

func newRecordConsumer(ctx context.Context, cfg config) (*RecordConsumer, error) {
	c := &RecordConsumer{ctx: ctx, cfg: &cfg, done: make(chan struct{}), errs: make(chan error, 100)}
	defaults := []kgo.Opt{
		//kgo.DisableAutoCommit(),
		kgo.BlockRebalanceOnPoll(),
		kgo.AutoCommitMarks(),
		kgo.OnPartitionsRevoked(c.revoked),
		kgo.RequireStableFetchOffsets(),
	}
	defaults = append(defaults, cfg.opts...)
	client, err := kgo.NewClient(defaults...)
	if err != nil {
		return nil, err
	}
	c.cl = client
	if err = c.cl.Ping(c.ctx); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *RecordConsumer) stop() {
	close(c.done)
}

func (c *RecordConsumer) poll() {
	for {
		fetches := c.cl.PollRecords(context.Background(), c.cfg.fetchCount)
		if fetches.IsClientClosed() {
			c.errs <- fmt.Errorf("client is closed")
			time.Sleep(time.Second)
		}
		fetches.EachError(func(topic string, partition int32, err error) {
			c.errs <- fmt.Errorf("some error occur on pollRecords,topic: %s, partition: %d, errMsg: %s", topic, partition, err)
		})
		fetches.EachRecord(func(record *kgo.Record) {
			select {
			case <-c.ctx.Done():
				return
			case <-c.done:
				return
			default:
				c.cfg.processor.Process(record.Value)
				c.cl.MarkCommitRecords(record)
			}
		})
		c.cl.AllowRebalance()
	}
}

func (c *RecordConsumer) Start() {
	c.poll()
}

func (c *RecordConsumer) Stop() {
	c.stop()
}

func (c *RecordConsumer) ListenErrs() <-chan error {
	return c.errs
}
