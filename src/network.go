package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)


type Network struct {
    kademliaNodes *Kademlia
    srv *net.UDPAddr
}

type RPCMessage string

const (
    Ping RPCMessage = "PING"
    Store RPCMessage = "STORE"
    FindNode RPCMessage = "FIND_NODE"
    FindValue RPCMessage = "FIND_VALUE"
)

type MessageContactBuilder struct {
    msg RPCMessage
    contact *Contact
}

func CreateNewMessage(contact *Contact, msgType RPCMessage) MessageContactBuilder {

    msgBuilder := MessageContactBuilder{
        msg: msgType,
        contact: contact,
    }
    return msgBuilder
}

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
    udp_addr, err := network.SetPortBootstrap()
    //udp_addr, err := net.ResolveUDPAddr("udp", network.srv.AddrPort().String())
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
        network.HandleRPC(udp_connection, buffer)
    }
}

func (network *Network) SetPortBootstrap() (*net.UDPAddr, error){
    bootNodevar := FetchEnvVar("BN")
    if bootNodevar == 1 {
        boot_port := "5678"
        boot_ip, _ := net.LookupHost("bootNode")
        boot_server := boot_ip[0]+":"+boot_port
        udp_addr, err := net.ResolveUDPAddr("udp", boot_server)
        return udp_addr, err 
    }else {
        udp_addr, err := net.ResolveUDPAddr("udp", network.srv.AddrPort().String())
        return udp_addr, err
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

func BootnodeConnect(boot_addr *net.UDPAddr) (*net.UDPConn, error){

    // 1. Connect to boot node 
    conn, err := net.DialUDP("udp", nil ,boot_addr)
    if err != nil {
            fmt.Println("Error creating UDP connection: ", err)
            return conn, err
    }
    //defer conn.Close()
    return conn, err 
        
} 

func (network *Network) ListenJoin() {
	// TODO 
    bootNodevar := FetchEnvVar("BN")
    if bootNodevar == 1 {
        println("This is boot node!")
    } else if bootNodevar == 0 {
        boot_address, _ := net.LookupHost("bootNode")
        boot_port := os.Getenv("BNPT")
        boot_server := boot_address[0]+":"+boot_port
        boot_addr, _ := net.ResolveUDPAddr("udp", boot_server)
       
        // 1. Connect to boot node 
        conn, _ := BootnodeConnect(boot_addr) 
       
        // 2. add bootnode to k-bucket 
        network.SendRPC(Ping, conn)

        defer conn.Close()

        /*
        msg_test := []byte("hello boot node!")
        _, errs := conn.Write(msg_test)
        if errs != nil {
            fmt.Println("Error sending msg:", errs)
            return 
        }
        */
    }

}


func (network *Network) ShowNodeStatus(){
    // method for showing node status and its data state
    println(network.kademliaNodes.node_contact.me.String())
}


// RPC messages and RPC send manager: 

func (network *Network) SendRPC(rpcMessageType RPCMessage, connection *net.UDPConn){
    switch rpcMessageType{
    case Ping:
        // Send Ping RPC call to a specific node
        msg_ping := network.SendPingMessage(&network.kademliaNodes.node_contact.me, Ping)
        _, errs := connection.Write(msg_ping)
        if errs != nil {
            fmt.Println("Error sending msg: ", errs)
        }
    case Store:
        // Send Store RPC package
    case FindNode: 
        // Send FIND_NODE RPC package to specific node
        //connection.Write("SEND ME YOUR CONTACT")
    case FindValue:
        // Send FIND_VALUE RPC package to specific node client
    default:
        fmt.Println("Unknown RPC message type!")
    }

}

// This function is to handle RPC messages from the receiver side
func (network *Network) HandleRPC(connection *net.UDPConn, buffer []byte){

    n, client, err := connection.ReadFromUDP(buffer)
        if err != nil {
            //return err 
            fmt.Println(err)
        }
        returned_msg := MessageContactBuilder{}

        decoded_json_err := json.Unmarshal(buffer[:n], &returned_msg)
        if decoded_json_err != nil {
            fmt.Println(decoded_json_err)
        }
        //fmt.Printf("received %s from %s \n", string(buffer[:n]), client)
        fmt.Printf("Received: %+v from %s: ", returned_msg, client) 
}

// TODO receiver method for handling received UDP messages

func (network *Network) SendPingMessage(contact *Contact, msgType RPCMessage) []byte {
	// TODO
    //"GET_CONTACT_FROM_BN"
    //SERALIZE (CONTACT) ---> NODE THAT WANTS TO JOIN! 
    //msg_ping := []byte(Ping)
    new_msg := CreateNewMessage(contact, msgType)
    
    json_msg, err := json.Marshal(new_msg)
    if err != nil {
        return json_msg
    }
    fmt.Printf("Message to send: %v", json_msg)
    return json_msg


    
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// takes a 160-bit ID as an argument. 
    
    // NOTE: Whenever a node receives a communication from another, it updates the corresponding bucket.
    // If the contact already exists, it is moved to the end of the bucket.
   
    //1. Send byte of string to contact the boot node

    //2. The boot nodes contact is sent as response to calling node
    //network.kademliaNodes.node_contact.AddContact()

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
