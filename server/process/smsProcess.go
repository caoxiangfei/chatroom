package process

import (
	"chatroom/common/message"
	"chatroom/common/utils"
	"encoding/json"
	"fmt"
	"net"
)

type SmsProcess struct {
}

//转发消息
func (sm *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器端的onlineUsers map[int]*UserProcess
	//将消息转发取出
	//取出mes的内容SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json unmarshal err=", err)
		return
	}
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal err=", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserId {
			continue
		}
		sm.SendMesToEachOnlineUser(data, up.Conn)
	}
}

func (sm *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {

	//创建一个Transfer
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("send message failed err=", err)
	}
}
