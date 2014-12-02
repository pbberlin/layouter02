package main

import (
	"fmt"
)

var spf func(format string, a ...interface{}) string = fmt.Sprintf
var pf func(format string, a ...interface{}) (int, error) = fmt.Printf

func main() {
	pf("main() finished\n")
}
