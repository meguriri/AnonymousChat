package logic

import (
	"github.com/meguriri/AnonymousChat/dao"
	"net/http"
)

func CheckOrigin(r *http.Request) bool { //判断跨域
	return true
}

func CheckDup(nickname string) bool { //判断重名
	userlist := dao.MyManager.GetUserList()
	for _, v := range userlist {
		if v.Nickname == nickname {
			return false
		}
	}
	return true
}
