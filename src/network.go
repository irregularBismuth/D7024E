package src

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
    "context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/integration-cli/cli"
)


type Network struct {
    kademliaNodes *Kademlia
    srv *net.UDPAddr
}

type RPCMessage int

const (
    Unknown RPCMessage = iota
    Ping 
    Store
    FindNode
    FindValue
)

// Get preferred outbound ip of this machine - retreving the local (source) address 
func GetOutboundIP() net.Addr {
    conn, err := net.Dial("udp", "8.8.8.8:80") //The dial function connects to a server (CLIENT)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    localAddr := conn.LocalAddr()
    return localAddr
}

func (network *Network) InvokeServer() error{
   
    // binding UDP server through resolve (what port to listen to...)
    udp_addr, err := net.ResolveUDPAddr("udp", network.srv.AddrPort().String())
    if err != nil {
        return err
    }
    
    // The ListenUDP method creates the server
    udp_connection, err := net.ListenUDP("udp", udp_addr) // (SERVER-SIDE)
    if err != nil {
        fmt.Println("Error creating UDP connection:", err)
        return err
    }
    fmt.Printf("udp_connection established: %v\n", udp_connection.LocalAddr().String())
    defer udp_connection.Close()
   
    // HandleConnection logic should go here?
    // Need to add logic for incoming and outgoing packets handling - channels?
    buffer := make([]byte, 1024)
    for {
        _, _, err := udp_connection.ReadFromUDP(buffer)
        if err != nil {
            return err 
        }
    }
}

func InitNodeNetwork() Network{
    // Initialize a new node UDP network 
    local_address := GetOutboundIP() // returns LocalAddr = ip address : port 
    local_addr := local_address.(*net.UDPAddr) // creates UDPAddr struct object 
    new_node := InitNode(local_address)
   
    // Create new network object struct containing its node and UDPAddr data
    new_network := Network{
        kademliaNodes: &new_node,
        srv: local_addr,
    }
    return new_network

} 

// Server side - handle receving incoming messages
func HandleConnection(connection net.Conn) {
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

func FetchEnvVar(envvar string) (bootNodevar int) {
    envVar:=os.Getenv(envvar)
    bootNodevar,err := strconv.Atoi(envVar)
    if err!=nil {
        fmt.Println("error converting BN to int",err)
    }
    return 
}

func Listen(ip string, port int) {
	// TODO 
    bootstrapNode:= FetchEnvVar("BN")
    if bootstrapNode == 1 {
        fmt.Println("peepo peepo")
    }
}


func (network *Network) ShowNodeStatus(){
    // method for showing node status and its data state
    println(network.kademliaNodes.node_contact.me.String())
}

func (network *Network) BootstrapConnect(){
    bs_address, _ := net.LookupHost("bootNode")
    fmt.Println(bs_address)
    bn_container_id := "bootNode"

    //container_info, err := 
    //conn, _ := net.DialUDP("udp", nil, bs_address[0])
    
}

// RPC messages and RPC send manager: 

func (network *Network) SendRPC(rpcMessageType RPCMessage, connection *net.UDPConn, address *net.UDPAddr){
    switch rpcMessageType{
    case Ping:
        // Send Ping RPC call to a specific node
    case Store:
        // Send Store RPC package
    case FindNode: 
        // Send FIND_NODE RPC package to specific node
    case FindValue:
        // Send FIND_VALUE RPC package to specific node client
    default:
        fmt.Println("Unknown RPC message type!")
    }

}

// TODO receiver method for handling received UDP messages

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// takes a 160-bit ID as an argument. 
    
    // NOTE: Whenever a node receives a communication from another, it updates the corresponding bucket.
    // If the contact already exists, it is moved to the end of the bucket.
    // If bucket is not full, the new contact is added at the end.
    //network.kademliaNodes.node_contact.AddContact(*contact)
    // Bootstrap node logic for initializing contact to the network

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
