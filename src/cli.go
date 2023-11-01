package src

import (
	"bufio"
	"fmt"
	"net"

	//"net"
	"os"
	"regexp"
)

func (network *Network) FindCommands(cmd string){
    //https://stackoverflow.com/questions/13737745/split-a-string-on-whitespace-in-go
    r:= regexp.MustCompile("[^\\s]+")
    w:=r.FindAllString(cmd,-1)
    if(w[0]=="Help" || w[0]=="help") {
        fmt.Println("List of available commands is : \n\n Put takes one argument object as UTF-8 format you want to store on the network e.g \n Put str123")
    } else if(w[0]== ("Put")) || (w[0]==("put")){
        if(len(w)==2){
            //check 
            fmt.Println("1. Running PUT cmd")
            network.node.Store([]byte(w[1]))
            //target_addr, _ := net.ResolveUDPAddr("udp", network.node.node_contact.me.Address)
            target_contact := network.node.node_contact.me
            k_targets := network.node.node_contact.FindClosestContacts(network.node.node_contact.me.ID, 3)
            for i := 0; i < len(k_targets); i++ {
                k_target := k_targets[i]
                target_addr, _ := net.ResolveUDPAddr("udp", k_target.Address)
                network.FetchRPCResponse(Store, "", &target_contact, target_addr, w[1])
            }
        }else{ 
            fmt.Println("The put command takes in only 1 Argument of the object you want to store on the kademlia network")
        }
    } else if (w[0]==("Get") || (w[0]==("get"))) {
        if(len(w)==2){
            byte := []byte(w[1])
            hash := network.node.Hash(byte) 
            fmt.Println("1. Running GET cmd")
            original, exists := network.node.LookupData(network, hash)
            fmt.Println(string(original), exists)
           /* target_contact := network.node.node_contact.me
            k_targets := network.node.node_contact.FindClosestContacts(network.node.node_contact.me.ID, 3)
            for i := 0; i < len(k_targets); i++ {
                k_target := k_targets[i]
                target_addr, _ := net.ResolveUDPAddr("udp", k_target.Address)
                network.FetchRPCResponse(FindValue, "", &target_contact, target_addr, w[1])
            } */
        }else{
            fmt.Println("The Get command takes in only 1 Argument of the object you want to download from the kademlia network")
        }
    } else if (w[0]== ("Exit")) || (w[0]==("exit")){
        fmt.Println("Terminate node")
    }else{
        fmt.Println("Command not found")
    }

}


func RunCLI(network *Network) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Enter CLI: via Docker attach Node")
    for{
        input,_ := reader.ReadString('\n')
        if(len(input) > 0){ network.FindCommands(input[:len(input)-1]) }
    }
}


