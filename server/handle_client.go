package server

import "fmt"

func (c *Client) Handle() {
loop:
	for req := range c.In {
		switch req.Type {
		case "read":
			fmt.Println("client reads file", req.ID)
		case "write":
			fmt.Println("client writes", req.Data, "to file", req.ID)
		default:
			c.Out <- map[string]interface{}{
				"error": "invalid-request-type",
			}
			break loop
		}
	}

	c.Wg.Wait()
	if err := c.Conn.Close(); err != nil {
		fmt.Println("Error while closing connection:", err)
	}
}
