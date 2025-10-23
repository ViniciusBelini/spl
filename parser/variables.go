package parser

import(
	// "fmt"

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

	var varAst []ast.AssignNode

	if tok.Type == models.TokenType{
		varData.Type = tok.Value
		p.next()
		if p.eof(){
			p.unexpected(fileName)
			return varAst
		}
		tok = p.peek()
	}

	if tok.Type == models.TokenIdent{
		varData.Name = tok.Value

		p.next()
		if p.eof(){
			return varAst
		}
		tok = p.peek()
	}else{
		p.generic("Unexpected token '"+tok.Value+"' ("+tok.Type+"), missing variable name", "S1003", fileName) // Error
		return varAst
	}

	if tok.Type != models.TokenAssign || (tok.Value != "=" && tok.Value != ":="){
		if tok.Type == models.TokenAssign && (tok.Value == "++" || tok.Value == "--"){
			p.back()
			p.back()
			if !p.sof(){
				tok = p.peek()
				if tok.Type == models.TokenType{
					p.unexpected(fileName) // Error
					return varAst
				}
			}
			p.next()
			p.next()

			tok = p.peek()

			method := "+"
			if tok.Value == "++"{
				method = "+"
			}else{
				method = "-"
			}

			varDataValue := ast.BinaryOpNode{
				Left: ast.IdentNode{
					Name: varData.Name,
					Line: varData.Line,
					Pos: varData.Pos,
				},
				Right: ast.LiteralNode{
					Value: "1",
					Type: "int",
					Line: varData.Line,
					Pos: varData.Pos,
				},
				Operator: method, Line: varData.Line, Pos: varData.Pos,
			}

			varAst = append(varAst, ast.AssignNode{
				Name: varData.Name,
				Type: varData.Type,
				Value: varDataValue,
				Method: method+method,
				Line: varData.Line,
				Pos: varData.Pos,
			})

			p.next()

			return varAst
		}else if tok.Type == models.TokenAssign && (tok.Value == "+=" || tok.Value == "+="){
			// To do
		}

		p.In = startIn
		return varAst
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

	varValueVerify := Astnize(varData.Value, fileName, varData.Name, true).([]ast.Node)
	var varValue ast.Node
	if len(varValueVerify) == 0{
		if config.Config["mode"] == "strict"{
			p.back()
			p.generic("Variable declaration must have an initial value in strict mode", "S1009", fileName) // Error
		}
		varValue = ast.NullNode{Line: varData.Line, Pos: varData.Pos,}
	}else{
		varValue = varValueVerify
	}

	varAst = append(varAst, ast.AssignNode{
		Name: varData.Name,
		Type: varData.Type,
		Value: varValue,
		Method: method,
		Line: varData.Line,
		Pos: varData.Pos,
	})

	return varAst
}
