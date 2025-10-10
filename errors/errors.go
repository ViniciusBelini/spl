package errors

import(
	"fmt"
	"os"
)

// type ParserError struct{
// 	Line		int
// 	Pos		int
// 	Message		string
// 	Id		string
// } Soon

func ParserError(msg string, faltalError bool) bool{ // temp
	fmt.Println(msg)

	if faltalError{
		os.Exit(1)
	}

	return true
}
