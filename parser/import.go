package parser

import(
	"path/filepath"
	"strings"

	"SPL/models"
	"SPL/ast"
)

func (p *Parser) ImportSys(fileName string) []ast.ImportNode{
	var importAst []ast.ImportNode

	tok := p.peek()

	if !(tok.Type == models.TokenImport && tok.Value == "import") || !p.canNext(){
		return importAst
	}

	p.next();tok = p.peek()

	if tok.Type != models.TokenString && tok.Type != models.TokenIdent{
		p.generic("Import require a valid namespace", "S1015", fileName) // Error
		return importAst
	}

	as := tok.Value
	stringM := false
	if tok.Type == models.TokenString{
		stringM = true
		tok.Value = tok.Value[1:len(tok.Value)-1]
		as = strings.TrimSuffix(filepath.Base(tok.Value), filepath.Ext(tok.Value))
	}

	if p.canNext() && p.peekNext().Type == models.TokenImport && p.peekNext().Value == "as"{
		p.next()

		if p.canNext() && p.peekNext().Type == models.TokenIdent{
			as = p.peekNext().Value
			p.next()
		}else{
			p.generic("Import as require a valid namespace", "S1016", fileName) // Error
			return importAst
		}
	}

	importAst = append(importAst, ast.ImportNode{Path: tok.Value, As: as, String: stringM, Line: tok.Line, Pos: tok.Pos})

	p.next()
	return importAst
}
