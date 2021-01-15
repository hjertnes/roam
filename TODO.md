# Plan to replace org-roam with something of my own
I feel more and more like I want to replace org-roam with something of my own

## TODO
- Create template
- Edit template
- Upgrade: migrates config
- Interactive mode that will give you a dropdown of options to replace it with and edit file option
- Add flag to getbacklinks and getlinks for private or not
- Add a configuration thing for what will be moved during publish e.g css/js etc
- partial sync
- Create is a mess
- Lint more

## Upcoming versions
### v0.2.12
- Error if --links and --backlinks are set
- Find crashes on empty database
- Some input validating msising on some commands
### v0.2.13
- Find things that belong in utilslib
- Utils should be broken up
- The path splitting stuff should use this regex: (?<!\\)/
### v0.2.14
- Find replace the ugly output with markdown
- Diagnostic v2: check that line 1 is --- and strip out multiple newlines in a row
### v0.2.15
- Find a way to tell the user there are problems they should deal with
### v0.2.16
- Delete folders in sync and links











