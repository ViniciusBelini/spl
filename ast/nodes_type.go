package ast

type Node interface{}

type AssignNode struct{
	Name		string
	Type		string
	Value		Node
	Method		string
	Line		int
	Pos		int
}

type IdentNode struct{
	Name		string
	Line		int
	Pos		int
}

type LiteralNode struct{
	Value		string
	Type		string
	Line		int
	Pos		int
}

type BinaryOpNode struct{
	Left		Node
	Right		Node
	Operator	string
	Line		int
	Pos		int
}

type NullNode struct{
	Args		string
	Line		int
	Pos		int
}

type Debug struct{
	Value		string
}
