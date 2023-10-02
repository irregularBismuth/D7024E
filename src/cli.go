package src 

import(
"os"
"fmt"
"bufio"
"regexp"
"strings"
)

//had to have this for it to run properly
var exitnooode=os.Exit

func FindCommands(cmd string,network *Network){
    //https://stackoverflow.com/questions/13737745/split-a-string-on-whitespace-in-go
    r:= regexp.MustCompile("[^\\s]+")
    w:=r.FindAllString(cmd,-1)

    for i,s := range w {
        w[i] = strings.ToLower(s)
    }


    if(w[0]=="help") {
        fmt.Println("List of available commands is : \n\n Put takes one argument object as UTF-8 format you want to store on the network e.g \n Put str123")
    } else if(w[0]=="put"){
        if(len(w)==2){
            network.SendStoreMessage(w[1])
            fmt.Println("Run put cmd")
        }else{ 
            fmt.Println("The put command takes in only 1 Argument of the object you want to store on the kademlia network\n")
        }
    } else if(w[0]=="get"){
        if(len(w)==2){
        network.SendFindDataMessage(w[1])
        }else{
            fmt.Println("You need to provide the hash of the object you want to retrieve \n")
        }
    }else if (w[0]=="exit"){
  
        fmt.Println("Terminate node")
        exitnooode(0)


    }else{
        fmt.Println("Command not found")
    }

}

func RunCLI(network *Network) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Enter CLI: via Docker attach Node")
    for{
        input,_ := reader.ReadString('\n')
        if(len(input) > 0){ FindCommands(input[:len(input)-1],network) }
    }
}
