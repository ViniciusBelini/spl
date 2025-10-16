package lexer

import(
	"regexp"

	"SPL/models"
)

func Tokenize(input string) []models.Token{
	patterns := []struct{
		Type string
		Re *regexp.Regexp
	}{
		{models.TokenNewLine, regexp.MustCompile(`\n`)},
		{models.TokenString, regexp.MustCompile(`"((?:\\.|[^"\\])*)"`)},
		{models.TokenString, regexp.MustCompile(`'((?:\\.|[^'\\])*)'`)},
		{models.TokenFloat, regexp.MustCompile(`[0-9]+\.[0-9]+`)},
		{models.TokenNumber, regexp.MustCompile(`[0-9]+`)},
		{models.TokenBoolean, regexp.MustCompile(`(true|false)`)},
		{models.TokenIfStatement, regexp.MustCompile(`(if|else|elif|elseif)`)},
		{models.TokenType, regexp.MustCompile(`\<(int|str|bool|float)\>`)},
		{models.TokenComment, regexp.MustCompile(`\/\/(.*?)$`)},
		{models.TokenCall, regexp.MustCompile(`([a-zA-Z0-9_]+)\((.*?)\)`)},
		{models.TokenParentheses, regexp.MustCompile(`(\((.*?)\))`)},
		{models.TokenBinOp, regexp.MustCompile(`(==|!=|>=|<=|>|<|and|or|\|\||&&|!)`)},
		{models.TokenArrayAccess, regexp.MustCompile(`[a-zA-Z0-9_]+\[(.*?)\]`)},
		{models.TokenNull, regexp.MustCompile(`null`)},
		{models.TokenIdent, regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)},
		{models.TokenAssign, regexp.MustCompile(`(=|:=|-=|\+=|--|\+\+)`)},
		{models.TokenOperator, regexp.MustCompile(`[+\-*/]`)},
		{models.TokenDelimiter, regexp.MustCompile(`(;|end)`)},
		{models.TokenSpace, regexp.MustCompile(`\s+`)},
		{models.TokenUnknown, regexp.MustCompile(`(.*)`)},
	}

	tokens := []models.Token{}
	i := 0
	line := 1
	pos := 1
	for i < len(input){
		match := false
		for _, p := range patterns{
			if loc := p.Re.FindStringIndex(input[i:]); loc != nil && loc[0] == 0{
				val := input[i+loc[0] : i+loc[1]]

				if p.Type == models.TokenNewLine{
					line++
				}

				if p.Type != models.TokenSpace{
					tokens = append(tokens, models.Token{p.Type, val, line, pos})
				}

				pos = pos + len(val)

				i += loc[1]
				match = true
				break
			}
		}
		if !match{
			panic("Undexpected character: " + string(input[i]))
		}
	}
	return tokens
}
