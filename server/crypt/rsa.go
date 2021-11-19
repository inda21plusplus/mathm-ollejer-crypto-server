package crypt

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"strings"
)

const signedServerEnc string = "V1kLoNmlkfOAd4TnPbKyhJzGjRczMWkNKpEiP1pdaXhQD2+RY2MjLsdcfVA0cLTdkzXqrivIw6/BMsNfJuwwF4dZGBiuhxA56AAQDCbDPiDA+KSN62fhmLv2iKlwJajnht2PIKJmrUo6N0NkvWLoG6TAMEHy7rjTJEbPd5ZIaLxlGGHUPTf0LwsVh5t+Tn26wFvToB8ndxlK/U/LdAY1l4qgtAKjt285KIfMMZ3Fw7LmTJSI/EGCPcE3QdiBI5wQgClG9AzcPl+oh7jp0x5WX2AMqHUxnESpnywKuBo3oPlmmU9op3tPPv9APHPyIur0zjuNo1d7dDXm4jksHe8/IA=="

func SignServerKey() {
	signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
	Es, _, err := getServerKeys()

	h := sha256.New()
	es, err := encodePublicKey(&Es.PublicKey)
	if err != nil { panic(err) }
	h.Write(es)
	if err != nil { panic(err) }
	hashed := h.Sum([]byte{})

	signBytes, err := rsa.SignPKCS1v15(rand.Reader, signingKey, crypto.SHA256, hashed)
	if err != nil { panic(err) }
	sign := base64.StdEncoding.EncodeToString(signBytes)


	a, _ := encodePublicKey(&signingKey.PublicKey)
	fmt.Println(string(a))
	fmt.Println(sign)

	os.Exit(0)
}

func TestSignServerKey() {
	signingKeyPEM, _ := pem.Decode([]byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzMAhVgOov2aBjkocReUP
MR9fbdUqpoKNBOVkXokpZCyGUQbnjaeUHsbqbmoudDLYo5xy7LpO7EcjzcTHGIyT
zYmnTxZc86WvxxwTWFS3oJoRTuVbUSg3Yz+1tkspLVgAJJYA7kofgCitrzbCck3X
00bTRE1Pbc8klCXsJiWN32EtH/6Dt7EuxrHXuF+agRLd3dvix7LdVZbEXhu1bPCT
OBPWXC6TXRBI9gTsk8fz0OgTz4CUzfLWvkDauYJqQTR1OOZXxrT2qwUA3k5X6d6Y
oAGb3YFZbIJRkQhXe9aLmwWONIxc9kuzwdOG/CcEsoSYqA/UEzXwRrOSFfjgYe6e
sQIDAQAB
-----END PUBLIC KEY-----`))
	signingKey, err := x509.ParsePKIXPublicKey(signingKeyPEM.Bytes)
	s := signingKey.(*rsa.PublicKey)
	if err != nil { panic(err) }
	Es, _, _ := getServerKeys()

	h := sha256.New()
	es, err := encodePublicKey(&Es.PublicKey)
	if err != nil { panic(err) }
	h.Write(es)
	if err != nil { panic(err) }
	hashed := h.Sum([]byte{})

	signed, err := base64.StdEncoding.DecodeString(signedServerEnc)
	if err != nil { panic(err) }
	err = rsa.VerifyPKCS1v15(s, crypto.SHA256, hashed, signed)
	fmt.Println(err)
}

// Returns client's public RSA key
func KeyExchange(conn net.Conn) (*big.Int, []byte, error) {
	in := bufio.NewReader(conn)

	serverEnc, serverSign, err := getServerKeys()
	if err != nil { return nil, nil, err }

	fmt.Println("reading client keys")
	clientEnc, clientSign, err := readClientKeys(in)
	if err != nil { return nil, nil, err }
	fmt.Println("sending server keys")
	err = sendServerKeys(conn, clientEnc, serverEnc, serverSign, signedServerEnc)
	if err != nil { return nil, nil, err }
	fmt.Println("reading client halfrand")
	clientHalfRand, err := readClientHalfRand(in, serverEnc, clientSign)
	if err != nil { return nil, nil, err }
	var serverHalfRand [16]byte
	_, err = rand.Read(serverHalfRand[:])
	if err != nil { return nil, nil, err }
	fmt.Println("sending server halfrand")
	err = sendServerHalfRand(conn, serverSign, clientEnc, serverHalfRand[:])
	if err != nil { return nil, nil, err }

	symmetricKey := append(clientHalfRand, serverHalfRand[:]...)

	return clientSign.N, symmetricKey, nil
}

func readClientKeys(in *bufio.Reader) (*rsa.PublicKey, *rsa.PublicKey, error) {
	var clientKeys struct{
		EncB64 string `json:"Ec"`
		SigB64 string `json:"Vc"`
	}
	err := json.NewDecoder(in).Decode(&clientKeys)
	if err != nil { return nil, nil, err }

	encPem, err := base64.StdEncoding.DecodeString(clientKeys.EncB64)
	if err != nil { return nil, nil, err }
	sigPem, err := base64.StdEncoding.DecodeString(clientKeys.SigB64)
	if err != nil { return nil, nil, err }
	encBlock, _ := pem.Decode(encPem)
	if encBlock == nil { return nil, nil, errors.New("no key received from client") }
	sigBlock, _ := pem.Decode(sigPem)
	if sigBlock == nil { return nil, nil, errors.New("no key received from client") }
	enc, err := x509.ParsePKIXPublicKey(encBlock.Bytes)
	if err != nil { return nil, nil, err }
	enc2, ok := enc.(*rsa.PublicKey)
	if !ok { return nil, nil, errors.New("not rsa keys") }
	sig, err := x509.ParsePKIXPublicKey(sigBlock.Bytes)
	if err != nil { return nil, nil, err }
	sig2, ok := sig.(*rsa.PublicKey)
	if !ok { return nil, nil, errors.New("not rsa keys") }

	return enc2, sig2, nil
}

func sendServerKeys(conn net.Conn, clientEnc *rsa.PublicKey, serverEnc *rsa.PrivateKey, serverSign *rsa.PrivateKey, signedServerEnc string) error {
	es, err := encodePublicKey(&serverEnc.PublicKey)
	if err != nil { return err }
	vs, err := encodePublicKey(&serverSign.PublicKey)
	if err != nil { return err }
	esB64 := base64.StdEncoding.EncodeToString(es)
	vsB64 := base64.StdEncoding.EncodeToString(vs)

	cleartext, err := json.Marshal(map[string]interface{}{
		"Es": esB64,
		"Vs": vsB64,
		"Sa(Es)": signedServerEnc,
	})
	if err != nil { return err }

	var packet string
	for i := 0; i < len(cleartext); i += 128 {
		end := i + 128
		if end > len(cleartext) {
			end = len(cleartext)
		}
		chunk := cleartext[i:end]

		payload, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, clientEnc, chunk, []byte{})
		if err != nil { return err }

		if i != 0 {
			packet += ","
		}
		packet += base64.StdEncoding.EncodeToString(payload)
	}

	packet += "\n"

	n, err := conn.Write([]byte(packet))
	if err != nil { return err }

	if n != len(packet) {
		return errors.New("everything not written")
	}

	return nil
}

func encodePublicKey(key *rsa.PublicKey) ([]byte, error) {
	b, err := x509.MarshalPKIXPublicKey(key)
	if err != nil { return nil, err }
	block := &pem.Block{
		Type: "PUBLIC KEY",
		Bytes: b,
	}
	var out bytes.Buffer
	if err := pem.Encode(&out, block); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func readClientHalfRand(in *bufio.Reader, serverEnc *rsa.PrivateKey, clientSign *rsa.PublicKey) ([]byte, error) {
	packet, err := in.ReadString('\n')

	var cleartext []byte
	for _, chunk := range strings.Split(packet, ",") {
		ciphertext, err := base64.StdEncoding.DecodeString(chunk)
		if err != nil { return nil, err }
		if err != nil { return nil, err }
		c, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, serverEnc, ciphertext, []byte{})
		if err != nil { return nil, err }
		cleartext = append(cleartext, c...)
	}

	var payload struct{
		Data string `json:"data"`
		Sign string `json:"sign"`
	}
	json.Unmarshal(cleartext, &payload)

	data, err := base64.StdEncoding.DecodeString(payload.Data)
	if err != nil { return nil, err }

	hashed, err := hashSHA256(data)
	if err != nil { return nil, err }

	sign, err := base64.StdEncoding.DecodeString(payload.Sign)
	if err != nil { return nil, err }

	if err := rsa.VerifyPKCS1v15(clientSign, crypto.SHA256, hashed, sign); err != nil {
		return nil, err
	}

	if len(data) != 16 {
		return nil, errors.New("Invalid random length")
	}
	return data, nil
}

func sendServerHalfRand(conn net.Conn, serverSign *rsa.PrivateKey, clientEnc *rsa.PublicKey, serverHalfRand []byte) error {
	data := base64.StdEncoding.EncodeToString(serverHalfRand)

	hashed, err := hashSHA256(serverHalfRand)
	if err != nil { return err }

	sign, err := rsa.SignPKCS1v15(rand.Reader, serverSign, crypto.SHA256, hashed)
	if err != nil { return err }
	signB64 := base64.StdEncoding.EncodeToString(sign)

	cleartext, err := json.Marshal(map[string]string{
		"data": data,
		"sign": signB64,
	})

	if err != nil { return err }
	var packet string
	for i := 0; i < len(cleartext); i += 128 {
		end := i + 128
		if end > len(cleartext) {
			end = len(cleartext)
		}
		chunk := cleartext[i:end]

		payload, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, clientEnc, chunk, []byte{})
		if err != nil { return err }

		if i != 0 {
			packet += ","
		}
		packet += base64.StdEncoding.EncodeToString(payload)
	}

	packet += "\n"

	n, err := conn.Write([]byte(packet))
	if err != nil { return err }
	if n != len(packet) {
		return errors.New("not everything written")
	}

	if err != nil { return err }

	return nil
}

func getServerKeys() (*rsa.PrivateKey, *rsa.PrivateKey, error) {
	file, err := os.OpenFile("keys.pem", os.O_RDWR | os.O_CREATE, 0600)
	if err != nil { return nil, nil, err }
	defer file.Close()
	stat, err := file.Stat()
	if err != nil { return nil, nil, err }
	if stat.Size() > 0 {
		contents, err := ioutil.ReadAll(file)
		if err != nil { return nil, nil, err }
		block, rest := pem.Decode(contents)
		if block == nil || block.Type != "RSA PRIVATE KEY" {
			return nil, nil, errors.New("invalid key file")
		}
		encryptionKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil { return nil, nil, err }
		block, _ = pem.Decode(rest)
		if block == nil || block.Type != "RSA PRIVATE KEY" {
			return nil, nil, errors.New("invalid key file")
		}
		signingKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		return encryptionKey, signingKey, nil
	} else {
		encryptionKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil { return nil, nil, err }
		block := &pem.Block{
			Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(encryptionKey),
		}
		if err := pem.Encode(file, block); err != nil {
			return nil, nil, err
		}
		signingKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil { return nil, nil, err }
		block = &pem.Block{
			Type: "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(signingKey),
		}
		if err := pem.Encode(file, block); err != nil {
			return nil, nil, err
		}
		return encryptionKey, signingKey, nil
	}
}

func hashSHA256(data []byte) ([]byte, error) {
	h := sha256.New()
	_, err := h.Write([]byte(data))
	if err != nil { return nil, err }
	return h.Sum([]byte{}), nil
}
