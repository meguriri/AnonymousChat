package dao

var (
	MaxUser int
)

type User struct {
	//用户昵称
	Nickname string `json:"nickname"`
	//用户性别
	Gender uint `json:"gender"`
	//头像颜色
	Color [3]uint `json:"color"`
}
