package process

import (
	"chatroom/common/message"
	"chatroom/common/utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
	//暂时不需要字段
}

func (us *UserProcess) Login(userId int, userPwd string) (err error) {
	//1.链接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("conn net.Dial err=", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2. 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType

	//3. 创建loginMes
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	//4. 将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json marshal err=", err)
	}

	//5. 把data赋给mes.Data字段
	mes.Data = string(data)

	//6. 将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal err=", err)
	}

	//7. 发送data消息
	//7.1 先发送data的长度
	//先获取到data的长度》转成一个表示长度的切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var bytes [4]byte
	binary.BigEndian.PutUint32(bytes[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(bytes[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	fmt.Println("client, message len correct")

	//7.2 发送data
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail=", err)
		return
	}

	//休眠20s
	//time.Sleep(20 * time.Second)
	//fmt.Println("sleep 20s..")
	//处理服务器返回的消息
	tf := utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err=", err)
		return
	}
	//将mes的Data部分反序列化成 LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if err != nil {
		fmt.Println("json unmarshal err=", err)
		return
	}
	if loginResMes.Code == 200 {
		//fmt.Println("client login success")

		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		//显示在线用户列表 遍历loginResMes.UserId
		fmt.Println("These are online users:")
		fmt.Println("############################")
		for _, v := range loginResMes.UserId {
			fmt.Printf("userId: %v \n", v)

			//客户端的 onlineUsers 初始化
			user := &message.User{
				UserId:     v,
				UserPwd:    "",
				UserName:   "",
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user

		}
		fmt.Println("############################")

		//这里我们需要在客户端启动一个协程
		//该协程保持和服务端的通讯 如果服务器有数据推送
		//则接收并显示在客户端
		go clientProcessMes(conn)

		//1.显示登录成功的菜单 循环
		for {
			showMenu()
		}

	} else {
		fmt.Println(loginResMes.Error)
	}

	return nil
}

func (us *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	//1.链接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("conn net.Dial err=", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2. 准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType

	//3. 创建registerMes
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	//4.将registerMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}

	//5. 把data赋给mes.Data字段
	mes.Data = string(data)

	//6. 将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}

	//7. 创建一个transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	//8. 发送data给服务器
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("register transfer data err=", err)
	}

	mes, err = tf.ReadPkg() //mes 是服务器返回信息
	if err != nil {
		fmt.Println("readPkg(conn) err=", err)
		return
	}

	//9. 将mes的Data部分反序列化成RegisterResMes
	var registerResMes message.RegisterResMes

	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("register succeed, please login")
		os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}
