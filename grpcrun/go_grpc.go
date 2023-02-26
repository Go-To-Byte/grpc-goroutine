// Author: BeYoung
// Date: 2023/2/26 14:27
// Software: GoLand

package grpcrun

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	mu   sync.Mutex
	node *snowflake.Node
)

func init() {
	var err error
	mu = sync.Mutex{}
	if node, err = snowflake.NewNode(int64(time.Now().Day())); err != nil {
		panic(err)
	}
}

// GoGrpc is used to run some goroutine of grpc.
// the grpc's return will filled in response and error
// Example:
//
//	func example() {
//		run := GoGrpc{}
//		run.AddNewTask(nil, nil, nil)
//		run.Call()
//		run.Wait()
//	}
type GoGrpc struct {
	mu      sync.Mutex
	ctx     context.Context
	cancel  context.CancelFunc
	wait    sync.WaitGroup
	Timeout time.Duration
	Task    map[string]*GrpcTask
}

func NewGoGrpc() *GoGrpc {
	mu.Lock()
	defer mu.Unlock()
	g := GoGrpc{}
	g.ctx, g.cancel = context.WithTimeout(context.Background(), 3*time.Second)
	g.mu = sync.Mutex{}
	g.wait = sync.WaitGroup{}
	g.Task = make(map[string]*GrpcTask, 0)
	return &g
}

// SetTimeout reset timeout, replace default timeout with a special time duration
func (g *GoGrpc) SetTimeout(timeout time.Duration) {
	mu.Lock()
	mu.Unlock()
	g.ctx, g.cancel = context.WithTimeout(context.Background(), timeout)
}

func (g *GoGrpc) Run() {
	for _, t := range g.Task {
		go g.run(t)
	}
}

func (g *GoGrpc) Wait() {
	defer g.cancel()
	g.wait.Wait()
}

func (g *GoGrpc) AddTask(task *GrpcTask) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Task[task.Name] = task
	g.wait.Add(1)
}

func (g *GoGrpc) AddNewTask(grpcName string, grpcMethod any, request any) {
	g.mu.Lock()
	defer g.mu.Unlock()
	zap.S()
	task := GrpcTask{
		ctx:        &g.ctx,
		grpcMethod: grpcMethod,
		request:    request,
		Name:       grpcName,
		log:        zap.S().Named(grpcName),
	}

	g.Task[node.Generate().String()] = &task
	g.wait.Add(1)
	return
}

func (g *GoGrpc) run(t *GrpcTask) {
	defer g.wait.Done()
	for {
		select {
		case <-g.ctx.Done():
			t.Err = errors.New("context canceled")
			return
		default:
			t.Call()
			return
		}
	}
}

func example() {
	run := GoGrpc{}
	run.AddNewTask("nil", nil, nil)
	run.Run()
	run.Wait()
}
