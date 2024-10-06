package main

import (
	"fmt"
	"net"
	"os"

	ncat "ncat/functions"
)

func main() {
	Server := ncat.CreateNewServer()
	arg := os.Args[1:]
	port := ":8989"

	if len(arg) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
	} else if len(arg) == 1 {
		if ncat.ValidPort(arg[0]) {
			arg[0] = port
		} else {
			fmt.Printf("Invalid port: %s\n", arg[0])
            return
		}
		
	} else {
		fmt.Println("Default port 8989 is used.")
	}

	listner, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile("png.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	Server.Content = content
	Server.Listen(listner)
}
