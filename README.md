# Unofficial Go library for fetching info from BMW ConnectedDrive

## Installation

`go get github.com/sdrobov/connected-drive`

## Usage
```go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	connecteddrive "github.com/sdrobov/connected-drive"
	"net/http"
)

func main() {
	var authStorage []byte
	c := connecteddrive.NewClient("user@example.com", "userPassword", bytes.NewBuffer(authStorage), http.DefaultClient)
	v, e := c.GetVehicles(context.Background())
	if e != nil {
		panic(e)
	}

	j, _ := json.Marshal(v)
	fmt.Printf("%v\n", string(j))
}
```

## Legal
This library are in no way connected to the company BMW AG. BMW and ConnectedDrive are registered trademarks of BMW AG.
