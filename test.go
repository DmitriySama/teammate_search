package main

import "fmt"

func main() {

	inter := map[string]interface{}{
		"AUE": 1,
		"UEA": "u",
	}

	for k, v := range inter {
		if k == "AUE" {
			fmt.Println('q')
		}
		fmt.Println(k)
		fmt.Println(v)
	}
}