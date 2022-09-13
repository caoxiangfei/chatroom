package main

import (
	"chatroom/client/process"
	"fmt"
)
import "os"

//定义两个变量，一个表示用户id， 一个表示用户密码
var userId int
var userPwd string
var userName string

func main() {

	//接收用户的选择
	var key int
	//判断是否还继续显示菜单
	var loop = true

	for loop {
		fmt.Println("------------------welcome to our chatroom----------------------")
		fmt.Println("                     1 login chatroom")
		fmt.Println("                     2 registry")
		fmt.Println("                     3 exit")
		fmt.Println("                     4 please choose a number(1-3)")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("login")
			fmt.Println("please input userId")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("please input your password")
			fmt.Scanf("%s\n", &userPwd)
			//loop = false
			//完成登录
			//1. 创建一个UserProcess的实例
			us := &process.UserProcess{}
			us.Login(userId, userPwd)
		case 2:
			fmt.Println("registry")
			fmt.Println("please input userId:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("please input password:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("please input username:")
			fmt.Scanf("%s\n", &userName)

			//2. 调用UserProcess 完成注册
			us := &process.UserProcess{}
			us.Register(userId, userPwd, userName)
			//loop = false
		case 3:
			fmt.Println("exit")
			os.Exit(0)
		default:
			fmt.Println("invalid input, please try again")
		}

		//根据用户的输入 显示新的提示信息
		//if key == 1 {
		//	//用户登录
		//	fmt.Println("please input userId")
		//	fmt.Scanf("%d\n", &userId)
		//	fmt.Println("please input your password")
		//	fmt.Scanf("%s\n", &userPwd)
		//	//登录函数写到另一个文件 login.go
		//	//dao.Login(userId, userPwd)
		//	//if err != nil {
		//	//	fmt.Println("failed to login")
		//	//} else {
		//	//	fmt.Println("login success")
		//	//}
		//} else if key == 2 {
		//
		//}
	}
}
