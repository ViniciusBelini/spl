package interpreter

import(
	"fmt"
	"strconv"
)

func MathOp(x, y interface{}, op string) (interface{}, error){
	switch x := x.(type){
		case int:
			switch y := y.(type){
				case int:
					switch op{
						case "+":
							return x + y, nil
						case "-":
							return x - y, nil
						case "*":
							return x * y, nil
						case "/":
							if y == 0 {
								return nil, fmt.Errorf("Division by zero")
							}
							return float64(x) / float64(y), nil
						case "%":
							if y == 0 {
								return nil, fmt.Errorf("Division by zero")
							}
							return x % y, nil
						default:
							return nil, fmt.Errorf("Unsupported operation '%s'", op)
					}
				case float64:
					switch op{
						case "+":
							return float64(x) + y, nil
						case "-":
							return float64(x) - y, nil
						case "*":
							return float64(x) * y, nil
						case "/":
							if y == 0 {
								return nil, fmt.Errorf("Division by zero")
							}
							return float64(x) / y, nil
						case "%":
							return nil, fmt.Errorf("Modulo not supported for float")
						default:
							return nil, fmt.Errorf("Unsupported operation '%s'", op)
					}
				default:
					return nil, fmt.Errorf("Unsupported type '%T' for operation", x)
			}
		case float64:
			switch y := y.(type){
				case int:
					switch op{
						case "+":
							return x + float64(y), nil
						case "-":
							return x - float64(y), nil
						case "*":
							return x * float64(y), nil
						case "/":
							if y == 0 {
								return nil, fmt.Errorf("Division by zero")
							}
							return x / float64(y), nil
						case "%":
							return nil, fmt.Errorf("Modulo not supported for float")
						default:
							return nil, fmt.Errorf("Unsupported operation '%s'", op)
					}
				case float64:
					switch op{
						case "+":
							return x + y, nil
						case "-":
							return x - y, nil
						case "*":
							return x * y, nil
						case "/":
							if y == 0 {
								return nil, fmt.Errorf("Division by zero")
							}
							return x / y, nil
						case "%":
							return nil, fmt.Errorf("Modulo not supported for float")
						default:
							return nil, fmt.Errorf("Unsupported operation '%s'", op)
					}
				default:
					return nil, fmt.Errorf("Unsupported type '%T' for operation", x)
			}
		case string:
			if op == "+"{
				switch y := y.(type){
					case string:
						return x + y, nil
					default:
						return nil, fmt.Errorf("Unsupported type '%T' for string operation", y)
				}
			}
			return nil, fmt.Errorf("Operation '%s' not supported for strings", op)
		default:
			return nil, fmt.Errorf("Unsupported type '%T' for operation", x)
	}
}

func MathJoin(a, b interface{}) string{
	var strA, strB string

	switch v := a.(type){
		case int:
			strA = strconv.Itoa(v)
		case float64:
			strA = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			strA = strconv.FormatBool(v)
		case string:
			strA = v
		default:
			strA = "null"
	}

	switch v := b.(type){
		case int:
			strB = strconv.Itoa(v)
		case float64:
			strB = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			strB = strconv.FormatBool(v)
		case string:
			strB = v
		default:
			strB = "null"
	}

	return strA + strB
}
