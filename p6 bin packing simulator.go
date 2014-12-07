package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

const nArmorphs = 20

type amorph1 struct {
	amorph
	consumed int
}

var rndAmorphs []amorph1 = make([]amorph1, nArmorphs)

func generateRandomamorph1s() {

	for i := 0; i < nArmorphs; i++ {
		lp := &rndAmorphs[i]
		lp.Nrows = 1 + rand.Intn(3)
		lp.Ncols = 1 + rand.Intn(3)
		if rand.Intn(4) > 2 {
			lp.Nrows = 2 + rand.Intn(6)
			lp.Ncols = 2 + rand.Intn(6)
		}
	}

}

func generateCSS1() string {
	s := ""
	for i := 0; i < 10; i++ {
		s += spf("	.h%v {height:%vpx;}  .w%v { width:%vpx;} \n", i, 100*i, i, 80*i)
	}
	return "<style>" + s + "</style>"
}

func init() {
	generateRandomamorph1s()
	// s := util.IndentedDump(rndAmorphs)
	// pf("%v", *s)
}

func tryBinpack(w http.ResponseWriter, r *http.Request) {
	data := map[string]map[string]interface{}{}
	hdr := map[string]interface{}{"Msg": "this is header"}
	ftr := map[string]interface{}{"Msg": "this is footer"}
	cnt := map[string]interface{}{"Msg": "content Msg"}
	data["vpHeader"] = hdr
	data["vpFooter"] = ftr
	data["vpContent"] = cnt

	dataContent, err := ioutil.ReadFile("tpl-go/content01.html")
	if err != nil {
		fmt.Fprintf(w, "%v <br>\n", err)
		return
	}

	tplContent := `{{define "tplContent"}}` + string(dataContent) + `{{end}}`
	data["vpContent"]["Msg"] = "some Msg"
	data["vpContent"]["Amorphs"] = rndAmorphs
	data["vpContent"]["CSS1"] = generateCSS1()

	renderTemplate(w, r, []string{"empty-ng-page.html", tplContent}, data)
}
