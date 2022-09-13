package model

import (
	"chatroom/common/message"
	"net"
)

//因为在客户端很多地方都会用到CurUser将其作为全局的

type CurUser struct {
	Conn net.Conn
	message.User
}
