package main

import (
	"os"
	"testing"
)

func TestTemplatesLoading(t *testing.T) {
	templates := NewTemplates()
	err := templates.LoadTemplates()
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}
	
	// Check if we have at least some templates loaded
	if len(templates.templates) == 0 {
		t.Error("No templates were loaded")
	}
	
	// Test a few common templates that should be available
	commonTemplates := []string{"Go", "Python", "Node"}
	for _, name := range commonTemplates {
		_, found := templates.GetTemplate(name)
		if !found {
			t.Errorf("Template for '%s' not found", name)
		}
	}
}

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
