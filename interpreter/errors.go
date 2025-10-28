package interpreter

import(
	"strconv"
)

func TRunMakeError(id int, x string, y string, z string, fileName string, line int, pos int) string{
	errorStr := "[TypeError] "

	switch id{
		case 1:
			errorStr += "Unsupported operand type for "+z+": '"+x+"' and '"+y+"'"
		case 2:
			errorStr += "Cannot convert value to boolean in strict mode"
	}

	errorStr += " at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"

	return errorStr
}

func TGRunMakeError(id int, msg string, fileName string, line int, pos int) string{
	return "[TypeError] "+msg+" at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
}
