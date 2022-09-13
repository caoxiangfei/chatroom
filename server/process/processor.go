package process

import (
	"chatroom/common/message"
	"chatroom/common/utils"
	"fmt"
	"io"
	"net"
)

//先创建一个Processor 结构体
type Processor struct {
	Conn net.Conn
}

//编写一个ServerProcessMes
//功能：根据客户端发送消息种类不同，决定哪个函数来处理
func (pr *Processor) serverProcessMes(mes *message.Message) (err error) {

	switch mes.Type {
	case message.LoginMesType:
		//处理登录 TODO
		//创建一个UserProcess 实例
		us := &UserProcess{
			Conn: pr.Conn,
		}
		err = us.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册 TODO
		//创建一个UserProcess 实例
		us := &UserProcess{
			Conn: pr.Conn,
		}
		err = us.ServerProcessRegister(mes)
	case message.SmsMesType:
		//创建一个SmsProcess实例
		fmt.Println("######", mes.Data)
		smsProcess := &SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("message type not exist")
	}
	return
}

func (pr *Processor) Process2() (err error) {

	//循环读客户端发送的消息
	//创建一个transfer实例 完成读包任务
	tf := utils.Transfer{
		Conn: pr.Conn,
	}
	for {
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("client exited, server exit...")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}

		}

		fmt.Println("message=", mes)
		err = pr.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}

}
