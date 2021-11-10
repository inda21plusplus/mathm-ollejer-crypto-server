package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

type Client struct {
	Conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}

func (c *Client) Run() {
	decoder := json.NewDecoder(c.Conn)
	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			if !errors.Is(err, io.EOF) {
				panic(err)
			}
			break
		}
		res := req.Handle(c)
		if err := json.NewEncoder(c.Conn).Encode(res); err != nil {
			panic(err)
			break
		}
	}

	if err := c.Conn.Close(); err != nil {
		fmt.Println("Error while closing connection:", err)
	}
}
