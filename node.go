package main

type Node struct {
	id      int64
	Lexical string
}

func (node *Node) ID() int64 {
	return node.id
}
