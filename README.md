# grpc goroutine

* 利用go的并发能力，将GRPC放置在后台运行。
* 统一封装，适用与所有的GRPC调用

## Usage

```go
func main()  {
	req := &loginReq{Username: "ciusyan", Password: "222"}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	call := grpcrun.NewGrpc(&ctx, Login, req)
	call.GrpcTask()

	if call.Err != nil {
		fmt.Println(call.Err)
		fmt.Println()
	}
	fmt.Println(call.Response.(*loginResp))
}
```