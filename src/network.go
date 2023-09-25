package src

import (
    "net"
    "fmt"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
)


type Network struct {
    kademliaNodes *Kademlia
    srv *UdpSocket
}

type RPCMessageTypes string

const (
    Ping RPCMessageTypes = "PING"
    Store RPCMessageTypes = "STORE"
    FindNode RPCMessageTypes = "FIND_NODE"
    FindValue RPCMessageTypes = "FIND_VALUE"
)

type RPCMessageBuilder struct {
    Msg RPCMessageTypes `json:"msg"` 
    Contact Contact `json:"contact"` // Source contact 
    ResponseData interface{} `json:"responseData"` //empty interface may hold values of any type
    Destination *net.UDPAddr `json:"destination"` // Destination address or what address to send the ResponseData too.
    IsRequest bool `json:"isRequest"`
}

// Should act as a thread safe communication channel
type UdpSocket struct {
    socketConnection *net.UDPConn
    serverAddress *net.UDPAddr
    response_channel chan RPCMessageBuilder
}

func NewUdpSocket(addr *net.UDPAddr) (UdpSocket){
    conn, err := net.ListenUDP("udp", addr)
    
    if err != nil {
        fmt.Println(err) 
    }
    fmt.Printf("udp_connection established: %v\n", conn.LocalAddr().String())

    return UdpSocket{
        socketConnection: conn,
        serverAddress: addr,
        response_channel: make(chan RPCMessageBuilder, 10),
    }
}

func CreateNewMessage(contact *Contact, msgType RPCMessageTypes, isRequest bool) RPCMessageBuilder {
    destination_addr,err := net.ResolveUDPAddr("udp",contact.Address)
    if err != nil {
        fmt.Println(err)
    }
    message := RPCMessageBuilder{}
    message.Contact = *contact
    message.Msg = msgType
    message.Destination = destination_addr
    message.IsRequest = isRequest
    return message
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

func (network *Network) ListenServer() error{
  
    // The ListenUDP method creates the server
    //udp_connection, err := net.ListenUDP("udp", network.srv.serverAddress) // (SERVER-SIDE)
    //udp_connection := network.srv.socketConnection // (SERVER-SIDE)

    //defer udp_connection.Close()
   
    // Need to add logic for incoming and outgoing packets handling - channels?
    buffer := make([]byte, 1024)
    for {
        network.RequestRPCHandler(buffer)
    }
}

func GetLocalAddr() (*net.UDPAddr, error){
    bootNodevar := FetchEnvVar("BN")
    if bootNodevar == 1 {
        udp_addr, err := GetBootnodeAddr()
        return udp_addr, err 
    }else {
        local_addr := GetOutboundIP()
        udp_addr, err := net.ResolveUDPAddr("udp", local_addr.String())
        return udp_addr, err
    }
}

func InitNodeNetwork() Network{
    // Initialize a new node and associated UDP socket network 
    local_address, _ := GetLocalAddr() // returns Local UDPAddr end point = ip address : port 
    new_node := InitNode(local_address)
    new_udpsocket := NewUdpSocket(local_address)

    // Create new network object struct containing its node and UDPAddr data
    new_network := Network{
        kademliaNodes: &new_node,
        srv: &new_udpsocket,
    }
    return new_network
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
    conn, err := net.DialUDP("udp", nil ,boot_addr)
    if err != nil {
            fmt.Println("Error creating UDP connection: ", err)
            return conn, err
    }
    //defer conn.Close()
    return conn, err 
        
}

func GetBootnodeAddr() (*net.UDPAddr, error){
    boot_address, _ := net.LookupHost("bootNode")
    boot_port := "5678"
    boot_server := boot_address[0]+":"+boot_port
    boot_addr, err := net.ResolveUDPAddr("udp", boot_server) 
    return boot_addr, err
}

func (network *Network) BootstrapJoinProcess(){
    // known bootstrap node connection
    boot_addr, _ := GetBootnodeAddr() 
    conn, _ := BootnodeConnect(boot_addr)
    contact := network.kademliaNodes.node_contact.me
    network.SendRPC(FindNode, &contact, conn)
    
    // Initialize new node that should join the network
    // 1. If it does not already have a NodeID N, it generates one
    // 2. It inserts the value of some known node C into appropriate bucket as its first contact
    // 3. It performs an iterative FIND_NODE for N 
    // 4. It refresh all buckets further away than its closes neighbor, which will be in the occupied bucket with the lowest index

}

// Should we return node object here or not? 
func (network *Network) JoinNetwork() {
	// TODO 
    bootNodevar := FetchEnvVar("BN")
    if bootNodevar == 1 {
        println("Initialiazing network - Creating bootnode!")
    } else if bootNodevar == 0 {
        println("This is not a boot node - Starting bootstraping join process!")
        network.BootstrapJoinProcess()
        //network.SendRPC(Ping, conn)
    }

}

func (network *Network) SendRPC(rpcMessageType RPCMessageTypes, contact *Contact, connection *net.UDPConn){
    switch rpcMessageType{
    case Ping:
        // Send Ping RPC call to a specific node
        //contact := network.kademliaNodes.node_contact.me
        msg_ping := network.SendPingMessage(contact, Ping)
        fmt.Printf("This is the contact: %v",contact)
        _, errs := connection.Write(msg_ping)
        if errs != nil {
            fmt.Println("Error sending msg: ", errs)
        }
    case Store:
        // Send Store RPC package
    case FindNode: 
        // Send FIND_NODE RPC package to specific node
        //sender_contact := network.kademliaNodes.node_contact.me
        msg_findnode := network.SendFindContactMessage(contact)
        _, err := connection.Write(msg_findnode)
        if err != nil {
            fmt.Println("Error sending msg: ", err)
        }
        //connection.Write("SEND ME YOUR CONTACT")
    case FindValue:
        // Send FIND_VALUE RPC package to specific node client
    default:
        fmt.Println("Unknown RPC message type!")
    }

}

// This function is to handle RPC messages from the receiver side
func (network *Network) RequestRPCHandler(buffer []byte){

    connection := network.srv.socketConnection // the request clients connection object
    _, _, err := connection.ReadFromUDP(buffer)
        if err != nil {
            //return err 
            fmt.Println(err)
        }
        var returned_msg MessageContactBuilder
        buffer_result := bytes.Trim(buffer, "\x00")
        decoded_json_err := json.Unmarshal(buffer_result, &returned_msg)
        if decoded_json_err != nil {
            fmt.Println(decoded_json_err)
        }
        if returned_msg.Msg == "PING"{
            fmt.Printf("Received: %#v", returned_msg) 
        } else if returned_msg.Msg == "PONG"{
            msg_type := returned_msg.Msg
            msg_address := returned_msg.ContactAddress
            fmt.Printf("Received: %#v from: %s", msg_type, msg_address) 
        }
    if err != nil {
        //return err 
        fmt.Println(err)
    }
    var returned_msg RPCMessageBuilder //deseralized json to struct object

    buffer_result := bytes.Trim(buffer,"\x00")
    decoded_json_err := json.Unmarshal(buffer_result, &returned_msg) //deseralize json 
    if decoded_json_err != nil {
            fmt.Println(decoded_json_err)
    }

    if returned_msg.IsRequest{
       // Handle and assign the request response and where to send the response too:
        switch returned_msg.Msg{
        case Ping:
            fmt.Printf("Received: %#v\n", returned_msg) 
            contact := network.kademliaNodes.node_contact.me
            //pong_json_msg := fmt.Sprintf(`{"Msg": "PONG", "Address": "%v"}`,contact)
            returned_msg.Contact = contact
            returned_msg.ResponseData = "PONG"
            returned_msg.IsRequest = false
            network.SendResponse(returned_msg)

            response_addr, _ := net.ResolveUDPAddr("udp", returned_msg.ContactAddress)

            conn, _ := BootnodeConnect(response_addr)
            conn.Write(pong_msg_response)
        }
        case FindNode:
            fmt.Printf("Received: %#v\n", returned_msg) 
            returned_msg.Contact = network.kademliaNodes.node_contact.me
            //target_id := returned_msg.Contact.ID
            //k_closest_nodes := network.kademliaNodes.node_contact.FindClosestContacts(target_id,3)
            returned_msg.ResponseData = network.kademliaNodes.node_contact.me
            returned_msg.IsRequest = false
            fmt.Println("Sending back my contact with KademliaID: ",returned_msg.Contact.ID.String())
            network.SendResponse(returned_msg)
        } 
    
    }else{
        // Process handler for message as a response by adding response message to the channel
        network.srv.response_channel<-returned_msg
    }
   
}

func (network *Network) SendResponse(response_msg RPCMessageBuilder){
    response_json, err := json.Marshal(response_msg)
    if err !=nil {
        fmt.Println("Error seralizing response message",err)
        return
    }
    connection := network.srv.socketConnection
    connection.WriteTo(response_json, response_msg.Destination)
}

// Process channel in a gorotuine - Response message 
func (network *Network) ResponseRPCHandler(msg RPCMessageBuilder){
    
    switch msg.Msg{
    case Ping:
        // Run ping logic 
        fmt.Printf("Response data: %s, from: %v ",msg.ResponseData, msg.Contact.Address)
    
    case FindNode:
        // Run kademlia logic
        fmt.Println("Response data: ", msg.ResponseData)
        recipent_contact := msg.Contact
        fmt.Println("Received a contact with KademliaID: ",recipent_contact.ID.String())
        network.kademliaNodes.node_contact.AddContact(recipent_contact)
    case FindValue:
        //
    }
     
}

func (network *Network) HandleResponseChannel(){
    for {
        select{
        case msg, ok := <-network.srv.response_channel:
            if !ok {
                fmt.Println("Channel closed")
                return 
            }
            //handle messages concurrently
            go network.ResponseRPCHandler(msg)
        default:
            time.Sleep(100*time.Millisecond)
        }
    }
}


// TODO receiver method for handling received UDP messages

func (network *Network) SendPingMessage(contact *Contact, msgType RPCMessageTypes) []byte {
	// TODO
    new_msg := CreateNewMessage(contact, msgType)


    fmt.Println("new message contact was created: ",new_msg.ContactID)

    isRequest := true
    new_msg := CreateNewMessage(contact, msgType, isRequest)
    json_msg, err := json.Marshal(new_msg)

    if err != nil {
        return json_msg
    }
    fmt.Printf("Message to send: %s\n", string(json_msg))
    return json_msg
    
}

func (network *Network) SendFindContactMessage(contact *Contact) []byte {
	// takes a 160-bit ID as an argument. 
    
    // NOTE: Whenever a node receives a communication from another, it updates the corresponding bucket.
    // If the contact already exists, it is moved to the end of the bucket.
   
    //1. Send byte of string to contact the boot node

    //2. The boot nodes contact is sent as response to calling node
    isRequest := true
    findcontact_msg := CreateNewMessage(contact, FindNode, isRequest)
    json_msg, err := json.Marshal(findcontact_msg)

    if err != nil {
        return json_msg
    }
    fmt.Printf("Message to send: %s\n", string(json_msg))
    return json_msg
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
