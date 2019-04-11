# Custom Tags

[README.md](README.md)

## Write a tag handler
Tag handler functions follow the `TagHandler` type:
```go
type TagHandler func(srvr * Server, args []string, w http.ResponseWriter, req * http.Request) error
```
`srvr` is the `Server` that called the tag handler function.
`args` are the arguments that were stored in the tagfile following the tag name and selector.

Look at default tag handler functions in [taghandlers.go](taghandlers.go) or the custom tag handlers in [request](request/request.go) or [reverseproxy](reverseproxy/reverseproxy.go) as examples of how to write a tag handler function.

From [taghandlers.go](taghandlers.go):
```go
func pseudo(srvr * Server, args []string, w http.ResponseWriter, req * http.Request) error {
	if len(args) < 1 {
		return ErrorNotEnoughArguments
	}
	name := filepath.Join(filepath.Dir(req.URL.Path), args[0])
	srvr.ServeFile(name, w, req)
	return nil
}
```

## Attach a tag handler to your server
Using `fileserve/request` as an example.
Import the `request` package and attach the custom tag handler to your server like this:
```go
package main

import (
	"log"
	"net/http"
	"github.com/rbxb/fileserve"
	"github.com/rbxb/fileserve/request"
)

func main() {
	server := fileserve.NewServer("./root", map[string]fileserve.TagHandler{
		"request": request.RequestTagHandler,
	})
	log.Fatal(http.ListenAndServe(":8080", server))
}
```
Try it by adding this to your tagfile:
```
#request google https://google.com
```
`localhost:8080/google` should show the html from `google.com`.

When attaching the tag handler function to your server, you can make the map key any value you want so long as it does not include `#` characters, spaces, or newline characters. The map key is the name that must appear in a tagfile following a `#` to indicate a tag name. The default tag handlers use the names `ignore`, `pseudo`, `redirect`, and `default`. You may overwrite these by attaching a tag handler function and using the same name in the key.