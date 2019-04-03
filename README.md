# fileserve

This package will help you quickly setup a basic server for serving webpages. The server uses configuration files (called tagfiles) to determine how requests are handled. These files can be modified live.

## Tutorial
[tutorial.md](tutorial.md)

## Usage Example
```shell
$ go get github.com/rbxb/fileserve
```
```go
package main

import (
	"net/http"
	"github.com/rbxb/fileserve"
)

func main() {
	server := fileserve.NewServer("./root", nil)
	http.ListenAndServe(":8080", server)
}
```

## Using the cmd
A simple implementation of the package to get you started.

```shell
$ go install github.com/rbxb/fileserve/cmd/fileserve
$ fileserve -port :8080 -directory ./root
```

### Flags

#### `-port`
The address:port the fileserver runs on. (:8080)

#### `-directory`
The directory to serve files from. (./root)

## Tagfiles
- By default, these are files named `_tags.txt`.
- When a file is requested, the server will first check for a tagfile in the parent directory of the requested file. The server will look for a tag where the first value matches the name of the requested file. The request is handled based on the tag name.
- Only one tag handler function will be executed per request.
- Tag files can be modified without restarting the server.

## Default Tags

#### `#ignore value1`
Trying to access the file named `value1` will result in a 404 error. Use this to prevent people from accessing certain files.

#### `#pseudo value1 value2`
Getting the file named `value1` will instead get the file named `value2`.

#### `#redirect value1 value2`
A request for `value1` will redirect the request to the URL `value2`.

## Custom Tags
Look at `request` and `reverseproxy` as examples.
Import the `request` package and attach the custom tag handler to your server like this:
```go
package main

import (
	"net/http"
	"github.com/rbxb/fileserve"
	"github.com/rbxb/fileserve/request"
)

func main() {
	server := fileserve.NewServer("./root", map[string]fileserve.TagHandler{
		"request": request.RequestTagHandler,
	})
	http.ListenAndServe(":8080", server)
}
```
Now **`#request value1 value2`** will perform a GET request to the URL in `value2` and respond with the response.
Try it by adding this to your tagfile:
```
#request google https://google.com
```
`localhost:8080/google` should show the html from `google.com`.

Write your own custom tag handlers in Go. A tag handler function follows this type:
```go
type TagHandler func(* Server, []string, http.ResponseWriter, * http.Request) error
```

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
- You can use `*` to select all files, e.g. `#ignore *` will hide all the files in the directory.
- You can change the name of the tagfiles. By default the tagfile name is `_tags.txt`.
```go
Server.SetTagfileName(name string)
```
- You can overwrite the default tag handler functions (`#ignore`, `#pseudo`, `#redirect`, and `#default`) with your own handlers.