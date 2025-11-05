package interpreter

import(
	"SPL/ast"
)

type Env struct{
	Return interface{}
	RealReturn interface{}
	Variables map[string]*Vars
	GlobalVars map[string]*Vars
	Functions map[string]*Func
	GlobalAccess bool
	Outer *Env
}

type Vars struct{
	Value interface{}
	Type string
}

type Func struct{
	Outer *Env
	Point *ast.FuncStatement
}
