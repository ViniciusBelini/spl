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
	In int
}

// Main func ASTNIZE
func Astnize(allTokens []models.Token, fileName string) ast.Node{
	p := Parser{
		Tokens: allTokens,
		In: 0,
	}

	for !p.eof(){
		tok := p.peek()

		switch tok.Type{
			case models.TokenType:
				p.VariableAssignment(fileName)
			default:
				p.unexpected(fileName)
		}
	}

	return ast.Debug{
		Value: "asd",
	}
}

// peek, eof, next, back funcs
func (p *Parser) peek() models.Token{return p.Tokens[p.In]}
func (p *Parser) eof() bool{if p.In >= len(p.Tokens){return true}else{return false}}
func (p *Parser) next() bool{if p.In+1 > len(p.Tokens){return false};p.In++;return true}
func (p *Parser) back() bool{if p.In-1 <= 0{return false};p.In--;return true}

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
