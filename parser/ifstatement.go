package parser

import(
	// "fmt"

	"SPL/models"
	"SPL/ast"
)

// If Statement
func (p *Parser) IfStatement(fileName string) []ast.IfStatement{
	tok := p.peek()

	if tok.Type != models.TokenIfStatement{
		p.unexpected(fileName)
	}

	var ifAst []ast.IfStatement

	startIn := p.In
	startLine := tok.Line
	startPos := tok.Pos

	inlineExpr := false
	exprType := tok.Value
	var ifExpr []models.Token
	p.next()
	if exprType == "if"{
		for !p.eof(){
			tok = p.peek()

			if tok.Type == models.TokenNewLine || tok.Value == ":"{
				if tok.Value == ":"{
					inlineExpr = true
				}

				break
			}

			ifExpr = append(ifExpr, tok)
			p.next()
		}

		if len(ifExpr) == 0{
			p.In = startIn
			p.generic("Missing condition in 'if' statement", "S1004", fileName) // Error
		}
	}else{
		p.unexpected(fileName) // Error
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

			if (tok.Type == models.TokenDelimiter && tok.Value == "end") || (inlineExpr && tok.Type == models.TokenDelimiter && tok.Value == ";"){
				if (tok.Type == models.TokenDelimiter && tok.Value == "end") || (inlineExpr && tok.Type == models.TokenDelimiter && tok.Value == ";" && blockWithEnds == 1){
					blockWithEnds--

				}
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
			if inlineExpr{
				p.generic("[SyntaxError] Missing ';' of 'if' statement", "S1005", fileName) // Error
			}else{
				p.generic("[SyntaxError] Missing 'end' of 'if' statement", "S1005", fileName) // Error
			}
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

	// useless
	getFirst := func(nodes []ast.Node) ast.Node{
		if len(nodes) > 0{
			return nodes
		}
		return nil
	}
	// useless

	ifAst = append(ifAst, ast.IfStatement{
		Test:       getFirst(Astnize(ifExpr, fileName, "IfStatement").([]ast.Node)),
		Consequent: getFirst(Astnize(ifBlock, fileName, "IfStatement").([]ast.Node)),
		Alternate:  getFirst(Astnize(ifAlternate, fileName, "IfStatement").([]ast.Node)),
		Line:       startLine,
		Pos:        startPos,
	})

	p.next()
	return ifAst
}
