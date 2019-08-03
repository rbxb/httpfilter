# tutorial

[README.md](./README.md)

## Setup

Make sure you have installed the package and the sample program.
```shell
$ go get github.com/rbxb/httpfilter
$ go install github.com/rbxb/httpfilter/cmd/httpfilter
```

Find the example directory at `/httpfilter/example`.

Run the sample program from the example directory.
```shell
$ cd ./httpfilter/example
$ httpfilter
```

Now we have a server that will serve static files from the directory `/httpfiler/example/root`. You should disable page caching in your browser for the next steps.

### The `pseudo` operator

The `pseudo` operator replaces the query with the first argument. Use this if you want to make a file appear as if it has a different name.

Try going to `localhost:8080/home.html`. You should see "**Hello world!**". Try going to `localhost:8080/home`. You should get a `404 Not found` error.

Open the file named `_filters.txt`, which is inside `/httpfilter/example/root`. We want the server to serve the file `home.html` when someone requests `/home`. Add this line to the filter file:
```
#pseudo home home.html
```
Save the filter file and go to `localhost:8080/home`. You should see "**Hello world!**".

### The `ignore` operator

The `ignore` operator responds with a `404 Not Found` error. Use this if you want to prevent access to a specific file.

Try going to `localhost:8080/secret.txt`.  
Add the line `#ignore secret.txt` to the filter file. Your filter file should now look like this:
```
#pseudo home home.html
#ignore secret.txt
```
Save and go to `localhost:8080/secret.txt`. You should get a `404 Not found` error.

### The `redirect` operator

The `redirect` operator redirects a request to the URL or path in the first argument.

Try going to `localhost:8080/index.html`. You should get a `404 Not found` error. We want this page to redirect the client to `/home`.
```
#pseudo home home.html
#ignore secret.txt
#redirect index.html home
```
Save the filter file and go to `localhost:8080/index.html`. You should be redirected to `/home`.
