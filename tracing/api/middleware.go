package apitracing

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/Cloudera-Sz/golang-micro/tracing"
)

// ApiTracer;
func ApiTracer(operationPrefix []byte) gin.HandlerFunc {
	if operationPrefix == nil {
		operationPrefix = []byte("api-request-")
	}
	return func(c *gin.Context) {
		// all before request is handled
		var span opentracing.Span
		if cspan, ok := c.Get(tracing.TracingContext); ok {
			span = StartSpanWithParent(cspan.(opentracing.Span).Context(), string(operationPrefix)+c.Request.Method, c.Request.Method, c.Request.URL.Path)

		} else {
			span = StartSpanWithHeader(&c.Request.Header, string(operationPrefix)+c.Request.Method, c.Request.Method, c.Request.URL.Path)
		}
		defer span.SetTag(string(ext.HTTPStatusCode), c.Writer.Status())
		defer span.Finish()
		c.Set(tracing.TracingContext, span)
		c.Next()
	}
}

func OpenTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		wireCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		serverSpan := opentracing.StartSpan(c.Request.URL.Path, ext.RPCServerOption(wireCtx))
		defer serverSpan.Finish()
		//span.SetTag("spanContext", span.Context())
		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), serverSpan))
		c.Next()
	}
}

func OpenTracing2() gin.HandlerFunc {
	return func(c *gin.Context) {
		var parentCtx opentracing.SpanContext
		parentSpan := opentracing.SpanFromContext(c)
		if parentSpan != nil {
			parentCtx = parentSpan.Context()
		}

		span := opentracing.GlobalTracer().StartSpan(
			c.Request.URL.Path,
			opentracing.ChildOf(parentCtx),
			opentracing.Tag{Key: string(ext.Component), Value: "api"},
			ext.SpanKindRPCClient,
		)
		defer span.Finish()

		err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span.LogFields(log.String("inject-error", err.Error()))
		}
		c.Next()
	}
}
