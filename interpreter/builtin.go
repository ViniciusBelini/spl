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
			return nil, errors.New(TRunMakeError(3, name, typeParam, paramsFn[i], fileName, line, pos))
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
	checkParams, err := BuiltInFuncsVerifyType("type_of", []string{models.TokenNumber, models.TokenNumber}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	min := checkParams[0].(int)
	max := checkParams[1].(int)
	number := rand.Intn(max-min+1) + min
	return number, nil
}
///////////// system math int
func BUILT_IN_SYSTEM_int(node ast.FuncCall, outer *Env, fileName string) (interface{}, error){
	checkParams, err := BuiltInFuncsVerifyType("type_of", []string{models.TokenString}, node.Param, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	str := checkParams[0].(string)
	number, _ := strconv.Atoi(str)
	return number, nil
}
