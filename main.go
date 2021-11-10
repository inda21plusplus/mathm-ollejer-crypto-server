package main

import (
	"fmt"
	"net"
	"os"

	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server"
)

func main() {
	listener, err := net.Listen("tcp", ":10000")
	if err != nil {
		fmt.Println("Error binding socket:", err)
		os.Exit(1)
	}

	fmt.Println("Big chungus hack your file on", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
		}

		go server.NewClient(conn).Run()
	}
}
