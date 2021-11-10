package merkel

import "crypto/sha256"

var GlobalTree = generateTree()

func generateTree() *Tree {
	root, lookup := genNode(2)

	return &Tree{root, lookup}
}

func genNode(depth uint) (Node, map[string]uint64) {
	if depth == 0 {
		return genLeafNode(), nil
	} else {
		return genBranchNode(depth)
	}
}

func genLeafNode() *LeafNode {
	sig := sha256.Sum256([]byte("hej"))
	Sig := make([]byte, sha256.Size)
	copy(Sig[:], sig[:])
	return &LeafNode{
		Sig:      Hash(string(Sig)),
		FileData: []byte("hej"),
	}
}

func genBranchNode(depth uint) (*BranchNode, map[string]uint64) {
	left, leftMap := genNode(depth - 1)
	right, rightMap := genNode(depth - 1)
	m := make(map[string]uint64, len(leftMap)+len(rightMap))
	for key, val := range leftMap {
		m[key] = val << 1
	}
	for key, val := range rightMap {
		m[key] = (val << 1) & 1
	}
	cat := append(left.GetHash(), right.GetHash()...)
	hasher := sha256.New()
	hasher.Write(cat)
	hash := hasher.Sum([]byte{})
	return &BranchNode{hash, left, right}, m
}
