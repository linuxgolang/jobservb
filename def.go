package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"gopkg.in/snksoft/crc.v1"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	HEADER = "250"//数据包开始标志
	PAYLOAD_SIZE = 1//数据长度占用字节,这里其实只保存20这个值,所以1个字节就够了
)

var (
	ErrLoginData = errors.New("错误的登陆数据")
)

func readSomething(conn *net.Conn) []byte {
	buffer := make([]byte, 512)
	for{
		_, err := (*conn).Read(buffer)
		if isErrAPrint(err) {
			//fmt.Println("aaaaaa")
			return nil
		}
		//if err != nil{
		//	//连接出现问题或服务器关闭连接,退出.
		//	panic(fmt.Sprintf("Read error: %s", err))
		//}
		//这里处理服务器返回的服务数据
		fmt.Println(buffer)
	}
}

/**
 * 判断登陆数据格式是否正确、完整
 */
func checkLoginData(data []byte)(bool,error){
	total := len(data)//总长(26)
	headerLen := len(HEADER)//固定头部占用字节数(3)
	headerPayloadLen := headerLen+PAYLOAD_SIZE//固定头部占用字节数和payload长度数字占用字节(4)
	if len(data) <= headerPayloadLen {
		return false,nil
	}
	payloadAndCrc := data[headerPayloadLen:]//payload和crc所有数据(22)
	payloadLen := data[headerLen:headerPayloadLen][0]//payload数据应该占用的字节数(20)
	if string(data[:headerLen]) != HEADER {
		return false,ErrLoginData
	}
	if (uint64(len(payloadAndCrc)) - uint64(payloadLen)) != 2{
		return false,nil
	}
	payload := data[headerPayloadLen:total-2]
	rck := data[total-2:]

	ck := make([]byte, 2)
	hash := crc.NewHash(crc.X25)
	x25Crc := hash.CalculateCRC(payload)
	binary.LittleEndian.PutUint16(ck, uint16(x25Crc))
	if ck[0] == rck[0] && ck[1] == rck[1] {
		return true,nil
	}

	return false,ErrLoginData
}

func writeSomething(conn *net.Conn, data []byte) bool {
	err := (*conn).SetWriteDeadline(time.Now().Add(5*time.Second))
	isErrAPrint(err)
	_, err = (*conn).Write(data)
	if isErrAPrint(err) {
		//fmt.Println("bbbbbbb")
		return false
		//连接出现问题,退出.
		//panic(fmt.Sprintf("Write error: %s", err))
	}
	return true
}

func isErrAPrint(err error) bool {
	if err != nil{
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return true
	}
	return false
}

func watch() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs,syscall.SIGINT)
	select {
	case sig := <- sigs:
		if sig == syscall.SIGINT{
			os.Exit(0)
		}
	}
}
