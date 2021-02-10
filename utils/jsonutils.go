package utils

import (
	"encoding/json"
	"reflect"
)

var (
	typeMapStringInterface reflect.Type = reflect.TypeOf(map[string]interface{}{})
	typeArrayOfInterface   reflect.Type = reflect.TypeOf([]interface{}{})
	typeString             reflect.Type = reflect.TypeOf("")
	typeFloat64            reflect.Type = reflect.TypeOf(0.0)
)

//FilterJSON filter json data based on given template
//e.g (here only raw json is referred for illustration, this json should be converted to map[string]interface or []interface{}
// before passing to this function)
//-------------------------------------------------------------
// In template 1 is defined for non leaf node
// filter:  {"a": {"b": 1}}
// input:  {"a": {"b": {"c": 5, "d": 5}, "c":5}}
// {"a":{"b":{"c":5,"d":5}}}
// ------------------------------------------------------------
// In template 1 is defined for leaf node
// filter: {"a": {"b": {"c" : 1}}}
// input:  {"a": {"b": {"c": 5, "d": 5}, "c":5}}
// output: {"a":{"b":{"c":5}}}
// ------------------------------------------------------------
// partial mismatch in template and the data (for key b)
// filter: {"a": {"b": {"c" : 1}, "c":1}}
// input:  {"a": {"b": [{"c": 5, "d": 5}], "c":5}}
// output: {"a":{"c":5}}
// ------------------------------------------------------------
// filter on array of objects
// filter: {"a": [{"b": 1}]}
// input:  {"a": [{"b": 5, "c":5}, {"c": 5}]}
// output: {"a":[{"b":5}]}
// ------------------------------------------------------------
// Apply wild card on key match
// filter: {"a": [{"*": 1}]}
// input:  {"a": [{"b": 5, "c":5}, {"c": 5}]}
// output: {"a":[{"b":5,"c":5},{"c":5}]}
// ------------------------------------------------------------
// key value  is stringified and deepParse is true
// filter: {"a": [{"*": 1}]}
// input:  {"a": "[{\"b\": 5, \"c\":5}, {\"c\": 5}]"}
// output: {"a":[{"b":5,"c":5},{"c":5}]}
// -------------------------------------------------------------
func FilterJSON(template interface{}, data interface{}, deepParse bool) (result interface{}) {
	typeTemplate := reflect.TypeOf(template)
	typeData := reflect.TypeOf(data)

	//If leafnode then return the data as is
	if typeTemplate == typeFloat64 {
		return data
	}

	// Fist check template and data should be of same type
	if typeTemplate != typeData {
		if deepParse && typeData == typeString {
			rawData := []byte(data.(string))
			switch typeTemplate {
			case typeMapStringInterface:
				data = map[string]interface{}{}
				err := json.Unmarshal(rawData, &data)
				if err != nil {
					if Logger != nil {
						Logger.Errorf("data is not a valid json object, %s", string(rawData))
					}
					return nil
				}

			case typeArrayOfInterface:
				data = []interface{}{}
				err := json.Unmarshal(rawData, &data)
				if err != nil {
					if Logger != nil {
						Logger.Errorf("data is not a valid json array, %s", string(rawData))
					}
					return nil
				}
			default:
				if Logger != nil {
					Logger.Errorf("%s is not a supported template format", typeTemplate)
				}
				return nil
			}
		} else {
			if Logger != nil {
				Logger.Errorf("Template(%s) and data(%T) are of mismatch type", typeTemplate, typeData)
			}
			return nil
		}
	}

	// If template interface is Array
	if typeTemplate == typeArrayOfInterface {
		arrTempl := template.([]interface{})
		arrData := data.([]interface{})
		lenArrTempl := len(arrTempl)
		if lenArrTempl > 1 {
			if Logger != nil {
				Logger.Errorf("Template(%+v) should have less than one element for array", arrTempl)
			}
			return nil
		}
		if lenArrTempl == 0 {
			return arrData
		}
		for _, d := range arrData {
			res := FilterJSON(arrTempl[0], d, deepParse)
			if res != nil {
				if result == nil {
					result = []interface{}{res}
				} else {
					result = append(result.([]interface{}), res)
				}
			}
		}
		return
	}

	//If template interface is Object
	if typeTemplate == typeMapStringInterface {
		mapTempl := template.(map[string]interface{})
		mapData := data.(map[string]interface{})
		for tk, tv := range mapTempl {
			if tk == "*" {
				for k, val := range mapData {
					res := FilterJSON(tv, val, deepParse)
					if res != nil {
						if result == nil {
							result = map[string]interface{}{k: res}
						} else {
							result.(map[string]interface{})[k] = res
						}
					}
				}
				return result

			}
			val, ok := mapData[tk]
			if !ok {
				return
			}
			// Check if its a leaf node
			if _, ok := tv.(float64); ok {
				if ok {
					return map[string]interface{}{tk: val}
				}
				return nil
			}

			res := FilterJSON(tv, val, deepParse)
			if res != nil {
				if result == nil {
					result = map[string]interface{}{tk: res}
				} else {
					result.(map[string]interface{})[tk] = res
				}
			}
		}
	}
	return
}
