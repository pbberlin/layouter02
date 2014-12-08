package main

import (
	"fmt"
	tt "html/template"
	"net/http"
	"strings"

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

	var err error

	// base template
	urlBaseTemplate := r.FormValue("tplName") // nails the *first* instance
	if urlBaseTemplate != "" {
		tf[0] = urlBaseTemplate // override with URL param
	}
	tBase := tt.New("tplBase").Funcs(funcMap)                  // func map forces previous naming, hooks us to some 'tplBase'
	tBase, err = tBase.ParseFiles("tpl-go/" + tf[0] + ".html") // must begin with {{define "tplBase"}}
	if err != nil {
		fmt.Fprintf(w, "%v <br>\n", err)
		return
	}

	//
	// now additional templates -
	nUrlArgs := len(r.Form["tplName"])
	for i := 1; i < len(tf); i++ {
		// override with URL param
		if i < nUrlArgs {
			urlContentTemplate := r.Form["tplName"][i] // nails the *first* instance
			if urlContentTemplate != "" {
				tf[i] = urlContentTemplate
			}
		}
		if strings.HasPrefix(tf[i], "{{define") {
			tBase, err = tBase.Parse(tf[i]) // definitions must appear at top level - but not at the start
		} else {
			// assuming filename, trying to load it
			tBase, err = tBase.ParseFiles("tpl-go/" + tf[i] + ".html")
		}
		if err != nil {
			fmt.Fprintf(w, "%v (idx %v) <br>\n", err, i)
			return
		}
	}

	// finally, combine template with data
	err = tBase.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "%v (tpl execution)<br>\n", err)
		return
	}

}

// example
func testPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]map[string]interface{}{}
	hdr := map[string]interface{}{"Msg": "header Msg"}
	ftr := map[string]interface{}{"Msg": "footer Msg"}
	cnt := map[string]interface{}{"Msg": "content Msg"}
	data["vpHeader"] = hdr
	data["vpFooter"] = ftr
	data["vpContent"] = cnt

	nOptions := 4
	r.ParseForm()
	data["vpContent"]["tplName"] = spf("%v", r.Form["tplName"])
	wasSet := util.StringSliceToMapKeys(r.Form["tplName"])
	sel := map[string]string{}
	for i := 0; i < nOptions; i++ {
		opt := spf("content%02d", i)
		sel[opt] = "nosel"
		if wasSet[opt] {
			sel[opt] = " selected "
		}
	}
	data["vpContent"]["tplNameSel"] = sel

	tplContent := `{{define "tplContent"}}
		here is content template msg {{.Msg}}
		<form>
			{{.tplName}}<br>
			{{.tplNameSel}}<br>
			<select multiple="multiple" name="tplName" size="2" style='font-size:10px;'>
				<option value="base-01-ng" selected="" >base-01-ng.html.html </option><br>
				<option value="base-02"                >&nbsp; base-02.html &nbsp; &nbsp; </option>  
			</select><br>
			<select multiple="multiple" name="tplName" size="2" style='font-size:10px;'>
				<option value="content01" selected="" >content01.html </option><br>
				<option value="content02"  >&nbsp; content02.html &nbsp; &nbsp; </option>  
			</select><br>
			<input type="input" name="tplName" value="another content04" />
			<input type='sUbmit' accesskey='u' /> 
		</form>
	{{end}}
	`
	renderTemplate(w, r, []string{"base-01-ng", tplContent}, data)
}

func init() {
	http.HandleFunc("/test-page", testPage)
}
