package types

import (
	"crypto/sha256"
	"hash"
)

type NodeValue interface {
	CalculateHash() ([]byte, error)
	Equals(other NodeValue) (bool, error)
}

type MerkleTree[T NodeValue] struct {
	Root         *Node[T]
	merkleRoot   []byte
	Leafs        []*Node[T]
	hashStrategy func() hash.Hash
}

type Node[T NodeValue] struct {
	Tree   *MerkleTree[T]
	Parent *Node[T]
	Left   *Node[T]
	Right  *Node[T]
	level  int64
	Hash   []byte
	Value  *T
}

func (n *Node[T]) leaf() bool {
	return n.level == 0
}

func (n *Node[T]) verifyNode() ([]byte, error) {
	if n.leaf() {
		return (*n.Value).CalculateHash()
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

func (n *Node[T]) calculateNodeHash() ([]byte, error) {
	if n.leaf() {
		return (*n.Value).CalculateHash()
	}

	h := n.Tree.hashStrategy()
	if _, err := h.Write(append(n.Left.Hash, n.Right.Hash...)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func NewMerkleTree[T NodeValue]() *MerkleTree[T] {
	merkleTree := &MerkleTree[T]{
		merkleRoot:   nil,
		Leafs:        make([]*Node[T], 0, 64),
		hashStrategy: sha256.New,
	}

	merkleTree.Root = &Node[T]{
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

func (m *MerkleTree[T]) Append(value T) (int, error) {
	root := m.Root

	valueHash, err := value.CalculateHash()
	var previousHash []byte

	if len(m.Leafs) > 0 {
		previousHash = m.Leafs[len(m.Leafs)-1].Hash
	} else {
		previousHash = make([]byte, 0)
	}

	if err != nil {
		return 0, err
	}

	h := m.hashStrategy()

	_, err = h.Write(valueHash)
	valueHash = h.Sum(previousHash)

	if err != nil {
		return 0, err
	}

	node := &Node[T]{
		Tree:  m,
		level: 0,
		Hash:  valueHash,
		Value: &value,
	}

	m.Leafs = append(m.Leafs, node)

	if root.Left == nil {
		root.Left = node
		node.Parent = root
		return len(m.Leafs), m.updateParentHashes(node)
	}

	if root.Right == nil {
		root.Right = node
		node.Parent = root
		return len(m.Leafs), m.updateParentHashes(node)
	}

	success := appendToSubroot(node, root)

	if !success {
		newRoot := &Node[T]{
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
	return len(m.Leafs), m.updateParentHashes(node)
}

func appendToSubroot[T NodeValue](node *Node[T], subroot *Node[T]) bool {
	if subroot.leaf() {
		return false
	}

	if subroot.level-subroot.Right.level == 1 {
		return appendToSubroot(node, subroot.Right)
	}

	child := subroot.Right
	parent := &Node[T]{
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

func (m *MerkleTree[T]) updateParentHashes(node *Node[T]) error {
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
