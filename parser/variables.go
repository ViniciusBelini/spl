package parser

import(
	"SPL/models"
	"SPL/ast"
)

// Variable Assignment
func (p *Parser) VariableAssignment(fileName string) bool{
	tok := p.peek()

	type VarData struct{
		Type string
		Name string
		Line int
		Pos int
		Value []models.Token
	}
	varData := VarData{
		Type: "dynamic",
		Name: "null",
		Line: tok.Line,
		Pos: tok.Pos,
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

	if tok.Type != models.TokenAssign || (tok.Value != "=" && tok.Value != ":="){
		// p.expected("one of '=', ':='", fileName) // Error
		p.back()
		return false
	}

	canBreakLine := true
	returnBreakLine := 0
	for !p.eof(){
		p.next()
		if !p.eof(){
			tok = p.peek()
		}

		if (tok.Type == models.TokenNewLine || tok.Value == ";") && canBreakLine{
			break
		}

		if tok.Value == "("{
			canBreakLine = false
			returnBreakLine++
		}else if tok.Value == ")"{
			returnBreakLine--

			if returnBreakLine <= 0{
				canBreakLine = true
			}
		}

		varData.Value = append(varData.Value, tok)
	}

	varAst := ast.AssignNode{
		Name: varData.Name,
		Type: varData.Type,
		Value: Astnize(varData.Value, fileName),
		Line: varData.Line,
		Pos: varData.Pos,
	}

	p.Ast = append(p.Ast, varAst)

	return true
}
