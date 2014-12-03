package main

import (
	"fmt"
	"net/http"
)

var spf func(format string, a ...interface{}) string = fmt.Sprintf
var pf func(format string, a ...interface{}) (int, error) = fmt.Printf

func main() {
	pf("start listening on 4001\n")
	http.ListenAndServe("localhost:4001", nil) // dont put into http server init() - it blocks the other init() funcs
}
