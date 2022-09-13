package process

import (
	"chatroom/common/message"
	"chatroom/common/utils"
	"chatroom/server/model"
	"encoding/json"
	"fmt"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加一个字段表示该连接来自哪个用户
	UserId int
}

//编写一个serverProcessLogin 函数
func (us *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	//核心代码
	//1. 先从mes 中取出 mes.Data, 并直接反序列化成LoginMes
	loginMes := message.LoginMes{
		UserId:   0,
		UserPwd:  "",
		UserName: "",
	}
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//fmt.Println("######################################")
	//fmt.Println(loginMes)

	//先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//声明一个LoginResMes
	var loginResMes message.LoginResMes

	//去redis数据库完成验证
	//1. 使用model.MyUserDao 到redis去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	fmt.Printf("user %v login success..\n", user)
	if err != nil {
		loginResMes.Code = 500
		loginResMes.Error = "client not exist, please registry"

	} else {
		loginResMes.Code = 200
		//将登录成功用户的userId 赋给 us
		us.UserId = loginMes.UserId
		userMgr.AddOnlineUser(us)

		//通知其他在线用户 我上线了
		us.NotifyOthersOnlineUser(loginMes.UserId)

		//将当前在线用户的id 放入到loginResMes.UserId
		//遍历userMgr.onlineUsers
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UserId = append(loginResMes.UserId, id)
		}
	}

	//如果用户的id=100 password=123456 合法
	//if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
	//	//合法
	//	loginResMes.Code = 200
	//} else {
	//	//不合法
	//	loginResMes.Code = 500
	//	loginResMes.Error = "client not exist, please registry"
	//	//return
	//}

	//3 将loginResMes 序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("login.Marshal fail err=", err)
		return
	}
	//4 将data赋值给resMes
	resMes.Data = string(data)
	//5 对resMes序列化
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("login.Marshal fail err=", err)
		return
	}
	//6. 发送data 编写一个write函数 TODO
	//因为使用的是分层模式 先创建一个transfer实列
	tf := &utils.Transfer{
		Conn: us.Conn,
	}
	err = tf.WritePkg(data)
	return
}

func (us *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//1. 先从mes 中取出 mes.Data, 并直接反序列化成RegisterMes
	registerMes := message.RegisterMes{User: message.User{
		UserId:   0,
		UserPwd:  "",
		UserName: "",
	}}
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}

	//先声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterMesType
	//声明一个LoginResMes
	var registerResMes message.RegisterResMes

	//我们需要到redis数据库完成注册
	//1. 使用model.MyUserDao 到redis去验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 506
			registerResMes.Error = "unknown err occurred"
		}
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	//4. 将data赋值给resMes
	resMes.Data = string(data)

	//5. 对resMes进行序列化 准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}

	//6. 发送data
	tf := &utils.Transfer{
		Conn: us.Conn,
	}

	err = tf.WritePkg(data)
	return
}

//编写一个通知所有在线用户的方法
func (us *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历 onlineUsers， 然后一个一个的发送
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		//开始通知 单独写一个方法
		up.NotifyMeOnline(userId)
	}

}

//notify
func (us *UserProcess) NotifyMeOnline(userId int) {
	//1. 组装notify mes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//2. 将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}
	//将序列化后的notifyUserStatusMes的值赋给 mes.Data
	mes.Data = string(data)

	//对mes序列化 准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}

	//3.准备发送mes
	//创建Transfer实例
	tf := &utils.Transfer{
		Conn: us.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}
