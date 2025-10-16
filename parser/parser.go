package parser

import(
	"fmt"
	"strconv"

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
func Astnize(allTokens []models.Token, fileName string, inside string) ast.Node{
	p := Parser{
		Tokens: allTokens,
		Ast: nil,
		In: 0,
		Inside: inside,
	}

	for !p.eof(){
		tok := p.peek()

		switch tok.Type{
			case models.TokenType:
				if p.VariableAssignment(fileName){
					continue
				}
			case models.TokenIdent:
				if p.VariableAssignment(fileName){
					continue
				}

				p.Ast = append(p.Ast, ast.IdentNode{Name: tok.Value, Line: tok.Line, Pos: tok.Pos})
				p.next()
			case models.TokenString, models.TokenNumber, models.TokenFloat, models.TokenBoolean, models.TokenParentheses:
				if p.ParseOperators(fileName){
					continue
				}

				p.Ast = append(p.Ast, ast.LiteralNode{Value: tok.Value, Type: tok.Type, Line: tok.Line, Pos: tok.Pos})
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

// Perser Operators
func (p *Parser) ParseOperators(fileName string) bool{
	var stack []ast.Node
	var operatorStack []string

	if p.eof(){
		return false
	}

	tok := p.peek()


	if tok.Type == models.TokenString || tok.Type == models.TokenNumber || tok.Type == models.TokenFloat{
		if !p.canNext() || p.peekNext().Type != models.TokenOperator{
			return false
		}
	}

	for !p.eof(){
		tok = p.peek()

		switch tok.Type{
			case models.TokenString, models.TokenNumber, models.TokenFloat, models.TokenParentheses:
				var currentAst ast.Node
				if tok.Type == models.TokenParentheses{
					currentAstTemp := lexer.Tokenize(tok.Value[1 : len(tok.Value)-1])
					currentAst = Astnize(currentAstTemp, fileName, p.Inside)
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
				if p.canNext() && (p.peekNext().Type == models.TokenString || p.peekNext().Type == models.TokenNumber || p.peekNext().Type == models.TokenFloat || p.peekNext().Type == models.TokenParentheses){
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

	fmt.Println(stack)
	fmt.Println(operatorStack)

	if len(stack) == 1{
		p.Ast = append(p.Ast, stack[0])
	}else{
		p.unexpected(fileName)
	}

	return true
}

// Operator precedence

func precedence(op string) int{
	switch op{
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}
