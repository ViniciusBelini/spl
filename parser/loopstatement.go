package parser

import(
	// "fmt"

	"SPL/models"
	"SPL/ast"
)

// Loop Statement Parser
func (p *Parser) LoopStatementParser(fileName string){
	tok := p.peek()

	switch tok.Value{
		case "while":
			pTemp := p
			tempAST := p.WhileStatement(fileName)
			if len(tempAST) > 0{
				p.Ast = append(p.Ast, tempAST[0])
			}
			p = pTemp
		default:
			p.unexpected(fileName)
	}
}

// While Statement
func (p *Parser) WhileStatement(fileName string) []ast.LoopStatement{
	tok := p.peek()

	var loopAst []ast.LoopStatement
	startIn := p.In
	var loopExpr []models.Token
	//inlineExpr := false
	if tok.Type == models.TokenLoopStatement && tok.Value == "while"{
		p.next()
		for !p.eof(){
			tok = p.peek()

			if tok.Type == models.TokenNewLine || tok.Value == ":"{
				if tok.Value == ":"{
					//inlineExpr := true
				}

				break
			}

			loopExpr = append(loopExpr, tok)
			p.next()
		}

		if len(loopExpr) == 0{
			p.In = startIn
			p.generic("Missing condition in 'while' statement", "S1004", fileName) // Error
		}
	}else{
		p.unexpected(fileName) // Error
	}
	p.next()

	//loopBlock := p.GetBlock(fileName)

	return loopAst
}
