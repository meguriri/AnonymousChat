package dao

type User struct {
	//用户昵称
	Nickname string `json:"nickname"`
	//用户性别
	Gender uint `json:"gender"`
	//头像颜色
	Color [3]uint `json:"color"`
}

var Users map[string]*User //

var UserList = make([]User, 0, len(Users))

type Message struct {
	//头像颜色
	//性别
	//昵称
	//发送时间
	//内容
}

var (
	MaxUser int
)
