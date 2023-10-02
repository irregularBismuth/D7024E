package main

import (

    "github.com/irregularBismuth/D7024E/src"
    "math/rand"
    "time"
)


func init() { rand.Seed(time.Now().UTC().UnixNano())}

func main(){

    // TODO...
    
    //kademliaNode:=src.InitNode(src.GetOutboundIP())
    kademliaNetwork := src.InitNodeNetwork()
    go kademliaNetwork.ProcessRequestChannel()
    go kademliaNetwork.ListenServer()
    
    kademliaNetwork.JoinNetwork()

    //kademliaNetwork.SendStoreMessage("test")
    //kademliaNetwork.SendFindDataMessage("a94a8fe5ccb19ba61c4c0873d391e987982fbbd3")

 /*   var id = src.NewRandomKademliaID() 
    a:=src.NewContact(id,GetOutboundIP().String());
    
    fmt.Println("This is random Node: ", id)
        
   // fmt.Println("BN : ",os.Getenv("BN"))
    fmt.Println(a.Address);
    src.Listen("localhost",8888)
    fmt.Println("132 3212")*/


    for {}
}
