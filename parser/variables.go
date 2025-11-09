package parser

import(
	// "fmt"
	// "os"

	"SPL/config"
	"SPL/models"
	"SPL/ast"
)

// Variable Assignment
func (p *Parser) VariableAssignment(fileName string) []ast.AssignNode{
	tok := p.peek()
	startIn := p.In

	type VarData struct{
		Type string
		Name string
		NameReal ast.Node
		Line int
		Pos int
		Value []models.Token
	}
	varData := VarData{
		Type: "dynamic",
		Name: "null",
		NameReal: nil,
		Line: tok.Line,
		Pos: tok.Pos,
		Value: nil,
	}

	var varAst []ast.AssignNode

	if tok.Type == models.TokenType{
		varData.Type = tok.Value
		if tok.Value == "dynamic"{
			varData.Type = "<dynamic>"
		}
		p.next()
		if p.eof(){
			p.unexpected(fileName)
			return varAst
		}
		tok = p.peek()
	}

	var tmpIdent models.Token
	if tok.Type == models.TokenIdent || tok.Type == models.TokenArrayAccess{
		varData.Name = tok.Value
		tmpIdent = tok

		p.next()
		if p.eof(){
			return varAst
		}
		tok = p.peek()
	}else{
		p.generic("Unexpected token '"+tok.Value+"' ("+tok.Type+"), missing variable name", "S1003", fileName) // Error
		return varAst
	}

	if tok.Type != models.TokenAssign{
		p.In = startIn
		return varAst
	}

	if tmpIdent.Type == models.TokenArrayAccess{
		tmpVname := Astnize([]models.Token{tmpIdent}, fileName, varData.Name, true)
		varData.NameReal = tmpVname[0]
	}

	if tok.Type == models.TokenAssign && (tok.Value == "+=" || tok.Value == "-="){
		varData.Value = append(varData.Value, tmpIdent)
		method := "+"
		if tok.Value == "-="{
			method = "-"
		}

		varData.Value = append(varData.Value, models.Token{Type: models.TokenOperator, Value: method, Line: tmpIdent.Line, Pos: tmpIdent.Pos})
	}

	method := tok.Value
	p.next()

	canBreakLine := true
	returnBreakLine := 0
	for !p.eof(){
		tok = p.peek()
		p.next()


		if (tok.Type == models.TokenNewLine || tok.Value == ";") && canBreakLine{
			break
		}

		if tok.Value == "("{
			canBreakLine = false
			returnBreakLine++
		} else if tok.Value == ")"{
			returnBreakLine--
			if returnBreakLine <= 0{
				canBreakLine = true
			}
		}

		varData.Value = append(varData.Value, tok)
	}

	varValueVerify := Astnize(varData.Value, fileName, varData.Name, true)
	var varValue ast.Node
	if len(varValueVerify) == 0{
		if config.Config["mode"] == "strict"{
			p.back()
			p.generic("Variable declaration must have an initial value in strict mode", "S1009", fileName) // Error
		}
		varValue = ast.NullNode{Line: varData.Line, Pos: varData.Pos,}
	}else{
		varValue = varValueVerify[0]
	}

	varAst = append(varAst, ast.AssignNode{
		Name: varData.Name,
		NamePonter: varData.NameReal,
		Type: varData.Type,
		Value: varValue,
		Method: method,
		Line: varData.Line,
		Pos: varData.Pos,
	})

	return varAst
}
