package server

import (
	"encoding/base64"

	e "github.com/inda21plusplus/mathm-ollejer-crypto-server/server/errors"
	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/merkle"
)

type Request struct {
	Kind    string `json:"type"`
	IDB64   string `json:"id"`
	SigB64  string `json:"signature"`
	DataB64 string `json:"data"`
}

func b64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func b64d(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func (req *Request) Handle(c *Client) interface{} {
	switch req.Kind {
	case "list":
		ids := merkle.GlobalTree.GetIDs()
		return map[string][]string{
			"ids": ids,
		}
	case "read":
		if len(req.IDB64) == 0 {
			return e.MissingParam("id")
		}
		sig, data, hashes, err := merkle.GlobalTree.ReadFile(req.IDB64)
		if err != nil {
			return err
		}
		res := struct {
			Sig        string                  `json:"signature"`
			Data       string                  `json:"data"`
			Validation []merkle.HashValidation `json:"validation"`
		}{
			string(sig),
			b64(data),
			make([]merkle.HashValidation, 0, len(hashes)),
		}
		for _, hash := range hashes {
			res.Validation = append(res.Validation, hash)
		}
		return res
	case "write":
		if len(req.IDB64) == 0 {
			return e.MissingParam("id")
		}
		if len(req.SigB64) == 0 {
			return e.MissingParam("signature")
		}
		data, err := b64d(req.DataB64)
		if err != nil {
			return e.BadRequest(err)
		}
		hashes, err := merkle.GlobalTree.WriteFile(req.IDB64, req.SigB64, data)
		if err != nil {
			return err
		}
		res := struct {
			Validation []merkle.HashValidation `json:"validation"`
		}{
			make([]merkle.HashValidation, 0, len(hashes)),
		}
		for _, hash := range hashes {
			res.Validation = append(res.Validation, hash)
		}
		return res
	}
	return e.BadRequest(nil)
}
