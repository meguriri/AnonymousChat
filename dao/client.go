package dao

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
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
	HeartBeatTime  time.Time       //心跳时间
	HeartBeat      chan struct{}   //心跳
}

func (c *Client) OffLine() {
	close(c.MessageChan)
	close(c.LogoutSignal)
	close(c.UserListSignal)
	c.LConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	c.MConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	log.Printf("[client.Offine] %s has offlined\n", c.Id)
}

func (c *Client) ClientHandler() {
	//go c.HeartBeatCheck()
	for {
		select {
		case <-c.UserListSignal: //用户列表更新
			log.Printf("[client.ClientHandler] %s's userList has changed\n", c.Id)
			err := c.LConn.WriteJSON(MyManager.GetUserList())
			if err != nil {
				log.Printf("[client.ClientHandler error] %s writeJSON error: %v\n", c.Id, err.Error())
			}

		case <-c.HeartBeat: //心跳
			log.Printf("[client.ClientHandler] %s beated %v\n", c.Id, time.Now())
			c.HeartBeatTime = time.Now().Add(time.Duration(1e9 * 20))

		case msg := <-c.MessageChan: //读取消息
			log.Printf("[client.ClientHandler] %s get message %v\n", c.Id, msg)
			c.MConn.WriteJSON(msg)

		case <-c.LogoutSignal: //下线
			c.OffLine()
			return

		case <-time.After(time.Until(c.HeartBeatTime)): //心跳
			log.Printf("[client.ClientHandler] %s is death %v\n", c.Id, time.Now())
			//删除客户端
			MyManager.DeleteClient(c)
			//redis中删除session
			if err := Del(c.Id); err != nil {
				log.Printf("[client.ClientHandler error] %s delete session error: %v\n", c.Id, err.Error())
			}
			//下线
			c.OffLine()
			return
		}
	}
}
