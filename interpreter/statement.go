package interpreter

import(
	// "fmt"

	"SPL/config"
	"SPL/models"
	"SPL/ast"
)

func IfStatement(node ast.IfStatement, outer *Env, fileName string) (interface{}, error){
	env := outer
	if config.Config["mode"] == "strict"{
		env = NewEnv(outer)
	}

	value, err := Run([]ast.Node{node.Test}, env, fileName, false)
	if err != nil{
		return nil, err
	}

	_, valueType := GetTypeData(value)

	if value == true && valueType == models.TokenBoolean{
		if node.Consequent != nil{
			_, err := Run(node.Consequent.([]ast.Node), env, fileName, false)
			if err != nil{
				return nil, err
			}
		}

		return value, nil
	}else if node.Alternate != nil{
		return Run(node.Alternate.([]ast.Node), env, fileName, false)
	}

	return value, nil
}
