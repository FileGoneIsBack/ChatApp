package translator

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Translate(source, sourceLang, targetLang string) (string, error) {
	var text []string
	var result []interface{}

	encodedSource := encodeURI(source)
	url := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=" +
		sourceLang + "&tl=" + targetLang + "&dt=t&q=" + encodedSource

	r, err := http.Get(url)
	if err != nil {
		return "", errors.New("Error getting translate.googleapis.com")
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("Error reading response body")
	}
	fmt.Printf("Translating text: '%s' from '%s' to '%s'", source, sourceLang, targetLang)

	// Log the raw response body for debugging
	fmt.Println("API Response Body:", string(body))

	bReq := strings.Contains(string(body), `<title>Error 400 (Bad Request)`)
	if bReq {
		return "", errors.New("Error 400 (Bad Request)")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", errors.New("Error unmarshaling data")
	}

	// Ensure result[0] exists and is a slice of interface{}
	if len(result) == 0 || result[0] == nil {
		return "", errors.New("No translated data in response")
	}

	inner, ok := result[0].([]interface{})
	if !ok {
		return "", errors.New("Unexpected data structure")
	}

	for _, slice := range inner {
		if translatedSlice, ok := slice.([]interface{}); ok && len(translatedSlice) > 0 {
			if translatedText, ok := translatedSlice[0].(string); ok {
				text = append(text, translatedText)
			} else {
				return "", errors.New("Unexpected data structure in translated text")
			}
		}
	}

	cText := strings.Join(text, "")
	return cText, nil
}

