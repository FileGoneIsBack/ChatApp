package functions

import (
	"fmt"
	"io"
	"log"
	"login/core/database"
	"login/core/models/server"
	"login/core/sessions"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	Route.NewSub(server.NewRoute("/avatar", func(w http.ResponseWriter, r *http.Request) {
		ok, user := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodPost {
			err := r.ParseMultipartForm(10 << 20) // 10 MB limit
			if err != nil {
				log.Printf("unable to parse form: %v", err)
				http.Error(w, "Unable to parse form", http.StatusBadRequest)
				return
			}

			file, fileHeader, err := r.FormFile("file")
			if err != nil {
				log.Printf("unable to get file from form: %v", err)
				http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
				return
			}
			defer file.Close()

			// Determine the file extension from MIME type
			ext := getFileExtension(fileHeader.Header.Get("Content-Type"))

			// Define the directory and create it if necessary
			dirPath := "assets/branding/media/avatar"
			err = os.MkdirAll(dirPath, os.ModePerm)
			if err != nil {
				log.Printf("unable to create directory: %v", err)
				http.Error(w, "Unable to create directory", http.StatusInternalServerError)
				return
			}

			// Create a unique file name with the correct extension
			fileName := fmt.Sprintf("%d_USER_%s%s", user.ID, user.Username, ext)
			filePath := filepath.Join(dirPath, fileName)

			// Create the new file
			outFile, err := os.Create(filePath)
			if err != nil {
				log.Printf("unable to create file: %v", err)
				http.Error(w, "Unable to create file", http.StatusInternalServerError)
				return
			}
			defer outFile.Close()
			database.Container.UpdateAvatar(fileName, user.ID)
			_, err = io.Copy(outFile, file)
			if err != nil {
				log.Printf("unable to save file: %v", err)
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}

			// Respond with success
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("File uploaded successfully"))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))
}

// getFileExtension determines the file extension based on MIME type
func getFileExtension(mimeType string) string {
	switch mimeType {
	case "image/png":
		return ".png"
	case "image/jpeg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	default:
		return ".bin" // Default to binary file if MIME type is unknown
	}
}
