package iface

import (
	"github.com/xtaci/kcp-go"
	"net"
)


// 此处用于tcp链接
type IConnection interface {
	// 启动连接，当前连接开始工作
	Start()
	// 停止连接，结束当前连接状态M
	Stop()
	// 从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接id
	GetConnID() uint32
	// 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// 直接将Message数据发送给远程的TCP客户端(无缓冲)
	SendMsg(MsgId uint32, data []byte) error
	// 直接将Message数据发送给远程的TCP客户端（有缓冲）
	SendBuffMsg(MsgId uint32, data []byte) error	// 添加带缓冲发送消息接口
	// 设置连接属性
	SetProperty(key string, value interface{})
	// 获取链接属性
	GetProperty(key string) (interface{}, error)
	// 移除连接属性
	RemoveProperty(key string)
}

// 定义一个统一处理链接业务的接口
//type HandFunc func(conn *mnet.TCPConn,bytes []byte,event int) error

// 此处用于kcp链接
type IKConnection interface {
	// 启动kcp链接
	Start()
	// 停止kcp链接
	Stop()
	// 获得udp会话信息
	GetUdpSession() *kcp.UDPSession
	// 获得kcp链接id
	GetKcpId() uint32
	// 获得链接地址
	GetConnectionAddr() net.Addr
	// 发送消息(无缓冲)
	SendMsg(msg uint32, data []byte)  (int, error)
}