package lexer

import(
	// "fmt"
	"regexp"

	"SPL/models"
)

func Tokenize(input string) []models.Token{
	patterns := []struct{
		Type string
		Re *regexp.Regexp
	}{
		{models.TokenComment, regexp.MustCompile(`(\/\/|\/\*|\*\/)`)},
		{models.TokenNewLine, regexp.MustCompile(`\r?\n`)},
		{models.TokenString, regexp.MustCompile(`"((?:\\.|[^"\\])*)"`)},
		{models.TokenString, regexp.MustCompile(`'((?:\\.|[^'\\])*)'`)},
		{models.TokenFloat, regexp.MustCompile(`[0-9]+\.[0-9]+`)},
		{models.TokenNumber, regexp.MustCompile(`[0-9]+`)},
		{models.TokenBoolean, regexp.MustCompile(`(true|false)`)},
		{models.TokenIfStatement, regexp.MustCompile(`(if|else)`)},
		{models.TokenType, regexp.MustCompile(`\<(int|str|bool|float)\>`)},
		{models.TokenCall, regexp.MustCompile(`([a-zA-Z0-9_]+)\((.*?)\)`)},
		{models.TokenParentheses, regexp.MustCompile(`(\((.*?)\))`)},
		{models.TokenBinOp, regexp.MustCompile(`(==|!=|>=|<=|>|<|and|or|\|\||&&|!)`)},
		{models.TokenArrayAccess, regexp.MustCompile(`[a-zA-Z0-9_]+\[(.*?)\]`)},
		{models.TokenNull, regexp.MustCompile(`null`)},
		{models.TokenAssign, regexp.MustCompile(`(=|:=|-=|\+=|--|\+\+)`)},
		{models.TokenOperator, regexp.MustCompile(`[+\-*/]`)},
		{models.TokenDelimiter, regexp.MustCompile(`(;|end)`)},
		{models.TokenIdent, regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)},
		{models.TokenSpace, regexp.MustCompile(`\s+`)},
		{models.TokenUnknown, regexp.MustCompile(`(.*)`)},
	}

	tokens := []models.Token{}
	i := 0
	line := 1
	pos := 1
	running := true
	broken := ""
	for i < len(input){
		match := false
		for _, p := range patterns{
			if loc := p.Re.FindStringIndex(input[i:]); loc != nil && loc[0] == 0{
				val := input[i+loc[0] : i+loc[1]]

				if p.Type == models.TokenComment{
					if running && val == "//" || val == "/*"{
						running = false
						broken = val
					}
				}

				if !running && p.Type == models.TokenNewLine && broken == "//"{
					running = true
					broken = ""
				}else if !running && p.Type == models.TokenComment && val == "*/"{
					running = true
					broken = ""
				}

				if p.Type == models.TokenNewLine{
					val = "null"
					line++
				}

				if p.Type != models.TokenSpace && running{
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
