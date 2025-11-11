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
					if arr[0].(string) == "continue"{
						continue
					}else if arr[0].(string) == "break"{
						break
					}else if arr[0].(string) == "return"{
						return arr, nil
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
func AssignFunc(name string, point *ast.FuncStatement, outer *Env, fileName string, line int, pos int) error{
	_, err := GetFunc(name, outer, fileName, line, pos)
	if err == nil{
		return errors.New(NRunMakeError(2, name, fileName, line, pos))
	}

	outer.Functions[name] = &Func{
		Outer: outer,
		Point: point,
	}

	return nil
}
func GetFunc(name string, outer *Env, fileName string, line int, pos int) (*ast.FuncStatement, error){
	if varVal, exists := outer.Functions[name];exists{
		return varVal.Point, nil
	}else{
		if outer.Outer != nil{
			varVal, err := GetFunc(name, outer.Outer, fileName, line, pos)
			if err == nil{
				return varVal, nil
			}
		}
		return nil, errors.New(NRunMakeError(1, name, fileName, line, pos))
	}
}
func CallFunc(name string, params []ast.Node, outer *Env, fileName string, line int, pos int) (interface{}, error){
	funcP, err := GetFunc(name, outer, fileName, line, pos)
	if err != nil{
		return nil, err
	}

	if len(params) < len(funcP.Param) || len(params) > len(funcP.Param){
		return nil, errors.New(TRunMakeError(8, name, strconv.Itoa(len(funcP.Param)), strconv.Itoa(len(params)), fileName, line, pos))
	}

	env := NewEnv(outer)
	env.GlobalVars = outer.GlobalVars
	env.GlobalAccess = false

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
