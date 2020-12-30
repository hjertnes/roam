// Package models contain shared models
package models

// TemplateFile is the type for each item in the list of templates in the config.
type TemplateFile struct {
	Filename string
	Title    string
}

// Configuration is the type for the config file model.
type Configuration struct {
	TimeFormat               string
	DateFormat               string
	DateTimeFormat           string
	DatabaseConnectionString string
	DefaultFileExtension     string
	Version                  int
	Templates                []TemplateFile
}

type Folder struct {
	ID    string `db:"id"`
	Path  string `db:"path"`
}

// Frontmatter is the type for the metadata in a file.
type Frontmatter struct {
	Title   string `fm:"title"`
	Private bool   `fm:"private"`
	Content string `fm:"content"`
}

type PublishFrontmatter struct {
	Title   string `fm:"title"`
	Url string  `fm:"url"`
	Type string  `fm:"type"`
	Content string `fm:"content"`
}

// File is the database model for a file or note.
type File struct {
	ID    string `db:"id"`
	Title string `db:"title"`
	Private bool `db:"private"`
	Path  string `db:"path"`
}
