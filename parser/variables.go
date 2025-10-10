package parser

import(
	"SPL/models"
)

// Variable Assignment
func (p *Parser) VariableAssignment(fileName string) bool{
	tok := p.peek()

	type VarData struct{
		Type string
		Name string
		Value []models.Token
	}
	varData := VarData{
		Type: "dynamic",
		Name: "null",
		Value: nil,
	}

	if tok.Type == models.TokenType{
		varData.Type = tok.Value
		p.next()
		tok = p.peek()
	}


	if tok.Type == models.TokenIdent{
		varData.Name = tok.Value
		p.next()
		tok = p.peek()
	}else{
		p.generic("Unexpected token '"+tok.Value+"' ("+tok.Type+"), missing variable name", "S1003", fileName) // Error
		return false
	}



	p.next()

	return true
}
