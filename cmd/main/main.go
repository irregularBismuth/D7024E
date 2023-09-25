package main

import (
	"math/rand"
	"time"

	"github.com/irregularBismuth/D7024E/src"
)

func init() { rand.Seed(time.Now().UTC().UnixNano())}

func main(){

    // TODO...
    
    //kademliaNode:=src.InitNode(src.GetOutboundIP())
    kademliaNetwork := src.InitNodeNetwork()
    go kademliaNetwork.ListenServer()
    go kademliaNetwork.HandleResponseChannel() 
    //kademliaNetwork.BootstrapConnect()
    kademliaNetwork.JoinNetwork()

    kademliaNetwork.SendStoreMessage("Hello World")

    src.RunCLI()
    

    //kademliaNetwork.BootstrapConnect()
   // go kademliaNetwork.ListenJoin()
    for {}
}
