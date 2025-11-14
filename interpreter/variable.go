package interpreter

import(
	// "fmt"
	// "runtime"
	"errors"
	// "reflect"

	"SPL/config"
	"SPL/ast"
)

// Assign variable
func AssignVariable(node ast.AssignNode, outer *Env, fileName string) (interface{}, error){
	if node.Method == ":=" || (node.Method == "=" && config.Config["mode"] == "dynamic" && node.NamePonter == nil){
		value, err := Run([]ast.Node{node.Value}, outer, fileName, false)
		if err != nil{
			return value, err
		}

		if arr, ok := value.([2]any);ok{
			value = arr[0]
		}

		namePointer, err := Run([]ast.Node{node.NamePonter}, outer, fileName, false)
		if err != nil{
			return nil, err
		}

		if arr, ok := namePointer.([2]any);ok{
			namePointer = arr[0]
		}

		if node.Type == "dynamic" && config.Config["mode"] == "strict"{
			return nil, errors.New(TRunMakeError(4, node.Name, "null", "null", fileName, node.Line, node.Pos))
		}

		if node.Type == "<dynamic>"{
			node.Type = "dynamic"
		}

		newVar, err := DefineVariable(node.Name, value, node.Type, outer, fileName, node.Line, node.Pos, true)
		_, err2 := GetVariable(node.Name, outer, fileName, node.Line, node.Pos, true)

		if err != nil && err2 != nil{
			return nil, err
		}

		if err != nil && err2 == nil{
			if config.Config["mode"] == "dynamic"{
				_, errs := SetVariable(node.Name, value, outer, fileName, node.Line, node.Pos)
				if errs != nil{
					return nil, errs
				}
				return value, nil
			}
		}

		return newVar, err
	}else if node.Method == "=" || node.Method == "+=" || node.Method == "-="{
		value, err := Run([]ast.Node{node.Value}, outer, fileName, false)
		if err != nil{
			return value, err
		}

		if arr, ok := value.([2]any);ok{
			value = arr[0]
		}

		if node.NamePonter != nil{
			switch node.NamePonter.(type){
				case ast.ArrayAccess:
					result, err := SetArrValue(node.NamePonter.(ast.ArrayAccess).Base, node.NamePonter.(ast.ArrayAccess).Key, value, node.NamePonter.(ast.ArrayAccess).Line, node.NamePonter.(ast.ArrayAccess).Pos, outer, fileName)
					if err != nil{
						return nil, err
					}

					return result, nil
				case ast.ObjCall:
					// fmt.Printf("%v", node.NamePonter.(ast.ObjCall).Obj)
					obj, err := Run([]ast.Node{node.NamePonter.(ast.ObjCall).Obj}, outer, fileName, false)
					if err != nil{
						return nil, err
					}

					if arr, ok := obj.([2]any);ok{
						obj = arr[0]
					}

					typeObj, _ := GetTypeData(obj)
					if typeObj == 4{
						_, err := Run([]ast.Node{
							ast.AssignNode{
								Name: node.NamePonter.(ast.ObjCall).Consequent.(ast.IdentNode).Name,
								NamePonter: nil,
								Type: node.Type,
								Value: node.Value,
								Method: node.Method,
								Line: node.Line,
								Pos: node.Pos,
							},
						}, obj.(*Env), fileName, false)

						if err != nil{
							return nil, err
						}
						return value, nil
					}/*else{
						_, err = SetVariable(node.NamePonter.(ast.ObjCall).Obj.(ast.IdentNode).Name, value, obj, fileName, node.Line, node.Pos)
						if err != nil{
							return nil, err
						}
						return value, nil
					}*/
			}
		}

		_, err = SetVariable(node.Name, value, outer, fileName, node.Line, node.Pos)
		if err != nil{
			return nil, err
		}
		return value, nil
	}

	return nil, nil
}

func GetVariable(name string, outer *Env, fileName string, line int, pos int, define bool) (*Vars, error){
	if varVal, exists := outer.Variables[name];exists{
		return varVal, nil
	}else{
		if varVal, exists := outer.GlobalVars[name];exists{
			return varVal, nil
		}else if !define{
			if outer.Outer != nil && outer.GlobalAccess{
				varVal, err := GetVariable(name, outer.Outer, fileName, line, pos, define)
				if err == nil{
					return varVal, nil
				}
			}
		}
		return nil, errors.New(NRunMakeError(1, name, fileName, line, pos))
	}
}

func DefineVariable(name string, value interface{}, vType string, outer *Env, fileName string, line int, pos int, define bool) (*Vars, error){
	_, err := GetVariable(name, outer, fileName, line, pos, true)
	if err == nil{
		return nil, errors.New(NRunMakeError(2, name, fileName, line, pos))
	}

	_, nType := GetTypeData(value)
	if vType != nType && vType != "dynamic"{
		return nil, errors.New(TRunMakeError(3, name, nType, vType, fileName, line, pos))
	}

	outer.Variables[name] = &Vars{
		Value: value,
		Type: vType,
	}

	return outer.Variables[name], nil
}
func ForceDefineVariable(name string, value interface{}, vType string, outer *Env, fileName string, line int, pos int) (*Vars, error){
	_, nType := GetTypeData(value)
	if vType != nType && vType != "dynamic"{
		return nil, errors.New(TRunMakeError(3, name, nType, vType, fileName, line, pos))
	}

	outer.Variables[name] = &Vars{
		Value: value,
		Type: vType,
	}

	return outer.Variables[name], nil
}
func DefineGlobalVariable(name string, value interface{}, vType string, outer *Env, fileName string, line int, pos int) (*Vars, error){
	_, err := GetVariable(name, outer, fileName, line, pos, true)
	if err == nil{
		return nil, errors.New(NRunMakeError(2, name, fileName, line, pos))
	}

	_, nType := GetTypeData(value)
	if vType != nType && vType != "dynamic"{
		return nil, errors.New(TRunMakeError(3, name, nType, vType, fileName, line, pos))
	}

	outer.GlobalVars[name] = &Vars{
		Value: value,
		Type: vType,
	}

	return outer.GlobalVars[name], nil
}

func SetVariable(name string, value interface{}, outer *Env, fileName string, line int, pos int) (*Vars, error){
	varVal, err := GetVariable(name, outer, fileName, line, pos, false)
	if err != nil{
		return nil, errors.New(NRunMakeError(1, name, fileName, line, pos))
	}

	_, nType := GetTypeData(value)
	if varVal.Type == nType || varVal.Type == "dynamic"{
		varVal.Value = value
	}else{
		return nil, errors.New(TRunMakeError(3, name, nType, varVal.Type, fileName, line, pos))
	}

	return varVal, nil
}
