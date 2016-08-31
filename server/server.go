package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	conn    *net.UDPConn
	message chan string
	clients map[int]Client
}

type Client struct {
	userId   int
	userName string
	userAddr *net.UDPAddr
}

type Message struct {
	status   int
	userId   int
	userName string
	content  string
}

func (s *Server) printMsg() {
	msg := <-s.message
	fmt.Println(msg)
}

func (s *Server) sendMsg() {
	for {
		daytime := time.Now().String()
		time.Sleep(1 * time.Second)
		for _, c := range s.clients {
			fmt.Println(c.userName)
			n, err := s.conn.WriteToUDP([]byte(daytime), c.userAddr)
			checkError(err, "sendMsg")
			fmt.Printf("send %d bytes to %s\n", n, c.userName)
		}
	}
}

func (s *Server) decoder(msg string, m *Message) {
	s1 := strings.Split(msg, "###")
	s2 := strings.Split(s1[1], "##") //s1[1]= 1##123##abc##enter chatroom
	switch s2[0] {
	case "1":
		m.status, _ = strconv.Atoi(s2[0])
		m.userId, _ = strconv.Atoi(s2[1])
		m.userName = s2[2]
		m.content = s2[3]
		return
	case "2":
		m.status, _ = strconv.Atoi(s2[0])
		m.userId, _ = strconv.Atoi(s2[1])
		m.userName = s2[2]
		m.content = s2[3]
		return
	case "3":
		m.status, _ = strconv.Atoi(s2[0])
		m.userId, _ = strconv.Atoi(s2[1])
		m.userName = s2[2]
		m.content = s2[3]
		return
	default:
		fmt.Println("unkown msg=====", msg)
		return
	}
	return
}

func (s *Server) handle_msg() {
	for {
		var buf [512]byte

		n, addr, err := s.conn.ReadFromUDP(buf[0:])
		checkError(err, "handle_msg")
		fmt.Println("ok_1")
		//decoder msg
		msg := string(buf[0:n])
		m := new(Message)
		s.decoder(msg, m)

		switch m.status {
		case 1:
			var c Client
			c.userAddr = addr
			c.userId = m.userId
			c.userName = m.userName
			s.clients[c.userId] = c
			fmt.Println("a new user to chatroom")
			s.message <- m.content
		case 2:
			fmt.Printf("this is a msg\n")
			fmt.Printf("recevie %d bytes msg from %s", n, addr)
			s.message <- m.content
		case 3:
			delete(s.clients, m.userId)
			fmt.Printf("the user: %s has leave from the chatroom", m.userName)
		default:
			fmt.Println("unkown data", msg)
		}
	} //for
}

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp4", ":1200")
	checkError(err, "main")
	var s Server
	s.message = make(chan string, 20)
	s.clients = make(map[int]Client, 0)

	s.conn, err = net.ListenUDP("udp4", udpAddr)
	checkError(err, "main")

	go s.sendMsg()
	go s.printMsg()
	s.handle_msg()
	fmt.Println("done!")
}

func checkError(err error, funcName string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %s----- in func: %s", err.Error(), funcName)
		os.Exit(1)
	}
}
