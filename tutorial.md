# Tutorial

### Install the server

Download the package and install the example server.
```shell
$ go get github.com/rbxb/fileserve
$ go install github.com/rbxb/fileserve/cmd/fileserve
```

Run the example server. Use `github.com/rbxb/fileserve/example/root` as the directory.
```shell
$ fileserve -port :8080 -directory ./example/root
```

### The pseudo tag

In a browser, go to `http://localhost:8080/home.html`.  
You should see **Hello world!**

Try going to `localhost:8080/home`.  
You should get a 404 Not Found error.

We want `/home` to serve the `home.html` file.  
Open the `_tags.txt` file inside the `example/root/` directory and add this line:
```
#pseudo home home.html
```

Save the tagfile and refresh `localhost:8080/home`.  
You should see **Hello world!**.

### The ignore tag

Now try going to `localhost:8080/secret.txt`.  
We don't want people to be able to access this file.

Add this line to the `_tags.txt` file:
```
#ignore secret.txt
```

Save and refresh `localhost:8080/secret.txt`.  
You should get a 404 Not Found error.

### The redirect tag

We want `localhost:8080` to redirect to `localhost:8080/home`.  
Create a new tagfile named `_tags.txt` in the `example/` directory.
```
#redirect root home
```
`root` will select a request with an empty path.  
If the redirect path is not an absolute URL, the redirect path will be relative to the current directory. Since `root/` is the base directory use only `home` as the redirect path because all paths will already be inside `root/`.  
Go to `localhost:8080` and you should be redirected to `localhost:8080/home`.

Try setting up a redirect to an external website.
```
#redirect root https://google.com
```
`localhost:8080` should redirect to `google.com`.
