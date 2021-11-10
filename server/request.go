package server

import (
	"encoding/base64"

	e "github.com/inda21plusplus/mathm-ollejer-crypto-server/server/errors"
	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/merkel"
)

type Request struct {
	Kind    string `json:"type"`
	IDB64   string `json:"id"`
	DataB64 string `json:"data"`
	SigB64  string `json:"signature"`
}

func (req *Request) Handle(c *Client) interface{} {
	switch req.Kind {
	case "read":
		id, err := base64.StdEncoding.DecodeString(req.IDB64)
		if err != nil {
			return e.BadRequest(err)
		}
		data, sig, hashes, err := merkel.GlobalTree.ReadFile(id)
		if err != nil {
			return err
		}
		res := struct {
			Data   string   `json:"data"`
			Sig    string   `json:"signature"`
			Hashes []string `json:"hashes"`
		}{
			base64.StdEncoding.EncodeToString(data),
			base64.StdEncoding.EncodeToString(sig),
			make([]string, 0, len(hashes)),
		}
		for _, hash := range hashes {
			res.Hashes = append(res.Hashes, base64.StdEncoding.EncodeToString(hash))
		}
		return res
	case "write":
		id, err := base64.StdEncoding.DecodeString(req.IDB64)
		if err != nil {
			return e.BadRequest(err)
		}
		data, err := base64.StdEncoding.DecodeString(req.DataB64)
		if err != nil {
			return e.BadRequest(err)
		}
		sig, err := base64.StdEncoding.DecodeString(req.SigB64)
		if err != nil {
			return e.BadRequest(err)
		}
		hashes, err := merkel.GlobalTree.WriteFile(id, data, sig)
		if err != nil {
			return err
		}
		res := struct {
			Hashes []string `json:"hashes"`
		}{
			make([]string, 0, len(hashes)),
		}
		for _, hash := range hashes {
			res.Hashes = append(res.Hashes, base64.RawStdEncoding.EncodeToString(hash))
		}
		return res
	}
	return e.BadRequest(nil)
}
