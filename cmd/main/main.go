package main

import (

	"fmt"
    "github.com/irregularBismuth/D7024E/src"
    "log"
    "net"
    "math/rand"
    "time"
)


func init() { rand.Seed(time.Now().UTC().UnixNano())}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}



func main(){

    // TODO...
    //var listener = src.Listen("127.0.0.1", )
    //fmt.Println("Socket open for listen at: ", )

    var id = src.NewRandomKademliaID() 
    a:=src.NewContact(id,GetOutboundIP().String());
    
    fmt.Println("This is random Node: ", id)
        

    fmt.Println(a.Address);
    src.Listen("localhost",8888)
    
    for {}
}
