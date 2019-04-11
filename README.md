# fileserve

This package will help you quickly setup a basic server for serving webpages. The server uses configuration files (called tagfiles) to determine how requests are handled. These files can be modified live.

## Tutorial
[tutorial.md](tutorial.md)

## Basic Usage Example
```shell
$ go get github.com/rbxb/fileserve
```
```go
package main

import (
	"log"
	"net/http"
	"github.com/rbxb/fileserve"
)

func main() {
	server := fileserve.NewServer("./root", nil)
	log.Fatal(http.ListenAndServe(":8080", server))
}
```

## cmd/fileserve
A simple implementation of the package to get you started.

```shell
$ go install github.com/rbxb/fileserve/cmd/fileserve
$ fileserve -port :8080 -directory ./root
```

### Flags

#### `-port`
The port the fileserver runs on. (:8080)

#### `-directory`
The directory to serve files from. (./root)

## Tagfiles
- By default, these are files named `_tags.txt`.
- When a file is requested, the server will first check for a tagfile in the parent directory of the requested file. The server will look for a tag where the selector matches the name of the requested file. The request is handled based on the tag name.
- Only one tag handler function will be executed per request.
- The server always reads the tags in the same order that they are written in the tagfile.
- Tagfiles can be modified without restarting the server.

### Anatomy of a tag
```
         tag
┏━━━━━━━━━┻━━━━━━━━━╍┅
#pseudo home home.html
 ┃      ┃    ┗━━━┳━━╍┅
 name   ┃    arguments
        ┃
        selector
```
- The name determines which tag handler function will be executed.
  - The name is always immedietly preceded by a `#`.
- The selector determines which requests will match with the tag.
  - The server only looks at the last element of the requested path when comparing the request to the selector.
- The additional arguments are passed to the tag handler function as an array of strings.
  - You may have any number of additional arguments.
  - Arguments are separated by spaces.

## Default Tags

#### `#ignore`
- Takes no additional arguments.
- Trying to access the selected file will result in a 404 error.
- Use this to prevent people from accessing certain files.

#### `#pseudo`
- Takes one additional argument.
- A request for the selected file will instead get the file named in the first argument.

#### `#redirect`
- Takes one additional argument.
- A request for the selected file will redirect the request to the URL in the first argument.

#### `#default`
- Takes no additional arguments.
- Serves the file at the requested path.
- If a request is not caught by any tags in the tagfile, it will be handled by the default tag handler function.

## Advanced selectors
- You can use `*` to select all files, e.g. 
  - `#ignore *` will hide all the files in the directory.
- You can select requests where only the name needs to match or where only the extension needs to match, e.g. 
  - `#ignore secret.*` will ignore all files named `secret` disregarding the extension.
	- `#ignore *.txt` will ignore all files that have the `.txt` extension.

## Custom Tags
Learn how to write custom tag handler functions at [customtags.md](customtags.md).

## Other notes

- It (should be) impossible to retrieve files that are outside of the root directory.
- You may have a tag file outside of the root directory in the same directory that root is in. You can use this to select requests which have no path, e.g. a request to `localhost:8080` with no path could be selected using the name `root`.
- Since all paths are already relative to `root/`, do not include `root` in paths in tagfiles.
- Only one tag handler will be executed per request. You *cannot* make an infinite loop of `#pseudo`, though you *can* make a redirect loop.
- When writing tag files, you don't have to specify the tag on every line--use linebreaks and indents, e.g.
```
#pseudo
	home home.html
	about about.html
	contact contact.html
#ignore
	secrets.txt
```
- You can change the name of the tagfiles. By default the tagfile name is `_tags.txt`.
```go
Server.SetTagfileName(name string)
```
- You can overwrite the default tag handler functions (`#ignore`, `#pseudo`, `#redirect`, and `#default`) with your own handlers.