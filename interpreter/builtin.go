package interpreter

import(
	"bufio"
	// "fmt"
	"os"
	"strings"
	"errors"
	"strconv"
	"math/rand"
	// "time"

	"SPL/models"
	"SPL/ast"
)

// Init
func GetBuiltIn() map[string]func(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	var BuiltInFuncs = map[string]func(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
		"len":  BUILT_IN_len,
		"type_of": BUILT_IN_type_of,

		"__SYSTEM__io_input": BUILT_IN_SYSTEM_io_input,
		"__SYSTEM__math_rand": BUILT_IN_SYSTEM_math_rand,

		"int": BUILT_IN_SYSTEM_int,
		"str": BUILT_IN_SYSTEM_str,
		"float": BUILT_IN_SYSTEM_float,
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
		param, err := Run([]ast.Node{paramsCl[i]}, outer, fileName, false)
		if err != nil{
			return nil, err
		}
		_, typeParam := GetTypeData(param)

		if paramsFn[i] == "dynamic" || paramsFn[i] == typeParam{
			paramsOrg = append(paramsOrg, param)
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
	_, typeParam := GetTypeData(getLen)
	switch getLen.(type){
		case string:
			return len(getLen.(string)), nil
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

///////////// system io input
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
