package main

import (
	"fmt"
	"os"
	"strings"
	"net/http"
	"io"
)
// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

    // Get the data
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create the file
    filepath += ".gitignore"
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}

func main() {
	// get project name
	project_name := os.Args[1:]
	if len(project_name) == 0 {
		fmt.Println("usage getignore [project_name]")
	} else {
		project := strings.Join(project_name, "")
		fmt.Println("getting gitignore for ", project)

		fileURL := "https://raw.githubusercontent.com/github/gitignore/master/" + project + ".gitignore"
		fmt.Println("url : ", fileURL)
		if err := DownloadFile(project, fileURL); err != nil {
			panic(err)
		}
	}
}