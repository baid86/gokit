package utils

import (
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
		  "ik4":  map[string]interface{}{
			  "iik4":  map[string]interface{}{
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