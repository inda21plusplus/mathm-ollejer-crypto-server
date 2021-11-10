package server

type WriteFileRequest struct {
	ID   []byte
	Data []byte
}

func (r WriteFileRequest) Handle(client *Client) {
}
