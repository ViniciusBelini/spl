package parser

import(
	// "fmt"
	"strings"

	"SPL/lexer"
	"SPL/config"
	"SPL/models"
	"SPL/ast"
)

func (p *Parser) FuncStatement(fileName string) []ast.FuncStatement{
	var FuncAST []ast.FuncStatement

	if p.eof(){
		return FuncAST
	}

	tok := p.peek()

	if tok.Type != models.TokenFuncStatement && tok.Value != "function"{
		p.unexpected(fileName) // Error
		return FuncAST
	}

	funLinePos := map[string]int{
		"line": tok.Line,
		"pos": tok.Pos,
	}

	if !p.canNext(){
		p.generic("Function declaration missing name and parameter list", "S1014", fileName) // Error
		return FuncAST
	}

	p.next();tok = p.peek()

	funcName := "__NULL_NAME__"
	var newParams []ast.ParamNode
	params := ""
	if tok.Type == models.TokenCall{
		openParen := strings.Index(tok.Value, "(")
		closeParen := strings.LastIndex(tok.Value, ")")

		if openParen == -1 || closeParen == -1 || closeParen <= openParen{
			p.unexpected(fileName) // Error
			return FuncAST
		}

		funcName = tok.Value[:openParen]
		params = tok.Value[openParen : closeParen+1]
		params = params[1 : len(params)-1]
	}else if tok.Type == models.TokenParentheses{
		if len(tok.Value) > 2{
			params = tok.Value[1:len(tok.Value)-1]
		}
	}else{
		p.generic("Function declaration missing parameter list", "S1014", fileName) // Error
		return FuncAST
	}

	funcParams := lexer.Tokenize(params, fileName, funLinePos["line"], funLinePos["pos"])
	for i := 0;i < len(funcParams);i++{
		paramToken := funcParams[i]

		if paramToken.Type == models.TokenIdent{
			typeParam := "dynamic"
			if i+1 < len(funcParams) && funcParams[i+1].Type == models.TokenType{
				typeParam = funcParams[i+1].Value
				i += 2
			}else if config.Config["mode"] == "strict"{
				p.generic("Function parameter must have an explicit type in strict mode", "S1012", fileName) // Error
				return FuncAST
			}

			newParams = append(newParams, ast.ParamNode{Name: paramToken.Value, Type: typeParam, Line: funLinePos["line"], Pos: funLinePos["pos"]})
		}else if paramToken.Type == models.TokenDelimiter && paramToken.Value == ","{
			continue
		}else{
			p.generic("Invalid parameter '"+paramToken.Value+"' for function", "S1013", fileName) // Error
			return FuncAST
		}
	}

	if !p.canNext(){
		p.unexpected(fileName) // Error
		return FuncAST
	}

	p.next();tok = p.peek()

	funcType := "dynamic"
	if tok.Type == models.TokenType{
		funcType = tok.Value

		if !p.canNext(){
			p.expected("function body", fileName)
			return FuncAST
		}
		p.next();tok = p.peek()
	}else if tok.Type != models.TokenType && config.Config["mode"] == "strict"{
		p.generic("Function must declare a return type in strict mode", "1011", fileName) // Error
		return FuncAST
	}

	if tok.Type != models.TokenNewLine && !(tok.Type == models.TokenDelimiter && tok.Value == ":"){
		p.unexpected(fileName) // Error
		return FuncAST
	}

	methodName := "function"
	if funcName != "__NULL_NAME__"{
		methodName = "'"+funcName+"' function"
	}
	funcBlock := p.GetBlock(fileName, methodName)

	// getFirst := func(nodes []ast.Node, returnFr bool) ast.Node{
	// 	if len(nodes) > 0{
	// 		if returnFr{
	// 			return nodes[0]
	// 		}
	// 		return nodes
	// 	}
	// 	return nil
	// }

	FuncAST = append(FuncAST, ast.FuncStatement{
		Name: funcName,
		Param: newParams,
		Type: funcType,
		Consequent: Astnize(funcBlock, fileName, "FuncStatement", false),
		Line: funLinePos["line"],
		Pos: funLinePos["pos"],
	})

	p.next()

	return FuncAST
}

func (p *Parser) FuncCall(fileName string) []ast.FuncCall{
	var callAst []ast.FuncCall

	if p.eof() || p.peek().Type != models.TokenCall{
		return callAst
	}

	tok := p.peek()
	tokInit := tok

	openParen := strings.Index(tok.Value, "(")
	closeParen := strings.LastIndex(tok.Value, ")")

	if openParen == -1 || closeParen == -1 || closeParen <= openParen{
		p.unexpected(fileName) // Error
		return callAst
	}

	funcName := tok.Value[:openParen]
	params := tok.Value[openParen : closeParen+1]
	params = params[1 : len(params)-1]

	funcParams := lexer.Tokenize(params, fileName, tokInit.Line, tokInit.Pos)
	var paramTokens []models.Token
	var newParams []ast.Node
	for i := 0;i < len(funcParams);i++{
		paramIn := funcParams[i]

		if !(paramIn.Type == models.TokenDelimiter && paramIn.Value == ","){
			paramTokens = append(paramTokens, paramIn)
		}

		if i+1 == len(funcParams) || paramIn.Type == models.TokenDelimiter && paramIn.Value == ","{
			varValueVerify := Astnize(paramTokens, fileName, funcName, true)
			var varValue ast.Node
			if len(varValueVerify) == 0{
				varValue = ast.NullNode{Line: tokInit.Line, Pos: tokInit.Pos,}
			}else{
				varValue = varValueVerify[0]
			}
			newParams = append(newParams, varValue)
			paramTokens = paramTokens[:0]

			if paramIn.Type == models.TokenDelimiter && paramIn.Value == ","{
				continue
			}
			break
		}
	}

	callAst = append(callAst, ast.FuncCall{Name: funcName, Param: newParams, Line: tokInit.Line, Pos: tokInit.Pos})

	p.next()
	return callAst
}

func (p *Parser) ObjCall(fileName string) []ast.ObjCall{
	var callAst []ast.ObjCall

	if p.eof() || p.peek().Type != models.TokenObj{
		return callAst
	}

	tok := p.peek()

	obj1, obj2 := splitFirstObject(tok.Value)

	obj1Token := lexer.Tokenize(obj1, fileName, tok.Line, tok.Pos)
	obj1Ast := Astnize(obj1Token, fileName, "object", true)

	obj2Token := lexer.Tokenize(obj2, fileName, tok.Line, tok.Pos)
	obj2Ast := Astnize(obj2Token, fileName, "object", true)

	getFirst := func(nodes []ast.Node, returnFr bool) ast.Node{
		if len(nodes) > 0{
			if returnFr{
				return nodes[0]
			}
			return nodes
		}
		return nil
	}

	callAst = append(callAst, ast.ObjCall{
		Obj: getFirst(obj1Ast, true),
		Consequent: getFirst(obj2Ast, true),
		Line: tok.Line,
		Pos: tok.Pos,
	})

	p.next()

	return callAst
}
func splitFirstObject(expr string) (string, string){ // helper
	bracketStack := []rune{}

	for i, r := range expr{
		switch r{
			case '(', '[':
				bracketStack = append(bracketStack, r)
			case ')':
				if len(bracketStack) > 0 && bracketStack[len(bracketStack)-1] == '('{
					bracketStack = bracketStack[:len(bracketStack)-1]
				}
			case ']':
				if len(bracketStack) > 0 && bracketStack[len(bracketStack)-1] == '['{
					bracketStack = bracketStack[:len(bracketStack)-1]
				}
			case '.':
				if len(bracketStack) == 0{
					return expr[:i], expr[i+1:]
				}
		}
	}

	return expr, ""
}

func (p *Parser) ArrayAccess(fileName string) []ast.ArrayAccess{
	var callAst []ast.ArrayAccess

	if p.eof() || p.peek().Type != models.TokenArrayAccess{
		return callAst
	}

	tok := p.peek()
	tokInit := tok

	openParen := strings.Index(tok.Value, "[")
	closeParen := strings.LastIndex(tok.Value, "]")

	if openParen == -1 || closeParen == -1 || closeParen <= openParen{
		p.unexpected(fileName) // Error
		return callAst
	}

	accessName := tok.Value[:openParen]
	params := tok.Value[openParen : closeParen+1]

	accessKey := lexer.Tokenize(params, fileName, tokInit.Line, tokInit.Pos)
	accessAst := Astnize([]models.Token{accessKey[len(accessKey)-1]}, fileName, accessName, true)

	base := accessName
	if len(accessKey) > 1{
		for i := 0;i < len(accessKey);i++{
			if i == len(accessKey)-1{
				continue
			}
			base += accessKey[i].Value
		}
	}
	accessResult := Astnize(lexer.Tokenize(base, fileName, tokInit.Line, tokInit.Pos), fileName, accessName, true)

	callAst = append(callAst, ast.ArrayAccess{Key: accessAst, Base: accessResult, Line: tokInit.Line, Pos: tokInit.Pos})

	p.next()
	return callAst
}
