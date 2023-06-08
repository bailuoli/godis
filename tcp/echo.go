package tcp

import (
	"bufio"
	"context"
	"godis/lib/logger"
	"godis/lib/sync/atomic"
	"godis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

// 客户端服务
type EchoHandler struct {
	activeConn sync.Map       //客户端连接
	closing    atomic.Boolean //是否已经关闭
}

// 客户端实体
type EchoClient struct {
	Conn net.Conn  //客户端连接
	wait wait.Wait //waitGroup
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

// 关闭客户端
func (e *EchoClient) Close() error {
	//等待客户端任务结束
	e.wait.WaitWithTimeout(7 * time.Second)
	//关闭客户端连接
	e.Conn.Close()
	return nil
}

func (e *EchoHandler) Handle(c context.Context, conn net.Conn) {
	//如果closing是已经关闭的状态 就关闭连接
	if e.closing.Get() {
		_ = conn.Close()
	}

	//如果未关闭 则实现下面业务
	//初始化客户端
	client := &EchoClient{
		Conn: conn,
	}

	//将客户端信息存入activeConn中
	e.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)

	//接受客户端消息
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("Conn close")
				e.activeConn.Delete(conn)
			} else {
				logger.Warn(err)
			}
			return
		}

		//发送消息
		client.wait.Add(1)

		msgByte := []byte(msg)
		//发送数据
		_, _ = client.Conn.Write(msgByte)
		client.wait.Done()
	}
}

func (e *EchoHandler) Close() error {
	logger.Info("shutting down")
	e.closing.Set(true)
	e.activeConn.Range(func(key, value any) bool {
		//将之前存入到activeConn中的客户端数据取出来
		client := key.(*EchoClient)
		//关闭客户端
		client.Conn.Close()
		//返回true遍历继续
		return true
	})
	return nil
}
