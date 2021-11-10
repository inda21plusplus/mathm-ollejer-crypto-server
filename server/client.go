package server

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/errors"
)

type Client struct {
	Conn net.Conn
	In   chan Request
	Out  chan interface{}
	Wg   sync.WaitGroup
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
		In:   make(chan Request),
		Out:  make(chan interface{}),
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
			if !goerrors.Is(err, io.EOF) {
				c.Out <- errors.BadRequest(err)
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
		if err := json.NewEncoder(c.Conn).Encode(res); err != nil {
			fmt.Println("Error sending data:", err)
		}
	}
}
