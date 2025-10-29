package interpreter

type Env struct{
	Return interface{}
	Variables map[string]*Vars
	GlobalVars map[string]*Vars
	Functions map[string]*Func
	Outer *Env
}

type Vars struct{
	Value interface{}
	Type string
}

type Func struct{
	Outer *Env
	Point int
}
