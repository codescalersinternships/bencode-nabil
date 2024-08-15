# bencode-nabil

Encoder and Decode of Bencode

## Installation

As a library

```shell
go get github.com/codescalersinternships/bencode-nabil/pkg
```

## Usage

in your Go app you can do something like

```go
package main

import (
	"fmt"

	bencoder "github.com/codescalersinternships/bencode-nabil/pkg"
)

func main() {
	s := "d3:bar3:moo4:spaml4:spam4:foooee"

	ret,_ := bencoder.Decoder(s)
	fmt.Println(ret)
	x,_ := bencoder.Encoder(ret)
	fmt.Println(string(x[:]))
}

```

## Testing

```shell
make test
```