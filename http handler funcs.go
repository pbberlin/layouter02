package main

import (
	"io"
	"net/http"

	"github.com/pbberlin/tools/util"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "showing current articles data without newly randomizing\n\n")
	s := util.IndentedDump(ArticlesRaw)
	io.WriteString(w, *s)
}

func randomizeArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "newly randomizing... \n\n")
	randomizeArticlesInternal()
	s := util.IndentedDump(ArticlesRaw)
	io.WriteString(w, *s)

}

func tokenizeArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "tokenizing articles ... \n\n")
	b := articlesToRawString(ArticlesRaw)
	w.Write(b.Bytes())
}

func tokenizedShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "result of tokenizing articles ... \n\n")
	s := util.IndentedDump(ArticlesAllTokenized)
	io.WriteString(w, *s)
}

func regenerateRandomAmorphs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	generateRandomAmorphs()

	L1.Init()
	L1.Seed()
	L1.Delimit()
	L1.OutlineDraw()

	io.WriteString(w, "result of randomized amorps ... \n\n")
	s := util.IndentedDump(AmorphsRandom)
	io.WriteString(w, *s)
}

//--------------------------------------
func pipelineAll(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "doing it all ... \n\n")
	randomizeArticles(w, r)
	tokenizeArticles(w, r)

	blockifyAll()
	io.WriteString(w, "---- ... \n\n")
	s1 := util.IndentedDump(articlesBlockified)
	io.WriteString(w, *s1)

	amorphifyAll()
	io.WriteString(w, "---- ... \n\n")
	s2 := util.IndentedDump(articlesAmorphified)
	io.WriteString(w, *s2)

}

func backend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, "--  <a href='/'                   target='b_out'>Raw Articles</a><br>\n")
	io.WriteString(w, "<a href='/randomize-articles' target='b_out'>Randomize Articles</a><br>\n")
	io.WriteString(w, "<a href='/tokenize-articles'  target='b_out'>Tokenize  Articles</a><br>\n")
	io.WriteString(w, "--  <a href='/tokenized-show'   target='b_out'>Tokenized Show</a><br>\n")
	io.WriteString(w, "--  <a href='/regenerate-random-amorphs'   target='b_out'>Randomized Amorphs regenerate</a><br>\n")

	io.WriteString(w, "--  <a href='/try-bin-pack'   target='b_out'>Binpack study 1</a><br>\n")
	io.WriteString(w, `--  <a href='/try-bin-pack?tplName=base-01-ng&tplName=content02'   
			target='b_out'>Binpack study 2</a><br>`+"\n")
	io.WriteString(w, `--  <a href='/try-bin-pack?tplName=base-01-ng&tplName=content03'   
			target='b_out'>Binpack - Random seeded</a><br>`+"\n")

	io.WriteString(w, "--  <a href='/matrix-raw'     target='b_out'>Matrix raw</a><br>\n")

	io.WriteString(w, "------------------------------------------<br>\n")
	io.WriteString(w, "--  <a href='/pipeline-all'   target='b_out'  accesskey='p'><b>P</b>ipeline All</a><br>\n")
}

func init() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/randomize-articles", randomizeArticles)
	http.HandleFunc("/tokenize-articles", tokenizeArticles)
	http.HandleFunc("/tokenized-show", tokenizedShow)
	http.HandleFunc("/regenerate-random-amorphs", regenerateRandomAmorphs)

	http.HandleFunc("/pipeline-all", pipelineAll)
	http.HandleFunc("/backend", backend)

	http.HandleFunc("/try-bin-pack", tryBinpack01)
	http.HandleFunc("/matrix-raw", matrixRaw)

	pf("http server init complete\n")
}
