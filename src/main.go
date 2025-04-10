package main

import (
	"encoding/json"
	"fmt"
	"vsixdm/input"
	"vsixdm/loader"
)

func main() {

	flags, err := input.GetFlags()
	if err != nil {
		panic(err)
	}

	extensionId, err := loader.GetExtensionIdByLink(flags.URI)
	if err != nil {
		fmt.Printf("Error when parse %v", err)
		return
	}

	extPath, err := loader.GetExtensionById(extensionId)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if flags.JsonOut {
		result, err := json.Marshal(extPath)
		if err != nil {
			panic(err)
		}
		fmt.Print(string(result))
		return
	}
	for _, p := range *extPath {
		fmt.Println(p)
	}
}
