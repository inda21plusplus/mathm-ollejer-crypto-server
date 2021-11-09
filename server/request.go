package server

type Request struct {
	Type string `json:"type"` // read, write, insert, delete
	ID   string `json:"id"`
	Data string `json:"data"`
}

// TODO: rename rawRequest, add Request interface and different Request structs
