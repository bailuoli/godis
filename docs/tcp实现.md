第一天、tcp实现

1、开启一个连接

```go
listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
}
```

2、使用一个死for循环 来连接客户端                                           

```go
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

```

