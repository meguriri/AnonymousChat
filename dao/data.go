package dao

var (
	MaxUser   int
	MyManager Manager
)

type User struct {
	//用户昵称
	Nickname string `json:"nickname"`
	//用户性别
	Gender uint `json:"gender"`
	//头像颜色
	Color [3]uint `json:"color"`
}

type Message struct {
	//发送用户
	SendUser User
	//发送时间
	SendTime string `json:"sendtime"`
	//内容
	Content string `json:"content"`
}

//var BroadcastChan = make(chan Message, 10000)
