package main

import (
	"math/rand"
	"time"

	"github.com/drhodes/golorem"
)

type Article1 struct {
	H1 string

	// This is an exercise in modification of pointers and values of slices
	H2s []string  // variant a - see below
	Ps  []*string // variant b - see below
}

var ArticlesRaw []Article1 = make([]Article1, 2)

func randomizeArticlesInternal() {

	ArticlesRaw = make([]Article1, 8+rand.Intn(4))

	pf("--\n")

	for i := 0; i < len(ArticlesRaw); i++ {

		articleSmall := false
		if t := rand.Intn(8); t > 1 {
			articleSmall = true // 75%  smaller sections
			pf("  small\n")
		} else {
			pf("large\n")
		}

		lpA := &ArticlesRaw[i]
		*lpA = Article1{}
		lpA.H1 = lorem.Sentence(4, 8)

		h2Sections := 3 + rand.Intn(3)
		if articleSmall {
			h2Sections = 1 + rand.Intn(2)
		}

		lpA.H2s = make([]string, h2Sections)
		lpA.Ps = make([]*string, h2Sections)
		for j := 0; j < h2Sections; j++ {

			{
				var lpPtrH2 *string              // variant a - line 1
				lpPtrH2 = &lpA.H2s[j]            // variant a - line 2
				*lpPtrH2 = lorem.Sentence(6, 14) // variant a - line 3

				lpValH2 := *lpPtrH2           // variant a - wrong approach
				lpValH2 = "futile assignment" // utterly futile
				_ = lpValH2                   // prevent "not used error"
			}

			{
				var stringInit string = "" // variant b - line 1
				lpA.Ps[j] = &stringInit    // variant b - line 2

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

}

func init() {
	rand.Seed(time.Now().UnixNano()) // before article generation!!!!
	randomizeArticlesInternal()
	pf("article generation init() finished\n")
}
