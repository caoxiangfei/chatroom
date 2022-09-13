package main

import (
	"chatroom/server/model"
	process2 "chatroom/server/process"
	"time"

	//"chatroom/server/process"
	"fmt"
	"net"
)

//处理和客户端的通讯
func process(conn net.Conn) {
	//这里需要延时关闭conn
	defer conn.Close()

	//这里调用总控 创建一个
	processor := &process2.Processor{
		Conn: conn,
	}
	err := processor.Process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误 err=", err)
		return
	}
}

func initUserDao() {
	//pool 本身就是一个全局变量
	//initUserDao在initPool之后初始化
	model.MyUserDao = model.NewUserDao(model.Pool)
}

func main() {

	//服务器启动时初始化redis连接池
	model.InitPool("47.112.112.32:6379", 16, 3, 100*time.Second)
	//初始化UserDao
	initUserDao()

	//提示信息
	fmt.Println("服务器(new)在8889端口监听...")
	listen, err := net.Listen("tcp", "127.0.0.1:8889")
	defer listen.Close()
	if err != nil {
		fmt.Println("net.listen err=", err)
		return
	}
	//一旦监听成功，就等待客户端来连接服务器
	for {
		fmt.Println("wait for client...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		//一旦连接成功， 则启动一个协程和客户端保持通讯...
		go process(conn)
	}
}
