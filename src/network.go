package src

import (
	"fmt"
	"log"
	"net"
	"os"
)


type Network struct {
    node *Kademlia 
}

func InitBootstrap(ip net.IP) {
    //bn := InitNode(ip)
    // call GetOutboundIP and assign the node based on the environment variable for the bootstrap node
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

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

func Listen(ip string, port int) {
	// TODO 
//    strCat:=ip + ":" + strconv.Itoa(port)
    master_node := os.Getenv("BN")
    if master_node == "1" {
        master_ip := GetOutboundIP()
        fmt.Println(master_ip)

    }
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
