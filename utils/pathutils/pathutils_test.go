package pathutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPath(t *testing.T){
	p := New("a/b/c.md")
	assert.Equal(t, "c.md", p.GetFilename())
	assert.Equal(t, "a/b", p.GetParent())

	p = New("a/b/c")
	assert.Equal(t, "c", p.GetFilename())
	assert.Equal(t, "a/b", p.GetParent())

}
