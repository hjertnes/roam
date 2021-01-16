# README
### Getting started 
#### Dependencies
- Golang
- PostgreSQL

I don't provide any binaries for this project at this point, so you need to install golang to build it. And I use PostgreSQL as the database so you also need to have that installed. There is a section about why further down on this page.  

#### Build & installation
Run `make install`and move the roam binary to somewhere on your path. 

#### Configuration
0. Set the ROAM env variable to where you want to store your files
1. Create a database
2. Run `roam init` this will create configuration files
3. Run `roam edit config` it will open the configuration in EDITOR and you at least need to set the databaseconnection string. Mine looks like this `host=localhost port=5432 dbname=wiki user=hjertnes`. It is a standard postgres connection string. 
4. Run `roam migrate` to create database tables and all that fun stuff
5. Run `roam sync` this will create a cache of your data folder in the database. 

## Configuration files
All of the config stuff are saved inside .config in your ROAM dir. 

### config.yaml
#### The date or time formats
This is written in golang, and I use the time.Time package for date time stuff. For details on how to change this format look here: <https://golang.org/pkg/time/#pkg-constants>

There are three settings for this:
- timeformat
- dateformat
- datetimeformat

They are kind of obvious, but dateformat are used when you only need a date, time is used when you only need time and datetime are used when you need both.

#### databaseconnectionstring
This is used to connect to the database. It is a standard postgres connection string.

#### defaultfileextension
Not in use yet. Might be removed in the future 
#### version
This tells what version of the config file. Not really relevant yet, but will be once I make incompatible changes

#### templates
This is a list of templates, they are used when you use the create command. You are free to add your own. You need to create another file in the .config/templates and add a entry to the templates list

More info on them in the templates section.
### templates
The Templates are just a simple text file. It supports the following variables:
- $$TIME$$
- $$DATE$$
- $$DATETIME$$
- $$TITLE$$

The time formats from your main config are used

### publish
Contains files used to publish a HTML verison of your roam. 

#### template.html
This is the HTML structure everything will be placed inside. 

The supported variables:
- $$TITLE$$
- $$TEXT$$
- $$BACKLINKS$$

They're probably obvious but title is the title of the page (usually note), text is the page content and bakclinks are files that have links to this file

In pages where one of them isn't relevant it will be replaced with a empty string. 

#### Other files
Files that ends with .jpg, .jpeg, .png, .css or .js will be included in output

## Formats
### Note format
This is the format of the notes
```
---
title: "Title"
private false
---

Note text
```

This is probably obvious but the Title is the title, private as described in the section about publish tells if it should be included by default in publishing or not. Only matters if you are going publish your database. 
### Import Format
This is the import format
```
---
title: "Note 1"
private: false
path: "a/b/c1.md"
---

Note text
---
title: "Note 2"
private: false
path: "a/b/c2.md"
---

Note text
```

The path is relative to the root of your ROAM dir. The format is more or less the same as the note format except that it also has a path property to tell where to save the file. And you can have as many notes as you want to in the import. 

## Usage
This section is structured a bit weird. But I have placed example commands at top and then a description below them 

### Clear
`roam clear database`
`roam clear config`

Used to remove database content or the config dir
### Log
`roam log clear`
`roam log list`

Lists or clears the log written after each sync. It tells if the sync was successful or not.
### Sync
`roam sync`

You should run the sync should be run after each time a file changes or a new file are added etc. I run it in a cron job. It is'nt cron because macOS, but the same concept. 
### Init
`roam init`

Creates the initial configuration and related files. It doesn't do anything if the files exist
### Migrate
`roam migrate`

Applies database migrations, needs to be executed when you set things up and when there are database changes
### Help
`roam help` or `roam whatever` (assuming I haven't created a command called whatever and forgotten to update this document).

Shows a (currently) crappy help message. This will also be shown when the commands you entered wasn't valid

### Version
`roam version` prints the current installed version

### Create
`roam create some/file.md`

The path in the command isrelative to the root of your roam directory. It will create a new document. Letting you specify what title you want it to have and what template it is based on. Look in the configuration section for more details on templates

Note, the directories doesn't have to exist, the command will create them.

Currently this only creates notes, but this will be extended to creating templates etc in the future 
### Edit
`roam edit config`

Opens the configuration in EDITOR. In the future this will be extended to more features like editing templates etc
### Import
`roam import ~/file/to/import.md`
`roam import ~/file/to/import.md --dry`

Look in the formats section for details on the format. But it will read the file and create notes based on it, if there isn't a file like it that exist yet. 

Use the --dry option to test before writing it to the database

### Daily
`roam daily`
`roam daily [date-string]`

This creates a new "daily note" in the Daily Notes directory. It uses the Daily Notes template
### Diagnostic
`roam diagnostic`

It runs a similar process to the sync, except that it will only warn about problems with your notes. Like invalid links or front matter problems etc. 

### Report
`roam report`

Shows a list of notes and their links and backlinks

### Find
`roam find`
`roam find query`
`roam find find`
`roam find edit`
`roam find --backlinks`
`roam find --links`

Find is the command I use the most, it is useful for searching after notes

Query (the default) will just print the matches out, or print out the list of links/backlinks if you add that parameter
Edit will let you search, select (and then select a link if --links or --backlinks is set) and then
View works the same way as edit except that it renders the note in the terminal instead of opening it in EDITOR
### Publish
`roam publish`
This will publish everything where private is false

`roam publish --include-private`
This will publish everything

`roam publish someOtherFolder`
This will make it publish to someOtherFolder instead of output

One of the core reasons I made this is that I want a good local "roam" like or wiki like thing. But I also want part of it to be available as a website. This does that. 

It takes your notes, and by default filters out all the notes with private set to true on the front matter (can be omitted by a flag) and it then it will use the template in .config/publish to build HTML. It will also copy any file in that directory that ends with:
- css
- js
- jpg
- jpeg
- png

This will most likely become a configuration in the future 

### Stats
`roam stats`

Stats just shows you how many files, links and backlinks are in your database
## Future
[TODO.md](TODO.md) contains information about stuff I want to fix and make. 
## Misc
### Why Postgres
This project could have used SQLite, but I have chosen to use PostgreSQL for a few different reasons.

1. I'm going to have the sync command running on a schedule in the background on my machines, and it is a lot easier to do that without having to deal with issues related to multiple processes trying to write to the same sqlite file 
2. I like the Postgres version of SQL much more than SQLite
3. I really hate the SQLite team for their so called code of conduct 

#### Why a database
I have a database with close to 1500 notes, to search them all would be dreadfully slow. And a few other things are also easier with a database. Like getting a list of backlinks etc. 
## Development

This is all written in Golang without a lot of deps. If you want a seperate roam for development you can create a file anywhere called `.roam` just put a path inside it. It will be checked before reading the ROAM env variable. 