package producer

import "github.com/twmb/franz-go/pkg/kgo"

type nopLogger struct{}

func (*nopLogger) Level() kgo.LogLevel { return kgo.LogLevelNone }
func (*nopLogger) Log(kgo.LogLevel, string, ...any) {
}
