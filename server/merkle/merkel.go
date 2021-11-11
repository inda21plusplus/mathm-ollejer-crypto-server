package merkle

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/errors"
)

func b64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func b64d(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func min(a, b uint) uint {
	if a < b {
		return a
	} else {
		return b
	}
}

type HashValidation struct {
	Hash      string `json:"hash"`
	Direction bool   `json:"is-right"`
}

const (
	Left  bool = false
	Right bool = true
)

type Tree struct {
	root              *Node
	traversion_lookup map[string][]bool
}

type Node struct {
	ID       string // is "" unless this is a leaf
	Hash     string // is signature of file if this is a leaf
	Left     *Node  // nil if this is a leaf
	Right    *Node  // nil if this is a leaf
	FileData []byte // nil unless this is a leaf
	MinDepth uint   // 0 if this is a leaf
}

func (t *Tree) Print() {
	fmt.Print("\n\n")
	for key, val := range t.traversion_lookup {
		fmt.Printf("%v: %v\n", key, val)
	}
	fmt.Print("\n")
	t.root.Print(0)
	fmt.Print("\n\n\n")
}

func (n *Node) Print(indent int) {
	if n.IsLeaf() {
		fmt.Println(strings.Repeat(" ", indent), n.ID)
	} else {
		n.Right.Print(indent + 4)
		fmt.Println(strings.Repeat(" ", indent), n.Hash)
		n.Left.Print(indent + 4)
	}
}

func BranchNode(left, right *Node) *Node {
	h := sha256.New()
	l, _ := b64d(string(left.Hash))
	r, _ := b64d(string(right.Hash))
	// TODO: handle errors here
	h.Write(append(l, r...))
	hash := b64(h.Sum([]byte{}))
	return &Node{
		"",
		hash,
		left,
		right,
		nil,
		1 + min(left.MinDepth, right.MinDepth),
	}
}

func LeafNode(id string, signature string, data []byte) *Node {
	return &Node{
		id,
		signature,
		nil,
		nil,
		data,
		0,
	}
}

func (n *Node) IsLeaf() bool {
	return n.MinDepth == 0
}

func (t *Tree) GetIDs() []string {
	ids := make([]string, 0, len(t.traversion_lookup))
	for id := range t.traversion_lookup {
		ids = append(ids, id)
	}
	return ids
}

func (t *Tree) ReadFile(id string) (string, []byte, []HashValidation, error) {
	leaf, hashes := t.traverse(id)
	if leaf == nil {
		return "", nil, nil, errors.FileNotFound()
	}

	return leaf.Hash, leaf.FileData, hashes, nil
}

func (t *Tree) WriteFile(id string, sig string, data []byte) ([]HashValidation, error) {
	leaf, hashes := t.traverse(id)
	if leaf == nil {
		leaf, hashes = t.createFile(id, sig, data)
	}

	return hashes, nil
}

func (t *Tree) createFile(id string, sig string, data []byte) (*Node, []HashValidation) {
	var node **Node = &t.root
	var hashes []HashValidation

	for !(*node).IsLeaf() {
		if (*node).Left.MinDepth <= (*node).Right.MinDepth {
			hashes = append(hashes, HashValidation{(*node).Right.Hash, Right})
			(*node).MinDepth = 1 + min((*node).Left.MinDepth+1, (*node).Right.MinDepth)
			node = &(*node).Left
		} else {
			hashes = append(hashes, HashValidation{(*node).Left.Hash, Left})
			(*node).MinDepth = 1 + min((*node).Left.MinDepth, (*node).Right.MinDepth+1)
			node = &(*node).Right
		}
	}

	hashes = append(hashes, HashValidation{(*node).Hash, Left})

	newLeaf := LeafNode(id, sig, data)

	branch := BranchNode(newLeaf, *node)

	t.traversion_lookup[newLeaf.ID] = make([]bool, len(t.traversion_lookup[(*node).ID]))
	copy(t.traversion_lookup[newLeaf.ID], t.traversion_lookup[(*node).ID])
	t.traversion_lookup[newLeaf.ID] = append(t.traversion_lookup[newLeaf.ID], Left)

	t.traversion_lookup[(*node).ID] = append(t.traversion_lookup[(*node).ID], Right)

	*node = branch

	return newLeaf, hashes
}

func (t *Tree) traverse(id string) (*Node, []HashValidation) {
	traversion, ok := t.traversion_lookup[id]
	if !ok {
		return nil, nil
	}
	node := t.root
	var hashes []HashValidation
	for _, dir := range traversion {
		if node.IsLeaf() {
			break
		} else {
			if dir == Right {
				hashes = append(hashes, HashValidation{node.Left.Hash, Left})
				node = node.Right
			} else {
				hashes = append(hashes, HashValidation{node.Right.Hash, Right})
				node = node.Left
			}
		}
	}
	return node, hashes
}
