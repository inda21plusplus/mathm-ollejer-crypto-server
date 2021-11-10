package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"

	e "github.com/inda21plusplus/mathm-ollejer-crypto-server/server/errors"
	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/merkle"
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
	encoder := json.NewEncoder(c.Conn)
	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			if !errors.Is(err, io.EOF) {
				encoder.Encode(e.BadRequest(err))
			}
			break
		}
		res := req.Handle(c)
		if err := encoder.Encode(res); err != nil {
			fmt.Println(err)
			break
		}
		merkle.GlobalTree.Print()
	}

	if err := c.Conn.Close(); err != nil {
		fmt.Println("Error while closing connection:", err)
	}
}
