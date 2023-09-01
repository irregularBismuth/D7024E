package main

import (
	"fmt"
    "github.com/irregularBismuth/D7024E/src"
)

func main(){
    // TODO...
    var listener = src.Listen("127.0.0.1", )
    fmt.Println("Socket open for listen at: ", )
    var id = src.NewRandomKademliaID() 
    fmt.Println("This is random Node: ", id)
}
