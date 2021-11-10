package merkle

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

var GlobalTree = genTree()

func genTree() *Tree {
	root := genNode(2)

	lookup := genMap(root, []bool{})
	fmt.Println(lookup)

	return &Tree{root, lookup}
}

func genMap(node *Node, traversion []bool) map[string][]bool {
	if node.IsLeaf() {
		return map[string][]bool{node.ID: traversion}
	} else {
		m := make(map[string][]bool)
		for id, val := range genMap(node.Left, append(traversion, Left)) {
			fmt.Println(traversion, id, val)
			m[id] = val
		}
		for id, val := range genMap(node.Right, append(traversion, Right)) {
			fmt.Println(traversion, id, val)
			m[id] = val
		}
		return m
	}
}

func genNode(depth uint) *Node {
	if depth == 0 {
		node := genLeafNode()
		return node
	} else {
		return genBranchNode(depth)
	}
}

func genLeafNode() *Node {
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
	return LeafNode(b64(id), b64(sig), []byte(b64(data)))
}

func genBranchNode(depth uint) *Node {
	left := genNode(depth - 1)
	right := genNode(depth - 1)
	return BranchNode(left, right)
}
