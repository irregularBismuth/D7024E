package src

import (
    "net"
 //   "strconv"
    "fmt"
)


type Network struct {
    
}

func handleConnection(connection net.Conn) {
    buf := make([]byte,1024)
    len,err := connection.Read(buf)
    if err!= nil{
        fmt.Printf("Error reading %#v\n",err)
        return
    }
    fmt.Printf("Message received %s\n",string(buf[:len]))
    connection.Write([]byte("Message received\n"))
    connection.Close()
}

func Listen(ip string, port int) {
	// TODO 
//    strCat:=ip + ":" + strconv.Itoa(port)
    addr:="localhost:8888"
    ln,err := net.Listen("tcp",addr)
    if err!=nil { }
    for { 
        conn,err := ln.Accept()
        if err!=nil {}
        go handleConnection(conn)
    }


}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
