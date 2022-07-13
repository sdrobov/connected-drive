package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	connecteddrive "github.com/sdrobov/connected-drive"
	"net/http"
)

func main() {
	var authStorage []byte
	c := connecteddrive.NewClient("user@example.com", "userPassword", bytes.NewBuffer(authStorage), http.DefaultClient)
	v, e := c.GetVehicles()
	if e != nil {
		panic(e)
	}

	j, _ := json.Marshal(v)
	fmt.Printf("%v\n", string(j))
}
