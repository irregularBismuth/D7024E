package src

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
    "bytes"
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
    Key string `json:"key"`
    Value string `json:"value"`
    StringMessage string `json:"stringMessage"`
}

type MessageBuilder struct{
    MessageType RPCTypes `json:"msg"`
    RequestID string `json:"requestID"` //Request response identifier
    Contact Contact `json:"contact"`
    Response PayloadData `json:"payloadData"`
    SourceAddress *net.UDPAddr `json:"srcAddress"`
    DestinationAddress *net.UDPAddr `json:"dstAddress"`
    //IsRequest bool `json:"isRequest"`
}

func CreateRPC(msg_type RPCTypes, request_id string, contact *Contact, payload PayloadData, src_addr net.UDPAddr, dst_addr net.UDPAddr) *MessageBuilder{
    new_message := MessageBuilder{
        MessageType: msg_type,
        RequestID: request_id,
        Contact: *contact,
        Response: payload,
        SourceAddress: &src_addr,
        DestinationAddress: &dst_addr,
        //IsRequest: is_request,
    }
    return &new_message
}

// Send the RPC message to queue or channel for processing
func (network *Network) AddToRequestChannel(request_msg *MessageBuilder){
    // should feed unique request ID for processing the correct response. 
    network.srv.request_channel<-*request_msg
    //close(network.srv.request_channel)
}

func (network *Network) ProcessRequestChannel(){
    for {
        select {
        case request_msg, ok := <-network.srv.request_channel:
            if !ok{
                fmt.Println("Channel closed. Exiting goroutine!")
                return
            }
            go network.SendRequestRPC(&request_msg)
        default: 
            // Handle empty channels
            time.Sleep(500*time.Millisecond)
        }
    }
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

func (network *Network) RequestResponseWorker(buffer []byte){
    var request_msg MessageBuilder //deseralized json to struct object
    var response_msg PayloadData // Payload response data to send back to client
    
    buffer_result := bytes.Trim(buffer,"\x00")
    decoded_json_err := json.Unmarshal(buffer_result, &request_msg) //deseralize json 
    if decoded_json_err != nil {
            fmt.Println(decoded_json_err)
    }
    // Here we define and assign the PayloadData struct and what response to send back to the client.
    response_msg.Contact = request_msg.Contact
    switch request_msg.MessageType{
    case Ping:
        // Send pong response
        fmt.Printf("Received RPC request: %#v from: %s\n", request_msg.MessageType, &request_msg.Contact.Address) 
        response_msg.Contact = network.node.node_contact.me
            
        response_msg.StringMessage = "PONG"
        //network.SendResponse(returned_msg)
    case Store: 
        // Send store response 
    case FindNode:
        // Send FIND NODE response 
    case FindValue:
        // Send FIND VALUE response 
    case JoinNetwork:
        // Perform join network response 
    }
}





