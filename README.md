# httpfilter

 - [Installation](#Installation)
 - [Tutorial](#Tutorial)
 - [How it Works](#How-it-Works)
 - [Usage Example](#Usage-Example)
 - [Standard Operators](#Standard-Operators)
 - [Writing Filters](#Writing-Filters)
 - [Attaching Additional Operators](#Attaching-Additional-Operators)
 - [Writing Operator Functions](#Writing-Operator-Functions)
 - [Fixed Filter File](#Fixed-Filter-File)

## Installation

```shell
$ go get github.com/rbxb/httpfilter
$ go install github.com/rbxb/httpfilter/cmd/httpfilter
```

## Tutorial

Playing around with the pre-made sample should help you get a feel for how httpfilter works.

[tutorial.md](./tutorial.md)

## How it Works

httpfilter uses scripts (*called filters*) nested in the website's static files directories to invoke Go functions that modify how the request is handled. For example, if you wanted to make the request path `/index.html` redirect to `/home.html`, you could use a filter like this:
```
#redirect index.html home.html
```
This filter file would be placed in the same directory as the supposed `index.html` file and the `home.html` file.

When the httpfilter server recieves a request for `/index.html`, it will first look for the filter in the same directory as the requested file. The server will read the filter from the top—down, searching for an entry where the **selector** matches the requested file name. If an entry's selector matches the requested file name, the server will call the Go function that is mapped to the **operator** of that entry. Any additional **arguments** are passed to the Go function as strings.

```
        entry
┏━━━━━━━━━┻━━━━━━━━━━━━━━━━╍┅
#redirect index.html home.html
 ┃          ┃       ┗━━━┳━━╍┅
 operator   ┃        arguments
            ┃
         selector
```

The httpfilter server passes the request and arguments to the operator function, which may write a response. The server will continue to read and execute entries from the filter until a response has been written.

If the server reaches the end of the filter and a response still hasn't been written, the server will call the `serve` operator, which will attempt to serve the file or respond with a `404 Not found` error.

## Usage Example

Create an httpfilter `Server` and pass it to `http.ListenAndServe`.

```go
package main

import (
	"log"
	"net/http"
	"github.com/rbxb/httpfilter"
)

func main() {
	server := httpfilter.NewServer("./root", "")
	log.Fatal(http.ListenAndServe(":8080", server))
}
```

## Standard Operators

#### `serve`

The `serve` operator attempts to serve the file named in the first argument or responds with a `404 Not found` error. E.g. this will serve the file `home.html` when the client requests `/home`:  
```
#serve home home.html
```
If the request is not fulfilled at the end of the filter file, the httpfilter serve will default to the serve operator to write a response.


#### `ignore`

The `ignore` operator responds with a `404 Not Found` error. Use this if you want to prevent access to a specific file, e.g.
```
#ignore secret.txt
```

#### `redirect`

The `redirect` operator redirects a request to the URL or path in the first argument. E.g. this will redirect `/index.html` to `/home`:
```
#redirect index.html home
```

## Writing Filters

 - Operators are always prefixed by a `#`.
 - Operators, selectors, and arguments are separated by spaces.
 - Entries are separated by line breaks.
 - The filter is read from the top—down and the server will never read upwards

### Selectors using `*`

You can use a `*` in the selector to select all queries.
 - `*` will match with all queries.
 - `a.*` will match with all queries where the name is `a` regardless of the extension.
 - `*.a` will match with all queries where the extension is `.a`.

For example, this filter will `ignore` all queries where the extension is `.txt`:
```
#ignore *.txt
```

### Selectors using `@`

Using the `@` symbol selects a request by its subdomain.   

For example, this filter will proxy the subdomain `service` to a local server and redirect the base domain to the `service` subdomain.
```
#proxy @service http://localhost:288
#redirect @ http://service.example.com
```

### Bulk Operator Syntax

If you have multiple entries that use the same operator repeatedly, e.g.
```
#serve home home.html
#serve about about.html
#serve contact contact.html
```
you can write the operator once and put the entries below it:
```
#serve
  home home.html
  about about.html
  contact contact.html
```

### Naming

Currently, filter files must be named `_filters.txt`.
Place the file in the directory that you want it to work in.
The httpfilter server will never serve a filter file to a client.

## Attaching Additional Operators

Non-standard operator functions can be attached to the server.
Pass them to `NewServer` as a `map[string]OpFunc`, where the map key is the operator name that should be used in the filter file to call the operator function.
```go
func NewServer(root string, ops ...map[string]OpFunc) * Server
```
Pass in your own operator functions:
```go
server := httpfilter.NewServer("./root", map[string]httpfilter.OpFunc{
		"myop": myOpFunc,
	})
```
The standard operators can be overwritten by passing in operators with the same key value.

## Writing Operator Functions

You can write your own operator functions and attach them to your server (see [Attaching Additional Operators](#Attaching-Additional-Operators)).

Operator functions follow this type:
```go
type OpFunc func(w http.ResponseWriter, req *http.Request, args ...string)
```

If the operator function calls `w.Write` or `w.WriteHeader`, the server will stop executing entries and the request/response is completed.

## Fixed Filter File

```go
server := httpfilter.NewServer("./root", "C:/_filter.txt")
```

Putting a path into the second argument of the server constructor will force the server to use that filter file for every request. You can use a fixed filter file and the `@` selector to route subdomains to other servers.
