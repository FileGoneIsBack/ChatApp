package apps

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"login/core/models/server"
	"login/core/models/translator"
	"login/core/sessions"
)

func init() {

	Route.NewSub(server.NewRoute("/translate", func(w http.ResponseWriter, r *http.Request) {
		ok, _ := sessions.IsLoggedIn(w, r)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if r.Method == http.MethodPost {
			var requestData struct {
				Q      string `json:"q"`
				Source string `json:"source"`
				Target string `json:"target"`
			}
			
			err := json.NewDecoder(r.Body).Decode(&requestData)
			if err != nil {
				log.Printf("unable to decode JSON: %v", err)
				http.Error(w, "Unable to decode JSON", http.StatusBadRequest)
				return
			}
		
			// Use the translator to translate the text
			translatedText, err := translator.Translate(requestData.Q, requestData.Source, requestData.Target)
			if err != nil {
				log.Printf("translation error: %v", err)
				http.Error(w, "Translation failed", http.StatusInternalServerError)
				return
			}
		
			// Respond with the translated text
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"translatedText": "%s"}`, translatedText)))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}		
	}))
}
