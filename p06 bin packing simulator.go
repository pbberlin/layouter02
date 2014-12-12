package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

const (
	numExampleAmorphs = 20 // randomly created, for exercise and study
	gridWidth         = 40
	gridHeight        = 50

	gridRows = 14
	gridCols = 20

	TotalGridWidth  = gridWidth * gridCols
	TotalGridHeight = gridHeight * gridRows
)

var AmorphsRandom []Amorph = make([]Amorph, numExampleAmorphs)

func generateRandomAmorphs() {

	for i := 0; i < numExampleAmorphs; i++ {
		lp := &AmorphsRandom[i]
		lp.Nrows = 1 + rand.Intn(3)
		lp.Ncols = 1 + rand.Intn(3)
		if rand.Intn(4) > 2 {
			lp.Nrows = 2 + rand.Intn(6)
			lp.Ncols = 2 + rand.Intn(6)
		}
		lp.IdxArticle = i
	}

}

func generateCSS1() string {
	s := ""
	for i := 0; i < 30; i++ {
		s += spf("	.h%v {height:%vpx;}  .w%v { width:%vpx;} \n", i, gridHeight*i,
			i, gridWidth*i)
		s += spf("	.t%v {top:%vpx;}  .l%v { left:%vpx;} \n", i, gridHeight*i,
			i, gridWidth*i)

	}
	return "<style>" + s + "</style>"
}

func init() {
	generateRandomAmorphs()
}

func tryBinpack01(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	hdr := map[string]interface{}{"Msg": "this is header"}
	ftr := map[string]interface{}{"Msg": "this is footer"}
	cnt := map[string]interface{}{"Msg": "content Msg"}
	data["vpHeader"] = hdr
	data["vpFooter"] = ftr
	data["vpContent"] = cnt
	data["HeadCSSLink"] = `<link rel="stylesheet" href="/css/bin-pack-grid.css" media="screen" type="text/css" />`
	data["HeadTitle"] = `Bin packing study`

	dataContent, err := ioutil.ReadFile("tpl-go/content01.html")
	if err != nil {
		fmt.Fprintf(w, "%v <br>\n", err)
		return
	}

	tplContent := string(dataContent)
	cnt["Msg"] = "some Msg"
	cnt["Amorphs"] = AmorphsRandom
	cnt["Layout1"] = L1
	cnt["TotalGridWidth"] = TotalGridWidth
	cnt["TotalGridHeight"] = TotalGridHeight
	cnt["CSS1"] = generateCSS1()

	renderTemplate(w, r, []string{"base-01-ng", tplContent, "repository-of-amorphs"}, data)
}
