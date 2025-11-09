package interpreter

import(
	"fmt"
)

func ConvToString(value any) string{
	switch v := value.(type){
		case string:
			return v
		case fmt.Stringer:
			return v.String()
		default:
			return fmt.Sprintf("%v", v)
	}
}

func convertToMapAnyAny(data interface{}) (interface{}, error){
	switch v := data.(type) {
	case map[string]interface{}:
		dataAny := make(map[any]any)

		for key, value := range v {
			convertedValue, err := convertToMapAnyAny(value)
			if err != nil {
				return nil, err
			}
			dataAny[key] = convertedValue
		}

		return dataAny, nil

	case []interface{}:
		var newArray []interface{}
		for _, item := range v {
			convertedItem, err := convertToMapAnyAny(item)
			if err != nil {
				return nil, err
			}
			newArray = append(newArray, convertedItem)
		}
		return newArray, nil

	default:
		return data, nil
	}
}
