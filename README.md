<!--
 * @Author: BeYoung
 * @Date: 2023-02-26 14:25:48
 * @LastEditTime: 2023-02-26 22:27:44
-->
# grpc goroutine
Running multiple grpc calls separately into goroutine.

If you want to run multiple grpc calls into goroutine, you might be like this:
```go
// goroutine 1
go func(c context.Context) {
    defer wait.Done()
    for {
	    select {
	    case <-ctx.Done():
		    return
	    default:
		    grpc.ServiceClient()
	    }
	}
}(ctx)
// goroutine 2
go func(c context.Context) {
    defer wait.Done()
    for {
	    select {
	    case <-ctx.Done():
		    return
	    default:
		    grpc.ServiceClient()
	    }
	}
}(ctx)
···
```
but now you can like this:
```go
run := GoGrpc{}
run.AddNewTask("grpcName", grpcMethod, grpcRequest)
run.Run()
run.Wait()
```
**Note: grpcName must is a unique value**

## use
Simple example:
```go
func example() {
    run := GoGrpc{}
    run.AddNewTask("grpcName", grpcMethod, grpcRequest)
    run.Run()
    run.Wait()
}
```
Or you could use NewGrpcTask creat a grpc task
```go
func example() {
    run := GoGrpc{}
    task := &NewGrpcTask{ctx, "grpcName", grpcMethod, grpcRequest}
    run.AddTask(task)
    run.Run()
    run.Wait()
}
```

