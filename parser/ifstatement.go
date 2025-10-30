package parser

import(
	// "fmt"

	"SPL/models"
	"SPL/ast"
)

// If Statement
func (p *Parser) IfStatement(fileName string) []ast.IfStatement{
	var ifAst []ast.IfStatement

	if p.eof(){
		return ifAst
	}
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
	if exprType == "if" || exprType == "else if"{
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
				if blockWithEnds-1 <= 0 && tok.Type == models.TokenIfStatement && (tok.Value == "else" || tok.Value == "else if"){
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

	if p.eof(){
		return ifAst
	}
	tok = p.peek()

	ifBlock := loopGetBlock(false)

	var ifAlternate []models.Token
	elseOnly := true
	if tok.Type == models.TokenIfStatement{
		if tok.Value == "else" || tok.Value == "else if"{
			if tok.Value == "else"{
				p.next()
				elseOnly = true
			}else{
				elseOnly = false
			}
			ifAlternate = p.GetBlock(fileName, "if")
		}else{
			p.unexpected(fileName)
		}
	}

	tok = p.peek()

	if tok.Type != models.TokenDelimiter && tok.Value != "end"{
		if p.eof(){
			p.back()
		}
		p.generic("[SyntaxError] Missing 'end' of 'if' statement", "S1005", fileName) // Error
	}

	if !elseOnly{
		ifAlternate = append(ifAlternate, tok)
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

	inside := "IfStatement/"+p.Inside

	ifAst = append(ifAst, ast.IfStatement{
		Test:       getFirst(Astnize(ifExpr, fileName, inside, true), true),
		Consequent: getFirst(Astnize(ifBlock, fileName, inside, false), false),
		Alternate:  getFirst(Astnize(ifAlternate, fileName, inside, false), false),
		Line:       startLine,
		Pos:        startPos,
	})

	p.next()
	return ifAst
}
