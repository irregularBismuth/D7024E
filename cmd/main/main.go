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
<<<<<<< HEAD
    go kademliaNetwork.ListenServer()
    go kademliaNetwork.HandleResponseChannel() 
    //kademliaNetwork.BootstrapConnect()
    kademliaNetwork.JoinNetwork()

 /*   var id = src.NewRandomKademliaID() 
    a:=src.NewContact(id,GetOutboundIP().String());
    
    fmt.Println("This is random Node: ", id)
        
   // fmt.Println("BN : ",os.Getenv("BN"))
    fmt.Println(a.Address);
    src.Listen("localhost",8888)
    fmt.Println("132 3212")*/


=======
    go kademliaNetwork.InvokeServer()
    src.RunCLI()
    //kademliaNetwork.BootstrapConnect()
   // go kademliaNetwork.ListenJoin()
>>>>>>> bf87628ae14c1b1c921e773143154eb6233a96c6
    for {}
}
