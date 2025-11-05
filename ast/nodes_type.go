package ast

type Node interface{}

type ImportNode struct{
	Path		string
	As		string
	String		bool
	Line		int
	Pos		int
}

type AssignNode struct{
	Name		string
	Type		string
	Value		Node
	Method		string
	Line		int
	Pos		int
}

type IfStatement struct{
	Test		Node
	Consequent	Node
	Alternate	Node
	Line		int
	Pos		int
}

type LoopStatement struct{
	Method		string
	Init		Node
	Test		Node
	Update		Node
	Consequent	Node
	Line		int
	Pos		int
}

type FuncStatement struct{
	Name		string
	Param		[]ParamNode
	Type		string
	Consequent	[]Node
	Line		int
	Pos		int
}

type ParamNode struct{
	Name		string
	Type		string
	Line		int
	Pos		int
}

type FuncCall struct{
	Name		string
	Param		[]Node
	Line		int
	Pos		int
}

type ObjCall struct{
	Obj		Node
	Consequent	Node
	Line		int
	Pos		int
}

type IdentNode struct{
	Name		string
	Line		int
	Pos		int
}

type NativeSugarNode struct{
	Name		string
	Value		Node
	Line		int
	Pos		int
}

type LiteralNode struct{
	Value		interface{}
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

type UnaryOpNode struct{
	Right		Node
	Operator	string
	Line		int
	Pos		int
}

type ControlFlowNode struct{
	Method		string
	Argument	Node
	Line		int
	Pos		int
}

type NullNode struct{
	Line		int
	Pos		int
}

type Debug struct{
	Value		string
}
