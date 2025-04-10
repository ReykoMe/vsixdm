package input

import (
	"flag"
	"fmt"
)

var example string = "https://marketplace.visualstudio.com/items?itemName=publisher.extensionId"

type Flags struct {
	URI     string
	JsonOut bool
}

func GetFlags() (*Flags, error) {
	var URI string
	var jsonOut bool
	flag.StringVar(&URI, "src", "", fmt.Sprintf("vscode marketplace link. Example:\nvsixdm.exe --src %s", example))
	flag.BoolVar(&jsonOut, "jsonOut", false, fmt.Sprintf("vscode marketplace link. Example:\nvsixdm.exe --src %s", example))
	flag.Parse()
	if len(URI) == 0 {
		return &Flags{URI: "", JsonOut: jsonOut}, fmt.Errorf("--src flag is required")
	}
	return &Flags{URI: URI, JsonOut: jsonOut}, nil
}
