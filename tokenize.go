package main

import (
	// "github.com/pbberlin/tools/util"
	"strings"
)

type WordGroupToken struct {
	Size          int
	SemanticStart string // possible values: p  br  h1  h2  ...
	SemanticEnd   string // possible values: p  br  h1  h2  ...
	Text          string
}

func (b *article) NormalizeForTokenize() {
	// todo: make all tags lowercase, remove some attributes, normalize <BR >  to <br/>
}

// Tokenize() splits articles into smallest semantically contingent groups of words.
// The resulting tokens need to be "stuffed" into layout blocks,
// and can not be ripped apart any further.
// Preconditions for Tokenize():
// Any attribuged inline markup is allowed i.e. <b class='blinking'>group of words</b>.
// Block markup is only allowed as follows: <br/> <p></p>, <h1></h1>, <h2></h2>.
// All lower case, all without attributes
func (b *article) Tokenize() ([]WordGroupToken, int) {

	if b.Ps == nil || len(b.Ps) == 0 {
		panic("Fill Article before tokenizing it")
	}

	cntrTokens := 0
	sumSize := 0
	var all []WordGroupToken

	for i := 0; i < len(b.Ps); i++ {

		// Even longer headlines should never by broken up
		if len(b.H2s[i]) > 1 {
			all = append(all, make([]WordGroupToken, 1)...)
			all[cntrTokens].SemanticStart = "h2"
			all[cntrTokens].SemanticEnd = "h2"
			all[cntrTokens].Text = b.H2s[i]
			all[cntrTokens].Size = len(b.H2s[i])
			sumSize += len(b.H2s[i])
			cntrTokens++
		}

		// todo: This inefficient
		// we rather want to split in *one* go
		// and by multiple separators
		// and want to retain the separators too
		var sentences []string
		sentences = strings.SplitAfter(*b.Ps[i], ".")
		sentences = RecombineShortTokens(sentences, 15)
		sentences = SplitFurther(sentences, ",")
		sentences = RecombineShortTokens(sentences, 15)
		sentences = SplitFurther(sentences, ";")
		sentences = SplitFurther(sentences, "!")
		sentences = SplitFurther(sentences, "?")
		sentences = SplitFurther(sentences, "<br/>")
		sentences = RecombineShortTokens(sentences, 15)

		// dump := util.IndentedDump(sentences)
		// pf((*dump))

		numRefinedSentences := len(sentences)
		all = append(all, make([]WordGroupToken, numRefinedSentences)...)
		for j := 0; j < numRefinedSentences; j++ {

			s := strings.TrimSpace(sentences[j])

			//
			// overhung from last block
			if strings.HasPrefix(s, "</p>") {
				if cntrTokens > 0 {
					all[cntrTokens-1].SemanticEnd = "p"
				}
				s = strings.TrimPrefix(s, "</p>")
			}

			//
			// regular
			if strings.HasPrefix(s, "<p>") {
				all[cntrTokens].SemanticStart = "p"
				s = strings.TrimPrefix(s, "<p>")
			}
			if strings.HasSuffix(s, "</p>") {
				all[cntrTokens].SemanticEnd = "p"
				s = strings.TrimSuffix(s, "</p>")
			}
			if strings.HasSuffix(s, "<br/>") {
				all[cntrTokens].SemanticStart = "br"
				all[cntrTokens].SemanticEnd = "br"
				s = strings.TrimSuffix(s, "<br/>")
			}

			all[cntrTokens].Text = s
			all[cntrTokens].Size = len(s)
			sumSize += len(s)

			cntrTokens++
		}
	}

	// dump := util.IndentedDump(all)
	// pf((*dump))

	return all, sumSize
}

// RecombineShortTokens removes empty tokens
// It also recombines all elements shorter than "atLeast"
func RecombineShortTokens(sentencesRaw []string, atLeast int) []string {

	sentencsRefined := make([]string, len(sentencesRaw))
	idx2 := 0
	hangover := ""
	for idxRaw, v := range sentencesRaw {

		if len(sentencesRaw[idxRaw]) < 1 {
			continue // remove empty tokens
		}
		if len(sentencesRaw[idxRaw]) < atLeast {
			hangover = hangover + v
			continue
		}

		sentencsRefined[idx2] = hangover + v
		hangover = ""
		idx2++
	}

	// a last hangover might be un-recombined:
	if len(hangover) > 0 {
		sentencsRefined[idx2-1] += hangover
		//pf("%s\n%s\n%s\n", hangover, sentencsRefined[idx2-1], sentencsRefined[idx2])
	}

	// constrain newly refined slice to non empty elements
	sentencsRefined = sentencsRefined[:idx2]

	return sentencsRefined

}

// SplitFurther takes an slice of strings
// and splits them by sep
// and returns a flattened array of string
func SplitFurther(sIn []string, sep string) []string {
	var sOut []string
	cntr := 0
	for _, v := range sIn {
		sLp := strings.SplitAfter(v, sep)
		sOut = append(sOut, make([]string, len(sLp))...)
		for j, k := range sLp {
			sOut[cntr+j] = k
		}
		cntr += len(sLp)
	}
	return sOut
}
