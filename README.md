
A lightweight proxy server with URL-based replacement.

`go get github.com/pranavraja/replaceproxy`

# Usage

Assuming `$GOPATH/bin` is in your `$PATH`

    replaceproxy [-repl url_prefix -with file_or_url]

Start a proxy server on `0.0.0.0:8080`. Log each response to stdout, color-coded by status.

e.g. `replaceproxy -repl http://google.com -with http://myotherserver.com/google.txt`

Replace all requests under the google.com domain with the response from myotherserver.com/google.txt. Preserve headers and status code of the original response.

e.g. `replaceproxy -repl http://myapi.com/user/1/feeds -with feeds.json`

Replace requests to myapi.com/user/1/feeds with the canned response stored in feeds.json. Preserve headers and status code of the original response.
