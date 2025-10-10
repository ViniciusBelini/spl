package ast

type Node interface{}

type AssignNode struct{
	Name		string
	Type		string
	Value		Node
	Line		int
	Pos		int
}

type IdentNode struct{
	Name		string
	Line		int
	Pos		int
}

type StringLiteralNode struct{
	Value		string
	Line		int
	Pos		int
}

type Debug struct{
	Value		string
}
