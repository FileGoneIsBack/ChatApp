package functions

import (
    "fmt"
    "os"
    "net/http"
	"path/filepath"
    "login/core/models/validation"
)
// Render executes the main template and includes any required partials or components.
func Render(v interface{}, w http.ResponseWriter, file string, cache *TemplateCache) {
    baseName := filepath.Base(file)

    // Get the main template with cached partials
    tmpl, err := cache.GetCachedTemplates(baseName)
    if err != nil {
        Sanitize.LogMessage("error", "Error retrieving templates:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    // Execute the main template with the partials included
    err = tmpl.ExecuteTemplate(w, baseName, v)
    if err != nil {
        Sanitize.LogMessage("error", "Error executing template:", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

// Function to find the correct file extension
func GetUserAvatarFilePath(userID int, username string) (string, error) {
    baseDir := "./assets/branding/media/avatar/"
    baseName := fmt.Sprintf("%d_USER_%s", userID, username)
    extensions := []string{".png", ".jpg", ".gif"}

    for _, ext := range extensions {
        filePath := filepath.Join(baseDir, baseName+ext)

        if _, err := os.Stat(filePath); err == nil {
            urlPath := fmt.Sprintf("/_assets/media/avatar/%d_USER_%s%s", userID, username, ext)
            return urlPath, nil
        }
    }
    
    // If no file is found, return a default avatar or an error
    defaultAvatarPath := "/_assets/media/avatar/default_avatar.png"
    return defaultAvatarPath, nil
}

