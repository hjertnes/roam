package pathutils

import (
	"strings"
)

type Path struct {
	path string
	elements []string
}

func New(path string) *Path{
	return &Path{
		path: path,
		elements: strings.Split(path, "/"),
	}
}

func (p *Path) GetPathWithoutFilename() string{
	return strings.Join(p.elements[0:len(p.elements)-1], "/")
}

func (p *Path) GetParent() string{
	return strings.Join(p.elements[0:len(p.elements)-2], "/")
}

func (p *Path) GetFilename() string{
	return p.elements[len(p.elements)-1]
}

func (p *Path) Destruct() (string, string){
	return p.GetPathWithoutFilename(), p.GetFilename()
}