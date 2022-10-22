package dao

import (
	"fmt"
	"github.com/gorilla/websocket"
)

// Client 客户端
type Client struct {
	Id             string          //客户端ID
	LConn          *websocket.Conn //user list websocket connect
	MConn          *websocket.Conn //message websocket connect
	UserPtr        *User           //用户指针
	MessageChan    chan Message    //消息通道
	UserListSignal chan struct{}   //用户列表更新信号
	LogoutSignal   chan struct{}   //用户下线信号
}

func (c *Client) ClientHandler() {
	for {
		select {
		case <-c.UserListSignal: //用户列表更新
			fmt.Println(c.Id, " user list changed")
			err := c.LConn.WriteJSON(MyManager.GetUserList())
			if err != nil {
				fmt.Println(err)
			}

		case msg := <-c.MessageChan: //读取消息
			fmt.Println(msg)
			c.MConn.WriteJSON(msg)

		case <-c.LogoutSignal: //下线
			close(c.MessageChan)
			close(c.LogoutSignal)
			close(c.UserListSignal)
			c.LConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			c.MConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		}
	}
}
