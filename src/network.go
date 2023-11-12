package src

import (
    "net"
    "fmt"
	"log"
	"os"
	"strconv"
)

type NetworkInterface interface {
    AddToRequestChannel(request_msg *MessageBuilder, request_err error)
    AddToResponseChannel(response_msg *PayloadData, response_err error)
    ReadResponseChannel(request MessageBuilder) *PayloadData
    
    BootstrapJoinProcess()
    JoinNetwork() 
    ProcessRequestChannel()
    
    FetchRPCResponse(rpc_type RPCTypes, rpc_id string, contact *Contact, dst_addr *net.UDPAddr) (*PayloadData, error) 
    SendRequestRPC(msg_payload *MessageBuilder) error
    SendResponseReply(response_msg *MessageBuilder)
    RequestResponseWorker(buffer []byte)
    AsynchronousFindNode(target_node *Contact, dst_addr *net.UDPAddr, response_ch chan<- PayloadData)
}

type Network struct {
    Kademlia *Kademlia
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

func InitNodeNetwork() *Network{
    // Initialize a new node and associated UDP socket network 
    local_address, _ := GetLocalAddr() // returns Local UDPAddr end point = ip address : port 
    new_node := InitNode(local_address)
    new_udpsocket := NewUdpSocket(local_address)

    // Create new network object struct containing its node and UDPAddr data
    new_network := &Network{
        Kademlia: &new_node,
        srv: &new_udpsocket,
    }
    new_network.Kademlia.SetNetworkInterface(new_network)
    fmt.Println("Comparing Network instances: ",new_network)
    fmt.Println("Comparing Network instances inside Kademlia: ",new_network.Kademlia.NetworkInterface)
    
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

func GetNodeUDPAddr(address string) (*net.UDPAddr, error){
    udp_addr, err := net.ResolveUDPAddr("udp", address)
    return udp_addr, err
}

func GetBootnodeAddr() (*net.UDPAddr, error){
    boot_address, _ := net.LookupHost("bootNode")
    boot_port := "5678"
    boot_server := boot_address[0]+":"+boot_port
    boot_addr, err := net.ResolveUDPAddr("udp", boot_server) 
    return boot_addr, err
}

// Send the RPC message to queue or channel for processing
func (network *Network) AddToRequestChannel(request_msg *MessageBuilder, request_err error){
    // should feed unique request ID for processing the correct response. 
    fmt.Printf("Received RPC request: %s from: %s with request id: %s\n", request_msg.MessageType, request_msg.Response.Contact.Address, request_msg.RequestID) 
    request_msg.Response.Error = request_err
    network.srv.request_channel<-*request_msg
}

// Send the RPC response to channel for processing
func (network *Network) AddToResponseChannel(response_msg *PayloadData, response_err error){
    // should feed unique request ID for processing the correct response. 
    response_msg.Error = response_err
    network.srv.response_channel<-*response_msg
}

func (network *Network) ReadResponseChannel(request MessageBuilder) *PayloadData{
    for response := range network.srv.response_channel{
        if response.ResponseID == request.RequestID{
            //response_ch <- response // transfer/send response payload to other channel
            fmt.Printf("Received RPC response: %s to: %s with response id: %s\n", response.StringMessage, response.Contact.Address, response.ResponseID)
            //fmt.Println("Response received updating contact bucket now!")
            //network.Kademlia.UpdateHandleBuckets(response.Contact) // response bucket update 
            //return &response, request_err
            return &response
        }
    }
    return nil
}

func (network *Network) BootstrapJoinProcess(){
    fmt.Println("Steps to perform:\n ☐ Add contact\n ☐ Node lookup on itself\n ☐ Refresh")
    
    boot_addr, _ := GetBootnodeAddr() 
    self_contact := network.Kademlia.node_contact.me
    
    // 1. Add contact 
    rpc_response, request_err := network.FetchRPCResponse(JoinNetwork, "my_rpc_id", &self_contact, boot_addr) 
    
    if request_err != nil {
        fmt.Println("JoinNetwork request unsuccessful: ", request_err.Error())
    }else{
        //network.node.node_contact.AddContact(rpc_response.Contact)
        //add if and else branch to check if the added contact was successful in the bootstrap process
        
        bucket_index := network.Kademlia.node_contact.getBucketIndex(rpc_response.Contact.ID)
        target_bucket := network.Kademlia.node_contact.buckets[bucket_index]

        if (target_bucket.DoesBucketContactExist(rpc_response.Contact)){
            fmt.Printf("☑ The contact: %+v from response node was added.\n ",rpc_response.Contact)
            fmt.Println("Proceeding with step 2 and 'Node lookup on itself'...")
            network.Kademlia.AsynchronousLookupContact(&self_contact)
            
        }else {
            fmt.Println("☐ The contact from response node was not added successfully!")
        }
 
        // 2. Node lookup on itself 
        //network.Kademlia.LookupContact(&self_contact)

        // Refresh step - picking a random ID in the bucket's range and perform a node search for that ID
    }
    
}

// Should we return node object here or not? 
func (network *Network) JoinNetwork() {
	// TODO 
    bootNodevar := FetchEnvVar("BN")
    if bootNodevar == 1 {
        println("Initialiazing network - Creating bootnode!\n")
    } else if bootNodevar == 0 {
        //println("Starting bootstraping join process!\n")
        network.BootstrapJoinProcess()

        // RPC tests here:
        
        /*
        boot_addr, _ := GetBootnodeAddr() 
        //conn, _ := BootnodeConnect(boot_addr)
        contact := network.Kademlia.node_contact.me
        rpc_response, request_err := network.FetchRPCResponse(Ping, "my_rpc_ping_id", &contact, boot_addr)
        if request_err != nil || rpc_response.Error != nil {
            fmt.Printf("Request unsuccessful: %s or: %s\n",request_err.Error(),rpc_response.Error.Error())
        }else{
            fmt.Println("Controll response received: ", rpc_response.StringMessage) 
        }
        */

    }
 
}

/*

// This function update appropriate k-bucket for the sender's node ID.
// The argument takes the target contacts bucket received from requests or response
func (network *Network) UpdateHandleBuckets(target_contact Contact){
   
    // Fetch the correct bucket location based on bucket index
    bucket_index := network.node.node_contact.getBucketIndex(target_contact.ID)
    target_bucket := network.node.node_contact.buckets[bucket_index]
    
    // if bucket is not full = add the node to the bucket 
    if target_bucket.Len() < GetMaximumBucketSize() {
        fmt.Printf("Bucket wasn't full adding contact: %s\n", target_contact.Address)
        network.node.node_contact.AddContact(target_contact)
    }else if target_bucket.Len() == GetMaximumBucketSize() {
        // If bucket is full - ping the k-bucket's least-recently seen node
        // least-recently node at the head & most-recently node at the tail 
        
        least_recently_node := target_bucket.GetLeastRecentlyNode()
        least_recently_addr, _ := net.ResolveUDPAddr("udp", least_recently_node.Address)
        my_contact := network.node.node_contact.me
        fmt.Printf("Bucket was full trying to ping recently-seen node: %s\n",least_recently_node.Address)
        rpc_ping, request_err := network.FetchRPCResponse(Ping, "bucket_full_ping_id", &my_contact, least_recently_addr)
        
        if request_err != nil || rpc_ping.Error != nil {
            //failed to response - removed from the k-bucket and new sender inserted at the tail
            fmt.Println("Request was unsuccessful, removing least-recently seen node from bucket: ", least_recently_node.Address)
            network.node.node_contact.RemoveTargetContact(*least_recently_node)
            
        }else if rpc_ping.Error == nil{
            //successful response -  contact is moved to the tail of the list 
            fmt.Printf("Ping was successful adding/moving contact: %s to tail\n", rpc_ping.Contact.Address)
            network.node.node_contact.AddContact(rpc_ping.Contact)
        }

    }
    network.ShowNodeBucketStatus()
    
}

// This function is to handle RPC messages from the receiver side

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
*/
func (network *Network) SendFindDataMessage(hash string) {
	fmt.Println("1. Reached Send GET message")
    //network.node.LookupData(hash)
}

func (network *Network) SendStoreMessage(data string) {
    byteString := []byte(data)    
    network.Kademlia.Store(byteString)
}


//func (network *Network) SendGetMessage(hash string) {
//    //fmt.Println("1. Reached Send GET message")
//    network.kademliaNodes.LookupData(hash)
//}

