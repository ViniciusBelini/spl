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
		p.unexpected(fileName) // Error
		return FuncAST
	}

	p.next();tok = p.peek()

	if tok.Type != models.TokenCall{
		p.unexpected(fileName) // Error
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
		p.generic("[SyntaxError] Function must declare a return type in strict mode", "1011", fileName) // Error
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
		Param: lexer.Tokenize(params, fileName, funLinePos["line"], funLinePos["pos"]),
		Consequent: getFirst(Astnize(funcBlock, fileName, "FuncStatement", false).([]ast.Node), false),
		Line: funLinePos["line"],
		Pos: funLinePos["pos"],
	})

	p.next()

	return FuncAST
}
