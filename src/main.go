package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
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

		go handleConnection(conn)
	}
}

type Request struct {
	Type string `json:"type"` // read, write, insert, delete
	ID   string `json:"id"`
	Data string `json:"data"`
}

func readPacketsFromConn(conn net.Conn, wg *sync.WaitGroup, ch chan Request, out chan interface{}) {
	wg.Add(1)
	defer wg.Done()
	decoder := json.NewDecoder(conn)
	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Println("Invalid data recieved:", err)
				out <- map[string]string{
					"error": "bad-request",
				}
			}
			close(ch)
			return
		}
		ch <- req
	}
}

func sendPacketsToConn(conn net.Conn, wg *sync.WaitGroup, ch chan interface{}) {
	wg.Add(1)
	defer wg.Done()
	encoder := json.NewEncoder(conn)
	for packet := range ch {
		if err := encoder.Encode(packet); err != nil {
			fmt.Println("Error while sending data:", err)
		}
	}
}

func handleConnection(conn net.Conn) {
	in, out := make(chan Request), make(chan interface{})
	var wg sync.WaitGroup
	go sendPacketsToConn(conn, &wg, out)
	go readPacketsFromConn(conn, &wg, in, out)
loop:
	for req := range in {
		switch req.Type {
		case "read":
			fmt.Println("client reads file", req.ID)
		case "write":
			fmt.Println("client writes", req.Data, "to file", req.ID)
		default:
			out <- map[string]interface{}{
				"error": "invalid-request-type",
			}
			break loop
		}
	}

	wg.Wait()
	if err := conn.Close(); err != nil {
		fmt.Println("Error while closing connection:", err)
	}
}
