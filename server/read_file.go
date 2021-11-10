package server

import (
	"encoding/base64"

	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/merkel"
)

type ReadFileRequest struct {
	ID []byte
}

type ReadFileResponse struct {
	Data      string   `json:"data"`
	Signature string   `json:"signature"`
	Hashes    []string `json:"hashes"`
}

func (r ReadFileRequest) Handle(client *Client) {
	data, sig, hashes, err := merkel.GlobalTree.ReadFile(r.ID)
	if err != nil {
		client.Out <- err
	} else {
		data := base64.StdEncoding.EncodeToString(data)
		sig := base64.StdEncoding.EncodeToString(sig)
		var h []string
		for _, hash := range hashes {
			h = append(h, base64.StdEncoding.EncodeToString(hash))
		}
		client.Out <- ReadFileResponse{data, sig, h}
	}
}
