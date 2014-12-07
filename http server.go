package main

import "net/http"

func serveSingleRootFile(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func init() {
	// static resources - Mandatory root-based
	serveSingleRootFile("/sitemap.xml", "./sitemap.xml")
	serveSingleRootFile("/favicon.ico", "./favicon.ico")
	serveSingleRootFile("/robots.txt", "./robots.txt")
	// static resources - other
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))
	http.Handle("/tpl-ng/", http.StripPrefix("/tpl-ng/", http.FileServer(http.Dir("./tpl-ng/"))))

	// => do this in main http.ListenAndServe("localhost:4000", nil)
}
