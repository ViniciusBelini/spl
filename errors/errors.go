package errors

import(
	"fmt"
	"os"
	// "runtime"
)

func ParserError(msg string, faltalError bool) bool{ // temp
	// file, lineC, line, ok := runtime.Caller(1)
	// if ok {
	// 	fmt.Println(file)
	// 	fmt.Println(lineC)
	// 	fmt.Println(line)
	// } else {
	// 	fmt.Println("Ooops!")
	// }

	fmt.Println(msg)

	if faltalError{
		os.Exit(1)
	}

	return true
}
