package dao

import (
	"fmt"
	"github.com/gorilla/websocket"
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
	fmt.Println("xia xian")
}

func (c *Client) ClientHandler() {
	//go c.HeartBeatCheck()
	for {
		select {
		case <-c.UserListSignal: //用户列表更新
			fmt.Println(c.Id, " user list changed")
			err := c.LConn.WriteJSON(MyManager.GetUserList())
			if err != nil {
				fmt.Println(err)
			}

		case <-c.HeartBeat: //心跳
			fmt.Println(c.Id, " heartbeat ", time.Now())
			c.HeartBeatTime = time.Now().Add(time.Duration(1e9 * 20))

		case msg := <-c.MessageChan: //读取消息
			fmt.Println(msg)
			c.MConn.WriteJSON(msg)

		case <-c.LogoutSignal: //下线
			c.OffLine()
			return

		case <-time.After(time.Until(c.HeartBeatTime)): //心跳
			fmt.Println(c.Id, " is death ", time.Now())
			//删除客户端
			MyManager.DeleteClient(c)
			//redis中删除session
			if err := Del(c.Id); err != nil {
				fmt.Println(err)
			}
			//下线
			c.OffLine()
			return
		}
	}
}
