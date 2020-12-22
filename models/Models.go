package models

type TemplateFile struct {
	Filename string
	Title string
}

type Configuration struct {
	TimeFormat string
	DateFormat string
	DateTimeFormat string
	DatabaseConnectionString string
	DefaultFileExtension string
	Version int
	Templates []TemplateFile
}

type Fm struct {
	Title string `fm:"title"`
	Private bool `fm:"private"`
	Content string `fm:"content"`
}

type File struct {
	Title string `db:"title"`
	Path string `db:"path"`
}