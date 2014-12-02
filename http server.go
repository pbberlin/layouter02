package main

import (
	"github.com/pbberlin/tools/util"
	"io"
	"net/http"
)

func randomizeArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	randomizeArticlesInternal()
	s := util.IndentedDump(articles)
	io.WriteString(w, *s)

}
func tokenizeArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	sizeAll := 0
	for i := 0; i < len(articles); i++ {
		lpA := &articles[i]
		xx := *lpA
		tokens, size := xx.Tokenize()
		sizeAll += size

		io.WriteString(w, spf("\n\narticle size is %v\n", size))
		for i := 0; i < len(tokens); i++ {
			x := tokens[i]
			s := spf("%5v %2v %2v %v\n", x.Size, x.SemanticStart, x.SemanticEnd, x.Text)
			io.WriteString(w, s)
		}
	}
	io.WriteString(w, spf("\n\ntotal size  %v\n", sizeAll))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	s := util.IndentedDump(articles)
	io.WriteString(w, *s)
}

func init() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/randomize-articles", randomizeArticles)
	http.HandleFunc("/tokenize-articles", tokenizeArticles)
	pf("listening on 4000")
	http.ListenAndServe("localhost:4001", nil)
}
