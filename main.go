package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Templates struct to hold all the gitignore templates
type Templates struct {
	templates map[string]string
}

// NewTemplates creates a new Templates instance
func NewTemplates() *Templates {
	return &Templates{
		templates: make(map[string]string),
	}
}

// LoadTemplates loads all the gitignore templates from the provided directory
func (t *Templates) LoadTemplates() error {
	// Load templates from the current directory
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}

	// Also load from "Global" and "community" directories if they exist
	dirNames := []string{".", "Global", "community"}
	
	for _, dirName := range dirNames {
		if dirName != "." {
			// Check if the directory exists
			if _, err := os.Stat(dirName); os.IsNotExist(err) {
				continue
			}
		}
		
		// Get files from the directory
		var dirFiles []os.FileInfo
		if dirName == "." {
			dirFiles = files
		} else {
			dirFiles, err = ioutil.ReadDir(dirName)
			if err != nil {
				continue // Skip if can't read directory
			}
		}
		
		for _, file := range dirFiles {
			if file.IsDir() {
				if dirName == "." && (file.Name() == "Global" || file.Name() == "community") {
					// We handle these directories separately
					continue
				}
				
				// Check for subdirectories in community
				if dirName == "community" || dirName == "Global" {
					subDirPath := filepath.Join(dirName, file.Name())
					subFiles, err := ioutil.ReadDir(subDirPath)
					if err != nil {
						continue
					}
					
					for _, subFile := range subFiles {
						if !subFile.IsDir() && strings.HasSuffix(subFile.Name(), ".gitignore") {
							fullPath := filepath.Join(subDirPath, subFile.Name())
							templateName := strings.TrimSuffix(subFile.Name(), ".gitignore")
							qualifiedName := file.Name() + "/" + templateName
							t.loadTemplate(fullPath, qualifiedName)
						}
					}
				}
				
				continue
			}
			
			if strings.HasSuffix(file.Name(), ".gitignore") {
				var templatePath string
				if dirName == "." {
					templatePath = file.Name()
				} else {
					templatePath = filepath.Join(dirName, file.Name())
				}
				
				templateName := strings.TrimSuffix(file.Name(), ".gitignore")
				if dirName != "." {
					templateName = dirName + "/" + templateName
				}
				t.loadTemplate(templatePath, templateName)
			}
		}
	}
	
	return nil
}

// loadTemplate reads a template file and adds it to the templates map
func (t *Templates) loadTemplate(path, name string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	
	// Handle template references (e.g., C++.gitignore inside Fortran.gitignore)
	trimmedContent := strings.TrimSpace(string(content))
	if !strings.Contains(trimmedContent, "\n") && strings.HasSuffix(trimmedContent, ".gitignore") {
		referencedTemplate := strings.TrimSuffix(trimmedContent, ".gitignore")
		if referenced, ok := t.templates[referencedTemplate]; ok {
			t.templates[name] = referenced
			return
		}
	}
	
	t.templates[name] = string(content)
}

// GetTemplate returns the template for the given framework
func (t *Templates) GetTemplate(framework string) (string, bool) {
	// Try exact match
	if template, ok := t.templates[framework]; ok {
		return template, true
	}
	
	// Try case-insensitive match
	lowerFramework := strings.ToLower(framework)
	for name, template := range t.templates {
		if strings.ToLower(name) == lowerFramework {
			return template, true
		}
	}
	
	// No match found
	return "", false
}

// ListTemplates returns a list of all available templates
func (t *Templates) ListTemplates() []string {
	var templates []string
	for name := range t.templates {
		if !strings.Contains(name, "/") {
			templates = append(templates, name)
		}
	}
	return templates
}

// WriteGitignore writes the gitignore template to the specified file
func WriteGitignore(template, outputPath string) error {
	return ioutil.WriteFile(outputPath, []byte(template), 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gitignore <framework-name>")
		fmt.Println("Or use: gitignore list - to list all available templates")
		os.Exit(1)
	}
	
	// Initialize and load templates
	templates := NewTemplates()
	err := templates.LoadTemplates()
	if err != nil {
		fmt.Printf("Error loading templates: %v\n", err)
		os.Exit(1)
	}
	
	command := strings.ToLower(os.Args[1])
	
	if command == "list" {
		// List all available templates
		templateList := templates.ListTemplates()
		fmt.Println("Available templates:")
		for _, t := range templateList {
			fmt.Printf("- %s\n", t)
		}
		return
	}
	
	// Get the requested template
	template, found := templates.GetTemplate(os.Args[1])
	if !found {
		fmt.Printf("No template found for '%s'\n", os.Args[1])
		fmt.Println("Try 'gitignore list' to see all available templates")
		os.Exit(1)
	}
	
	// Write to .gitignore in current directory
	outputPath := ".gitignore"
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}
	
	// Check if file exists and confirm overwrite
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("File '%s' already exists. Overwrite? (y/n): ", outputPath)
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled")
			os.Exit(0)
		}
	}
	
	err = WriteGitignore(template, outputPath)
	if err != nil {
		fmt.Printf("Error writing gitignore: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Successfully created gitignore for '%s' at '%s'\n", os.Args[1], outputPath)
}
