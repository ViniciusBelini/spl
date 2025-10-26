package parser

import(
	// "fmt"
	// "strings"

	// "SPL/lexer"
	// "SPL/config"
	"SPL/models"
	"SPL/ast"
)

// Native no params functions - sugar functions
func (p *Parser) NativeSugar(fileName string){
	tok := p.peek()

	switch tok.Value{
		case "show":
			pTemp := p
			tempAST := p.ShowSugar(fileName)
			if len(tempAST) > 0{
				p.Ast = append(p.Ast, tempAST[0])
				return
			}
			p = pTemp
			p.unexpected(fileName)
		default:
			p.unexpected(fileName)
	}

	return
}

// Show sugar
func (p *Parser) ShowSugar(fileName string) []ast.NativeSugarNode{
	var astSugar []ast.NativeSugarNode

	tok := p.peek()
	tokInit := tok

	if tok.Value != "show"{
		return astSugar
	}

	p.next()
	var sugarTokens []models.Token
	for !p.eof(){
		tok = p.peek()

		if tok.Type == models.TokenNewLine{
			break
		}

		sugarTokens = append(sugarTokens, tok)
		p.next()
	}

	getFirst := func(nodes []ast.Node, returnFr bool) ast.Node{
		if len(nodes) > 0{
			if returnFr{
				return nodes[0]
			}
			return nodes
		}
		return nil
	}

	astSugar = append(astSugar, ast.NativeSugarNode{
		Name: "show",
		Value: getFirst(Astnize(sugarTokens, fileName, "FuncStatement", true).([]ast.Node), true),
		Line: tokInit.Line,
		Pos: tokInit.Pos,
	})

	p.next()
	return astSugar
}
