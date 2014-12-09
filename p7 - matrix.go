package main

import (
	"io"
	"net/http"

	"github.com/pbberlin/tools/util"
)

type Slot struct {
	A         *Amorph
	IsLeftTop bool
}

type Matrix [][]Slot

type Layout struct {
	M                        Matrix
	West, South, East, North int
}

var L1 Layout

func (l *Layout) Init() {
	newM := make([][]Slot, gridRows) // top-left; rows-cols; x=> rows, y=> cols, *not* Ordinate/Abszisse convention
	for i := 0; i < len(newM); i++ {
		newM[i] = make([]Slot, gridCols)
	}
	l.M = newM

}

// Plaster some randomly generated amorphs onto the layout
func (l *Layout) Seed() {
	seeds := 3
	lastColPos := 5

	for i := 0; i < seeds; i++ {
		L1.Plaster(&AmorphsRandom[i], 5, lastColPos)
		lastColPos += AmorphsRandom[i].Ncols + 1
	}
}

// Paste an amorph onto a matrix
func (l *Layout) Plaster(a *Amorph, top, left int) {
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

// Redraw borders
func (l *Layout) Delimit() {

	w := MaxInt
	e := 0
	n := MaxInt
	s := 0
	for row := 0; row < len(l.M); row++ {
		for col := 0; col < len(l.M[row]); col++ {
			if l.M[row][col].A != nil {
				if col < w {
					w = col
				}
				if col > e {
					e = col
				}

				if row < n {
					n = row
				}
				if row > s {
					s = row
				}
			}
		}
	}
	l.West = w
	l.East = e + 1
	l.North = n
	l.South = s + 1

}

func matrixRaw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "showing current (layout) matrix \n\n")
	s2 := util.IndentedDump(L1)
	io.WriteString(w, *s2)
}

func init() {
	L1.Init()
	L1.Seed()
	// L1.Plaster(&AmorphsRandom[0], 7, 8)
	// L1.Plaster(&AmorphsRandom[1], 7, 12)
	L1.Delimit()
}
