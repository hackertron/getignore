package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestCaseInsensitiveMatching tests case-insensitive template matching
func TestCaseInsensitiveMatching(t *testing.T) {
	templates := NewTemplates()

	// Add a test template
	templates.templates["Go"] = "test content"

	// Test case-insensitive matching
	content, found := templates.GetTemplate("go")
	if !found {
		t.Error("Case-insensitive template matching failed")
	}

	if content != "test content" {
		t.Error("Template content mismatch")
	}
}

// TestTemplateLoadingFromFile tests template loading from file
func TestTemplateLoadingFromFile(t *testing.T) {
	// Create a temporary directory for test templates
	tempDir, err := ioutil.TempDir("", "gitignore-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test template file
	testContent := "# Test gitignore template\n*.log\n.DS_Store\n"
	testFile := filepath.Join(tempDir, "Test.gitignore")
	err = ioutil.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a subdirectory with a template
	subDir := filepath.Join(tempDir, "Global")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	subContent := "# Global test gitignore\n.vscode/\n*.tmp\n"
	subFile := filepath.Join(subDir, "Editor.gitignore")
	err = ioutil.WriteFile(subFile, []byte(subContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create subdir test file: %v", err)
	}

	// Create a templates instance and load from our test dir
	templates := NewTemplates()
	templates.loadTemplate(testFile, "Test")
	templates.loadTemplate(subFile, "Global/Editor")

	// Test if templates were loaded correctly
	testTemplate, found := templates.GetTemplate("Test")
	if !found {
		t.Error("Failed to get Test template")
	} else if testTemplate != testContent {
		t.Errorf("Template content mismatch. Expected '%s', got '%s'", testContent, testTemplate)
	}

	globalTemplate, found := templates.GetTemplate("Global/Editor")
	if !found {
		t.Error("Failed to get Global/Editor template")
	} else if globalTemplate != subContent {
		t.Errorf("Template content mismatch. Expected '%s', got '%s'", subContent, globalTemplate)
	}
}

// TestGetTemplatesDir tests the template directory logic
func TestGetTemplatesDir(t *testing.T) {
	// Just test that we get a non-empty string and no error
	templatesDir, err := getTemplatesDir()
	if err != nil {
		t.Fatalf("getTemplatesDir returned error: %v", err)
	}

	if templatesDir == "" {
		t.Error("getTemplatesDir returned empty string")
	}
}

// TestTemplateReference tests the template reference handling
func TestTemplateReference(t *testing.T) {
	templates := NewTemplates()

	// Add a template that will be referenced
	templates.templates["Referenced"] = "# Referenced template content"

	// Create a temporary file with a reference
	tempFile, err := ioutil.TempFile("", "gitignore-ref-*.gitignore")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write just a reference to another template
	_, err = tempFile.WriteString("Referenced.gitignore")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	// Load the template with reference
	templates.loadTemplate(tempFile.Name(), "Referencing")

	// Check if the reference was resolved
	content, found := templates.GetTemplate("Referencing")
	if !found {
		t.Error("Failed to get Referencing template")
	} else if content != "# Referenced template content" {
		t.Errorf("Template reference not resolved correctly. Got: '%s'", content)
	}
}
