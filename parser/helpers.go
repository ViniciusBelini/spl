package parser

import(
	"SPL/models"
)

func (p *Parser) GetBlock(fileName string, method string) []models.Token{
	var varBlock []models.Token
	blockWithEnds := 1

	if p.eof(){
		return varBlock
	}

	tok := p.peek()
	for !p.eof(){
		tok = p.peek()

		if hasEndDelimiter(tok.Value){
			blockWithEnds++
		}

		if tok.Type == models.TokenDelimiter && tok.Value == "end"{
			blockWithEnds--

			if blockWithEnds == 0{
				// if full{
				// 	varBlock = append(varBlock, tok)
				// }
				break
			}
		}

		varBlock = append(varBlock, tok)
		p.next()
	}

	if blockWithEnds != 0{
		if p.eof(){
			p.back()
		}

		p.generic("[SyntaxError] Missing 'end' of '"+method+"' statement", "S1005", fileName) // Error
	}

	return varBlock
}
