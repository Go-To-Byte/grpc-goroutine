// @Author: Ciusyan 2023/2/26
package grpcrun_test

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"strconv"
	"testing"
	"time"

	"github.com/Go-To-Byte/grpc-goroutine/grpcrun"
)

// 作为GRPC请求的参数
type loginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 作为GRPC响应的参数
type loginResp struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

// GRPC请求的方法参数
func Login(ctx context.Context, req *loginReq) (*loginResp, error) {
	if req.Username != "ciusyan" || req.Password != "222" {
		return nil, fmt.Errorf("登录失败")
	}
	fmt.Println("登录成功")
	return &loginResp{UserId: 21, Token: "test grpc call success"}, nil
}

// 参数数量不正常、返回值正常
func Login1() (*loginResp, error) {
	fmt.Println("登录成功")
	return &loginResp{UserId: 21, Token: "test grpc call success"}, nil
}

// 第一个参数 不是 context 类型
func Login2(ctx int, req *loginReq) (*loginResp, error) {

	if req.Username != "ciusyan" && req.Password != "222" {
		return nil, fmt.Errorf("登录失败")
	}
	fmt.Println("登录成功")
	return &loginResp{UserId: 21, Token: "test grpc call success"}, nil
}

// 参数正常，返回值数量不正常
func Login3(ctx context.Context, req *loginReq) {
	if req.Username != "ciusyan" && req.Password != "222" {
		return
	}
	fmt.Println("登录成功")
}

// 参数正常，第二个返回值不是 error
func Login4(ctx context.Context, req *loginReq) (*loginResp, int) {
	if req.Username != "ciusyan" && req.Password != "222" {
		return nil, 0
	}
	fmt.Println("登录成功")
	return &loginResp{UserId: 21, Token: "test grpc call success"}, 1
}

var (
	datas []*data
)

func TestGrpcTask(t *testing.T) {

	for i, d := range datas {
		call := grpcrun.NewGrpcTask(&d.ctx, "test{"+strconv.Itoa(i)+"}", d.method, d.req)
		call.Call()

		t.Logf("第 %d 次执行\n", i+1)
		if call.Err != nil {
			fmt.Println(call.Err)
			fmt.Println()
			continue
		}
		// if should.NoError(call.Err) {
		//	fmt.Println(call.Response.(*loginResp))
		// }
		fmt.Println(call.Response.(*loginResp))
		fmt.Println()

	}
}

type data struct {
	ctx    context.Context
	method any
	req    any
}

func newData(ctx context.Context, method any, req any) *data {
	return &data{ctx: ctx, method: method, req: req}
}

func init() {
	req := &loginReq{Username: "ciusyan", Password: "222"}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 测试表格
	datas = []*data{
		newData(ctx, Login, req),        // 正常
		newData(ctx, Login1, req),       // [grpcMethod]必须有2个参数(context.Context, *request)
		newData(ctx, Login2, req),       // [grpcMethod]的第1个参数必须是：context.Context
		newData(ctx, Login3, req),       // [grpcMethod]必须有2个返回值(*Response, error)
		newData(ctx, Login4, req),       // [grpcMethod]的第2个返回值必须是：error
		newData(nil, Login, req),        // 请正确的传递[Context]，不支持：nil
		newData(ctx, nil, req),          // [grpcMethod]必须是一个GRPC的函数类型，现在是：invalid
		newData(ctx, Login, nil),        // 请正确的传递[request]，不支持：invalid
		newData(ctx, "其他类型", req),   // [grpcMethod]必须是一个GRPC的函数类型，现在是：string
		newData(ctx, Login, "其他类型"), // 请正确的传入[request]，不支持：string
		newData(ctx, Login, zap.S()),    // [request]的参数与[grpcMethod]的参数不匹配：grpcMethod = v3_test.loginReq, request = zap.SugaredLogger

	}
}

func TestGoGrpc_AddNewTask(t *testing.T) {
	run := grpcrun.NewGoGrpc()
	for i, d := range datas {
		run.AddNewTask("test{"+strconv.Itoa(i)+"}", d.method, d.req)
	}
}

func TestGoGrpc_Run(t *testing.T) {
	run := grpcrun.NewGoGrpc()
	for i, d := range datas {
		run.AddNewTask("test{"+strconv.Itoa(i)+"}", d.method, d.req)
	}
	run.Run()
	run.Wait()

	for _, t := range run.Task {
		if t.Err != nil {
			fmt.Println(t.Err)
			fmt.Println()
			continue
		}
		fmt.Println(t.Response.(*loginResp))
		fmt.Println()
	}
}
