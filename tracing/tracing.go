package tracing

import (
	"context"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var TracingContext = "tracing-context"

// 创建追踪器;
func NewTracer(serviceName string, jagentHost string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	// 设置jaeger;
	jcfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  jagentHost,
		},
	}
	// 创建追踪器;
	tracer, closer, err = jcfg.New(
		serviceName,
		jaegercfg.Logger(jaeger.StdLogger),
	)
	if err != nil {
		return
	}
	// 设置追踪器;
	opentracing.SetGlobalTracer(tracer)
	return
}

func GetSpanID(c context.Context) opentracing.SpanContext {
	if cspan := c.Value(TracingContext); cspan != nil {
		return cspan.(opentracing.Span).Context()
	}
	return nil
}
