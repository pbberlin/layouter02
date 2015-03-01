package main

import (
	"github.com/pbberlin/tools/subsort"
	"github.com/pbberlin/tools/util"
)

// The number of blocks is cast into several bounding rectangles.
// I.e. five blocks might be cast like this:
//
// xx xx xx xx xx	OR	xx xx xx	OR  xx xx
// xx xx xx xx xx      	xx xx xx		xx xx
//
//                      xx xx			xx xx
//                      xx xx			xx xx
//
//                         				xx
//                         				xx
//
// We put them into slices of amorphs
// a := []amorph{amorph{4,1},amorph{3,2},amorph{2,3}}
// The term is chosen to evoke biological amoeba,
// who can assume various forms
//

const (
	goldenRatio = 1 / 1.41 // most aesthetically pleasing, i.e. 3 rows / 4 cols
	fattest     = 0.25     // each row with four times columns
	slimmest    = 4        // four rows - one col - eight rows - two cols
)

type Amorph struct {
	Nrows          int
	Ncols          int
	Slack          int
	AestheticValue string // 100 - abs((ratio -  goldenRatio)*100) - string formatted %04d for generic sorting
	IdxArticle     int    // reference 'upward', since this type is also used in the 'matrix' later on
}

type ArticleAmorphified struct {
	*ArticleBlockified // embedding
	Amorphs            []*Amorph
	AmorphsByAesth     []subsort.SortedByStringVal // leightweight, therefore no pointer
	AmorphsBySlack     []subsort.SortedByIntVal
	ArticleConsumed    bool // any of my morphs admitted into the layout?
}

var articlesAmorphified []ArticleAmorphified

func (a *ArticleAmorphified) amorphify() {

	nBlocks := len(a.Blocks)
	previousCols := 0
	previousSlack := 0

	for row := 1; row <= nBlocks; row++ {

		cols := nBlocks / row
		if nBlocks%row != 0 {
			cols = nBlocks/row + 1
		}

		slack := row*cols - nBlocks
		if cols == previousCols &&
			slack >= previousSlack {
			continue
		}

		ratio := float64(row) / float64(cols)
		if ratio > slimmest || ratio < fattest {
			continue
		}

		a.Amorphs = append(a.Amorphs, new(Amorph))
		topAmorph := a.Amorphs[len(a.Amorphs)-1]
		topAmorph.Ncols = cols
		topAmorph.Slack = slack
		topAmorph.Nrows = row
		topAmorph.AestheticValue = spf("%04d", 1000-util.Abs(int((ratio-goldenRatio)*100)))

		previousCols = cols
		previousSlack = slack
	}

}

func amorphifyAll() {
	for i := 0; i < len(articlesBlockified); i++ {
		articlesAmorphified = append(articlesAmorphified, make([]ArticleAmorphified, 1)...)
		articlesAmorphified[i].ArticleBlockified = &articlesBlockified[i]
		articlesAmorphified[i].amorphify()

		f1 := func(k int) string {
			return articlesAmorphified[i].Amorphs[k].AestheticValue
		}
		f2 := func(k int) int {
			return articlesAmorphified[i].Amorphs[k].Slack
		}
		bA2 := subsort.SortByStringValDesc(len(articlesAmorphified[i].Amorphs), f1)
		bS2 := subsort.SortByIntValAsc(len(articlesAmorphified[i].Amorphs), f2)
		articlesAmorphified[i].AmorphsByAesth = bA2
		articlesAmorphified[i].AmorphsBySlack = bS2

	}
}
