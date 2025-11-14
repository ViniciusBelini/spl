package interpreter

import(
	"bufio"
	"fmt"
	"os"
	"strings"
	"errors"
	"strconv"
	"math/rand"
	// "time"
	"encoding/json"
	// "reflect"

	"SPL/models"
	"SPL/ast"
	"SPL/modules"
)

// Init
func GetBuiltIn() map[string]func(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	var BuiltInFuncs = map[string]func(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
		"len":  BUILT_IN_len,
		"type_of": BUILT_IN_type_of,
		"append": BUILT_IN_append,
		"delete": BUILT_IN_delete,
		"has": BUILT_IN_has,

		"print": BUILT_IN_SYSTEM_print,
		"__SYSTEM__io_input": BUILT_IN_SYSTEM_io_input,
		"__SYSTEM__math_rand": BUILT_IN_SYSTEM_math_rand,

		"int": BUILT_IN_SYSTEM_int,
		"str": BUILT_IN_SYSTEM_str,
		"float": BUILT_IN_SYSTEM_float,

		"__SYSTEM__http_get": BUILT_IN_RequestHttpGet,

		"__SYSTEM__json_encode": BUILT_IN_json_encode,
		"__SYSTEM__json_decode": BUILT_IN_json_decode,

		"__SYSTEM__gui": BUILT_IN_gui,
	}

	return BuiltInFuncs
}

func BuiltInFuncsExists(name string) bool{
	BuiltInFuncs := GetBuiltIn()

	if _, ok := BuiltInFuncs[name]; ok{
		return true
	}
	return false
}
func BuiltInFuncsCall(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	if !BuiltInFuncsExists(node.Name){
		return nil, errors.New("Error: BUILT IN function not defined")
	}

	BuiltInFuncs := GetBuiltIn()
	return BuiltInFuncs[node.Name](node, outer, fileName)
}
func BuiltInFuncsVerifyType(name string, paramsFn []string, paramsCl []ast.Node, outer *Env, fileName string, line int, pos int) ([]interface{}, error){
	if len(paramsCl) < len(paramsFn) || len(paramsCl) > len(paramsFn){
		return nil, errors.New(TRunMakeError(8, name, strconv.Itoa(len(paramsFn)), strconv.Itoa(len(paramsCl)), fileName, line, pos))
	}

	var paramsOrg []interface{}
	for i := 0;i < len(paramsCl);i++{
		paramB, err := Run([]ast.Node{paramsCl[i]}, outer, fileName, false)
		if err != nil{
			return nil, err
		}

		param := paramB
		if arr, ok := paramB.([2]any);ok{
			param = arr[0]
		}

		_, typeParam := GetTypeData(param)

		if paramsFn[i] == "dynamic" || paramsFn[i] == typeParam{
			paramsOrg = append(paramsOrg, paramB)
		}else{
			return nil, errors.New(TRunMakeError(12, name, typeParam, paramsFn[i], fileName, line, pos))
		}
	}

	return paramsOrg, nil
}

// Built in functions
///////////////// len
func BUILT_IN_len(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("len", []string{models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	getLen := checkParams[0]

	if arr, ok := getLen.([2]any);ok{
		getLen = arr[0]
	}

	mayConv, err := convertToMapAnyAny(getLen)
	id, typeParam := GetTypeData(mayConv)
	switch id{
		case 0:
			strRune := []rune(getLen.(string))
			return len(strRune), nil
		case 6:
			return len(getLen.(map[any]any)), nil
		case 7:
			return len(getLen.([]any)), nil
		default:
			return nil, errors.New(TRunMakeError(11, typeParam, "len()", "null", fileName, node.Line, node.Pos))
	}

	return nil, errors.New(TRunMakeError(11, typeParam, "len()", "null", fileName, node.Line, node.Pos))
}

///////////// type of
func BUILT_IN_type_of(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("type_of", []string{models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	getLen := checkParams[0]

	if arr, ok := getLen.([2]any);ok{
		getLen = arr[0]
	}

	_, typeParam := GetTypeData(getLen)
	return typeParam, nil
}

///////////// append
func BUILT_IN_append(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("append", []string{models.TokenDynamic, models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	list := checkParams[0]

	if arr, ok := list.([2]any);ok{
		list = arr[0]
	}
	id, typeParam := GetTypeData(list)

	value := checkParams[1]

	if arr, ok := value.([2]any);ok{
		value = arr[0]
	}

	if id == 7{
		list = append(list.([]any), value)
	}else{
		return nil, errors.New(TRunMakeError(11, typeParam, "append()", "null", fileName, node.Line, node.Pos))
	}

	return list, nil
}

///////////// unset
func BUILT_IN_delete(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("delete", []string{models.TokenDynamic, models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	list := checkParams[0]

	if arr, ok := list.([2]any);ok{
		list = arr[0]
	}
	id, typeParam := GetTypeData(list)

	value := checkParams[1]

	if arr, ok := value.([2]any);ok{
		value = arr[0]
	}

	_, typeParamKey := GetTypeData(value)

	if id == 6{
		if vKey, ok := value.(string);ok{
			delete(list.(map[any]any), vKey)
		}else if vKey, ok := value.(int);ok{
			delete(list.(map[any]any), vKey)
		}else if vKey, ok := value.(float64);ok{
			delete(list.(map[any]any), vKey)
		}else{
			return nil, errors.New(TRunMakeError(15, typeParamKey, "null", "null", fileName, node.Line, node.Pos))
		}
	}else if id == 7{
		if vKey, ok := value.(int);ok{
			list = append(list.([]any)[:vKey], list.([]any)[vKey+1:]...)
		}else{
			return nil, errors.New(TRunMakeError(15, typeParamKey, "null", "null", fileName, node.Line, node.Pos))
		}
		return list, nil
	}else{
		return nil, errors.New(TRunMakeError(11, typeParam, "delete()", "null", fileName, node.Line, node.Pos))
	}

	return list, nil
}

///////////// has
func BUILT_IN_has(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("has", []string{models.TokenDynamic, models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	list := checkParams[0]

	if arr, ok := list.([2]any);ok{
		list = arr[0]
	}
	id, typeParam := GetTypeData(list)

	value := checkParams[1]

	if arr, ok := value.([2]any);ok{
		value = arr[0]
	}

	_, typeParamKey := GetTypeData(value)

	if id == 6{
		if vKey, ok := value.(string);ok{
			_, ok = list.(map[any]any)[vKey]
			if ok{
				return true, nil
			}else{
				return false, nil
			}
		}else if vKey, ok := value.(int);ok{
			_, ok = list.(map[any]any)[vKey]
			if ok{
				return true, nil
			}else{
				return false, nil
			}
		}else if vKey, ok := value.(float64);ok{
			_, ok = list.(map[any]any)[vKey]
			if ok{
				return true, nil
			}else{
				return false, nil
			}
		}else{
			return nil, errors.New(TRunMakeError(15, typeParamKey, "null", "null", fileName, node.Line, node.Pos))
		}
	}else if id == 7{
		if vKey, ok := value.(int);ok{
			if vKey >= 0 && vKey < len(list.([]any)){
				return true, nil
			}
			return false, nil
		}else{
			return nil, errors.New(TRunMakeError(15, typeParamKey, "null", "null", fileName, node.Line, node.Pos))
		}
	}else{
		return nil, errors.New(TRunMakeError(11, typeParam, "has()", "null", fileName, node.Line, node.Pos))
	}

	return false, nil
}

///////////// system io input
func BUILT_IN_SYSTEM_print(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("print", []string{models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return false, err
	}

	message := checkParams[0]

	if arr, ok := message.([2]any);ok{
		fmt.Printf("%v\n", arr[1])
	}else{
		fmt.Printf("%v\n", message)
	}

	return true, nil
}
func BUILT_IN_SYSTEM_io_input(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	return text, nil
}

///////////// system math rand
func BUILT_IN_SYSTEM_math_rand(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("math_rand", []string{models.TokenNumber, models.TokenNumber}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	min := checkParams[0].(int)
	max := checkParams[1].(int)
	number := rand.Intn(max-min+1) + min
	return number, nil
}
///////////// system int - str - float
func BUILT_IN_SYSTEM_int(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("int", []string{models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	str := checkParams[0]
	_, typeStr := GetTypeData(str)
	switch str.(type){
		case string:
			r, err := strconv.Atoi(str.(string))
			if err != nil{
				return nil, errors.New(TRunMakeError(13, "null", typeStr, models.TokenNumber, fileName, node.Line, node.Pos))
			}
			return r, nil
		case int:
			return str.(int), nil
		case float64:
			return int(str.(float64)), nil
		default:
			return nil, errors.New(TRunMakeError(13, "null", typeStr, models.TokenNumber, fileName, node.Line, node.Pos))
	}
}
func BUILT_IN_SYSTEM_str(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("str", []string{models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	str := checkParams[0]
	switch str.(type){
		case string:
			return str.(string), nil
		case int:
			return strconv.Itoa(str.(int)), nil
		case float64:
			return strconv.FormatFloat(str.(float64), 'f', 2, 64), nil
		default:
			_, typeStr := GetTypeData(str)
			return nil, errors.New(TRunMakeError(13, "null", typeStr, models.TokenString, fileName, node.Line, node.Pos))
	}
}
func BUILT_IN_SYSTEM_float(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("float", []string{models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	str := checkParams[0]
	_, typeStr := GetTypeData(str)
	switch str.(type){
		case string:
			r, err := strconv.ParseFloat(str.(string), 64)
			if err != nil{
				return nil, errors.New(TRunMakeError(13, "null", typeStr, models.TokenFloat, fileName, node.Line, node.Pos))
			}
			return r, nil
		case int:
			return float64(str.(int)), nil
		case float64:
			return str.(float64), nil
		default:
			return nil, errors.New(TRunMakeError(13, "null", typeStr, models.TokenFloat, fileName, node.Line, node.Pos))
	}
}

func BUILT_IN_RequestHttpGet(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("__SYSTEM__http_get", []string{models.TokenString, "map<str:str>", models.TokenString, models.TokenString}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	url := checkParams[0]
	if arr, ok := url.([2]any);ok{
		url = arr[0]
	}

	headers := checkParams[1]
	if arr, ok := headers.([2]any);ok{
		headers = arr[0]
	}

	method := checkParams[2]
	if arr, ok := method.([2]any);ok{
		method = arr[0]
	}

	jsonData := checkParams[3]
	if arr, ok := jsonData.([2]any);ok{
		jsonData = arr[0]
	}

	response, err := modules.HttpGet(url.(string), headers.(map[any]any), method.(string), jsonData.(string))
	if err != nil{
		return nil, err
	}

	return response.Body, nil
}

func BUILT_IN_json_encode(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("__SYSTEM__json_encode", []string{models.TokenDynamic}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	content := checkParams[0]
	if arr, ok := content.([2]any); ok {
		content = arr[0]
	}

	if m, ok := content.(map[any]any); ok{
		newMap := loopConvertJSONAcept(m)

		result, err := json.Marshal(newMap)
		if err != nil {
			return nil, err
		}
		return string(result), nil
	}

	result, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}

	return string(result), err
}
func loopConvertJSONAcept(data interface{}) interface{}{
	switch v := data.(type) {
	case map[any]any:
		newMap := make(map[string]interface{})
		for k, val := range v {
			key := fmt.Sprintf("%v", k)
			newMap[key] = loopConvertJSONAcept(val)
		}
		return newMap
	case []any:
		newSlice := make([]interface{}, len(v))
		for i, val := range v {
			newSlice[i] = loopConvertJSONAcept(val)
		}
		return newSlice
	case [2]any:
		newSlice := make([]interface{}, len(v))
		for i, val := range v {
			newSlice[i] = loopConvertJSONAcept(val)
		}
		return newSlice
	default:
		return v
	}
}
func BUILT_IN_json_decode(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("__SYSTEM__json_decode", []string{models.TokenString}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	content := checkParams[0]
	if arr, ok := content.([2]any); ok {
		content = arr[0]
	}

	var result interface{}
	err = json.Unmarshal([]byte(content.(string)), &result)
	if err != nil {
		return nil, errors.New("Invalid JSON format")
	}

	cnvtCase, err := convertToMapAnyAny(result)
	if err != nil{
		return nil, err
	}
	return cnvtCase, nil
}

func BUILT_IN_gui(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	var lParamList []string
	for i := 0;i < len(node.Param);i++{
		lParamList = append(lParamList, models.TokenDynamic)
	}

	checkParams, err := BuiltInFuncsVerifyType("__SYSTEM__gui", lParamList, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	result, err := modules.GUI_run(checkParams)
	if err != nil{
		return nil, err
	}

	return result, nil
}
