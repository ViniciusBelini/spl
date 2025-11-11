package lexer

import(
	// "fmt"
	"regexp"
	"strconv"

	"SPL/errors"
	"SPL/models"
)

func Tokenize(input string, fileName string, line int, pos int) []models.Token{
	patterns := []struct{
		Type string
		Re *regexp.Regexp
	}{
		{models.TokenList, regexp.MustCompile(`(\{|\})`)},
		{models.TokenImport, regexp.MustCompile(`(\bimport\b|\bas\b)`)},
		{models.TokenComment, regexp.MustCompile(`(\/\/|\/\*|\*\/|#)`)},
		{models.TokenNewLine, regexp.MustCompile(`\r?\n`)},
		{"QUOTE", regexp.MustCompile(`("|')`)},
		{"BACK_SLASH", regexp.MustCompile(`\\`)},
		{models.TokenFloat, regexp.MustCompile(`-?[0-9]+\.[0-9]+`)},
		{models.TokenList, regexp.MustCompile(`(\{|\})`)},
		{models.TokenNumber, regexp.MustCompile(`-?[0-9]+`)},
		{models.TokenFuncStatement, regexp.MustCompile(`(\bfunction\b)`)},
		{models.TokenNativeSugar, regexp.MustCompile(`(\bshow\b)`)},
		{models.TokenBoolean, regexp.MustCompile(`(\btrue\b|\bfalse\b)`)},
		{models.TokenControlFlow, regexp.MustCompile(`(\bbreak\b|\bcontinue\b|\breturn\b)`)},
		{models.TokenIfStatement, regexp.MustCompile(`(\bif\b|\belse if\b|\belse\b)`)},
		{models.TokenLoopStatement, regexp.MustCompile(`(\bwhile\b)`)},
		{models.TokenType, regexp.MustCompile(`(\bmap|\barray)?<(int|str|bool|float)(:(.*))?>|\bdynamic\b`)},
		{"PARENTHESE", regexp.MustCompile(`(\(|\))`)},
		{models.TokenBinOp, regexp.MustCompile(`(==|!=|>=|<=|>|<|\|\||&&)`)},
		{models.TokenNull, regexp.MustCompile(`\bnull\b`)},
		{models.TokenAssign, regexp.MustCompile(`(=|:=|-=|\+=|\.\.=)`)},
		{models.TokenUnOp, regexp.MustCompile(`(!|\+\+|--)`)},
		{models.TokenOperator, regexp.MustCompile(`(\+|\-|\*|\/|%|\.\.)`)},
		{models.TokenDelimiter, regexp.MustCompile(`(;|\bend\b|:|,)`)},
		// {models.TokenObj, regexp.MustCompile(`[a-zA-Z_]\w*(?:(?:\.\w+(?:\([^()]*?(?:\([^()]*\)[^()]*)*\)|\[[^\[\]]*?(?:\[[^\[\]]*\][^\[\]]*)*\])?)|(?:\([^()]*?(?:\([^()]*\)[^()]*)*\)\.\w+)|(?:\[[^\[\]]*?(?:\[[^\[\]]*\][^\[\]]*)*\]\.\w+))+(?:\([^()]*?(?:\([^()]*\)[^()]*)*\)|\[[^\[\]]*?(?:\[[^\[\]]*\][^\[\]]*)*\])?`)},
		// {models.TokenCall, regexp.MustCompile(`([a-zA-Z0-9_]+)\((.*?)\)`)},
		// {models.TokenArrayAccess, regexp.MustCompile(`([a-zA-Z0-9_]+)\[(.*?)\]`)},
		{models.TokenArray, regexp.MustCompile(`(\[|\])`)},
		{models.TokenIdent, regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`)},
		{models.TokenSpace, regexp.MustCompile(`\s+`)},
		{models.TokenUnknown, regexp.MustCompile(`^.`)},
	}

	tokens := []models.Token{}
	i := 0
	running := true
	tempLine := 1
	tempPos := 1
	for i < len(input){
		match := false
		for _, p := range patterns{
			if loc := p.Re.FindStringIndex(input[i:]); loc != nil && loc[0] == 0{
				val := input[i+loc[0] : i+loc[1]]

				if p.Type == models.TokenNewLine{
					line++
					pos = 0
				}

				if running{
					tokens = append(tokens, models.Token{p.Type, val, line, pos})
				}

				pos = pos + len(val)

				i += loc[1]
				match = true
				break
			}
		}
		if !match{
			errors.ParserError("Undexpected character: " + string(input[i]), true)
		}
	}

	running = true
	runner := map[string]string{
		"breaker": "null",
		"helper": "null",
		"helper_2": "null",
	}
	tempLine = 1
	tempPos = 1
	tempOnce := 0
	var n_tokens []models.Token
	for i = 0;i < len(tokens);i++{
		tok := tokens[i]

		if running || running == false && runner["breaker"] == "SINGLE_COMMENT"{
			if running && tok.Type == models.TokenComment && (tok.Value == "//" || tok.Value == "#"){
				running = false
				runner["breaker"] = "SINGLE_COMMENT"
				continue
			}

			if !running{
				if tok.Type == models.TokenNewLine{
					running = true
					runner["breaker"] = "null"
				}
			}
		}

		if running || running == false && runner["breaker"] == "MULTI_COMMENT"{
			if running && tok.Type == models.TokenComment && tok.Value == "/*"{
				running = false
				runner["breaker"] = "MULTI_COMMENT"
				continue
			}

			if !running{
				if tok.Type == models.TokenComment && tok.Value == "*/"{
					running = true
					runner["breaker"] = "null"
				}

				continue
			}
		}

		if running || running == false && runner["breaker"] == "QUOTE"{
			if running && tok.Type == "QUOTE"{
				running = false
				runner["breaker"] = "QUOTE"
				runner["helper"] = tok.Value
				runner["helper_2"] = tok.Value
				tempLine = tok.Line
				tempPos = tok.Pos
				continue
			}

			if !running{
				if tok.Type == "BACK_SLASH" && i+1 < len(tokens) && tokens[i+1].Type == "QUOTE" && runner["helper_2"] == tokens[i+1].Value{
					continue
				}
				runner["helper"] += tok.Value
				if tok.Type == "QUOTE" && runner["helper_2"] == tok.Value{
					if i-1 >= 0 && tokens[i-1].Type == "BACK_SLASH"{
						continue
					}

					running = true
					n_tokens = append(n_tokens, models.Token{models.TokenString, runner["helper"], tempLine, tempPos})
				}
				continue
			}
		}

		if running || running == false && runner["breaker"] == "ARRAY_ACCESS"{
			if running && tok.Type == models.TokenIdent && i+1 < len(tokens) && tokens[i+1].Type == models.TokenArray && tokens[i+1].Value == "["{
				running = false
				runner["breaker"] = "ARRAY_ACCESS"
				runner["helper"] = tok.Value
				runner["helper_2"] = tok.Value
				tempLine = tok.Line
				tempPos = tok.Pos
				tempOnce = 0
				continue
			}

			if !running{
				runner["helper"] += tok.Value
				if tok.Type == models.TokenArray && tok.Value == "["{
					tempOnce++
				}else if tok.Type == models.TokenArray && tok.Value == "]"{
					tempOnce--

					if tempOnce == 0{
						if (i+1 < len(tokens) && tokens[i+1].Type == models.TokenArray && tokens[i+1].Value == "["){
							continue
						}
						running = true

						n_tokens = append(n_tokens, models.Token{models.TokenArrayAccess, runner["helper"], tempLine, tempPos})
					}
				}
				continue
			}
		}

		if running || running == false && runner["breaker"] == "CALL"{
			if running && tok.Type == models.TokenIdent && i+1 < len(tokens) && tokens[i+1].Type == "PARENTHESE" && tokens[i+1].Value == "("{
				running = false
				runner["breaker"] = "CALL"
				runner["helper"] = tok.Value
				runner["helper_2"] = tok.Value
				tempLine = tok.Line
				tempPos = tok.Pos
				tempOnce = 0
				continue
			}

			if !running{
				runner["helper"] += tok.Value
				if tok.Type == "PARENTHESE" && tok.Value == "("{
					tempOnce++
				}else if tok.Type == "PARENTHESE" && tok.Value == ")"{
					tempOnce--

					if tempOnce == 0{
						if (i+1 < len(tokens) && tokens[i+1].Type == "PARENTHESE" && tokens[i+1].Value == "["){
							continue
						}
						running = true

						n_tokens = append(n_tokens, models.Token{models.TokenCall, runner["helper"], tempLine, tempPos})
					}
				}
				continue
			}
		}

		if running || running == false && runner["breaker"] == "PARENTHESE"{
			if running && tok.Type == "PARENTHESE" && tok.Value == "("{
				running = false
				runner["breaker"] = "PARENTHESE"
				runner["helper"] = tok.Value
				runner["helper_2"] = tok.Value
				tempLine = tok.Line
				tempPos = tok.Pos
				tempOnce = 1
				continue
			}

			if !running{
				runner["helper"] += tok.Value
				if tok.Type == "PARENTHESE" && tok.Value == "("{
					tempOnce++
				}else if tok.Type == "PARENTHESE" && tok.Value == ")"{
					tempOnce--

					if tempOnce == 0{
						running = true
						n_tokens = append(n_tokens, models.Token{models.TokenParentheses, runner["helper"], tempLine, tempPos})
					}
				}
				continue
			}
		}
		if running || running == false && runner["breaker"] == "PARENTHESE_S"{
			if running && tok.Type == models.TokenArray && tok.Value == "["{
				running = false
				runner["breaker"] = "PARENTHESE_S"
				runner["helper"] = tok.Value
				runner["helper_2"] = tok.Value
				tempLine = tok.Line
				tempPos = tok.Pos
				tempOnce = 1
				continue
			}

			if !running{
				runner["helper"] += tok.Value
				if tok.Type == models.TokenArray && tok.Value == "["{
					tempOnce++
				}else if tok.Type == models.TokenArray && tok.Value == "]"{
					tempOnce--

					if tempOnce == 0{
						running = true
						n_tokens = append(n_tokens, models.Token{models.TokenParentheses, runner["helper"], tempLine, tempPos})
					}
				}
				continue
			}
		}

		if running || running == false && runner["breaker"] == "LIST_BLOCK"{
			if running && tok.Type == "LIST" && tok.Value == "{"{
				running = false
				runner["breaker"] = "LIST_BLOCK"
				runner["helper"] = tok.Value
				runner["helper_2"] = tok.Value
				tempLine = tok.Line
				tempPos = tok.Pos
				tempOnce = 1
				continue
			}

			if !running{
				runner["helper"] += tok.Value
				if tok.Type == "LIST" && tok.Value == "{"{
					tempOnce++
				}else if tok.Type == "LIST" && tok.Value == "}"{
					tempOnce--

					if tempOnce == 0{
						running = true

						n_tokens = append(n_tokens, models.Token{models.TokenListGroup, runner["helper"], tempLine, tempPos})
					}
				}
				continue
			}
		}

		if running && tok.Type != models.TokenSpace{
			n_tokens = append(n_tokens, models.Token{tok.Type, tok.Value, tok.Line, tok.Pos})
		}
	}

	if(!running){
		if runner["breaker"] == "QUOTE"{
			errors.ParserError("[SyntaxError] Unterminated string literal starting at "+fileName+":"+strconv.Itoa(tempLine)+":"+strconv.Itoa(tempPos)+" [S1007]", true)
		}else if runner["breaker"] == "PARENTHESE"{
			errors.ParserError("[SyntaxError] Expected ')' before end of input at "+fileName+":"+strconv.Itoa(tempLine)+":"+strconv.Itoa(tempPos)+" [S1008]", true)
		}else if runner["breaker"] == "LIST_BLOCK"{
			errors.ParserError("[SyntaxError] Expected '}' before end of list at "+fileName+":"+strconv.Itoa(tempLine)+":"+strconv.Itoa(tempPos)+" [S1008]", true)
		}else if runner["breaker"] == "ARRAY_ACCESS"{
			errors.ParserError("[SyntaxError] Expected ']' before end of input at "+fileName+":"+strconv.Itoa(tempLine)+":"+strconv.Itoa(tempPos)+" [S1008]", true)
		}else if runner["breaker"] == "CALL"{
			errors.ParserError("[SyntaxError] Expected ')' before end of function call at "+fileName+":"+strconv.Itoa(tempLine)+":"+strconv.Itoa(tempPos)+" [S1008]", true)
		}else{
			errors.ParserError("[SyntaxError] Unexpected token at "+fileName+":"+strconv.Itoa(tempLine)+":"+strconv.Itoa(tempPos)+" [S1007]", true)
		}
	}

	var nn_tokens []models.Token
	for i = 0;i < len(n_tokens);i++{
		tok := n_tokens[i]

		if running || running == false && runner["breaker"] == "obj"{
			if running && tok.Type == models.TokenIdent && i+1 < len(n_tokens) && n_tokens[i+1].Value == "."{
				running = false
				runner["breaker"] = "obj"
				runner["helper"] = tok.Value
				runner["helper_2"] = string(i)
				tempLine = tok.Line
				tempPos = tok.Pos
				tempOnce = 0
				continue
			}

			if !running{
				if !((tok.Type == models.TokenCall || tok.Type == models.TokenArrayAccess || tok.Type == models.TokenIdent) && i-1 >= 0 && n_tokens[i-1].Value == ".") && !(tok.Value == "." && i+1 < len(n_tokens) && (n_tokens[i+1].Type == models.TokenCall || n_tokens[i+1].Type == models.TokenArrayAccess || n_tokens[i+1].Type == models.TokenIdent)){
					running = true
					nn_tokens = append(nn_tokens, models.Token{models.TokenObj, runner["helper"], tempLine, tempPos})
					i--
					continue
				}

				runner["helper"] += tok.Value
				continue
			}
		}

		if running{
			nn_tokens = append(nn_tokens, models.Token{tok.Type, tok.Value, tok.Line, tok.Pos})
		}
	}

	if(!running){
		errors.ParserError("[SyntaxError] Unexpected token at "+fileName+":"+strconv.Itoa(tempLine)+":"+strconv.Itoa(tempPos)+" [S1007]", true)
	}

	return nn_tokens
}
