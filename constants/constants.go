// Package constants contains various constants
package constants

import "regexp"

// Version is the version of roam.
const Version = "0.2.18"

// FolderPermission is the file system permission used when creating new folders.
const FolderPermission = 0700

// FilePermission is the file system permission used when creating new files.
const FilePermission = 0600

// Zero is a constant for the value 0.
const Zero = 0

// One is a constant for the value 1.
const One = 1

// Two is a constant for the value 2.
const Two = 2

// Three is a constant for the value 3.
const Three = 3

// NoteLinkRegexp is the regular expression for detecting links.
var NoteLinkRegexp = regexp.MustCompile(`\[\[([\d\w\s./]+)]]`)