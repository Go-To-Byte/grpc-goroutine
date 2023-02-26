// Package grpcrun @Author: Ciusyan 2023/2/26
package grpcrun

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"reflect"
)

// GrpcTask 用于构建Grpc Task
type GrpcTask struct {
	// 必须符合GRPC的 Method 签名
	grpcMethod any

	// GRPC的调用参数
	ctx     *context.Context
	request any

	// GRPC的调用返回值
	Name     string
	Response any
	Err      error

	// 日志对象
	log *zap.SugaredLogger
}

// NewGrpcTask creates a new GrpcTask
//
// Note:
// @param grpcName string name of the grpc, this should be unique
func NewGrpcTask(ctx *context.Context, grpcName string, grpcMethod any, request any) *GrpcTask {
	mu.Lock()
	defer mu.Unlock()
	zap.S()

	if grpcName == "" {
		grpcName = node.Generate().String()
	}

	return &GrpcTask{
		ctx:        ctx,
		grpcMethod: grpcMethod,
		request:    request,
		Name:       grpcName,
		log:        zap.S().Named(grpcName),
	}
}

// Call 去调用GRPC的方法
func (c *GrpcTask) Call() {
	// 进行参数校验
	c.validate()
	if c.Err != nil {
		return
	}

	// 能来到这里，参数一定正确了，进行方法调用
	// 形如：Login(ctx context.Context, req *loginReq) (*loginResp, error)
	c.call()
}

func (c *GrpcTask) call() {
	v := reflect.ValueOf(c.grpcMethod)

	// 调用参数
	argv := make([]reflect.Value, 2)
	argv[0] = reflect.ValueOf(*c.ctx)
	argv[1] = reflect.ValueOf(c.request)

	// 反射调用
	res := v.Call(argv)

	// 给返回值 赋值
	c.Response = res[0].Interface()
	if res[1].Interface() != nil {
		c.Err = res[1].Interface().(error)
	}
}

// 校验结构体
func (c *GrpcTask) validate() {

	// 校验 req 类型
	reqV := reflect.ValueOf(c.request)
	if !reqV.IsValid() || reqV.Kind() != reflect.Ptr {
		c.Err = fmt.Errorf("请正确的传递[request]，不支持：%v", reqV.Kind())
		return
	}

	// 校验 ctx 类型
	ctxV := reflect.ValueOf(c.ctx).Elem()
	if ctxV.IsNil() {
		c.Err = fmt.Errorf("请正确的传递[Context]，不支持：nil")
		return
	}

	// 校验 grpcMethod 的信息
	methodV := reflect.ValueOf(c.grpcMethod)
	if methodV.Kind() != reflect.Func {
		c.Err = fmt.Errorf("[grpcMethod]必须是一个GRPC的函数类型，现在是：%v", methodV.Kind())
		return
	}

	// 简单校验参数类型
	methodT := methodV.Type()
	if methodT.NumIn() != 2 {
		c.Err = fmt.Errorf("[grpcMethod]必须有2个参数(context.Context, *request)")
		return
	}

	// 校验 grpcMethod 的第一个参数
	if methodT.In(0).Kind() != reflect.Interface || methodT.In(0).Name() != "Context" {
		c.Err = fmt.Errorf("[grpcMethod]的第1个参数必须是：context.Context")
		return
	}

	if reqV.Type().Elem() != methodT.In(1).Elem() {
		c.Err = fmt.Errorf("[request]的参数与[grpcMethod]的参数不匹配：grpcMethod = %v, request = %v",
			methodT.In(1).Elem(), reqV.Type().Elem())
		return
	}

	if methodT.NumOut() != 2 {
		c.Err = fmt.Errorf("[grpcMethod]必须有2个返回值(*Response, error)")
		return
	}

	if methodT.Out(1).Kind() != reflect.Interface || methodT.Out(1).Name() != "error" {
		c.Err = fmt.Errorf("[grpcMethod]的第2个返回值必须是：error")
		return
	}
}
