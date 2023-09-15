package src

import (
    "net"
    "strconv"
    "fmt"
    "os"
)


type Network struct {
    kademliaNodes *Kademlia 
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

func fetchEnvVar(envvar string) (bootNodevar int) {
    envVar:=os.Getenv(envvar)
    bootNodevar,err := strconv.Atoi(envVar)
    if err!=nil {
        fmt.Println("error converting BN to int",err)
    }
    return 
}

func Listen(ip string, port int) {
	// TODO 
    bootstrapNode:= fetchEnvVar("BN")
    if bootstrapNode == 1 {
        fmt.Println("peepo peepo")
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
    network.kademliaNodes.node_contact.AddContact(*contact)
    closest_contact := network.kademliaNodes.node_contact.FindClosestContacts(contact.ID, 3)
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
