package dao

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type Manager struct {
	Group                map[string]*Client //客户端列表
	Lock                 sync.Mutex         //互斥锁
	Register, UnRegister chan *Client       //上线，下线通道
	BroadCastChan        chan Message       //广播通道
	ClientCount          uint               //客户端数量
}

func (m *Manager) InitManager() {
	m.Register = make(chan *Client)
	m.UnRegister = make(chan *Client)
	m.BroadCastChan = make(chan Message, 10000)
	m.Group = make(map[string]*Client, 1000)
}

func (m *Manager) GetUserNumber() uint {
	return m.ClientCount
}

func (m *Manager) GetUserList() []User {
	userList := make([]User, 0, MaxUser)
	for _, v := range m.Group {
		userList = append(userList, *v.UserPtr)
	}
	return userList
}

func (m *Manager) ClientRegist(sid string, user *User) {
	//创建新的客户端
	client := &Client{
		Id:          sid,
		Conn:        new(websocket.Conn),
		UserPtr:     user,
		MessageChan: make(chan Message),
	}
	//注册
	m.Register <- client
}

func (m *Manager) ClientUnRegist(sid string) {
	m.UnRegister <- m.Group[sid]
}

func (m *Manager) Managed() {
	fmt.Println("start manager!!")
	for {
		select {
		case client := <-m.Register:
			//注册客户端
			m.Lock.Lock()
			m.Group[client.Id] = client
			m.ClientCount++
			fmt.Printf("客户端注册: 客户端id为%s,用户为%v\n", client.Id, *client.UserPtr)
			m.Lock.Unlock()

		case client := <-m.UnRegister:
			//注销客户端
			m.Lock.Lock()
			if _, ok := m.Group[client.Id]; ok {
				//关闭客户端消息通道
				close(client.MessageChan)
				//删除分组中客户端
				delete(m.Group, client.Id)
				//客户端数量减1
				m.ClientCount--
				fmt.Printf("客户端注销: 客户端id为%s", client.Id)
			}
			m.Lock.Unlock()

		case data := <-m.BroadCastChan:
			//将数据广播给所有客户端
			for _, conn := range m.Group {
				fmt.Println(conn.Id, data)
				conn.MessageChan <- data
			}

		}
	}
}
