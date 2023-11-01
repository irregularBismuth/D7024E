package src

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
    //"bytes"
)

type RPCTypes string

const (
    Ping RPCTypes = "PING"
    Store RPCTypes = "STORE"
    FindNode RPCTypes = "FIND_NODE"
    FindValue RPCTypes = "FIND_VALUE"
    JoinNetwork RPCTypes = "JOIN_NETWORK"
)

type PayloadData struct {
    Contacts []Contact `json:"contacts"`
    Contact Contact `json:"contact"`
    ResponseID string `json:"responseID"`
    Key string `json:"key"`
    Value []byte `json:"value"`
    StringMessage string `json:"stringMessage"`
}

type MessageBuilder struct{
    MessageType RPCTypes `json:"msg"`
    RequestID string `json:"requestID"` //Request response identifier
    Response PayloadData `json:"payloadData"`
    SourceAddress *net.UDPAddr `json:"srcAddress"`
    DestinationAddress *net.UDPAddr `json:"dstAddress"`
    IsRequest bool `json:"isRequest"`
}

func CreateRPC(msg_type RPCTypes, request_id string, payload PayloadData, src_addr net.UDPAddr, dst_addr net.UDPAddr) *MessageBuilder{
    new_message := MessageBuilder{
        MessageType: msg_type,
        RequestID: request_id,
        //Contact: contact,
        Response: payload,
        SourceAddress: &src_addr,
        DestinationAddress: &dst_addr,
        IsRequest: true, //Sets default to true 
    }
    return &new_message
}

// Send the RPC message to queue or channel for processing
func (network *Network) AddToRequestChannel(request_msg *MessageBuilder){
    // should feed unique request ID for processing the correct response. 
    fmt.Printf("Received RPC request: %s from: %s with request id: %s\n", request_msg.MessageType, request_msg.Response.Contact.Address, request_msg.RequestID) 
    network.srv.request_channel<-*request_msg
}

// Send the RPC response to channel for processing
func (network *Network) AddToResponseChannel(response_msg *PayloadData){
    // should feed unique request ID for processing the correct response. 
    //fmt.Printf("Adding RPC response to channel: %s to: %s with response id: %s\n", response_msg.StringMessage, response_msg.Contact.Address, response_msg.ResponseID) 
    network.srv.response_channel<-*response_msg
}

func (network *Network) ProcessRequestChannel(){
    for {
        select {
        case request_msg, ok := <-network.srv.request_channel:
            if !ok{
                fmt.Println("Channel closed. Exiting goroutine!")
                return
            }
            //go network.SendRequestRPC(&request_msg)
            go network.SendResponseReply(&request_msg)

        default: 
            // Handle empty channels
            time.Sleep(500*time.Millisecond)
        }
    }
}

// This will send a RPC request and wait for response value to return from the response channel 
func (network *Network) FetchRPCResponse(rpc_type RPCTypes, rpc_id string, contact *Contact, dst_addr *net.UDPAddr, hash string) *PayloadData{
    var payload PayloadData
    var default_byte []byte

    src_addr := network.srv.serverAddress
    switch rpc_type{
    case FindNode:

    case FindValue:
        var default_byte []byte
        payload = PayloadData{nil, *contact, "","",default_byte, hash}

    }
    if hash != ""{
        payload = PayloadData{nil, *contact, "","",default_byte, hash}
    }else {
        payload = PayloadData{nil, *contact, "","",default_byte, hash}
    }

    //src_payload := PayloadData{nil, *contact,"","","",""} //empty request payload 
    new_request := CreateRPC(rpc_type, rpc_id, payload, *src_addr, *dst_addr)
    network.SendRequestRPC(new_request)

    for response := range network.srv.response_channel{
        if response.ResponseID == rpc_id{
            fmt.Printf("Received RPC response: %s to: %s with response id: %s\n", response.StringMessage, response.Contact.Address, response.ResponseID)
            
            return &response
        }
    }
    return nil
}

func (network *Network) SendRequestRPC(msg_payload *MessageBuilder){
    dest_ip := msg_payload.DestinationAddress.IP.String()
    dest_port := msg_payload.DestinationAddress.Port
    
    request_json, err := json.Marshal(msg_payload)
    fmt.Printf("Sending RPC request: %s to client: %s:%d \n",string(request_json), dest_ip, dest_port)
    if err !=nil {
        fmt.Println("Error seralizing response message",err)
        return
    }
    network.srv.socketConnection.WriteTo(request_json, msg_payload.DestinationAddress)
}

// Takes unmarshalled request data and process the response payload to send back to the client
func (network *Network) SendResponseReply(response_msg *MessageBuilder){
    response_msg.IsRequest = false
    response_msg.DestinationAddress = response_msg.SourceAddress // Destination address = Source address
    response_msg.Response.ResponseID = response_msg.RequestID
    dest_ip := response_msg.DestinationAddress.IP.String()
    dest_port := response_msg.DestinationAddress.Port

    switch response_msg.MessageType{
    case Ping:
        response_msg.Response.Contact = network.node.node_contact.me
        response_msg.Response.StringMessage = "PONG"

    case Store:
        data := response_msg.Response.StringMessage
        network.node.Store([]byte(data))
        fmt.Printf("Download from node: %s , SUCCESS", response_msg.Response.Contact.Address)

    case FindNode:
        target_id := response_msg.Response.Contact.ID
        k_closest_nodes := network.node.node_contact.FindClosestContacts(target_id,3)
        response_msg.Response.Contacts = k_closest_nodes
        response_msg.Response.Contact = network.node.node_contact.me 

    case FindValue:
        hash := response_msg.Response.StringMessage
        response_hash_value := network.node.data[hash]
        
        response_msg.Response.Value = response_hash_value

    case JoinNetwork:
        response_msg.Response.Contact = network.node.node_contact.me
        response_msg.Response.StringMessage = "Bootstrap joining!"
    }
    response_json, err := json.Marshal(response_msg)
    fmt.Printf("Sending RPC response: %v to client: %s:%d \n",response_msg.Response, dest_ip, dest_port)
    if err !=nil {
        fmt.Println("Error seralizing response message",err)
        return
    }
    network.srv.socketConnection.WriteTo(response_json, response_msg.DestinationAddress)
}

// This is a listener for receiving RCP requests via ReadFromUDP(buffer) and sends back to client 
func (network *Network) RequestResponseWorker(buffer []byte){
    var request_msg MessageBuilder //deseralized json to struct object
    var response_msg PayloadData // Payload response data to send back to client
   
    // ##################################### Fetch from UDP socket 
    connection := network.srv.socketConnection // the request clients connection object
    n, _, err := connection.ReadFromUDP(buffer) 
    if err != nil {
        //return err 
        fmt.Println(err)
    }
    
    // ##################################### Unmarshal data
    //buffer_result := bytes.Trim(buffer,"\x00")
    decoded_json_err := json.Unmarshal(buffer[:n], &request_msg) //deseralize json 
    if decoded_json_err != nil {
            fmt.Println(decoded_json_err)
    }

    if request_msg.IsRequest{
        // If RPC request was received we add it to the request channel
        // this will be handled in a goroutine function "ProcessRequestChannel" which will send a response back
        go network.AddToRequestChannel(&request_msg)
       
    }else{
        // If RPC response was received we add it to the response channel 
        response_msg = request_msg.Response
        go network.AddToResponseChannel(&response_msg)
    }
}



