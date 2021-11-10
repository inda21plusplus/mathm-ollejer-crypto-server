package merkel

import (
	"fmt"

	"github.com/inda21plusplus/mathm-ollejer-crypto-server/server/errors"
)

type Hash []byte

type Tree struct {
	root   Node
	lookup map[string]uint64
}

func (t *Tree) traverseTo(id []byte) (leaf *LeafNode, hashes []Hash, err error) {
	branches, ok := t.lookup[string(id)]
	if !ok {
		err = errors.FileNotFound()
		return
	}
	node := t.root
	for i := 0; i < 64; i++ {
		goRight := (branches>>i)&1 == 1
		if n, ok := node.(*BranchNode); ok {
			if goRight {
				hashes = append(hashes, n.Left.GetHash())
				node = n.Right
			} else {
				hashes = append(hashes, n.Right.GetHash())
				node = n.Left
			}
		} else {
			break
		}
	}
	leaf, ok = node.(*LeafNode)
	if !ok {
		// This should never happen
		fmt.Println("Tree depth above 2^64")
		err = errors.FileNotFound()
		return
	}

	return
}

func (t *Tree) ReadFile(id []byte) (data []byte, sig Hash, hashes []Hash, err error) {
	var leaf *LeafNode
	leaf, hashes, err = t.traverseTo(id)
	if err != nil {
		return
	}

	data = leaf.FileData
	sig = leaf.Sig
	return
}

func (t *Tree) WriteFile(id []byte, data []byte, sig Hash) (hashes []Hash, err error) {
	var leaf *LeafNode
	leaf, hashes, err = t.traverseTo(id)
	if err != nil {
		return
	}

	leaf.FileData = data
	leaf.Sig = sig

	return
}

type Node interface {
	GetHash() Hash
}

type BranchNode struct {
	Hash  Hash
	Left  Node
	Right Node
}

func (node *BranchNode) GetHash() Hash {
	return node.Hash
}

type LeafNode struct {
	Sig      Hash
	FileData []byte
}

func (node *LeafNode) GetHash() Hash {
	return node.Sig
}
