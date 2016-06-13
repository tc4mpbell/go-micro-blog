package main

import "os"
import "github.com/tc4mpbell/go-micro-auth"

func main() {
    argsWithoutProg := os.Args[1:]

    if len(argsWithoutProg) > 0 {
        username := argsWithoutProg[0]
        password := argsWithoutProg[1]

        auth.CreateAccount(username, password)
    } else {
        SetupPostEndpoints()
    }
}