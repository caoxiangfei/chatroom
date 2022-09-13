package model

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type UserDao struct {
	pool *redis.Pool
}

var MyUserDao *UserDao

//使用工厂模式 创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//1.根据一个用户id 返回一个User实例+err

func (ud *UserDao) getUserById(conn redis.Conn, id int) (user *message.User, err error) {
	//通过给定id 去redis查询该用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil { //表示在users 哈希中，没有找到对应id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	//把res 反序列化成User实例
	user = &message.User{
		UserId:   0,
		UserPwd:  "",
		UserName: "",
	}
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	//fmt.Println("#################################")
	//fmt.Println(user.UserId)
	return user, nil
}

//完成登录校验
//1. 完成用户校验
//2. 如果用户的id和密码正确 则返回user 实例
//3. 如果id或pwd有误，则返回对应错误信息
func (ud *UserDao) Login(userId int, userPwd string) (user *message.User, err error) {
	//先从UserDao连接池中取出一个连接
	conn := ud.pool.Get()
	defer conn.Close()
	user, err = ud.getUserById(conn, userId)
	if err != nil {
		return
	}

	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return user, nil
}

func (ud *UserDao) Register(user *message.User) (err error) {

	//先从UserDao的连接池中取出一个连接
	conn := ud.pool.Get()
	defer conn.Close()
	_, err = ud.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}
	//这时说明id在redis中不存在 可以注册
	data, err := json.Marshal(user) //序列化
	if err != nil {
		return
	}
	//入库
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("save user data err=", err)
		return
	}
	return
}