package parser

import(
	// "fmt"

	// "SPL/config"
	"SPL/models"
	"SPL/ast"
)

// Unery expression starter
func (p *Parser) UnExpr(fileName string) []ast.UnaryOpNode{
	tok := p.peek()

	var UnExprAST []ast.UnaryOpNode

	operatorTok := tok
	valueTok := tok
	if tok.Type == models.TokenUnOp{
		if !p.canNext(){
			p.unexpected(fileName)
			return UnExprAST
		}
		valueTok = p.peekNext()
	}else if p.canNext() && p.peekNext().Type == models.TokenUnOp{
		operatorTok = p.peekNext()
	}else{
		p.unexpected(fileName)
		return UnExprAST
	}


	if operatorTok.Value == "!" || operatorTok.Value == "++" || operatorTok.Value == "--"{
		if !p.canNext(){
			p.unexpected(fileName)
			return UnExprAST
		}

		p.next()
		getFirst := func(nodes []ast.Node, returnFr bool) ast.Node{
			if len(nodes) > 0{
				if returnFr{
					return nodes[0]
				}
				return nodes
			}
			return UnExprAST
		}

		UnOpResult := operatorTok.Value
		if operatorTok.Value == "++"{
			UnOpResult = "+"
		}else if operatorTok.Value == "--"{
			UnOpResult = "+"
		}

		var nTok []models.Token
		nTok = append(nTok, valueTok)
		UnExprAST = append(UnExprAST, ast.UnaryOpNode{Right: getFirst(Astnize(nTok, fileName, "IfStatement", true), true), Operator: UnOpResult, Line: operatorTok.Line, Pos: operatorTok.Pos})

		p.next()

		return UnExprAST
	}
	return UnExprAST
}
