package utils

import (
	"chatroom/common/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

//这里将这些方法关联到结构体中
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte
}

func (tr *Transfer) ReadPkg() (mes message.Message, err error) {

	//buf := make([]byte, 8096)
	//fmt.Println("reading data ..")
	_, err = tr.Conn.Read(tr.Buf[:4])
	if err != nil {
		fmt.Println("conn.Read err=", err)
		return
	}
	//fmt.Println("read but=", buf[:4])
	//根据buf[:4]转换成unit32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(tr.Buf[:4])

	//根据pkgLen读取消息内容
	n, err := tr.Conn.Read(tr.Buf[:pkgLen])

	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Read fail err=", err)
		return
	}

	//把pkgLen 反序列化成 message.Message
	err = json.Unmarshal(tr.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json unmarshal err=", err)
		return
	}
	return
}

//编写一个writePkg
func (tr *Transfer) WritePkg(data []byte) (err error) {
	//先发送一个长度给对方
	//先获取到data的长度》转成一个表示长度的切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	//var bytes [4]byte
	binary.BigEndian.PutUint32(tr.Buf[0:4], pkgLen)
	//发送长度
	n, err := tr.Conn.Write(tr.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail err=", err)
		return
	}
	//fmt.Println("message len correct")
	//发送数据本身
	n, err = tr.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) fail err=", err)
		return
	}
	return
}
