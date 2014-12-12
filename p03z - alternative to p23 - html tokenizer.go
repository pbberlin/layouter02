package main

import (
	"bytes"
	"io"

	"code.google.com/p/go.net/html"
)

// src http://golang-examples.tumblr.com/page/2
func parseHtml(r io.Reader) *bytes.Buffer {

	b := new(bytes.Buffer)
	d := html.NewTokenizer(r)
	for {
		// token type
		tokenType := d.Next()
		if tokenType == html.ErrorToken {
			panic(spf("error in tokenizing html: %v", tokenType))
		}
		token := d.Token()
		switch tokenType {
		case html.StartTagToken: // <tag>
			// type Token struct {
			//     Type     TokenType
			//     DataAtom atom.Atom
			//     Data     string
			//     Attr     []Attribute
			// }
			//
			// type Attribute struct {
			//     Namespace, Key, Val string
			// }
		case html.TextToken: // text between start and end tag
			b.WriteString(spf("%v", token))
		case html.EndTagToken: // </tag>
		case html.SelfClosingTagToken: // <tag/>

		}
	}

	return b
}
