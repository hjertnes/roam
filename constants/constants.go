// Package constants contains various constants
package constants

import "regexp"

// Version is the version of roam.
const Version = "0.3.4"

const LastMigration = 3

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

const DefaultTemplate = `---
title: "$$TITLE$$"
private: false
---
`

const DailyTemplate = `---
title: "$$DATE$$"
private: false
---
`

const PublishTemplate = `<!DOCTYPE html>
<html lang="en-us">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>hjertnes.wiki</title>
	<link rel="stylesheet" href="/style.css">
</head>
<body>
    <header>
      <h1><a href="https://hjertnes.wiki/">hjertnes.wiki </a></h1>
    </a>
    </header>
	<main>
<article>
<h1>$$TITLE$$</h1>
<div>$$TEXT$$</div>
<div<$$BACKLINKS$$</div>

</article>
</main>
</body>
</html>`