# fileserve

`$ go get github.com/rbxb/fileserve`

## Usage example

```go
package main

import (
	"net/http"
	"github.com/rbxb/fileserve"
)

func main() {
	server := fileserve.NewServer("C:/mywebsite/root", nil)
	http.ListenAndServe(":8080", server)
}
```
  
See the `exmaple/directoryExample` folder for an example of how to setup a directory.

## Tag files
By default, these are files named `_tags.txt`.  
When a file is requested, the server will first check for a tag file in the parent directory of the requested file. The server will look for a tag where the first value matches the requested file, then modify the request using a handler function. Only one tag handler will be executed per request.  
Tag files can be changed without restarting the server.

## Tags

**`#ignore value1`**  
Trying to access the file named `value1` will result in a 404 error. Use this to prevent people from accessing certain files.  
  
**`#pseudo value1 value2`**  
Getting the file named `value1` will instead get the file named `value2`.  
  
**`#redirect value1 value2`**  
A request for `value1` will redirect the request to the URL `value2`.  

## Other notes

- It (should be) impossible to retrieve files that are outside of the root directory.  
- You may have a tag file outside of the root directory in the same directory that root is in.  
- Only one tag handler will be executed per request. You *cannot* make an infinite loop of `#pseudo`, though you *can* make a redirect loop. 
- When writing tag files, you don't have to specify the tag on every line. Make sure to use linebreaks, e.g.
```
#pseudo
	home home.html
	about about.html
	contact contact.html
#ignore
	secrets.txt
```
- You can create custom tags too. See `request` and `reverseproxy` as examples.
- You can overwrite the default tags (`#ignore`, `#pseudo`, `#redirect`, and `#default`) with your own.
- You can use `*` to select all files, e.g. `#ignore *`