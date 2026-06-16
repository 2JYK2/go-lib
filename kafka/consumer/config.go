package consumer

import (
	"github.com/2JYK2/go-lib/kafka/contract"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/plugin/kotel"
)

type config struct {
	processor contract.Processor
	opts      []kgo.Opt
	logger    kgo.Logger
	sasl      sasl.Mechanism
	tracer    *kotel.Tracer

	topic      []string
	fetchCount int
}

func newCfg() config {
	return config{
		opts: make([]kgo.Opt, 0),

		fetchCount: 500,
		logger:     new(nopLogger),
	}
}
