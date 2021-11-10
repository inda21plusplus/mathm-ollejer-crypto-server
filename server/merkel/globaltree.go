package merkel

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
)

var GlobalTree = generateTree()

func generateTree() *Tree {
	root, lookup := genNode(2)

	return &Tree{root, lookup}
}

func genNode(depth uint) (Node, map[string]uint64) {
	if depth == 0 {
		node, id := genLeafNode()
		m := map[string]uint64{
			id: 0,
		}
		return node, m
	} else {
		return genBranchNode(depth)
	}
}

func genLeafNode() (*LeafNode, string) {
	data := make([]byte, 20)
	if _, err := rand.Read(data); err != nil {
		panic(err)
	}
	hm := hmac.New(sha256.New, []byte("secretcode"))
	if _, err := hm.Write(data); err != nil {
		panic(err)
	}
	sig := hm.Sum([]byte{})
	id := make([]byte, 18)
	if _, err := rand.Read(id); err != nil {
		panic(err)
	}
	return &LeafNode{
		Sig:      Hash(sig),
		FileData: data,
	}, string(id)
}

func genBranchNode(depth uint) (*BranchNode, map[string]uint64) {
	left, leftMap := genNode(depth - 1)
	right, rightMap := genNode(depth - 1)
	m := make(map[string]uint64, len(leftMap)+len(rightMap))
	for key, val := range leftMap {
		m[key] = val << 1
	}
	for key, val := range rightMap {
		m[key] = (val << 1) | 1
	}
	cat := append(left.GetHash(), right.GetHash()...)
	hasher := sha256.New()
	hasher.Write(cat)
	hash := hasher.Sum([]byte{})
	return &BranchNode{hash, left, right}, m
}
