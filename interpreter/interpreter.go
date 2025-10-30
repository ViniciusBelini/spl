package interpreter

import(
	"fmt"
	"reflect"
	"strconv"
	"errors"

	"SPL/config"
	"SPL/models"
	"SPL/ast"
)

// Main interpreter routine
func Run(aAst []ast.Node, outer *Env, fileName string, newEnvS bool) (interface{}, error){
	env := outer
	if newEnvS{
		env = NewEnv(outer)
	}

	// First scan
	for i := 0;i < len(aAst);i++{
		node := aAst[i]
		switch node.(type){
			case ast.FuncStatement:
				// fmt.Printf("%#v\n", node)
				// fmt.Println(node.(ast.FuncStatement).Name)

				env.Functions[node.(ast.FuncStatement).Name] = &Func{
					Outer: env,
					Point: i,
				}
		}
	}

	// Second scan
	for i := 0;i < len(aAst);i++{
		node := aAst[i]
		switch node.(type){
			case ast.NullNode:
				env.Return = nil
			case ast.NativeSugarNode:
				if node.(ast.NativeSugarNode).Name == "show"{
					value, err := Run([]ast.Node{node.(ast.NativeSugarNode).Value}, env, fileName, false)
					if err != nil{
						return nil, err
					}
					fmt.Println(value)
					env.Return = true
				}
			case ast.IdentNode:
				varData, err := GetVariable(node.(ast.IdentNode).Name, env, fileName, node.(ast.IdentNode).Line, node.(ast.IdentNode).Pos)
				if err != nil{
					return nil, err
				}

				env.Return = varData.Value
			case ast.AssignNode:
				value, err := AssignVariable(node.(ast.AssignNode), env, fileName)
				if err != nil{
					return value, err
				}
			case ast.UnaryOpNode:
				value, err := Run([]ast.Node{node.(ast.UnaryOpNode).Right}, env, fileName, false)
				typeValue, typeName := GetTypeData(value)
				if err != nil{
					return value, err
				}

				if node.(ast.UnaryOpNode).Operator == "!"{
					if typeValue == 3 || config.Config["mode"] == "dynamic"{
						RetunUnary, err := UnaryOpConv(value)
						if err != nil{
							return RetunUnary, errors.New(TGRunMakeError(1, err.Error(), fileName, node.(ast.UnaryOpNode).Line, node.(ast.UnaryOpNode).Pos))
						}
						env.Return = RetunUnary
					}else{
						return nil, errors.New(TRunMakeError(2, "null", "null", "null", fileName, node.(ast.UnaryOpNode).Line, node.(ast.UnaryOpNode).Pos))
					}
				}else if node.(ast.UnaryOpNode).Operator == "+" || node.(ast.UnaryOpNode).Operator == "-"{
					varName := node.(ast.UnaryOpNode).Right.(ast.IdentNode).Name

					if typeValue == 1 || typeValue == 2{
						varNValue, ers := MathOp(value, 1, node.(ast.UnaryOpNode).Operator)
						if ers != nil{
							return nil, ers
						}

						_, err := SetVariable(varName, varNValue, env, fileName, node.(ast.UnaryOpNode).Line, node.(ast.UnaryOpNode).Pos)
						if err != nil{
							return nil, err
						}
					}else{
						return nil, errors.New(TRunMakeError(5, varName, "null", typeName, fileName, node.(ast.UnaryOpNode).Line, node.(ast.UnaryOpNode).Pos))
					}
				}
			case ast.BinaryOpNode:
				left, err := Run([]ast.Node{node.(ast.BinaryOpNode).Left}, env, fileName, false)
				if err != nil{
					return left, err
				}

				right, err := Run([]ast.Node{node.(ast.BinaryOpNode).Right}, env, fileName, false)
				if err != nil{
					return right, err
				}

				switch node.(ast.BinaryOpNode).Operator{
					case "+", "-", "/", "*", "%":
						if VerifyTypeData(left, right) || (TypeDataNumber(left, right) && config.Config["mode"] == "dynamic"){
							ReturnMath, err := MathOp(left, right, node.(ast.BinaryOpNode).Operator)
							if err != nil{
								return ReturnMath, errors.New(TGRunMakeError(1, err.Error(), fileName, node.(ast.BinaryOpNode).Line, node.(ast.BinaryOpNode).Pos))
							}
							env.Return = ReturnMath
						}else{
							if config.Config["mode"] == "strict"{
								_, typeLeft := GetTypeData(left)
								_, typeRight := GetTypeData(right)
								return nil, errors.New(TRunMakeError(1, typeLeft, typeRight, node.(ast.BinaryOpNode).Operator, fileName, node.(ast.BinaryOpNode).Line, node.(ast.BinaryOpNode).Pos))
							}else{
								env.Return = MathJoin(left, right)
							}
						}
					case "==", "!=", ">=", "<=", ">", "<", "||", "&&":
						if VerifyTypeData(left, right) || (TypeDataNumber(left, right) && config.Config["mode"] == "dynamic"){
							ReturnLogic, err := CompareOp(left, right, node.(ast.BinaryOpNode).Operator)
							if err != nil{
								return ReturnLogic, errors.New(TGRunMakeError(1, err.Error(), fileName, node.(ast.BinaryOpNode).Line, node.(ast.BinaryOpNode).Pos))
							}
							env.Return = ReturnLogic
						}else{
							if config.Config["mode"] == "strict"{
								_, typeLeft := GetTypeData(left)
								_, typeRight := GetTypeData(right)
								return nil, errors.New(TRunMakeError(1, typeLeft, typeRight, node.(ast.BinaryOpNode).Operator, fileName, node.(ast.BinaryOpNode).Line, node.(ast.BinaryOpNode).Pos))
							} // Error - maybe
						}
				}
			case ast.LiteralNode:
				tmpValue := node.(ast.LiteralNode).Value
				tmpType := node.(ast.LiteralNode).Type

				if tmpType == models.TokenString{
					env.Return = tmpValue.(string)[1 : len(tmpValue.(string))-1]
				}else if tmpType == models.TokenNumber{
					if v, err := strconv.Atoi(tmpValue.(string));err == nil{
						env.Return = v
					} // Error - maybe
				}else if tmpType == models.TokenFloat{
					if v, err := strconv.ParseFloat(tmpValue.(string), 64);err == nil{
						env.Return = v
					} // Error - maybe
				}else if tmpType == models.TokenBoolean{
					if v, err := strconv.ParseBool(tmpValue.(string));err == nil{
						env.Return = v
					} // Error - maybe
				}

				continue
			case ast.FuncStatement:
				// Just skip
				continue
			default:
				fmt.Println("Here we are:", reflect.TypeOf(node)) // Remove this XXXXXXXXXXXXX
		}
	}

	return env.Return, nil
}

func NewEnv(outer *Env) *Env{
	return &Env{
		Return: nil,
		Variables: make(map[string]*Vars),
		GlobalVars: make(map[string]*Vars),
		Functions: make(map[string]*Func),
		Outer:     outer,
	}
}

func VerifyTypeData(x interface{}, y interface{})bool{
	a, _ := GetTypeData(x)
	b, _ := GetTypeData(y)

	return a == b
}
func TypeDataNumber(x interface{}, y interface{}) bool{
	a, _ := GetTypeData(x)
	b, _ := GetTypeData(y)

	return (a == 1 || a == 2) && (b == 1 || b == 2)
}
func GetTypeData(x interface{})(int, string){
	switch x.(type){
		case string:
			return 0, models.TokenString
		case int:
			return 1, models.TokenNumber
		case float64:
			return 2, models.TokenFloat
		case bool:
			return 3, models.TokenBoolean
		default:
			return -1, models.TokenUnknown
	}
}
