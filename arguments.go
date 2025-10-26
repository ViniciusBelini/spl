package main

import(
	"os"

	"SPL/config"
	"SPL/errors"
)

func InterpretArgs(fileName string) bool{
	if len(os.Args) < 2{
		return true
	}

	i := 2
	for i != len(os.Args){
		if i >= len(os.Args) {
			return true
		}

		argument := os.Args[i]

		switch argument{
			case "--mode":

				if i+1 > len(os.Args) || (os.Args[i+1] != "strict" && os.Args[i+1] != "dynamic"){
					errors.ParserError("InvalidValue: Unsupported value for '--mode'. Expected 'dynamic' or 'strict'", true)
					// Give an error

					return false
				}

				config.Config["mode"] = os.Args[i+1]

				i += 2

				continue
			default:
				i++
				continue
		}

		i++
	}

	return true
}
