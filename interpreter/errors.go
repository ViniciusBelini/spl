package interpreter

import(
	// "fmt"
	// "runtime"
	"strconv"
)

func TRunMakeError(id int, x string, y string, z string, fileName string, line int, pos int) string{
	errorStr := ""

	switch id{
		case 1:
			errorStr += "[TypeError] Unsupported operand type for "+z+": '"+x+"' and '"+y+"' at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 2:
			errorStr += "[TypeError] Cannot convert value to boolean in strict mode at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 3:
			errorStr += "[TypeError] Cannot assign a '"+y+"' to '"+x+"' (expected type "+z+") at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 4:
			errorStr += "[TypeError] Variable '"+x+"' declaration must include an explicit type in strict mode at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]\nDid you mean: `<int> x := 5`?"
		case 5:
			errorStr += "[TypeError] Cannot increase to '"+x+"' (type "+z+") at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 6:
			errorStr += "[TypeError] Cannot use '+' to concatenate strings in strict mode at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]\nDid you mean to use '..' for string concatenation?"
		case 7:
			errorStr += "[TypeError] Cannot concatenate non-string operands: '"+x+"' and '"+y+"' at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 8:
			errorStr += "[TypeError] '"+x+"' expects "+y+" arguments, but "+z+" were given at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 9:
			errorStr += "[TypeError] Function expects a return type of '"+x+"', but returned type is '"+y+"' at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 10:
			errorStr += "[TypeError] Function is required to return a '"+x+"', but no value was returned at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 11:
			errorStr += "[TypeError] Object of type '"+x+"' has no "+y+" at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 12:
			errorStr += "[TypeError] Cannot pass a value of type '"+y+"' to parameter '"+x+"' (expected type "+z+") at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 13:
			errorStr += "[TypeError] Cannot convert value of type '"+y+"' to type '"+z+"' in strict mode at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
		case 14:
			errorStr += "[TypeError] Invalid collection initialization: cannot create " + x + " with values of type '" + z + "' at " + fileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(pos) + " [T" + strconv.Itoa(1000+id) + "]"
			errorStr += "\nArrays cannot contain associative types such as 'map', and maps must define key-value pairs."
		case 15:
			errorStr += "[TypeError] Invalid key for map: cannot use type " + x + " for key at " + fileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(pos) + " [T" + strconv.Itoa(1000+id) + "]"
		case 16:
			errorStr += "[NameError] Invalid variable name for "+x+" at " + fileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(pos) + " [T" + strconv.Itoa(1000+id) + "]"
		case 17:
			errorStr += "[ValueError] Item '"+x+"' not found in map at " + fileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(pos) + " [T" + strconv.Itoa(1000+id) + "]"
		case 18:
			errorStr += "[TypeError] '"+x+"' object is not subscriptable at " + fileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(pos) + " [T" + strconv.Itoa(1000+id) + "]"
		case 19:
			errorStr += "[ValueError] Array index out of range at " + fileName + ":" + strconv.Itoa(line) + ":" + strconv.Itoa(pos) + " [T" + strconv.Itoa(1000+id) + "]"
	}

	return errorStr
}

func TGRunMakeError(id int, msg string, fileName string, line int, pos int) string{
	return "[TypeError] "+msg+" at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [T"+strconv.Itoa(1000+id)+"]"
}

func NRunMakeError(id int, arg string, fileName string, line int, pos int) string{
	// file, lineC, line, ok := runtime.Caller(1)
	// if ok {
	// 	fmt.Println(file)
	// 	fmt.Println(lineC)
	// 	fmt.Println(line)
	// } else {
	// 	fmt.Println("Ooops!")
	// }

	errorStr := "[NameError] "

	switch id{
		case 1:
			errorStr += "name '"+arg+"' is not defined at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [N"+strconv.Itoa(1000+id)+"]"
		case 2:
			errorStr += "name '"+arg+"' redeclared at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [N"+strconv.Itoa(1000+id)+"]"
	}

	return errorStr
}

func MRunMakeError(id int, arg string, fileName string, line int, pos int) string{
	errorStr := "[ModuleError] "

	switch id{
		case 1:
			errorStr += "Module path error at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [N"+strconv.Itoa(1000+id)+"]"
		case 2:
			errorStr += "File path '"+arg+"' doesnt exists at "+fileName+":"+strconv.Itoa(line)+":"+strconv.Itoa(pos)+" [N"+strconv.Itoa(1000+id)+"]"
	}

	return errorStr
}
