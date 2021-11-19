package server

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"
)

type Client struct {
	Conn       net.Conn
	ID         *big.Int
	Key        []byte
}

func NewClient(conn net.Conn, clientID *big.Int, key []byte) *Client {
	return &Client{
		conn,
		clientID,
		key,
	}
}

func (c *Client) Run() {
	defer c.Conn.Close()

	cipher, err := chacha20poly1305.New(c.Key)
	if err != nil { return }

	reader := bufio.NewReader(c.Conn)
	for {
		var req Request

		{
			payload, err := reader.ReadSlice('\n')
			if err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println(err)
				}
				break
			}

			split := strings.Split(strings.Trim(string(payload), "\n"), " ")
			if len(split) != 2 { break }
			nonce, err := b64d(split[0])
			if err != nil { panic(err) }
			ciphertext, err := b64d(split[1])
			if err != nil { panic(err) }

			data, err := cipher.Open([]byte{}, nonce, ciphertext, []byte{})
			if err != nil {
				fmt.Println("open", err)
				break
			}

			if err := json.Unmarshal(data, &req); err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Println(err)
				}
				break
			}
		}

		res := req.Handle(c)
		var plaintext bytes.Buffer
		if err := json.NewEncoder(&plaintext).Encode(res); err != nil {
			fmt.Println(err)
			break
		}
		nonce := make([]byte, cipher.NonceSize())
		_, err = rand.Read(nonce)
		if err != nil {
			break
		}
		data := cipher.Seal([]byte{}, nonce, plaintext.Bytes(), []byte{})
		payload := []byte(fmt.Sprintf("%v %v\n", b64(nonce), b64(data)))

		fmt.Println(string(payload))

		n, err := c.Conn.Write(payload)
		if err != nil { break }
		if n != len(payload) { panic(n) }
	}
}
