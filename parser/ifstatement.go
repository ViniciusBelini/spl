package parser

import(
	// "fmt"

	"SPL/models"
	"SPL/ast"
)

// If Statement
func (p *Parser) IfStatement(fileName string) bool{
	tok := p.peek()

	if tok.Type != models.TokenIfStatement{
		p.unexpected(fileName)
	}

	startIn := p.In
	startLine := tok.Line
	startPos := tok.Pos

	exprType := tok.Value
	var ifExpr []models.Token
	p.next()
	if exprType == "if"{
		for !p.eof(){
			tok = p.peek()

			if tok.Type == models.TokenNewLine || tok.Value == ":"{
				break
			}

			ifExpr = append(ifExpr, tok)
			p.next()
		}

		if len(ifExpr) == 0{
			p.In = startIn
			p.generic("Missing condition in 'if' statement", "S1004", fileName) // Error
		}
	}
	p.next()

	loopGetBlock := func(full bool) []models.Token{
		var ifBlock []models.Token
		blockWithEnds := 1
		if full{
			blockWithEnds = 0
		}
		for !p.eof(){
			tok = p.peek()

			if hasEndDelimiter(tok.Value){
				blockWithEnds++
			}

			if tok.Type == models.TokenDelimiter && tok.Value == "end"{
				blockWithEnds--
				if blockWithEnds == 0{
					if full{
						ifBlock = append(ifBlock, tok)
					}
					break
				}
			}

			if !full{
				if blockWithEnds-1 <= 0 && tok.Type == models.TokenIfStatement && tok.Value == "else"{
					blockWithEnds = 0
					break
				}
			}

			ifBlock = append(ifBlock, tok)
			p.next()
		}

		if blockWithEnds != 0{
			if p.eof(){
				p.back()
			}
			p.generic("[SyntaxError] Missing 'end' of 'if' statement", "S1005", fileName) // Error
		}

		return ifBlock
	}

	if !p.eof(){
		tok = p.peek()
	}

	ifBlock := loopGetBlock(false)

	var ifAlternate []models.Token
	if tok.Type == models.TokenIfStatement && tok.Value == "else"{
		if p.canNext() && p.peekNext().Value == "if"{
			p.next()
			ifAlternate = loopGetBlock(true)
		}else if p.canNext(){
			p.next()
			ifAlternate = loopGetBlock(false)
		}else{
			p.unexpected(fileName)
		}
	}

	// tok = p.peek()

	if tok.Type != models.TokenDelimiter && tok.Value != "end"{
		if p.eof(){
			p.back()
		}
		p.generic("[SyntaxError] Missing 'end' of 'if' statement", "S1005", fileName) // Error
	}

	ifAst := ast.IfStatement{
		Test: Astnize(ifExpr, fileName, "IfStatement").([]ast.Node)[0],
		Consequent: Astnize(ifBlock, fileName, "IfStatement").([]ast.Node)[0],
		Alternate: Astnize(ifAlternate, fileName, "IfStatement").([]ast.Node)[0],
		Line: startLine,
		Pos: startPos,
	}

	p.Ast = append(p.Ast, ifAst)

	p.next()
	return true
}
