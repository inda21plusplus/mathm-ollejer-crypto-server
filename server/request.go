package server

import (
	"encoding/base64"
)

type rawRequest struct {
	Kind    string `json:"type"` // read, write, insert, delete
	IDB64   string `json:"id"`
	DataB64 string `json:"data"`
}

type Request interface {
	Handle(*Client)
}

func (r rawRequest) toRequest() (Request, *Error) {
	switch r.Kind {
	case "read":
		id, err := base64.StdEncoding.DecodeString(r.IDB64)
		if err != nil {
			return nil, BadRequest(err)
		}
		return ReadFileRequest{id}, nil
	case "write":
		id, err := base64.StdEncoding.DecodeString(r.IDB64)
		if err != nil {
			return nil, BadRequest(err)
		}
		data, err := base64.StdEncoding.DecodeString(r.DataB64)
		if err != nil {
			return nil, BadRequest(err)
		}
		return WriteFileRequest{id, data}, nil
	}
	return nil, BadRequest(nil)
}
