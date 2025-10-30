package interpreter

import(
	// "fmt"
	"errors"

	"SPL/config"
	"SPL/ast"
)

// Assign variable
func AssignVariable(node ast.AssignNode, outer *Env, fileName string) (interface{}, error){
	if node.Method == ":=" || (node.Method == "=" && config.Config["mode"] == "dynamic"){
		value, err := Run([]ast.Node{node.Value}, outer, fileName, false)
		if err != nil{
			return value, err
		}

		if node.Type == "dynamic" && config.Config["mode"] == "strict"{
			return nil, errors.New(TRunMakeError(4, node.Name, "null", "null", fileName, node.Line, node.Pos))
		}


		newVar, err := DefineVariable(node.Name, value, node.Type, outer, fileName, node.Line, node.Pos)
		_, err2 := GetVariable(node.Name, outer, fileName, node.Line, node.Pos)
		if err != nil || err2 == nil{
			if config.Config["mode"] == "dynamic"{
				_, errs := SetVariable(node.Name, value, outer, fileName, node.Line, node.Pos)
				if errs != nil{
					return nil, errs
				}
				return value, nil
			}
			return nil, err
		}
		return newVar, err
	}else if node.Method == "=" || node.Method == "+=" || node.Method == "-="{
		value, err := Run([]ast.Node{node.Value}, outer, fileName, false)
		if err != nil{
			return value, err
		}

		_, err = SetVariable(node.Name, value, outer, fileName, node.Line, node.Pos)
		if err != nil{
			return nil, err
		}
		return value, nil
	}

	return nil, nil
}

func GetVariable(name string, outer *Env, fileName string, line int, pos int) (*Vars, error){
	if varVal, exists := outer.Variables[name];exists{
		return varVal, nil
	}else{
		if varVal, exists := outer.GlobalVars[name];exists{
			return varVal, nil
		}
		return nil, errors.New(NRunMakeError(1, name, fileName, line, pos))
	}
}

func DefineVariable(name string, value interface{}, vType string, outer *Env, fileName string, line int, pos int) (*Vars, error){
	_, err := GetVariable(name, outer, fileName, line, pos)
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

func SetVariable(name string, value interface{}, outer *Env, fileName string, line int, pos int) (*Vars, error){
	varVal, err := GetVariable(name, outer, fileName, line, pos)
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
