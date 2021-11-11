package merkle

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
)

var GlobalTree = genTree()

func genTree() *Tree {
	root := genNode(2)

	lookup := make(map[string][]bool)
	genMap(root, make([]bool, 0), lookup)

	return &Tree{root, lookup}
}

func genMap(node *Node, traversion []bool, lookup map[string][]bool) {
	if node.IsLeaf() {
		lookup[node.ID] = traversion
	} else {
		l := make([]bool, len(traversion))
		copy(l, traversion)
		l = append(l, Left)
		genMap(node.Left, l, lookup)
		r := make([]bool, len(traversion))
		copy(r, traversion)
		r = append(r, Right)
		genMap(node.Right, r, lookup)
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
