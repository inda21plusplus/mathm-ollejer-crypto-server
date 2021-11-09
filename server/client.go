package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

type Client struct {
	Conn net.Conn
	In   chan Request
	Out  chan Response
	Wg   sync.WaitGroup
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
		In:   make(chan Request),
		Out:  make(chan Response),
		Wg:   sync.WaitGroup{},
	}
}

func (c *Client) Start() {
	go c.sendData()
	go c.readData()
}

func (c *Client) readData() {
	c.Wg.Add(1)
	defer c.Wg.Done()
	decoder := json.NewDecoder(c.Conn)
	for {
		var req Request
		if err := decoder.Decode(&req); err != nil {
			if !errors.Is(err, io.EOF) {
				c.Out <- map[string]string{
					"error": "bad-request",
				}
			}
			close(c.In)
			return
		}
		c.In <- req
	}
}

func (c *Client) sendData() {
	c.Wg.Add(1)
	defer c.Wg.Done()
	encoder := json.NewEncoder(c.Conn)
	for res := range c.Out {
		if err := encoder.Encode(res); err != nil {
			fmt.Println("Error while sending data", err)
		}
	}
}
