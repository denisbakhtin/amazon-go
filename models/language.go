package models

//Language stores info about language
type Language struct {
	Model
	Code        string //ISO 639-1 language code
	Title       string
	NativeTitle string
}
