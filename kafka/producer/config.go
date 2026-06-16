package producer

import (
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
)

type config struct {
	opts   []kgo.Opt
	logger kgo.Logger
	sasl   sasl.Mechanism
}

func newCfg() config {
	return config{
		opts: make([]kgo.Opt, 0),

		logger: new(nopLogger),
	}
}
