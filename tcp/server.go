package tcp

import (
	"context"
	"fmt"
	"godis/interface/tcp"
	"godis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Config struct {
	//监听地址
	Address string
}

// ClientCount 定义连接客户端数量
var ClientCount int

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	//创建关闭tcp连接信号量chan
	closeChan := make(chan struct{})

	//信号量chan
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sign := <-signalChan
		switch sign {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listen, handler, closeChan)
	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {

	// 定义错误 chan
	errChan := make(chan error, 1)
	defer close(errChan)

	go func() {
		//接受调度段的关闭信号
		select {
		//如果时关闭信号量
		case <-closeChan:
			logger.Info("get exit signal")
		//如果是错误信号量
		case err := <-errChan:
			logger.Error(fmt.Sprintf("accept error %s", err.Error()))
		}

		logger.Info("shutting down")
		//关闭连接
		_ = listener.Close()
		_ = handler.Close()
	}()

	var waitDone = sync.WaitGroup{}
	ctx := context.Background()

	//客户端连接
	for true {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accept link")
		ClientCount++
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
				ClientCount--
			}()
			handler.Handle(ctx, conn)
		}()
	}

	//等待所有go func 结束
	waitDone.Wait()

}
