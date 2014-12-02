package main

import (
	"github.com/drhodes/golorem"
	"math/rand"
	"time"
)

type layout struct {
}

type article struct {
	H1  string
	H2s []string
	Ps  []*string // this is an exercise in modification of pointer and value slices
}

var articles []article = make([]article, 3+rand.Intn(4))

var data01 layout

func randomizeArticlesInternal() {

	articles = make([]article, 8+rand.Intn(4))

	pf("--\n")

	for i := 0; i < len(articles); i++ {

		articleSmall := false
		if t := rand.Intn(8); t > 1 {
			articleSmall = true // 75%  smaller sections
			pf("  small\n")
		} else {
			pf("large\n")
		}

		lpA := &articles[i]
		*lpA = article{}
		lpA.H1 = lorem.Sentence(4, 8)

		h2Sections := 3 + rand.Intn(3)
		if articleSmall {
			h2Sections = 1 + rand.Intn(2)
		}

		lpA.H2s = make([]string, h2Sections)
		lpA.Ps = make([]*string, h2Sections)
		for j := 0; j < h2Sections; j++ {
			var lpH2 *string
			lpH2 = &lpA.H2s[j]
			*lpH2 = lorem.Sentence(6, 14)

			var stringInit string = ""
			lpA.Ps[j] = &stringInit

			numParagraphs := 3 + rand.Intn(10)
			if articleSmall {
				numParagraphs = 1 + rand.Intn(3)
			}
			for k := 0; k < numParagraphs; k++ {
				var newPar string
				newPar = lorem.Paragraph(2, 6)
				*lpA.Ps[j] += "<p>" + newPar + "</p>"
			}

		}

	}

}

func init() {
	rand.Seed(time.Now().UnixNano())
	randomizeArticlesInternal()
	pf("data01.init() finished\n")
}
