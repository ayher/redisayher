package echo

import (
	"bufio"
	"context"
	"github.com/hdt3213/godis/lib/logger"
	"github.com/hdt3213/godis/lib/sync/atomic"
	"github.com/hdt3213/godis/lib/sync/wait"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// 客户端连接的抽象
type Client struct {
	// tcp 连接
	Conn net.Conn
	// 当服务端开始发送数据时进入waiting, 阻止其它goroutine关闭连接
	// wait.Wait是作者编写的带有最大等待时间的封装:
	// https://github.com/HDT3213/godis/blob/master/src/lib/sync/wait/wait.go
	Waiting wait.Wait
}

type EchoHandler struct {
	// 保存所有工作状态client的集合(把map当set用)
	// 需使用并发安全的容器
	activeConn sync.Map

	// 关闭状态标识位
	closing atomic.Boolean
}

func MakeEchoHandler()(*EchoHandler) {
	return &EchoHandler{}
}

func (h *EchoHandler)Handle(ctx context.Context, conn net.Conn) {
	// 关闭中的 handler 不会处理新连接
	if h.closing.Get() {
		conn.Close()
	}

	client := &Client{
		Conn: conn,
	}
	h.activeConn.Store(client, struct{}{}) // 记住仍然存活的连接

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("connection close")
				h.activeConn.Delete(conn)
			} else {
				logger.Warn(err)
			}
			return
		}
		// 发送数据前先置为waiting状态，阻止连接被关闭
		client.Waiting.Add(1)

		// 模拟关闭时未完成发送的情况
		//log.Println()("sleeping")
		//time.Sleep(10 * time.Second)

		b := []byte(msg)
		conn.Write(b)
		// 发送完毕, 结束waiting
		client.Waiting.Done()
	}
}

// 关闭客户端连接
func (c *Client)Close()error {
	// 等待数据发送完成或超时
	c.Waiting.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

// 关闭服务器
func (h *EchoHandler)Close()error {
	log.Println("handler shutting down...")
	h.closing.Set(true)
	// 逐个关闭连接
	h.activeConn.Range(func(key interface{}, val interface{})bool {
		client := key.(*Client)
		client.Close()
		return true
	})
	return nil
}