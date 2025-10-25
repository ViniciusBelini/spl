package parser

import(
	// "fmt"
	"strconv"
	// "runtime"

	"SPL/config"
	"SPL/lexer"
	"SPL/models"
	"SPL/ast"
	"SPL/errors"
)

type Parser struct{
	Tokens []models.Token
	Ast []ast.Node
	In int
	Inside string
}

// Main func ASTNIZE
func Astnize(allTokens []models.Token, fileName string, inside string, statementExpr bool) ast.Node{
	p := Parser{
		Tokens: allTokens,
		Ast: nil,
		In: 0,
		Inside: inside,
	}

	for !p.eof(){
		tok := p.peek()

		switch tok.Type{
			case models.TokenControlFlow:
				pTemp := p
				tempAST := p.ControlFlow(fileName)
				if len(tempAST) > 0 && !statementExpr{
					p.Ast = append(p.Ast, tempAST[0])
					continue
				}
				p = pTemp
				p.unexpected(fileName)
			case models.TokenObj, models.TokenArrayAccess, models.TokenCall:
				p.next()
				continue
			case models.TokenNewLine:
				p.next()
				continue
			case models.TokenType:
				pTemp := p
				tempAST := p.VariableAssignment(fileName)
				if len(tempAST) > 0 && (!statementExpr || statementExpr && config.Config["mode"] == "dynamic"){
					p.Ast = append(p.Ast, tempAST[0])
					continue
				}else if len(tempAST) > 0{
					p.back()
					p.generic("[SyntaxError] '=' (ASSIGN) is not valid here in strict mode – variable declarations must be top-level", "S1006", fileName) // Error
				}
				// REMIDER: allow variable declaration without value `<int> age;` - ; is optional
				p = pTemp
				p.unexpected(fileName)
			case models.TokenIfStatement:
				if tok.Value == "if"{
					pTemp := p
					tempAST := p.IfStatement(fileName)
					if len(tempAST) > 0 && !statementExpr{
						p.Ast = append(p.Ast, tempAST[0])
						continue
					}
					p = pTemp
					p.unexpected(fileName)
				}else{
					p.unexpected(fileName)
				}
			case models.TokenLoopStatement:
				p.LoopStatementParser(fileName, statementExpr)
			case models.TokenUnOp:
				pTemp := p
				tempAST := p.UnExpr(fileName)
				if len(tempAST) > 0{
					p.Ast = append(p.Ast, tempAST[0])
					continue
				}
				p = pTemp
				p.unexpected(fileName)
			case models.TokenIdent, models.TokenString, models.TokenNumber, models.TokenFloat, models.TokenBoolean, models.TokenParentheses:
				pTemp := p
				tempAST2 := p.ParserLogical(fileName)
				if len(tempAST2) > 0{
					p.Ast = append(p.Ast, tempAST2[0])
					continue
				}
				p = pTemp

				pTemp = p
				tempAST2 = p.ParseOperators(fileName)
				if len(tempAST2) > 0{
					p.Ast = append(p.Ast, tempAST2[0])
					continue
				}
				p = pTemp

				pTemp = p
				if tok.Type == models.TokenIdent{
					tempAST := p.VariableAssignment(fileName)
					if len(tempAST) > 0 && (!statementExpr || statementExpr && config.Config["mode"] == "dynamic"){
						p.Ast = append(p.Ast, tempAST[0])
						continue
					}else if len(tempAST) > 0{
						p.back()
						p.generic("'=' (ASSIGN) is not valid here in strict mode – variable declarations must be top-level", "S1006", fileName) // Error
					}
					p = pTemp
				}

				if tok.Type == models.TokenIdent{
					p.Ast = append(p.Ast, ast.IdentNode{Name: tok.Value, Line: tok.Line, Pos: tok.Pos})
				}else if tok.Type == models.TokenParentheses{
					p.Ast = append(p.Ast, Astnize(lexer.Tokenize(tok.Value[1 : len(tok.Value)-1], fileName, tok.Line, tok.Pos), fileName, p.Inside, statementExpr).([]ast.Node)[0])
				}else{
					p.Ast = append(p.Ast, ast.LiteralNode{Value: tok.Value, Type: tok.Type, Line: tok.Line, Pos: tok.Pos})
				}
				p.next()
			case models.TokenNull:
				p.Ast = append(p.Ast, ast.NullNode{Line: tok.Line, Pos: tok.Pos,})
				p.next()
			default:
				p.unexpected(fileName)
		}
	}

	return p.Ast
}

// peek, eof, sof, next, back funcs
func (p *Parser) peek() models.Token{return p.Tokens[p.In]}
func (p *Parser) peekNext() models.Token{return p.Tokens[p.In+1]}
func (p *Parser) peekBack() models.Token{return p.Tokens[p.In-1]}
func (p *Parser) eof() bool{if p.In >= len(p.Tokens){return true}else{return false}}
func (p *Parser) sof() bool{if p.In < 0{return true}else{return false}}
// func (p *Parser) next() bool{if p.In+1 > len(p.Tokens){return false};p.In++;return true}
func (p *Parser) next(){p.In++}
func (p *Parser) canNext() bool{if p.In+2 > len(p.Tokens){return false};return true}
// func (p *Parser) back() bool{if p.In-1 <= 0{return false};p.In--;return true}
func (p *Parser) back(){p.In--}
func (p *Parser) canBack() bool{if p.In-1 >= 0{return true};return false}

// ---
// errors
func (p *Parser) unexpected(fileName string){
	// file, lineC, line, ok := runtime.Caller(1)
	// if ok {
	// 	fmt.Println(file)
	// 	fmt.Println(lineC)
	// 	fmt.Println(line)
	// } else {
	// 	fmt.Println("Ooops!")
	// }

	ParserErrorMsg := "[SyntaxError] Unexpected token at "+fileName+" [S1001]" // Error

	if !p.eof(){
		tok := p.peek()

		ParserErrorMsg = "[SyntaxError] Unexpected token '"+tok.Value+"' ("+tok.Type+") at "+fileName+":"+strconv.Itoa(tok.Line)+":"+strconv.Itoa(tok.Pos)+" [S1001]" // Error
	}

	errors.ParserError(ParserErrorMsg, true)
}
func (p *Parser) generic(message string, id string, fileName string){
	tok := p.peek()

	ParserErrorMsg := "[SyntaxError] "+message+" at "+fileName+":"+strconv.Itoa(tok.Line)+":"+strconv.Itoa(tok.Pos)+" ["+id+"]" // Error

	errors.ParserError(ParserErrorMsg, true)
}
func (p *Parser) expected(token string, fileName string){
	tok := p.peek()

	ParserErrorMsg := "[SyntaxError] Unexpected token '"+tok.Value+"' ("+tok.Type+"), expected "+token+" at "+fileName+":"+strconv.Itoa(tok.Line)+":"+strconv.Itoa(tok.Pos)+" [S1002]" // Error

	errors.ParserError(ParserErrorMsg, true)
}

// Perser Logical
func (p *Parser) ParserLogical(fileName string) []ast.BinaryOpNode{
	var logicalAst []ast.BinaryOpNode

	if p.eof(){
		return logicalAst
	}

	tok := p.peek()

	if !isLiteral(tok.Type) && tok.Type != models.TokenParentheses && tok.Type != models.TokenIdent{
		return logicalAst
	}

	var stack []ast.Node
	var exprStack []models.Token
	var logicalStack []string

	startIn := p.In
	for !p.eof(){
		tok = p.peek()

		switch tok.Type{
			case models.TokenBinOp:
				logicalStack = append(logicalStack, tok.Value)
				currentAstTemp := Astnize(exprStack, fileName, p.Inside, true).([]ast.Node)[0]
				stack = append(stack, currentAstTemp)
				exprStack = []models.Token{}
				p.next()
			default:
				var currentAstTemp []models.Token
				if tok.Type == models.TokenParentheses{
					currentAstTemp = lexer.Tokenize(tok.Value[1 : len(tok.Value)-1], fileName, tok.Line, tok.Pos)
				}else{
					currentAstTemp = append(currentAstTemp, tok)
				}
				for i := 0;i < len(currentAstTemp);i++{
					exprStack = append(exprStack, currentAstTemp[i])
				}

				if !p.canNext() && len(stack) >= 1{
					currentAstTemp := Astnize(exprStack, fileName, p.Inside, true).([]ast.Node)[0]
					stack = append(stack, currentAstTemp)
					exprStack = []models.Token{}
				}
				p.next()
		}
	}

	if len(stack) < 2{
		p.In = startIn
		return logicalAst
	}

	for len(logicalStack) > 0{
		maxPrec := -1
		maxIndex := -1
		for i, op := range logicalStack{
			if precedence(op) > maxPrec{
				maxPrec = precedence(op)
				maxIndex = i
			}
		}

		op := logicalStack[maxIndex]
		logicalStack = append(logicalStack[:maxIndex], logicalStack[maxIndex+1:]...)

		left := stack[maxIndex]
		if maxIndex+1 >= len(stack){
			p.In = startIn
			return logicalAst
		}
		right := stack[maxIndex+1]

		stack = append(stack[:maxIndex], stack[maxIndex+2:]...)

		node := ast.BinaryOpNode{
			Left:     left,
			Right:    right,
			Operator: op,
			Line:     1,
			Pos:      1,
		}

		stack = append(stack[:maxIndex], append([]ast.Node{node}, stack[maxIndex:]...)...)
	}

	if len(stack) == 1{
		return []ast.BinaryOpNode{stack[0].(ast.BinaryOpNode)}
	}else{
		p.unexpected(fileName)
	}

	return logicalAst
}

// Perser Operators
func (p *Parser) ParseOperators(fileName string) []ast.BinaryOpNode{
	var operatorsAst []ast.BinaryOpNode

	if p.eof(){
		return operatorsAst
	}

	tok := p.peek()

	if isLiteral(tok.Type) || tok.Type == models.TokenParentheses || tok.Type == models.TokenIdent{
		if !p.canNext() || p.peekNext().Type != models.TokenOperator{
			return operatorsAst
		}
	}else{
		return operatorsAst
	}

	var stack []ast.Node
	var operatorStack []string

	for !p.eof(){
		tok = p.peek()

		switch tok.Type{
			case models.TokenIdent, models.TokenString, models.TokenNumber, models.TokenFloat, models.TokenParentheses:
				var currentAst ast.Node
				if tok.Type == models.TokenParentheses{
					currentAstTemp := lexer.Tokenize(tok.Value[1 : len(tok.Value)-1], fileName, tok.Line, tok.Pos)
					currentAst = Astnize(currentAstTemp, fileName, p.Inside, true).([]ast.Node)[0]
				}else if tok.Type == models.TokenIdent{
					currentAst = ast.IdentNode{
						Name: tok.Value, Line: tok.Line, Pos: tok.Pos,
					}
				}else{
					currentAst = ast.LiteralNode{
						Value: tok.Value, Type: tok.Type, Line: tok.Line, Pos: tok.Pos,
					}
				}

				stack = append(stack, currentAst)
				p.next()

				if p.canNext() && p.peekNext().Type == models.TokenOperator{
					continue
				}else{
					break
				}
			case models.TokenOperator:
				if p.canNext() && (isLiteral(p.peekNext().Type) || p.peekNext().Type == models.TokenParentheses || p.peekNext().Type == models.TokenIdent){
					operatorStack = append(operatorStack, tok.Value)
					p.next()
					continue
				}else{
					p.unexpected(fileName)
				}
			default:
				p.next()
				break
		}
	}

	for len(operatorStack) > 0{
		maxPrec := -1
		maxIndex := -1
		for i, op := range operatorStack{
			if precedence(op) > maxPrec{
				maxPrec = precedence(op)
				maxIndex = i
			}
		}

		op := operatorStack[maxIndex]
		operatorStack = append(operatorStack[:maxIndex], operatorStack[maxIndex+1:]...)

		left := stack[maxIndex]
		right := stack[maxIndex+1]

		stack = append(stack[:maxIndex], stack[maxIndex+2:]...)

		node := ast.BinaryOpNode{
			Left:     left,
			Right:    right,
			Operator: op,
			Line:     1,
			Pos:      1,
		}

		stack = append(stack[:maxIndex], append([]ast.Node{node}, stack[maxIndex:]...)...)
	}

	if len(stack) == 1{
		return []ast.BinaryOpNode{stack[0].(ast.BinaryOpNode)}
	}else{
		p.unexpected(fileName)
	}

	return operatorsAst // never reach here
}

// Operator precedence
func precedence(op string) int {
	switch op {
	case "||":
		return 1
	case "&&":
		return 2
	case "==", "!=":
		return 3
	case "<", ">", "<=", ">=":
		return 4
	case "+", "-":
		return 5
	case "*", "/", "%":
		return 6
	// case "!", "++", "--", "-u":
	// 	return 7
	default:
		return 0
	}
}

// Verify is literal
func isLiteral(token string) bool{
	if token == models.TokenString || token == models.TokenNumber || token == models.TokenFloat{
		return true
	}

	return false
}

// Verify if is a bloc with END delimiter
func hasEndDelimiter(token string) bool{
	if token == "if" || token == "while"{
		return true
	}

	return false
}

// Control flow
func (p *Parser) ControlFlow(fileName string) []ast.ControlFlowNode{
	var flowAst []ast.ControlFlowNode

	if p.eof(){
		return flowAst
	}

	tok := p.peek()

	tempLinePos := map[string]int{
		"line": tok.Line,
		"pos": tok.Pos,
	}
	switch tok.Value{
		case "continue", "break":
			method := tok.Value
			if p.Inside == "LoopStatement"{
				var controlTokens []models.Token
				p.next()
				for !p.eof(){
					tok = p.peek()

					if tok.Type == models.TokenNewLine{
						break
					}

					controlTokens = append(controlTokens, tok)
					p.next()
				}

				controlAstVerify := Astnize(controlTokens, fileName, "LoopStatement", true).([]ast.Node)
				var controlAst ast.Node
				if len(controlAstVerify) == 0{
					controlAst = ast.NullNode{Line: tempLinePos["line"], Pos: tempLinePos["pos"],}
				}else{
					controlAst = controlAstVerify
				}

				flowAst = append(flowAst, ast.ControlFlowNode{Method: method, Argument: controlAst, Line: tempLinePos["line"], Pos: tempLinePos["pos"],})
				return flowAst
			}else{
				p.generic("'"+tok.Value+"' cannot be used outside of a loop statement", "S1010", fileName)
			}
		default:
			return flowAst
	}

	return flowAst
}
