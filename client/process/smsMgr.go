package process

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
)

func outputGroupMes(mes *message.Message) {
	//显示即可
	//1. 发序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json unmarshal err=", err.Error())
		return
	}

	//显示信息
	info := fmt.Sprintf("userId: %d send message: %s",
		smsMes.UserId, smsMes.Content)
	fmt.Println("####################################")
	fmt.Println(info)
	fmt.Println("####################################")
	fmt.Println()
}
