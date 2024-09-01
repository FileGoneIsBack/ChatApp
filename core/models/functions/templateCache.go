package functions

import (
    "path/filepath"
    "text/template"
    "login/core/models/validation"
    "strings"
    "os"
    "log"
    "sync"
	"errors"
)
var ErrTemplateNotFound = errors.New("template not found in cache")
// TemplateCache stores the parsed templates
type TemplateCache struct {
    templates map[string]*template.Template
    mu        sync.RWMutex
}

// NewTemplateCache initializes a new TemplateCache.
func NewTemplateCache() *TemplateCache {
    return &TemplateCache{
        templates: make(map[string]*template.Template),
    }
}

func (tc *TemplateCache) LoadTemplates(directories ...string) error {
    tc.mu.Lock()
    defer tc.mu.Unlock()

    for _, dir := range directories {
        err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return err
            }
            if !info.IsDir() && filepath.Ext(path) == ".html" {
                // Load and parse the template file
                tmpl, err := template.New(filepath.Base(path)).ParseFiles(path)
                if err != nil {
                    return err
                }

                // Store the main template and its named templates
                for _, t := range tmpl.Templates() {
                    name := t.Name()
                    tc.templates[name] = t
                }
            }
            return nil
        }) 
        if err != nil {
            return err
        }
    }

    Sanitize.LogMessage("success", "All templates loaded and cached.", nil)
    return nil
}


// GetTemplate retrieves a specific template from the cache.
func (tc *TemplateCache) GetTemplate(name string) (*template.Template, bool) {
    tc.mu.RLock()
    defer tc.mu.RUnlock()

    tmpl, found := tc.templates[name]
    if !found {
        log.Printf("Template not found: %s", name)
    }
    return tmpl, found
}

// GetCachedTemplates retrieves both the main template and all cached partials.
func (tc *TemplateCache) GetCachedTemplates(mainFile string) (*template.Template, error) {
    tc.mu.RLock()
    defer tc.mu.RUnlock()

    // Retrieve the main template from the cache
    mainTemplate, found := tc.templates[mainFile]
    if !found {
        return nil, ErrTemplateNotFound
    }

    // Clone the main template so we can add partials
    clonedTemplate, err := mainTemplate.Clone()
    if err != nil {
        return nil, err
    }

    // Attach all cached templates to the cloned template
    for name, tmpl := range tc.templates {
        if name != mainFile {
            _, err := clonedTemplate.AddParseTree(name, tmpl.Tree)
            if err != nil {
                return nil, err
            }
        }
    }

    return clonedTemplate, nil
}

func (tc *TemplateCache) ListTemplates() {
    tc.mu.RLock()
    defer tc.mu.RUnlock()

    seenTemplates := make(map[string]bool)

    for name := range tc.templates {
        // Strip the file extension and get the base name of the template
        templateName := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
        
        // Log the template name only if it hasn't been seen yet
        if !seenTemplates[templateName] {
            Sanitize.LogMessage("info", "- "+templateName, "")
            seenTemplates[templateName] = true
        }
    }
}
