package server

type ReadFileRequest struct {
	ID []byte
}

func (r ReadFileRequest) Handle(client *Client) {
}
