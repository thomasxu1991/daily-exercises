//简单的server-Client通信

package main

import (
	"fmt"
	"net"
	"os"
)
//异常检查
func errorCheck(err error,info string) (res bool){
	if(err != nil){
		fmt.Println(info+":"+err.Error())
		return false
	}else {
		return true
	}
}
//接收消息
func messageHandler(conn net.Conn,messages chan string){
	fmt.Println("Client IP:",conn.RemoteAddr().String())
	buf :=make([]byte,1024)
	for{
		length, err :=conn.Read(buf)
		if !errorCheck(err,"Connection"){
			conn.Close()
			break
		}
		if length>0{
			buf[length]=0
		}
		reciveStr := string(buf[0:length])
		messages <-reciveStr
	}
}
//发送数据
func echoHandler(conns *map[string]net.Conn,messages chan string){
	for{
		msg := <-messages
		fmt.Println(msg)
		for key,value := range *conns{
			fmt.Printf("Connectioned from:",key)
			_,err :=value.Write([]byte(msg))
			if err !=nil{
				fmt.Println(err.Error())
				delete(*conns,key)
			}
		}
	}
}
//启动服务器
func runServer(port string){
	tcpAddr,err :=net.ResolveTCPAddr("tcp4",":"+port);
	errorCheck(err,"ResolveTCPAddr")
	listen,err :=net.ListenTCP("tcp",tcpAddr)
	errorCheck(err,"ListenTCP")
	conns :=make(map[string]net.Conn)
	messages :=make(chan string,10)
	go echoHandler(&conns,messages)

	for{
		fmt.Println("Listening")
		conn,err := listen.Accept()
		errorCheck(err,"Accept")
		fmt.Println("Accepting")
		conns[conn.RemoteAddr().String()]=conn
		go messageHandler(conn,messages)

	}
}
//客户端
func client(conn net.Conn){
	var input string
	username :=conn.LocalAddr().String()
	for{
		fmt.Scanln(&input)
		if input =="/exit"{
			fmt.Println("exit")
			conn.Close()
			os.Exit(0);
		}
		lens,err :=conn.Write([]byte(username+":"+input))
		fmt.Println(lens)
		if(err != nil){
			fmt.Println(err.Error())
			conn.Close();
			break
		}
	}
}
//启动客户端
func runClient(tcpaddr string){
	tcpAddr,err := net.ResolveTCPAddr("tcp4",":"+tcpaddr)
	errorCheck(err,"ResoloveTCPAddr")
	conn,err :=net.DialTCP("tcp",nil,tcpAddr)
	errorCheck(err,"DialTCP")
	go client(conn)
	buf :=make([]byte,1024)
	for{
		length,err :=conn.Read(buf)
		if !errorCheck(err,"Connection"){
			conn.Close()
			fmt.Println("Server closed!")
			os.Exit(0);
		}
		fmt.Println(string(buf[0:length]))
	}
}

func main(){
	if len(os.Args)!=4{
		fmt.Println("Wrong input")
		os.Exit(0)
	}
	if os.Args[2]=="Server" && len(os.Args)==4{
		runServer(os.Args[3])
	}
	if os.Args[2]=="Client" && len(os.Args)==4{
		runClient(os.Args[3])
	}
}