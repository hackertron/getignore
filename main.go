package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Templates struct to hold all the gitignore templates
type Templates struct {
	templates map[string]string
}

// TemplateFile represents a file from GitHub API
type TemplateFile struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	DownloadURL string `json:"download_url"`
	Type        string `json:"type"`
}

// NewTemplates creates a new Templates instance
func NewTemplates() *Templates {
	return &Templates{
		templates: make(map[string]string),
	}
}

// getTemplatesDir returns the path to the templates directory
func getTemplatesDir() (string, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Create .gitignore-cli directory in user's home if it doesn't exist
	templatesDir := filepath.Join(homeDir, ".gitignore-cli")
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		err = os.Mkdir(templatesDir, 0755)
		if err != nil {
			return "", err
		}
	}

	return templatesDir, nil
}

// LoadTemplates loads gitignore templates from local storage
func (t *Templates) LoadTemplates() error {
	// Get templates directory
	templatesDir, err := getTemplatesDir()
	if err != nil {
		return fmt.Errorf("error getting templates directory: %v", err)
	}

	// Create templates directory if it doesn't exist
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		err = os.Mkdir(templatesDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating templates directory: %v", err)
		}
	}

	// Load templates from files
	err = loadTemplatesFromDir(t, templatesDir, "")
	if err != nil {
		return fmt.Errorf("error loading templates from directory: %v", err)
	}

	return nil
}

// loadTemplatesFromDir loads templates from a directory and its subdirectories
func loadTemplatesFromDir(t *Templates, dir, prefix string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", dir, err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".gitignore") {
			templatePath := filepath.Join(dir, file.Name())
			templateName := strings.TrimSuffix(file.Name(), ".gitignore")
			if prefix != "" {
				templateName = prefix + "/" + templateName
			}
			t.loadTemplate(templatePath, templateName)
		} else if file.IsDir() {
			// Handle subdirectories
			subDir := filepath.Join(dir, file.Name())
			subPrefix := file.Name()
			if prefix != "" {
				subPrefix = prefix + "/" + subPrefix
			}
			loadTemplatesFromDir(t, subDir, subPrefix)
		}
	}

	return nil
}

// DownloadSingleTemplate downloads a specific template from GitHub
func DownloadSingleTemplate(framework string) (string, error) {
	templatesDir, err := getTemplatesDir()
	if err != nil {
		return "", fmt.Errorf("error getting templates directory: %v", err)
	}

	// Try different locations where the template might be
	possiblePaths := []string{
		framework + ".gitignore",                // Root directory
		"Global/" + framework + ".gitignore",    // Global directory
		"community/" + framework + ".gitignore", // Community directory
	}

	// Try to find subdirectories in community
	communityURL := "https://api.github.com/repos/github/gitignore/contents/community"
	resp, err := http.Get(communityURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		var files []TemplateFile
		err = json.NewDecoder(resp.Body).Decode(&files)
		if err == nil {
			for _, file := range files {
				if file.Type == "dir" {
					possiblePaths = append(possiblePaths,
						"community/"+file.Name+"/"+framework+".gitignore")
				}
			}
		}
	}

	// Try each possible location
	for _, path := range possiblePaths {
		downloadURL := "https://raw.githubusercontent.com/github/gitignore/main/" + path
		resp, err := http.Get(downloadURL)
		if err != nil {
			continue
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			continue
		}

		// We found the template!
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("error reading template content: %v", err)
		}

		// Save the template locally for future use
		dirPath := filepath.Dir(filepath.Join(templatesDir, path))
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				return "", fmt.Errorf("error creating directory: %v", err)
			}
		}

		templatePath := filepath.Join(templatesDir, path)
		err = ioutil.WriteFile(templatePath, content, 0644)
		if err != nil {
			return "", fmt.Errorf("error saving template: %v", err)
		}

		return string(content), nil
	}

	return "", fmt.Errorf("template not found")
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

// downloadTemplates downloads templates from GitHub
func downloadTemplates(templatesDir string) error {
	baseURL := "https://api.github.com/repos/github/gitignore/contents"

	// Download root templates
	err := downloadTemplatesFromPath(baseURL, "", templatesDir)
	if err != nil {
		return err
	}

	// Download Global templates
	err = downloadTemplatesFromPath(baseURL+"/Global", "Global", templatesDir)
	if err != nil {
		return err
	}

	// Download community templates
	return downloadTemplatesFromPath(baseURL+"/community", "community", templatesDir)
}

// downloadTemplatesFromPath downloads templates from a specific GitHub path
func downloadTemplatesFromPath(url, prefix, templatesDir string) error {
	// Get directory listing from GitHub
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download templates, status code: %d", resp.StatusCode)
	}

	var files []TemplateFile
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		return err
	}

	// Create subdirectory if needed
	if prefix != "" {
		subDir := filepath.Join(templatesDir, prefix)
		if _, err := os.Stat(subDir); os.IsNotExist(err) {
			err = os.Mkdir(subDir, 0755)
			if err != nil {
				return err
			}
		}
	}

	// Download each file
	for _, file := range files {
		if file.Type == "file" && strings.HasSuffix(file.Name, ".gitignore") {
			var targetPath string
			if prefix == "" {
				targetPath = filepath.Join(templatesDir, file.Name)
			} else {
				targetPath = filepath.Join(templatesDir, prefix, file.Name)
			}

			// Download the file
			err = downloadFile(file.DownloadURL, targetPath)
			if err != nil {
				fmt.Printf("Warning: failed to download %s: %v\n", file.Name, err)
				continue
			}

			// Give some feedback on progress
			fmt.Printf("Downloaded %s\n", file.Name)

			// Adding a small delay to avoid hitting rate limits
			time.Sleep(100 * time.Millisecond)
		} else if file.Type == "dir" && prefix == "community" {
			// For community subdirectories, we need to download their contents too
			subDirURL := url + "/" + file.Name
			subDirPrefix := prefix + "/" + file.Name

			// Create subdirectory
			subDir := filepath.Join(templatesDir, subDirPrefix)
			if _, err := os.Stat(subDir); os.IsNotExist(err) {
				err = os.Mkdir(subDir, 0755)
				if err != nil {
					return err
				}
			}

			// Download templates from subdirectory
			err = downloadTemplatesFromPath(subDirURL, subDirPrefix, templatesDir)
			if err != nil {
				fmt.Printf("Warning: failed to download from %s: %v\n", subDirURL, err)
			}
		}
	}

	return nil
}

// downloadFile downloads a file from URL to the specified local path
func downloadFile(url, targetPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	// Create the target file
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy content from response to file
	_, err = io.Copy(out, resp.Body)
	return err
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
		templates = append(templates, name)
	}
	return templates
}

// WriteGitignore writes the gitignore template to the specified file
func WriteGitignore(template, outputPath string) error {
	return ioutil.WriteFile(outputPath, []byte(template), 0644)
}

// printHelp prints the help information
func printHelp() {
	fmt.Println("Gitignore Generator - A tool to create .gitignore files for your projects")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  gitignore <command> [arguments] [options]")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  <framework-name>     Generate a .gitignore file for the specified framework")
	fmt.Println("  list                 List all available templates")
	fmt.Println("  update               Update templates from GitHub")
	fmt.Println("  download-all         Download all templates from GitHub")
	fmt.Println("  clean                Remove all locally stored templates")
	fmt.Println("  help, -h, --help     Show this help message")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --download-all       When used with a framework name, will download all templates")
	fmt.Println("                       instead of just the requested one")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  gitignore Go                 Create a .gitignore file for Go")
	fmt.Println("  gitignore Python output.txt  Create a Python .gitignore file named output.txt")
	fmt.Println("  gitignore list               Show all available templates")
	fmt.Println("  gitignore download-all       Download all templates from GitHub")
	fmt.Println()
	fmt.Println("The templates are stored in ~/.gitignore-cli directory.")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	command := strings.ToLower(os.Args[1])

	// Check if download-all flag is present
	downloadAllFlag := false
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--download-all" {
			downloadAllFlag = true
			break
		}
	}

	// Handle help command
	if command == "help" || command == "-h" || command == "--help" {
		printHelp()
		return
	}

	// Handle clean command
	if command == "clean" {
		templatesDir, err := getTemplatesDir()
		if err != nil {
			fmt.Printf("Error getting templates directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Print("Are you sure you want to remove all templates? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled")
			return
		}

		err = os.RemoveAll(templatesDir)
		if err != nil {
			fmt.Printf("Error removing templates: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Templates successfully removed")
		return
	}

	// Handle download-all command
	if command == "download-all" {
		templatesDir, err := getTemplatesDir()
		if err != nil {
			fmt.Printf("Error getting templates directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Downloading all templates from GitHub...")
		err = downloadTemplates(templatesDir)
		if err != nil {
			fmt.Printf("Error downloading templates: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("All templates downloaded successfully!")
		return
	}

	if command == "update" {
		// Force update templates
		templatesDir, err := getTemplatesDir()
		if err != nil {
			fmt.Printf("Error getting templates directory: %v\n", err)
			os.Exit(1)
		}

		// Remove existing templates
		err = os.RemoveAll(templatesDir)
		if err != nil {
			fmt.Printf("Error removing existing templates: %v\n", err)
			os.Exit(1)
		}

		// Create templates directory again
		err = os.Mkdir(templatesDir, 0755)
		if err != nil {
			fmt.Printf("Error creating templates directory: %v\n", err)
			os.Exit(1)
		}

		// Download templates
		fmt.Println("Updating templates from GitHub...")
		err = downloadTemplates(templatesDir)
		if err != nil {
			fmt.Printf("Error downloading templates: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Templates updated successfully!")
		return
	}

	// Initialize and load templates
	templates := NewTemplates()
	err := templates.LoadTemplates()
	if err != nil {
		fmt.Printf("Error loading templates: %v\n", err)
		os.Exit(1)
	}

	if command == "list" {
		// List all available templates
		templateList := templates.ListTemplates()
		fmt.Printf("Available templates (%d):\n", len(templateList))

		// Group templates by directory
		groups := make(map[string][]string)
		for _, t := range templateList {
			if strings.Contains(t, "/") {
				parts := strings.SplitN(t, "/", 2)
				groups[parts[0]] = append(groups[parts[0]], parts[1])
			} else {
				groups["Main"] = append(groups["Main"], t)
			}
		}

		// Print templates by group
		for group, templates := range groups {
			fmt.Printf("\n%s:\n", group)
			for _, t := range templates {
				fmt.Printf("  - %s\n", t)
			}
		}
		return
	}

	// Get the requested template
	templateContent := ""
	found := false

	// If download-all flag is present, always download all templates
	if downloadAllFlag {
		templatesDir, err := getTemplatesDir()
		if err != nil {
			fmt.Printf("Error getting templates directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Downloading all templates from GitHub...")
		err = downloadTemplates(templatesDir)
		if err != nil {
			fmt.Printf("Error downloading templates: %v\n", err)
			os.Exit(1)
		}

		// Reload templates
		templates = NewTemplates()
		err = templates.LoadTemplates()
		if err != nil {
			fmt.Printf("Error loading templates: %v\n", err)
			os.Exit(1)
		}

		templateContent, found = templates.GetTemplate(os.Args[1])
	} else {
		// Try to get template from local cache first
		templateContent, found = templates.GetTemplate(os.Args[1])

		// If not found locally, try to download just this template
		if !found {
			fmt.Printf("Template for '%s' not found locally. Trying to download...\n", os.Args[1])
			var err error
			templateContent, err = DownloadSingleTemplate(os.Args[1])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				fmt.Println("Try 'gitignore list' to see available templates")
				fmt.Println("or 'gitignore download-all' to download all templates")
				os.Exit(1)
			}
			found = true
			fmt.Printf("Template for '%s' downloaded successfully\n", os.Args[1])
		}
	}

	if !found {
		fmt.Printf("No template found for '%s'\n", os.Args[1])
		fmt.Println("Try 'gitignore list' to see all available templates")
		os.Exit(1)
	}

	// Write to .gitignore in current directory
	outputPath := ".gitignore"
	if len(os.Args) > 2 && os.Args[2] != "--download-all" {
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

	err = WriteGitignore(templateContent, outputPath)
	if err != nil {
		fmt.Printf("Error writing gitignore: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully created gitignore for '%s' at '%s'\n", os.Args[1], outputPath)
}
