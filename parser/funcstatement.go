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

	if tok.Type != models.TokenCall{
		p.generic("Function declaration missing name and parameter list", "S1014", fileName) // Error
		return FuncAST
	}

	openParen := strings.Index(tok.Value, "(")
	closeParen := strings.LastIndex(tok.Value, ")")

	if openParen == -1 || closeParen == -1 || closeParen <= openParen{
		p.unexpected(fileName) // Error
		return FuncAST
	}

	funcName := tok.Value[:openParen]
	params := tok.Value[openParen : closeParen+1]
	params = params[1 : len(params)-1]

	funcParams := lexer.Tokenize(params, fileName, funLinePos["line"], funLinePos["pos"])
	var newParams []ast.ParamNode
	for i := 0;i < len(funcParams);i++{
		paramToken := funcParams[i]

		if paramToken.Type == models.TokenIdent{
			typeParam := "dynamic"
			if i+1 < len(funcParams) && funcParams[i+1].Type == models.TokenType{
				typeParam = funcParams[i+1].Value
				i += 2
			}else{
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

	// funcType := "dynamic"
	if tok.Type == models.TokenType{
		// funcType = tok.Value

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

	funcBlock := p.GetBlock(fileName, "function")

	getFirst := func(nodes []ast.Node, returnFr bool) ast.Node{
		if len(nodes) > 0{
			if returnFr{
				return nodes[0]
			}
			return nodes
		}
		return nil
	}

	FuncAST = append(FuncAST, ast.FuncStatement{
		Name: funcName,
		Param: newParams,
		Consequent: getFirst(Astnize(funcBlock, fileName, "FuncStatement", false).([]ast.Node), false),
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
			varValueVerify := Astnize(paramTokens, fileName, funcName, true).([]ast.Node)
			var varValue ast.Node
			if len(varValueVerify) == 0{
				varValue = ast.NullNode{Line: tokInit.Line, Pos: tokInit.Pos,}
			}else{
				varValue = varValueVerify[0]
			}
			newParams = append(newParams, varValue)
			paramTokens = paramTokens[:0]

			break
		}
	}

	callAst = append(callAst, ast.FuncCall{Name: funcName, Param: newParams, Line: tokInit.Line, Pos: tokInit.Pos})

	p.next()
	return callAst
}
