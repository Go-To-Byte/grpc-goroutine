// Package grpcrun @Author: Ciusyan 2023/2/26
package grpcrun

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"reflect"
)

// Grpc 用于构建Grpc请求
type Grpc struct {
	// 必须符合GRPC的 Method 签名
	grpcMethod any

	// GRPC的调用参数
	ctx     *context.Context
	request any

	// GRPC的调用返回值
	Response any
	Err      error

	// 日志对象
	log *zap.SugaredLogger
}

func NewGrpc(ctx *context.Context, grpcMethod any, req any) *Grpc {
	zap.S()
	return &Grpc{
		ctx:        ctx,
		grpcMethod: grpcMethod,
		request:    req,
		log:        zap.S().Named("Grpc-Task"),
	}
}

// GrpcTask ：去调用GRPC的方法
func (c *Grpc) GrpcTask() {
	v := reflect.ValueOf(c.grpcMethod)
	if v.Kind() != reflect.Func {
		c.Err = fmt.Errorf("[grpcMethod]必须是一个GRPC的函数类型")
		return
	}

	// 进行方法调用
	argv := make([]reflect.Value, 2)
	argv[0] = reflect.ValueOf(*c.ctx)
	argv[1] = reflect.ValueOf(c.request)
	res := v.Call(argv)

	c.Response = res[0].Interface()
	if res[1].Interface() != nil {
		c.Err = res[1].Interface().(error)
	}
}
