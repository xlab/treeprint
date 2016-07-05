// Package treeprint provides a simple ASCII tree composing tool.
package treeprint

import (
	"bytes"
	"fmt"
	"io"
)

type Value interface{}
type MetaValue interface{}

type Tree interface {
	AddNode(v Value) Tree
	AddMetaNode(meta MetaValue, v Value) Tree
	AddBranch(v Value) Tree
	AddMetaBranch(meta MetaValue, v Value) Tree

	String() string
	Bytes() []byte
}

type node struct {
	Branch *node
	Meta   MetaValue
	Value  Value
	Nodes  []*node
}

func (n *node) AddNode(v Value) Tree {
	n.Nodes = append(n.Nodes, &node{
		Branch: n,
		Value:  v,
	})
	if n.Branch != nil {
		return n.Branch
	}
	return n
}

func (n *node) AddMetaNode(meta MetaValue, v Value) Tree {
	n.Nodes = append(n.Nodes, &node{
		Branch: n,
		Meta:   meta,
		Value:  v,
	})
	if n.Branch != nil {
		return n.Branch
	}
	return n
}

func (n *node) AddBranch(v Value) Tree {
	branch := &node{
		Value: v,
	}
	n.Nodes = append(n.Nodes, branch)
	return branch
}

func (n *node) AddMetaBranch(meta MetaValue, v Value) Tree {
	branch := &node{
		Meta:  meta,
		Value: v,
	}
	n.Nodes = append(n.Nodes, branch)
	return branch
}

func (n *node) Bytes() []byte {
	buf := new(bytes.Buffer)
	level := 0
	levelEnded := make(map[int]bool)
	if n.Branch == nil {
		buf.WriteString(string(EdgeTypeStart))
		buf.WriteByte('\n')
	} else {
		edge := EdgeTypeMid
		if len(n.Nodes) == 0 {
			edge = EdgeTypeEnd
			levelEnded[level] = true
		}
		printValues(buf, 0, levelEnded, edge, n.Meta, n.Value)
	}
	if len(n.Nodes) > 0 {
		printNodes(buf, level, levelEnded, n.Nodes)
	}
	return buf.Bytes()
}

func (n *node) String() string {
	return string(n.Bytes())
}

func printNodes(wr io.Writer,
	level int, levelEnded map[int]bool, nodes []*node) {

	for i, node := range nodes {
		edge := EdgeTypeMid
		if i == len(nodes)-1 {
			levelEnded[level] = true
			edge = EdgeTypeEnd
		}
		printValues(wr, level, levelEnded, edge, node.Meta, node.Value)
		if len(node.Nodes) > 0 {
			printNodes(wr, level+1, levelEnded, node.Nodes)
		}
	}
}

func printValues(wr io.Writer,
	level int, levelEnded map[int]bool, edge EdgeType, meta MetaValue, v Value) {

	for i := 0; i < level; i++ {
		if levelEnded[i] {
			fmt.Fprint(wr, "    ")
			continue
		}
		fmt.Fprintf(wr, "%s   ", EdgeTypeLink)
	}
	if meta != nil {
		fmt.Fprintf(wr, "%s [%v]  %v\n", edge, meta, v)
		return
	}
	fmt.Fprintf(wr, "%s %v\n", edge, v)
}

type EdgeType string

const (
	EdgeTypeStart EdgeType = "."
	EdgeTypeLink  EdgeType = "│"
	EdgeTypeMid   EdgeType = "├──"
	EdgeTypeEnd   EdgeType = "└──"
)

func New() Tree {
	return &node{}
}
