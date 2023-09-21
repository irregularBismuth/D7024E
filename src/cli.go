package src 

import(
"os"
"fmt"
"bufio")


func RunCLI() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("Enter CLI: press any key to continue")
    for{
        input,_ := reader.ReadString('\n')
        fmt.Println(input)
    }
}


