package main

import (
	"fmt"
	"os"

	"/Users/g1/test_1/drone-httprequest-plugin/plugin" // Import the plugin package
)

func main() {
	fmt.Println("Starting HTTP Request Plugin...")
	plugin.ExecuteRequest()
	os.Exit(0)
}
