# httpfilter

 - [Installation](#Installation)
 - [Tutorial](#Tutorial)
 - [How it Works](#How-it-Works)
 - [Usage Example](#Usage-Example)
 - [Standard Operators](#Standard-Operators)
 - [Writing Filters](#Writing-Filters)
 - [Attaching Additional Operators](#Attaching-Additional-Operators)
 - [Writing Operator Functions](#Writing-Operator-Functions)

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

When the httpfilter server recieves a request for `/index.html`, it will first look for the filter in the same directory as the requested file. The server will read the filter from the top—down, searching for an entry where the **selector** matches the requested file name (*also called the query*). If an entry's selector matches the query, the server will call the Go function that is mapped to the **operator** of that entry. Any additional **arguments** are passed to the Go function as a slice of strings.

```
        entry
┏━━━━━━━━━┻━━━━━━━━━━━━━━━━╍┅
#redirect index.html home.html
 ┃          ┃       ┗━━━┳━━╍┅
 operator   ┃        arguments
            ┃
         selector
```

The httpfilter server passes the query to the operator function, which will either write a response or return a modified query. The server will use the modified query to compare against the next selectors in the filter and the server will pass the modified query to the next operator function. The server will continue to read and execute entries from the filter until a response has been written.

If the server reaches the end of the filter and a response still hasn't been written, the server will call the default operator, which will attempt to serve the file named *query* or respond with a `404 Not found` error.

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
	server := httpfilter.NewServer("./root")
	log.Fatal(http.ListenAndServe(":8080", server))
}
```

## Standard Operators

#### `deft`

*The default operator.*  
The `deft` operator attempts to serve the file named *query* or responds with a `404 Not found` error.
A `#deft *` will be automatically appended to the end of every filter.

#### `ignore`

The `ignore` operator responds with a `404 Not Found` error. Use this if you want to prevent access to a specific file, e.g.
```
#ignore secret.txt
```

#### `pseudo`

The `pseudo` operator replaces the query with the first argument. Use this if you want to make a file appear as if it has a different name.
```
#pseudo home home.html
```
The client may request either `/home` or `/home.html` and they will both serve `home.html`.

#### `redirect`

The `redirect` operator redirects a request to the URL or path in the first argument. E.g. this will redirect `/index.html` to `/home.html`:
```
#redirect index.html home.html
```

## Writing Filters

Operators are always prefixed by a `#`.
Operators, selectors, and arguments are separated by spaces.
Entries are separated by line breaks.

The default operators can be very powerful if you use them in combination.
A few reminders:
 - The filter is read from the top—down and the server will never read upwards
 - The query may be changed by operator functions

Here is an example showing how the query can change:
```
#pseudo secret.txt a
#ignore secret.txt
#pseudo a secret.txt
```
In this example, clients *will* be able to access `secret.txt`. When the client requests `/secret.txt`, the first `pseudo` call renames the query to `a`. The `ignore` has no effect on the query `a`. The second `pseudo` renames the query back to `secret.txt`. Finally, the default operator will recieve the query `secret.txt` and serve the file.

### Selectors with `*`

You can use a `*` in the selector to select all queries.
 - `*` will match with all queries.
 - `a.*` will match with all queries where the name is `a` regardless of the extension.
 - `*.a` will match with all queries where the extension is `.a`.

For example, this filter will `ignore` all queries where the extension is `.txt`:
```
#ignore *.txt
```

### Bulk Operator Syntax

If you have multiple entries that use the same operator repeatedly, e.g.
```
#pseudo home home.html
#pseudo about about.html
#pseudo contact contact.html
```
you can write the operator once and put the entries below it:
```
#pseudo
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
type OpFunc func(w http.ResponseWriter, req httpfilter.FilterRequest) string
```

The `FilterRequest` object looks like this:
```go
type FilterRequest struct {
	*http.Request
	Query      string
	Args       []string
	...
}
```

You may attach a generic session interface to the request using `FilterRequest.SetSession(interface{})`.  
This session can be accessed in any future operator functions that handle the request by calling `FilterRequest.GetSession()`.

When the operator function returns, the server will replace the query with the returned string.
If the operator function calls `w.WriteHeader`, the server will stop executing entries and the request/response is completed.
