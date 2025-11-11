package interpreter

import(
	"fmt"
	"reflect"
	"strconv"
	"errors"
	// "strings"
	// "encoding/json"

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
				funcStat := node.(ast.FuncStatement)
				err := AssignFunc(funcStat.Name, &funcStat, env, fileName, funcStat.Line, funcStat.Pos)
				if err != nil{
					return nil, err
				}
		}
	}

	// Second scan
	for i := 0;i < len(aAst);i++{
		node := aAst[i]
		switch node.(type){
			case ast.ArrayAccess:
				value, err := GetArr(node.(ast.ArrayAccess).Base, node.(ast.ArrayAccess).Key, node.(ast.ArrayAccess).Line, node.(ast.ArrayAccess).Pos, env, fileName)
				if err != nil{
					return nil, err
				}
				env.Return = value
			case []ast.ArrayOneItem:
				result, err := ListForm(node.([]ast.ArrayOneItem), env, fileName)
				if err != nil{
					return nil, err
				}

				splData, err := convertSplData(result, "")
				if err != nil{
					return nil, err
				}

				env.Return = [2]any{result, splData}
			case ast.ObjCall:
				obj1, err := Run([]ast.Node{node.(ast.ObjCall).Obj}, env, fileName, false)
				if err != nil{
					return nil, err
				}

				if arr, ok := obj1.([2]any);ok{
					obj1 = arr[0]
				}

				_, typeObj := GetTypeData(obj1)
				if typeObj == models.TokenModule{
					tmpScope := obj1.(*Env).Outer
					obj1.(*Env).Outer = env
					result, err := Run([]ast.Node{node.(ast.ObjCall).Consequent}, obj1.(*Env), fileName, false)
					obj1.(*Env).Outer = tmpScope
					if err != nil{
						return nil, err
					}

					if arr, ok := result.([2]any);ok{
						result = arr[0]
					}

					env.Return = result
				}
			case ast.ImportNode:
				result, err := ImportFunc(node.(ast.ImportNode), env, fileName)
				if err != nil{
					return nil, err
				}
				env.Return = result
			case ast.FuncCall:
				if BuiltInFuncsExists(node.(ast.FuncCall).Name){
					result, err := BuiltInFuncsCall(node.(ast.FuncCall), env, fileName)
					if err != nil{
						return nil, err
					}
					env.Return = result
				}else{
					result, err := CallFunc(node.(ast.FuncCall).Name, node.(ast.FuncCall).Param, env, fileName, node.(ast.FuncCall).Line, node.(ast.FuncCall).Pos)
					if err != nil{
						return nil, err
					}
					env.Return = result
				}
			case ast.ControlFlowNode:
				method := node.(ast.ControlFlowNode).Method
				argument, err := Run([]ast.Node{node.(ast.ControlFlowNode).Argument}, env, fileName, false)
				if err != nil{
					return nil, err
				}

				if arr, ok := argument.([2]any);ok{
					argument = arr[0]
				}

				return [2]any{method, argument}, nil
			case ast.IfStatement, ast.LoopStatement:
				switch node.(type){
					case ast.IfStatement:
						stm, err := IfStatement(node.(ast.IfStatement), env, fileName)
						if err != nil{
							return nil, err
						}
						env.Return = stm

						if _, ok := stm.([2]any);ok{
							return stm, nil
						}
					case ast.LoopStatement:
						stm, err := LoopStatement(node.(ast.LoopStatement), env, fileName)
						if err != nil{
							return nil, err
						}
						env.Return = stm

						if _, ok := stm.([2]any);ok{
							return stm, nil
						}
				}
			case ast.NullNode:
				env.Return = [2]any{nil, "null"}
			case ast.NativeSugarNode:
				if node.(ast.NativeSugarNode).Name == "show"{
					value, err := Run([]ast.Node{node.(ast.NativeSugarNode).Value}, env, fileName, false)
					if err != nil{
						return nil, err
					}

					if arr, ok := value.([2]any);ok{
						fmt.Printf("%v", arr[1])
					}else{
						fmt.Printf("%v", value)
					}
					env.Return = true
				}
			case ast.IdentNode:
				varData, err := GetVariable(node.(ast.IdentNode).Name, env, fileName, node.(ast.IdentNode).Line, node.(ast.IdentNode).Pos, false)
				if err != nil{
					return nil, err
				}

				id, typeData := GetTypeData(varData.Value)
				if typeData == models.TokenModule{
					env.Return = [2]any{varData.Value, models.TokenModule+"("+node.(ast.IdentNode).Name+")"}
				}else if id == 6 || id == 7{
					splData, err := convertSplData(varData.Value, "")
					if err != nil{
						return nil, err
					}

					env.Return = [2]any{varData.Value, splData}
				}else{
					env.Return = varData.Value
				}
			case ast.AssignNode:
				value, err := AssignVariable(node.(ast.AssignNode), env, fileName)
				if err != nil{
					return value, err
				}
				env.Return = value
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
				if arr, ok := left.([2]any);ok{
					left = arr[0]
				}

				right, err := Run([]ast.Node{node.(ast.BinaryOpNode).Right}, env, fileName, false)
				if err != nil{
					return right, err
				}
				if arr, ok := right.([2]any);ok{
					right = arr[0]
				}

				switch node.(ast.BinaryOpNode).Operator{
					case "..":
						if 1 == 1 || TypeDataString(left, right) || config.Config["mode"] == "dynamic"{
							env.Return = MathJoin(left, right)
						}else{
							_, typeLeft := GetTypeData(left)
							_, typeRight := GetTypeData(right)

							return nil, errors.New(TRunMakeError(7, typeLeft, typeRight, node.(ast.BinaryOpNode).Operator, fileName, node.(ast.BinaryOpNode).Line, node.(ast.BinaryOpNode).Pos))
						}
					case "+", "-", "/", "*", "%":
						if VerifyTypeData(left, right) || (TypeDataNumber(left, right) && config.Config["mode"] == "dynamic"){
							if TypeDataString(left, right) && config.Config["mode"] == "strict"{
								return nil, errors.New(TRunMakeError(6, "null", "null", "null", fileName, node.(ast.BinaryOpNode).Line, node.(ast.BinaryOpNode).Pos))
							}

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
					// env.Return = strings.ReplaceAll(tmpValue.(string)[1 : len(tmpValue.(string))-1], "\b\\n", "\n")
					realText, err := strconv.Unquote(tmpValue.(string))
					if err != nil{
						return nil, err
					}

					env.Return = realText
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
			case ast.FuncStatement, nil:
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
		GlobalAccess: false,
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
func TypeDataString(x interface{}, y interface{}) bool{
	a, _ := GetTypeData(x)
	b, _ := GetTypeData(y)

	return a == 0 && b == 0
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
		case *Env:
			return 4, models.TokenModule
		case nil:
			return 5, models.TokenNull
		case map[any]any:
			keyType := ""
			valueType := ""
			for key, value := range x.(map[any]any){
				_, getKeyT := GetTypeData(key)
				_, getValueT := GetTypeData(value)

				if keyType == "" && valueType == ""{
					keyType = getKeyT
					valueType = getValueT
				}

				if keyType != getKeyT || valueType != getValueT{
					keyType = "dynamic"
					valueType = "dynamic"
				}
			}

			if keyType != "dynamic"{
				keyType = keyType[1: len(keyType)-1]
			}
			if valueType == models.TokenString || valueType == models.TokenNumber || valueType == models.TokenFloat || valueType == models.TokenBoolean{
				valueType = valueType[1: len(valueType)-1]
			}

			return 6, "map<"+keyType+":"+valueType+">"
		case []any:
			valueType := ""
			for i := 0;i < len(x.([]any));i++{
				value := x.([]any)[i]

				_, getValueT := GetTypeData(value)

				if valueType == ""{
					valueType = getValueT
				}

				if valueType != getValueT{
					valueType = "dynamic"
				}
			}

			if valueType[0:1] == "<" && valueType[len(valueType)-1:len(valueType)] == ">"{
				valueType = valueType[1:len(valueType)-1]
			}

			return 7, "array<"+valueType+">"
		default:
			// fmt.Println(reflect.TypeOf(x))
			return -1, models.TokenUnknown
	}
}
