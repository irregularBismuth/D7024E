package src

import (
<<<<<<< HEAD
	"fmt"
	"log"
	"net"
	"os"
=======
    "net"
    "strconv"
    "fmt"
    "os"
>>>>>>> 235f8d44dd5b61580813b9065e9f0050ff3be5f4
)


type Network struct {
<<<<<<< HEAD
    node *Kademlia 
}

func InitBootstrap(ip net.IP) {
    //bn := InitNode(ip)
    // call GetOutboundIP and assign the node based on the environment variable for the bootstrap node
=======
    kademliaNodes *kademlia     
>>>>>>> 235f8d44dd5b61580813b9065e9f0050ff3be5f4
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

<<<<<<< HEAD
// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
=======
func fetchEnvVar(envvar string) (bootNodevar int) {
    envVar:=os.Getenv(envvar)
    bootNodevar,err := strconv.Atoi(envVar)
    if err!=nil {
        fmt.Println("error converting BN to int",err)
    }
    return 
>>>>>>> 235f8d44dd5b61580813b9065e9f0050ff3be5f4
}

func Listen(ip string, port int) {
	// TODO 
<<<<<<< HEAD
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
=======
    bootstrapNode:= fetchEnvVar("BN")
    if bootstrapNode == 1 {
        fmt.Println("peepo peepo")
>>>>>>> 235f8d44dd5b61580813b9065e9f0050ff3be5f4
    }
    

}

// RPC calls: 

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// takes a 160-bit ID as an argument. 
    
    // NOTE: Whenever a node receives a communication from another, it updates the corresponding bucket.
    // If the contact already exists, it is moved to the end of the bucket.
    // If bucket is not full, the new contact is added at the end.
    network.node.node_contact.AddContact(*contact)
    closest_contact := network.node.node_contact.FindClosestContacts(contact.ID, 3)
    for _, contact := range closest_contact{
        // Call RPC "FIND_NODE" here
    }

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
