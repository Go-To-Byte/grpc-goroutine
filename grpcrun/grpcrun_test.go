// @Author: Ciusyan 2023/2/26
package grpcrun_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	if req.Username != "ciusyan" && req.Password != "222" {
		return nil, fmt.Errorf("登录失败")
	}
	fmt.Println("登录成功")
	return &loginResp{UserId: 21, Token: "testgrpccallsuccess"}, nil
}

func TestGrpcTask(t *testing.T) {

	should := assert.New(t)
	req := &loginReq{Username: "ciusyan", Password: "222"}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := grpcrun.GrpcTask(Login, &ctx, req)

	if should.NoError(err) {
		t.Log(resp.(*loginResp))
	}
}
