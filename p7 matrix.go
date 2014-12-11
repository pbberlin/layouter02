package main

/*


*/

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

func (pm *Matrix) Filled(row, col int) bool {
	m := *pm
	// pf("---%v %v\n", row, col)
	if row < 0 || col < 0 || row > len(m)-1 || col > len(m[row])-1 { // beyound slice boundary => for outline drawing this means not filled
		return false
	}
	if m[row][col].A == nil {
		return false
	}
	return true
}

type Layout struct {
	M                            Matrix
	CRow, CCol                   int // center row, center column
	West, South, East, North     int // outer border of inscribed amorphs
	IWest, ISouth, IEast, INorth int // inner border of inscribed amorphs
	Outline                      []Line
	OutlineConcavity             bool
}

func (l *Layout) CheckConcavity() {
	curTrend := jsonDirection(0)
	alternations := 0
	for i := 0; i < len(l.Outline); i++ {
		if l.Outline[i].Direction == RIGHT || l.Outline[i].Direction == LEFT {
			// nothing
		}
		if l.Outline[i].Direction == UP || l.Outline[i].Direction == DOWN {
			lpDir := jsonDirection(l.Outline[i].Direction)
			if curTrend == jsonDirection(0) { // first trend
				curTrend = lpDir
				continue
			}
			if curTrend == lpDir {
				// fine
			}
			if curTrend != lpDir {
				curTrend = lpDir
				alternations++
			}
		}
	}
	l.OutlineConcavity = false
	if alternations > 1 {
		l.OutlineConcavity = true
	}

}

var L1 Layout

func (l *Layout) Init() {
	newM := make([][]Slot, gridRows) // top-left; rows-cols; x=> rows, y=> cols, *not* Ordinate/Abszisse convention
	for i := 0; i < len(newM); i++ {
		newM[i] = make([]Slot, gridCols)
	}
	l.M = newM

	l.CRow = len(l.M) / 2
	l.CCol = len(l.M[0]) / 2

}

// Plaster some randomly generated amorphs onto the layout
func (l *Layout) Seed() {

	seeds := 4

	// if AmorphsRandom[1].Nrows < 2 {
	// 	AmorphsRandom[1].Nrows = 2 + rand.Intn(4)
	// }
	// if AmorphsRandom[1].Ncols < 2 {
	// 	AmorphsRandom[1].Ncols = 2 + rand.Intn(4)
	// }

	lastRowPos := l.CRow
	lastColPos := l.CCol - AmorphsRandom[0].Ncols - AmorphsRandom[1].Ncols

	for i := 0; i < seeds; i++ {
		L1.Plaster(&AmorphsRandom[i], lastRowPos-AmorphsRandom[i].Nrows/2, lastColPos)
		// lastColPos += AmorphsRandom[i].Ncols + 1
		lastColPos += AmorphsRandom[i].Ncols
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

	iw := w
	ie := e
	in := n
	is := s
	for row := n; row < s+1; row++ {
		for col := w; col < e+1; col++ {
			if l.M[row][col].A == nil {
				if col >= iw && col < l.CCol {
					iw = col + 1
				}
				if col <= ie && col >= l.CCol {
					ie = col - 1
				}

				if row >= in && row < l.CRow {
					in = row + 1
					// pf("%v %v %v\n", row, col, is)
				}
				if row <= is && row >= l.CRow {
					is = row - 1
					// pf("%v %v %v\n", row, col, is)
				}
			}
		}
	}
	l.IWest = iw
	l.IEast = ie + 1

	l.INorth = in
	l.ISouth = is + 1

}

func matrixRaw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "showing current (layout) matrix \n\n")
	s2 := util.IndentedDump(L1)
	io.WriteString(w, *s2)
}

func layoutPipeline() {
	L1.Init()
	L1.Seed()
	// L1.Plaster(&AmorphsRandom[0], 7, 8)
	// L1.Plaster(&AmorphsRandom[1], 7, 12)
	L1.Delimit()
	L1.OutlineDrawN()
	L1.OutlineDrawS()
	L1.CheckConcavity()

}

func init() {
	layoutPipeline()
}
