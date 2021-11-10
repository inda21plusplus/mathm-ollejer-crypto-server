package server

import (
	"encoding/base64"

	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/merkel"
)

type WriteFileRequest struct {
	ID   []byte
	Data []byte
	Sig  []byte
}

type WriteFileResponse struct {
	Hashes []string `json:"hashes"`
}

func (r WriteFileRequest) Handle(client *Client) {
	hashes, err := merkel.GlobalTree.WriteFile(r.ID, r.Data, r.Sig)
	if err != nil {
		client.Out <- err
	} else {
		var h []string
		for _, hash := range hashes {
			h = append(h, base64.StdEncoding.EncodeToString(hash))
		}
		client.Out <- WriteFileResponse{h}
	}
}
