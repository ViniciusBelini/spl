package parser

import(
	// "fmt"

	"SPL/models"
	"SPL/ast"
)

func (p *Parser) ArrayOne(fileName string) []ast.ArrayOneItem{
	var items []ast.ArrayOneItem

	if p.eof() || p.peek().Type != models.TokenList || p.peek().Value != "{"{
		p.unexpected(fileName)
	}

	tok := p.peek()
	bcTok := tok

	closeMain := 1
	loopGetKey := func(typeExpected int) []models.Token{
		var alTk []models.Token
		for !p.eof(){
			tok = p.peek()

			if tok.Value == "{"{
				closeMain++
			}else if tok.Value == "}"{
				closeMain--
			}

			if tok.Type == models.TokenNewLine{
				p.next()
				continue
			}

			if typeExpected == 1 && tok.Value == ":"{
				p.unexpected(fileName)
			}

			if tok.Value == ":" || tok.Value == "," || tok.Value == "}"{
				if (closeMain == 1 && tok.Value != "}") || (closeMain == 0 && tok.Value == "}"){
					if len(alTk) > 0{
						break
					}
				}
			}

			alTk = append(alTk, tok)

			p.next()
		}
		return alTk
	}
	p.next()

	var firstValue []models.Token
	var lastValue []models.Token
	for !p.eof(){
		tok = p.peek()

		for true{
			firstValue = loopGetKey(0)
			if len(firstValue) > 0{
				break
			}
		}

		if len(firstValue) == 0{
			if p.eof(){
				p.back()
			}
			p.generic("Missing value in list declaration", "S1018", fileName) // Error
		}

		if !p.eof() && p.peek().Value == ":"{
			if p.canNext(){
				p.next()
				for true{
					lastValue = loopGetKey(1)
					if len(lastValue) > 0{
						break
					}
				}

				if len(lastValue) == 0{
					if p.eof(){
						p.back()
					}
					p.generic("Missing value in list declaration", "S1018", fileName) // Error
				}
			}else{
				p.unexpected(fileName)
			}
		}else{
			lastValue = []models.Token{}
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

		left := Astnize(firstValue, fileName, p.Inside, true)
		if len(left) > 1{
			p.generic("Expected ',' before end of value in list item", "S1017", fileName) // Error
		}

		right := Astnize(lastValue, fileName, p.Inside, true)
		if len(right) > 1{
			p.generic("Expected ',' before end of value in list item", "S1017", fileName) // Error
		}

		items = append(items, ast.ArrayOneItem{
			Left: getFirst(left, true),
			Right: getFirst(right, true),
			Line: bcTok.Line,
			Pos: bcTok.Pos,
		})

		if !p.eof() && p.peek().Value == ","{
			p.next()
			continue
		}
		break
	}

	if p.eof(){
		p.back()
		p.generic("Expected ',' before end of value in list", "S1017", fileName) // Error
	}
	if p.peek().Value != "}"{
		p.generic("Expected ',' before end of value in list", "S1017", fileName) // Error
	}
	p.next()

	return items
}
