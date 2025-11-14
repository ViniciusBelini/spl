package interpreter

import(
	// "fmt"
	"errors"
	// "runtime"

	"SPL/models"
	"SPL/ast"
)

func ListForm(node []ast.ArrayOneItem, outer *Env, fileName string) (interface{}, error){
	listType := 0
	var listA []any
	var listB map[any]any
	listB = make(map[any]any)
	for i := 0;i < len(node);i++{
		list := node[i]

		if listType == 0{
			if list.Right == nil{
				listType = 1
			}else{
				listType = 2
			}
		}

		if listType == 1{
			if list.Right != nil{
				return nil, errors.New(TRunMakeError(14, "array", "null", "map", fileName, list.Line, list.Pos))
			}

			responseA, err := Run([]ast.Node{list.Left}, outer, fileName, false)
			if err != nil{
				return nil, err
			}

			if arr, ok := responseA.([2]any);ok{
				responseA = arr[0]
			}

			listA = append(listA, responseA)
		}else if listType == 2{
			if list.Right == nil{
				return nil, errors.New(TRunMakeError(14, "map", "null", "array", fileName, list.Line, list.Pos))
			}

			responseA, err := Run([]ast.Node{list.Left}, outer, fileName, false)
			if err != nil{
				return nil, err
			}

			if arr, ok := responseA.([2]any);ok{
				responseA = arr[0]
			}

			_, TypeOfA := GetTypeData(responseA)
			if TypeOfA != models.TokenString && TypeOfA != models.TokenNumber && TypeOfA != models.TokenFloat{
				return nil, errors.New(TRunMakeError(15, TypeOfA, "null", "array", fileName, list.Line, list.Pos)) // Reminder!!!! Change this error
			}

			responseB, err := Run([]ast.Node{list.Right}, outer, fileName, false)
			if err != nil{
				return nil, err
			}

			if arr, ok := responseB.([2]any);ok{
				responseB = arr[0]
			}

			listB[responseA] = responseB
		}
	}

	if listType == 1{
		return listA, nil
	}else if listType == 2{
		return listB, nil
	}

	return nil, nil
}

func convertSplData(data interface{}, tabs string) (string, error){
	typeId, _ := GetTypeData(data)

	if typeId != 6 && typeId != 7 && typeId != 8{
		return ConvToString(data), nil
	}

	result := "{\n"
	i := 0
	if typeId == 6{
		for key, value := range data.(map[any]any){
			i++

			valueResult, err := convertSplData(value, tabs+"\t")
			if err != nil{
				return "null", err
			}

			result += tabs+"\t"+ConvToString(key)+": "+valueResult
			if i < len(data.(map[any]any)){
				result += ",\n"
			}
		}
	}else if typeId == 7{
		for i := 0;i < len(data.([]any));i++{
			resolveV, err := convertSplData(data.([]any)[i], tabs+"\t")
			if err != nil{
				return "null", err
			}

			result += tabs+"\t"+resolveV
			if i+1 < len(data.([]any)){
				result += ",\n"
			}
		}
	}
	result += "\n"+tabs+"}"

	return result, nil
}


func GetArr(baseE []ast.Node, keyE []ast.Node, line int, pos int, outer *Env, fileName string) (interface{}, error){
	value, err := Run(baseE, outer, fileName, false)
	if err != nil{
		return nil, err
	}

	if arr, ok := value.([2]any);ok{
		value = arr[0]
	}

	key, err := Run(keyE, outer, fileName, false)
	if err != nil{
		return nil, err
	}

	if arr, ok := key.([2]any);ok{
		key = arr[0]
	}

	idType, idTypeName := GetTypeData(value)
	idTypeKey, _ := GetTypeData(key)

	var resultE interface{}
	if idType == 6{
		if m, ok := value.(map[any]any); ok{
			if v, exists := m[key]; exists{
				resultE = v
			}else{
				return nil, errors.New(TRunMakeError(17, ConvToString(key), "null", "null", fileName, line, pos))
			}
		}
	}else if idType == 7{
		if idTypeKey == 1{
			if s, ok := value.([]any); ok{
				if key.(int) >= 0 && key.(int) < len(s){
					resultE = s[key.(int)]
				}else{
					return nil, errors.New(TRunMakeError(19, ConvToString(key), "null", "null", fileName, line, pos))
				}
			}
		}
	}else if idType == 0{
		if idTypeKey == 1{
			if s, ok := value.(string);ok{
				strRune := []rune(s)
				if key.(int) >= 0 && key.(int) < len(strRune){
					return string(strRune[key.(int)]), nil
				}
			}
			return nil, errors.New(TRunMakeError(19, ConvToString(key), "null", "null", fileName, line, pos))
		}
	}

	if resultE == nil{
		return nil, errors.New(TRunMakeError(18, idTypeName, "null", "null", fileName, line, pos))
	}

	cnv, err := convertSplData(resultE, "")
	if err != nil{
		return nil, err
	}

	resultE = [2]any{resultE, cnv}
	return resultE, nil
}
func SetArrValue(baseE []ast.Node, keyE []ast.Node, newValue interface{}, line int, pos int, outer *Env, fileName string) (interface{}, error){
	value, err := Run(baseE, outer, fileName, false)
	if err != nil{
		return nil, err
	}

	if arr, ok := value.([2]any);ok{
		value = arr[0]
	}

	key, err := Run(keyE, outer, fileName, false)
	if err != nil{
		return nil, err
	}

	if arr, ok := key.([2]any);ok{
		key = arr[0]
	}

	idType, _ := GetTypeData(value)
	idTypeKey, nameTypeKey := GetTypeData(key)

	if idType == 6{
		if _, ok := value.(map[any]any); ok{
			// if _, exists := m[key]; exists{
			if vKey, ok := key.(string);ok{
				value.(map[any]any)[vKey] = newValue
			}else if vKey, ok := key.(int);ok{
				value.(map[any]any)[vKey] = newValue
			}else if vKey, ok := key.(float64);ok{
				value.(map[any]any)[vKey] = newValue
			}else{
				return nil, errors.New(TRunMakeError(15, nameTypeKey, "null", "null", fileName, line, pos))
			}
			// }else{
				// value.(map[any]any)[key] = newValue // double code LOL
				// return nil, errors.New(TRunMakeError(17, ConvToString(key), "null", "null", fileName, line, pos))
			// }
		}
	}else if idType == 7{
		if idTypeKey == 1{
			if s, ok := value.([]any); ok{
				if key.(int) >= 0 && key.(int) < len(s){
					value.([]any)[key.(int)] = newValue
				}else{
					return nil, errors.New(TRunMakeError(19, ConvToString(key), "null", "null", fileName, line, pos))
				}
			}
		}
	}

	cnv, err := convertSplData(newValue, "")
	if err != nil{
		return nil, err
	}

	newValue = [2]any{newValue, cnv}
	return newValue, nil
}
