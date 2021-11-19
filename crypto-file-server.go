package main

import (
	"fmt"
	"net"
	"os"

	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server"
	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/crypt"
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
		clientSign, symmetricKey, err := crypt.KeyExchange(conn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		go server.NewClient(conn, clientSign, symmetricKey).Run()
	}
}
