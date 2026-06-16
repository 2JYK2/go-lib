package consumer

import (
	"context"
	"fmt"
	"github.com/2JYK2/go-lib/kafka/contract"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kotel"
	"sync"
	"time"
)

type PartitionConsumer struct {
	ctx       context.Context
	cfg       *config
	cl        *kgo.Client
	tracer    *kotel.Tracer
	instances map[tp]*instance
	lock      sync.Mutex
	done      chan struct{}
	errs      chan error
}

type (
	instance struct {
		processor contract.Processor
		cfg       *config
		cl        *kgo.Client
		p         *PartitionConsumer

		topic     string
		partition int32
		recv      chan kgo.FetchTopicPartition

		quit chan struct{}
		done chan struct{}
	}

	tp struct {
		t string
		p int32
	}
)

func (i *instance) consume() {
	defer close(i.done)
	i.cfg.logger.Log(kgo.LogLevelDebug, "starting, t %s p %d\n", i.topic, i.partition)
	defer func() {
		i.cfg.logger.Log(kgo.LogLevelDebug, "killing, t %s p %d\n", i.topic, i.partition)
	}()
	for {
		select {
		case <-i.quit:
			return
		case p := <-i.recv:
			for _, record := range p.Records {
				func() {
					if i.p.tracer != nil {
						_, span := i.p.tracer.WithProcessSpan(record)
						//fmt.Println("#########", string(record.Value))
						defer span.End()
					}
					i.processor.Process(record.Value)
				}()
			}
			i.cl.MarkCommitRecords(p.Records...)
		}
	}
}

func (c *PartitionConsumer) assigned(_ context.Context, client *kgo.Client, m map[string][]int32) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for topic, partitions := range m {
		for _, partition := range partitions {
			pc := &instance{
				cl:        client,
				p:         c,
				topic:     topic,
				partition: partition,
				cfg:       c.cfg,
				processor: c.cfg.processor,

				done: make(chan struct{}),
				quit: make(chan struct{}),
				recv: make(chan kgo.FetchTopicPartition, 10),
			}
			c.instances[tp{topic, partition}] = pc
			go pc.consume()
		}
	}
}

func (c *PartitionConsumer) lost(_ context.Context, _ *kgo.Client, m map[string][]int32) {
	c.errs <- fmt.Errorf("function lost be called, check network")
	// if network down, lost will be called, function poll will panic because tp be remmoved
	//c.killInstance(m)
}

func (c *PartitionConsumer) revoked(ctx context.Context, client *kgo.Client, m map[string][]int32) {
	c.killInstance(m)
	if err := client.CommitMarkedOffsets(ctx); err != nil {
		c.errs <- fmt.Errorf("commitMarkedOffsets failed in OnPartitionsRevoked, errMsg: %s", err)
		c.cfg.logger.Log(kgo.LogLevelError, "commitMarkedOffsets failed in OnPartitionsRevoked", "because", err)
	}
}

func (c *PartitionConsumer) killInstance(m map[string][]int32) {
	var wg sync.WaitGroup
	defer wg.Wait()
	c.lock.Lock()
	defer c.lock.Unlock()
	for t, pars := range m {
		for _, partition := range pars {
			key := tp{t, partition}
			pc := c.instances[key]
			delete(c.instances, key)
			close(pc.quit)
			c.cfg.logger.Log(kgo.LogLevelInfo, fmt.Sprintf("waiting for work to finish t %s p %d\n", t, partition))
			wg.Add(1)
			go func() { <-pc.done; wg.Done() }()
		}
	}
}

func New(ctx context.Context, opts ...Option) (*PartitionConsumer, error) {
	cfg := newCfg()
	for _, o := range opts {
		cfg = o.apply(cfg)
	}
	return newConsumer(ctx, cfg)
}

func newConsumer(ctx context.Context, cfg config) (*PartitionConsumer, error) {
	c := &PartitionConsumer{ctx: ctx, cfg: &cfg, tracer: cfg.tracer, done: make(chan struct{}), instances: make(map[tp]*instance), errs: make(chan error, 100)}
	defaults := []kgo.Opt{
		//kgo.DisableAutoCommit(),
		kgo.BlockRebalanceOnPoll(),
		kgo.AutoCommitMarks(),
		kgo.OnPartitionsRevoked(c.revoked),
		kgo.OnPartitionsAssigned(c.assigned),
		kgo.OnPartitionsLost(c.lost),
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

func (c *PartitionConsumer) stop() {
	close(c.done)
	// close partition consumer
	var wg sync.WaitGroup
	defer wg.Wait()
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, ins := range c.instances {
		close(ins.quit)
		wg.Add(1)
		go func(i *instance) { <-i.done; wg.Done() }(ins)
	}
}

func (c *PartitionConsumer) poll() {
	for {
		fetches := c.cl.PollRecords(context.Background(), c.cfg.fetchCount)
		if fetches.IsClientClosed() {
			c.errs <- fmt.Errorf("client is closed")
			time.Sleep(time.Second)
		}
		//fetches.EachError(func(topic string, partition int32, err error) {
		//c.errs <- fmt.Errorf("some error occur on pollRecords,topic: %s, partition: %d, errMsg: %s", topic, partition, err)
		//})
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			select {
			case <-c.ctx.Done():
				return
			case <-c.done:
				return
			default:
				if p.Err != nil {
					c.errs <- fmt.Errorf("some error occur on pollRecords,topic: %s, partition: %d, errMsg: %s", p.Topic, p.Partition, p.Err)
				} else {
					t := tp{p.Topic, p.Partition}
					c.instances[t].recv <- p
				}
			}
			c.cl.AllowRebalance()
		})
	}
}

func (c *PartitionConsumer) Start() {
	c.poll()
}

func (c *PartitionConsumer) Stop() {
	c.stop()
}

func (c *PartitionConsumer) ListenErrs() <-chan error {
	return c.errs
}
