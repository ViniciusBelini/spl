package parser

import(
//	"fmt"
	"strconv"

	"SPL/models"
	"SPL/ast"
	"SPL/errors"
)

type Parser struct{
	Tokens []models.Token
	Ast []ast.Node
	In int
}

// Main func ASTNIZE
func Astnize(allTokens []models.Token, fileName string) ast.Node{
	p := Parser{
		Tokens: allTokens,
		Ast: nil,
		In: 0,
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
			case models.TokenString:
				p.Ast = append(p.Ast, ast.StringLiteralNode{Value: tok.Value, Line: tok.Line, Pos: tok.Pos})
				p.next()
			default:
				p.unexpected(fileName)
		}
	}

	return p.Ast
}

// peek, eof, next, back funcs
func (p *Parser) peek() models.Token{return p.Tokens[p.In]}
func (p *Parser) eof() bool{if p.In >= len(p.Tokens){return true}else{return false}}
// func (p *Parser) next() bool{if p.In+1 > len(p.Tokens){return false};p.In++;return true}
func (p *Parser) next(){p.In++}
// func (p *Parser) back() bool{if p.In-1 <= 0{return false};p.In--;return true}
func (p *Parser) back(){p.In--}

// ---
// errors
func (p *Parser) unexpected(fileName string){
	tok := p.peek()

	ParserErrorMsg := "[SyntaxError] Unexpected token '"+tok.Value+"' ("+tok.Type+") at "+fileName+":"+strconv.Itoa(tok.Line)+":"+strconv.Itoa(tok.Pos)+" [S1001]" // Error

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
