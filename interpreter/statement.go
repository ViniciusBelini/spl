package interpreter

import(
	// "fmt"
	"errors"
	"strconv"
	// "reflect"

	"SPL/config"
	"SPL/models"
	"SPL/ast"
)

// If statement
func IfStatement(node ast.IfStatement, outer *Env, fileName string) (interface{}, error){
	env := outer
	if config.Config["mode"] == "strict"{
		env = NewEnv(outer)
		env.GlobalAccess = true
		env.GlobalVars = outer.GlobalVars
	}

	value, err := Run([]ast.Node{node.Test}, env, fileName, false)
	if err != nil{
		return nil, err
	}

	_, valueType := GetTypeData(value)

	if value == true && valueType == models.TokenBoolean{
		if node.Consequent != nil{
			code, err := Run(node.Consequent.([]ast.Node), env, fileName, false)
			if err != nil{
				return nil, err
			}
			return code, nil
		}

		return value, nil
	}else if node.Alternate != nil{
		return Run(node.Alternate.([]ast.Node), env, fileName, false)
	}

	return value, nil
}

// Loop statement
func LoopStatement(node ast.LoopStatement, outer *Env, fileName string) (interface{}, error){
	for true{
		env := outer
		if config.Config["mode"] == "strict"{
			env = NewEnv(outer)
			env.GlobalAccess = true
			env.GlobalVars = outer.GlobalVars
		}

		value, err := Run([]ast.Node{node.Test}, env, fileName, false)
		if err != nil{
			return nil, err
		}

		_, valueType := GetTypeData(value)
		if value == true && valueType == models.TokenBoolean{
			if node.Consequent != nil{
				result, err := Run(node.Consequent.([]ast.Node), env, fileName, false)
				if err != nil{
					return nil, err
				}

				if arr, ok := result.([2]any);ok{
					if v, ok := arr[0].(string);ok{
						if v == "continue"{
							continue
						}else if v == "break"{
							break
						}else if v == "return"{
							return arr, nil
						}
					}
				}

				continue
			}

			return value, nil
		}

		return value, nil
	}

	return false, nil
}

// Func statement
func AssignFunc(name string, point *ast.FuncStatement, outer *Env, fileName string, line int, pos int) (*Func, error){
	_, err := GetFunc(name, outer, fileName, line, pos)
	if err == nil{
		return nil, errors.New(NRunMakeError(2, name, fileName, line, pos))
	}

	funcP := &Func{
		Outer: outer,
		Point: point,
		FileName: fileName,
	}

	if name != "__NULL_NAME__"{
		_, err = DefineGlobalVariable(name, funcP, models.TokenFunction, outer, fileName, line, pos)
		if err != nil{
			return nil, err
		}
	}

	return funcP, nil
}
func GetFunc(name string, outer *Env, fileName string, line int, pos int) (*ast.FuncStatement, error){
	varVal, err := GetVariable(name, outer, fileName, line, pos, false)
	if err != nil{
		if outer.Outer != nil{
			varVal, err := GetFunc(name, outer.Outer, fileName, line, pos)
			if err == nil{
				return varVal, nil
			}
		}
		return nil, errors.New(NRunMakeError(1, name, fileName, line, pos))
	}

	if f, ok := varVal.Value.(*Func);ok{
		return f.Point, nil
	}
	return nil, errors.New(TRunMakeError(20, name, "null", "null", fileName, line, pos))
}
func CallFunc(name string, nameP *Func, params []ast.Node, outer *Env, fileName string, line int, pos int) (interface{}, error){
	var funcP *ast.FuncStatement

	if nameP == nil{
		funcPtmp, err := GetFunc(name, outer, fileName, line, pos)
		if err != nil{
			return nil, err
		}
		funcP = funcPtmp
	}else{
		funcP = nameP.Point
	}

	if len(params) < len(funcP.Param) || len(params) > len(funcP.Param){
		return nil, errors.New(TRunMakeError(8, name, strconv.Itoa(len(funcP.Param)), strconv.Itoa(len(params)), fileName, line, pos))
	}

	env := NewEnv(outer)
	env.GlobalVars = outer.GlobalVars
	env.GlobalAccess = false
	if nameP != nil{
		env.GlobalAccess = true
	}

	for i := 0;i < len(funcP.Param);i++{
		funcParam := funcP.Param[i]
		callParam := params[i]

		callResult, err := Run([]ast.Node{callParam}, outer, fileName, false)
		if err != nil{
			return nil, err
		}
		if arr, ok := callResult.([2]any);ok{
			callResult = arr[0]
		}

		_, typeParam := GetTypeData(callResult)

		if funcParam.Type == "dynamic" || typeParam == funcParam.Type{
			_, err := ForceDefineVariable(funcParam.Name, callResult, funcParam.Type, env, fileName, line, pos)
			if err != nil{
				return nil, err
			}
		}else{
			return nil, errors.New(TRunMakeError(12, funcParam.Name, typeParam, funcParam.Type, fileName, line, pos))
		}
	}

	if env.Outer != nil{
		varPath, err := GetVariable("__PATH__", env.Outer, fileName, line, pos, false)
		if err == nil{
			if varPath.Type == models.TokenString{
				fileName = varPath.Value.(string)
			}
		}
	}

	returnFunc, err := Run(funcP.Consequent, env, fileName, false)

	if err != nil{
		return nil, err
	}

	if arr, ok := returnFunc.([2]any);ok{
		if arr[0].(string) == "return"{
			_, typeResult := GetTypeData(arr[1])
			if funcP.Type == "dynamic" || funcP.Type == typeResult{
				return arr[1], nil
			}else{
				return nil, errors.New(TRunMakeError(9, funcP.Type, typeResult, "null", fileName, funcP.Line, funcP.Pos))
			}
		}
	}

	if funcP.Type == "dynamic" || config.Config["mode"] == "dynamic"{
		return true, nil
	}

	return nil, errors.New(TRunMakeError(10, funcP.Type, "null", "null", fileName, line, pos))
}
