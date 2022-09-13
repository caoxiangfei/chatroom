package process

import "fmt"

//在很多地方都会用到 将其定义为全局变量
var userMgr *UserMgr

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

//初始化userMgr
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

//onlineUsers添加
func (um *UserMgr) AddOnlineUser(up *UserProcess) {
	um.onlineUsers[up.UserId] = up
}

//onlineUsers删除
func (um *UserMgr) DelOnlineUser(userId int) {
	delete(um.onlineUsers, userId)
}

//返回当前所有用户
func (um *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return um.onlineUsers
}

//根据id返回对应的值
func (um *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	//如何从map中取出一个值 待检测方式
	up, ok := um.onlineUsers[userId]
	if !ok {
		err = fmt.Errorf("user %d not exists..\n", userId)
		return
	}
	return
}
