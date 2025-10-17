package main

import(
	"fmt"
	"os"
	// "bufio"

	// "encoding/json"

	"SPL/models"
	"SPL/lexer"
	"SPL/parser"
)

// Config var

var Config = map[string]interface{}{
	"mode":		"dynamic",		// dynamic - strict
	"warnings":	true,
	"version":	"0.0.0",
	"name":		"Alpha",
}

// Main function - start point
func main(){
	// verifying arguments - must have minumum 2
	if len(os.Args) < 2{
		fmt.Println("Usage: spl <file_name> [arguments]") // change this later
		return
	}

	fileName := os.Args[1]

	interpretArgs := InterpretArgs(fileName)
	if interpretArgs == false{
		return
	}

	allTokens := readFileTokenize(fileName)
	run(allTokens, fileName)

	return
}

// input fileName read file and tokenize and return
func readFileTokenize(fileName string) []models.Token{
	var allTokens []models.Token
	data, err := os.ReadFile(fileName)
	if err != nil{
		fmt.Printf("Error: %v\n", err)
		return allTokens
	}
	allTokens = lexer.Tokenize(string(data))

	// file, err := os.Open(fileName)
	// if err != nil{
	// 	fmt.Printf("Error: %v\n", err)
	// 	return allTokens
	// }
	// defer file.Close()

	// scanner := bufio.NewScanner(file)

	// for scanner.Scan(){
	// 	line := scanner.Text()

	// 	tokens := lexer.Tokenize(line)

	// 	// allTokens = append(allTokens, models.Token{Type: "NEW_LINE", Value: "null", Line: 0, Pos: 0})
	// 	allTokens = append(allTokens, tokens...)
	// }
	// if err := scanner.Err(); err != nil{
	// 	fmt.Printf("Error: %v\n", err)
	// }

	// allTokens = append(allTokens, models.Token{Type: "NEW_LINE", Value: "null", Line: 0, Pos: 0})

	return allTokens
}

// run the program
func run(allTokens []models.Token, fileName string) bool{
	ast := parser.Astnize(allTokens, fileName, "null")

	//jsonData, _ := json.MarshalIndent(ast, "", "   ")
	fmt.Printf("%#v\n", ast)

	return true
}
