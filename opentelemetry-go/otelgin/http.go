package otelgin

import (
	"bytes"
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: time.Second * 5}

func SetHttpClient(c *http.Client) {
	client = c
}

func DoRequest(ctx context.Context, req *http.Request, cb func(body []byte) error) error {
	buf := bytes.NewBuffer([]byte{})
	if req.Body != nil {
		req.Body = io.NopCloser(io.TeeReader(req.Body, buf))
	}
	data := make([]byte, 0)

	err := func() error {
		pp := otel.GetTextMapPropagator()
		carrier := propagation.MapCarrier{}
		pp.Inject(getContext(ctx), carrier)
		for k, v := range carrier {
			req.Header.Set(k, v)
		}
		rsp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer rsp.Body.Close()
		if rsp.StatusCode != http.StatusOK {
			return fmt.Errorf(fmt.Sprintf("http request failed, statusCode is %d", rsp.StatusCode))
		}
		data, _ = io.ReadAll(rsp.Body)
		return cb(data)
	}()
	// add tracer info
	span := trace.SpanFromContext(getContext(ctx))
	span.SetAttributes(
		semconv.HTTPMethodKey.String(req.Method),
		semconv.HTTPURLKey.String(req.URL.String()),
	)
	if len(data) > 0 {
		span.SetAttributes(attribute.Key("http.response").String(string(data)))
	}
	if buf.Len() > 0 {
		span.SetAttributes(attribute.Key("http.payload").String(buf.String()))
	}
	return err
}
