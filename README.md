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
	s := "d3:bar4:spam3:fooi42ee"

	ret,_ := bencoder.Encoder(s)
	fmt.Println(ret)
}

```

## Testing

```shell
make test
```