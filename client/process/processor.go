package process

import (
	"chatroom/common/message"
	"chatroom/common/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

//显示登录成功后的界面
func showMenu() {

	fmt.Println("---------------Congratulations your Login succeed-------------------")
	fmt.Println("---------------1. show the users online----------------")
	fmt.Println("---------------2. send message         ----------------")
	fmt.Println("---------------3. message list         ----------------")
	fmt.Println("---------------4. exit                 ----------------")
	fmt.Println("please choose a number among 1 to 4")
	var key int
	var content string
	//因为我们总会用到SmsProcess实例 因此我们将其定义在switch外部
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		fmt.Println("show the users online-")
		outputOnlineUser()
	case 2:
		fmt.Println("input a message")
		fmt.Scanln(&content)
		fmt.Println("********", content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("message list")
	case 4:
		fmt.Println("exit....")
		os.Exit(0)
	default:
		fmt.Println("invalid input..")
	}
}

//keep the connection between server and client
func clientProcessMes(conn net.Conn) {
	//创建一个transfer实例 不断读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("client waiting message from server...")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err=", err)
			return
		}
		//如果读取到消息 进行处理
		switch mes.Type {
		case message.NotifyUserStatusMesType:
			//表示有人上线
			//1. 取出NotifyUsersStatusMes
			var notifyUsersStatusMes message.NotifyUserStatusMes
			err = json.Unmarshal([]byte(mes.Data), &notifyUsersStatusMes)
			if err != nil {
				fmt.Println("json unmarshal err=", err)
			}
			//2. 把用户信息保存到客户map中
			updateUserStatus(&notifyUsersStatusMes)
		case message.SmsMesType:
			outputGroupMes(&mes)
		default:
			fmt.Println("server send back unknown message type")
		}
		//显示消息
		//fmt.Printf("mes=%v\n", mes)
	}
}
