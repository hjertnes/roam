# Changelog
## v0.3.12
- Removed useless test from diagnostic
- Re-written publish to not be a gigantic function
## v0.3.11
- Broken up create 
## v0.3.10
- Split up find
- Added check to create to make sure it never overwrites anything
## v0.3.9
- Moved TODO.md to my task manager
- Added a EditFile method to utils and moved duplicated code from other places
- Removed newline check from diagnostic
- Import: added feedback on how many notes imported
- Import: dry run checks if files exist
- Import: dry checks if the filename ends in .md
## v0.3.8
- Added a warning if more than 5 min since last syncx
## v0.3.7
- Moved the creation of commands to a factory
- Fixed some crashes 
- Added handling for a weird case where I returned nil,nil
## v0.3.6
- Make a proper type out of diagnostics
- Split it into not being a huge function
- Moved the command builders in main to a sep. file. 
- And renamed them to "buildX" from "getX". 
## v0.3.5
- Split help into multiple files
## v0.3.4
- Changed reports to be a proper type
## v0.3.3
- Re-written sync to not be a giant function
## v0.3.2
- Cleaned up the template module a little
## v0.3.1
- Added an abstraction around os.Args instead of using it directly
- Replaced most of the uses of os,Exit
- Added a check that warns if you have pending migrations
## v0.3.0
Just bumped it to v0.3.0 because this .22 version number is dumb
## v0.2.22
Added a config for what is copied to output during publish
## v0.2.21
Added commands to manage templates (CRUD)
## v0.2.20
- Published wrote html to .md instead of .html
- Fixed selectinput because markdown changed it to using 1 indexed lists
## v0.2.19
Removed most of the stuff in utils
## v0.2.18
- Added options to getlinks and getbacklinks to if private should be included or not
## v0.2.17
- Delete folders in sync
## v0.2.16
- Added log for sync
- Message before each command telling if there was a problem with the last sync
- Command to show and clear the log
## v0.2.15
- Fixed crashing bug in pathuti;s
## v0.2.14
- Added more checks to diagnostics
- Replaced some output with rendered markdown
## v0.2.13
Replaced various path functions with a pathtuils package
## v0.2.12
More validation on some sub commands that lacked it
## v0.2.11
Linting
## v0.2.10
Added commands for clearing all data from database or removing the config directory
## v0.2.9
Added changelog
## v0.2.8
Removed too many newlines in import and some input validation + a dry run option
## v0.2.7
Changed file/folder permission to the lowest functioning
## v0.2.6
Added version command
## v0.2.5
Added support for overriding roam dir to make development less risky
## v0.2.4
Fixed output in diagnostic
## v0.2.3
Added loaders anywhere that might be slow
## v0.2.2
Improved help command
## v0.2.1
Added readme
## v0.2.0
Re-written find to make it easier to use
## v0.1.3
Splitting commands into sep. packages and replace repetitive dal construction with a state object
## v0.1.2
Don't show folder if index is private
## v0.1.1
Include all css / js files in .config/publish in publish output
## v0.1.0
Cleanup
## v0.0.17
Nested index and control over if private is included or not
## v0.0.16
Include links and backlinks in publish
## v0.0.15
First version of publish to html
## v0.0.13
Fixing lint issues
## v.0.0.12
More cleanup and a report command that lists notes and its links / backlinks
## v.0.0.11
Cleanup
## v.0.0.10
Command for showing backlinks
## v.0.0.9
Added indexes to database
## v.0.0.8
Sync now writes links to database
## v.0.0.7
Added command called diagnostics that looks for problems in your notes
## v.0.0.6
Cleanup
## v.0.0.5
Command for showing basic stats
## v0.0.4
Command for creating a note
## v0.0.3
Command for finding a note
## v0.0.2
Cleaned everything up a little
## v0.0.1
Basic database setup, configuration etc