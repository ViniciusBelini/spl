package ast

type Node interface{}

type AssignNode struct{
	Name		string
	Type		string
	Value		Node
	Line		int
	Pos		int
}

type LiteralNode struct{
	Value		int
	Line		int
	Pos		int
}

type Debug struct{
	Value		string
}
