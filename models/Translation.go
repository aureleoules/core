package models

// Translation struct
type Translation struct {
	LanguageName string `json:"language_name" bson:"language_name"`
	LanguageCode string `json:"language_code" bson:"language_code"`
	Content      string `json:"content" bson:"content"`
}
