package main

import (
	"fmt"
	"os"
	"strings"
	//"net/http"
)

func main() {
	// get project name
	project_name := os.Args[1:]
	if len(project_name) == 0 {
		fmt.Println("usage getignore [project_name]")
	} else {
		project := strings.Join(project_name, "")
		fmt.Println(project)
	}
}