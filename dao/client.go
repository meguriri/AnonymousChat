package dao

import (
	"github.com/gorilla/websocket"
)

// Client 客户端
type Client struct {
	Id          string          //客户端ID
	Conn        *websocket.Conn //websocket connect
	UserPtr     *User           //用户指针
	MessageChan chan Message    //消息通道
}
