package wrapper

import (
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// 获取拨号追踪设置;
func GetDialOption() []grpc.DialOption {
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	dialOpts = append(dialOpts, DialOption(opentracing.GlobalTracer()))
	return dialOpts
}

// 获取服务端追踪设置;
func GetServerOption() []grpc.ServerOption {
	var servOpts []grpc.ServerOption
	servOpts = append(servOpts, ServerOption(opentracing.GlobalTracer()))
	return servOpts
}
