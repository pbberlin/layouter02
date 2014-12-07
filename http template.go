package main

import (
	"fmt"
	tt "html/template"
	"net/http"

	"github.com/pbberlin/tools/colors"
	"github.com/pbberlin/tools/util"
)

var funcMap = tt.FuncMap{
	"fColorizer": colors.AlternatingColorShades,
	"fMakeRange": func(num int) []int {
		sl := make([]int, num)
		for i, _ := range sl {
			sl[i] = i
		}
		return sl
	},
	"fMult": func(x, y int) int {
		return x * y
	},
	"fAdd": func(x, y int) int {
		return x + y
	},
	"fHTML": func(s string) tt.HTML {
		// to CSS  - http://stackoverflow.com/questions/14765395/why-am-i-seeing-zgotmplz-in-my-go-html-template-output
		return tt.HTML(s)
	},
	"fCSS": func(s string) tt.CSS {
		// to CSS  - http://stackoverflow.com/questions/14765395/why-am-i-seeing-zgotmplz-in-my-go-html-template-output
		return tt.CSS(s)
	},
	"fAttr": func(s string) tt.HTMLAttr {
		return tt.HTMLAttr(s)
	},
	"fGlobId": func() int {
		return <-util.Counter
	},
}

func renderTemplate(w http.ResponseWriter, r *http.Request, tf []string, data interface{}) {

	pTemplateName := r.FormValue("t")
	if pTemplateName != "" {
		tf[0] = pTemplateName
	}

	var err error
	tBase := tt.New("tplBase").Funcs(funcMap)
	tBase, err = tBase.ParseFiles("tpl-go/" + tf[0])
	if err != nil {
		fmt.Fprintf(w, "%v <br>\n", err)
		return
	}

	for i := 1; i < len(tf); i++ {
		tBase, err = tBase.Parse(tf[i]) // definitions must appear at top level - but not at the start
		if err != nil {
			fmt.Fprintf(w, "%v <br>\n", err)
			return
		}
	}

	{
		err := tBase.Execute(w, data)
		if err != nil {
			fmt.Fprintf(w, "%v <br>\n", err)
			return
		}
	}

}

// example
func testPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]map[string]string{}
	hdr := map[string]string{"Msg": "header Msg"}
	ftr := map[string]string{"Msg": "footer Msg"}
	cnt := map[string]string{"Msg": "content Msg"}
	data["vpHeader"] = hdr
	data["vpFooter"] = ftr
	data["vpContent"] = cnt

	tplContent := `{{define "tplContent"}}
		here is content template msg {{.Msg}}
	{{end}}
	`
	renderTemplate(w, r, []string{"empty-ng-page.html", tplContent}, data)
}

func init() {
	http.HandleFunc("/test-page", testPage)
}
