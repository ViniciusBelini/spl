package parser

import(
	// "SPL/config"
	"SPL/models"
	"SPL/ast"
)

// Unery expression starter
func (p *Parser) UnExpr(fileName string) []ast.UnaryOpNode{
	tok := p.peek()

	var UnExprAST []ast.UnaryOpNode
	if tok.Type != models.TokenUnOp{
		p.unexpected(fileName)
		return UnExprAST
	}

	operatorTok := tok
	if tok.Value == "!"{
		if !p.canNext(){
			p.unexpected(fileName)
			return UnExprAST
		}

		p.next()
		tok = p.peek()

		getFirst := func(nodes []ast.Node, returnFr bool) ast.Node{
			if len(nodes) > 0{
				if returnFr{
					return nodes[0]
				}
				return nodes
			}
			return UnExprAST
		}

		var nTok []models.Token
		nTok = append(nTok, tok)
		UnExprAST = append(UnExprAST, ast.UnaryOpNode{Right: getFirst(Astnize(nTok, fileName, "IfStatement", true).([]ast.Node), true), Operator: "!", Line: operatorTok.Line, Pos: operatorTok.Pos})

		p.next()

		return UnExprAST
	}
	return UnExprAST
}
