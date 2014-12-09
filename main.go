package main

import (
	"fmt"
	"net/http"
)

var spf func(format string, a ...interface{}) string = fmt.Sprintf
var pf func(format string, a ...interface{}) (int, error) = fmt.Printf

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

func main() {
	pf("start listening on 4001\n")
	http.ListenAndServe("localhost:4001", nil) // dont put into http server init() - it blocks the other init() funcs
}
