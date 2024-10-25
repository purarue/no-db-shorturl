# no-db-shorturl

It should not be this hard to have a URL shortened. I don't want to configure a SQL database, run a docker container, install thousands of NPM packages, or configure PHP to redirect URLs.

There are lots of URL shorteners out there, but they mostly use a database as a backend. This has:

- No Configuration Files
- No Database
- No Web Front-End
- No Dependencies other than the go stdlib.
- Statically built binary downloadable from [here](https://sean.fish/p/no-db-shorturl-builds/)

This stores each link in its own individual file, in the `./data` directory.

- To add a URL, you create a file with the hash (name) in the `./data` directory, with the URL to redirect to on the first line. This can also be done with a POST request to the `/` endpoint (see below)
- To change the URL hash for a shorturl, rename the file
- To delete a shorturl, delete the file
- To change what URL a shorturl redirects to, edit the file and change the contents

That is all I want.

This allows for programmatic shorturl generation on my server, using anything that can read/write to files, as well as simple aliases to grab one from my machine by `ssh`ing to my server, like:

```bash
# list all shorturls
alias shorturls="ssh myserver 'ls shorturls'"  # directory name
# copy a shorturl to my clipboard
alias shz="shorturls | fzf | sed -e 's|^|https://sean.fish/s/|' | tee /dev/tty | clipcopy"
```

You can also build this from source instead:

`go install -v "github.com/purarue/no-db-shorturl@latest"`

Usage:

```
Usage of no-db-shorturl:
  -data-folder string
    	directory to store data in (default "./data")
  -port int
    	port to serve shorturl on (default 8040)
  -secret-key string
    	secret key to authenticate POST requests
```

By default, anyone can create a URL, so I'd recommend providing a secret key with the `-secret-key` flag.

You can also set the `SHORTURL_KEY` environment variable, instead of `-secret-key` flag.

To add a shortened URL, make a POST request to the `/` endpoint

To use `curl` as an example:

`curl --header "Content-Type: application/json" --request POST --data '{"key":"your_secret_key","url":"https://sean.fish"}' http://localhost:8040`

or to specify the path to create the shortcut on:

`curl --header "Content-Type: application/json" --request POST --data '{"key":"your_secret_key","url":"https://sean.fish","hash":"short"}' http://localhost:8040`

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

I have a script [here](https://github.com/purarue/vps/blob/master/bin/shorten) which I use to create new URLs.
