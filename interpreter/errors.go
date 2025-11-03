package interpreter

import(
	"strconv"
)

func TRunMakeError(id int, x string, y string, z string, fileName string, line int, pos int) string{
	errorStr := "[TypeError] "

	switch id{
		case 1:
			errorStr += "Unsupported operand type for "+z+": '"+x+"' and '"+y+"' at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 2:
			errorStr += "Cannot convert value to boolean in strict mode at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 3:
			errorStr += "Cannot assign a '"+y+"' to '"+x+"' (type "+z+") at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 4:
			errorStr += "Variable '"+x+"' declaration must include an explicit type in strict mode  at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]\nDid you mean: `<int> x := 5`?"
		case 5:
			errorStr += "Cannot increase to '"+x+"' (type "+z+") at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 6:
			errorStr += "Cannot use '+' to concatenate strings in strict mode at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]\nDid you mean to use '..' for string concatenation?"
		case 7:
			errorStr += "Cannot concatenate non-string operands: '"+x+"' and '"+y+"' at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 8:
			errorStr += "'"+x+"' expects "+y+" arguments, but "+z+" were given at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 9:
			errorStr += "Function expects a return type of '"+x+"', but returned type is '"+y+"' at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 10:
			errorStr += "Function is required to return a '"+x+"', but no value was returned at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 11:
			errorStr += "Object of type '"+x+"' has no "+y+" at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
	}

	return errorStr
}

func TGRunMakeError(id int, msg string, fileName string, line int, pos int) string{
	return "[TypeError] "+msg+" at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
}

func NRunMakeError(id int, arg string, fileName string, line int, pos int) string{
	errorStr := "[NameError] "

	switch id{
		case 1:
			errorStr += "name '"+arg+"' is not defined at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [N"+strconv.Itoa(1000+id)+"]"
		case 2:
			errorStr += "name '"+arg+"' redeclared at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [N"+strconv.Itoa(1000+id)+"]"
	}

	return errorStr
}
