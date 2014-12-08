package main

import (
	"io"
	"net/http"

	"github.com/pbberlin/tools/util"
)

// var Matrix [][]*Amorph

type Slot struct {
	A         *Amorph
	IsLeftTop bool
}

type Matrix [][]Slot

type Layout struct {
	M          Matrix
	WestLimes  int
	WestWall   int
	SouthLimes int
	SouthWall  int
	NorthLimes int
	NorthWall  int
	EastLimes  int
	EastWall   int
}

var L1 Layout

func matrixRaw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// s1 := util.IndentedDump(AmorphsRandom)
	// io.WriteString(w, *s1)

	io.WriteString(w, "showing current (layout) matrix \n\n")
	s2 := util.IndentedDump(L1.M)
	io.WriteString(w, *s2)
}

func init() {
	L1.Init()
	L1.Seed()
	// L1.Fill(&AmorphsRandom[0], 7, 8)
	// L1.Fill(&AmorphsRandom[1], 7, 12)
	L1.Delimit()

	// L1.M[4][5].IsLeftTop = true
}
func (l *Layout) Init() {
	newM := make([][]Slot, gridRows) // top-left; rows-cols; x=> rows, y=> cols, *not* Ordinate/Abszisse convention
	for i := 0; i < len(newM); i++ {
		newM[i] = make([]Slot, gridCols)
	}
	l.M = newM
}

func (l *Layout) Seed() {

	seeds := 3
	lastColPos := 5

	for i := 0; i < seeds; i++ {
		L1.Fill(&AmorphsRandom[i], 5, lastColPos)
		lastColPos += AmorphsRandom[i].Ncols + 1
	}

}

func (l *Layout) Fill(a *Amorph, top, left int) {

	for i := top; i < top+a.Nrows; i++ {
		for j := left; j < left+a.Ncols; j++ {
			if i < len(l.M) && j < len(l.M[0]) {
				l.M[i][j].A = a
			} else {
				pf("Matrix TOO SMALL %v:%v  into  %v:%v   \n", i, j, len(l.M), len(l.M[0]))
			}
		}
	}
	l.M[top][left].IsLeftTop = true

}

func (l *Layout) Delimit() {
	if l.M == nil || l.M[0] == nil {
		panic("Layout struct must be initialized before limiting it")
	}
	m := l.M
	cy := len(m) / 2
	cx := len(m[0]) / 2
	pf("heart of the layout %v %v\n", cx, cy)
}
