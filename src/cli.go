package src 

import(
"os"
"fmt"
"bufio"
"regexp"
)

func FindCommands(cmd string){
    //https://stackoverflow.com/questions/13737745/split-a-string-on-whitespace-in-go
    r:= regexp.MustCompile("[^\\s]+")
    w:=r.FindAllString(cmd,-1)
    if(w[0]=="Help" || w[0]=="help") {
        fmt.Println("List of available commands is : ")
    } else if(w[0]=="Put"){
        if(len(w)==2){
            //check 
            fmt.Println("Run put cmd")
        }else{ 
            fmt.Println("The put command takes in only 1 Argument of the object you want to store on the kademlia network\n")
        }
    }else{
        fmt.Println("Command not found")
    }

}


func RunCLI() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Enter CLI: via Docker attach Node")
    for{
        input,_ := reader.ReadString('\n')
        if(len(input) > 0){ FindCommands(input[:len(input)-1]) }
    }
}


