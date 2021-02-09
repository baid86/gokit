package utils

import (
	"strings"
	"encoding/json"
)

// StringArrayContains find in string array
//
// This routine looks for a word in string array. Case sensitive and insensitive
// search can be done
//
func StringArrayContains(s []string, e string, casecare bool) bool {
	if !casecare {
		e = strings.ToLower(e)
	}
	for _, a := range s {
		if !casecare {
			a = strings.ToLower(a)
		}
		if a == e {
			return true
		}
	}
	return false
}


// GetFromMap look for the keys specified by []string. 
// it will iterate through each keys recursively and return the final value as interface
// if the value is not found it will return nil.
func GetFromMap(m map[string]interface{}, k []string, deepParse bool) (interface{}, bool) {
	if len(k) == 0 {
		return nil, false
	}
	
	key := k[0]
	k = k[1:]
	
	ret, ok := m[key]
	
	if len(k) == 0 {
		return ret, ok
	}

	ret1, isMap := ret.(map[string]interface{})
	if !isMap {
		if !deepParse {
			return nil, false
		}
		str, isString := ret.(string)
		if !isString {
			return nil, false
		} 
		var ret1 map[string]interface{}
		if err := json.Unmarshal([]byte(str), &ret1); err != nil {
			return nil, false
		}
		return GetFromMap(ret1, k, deepParse)
	}
	return GetFromMap(ret1, k, deepParse)	

}
