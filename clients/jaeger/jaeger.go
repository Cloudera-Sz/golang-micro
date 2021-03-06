package jaeger

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/Cloudera-Sz/golang-micro/clients/config"
	"github.com/Cloudera-Sz/golang-micro/clients/etcd"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-client-go/rpcmetrics"
	"github.com/uber/jaeger-lib/metrics"

	tags "github.com/opentracing/opentracing-go/ext"
)

//InitTracer ..
func InitTracer(config *config.JaegerConfig, serviceName string) (opentracing.Tracer, io.Closer, error) {
	if config == nil {
		return nil, nil, errors.New("jaeger config is not exist")
	}
	if serviceName == "" {
		return nil, nil, errors.New("service name is empty")
	}
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:                    jaeger.SamplerTypeConst,
			Param:                   1,
			SamplingServerURL:       config.Sampler.HostPort,
			SamplingRefreshInterval: config.Sampler.RefreshInterval,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            config.Reporter.LogSpans,
			LocalAgentHostPort:  config.Reporter.HostPort,
			QueueSize:           config.Reporter.QueueSize,
			BufferFlushInterval: config.Reporter.FlushInterval,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	return cfg.New(
		serviceName,
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Observer(rpcmetrics.NewObserver(jMetricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
}

//InitGlobalTracer ...
func InitGlobalTracer(config *config.JaegerConfig) (io.Closer, error) {
	if config == nil {
		return nil, errors.New("jaeger config is not exist")
	}
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:                    jaeger.SamplerTypeConst,
			Param:                   1,
			SamplingServerURL:       config.Sampler.HostPort,
			SamplingRefreshInterval: config.Sampler.RefreshInterval,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            config.Reporter.LogSpans,
			LocalAgentHostPort:  config.Reporter.HostPort,
			QueueSize:           config.Reporter.QueueSize,
			BufferFlushInterval: config.Reporter.FlushInterval,
		},
	}
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory
	return cfg.InitGlobalTracer(
		config.ServiceName,
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Observer(rpcmetrics.NewObserver(jMetricsFactory, rpcmetrics.DefaultNameNormalizer)),
	)
}

//NewJaegerConfig get config from etcd
func NewJaegerConfig(etcdCli *etcd.Client, appName, profile string) (string, *config.JaegerConfig) {
	if appName == "" {
		appName = os.Getenv("APP_NAME")
	}
	if profile == "" {
		profile = os.Getenv("PROFILE")
	}
	jaegerKey := etcdCli.GetEtcdKey(profile, appName, "jaeger")
	jaegerConfig := new(config.JaegerConfig)
	err := etcdCli.GetValue(5*time.Second, jaegerKey, jaegerConfig)
	if err != nil {
		log.Println("jaeger config is not exist:", err)
		jaegerConfig = nil
	}
	return jaegerKey, jaegerConfig
}

//InitTracerFromEtcd ...
func InitTracerFromEtcd(etcdCli *etcd.Client, appName, profile string, serviceName string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	jaegerKey, jaegerConfig := NewJaegerConfig(etcdCli, appName, profile)
	tracer, c, err := InitTracer(jaegerConfig, serviceName)
	if err != nil {
		log.Println("jaeger connect failed")
	}
	closer = c
	log.Println("jaeger inited")
	//change to new db connection when  config changed
	go etcdCli.WatchKey(jaegerKey, jaegerConfig, func() {
		t, c, err := InitTracer(jaegerConfig, serviceName)
		if err != nil {
			log.Println("jaeger connect failed")
			return
		}
		closer.Close()
		closer = c
		tracer = t
		log.Println("jaeger changed")
	})
	return tracer, closer, err
}

//InitGlobalTracerFromEtcd ...
func InitGlobalTracerFromEtcd(etcdCli *etcd.Client, appName, profile string) (closer io.Closer, err error) {
	jaegerKey, jaegerConfig := NewJaegerConfig(etcdCli, appName, profile)
	c, err := InitGlobalTracer(jaegerConfig)
	if err != nil {
		log.Println("jaeger connect failed")
	}
	closer = c
	log.Println("jaeger inited")
	//change to new db connection when  config changed
	go etcdCli.WatchKey(jaegerKey, jaegerConfig, func() {
		c, err := InitGlobalTracer(jaegerConfig)
		if err != nil {
			log.Println("jaeger connect failed")
			return
		}
		closer.Close()
		closer = c
		log.Println("jaeger changed")
	})
	return closer, err
}

//GetSpan 获取span
func GetSpan(ctx context.Context, name string, f func(span opentracing.Span)) (context.Context, opentracing.Span) {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan(name, opentracing.ChildOf(span.Context()))
	} else {
		span = opentracing.StartSpan(name)
	}
	f(span)
	ctx = opentracing.ContextWithSpan(ctx, span)
	return ctx, span
}

//GetDefaultSpan 获取默认 span
func GetDefaultSpan(ctx context.Context, name string) (context.Context, opentracing.Span) {
	return GetSpan(ctx, name, func(span opentracing.Span) {
	})
}

//GetServerSpan 获取server span
func GetServerSpan(ctx context.Context, name string) (context.Context, opentracing.Span) {
	return GetSpan(ctx, name, func(span opentracing.Span) {
		tags.SpanKindRPCServer.Set(span)
	})
}
