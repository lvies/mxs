package znet

import (
	"net"
	"zinx/ziface"
	"zinx/log"
)

type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn
	// 当前连接的ID，可以称作为SessionID,ID全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool
	// 告知该连接已经退出/停止的channel
	ExitBuffChan chan bool
	// 该连接的处理方法router
	Router  ziface.IRouter
}

// 创建连接的方法
func NewConnecion(conn *net.TCPConn, connId uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connId,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
		Router: router,
	}
	return  c
}

// 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	log.Debug("Reader Goroutine is running")
	defer log.Debug("%s conn reader exit!",c.RemoteAddr().String())
	defer c.Stop()
	for {
		// 读取我们最大的数据到buf中
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil{
			log.Error("recv buf err %v", err)
			c.ExitBuffChan <- true
			continue
		}
		// 得到当前客户端请求的Request数据
		req := Request{
			conn:c,
			data:buf,
		}
		// 从路由Routers 中找到注册绑定Conn对应的Handle
		go func (request ziface.IRequest) {
			// 执行注册的路由方法
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	go c.StartReader()
	for {
		select {
		case <- c.ExitBuffChan:
		// 得到退出消息,不在阻塞
			return
		}
	}
}

// 获取远程客户端的地址信息
func (c *Connection) RemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}

func (c *Connection) Stop() {
	// 1.如果当前连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// TODO Connection Stop() 如果用户注册了该连接的关闭回调业务，那么此刻应该显示调用

	// 关闭socket连接
	c.Conn.Close()

	// 通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true

	// 关闭该连接的全部管道
	close(c.ExitBuffChan)
}

// 获取当前连接id
func (c *Connection) GetConnID() uint32{
	return c.ConnID
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}