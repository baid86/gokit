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
