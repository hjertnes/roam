// Package models contain shared models
package models

import "time"

// TemplateFile is the type for each item in the list of templates in the config.
type TemplateFile struct {
	Filename string
	Title    string
}

type Publish struct {
	FilesToCopy []string

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
	Publish *Publish
}

type Log struct {
	ID string `db:"id"`
	Timestamp time.Time `db:"created_at"`
	Failure bool `db:"failure"`
}

// Folder is the database model for a folder.
type Folder struct {
	ID   string `db:"id"`
	Path string `db:"path"`
}

// Frontmatter is the type for the metadata in a file.
type Frontmatter struct {
	Title   string `fm:"title"`
	Private bool   `fm:"private"`
	Content string `fm:"content"`
}

// ImportFrontmatter is the metadata for import files.
type ImportFrontmatter struct {
	Title   string `fm:"title"`
	Private bool   `fm:"private"`
	Path    string `fm:"path"`
	Content string `fm:"content"`
}

// File is the database model for a file or note.
type File struct {
	ID      string `db:"id"`
	Title   string `db:"title"`
	Private bool   `db:"private"`
	Path    string `db:"path"`
}

// Choice is the type for the options.
type Choice struct {
	Title string
	Value string
}