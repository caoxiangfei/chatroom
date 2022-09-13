package process

import (
	"chatroom/client/model"
	"chatroom/common/message"
	"fmt"
)

//客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser //在用户登录成功后 完成对CurUser初始化

//在客户端显示当前的用户
func outputOnlineUser() {
	//遍历一遍 onlineUsers
	fmt.Println("----Users online----")
	for id, _ := range onlineUsers {
		fmt.Printf("UserId: %v\n", id)
	}
	fmt.Println("--------------------")
}

//编写一个方法 处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {

	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		user = &message.User{UserId: notifyUserStatusMes.UserId}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user

	outputOnlineUser()
}
