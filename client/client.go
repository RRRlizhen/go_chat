package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	conn           *net.UDPConn
	msg_type       int
	gkey           bool
	userId         int
	userName       string
	sendMessage    chan string
	receiveMessage chan string
}

func (c *Client) sendMsg() {
	for c.gkey {
		msg := <-c.sendMessage
		//处理向server端 发送msg的动作
		str := fmt.Sprintf("###2##%d##%s##%s###", c.userId, c.userName, msg)
		_, err := c.conn.Write([]byte(str))
		checkError(err, "sendMsg")
	}
}

func (c *Client) receiveMsg() {
	var buf [512]byte
	for c.gkey {
		//处理从server端 接受的动作
		n, err := c.conn.Read(buf[0:])
		checkError(err, "receiveMsg")
		c.receiveMessage <- string(buf[0:n])
	}
}

func (c *Client) inputMsg() {
	//这里处理 用户从stdin 输入的msg
	//处理断开的动作
	var msg string
	for c.gkey {
		fmt.Println("msg: ")
		_, err := fmt.Scanln(&msg)
		checkError(err, "inputMsg")

		if msg == "quit" {
			c.gkey = false
			return
		} else {
			c.sendMessage <- encodeMessage(msg)
		}
	}
}

func encodeMessage(msg string) string {
	return strings.Join(strings.Split(strings.Join(strings.Split(msg, "\\"), "\\\\"), "#"), "\\#")
}

func (c *Client) printMsg() {
	//将 从server接收到的msg 打印到stdout
	for c.gkey {
		msg := <-c.receiveMessage
		fmt.Println(msg)
	}
}

func main() {
	service := "localhost:1200"

	udpAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err, "main")

	var c Client
	c.gkey = true
	c.sendMessage = make(chan string)
	c.receiveMessage = make(chan string)

	fmt.Println("input id:")
	_, err = fmt.Scanln(&c.userId)
	checkError(err, "main")

	fmt.Println("input username:")
	_, err = fmt.Scanln(&c.userName)
	checkError(err, "main")

	//connect the udp server
	c.conn, err = net.DialUDP("udp4", nil, udpAddr)
	checkError(err, "main")
	defer c.conn.Close()

	//开始发送信息，
	//msg_type = 1, 进入chatroom
	//msg_type = 2，发送msg
	//msg_type = 3，离开chatroom
	c.msg_type = 1
	str := fmt.Sprintf("###%d##%d##%s##%s###", c.msg_type, c.userId, c.userName, "进入chatroom")
	_, err = c.conn.Write([]byte(str))
	checkError(err, "main")

	go c.printMsg()   //  <-c.receiv
	go c.receiveMsg() // c.receiv<-

	go c.sendMsg() // <-c.send
	c.inputMsg()   // c.send<-

	str_exit := fmt.Sprintf("###%d##%d##%s##%s###", 3, c.userId, c.userName, "leave the chatroom")

	_, err = c.conn.Write([]byte(str_exit))
	checkError(err, "main")
	os.Exit(0)
}

func checkError(err error, funcName string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s----- in func: %s", err.Error(), funcName)
		os.Exit(1)
	}
}
