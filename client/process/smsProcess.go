package process

import (
	"chatroom/common/message"
	"chatroom/common/utils"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

//发送群聊的消息
func (sm *SmsProcess) SendGroupMes(content string) (err error) {
	//1 创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType

	//2. 创建一个SmsMes 实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus

	//3. 序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json marshal err=", err.Error())
		return
	}
	mes.Data = string(data)
	//4. 对mes再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json marshal err=", err.Error())
		return
	}

	//5. 将mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//6. 发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("sendGroupMes err=", err.Error())
		return
	}
	return
}
