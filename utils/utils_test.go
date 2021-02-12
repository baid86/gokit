package utils

import (
	"encoding/json"
	"fmt"
)

func Example_StringArrayContains_CaseSensitive() {
	arr := []string{"alpha", "beta", "echo", "delta"}
	fmt.Println(StringArrayContains(arr, "alpha", true))
	fmt.Println(StringArrayContains(arr, "Alpha", true))
	//output:
	// true
	// false
}

func Example_StringArrayContains_IgnoreCase() {
	arr := []string{"alpha", "beta", "echo", "delta"}
	fmt.Println(StringArrayContains(arr, "alpha", false))
	fmt.Println(StringArrayContains(arr, "Alpha", false))
	fmt.Println(StringArrayContains(arr, "zulu", false))
	//output:
	// true
	// true
	// false
}

func Example_GetFromMap() {
	in := map[string]interface{}{
		"k1": "v1",
		"k2": map[string]interface{}{
			"ik2": "iv2",
		},
		"k3": `{"ik3":"iv3"}`,
		"k4": map[string]interface{}{
			"ik4": map[string]interface{}{
				"iik4": map[string]interface{}{
					"iiik4": "v4",
				},
			},
		},
		"k5": `{
		"ik5" : {
			"iik5": {
				"iiik5": "v5"
			}
		}}`,
	}
	fmt.Println(GetFromMap(in, []string{""}, false))
	fmt.Println(GetFromMap(in, []string{"k0"}, false))
	fmt.Println(GetFromMap(in, []string{"k1"}, false))
	fmt.Println(GetFromMap(in, []string{"k2", "ik2"}, false))
	fmt.Println(GetFromMap(in, []string{"k2", "ik3"}, false))
	fmt.Println(GetFromMap(in, []string{"k3", "ik3"}, false))
	fmt.Println(GetFromMap(in, []string{"k3", "ik3"}, true))
	fmt.Println(GetFromMap(in, []string{"k4", "ik4", "iik4", "iiik4"}, false))
	fmt.Println(GetFromMap(in, []string{"k5", "ik5", "iik5", "iiik5"}, true))
	//output:
	// <nil> false
	// <nil> false
	// v1 true
	// iv2 true
	// <nil> false
	// <nil> false
	// iv3 true
	// v4 true
	// v5 true
}

//
func Example_FilterJSON_1() {
	a := `{"a": {"b": 1}}`
	b := `{"a": {"b": {"c": 5, "d": 5}, "c":5}}`
	var data map[string]interface{}
	var tmpl map[string]interface{}
	json.Unmarshal([]byte(a), &tmpl)
	json.Unmarshal([]byte(b), &data)
	out := FilterJSON(tmpl, data, false)
	raw, _ := json.Marshal(&out)
	fmt.Println(string(raw))
	//output:
	// {"a":{"b":{"c":5,"d":5}}}

}

func Example_FilterJSON_2() {
	a := `{"a": {"b": {"c" : 1}}}`
	b := `{"a": {"b": {"c": 5, "d": 5}, "c":5}}`
	var data map[string]interface{}
	var tmpl map[string]interface{}
	json.Unmarshal([]byte(a), &tmpl)
	json.Unmarshal([]byte(b), &data)
	out := FilterJSON(tmpl, data, false)
	raw, _ := json.Marshal(&out)
	fmt.Println(string(raw))
	//output:
	// {"a":{"b":{"c":5}}}
}

//Example_FilterJSON_3 Template does not match the data
// The keys which are mismatched are not populated
func Example_FilterJSON_3() {
	a := `{"a": {"b": {"c" : 1}, "c":1}}`
	b := `{"a": {"b": [{"c": 5, "d": 5}], "c":5}}`
	var data map[string]interface{}
	var tmpl map[string]interface{}
	json.Unmarshal([]byte(a), &tmpl)
	json.Unmarshal([]byte(b), &data)
	out := FilterJSON(tmpl, data, false)
	raw, _ := json.Marshal(&out)
	fmt.Println(string(raw))
	//output:
	// {"a":{"c":5}}
}

// Json Contains Array
func Example_FilterJSON_4() {
	a := `{"a": [{"b": 1}]}`
	b := `{"a": [{"b": 5, "c":5}, {"c": 5}]}`
	var data map[string]interface{}
	var tmpl map[string]interface{}
	json.Unmarshal([]byte(a), &tmpl)
	json.Unmarshal([]byte(b), &data)
	out := FilterJSON(tmpl, data, false)
	raw, _ := json.Marshal(&out)
	fmt.Println(string(raw))
	//output:
	// {"a":[{"b":5}]}
}

func Example_FilterJSON_5() {
	a := `{"a": [{"*": 1}]}`
	b := `{"a": [{"b": 5, "c":5}, {"c": 5}]}`
	var data map[string]interface{}
	var tmpl map[string]interface{}
	json.Unmarshal([]byte(a), &tmpl)
	json.Unmarshal([]byte(b), &data)
	out := FilterJSON(tmpl, data, false)
	raw, _ := json.Marshal(&out)
	fmt.Println(string(raw))
	//output:
	// {"a":[{"b":5,"c":5},{"c":5}]}
}

func Example_FilterJSON_7() {
	a := `{"a": {}`
	b := `{"a": "[{\"b\": 5, \"c\":5}, {\"c\": 5}]"}`
	var data map[string]interface{}
	var tmpl map[string]interface{}
	json.Unmarshal([]byte(a), &tmpl)
	json.Unmarshal([]byte(b), &data)
	out := FilterJSON(tmpl, data, true)
	raw, _ := json.Marshal(&out)
	fmt.Println(string(raw))
	//output:
	// null
}
