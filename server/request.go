package server

import (
	"encoding/base64"
	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/errors"
)

type rawRequest struct {
	Kind    string `json:"type"` // read, write, insert, delete
	IDB64   string `json:"id"`
	DataB64 string `json:"data"`
	SigB64  string `json:"signature"`
}

type Request interface {
	Handle(*Client)
}

func (r rawRequest) toRequest() (Request, *errors.Error) {
	switch r.Kind {
	case "read":
		id, err := base64.StdEncoding.DecodeString(r.IDB64)
		if err != nil {
			return nil, errors.BadRequest(err)
		}
		return ReadFileRequest{id}, nil
	case "write":
		id, err := base64.StdEncoding.DecodeString(r.IDB64)
		if err != nil {
			return nil, errors.BadRequest(err)
		}
		data, err := base64.StdEncoding.DecodeString(r.DataB64)
		if err != nil {
			return nil, errors.BadRequest(err)
		}
		sig, err := base64.StdEncoding.DecodeString(r.SigB64)
		if err != nil {
			return nil, errors.BadRequest(err)
		}
		return WriteFileRequest{id, data, sig}, nil
	}
	return nil, errors.BadRequest(nil)
}
