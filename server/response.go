package server

import (
	"encoding/json"
	"fmt"
)

type Response interface {
	Respond(client *Client)
}

func (e Error) Respond(client *Client) {
	if err := json.NewEncoder(client.Conn).Encode(map[string]interface{}{
		"error": e.message,
		"inner": e.inner,
	}); err != nil {
		fmt.Println("Error sending data:", err)
	}
}
