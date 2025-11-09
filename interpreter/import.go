package interpreter

import(
	// "fmt"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"SPL/models"
	"SPL/lexer"
	"SPL/parser"
	"SPL/ast"
)

func ImportFunc(node ast.ImportNode, outer *Env, fileName string) (*Env, error){ ////// TEMPORARY
	if len(node.Path) < 2{
		return nil, errors.New(MRunMakeError(1, "null", fileName, node.Line, node.Pos))
	}

	path := node.Path

	exeDir, _ := filepath.Split(os.Args[0])
	cwd, _ := os.Getwd()

	baseDir := cwd
	if !strings.Contains(exeDir, "go-build") && exeDir != ""{
		baseDir = exeDir
	}
	modulesPath := filepath.Join(baseDir, "modules")

	if ImportFileExists(modulesPath+"/"+path+".spl"){
		path = modulesPath+"/"+path+".spl"
	}else{
		if !node.String{
			path = filepath.Dir(fileName)+"/"+path+".spl"
		}
	}

	if !ImportFileExists(path){
		return nil, errors.New(MRunMakeError(2, path, fileName, node.Line, node.Pos))
	}

	allTokens, err := ImportReadFIleTokens(path, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}
	importNode := parser.Astnize(allTokens, fileName, "null", false)

	newEnv := NewEnv(outer)
	Run(importNode, newEnv, fileName, false)

	_, err = DefineGlobalVariable(node.As, newEnv, models.TokenModule, outer, fileName, node.Line, node.Pos)
	if err != nil{
		return nil, err
	}

	return newEnv, nil
}
func ImportFileExists(fileName string) bool{
	if _, err := os.Stat(fileName); err != nil{
		if os.IsNotExist(err){
			return false
		}else{
			return false
		}
	}
	return true
}
func ImportReadFIleTokens(fileName string, line int, pos int) ([]models.Token, error){
	var allTokens []models.Token
	data, err := os.ReadFile(fileName)
	if err != nil{
		return nil, errors.New(MRunMakeError(1, "null", fileName, line, pos))
	}
	allTokens = lexer.Tokenize(string(data), fileName, 1, 1)

	return allTokens, nil
}
