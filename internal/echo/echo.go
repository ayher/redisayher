package echo

import (
	"context"
	"fmt"
	"github.com/hdt3213/godis/interface/tcp"
	"github.com/hdt3213/godis/lib/logger"
	roottcp "github.com/hdt3213/godis/tcp"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)


// HandleFunc represents application handler function
type HandleFunc func(ctx context.Context, conn net.Conn)

// 监听并提供服务，并在收到 closeChan 发来的关闭通知后关闭
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	// 监听关闭通知
	go func() {
		<-closeChan
		logger.Info("shutting down...")
		// 停止监听，listener.Accept()会立即返回 io.EOF
		_ = listener.Close()
		// 关闭应用层服务器
		_ = handler.Close()
	}()

	// 在异常退出后释放资源
	defer func() {
		// close during unexpected error
		_ = listener.Close()
		_ = handler.Close()
	}()
	ctx := context.Background()
	var waitDone sync.WaitGroup
	for {
		// 监听端口, 阻塞直到收到新连接或者出现错误
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		// 开启 goroutine 来处理新连接
		logger.Info("accept link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}

// ListenAndServeWithSignal 监听中断信号并通过 closeChan 通知服务器关闭
func ListenAndServeWithSignal(cfg *roottcp.Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listener, handler, closeChan)
	return nil
}
