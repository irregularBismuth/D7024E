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
    go kademliaNetwork.InvokeServer()
    src.RunCLI()
    //kademliaNetwork.BootstrapConnect()
   // go kademliaNetwork.ListenJoin()
    for {}
}
