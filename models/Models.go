package models

type Configuration struct {
	DatabaseConnectionString string
	DefaultFileExtension string
	Version int
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