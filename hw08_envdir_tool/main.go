package main

import (
	"fmt"
	"os"
)

func main() {
	dir := os.Args[1]

	environments, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)

		return
	}

	returnCode := RunCmd(os.Args[2:], environments)

	os.Exit(returnCode)
}
