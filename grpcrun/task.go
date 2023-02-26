// Package grpcrun @Author: Ciusyan 2023/2/26
package grpcrun

import (
	"context"
	"fmt"
	"reflect"
)

// GrpcTask ：去调用GRPC的方法
func GrpcTask(f any, ctx *context.Context, req any) (any, error) {

	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return nil, fmt.Errorf("参数错误")
	}

	// 调用参数
	argv := make([]reflect.Value, 2)
	argv[0] = reflect.ValueOf(*ctx)
	argv[1] = reflect.ValueOf(req)

	// 具体调用
	res := v.Call(argv)

	// 返回错误
	errV := res[1].Interface()
	if errV != nil {
		return res[0].Interface(), errV.(error)
	}

	return res[0].Interface(), nil
}
