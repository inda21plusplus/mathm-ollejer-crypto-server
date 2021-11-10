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

func (c *Client) Handle() {
	go c.sendData()
	go c.readData()
	for req := range c.In {
		req.Handle(c)
	}
	close(c.Out)

	c.Wg.Wait()
	if err := c.Conn.Close(); err != nil {
		fmt.Println("Error while closing connection:", err)
	}
}

func (c *Client) readData() {
	c.Wg.Add(1)
	defer c.Wg.Done()
	decoder := json.NewDecoder(c.Conn)
	for {
		var rawReq rawRequest
		if err := decoder.Decode(&rawReq); err != nil {
			if !errors.Is(err, io.EOF) {
				c.Out <- BadRequest(err)
			}
			close(c.In)
			return
		}

		req, err := rawReq.toRequest()
		if err != nil {
			c.Out <- err
		} else {
			c.In <- req
		}
	}
}

func (c *Client) sendData() {
	c.Wg.Add(1)
	defer c.Wg.Done()
	for res := range c.Out {
		res.Respond(c)
	}
}
