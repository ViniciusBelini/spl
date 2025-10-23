package parser

import(
	// "fmt"

	"SPL/models"
	"SPL/ast"
)

// Loop Statement Parser
func (p *Parser) LoopStatementParser(fileName string, statementExpr bool){
	tok := p.peek()

	switch tok.Value{
		case "while":
			pTemp := p
			tempAST := p.WhileStatement(fileName)
			if len(tempAST) > 0 && !statementExpr{
				p.Ast = append(p.Ast, tempAST[0])
				return
			}
			p = pTemp
			p.unexpected(fileName)
		default:
			p.unexpected(fileName)
	}
}

// While Statement
func (p *Parser) WhileStatement(fileName string) []ast.LoopStatement{
	tok := p.peek()

	var loopAst []ast.LoopStatement

	startIn := p.In
	startLine := tok.Line
	startPos := tok.Pos

	var loopExpr []models.Token
	if tok.Type == models.TokenLoopStatement && tok.Value == "while"{
		p.next()
		for !p.eof(){
			tok = p.peek()

			if tok.Type == models.TokenNewLine || tok.Value == ":"{
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

	loopBlock := p.GetBlock(fileName, "while")

	if len(loopBlock) == 0{
		p.In = startIn
		p.generic("[SyntaxError] Missing 'end' of 'while' statement", "S1005", fileName) // Error
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

	loopAst = append(loopAst, ast.LoopStatement{
		Method:		"while",
		Init:		nil,
		Test:		getFirst(Astnize(loopExpr, fileName, "LoopStatement", true).([]ast.Node), true),
		Update:		nil,
		Consequent:	getFirst(Astnize(loopBlock, fileName, "LoopStatement", false).([]ast.Node), false),
		Line:		startLine,
		Pos:		startPos,
	})

	p.next()

	return loopAst
}
