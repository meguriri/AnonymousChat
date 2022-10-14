package dao

import (
	"github.com/gorilla/websocket"
	"sync"
)

//客户端
type Client struct {
	Id          string          //客户端ID
	Socket      *websocket.Conn //websocket Conn
	MessageChan chan Message    //消息管道
}

type Manager struct {
	Group                map[string]*Client
	Lock                 sync.Mutex
	Register, UnRegister chan *Client
	BroadCastChan        chan Message
	ClientCount          uint
}

func (m *Manager) InitManager() {
	m.Register = make(chan *Client)
	m.UnRegister = make(chan *Client)
	m.BroadCastChan = make(chan Message, 1000)
	m.Group = make(map[string]*Client, 1000)
}
