package main

import "bytes"

/*
In this stage we inject another subgroup of WordGroupTokens, transforming
	ArticlesAllTokenized
	 	[article_1]...
	 	[article_2]...
			[WordGroupToken1]
	     	[WordGroupToken2]
	     	[WordGroupToken3]
			[WordGroupToken4]
	     	[WordGroupToken5]
	     	[WordGroupToken6]
	 	[article_n]
into
		...
	 	[article_2]...
	 		[block_1]
				[WordGroupToken1]
		     	[WordGroupToken2]
		     	[WordGroupToken3]
	 		[block_2]
				[WordGroupToken4]
		     	[WordGroupToken5]
	 		[block_3]
		     	[WordGroupToken6]
				...

*/

type BlockifiedContent struct {
	Size      int
	Overshoot int
	Tokens    []*WordGroupToken
}

type ArticleBlockified struct {
	*Article1 // embedding
	Size      int
	Blocks    []*BlockifiedContent
}

var articlesBlockified []ArticleBlockified

const CHARS_BLOCK = 200

func (a *ArticleBlockified) blockify(idxTokens int) *bytes.Buffer {

	b := new(bytes.Buffer)
	cntrBlock := 0   // number of chars
	cntrArticle := 0 // number of chars for all blocks
	tokens := ArticlesAllTokenized[idxTokens]

	for i := 0; i < len(tokens); i++ {

		if cntrBlock == 0 || cntrBlock > CHARS_BLOCK {
			cntrBlock = 0
			a.Blocks = append(a.Blocks, new(BlockifiedContent)) // this appends to the very slice of articleBlockified
		}

		topBlock := a.Blocks[len(a.Blocks)-1]
		topBlock.Tokens = append(topBlock.Tokens, &tokens[i])

		cntrBlock += tokens[i].Size
		cntrArticle += tokens[i].Size
		topBlock.Size = cntrBlock

		b.WriteString(spf("%v %v %v\n", idxTokens, i, cntrBlock))
	}
	a.Size = cntrArticle

	return b
}

func blockifyAll() {
	for i := 0; i < len(ArticlesAllTokenized); i++ {
		articlesBlockified = append(articlesBlockified, make([]ArticleBlockified, 1)...)
		articlesBlockified[i].Article1 = &ArticlesRaw[i]
		articlesBlockified[i].blockify(i)
	}
}
