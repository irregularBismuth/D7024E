package src

import (
	//"bytes"
	// "crypto/sha256"
	//"encoding/json"
	//"bytes"
	//"encoding/json"

	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	//"time"
)


type Network struct {
    node *Kademlia
    srv *UdpSocket
}

// Should act as a thread safe communication channel
type UdpSocket struct {
    socketConnection *net.UDPConn
    serverAddress *net.UDPAddr
    response_channel chan PayloadData
    request_channel chan MessageBuilder
    
}

func NewUdpSocket(addr *net.UDPAddr) (UdpSocket){
    conn, err := net.ListenUDP("udp", addr)
    
    if err != nil {
        fmt.Println(err) 
    }
    fmt.Printf(" Node network UDP connection established at: %v\n", conn.LocalAddr().String())

    return UdpSocket{
        socketConnection: conn,
        serverAddress: addr,
        response_channel: make(chan PayloadData, 10),
        request_channel: make(chan MessageBuilder, 10),
    }
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
        //network.RequestRPCHandler(buffer)
        network.RequestResponseWorker(buffer)
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
        node: &new_node,
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
    //conn, _ := BootnodeConnect(boot_addr)
    contact := network.node.node_contact.me
    
    // 1. Add contact 
    //response_contact := network.SendRPC(JoinNetwork, &contact, conn)
    rpc_response := network.FetchRPCResponse(JoinNetwork, "my_rpc_id", &contact, boot_addr, "") 
    network.node.node_contact.AddContact(rpc_response.Contact)
    fmt.Println("(1) Added the contact: ",rpc_response.Contact)
    // 2. Node lookup on itself 
    //self_contact := network.node.node_contact.me
    //resulting_lookup_contacts := network.node.LookupContact(network, &self_contact)
    //fmt.Println("Lookup result: ", resulting_lookup_contacts)
    
}

// Should we return node object here or not? 
func (network *Network) JoinNetwork() {
	// TODO 
    bootNodevar := FetchEnvVar("BN")
    if bootNodevar == 1 {
        println("Initialiazing network - Creating bootnode!\n")
    } else if bootNodevar == 0 {
        println("Starting bootstraping join process!\n")
        network.BootstrapJoinProcess()

        // RPC tests here:
        boot_addr, _ := GetBootnodeAddr() 
        //conn, _ := BootnodeConnect(boot_addr)
        contact := network.node.node_contact.me
        rpc_response := network.FetchRPCResponse(Ping, "my_rpc_ping_id", &contact, boot_addr, "") 
        fmt.Println("Controll response received: ", rpc_response.StringMessage) 
        //go network.SendRPC(Ping, &contact, conn)
    }

}

// This function is to handle RPC messages from the receiver side
/*
func (network *Network) RequestRPCHandler(buffer []byte){

    connection := network.srv.socketConnection // the request clients connection object
    _, _, err := connection.ReadFromUDP(buffer)
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
            fmt.Printf("Received RPC request: %#v from: %s\n", returned_msg.Msg, returned_msg.Contact.Address) 
            returned_msg.ResponseData.Contact = network.node.node_contact.me
            
            returned_msg.ResponseData.StringMessage = "PONG"
            returned_msg.IsRequest = false
            network.SendResponse(returned_msg)

        case FindNode:
            fmt.Printf("Received RPC request: %#v from: %s\n", returned_msg.Msg, returned_msg.Contact.Address) 
            target_id := returned_msg.Contact.ID
            k_closest_nodes := network.node.node_contact.FindClosestContacts(target_id,3)
            
            returned_msg.ResponseData.Contact = network.node.node_contact.me   
            returned_msg.ResponseData.Contacts = k_closest_nodes
            returned_msg.IsRequest = false
            network.SendResponse(returned_msg)
        
        case JoinNetwork:
            fmt.Printf("Received RPC request: %#v from: %s\n", returned_msg.Msg, returned_msg.Contact.Address) 
            
            returned_msg.ResponseData.Contact = network.node.node_contact.me
            returned_msg.ResponseData.StringMessage = "Bootstrap joining!"
            returned_msg.IsRequest = false
            network.SendResponse(returned_msg)
        } 
    
    }else{
        // Process handler for message as a response by adding response message to the channel
        network.srv.response_channel<-returned_msg
    }
   
}


func (network *Network) SendResponse(response_msg RPCMessageBuilder){
    response_json, err := json.Marshal(response_msg)
    dest_ip := response_msg.Destination.IP.String()
    dest_port := response_msg.Destination.Port
    fmt.Printf("Sending RPC response: %v to client: %s:%d \n",response_msg.ResponseData, dest_ip, dest_port)
    if err !=nil {
        fmt.Println("Error seralizing response message",err)
        return
    }
    connection := network.srv.socketConnection
    connection.WriteTo(response_json, response_msg.Destination)
}


// Process channel in a gorotuine - Response message 
func (network *Network) ResponseRPCHandler(msg RPCMessageBuilder) (*ResponseData) {
    // Receives the RPC response by reading the response channel
    // Should be able to return value for certain RPC calls...
    switch msg.Msg{
    case Ping:
        // Run ping logic 
        fmt.Printf("Response data: %v, from: %s\n",msg.ResponseData.StringMessage, msg.Contact.Address)
    
    case FindNode:
        // Run kademlia logic
        fmt.Printf("Response data: %v\n", msg.ResponseData.Contacts)
        recipent_contact := msg.Contact
        fmt.Printf("Received k-closest contacts from KademliaID: %s\n",recipent_contact.ID.String())
    case FindValue:
        //

    case JoinNetwork:
        fmt.Printf("Response data: %v\n", msg.ResponseData)
        fmt.Println("Steps to perform:\n 1. Add contact\n 2. Node lookup on itself\n 3. Refresh")

        //3. refresh 
    }
    return &msg.ResponseData
}

func (network *Network) HandleResponseChannel() *ResponseData {
    responseChannel := make(chan *ResponseData)

    go func() {
        for {
            select {
            case msg, ok := <-network.srv.response_channel:
                if !ok {
                    fmt.Println("Channel closed")
                    responseChannel <- nil
                    return
                }
                responsePayload := network.ResponseRPCHandler(msg)
                responseChannel <- responsePayload
            default:
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()

    // Read the response payload from the channel
    return <-responseChannel
}

func (network *Network) SendRequestMessage(connection *net.UDPConn, msg []byte){
     _, err := connection.Write(msg)
    if err != nil {
        fmt.Println("Error sending msg: ", err)
    }
}

// TODO receiver method for handling received UDP messages
func (network *Network) JoinNetworkMessage(contact *Contact, msgType RPCTypes) []byte{
    isRequest := true
    new_msg := CreateNewMessage(contact, msgType, isRequest)
    json_msg, err := json.Marshal(new_msg)
    if err != nil {
        return json_msg
    }
    fmt.Printf("RPC request message to send %s\n", string(json_msg))
    return json_msg
}

func (network *Network) SendPingMessage(contact *Contact, msgType RPCTypes) []byte {
	// TODO
    isRequest := true
    new_msg := CreateNewMessage(contact, msgType, isRequest)
    json_msg, err := json.Marshal(new_msg)

    if err != nil {
        return json_msg
    }
    fmt.Printf("RPC request message to send: %s\n", string(json_msg))
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
    fmt.Printf("RPC request message to send: %s\n", string(json_msg))
    return json_msg
}

func (network *Network) SendFindDataMessage(hash string) {
	//fmt.Println("1. Reached Send GET message")
    network.kademliaNodes.LookupData(hash)
}

func (network *Network) SendStoreMessage(data string) {
    byteString := []byte(data)    
    network.kademliaNodes.Store(byteString)
}

//func (network *Network) SendGetMessage(hash string) {
//    //fmt.Println("1. Reached Send GET message")
//    network.kademliaNodes.LookupData(hash)
//}
*/
