package types

import (
	"crypto/sha256"
	"hash"
)

type NodeValue interface {
	CalculateHash() ([]byte, error)
	Equals(other NodeValue) (bool, error)
}

type MerkleTree struct {
	Root         *Node
	merkleRoot   []byte
	Leafs        []*Node
	hashStrategy func() hash.Hash
}

type Node struct {
	Tree   *MerkleTree
	Parent *Node
	Left   *Node
	Right  *Node
	level  int64
	Hash   []byte
	Value  NodeValue
}

func (n *Node) leaf() bool {
	return n.level == 0
}

func (n *Node) verifyNode() ([]byte, error) {
	if n.leaf() {
		return n.Value.CalculateHash()
	}
	rightBytes, err := n.Right.verifyNode()
	if err != nil {
		return nil, err
	}

	leftBytes, err := n.Left.verifyNode()
	if err != nil {
		return nil, err
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(leftBytes, rightBytes...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (n *Node) calculateNodeHash() ([]byte, error) {
	if n.leaf() {
		return n.Value.CalculateHash()
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(n.Left.Hash, n.Right.Hash...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func NewMerkleTree() *MerkleTree {
	merkleTree := &MerkleTree{
		merkleRoot:   nil,
		Leafs:        make([]*Node, 0, 64),
		hashStrategy: sha256.New,
	}

	merkleTree.Root = &Node{
		Tree:   merkleTree,
		Parent: nil,
		Left:   nil,
		Right:  nil,
		level:  1,
		Hash:   nil,
		Value:  nil,
	}
	return merkleTree
}

func (m *MerkleTree) Append(value NodeValue) error {
	root := m.Root

	hash, err := value.CalculateHash()

	if err != nil {
		return err
	}

	node := &Node{
		Tree:  m,
		level: 0,
		Hash:  hash,
		Value: value,
	}

	m.Leafs = append(m.Leafs, node)

	if root.Left == nil {
		root.Left = node
		node.Parent = root
		return m.updateParentHashes(node)
	}

	if root.Right == nil {
		root.Right = node
		node.Parent = root
		return m.updateParentHashes(node)
	}

	success := appendToSubroot(node, root)

	if !success {
		newRoot := &Node{
			Tree:   m,
			Parent: nil,
			Left:   root,
			Right:  node,
			level:  root.level + 1,
			Hash:   nil,
			Value:  nil,
		}

		m.Root = newRoot
		root.Parent = newRoot
		node.Parent = newRoot
	}
	return m.updateParentHashes(node)
}

func appendToSubroot(node *Node, subroot *Node) bool {
	if subroot.leaf() {
		return false
	}

	if subroot.level-subroot.Right.level == 1 {
		return appendToSubroot(node, subroot.Right)
	}

	child := subroot.Right
	parent := &Node{
		Tree:   subroot.Tree,
		Parent: subroot,
		Left:   child,
		Right:  node,
		level:  child.level + 1,
		Hash:   nil,
		Value:  nil,
	}

	subroot.Right = parent
	child.Parent = parent
	node.Parent = parent
	return true
}

func (m *MerkleTree) updateParentHashes(node *Node) error {
	if node == nil {
		return nil
	}

	parent := node.Parent

	if parent == nil {
		return nil
	}

	h := m.hashStrategy()

	leftBytes, err := parent.Left.calculateNodeHash()
	if err != nil {
		return err
	}

	if parent.Right != nil {
		rightBytes, err := parent.Right.calculateNodeHash()
		if err != nil {
			return err
		}

		_, err = h.Write(append(leftBytes, rightBytes...))
	} else {
		_, err = h.Write(append(leftBytes))
	}

	if err != nil {
		return err
	}

	parent.Hash = h.Sum(nil)
	return m.updateParentHashes(parent)
}
