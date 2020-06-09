# no-db-static-shorturl

It should not be this hard to have a URL shortened. I don't to configure a SQL database, or run a docker container, or install thousands of NPM packages or configure PHP to redirect URLs.

There are lots of URL shorteners out there, but they mostly use a database as a backend. Here's what this does:

- No Configuration Files
- No Database
- No Web Front-End
- Statically built binary downloadable from [here]()

This stores each link in its own individual file, in the `./data` directory.

- To delete a shorturl, delete the corresponding file.
- To rename shorturls, rename the file.
- To change what URL a shorturl redirects to, edit the file and change the contents.

That is all I want.

Due to all these choices, the performance wont be as great as the database backed URL shorteners. Though this is just for my personal website, this isn't going to be used by thousands of people every day, and golang is no slacker.

```
Usage of no-db-static-shorturl:
  -data-folder string
    	directory to store data in (default "./data")
  -port int
    	port to serve shorturl on (default 8040)
  -secret-key string
    	secret key to authenticate POST requests
```

By default, anyone can create a URL, so I'd recommend providing a secret key with the `-secret-key` flag.

To add a shortened URL, make a POST request to this endpoint. Example:

```
curl --header "Content-Type: application/json" --request POST --data '{"key":"your_secret_key","url":"https://sean.fish"}' http://localhost:8040
```

or to specify the path to create the shortcut on:

```
curl --header "Content-Type: application/json" --request POST --data '{"key":"your_secret_key","url":"https://sean.fish","hash":"short"}' http://localhost:8040
```

I use this with `nginx`, like so:

```
server {
  listen 443 ssl;
  ssl_certificate ....

  location /s/ {
    proxy_pass https://127.0.0.1:8040/;
  }
}
```

Which makes this accessible at `https://mywebsite/s/`.
