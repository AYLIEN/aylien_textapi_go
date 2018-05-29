About
=====

This is the Go client library for AYLIEN's APIs. If you haven't already done so, you will need to [sign up](https://developer.aylien.com/signup).

Installation
============

To install, simply use `go get`:

```bash
$ go get github.com/AYLIEN/aylien_textapi_go
```

See the [Developers Guide](https://developer.aylien.com/docs) for additional documentation.

Example
=======

```go
package main

import (
	"fmt"
	textapi "github.com/AYLIEN/aylien_textapi_go"
)

func main() {
	auth := textapi.Auth{"YourApplicationId", "YourApplicationKey"}
	client, err := textapi.NewClient(auth, true)
	if err != nil {
		panic(err)
	}
	params := &textapi.SentimentParams{Text: "John is a very good football player!"}
	sentiment, err := client.Sentiment(params)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", sentiment)
}
```

Support
=======

support@aylien.com
