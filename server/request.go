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
	tree := merkle.GetTree(c.ID)
	switch req.Kind {
	case "list":
		ids := tree.GetIDs()
		return map[string][]string{
			"ids": ids,
		}
	case "exists":
		if len(req.IDB64) == 0 {
			return e.MissingParam("id")
		}
		exists := tree.Exists(req.IDB64)
		res := struct{
			Exists bool `json:"exists"`
		}{
			exists,
		}
		return res
	case "read":
		if len(req.IDB64) == 0 {
			return e.MissingParam("id")
		}
		sig, data, hashes, err := tree.ReadFile(req.IDB64)
		if err != nil {
			return err
		}
		res := struct {
			Sig        string                  `json:"signature"`
			Data       string                  `json:"data"`
			Validation []merkle.HashValidation `json:"validation"`
		}{
			string(sig),
			string(data),
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
		hashes, err := tree.WriteFile(req.IDB64, req.SigB64, []byte(req.DataB64))
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
