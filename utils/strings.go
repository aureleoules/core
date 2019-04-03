package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//GetStrings get strings
func GetStrings(r *http.Request) map[string]string {
	lang := GetLang(r)
	jsonFile, err := os.Open("./strings.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	var data map[string]map[string]string

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &data)

	if lang != "fr" {
		lang = "en"
	}

	return data[lang]
}

//GetLang returns user language
func GetLang(r *http.Request) string {
	return r.Header.Get("Language")
}
