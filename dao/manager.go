package dao

import (
	"log"
	"sync"
	"time"

	"github.com/meguriri/AnonymousChat/redis"
)

var (
	MyManager Manager
)

type Manager struct {
	Group                map[string]*Client //客户端列表
	Lock                 sync.RWMutex       //读写锁
	Register, UnRegister chan *Client       //上线，下线通道
	BroadCastChan        chan Message       //广播通道
	ClientCount          uint               //客户端数量
}

func (m *Manager) InitManager() { //初始化管理器
	m.Register = make(chan *Client)
	m.UnRegister = make(chan *Client)
	m.BroadCastChan = make(chan Message, 10000)
	m.Group = make(map[string]*Client, 1024)
}

func (m *Manager) GetClient(sid string) *Client { //获取客户端
	return m.Group[sid]
}

func (m *Manager) GetUserNumber() uint { //获取用户人数
	return m.ClientCount
}

func (m *Manager) GetUserList() []User { //获取用户列表
	userList := make([]User, 0, MaxUser)
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	for _, v := range m.Group {
		userList = append(userList, *v.UserPtr)
	}
	return userList
}

func (m *Manager) ClientRegister(sid string, user *User) { //注册客户端
	//创建新的客户端
	client := &Client{
		Id:             sid,
		LConn:          nil,
		MConn:          nil,
		UserPtr:        user,
		MessageChan:    make(chan Message, 1024),
		UserListSignal: make(chan struct{}),
		LogoutSignal:   make(chan struct{}),
		HeartBeatTime:  time.Now().Add(time.Duration(MaxLifetime)),
		HeartBeat:      make(chan struct{}),
	}
	//注册
	m.Register <- client
}

func (m *Manager) ClientUnRegister(sid string) { //注销客户端
	m.UnRegister <- m.Group[sid]
}

func (m *Manager) DeleteClient(client *Client) { //删除客户端
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if _, ok := m.Group[client.Id]; ok {
		//删除分组中客户端
		delete(m.Group, client.Id)
		//客户端数量减1
		m.ClientCount--
		log.Printf("[manager.DeleteClient] clinet unregister: id: %s\n", client.Id)
	}
}

func (m *Manager) Managed() { //管理
	for {
		select {
		case client := <-m.Register:
			log.Printf("[manager.Managed] a new client  is coming: sid: %v\n", client.Id)
			//注册客户端
			m.Lock.Lock()
			m.Group[client.Id] = client
			m.ClientCount++
			log.Printf("[manager.Managed] client register: sid：%s, user:%v\n", client.Id, *client.UserPtr)
			m.Lock.Unlock()

			//开启客户端监听协程
			log.Printf("[manager.Managed] %v start a new ClientHandler\n", client.Id)
			go client.ClientHandler()

			//向除了刚注册的客户端以外的所有客户端通知更新用户列表
			for _, v := range m.Group {
				if v.Id != client.Id {
					v.UserListSignal <- struct{}{}
				}
			}

		case client := <-m.UnRegister:
			log.Printf("[manager.Managed] %v a client will leave\n", client.Id)

			//注销客户端
			m.DeleteClient(client)

			//向客户端handler发送注销信号
			client.LogoutSignal <- struct{}{}

			//向所有客户端通知更新用户列表
			for _, v := range m.Group {
				v.UserListSignal <- struct{}{}
			}

		case msg := <-m.BroadCastChan:
			//存入redis
			if err := msg.Save(redis.Rdb); err != nil {
				log.Printf("[manager.Managed]message save err: %v\n", err.Error())
			}
			//将数据广播给所有客户端
			for _, v := range m.Group {
				//将消息广播给每个客户端
				if v.UserPtr.Nickname != msg.SendUser.Nickname {
					log.Println(v.Id, msg)
					v.MessageChan <- msg
				}
			}
		}
	}
}
