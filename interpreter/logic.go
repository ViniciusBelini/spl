package interpreter

import(
	"fmt"
)

func CompareOp(x, y interface{}, op string) (interface{}, error){
	switch x := x.(type){
	case int:
		switch y := y.(type){
		case int:
			switch op{
			case "==":
				return x == y, nil
			case "!=":
				return x != y, nil
			case ">":
				return x > y, nil
			case "<":
				return x < y, nil
			case ">=":
				return x >= y, nil
			case "<=":
				return x <= y, nil
			case "&&":
				return nil, fmt.Errorf("Operator '&&' not supported for integers")
			case "||":
				return nil, fmt.Errorf("Operator '||' not supported for integers")
			default:
				return nil, fmt.Errorf("Unsupported operation '%s' for type '%T'", op, x)
			}
		case float64:
			switch op{
			case "==":
				return float64(x) == y, nil
			case "!=":
				return float64(x) != y, nil
			case ">":
				return float64(x) > y, nil
			case "<":
				return float64(x) < y, nil
			case ">=":
				return float64(x) >= y, nil
			case "<=":
				return float64(x) <= y, nil
			case "&&":
				return nil, fmt.Errorf("Operator '&&' not supported for int and float")
			case "||":
				return nil, fmt.Errorf("Operator '||' not supported for int and float")
			default:
				return nil, fmt.Errorf("Unsupported operation '%s' for type '%T' and '%T'", op, x, y)
			}
		case string:
			switch op {
			case "==":
				return fmt.Sprintf("%d", x) == y, nil
			case "!=":
				return fmt.Sprintf("%d", x) != y, nil
			case "&&":
				return nil, fmt.Errorf("Operator '&&' not supported for int and string")
			case "||":
				return nil, fmt.Errorf("Operator '||' not supported for int and string")
			default:
				return nil, fmt.Errorf("Unsupported operation '%s' for type '%T' and '%T'", op, x, y)
			}
		default:
			return nil, fmt.Errorf("Unsupported type '%T' for operation", x)
		}
	case float64:
		switch y := y.(type) {
		case int:
			switch op {
			case "==":
				return x == float64(y), nil
			case "!=":
				return x != float64(y), nil
			case ">":
				return x > float64(y), nil
			case "<":
				return x < float64(y), nil
			case ">=":
				return x >= float64(y), nil
			case "<=":
				return x <= float64(y), nil
			case "&&":
				return nil, fmt.Errorf("Operator '&&' not supported for float and int")
			case "||":
				return nil, fmt.Errorf("Operator '||' not supported for float and int")
			default:
				return nil, fmt.Errorf("Unsupported operation '%s' for type '%T' and '%T'", op, x, y)
			}
		case float64:
			switch op {
			case "==":
				return x == y, nil
			case "!=":
				return x != y, nil
			case ">":
				return x > y, nil
			case "<":
				return x < y, nil
			case ">=":
				return x >= y, nil
			case "<=":
				return x <= y, nil
			case "&&":
				return nil, fmt.Errorf("Operator '&&' not supported for floats")
			case "||":
				return nil, fmt.Errorf("Operator '||' not supported for floats")
			default:
				return nil, fmt.Errorf("Unsupported operation '%s' for type '%T' and '%T'", op, x, y)
			}
		case string:
			switch op {
			case "==":
				return fmt.Sprintf("%f", x) == y, nil
			case "!=":
				return fmt.Sprintf("%f", x) != y, nil
			case "&&":
				return nil, fmt.Errorf("Operator '&&' not supported for float and string")
			case "||":
				return nil, fmt.Errorf("Operator '||' not supported for float and string")
			default:
				return nil, fmt.Errorf("Unsupported operation '%s' for type '%T' and '%T'", op, x, y)
			}
		default:
			return nil, fmt.Errorf("Unsupported type '%T' for operation", x)
		}
	case string:
		switch y := y.(type) {
		case string:
			switch op {
			case "==":
				return x == y, nil
			case "!=":
				return x != y, nil
			case ">":
				return x > y, nil
			case "<":
				return x < y, nil
			case ">=":
				return x >= y, nil
			case "<=":
				return x <= y, nil
			case "&&":
				return nil, fmt.Errorf("Operator '&&' not supported for strings")
			case "||":
				return nil, fmt.Errorf("Operator '||' not supported for strings")
			default:
				return nil, fmt.Errorf("Unsupported operation '%s' for type '%T'", op, x)
			}
		default:
			return nil, fmt.Errorf("Unsupported type '%T' for operation", x)
		}
	case bool:
		switch y := y.(type) {
		case bool:
			switch op {
			case "==":
				return x == y, nil
			case "!=":
				return x != y, nil
			case "&&":
				if x == false{
					return false, nil
				}
				if y == false{
					return false, nil
				}
				return true, nil
			case "||":
				return x || y, nil
			default:
				return nil, fmt.Errorf("Unsupported operation '%s' for type '%T'", op, x)
			}
		default:
			return nil, fmt.Errorf("Unsupported type '%T' for operation", x)
		}
	default:
		return nil, fmt.Errorf("Unsupported type '%T' for operation", x)
	}
}

func UnaryOpConv(x interface{}) (interface{}, error){
	switch x := x.(type){
		case int:
			return x == 0, nil
		case float64:
			return x == 0.0, nil
		case bool:
			return x == false, nil
		case string:
			return x == "", nil
		default:
			return nil, fmt.Errorf("Unsupported type '%T' for operation", x)
	}
}
