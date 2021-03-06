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

type Layout struct {
	M                            Matrix
	CRow, CCol                   int // center row, center column
	West, South, East, North     int // outer border of inscribed amorphs
	IWest, ISouth, IEast, INorth int // inner border of inscribed amorphs
	OutlineN                     []Line
	OutlineS                     []Line
	ConcavitiesN                 []Line
	ConcavitiesS                 []Line
}

var L1 Layout

func (pm *Matrix) Filled(row, col int) bool {
	m := *pm
	// pf("---%v %v\n", row, col)
	if row < 0 || col < 0 || row > len(m)-1 || col > len(m[row])-1 { // beyound slice boundary => for OutlineN drawing this means not filled
		return false
	}
	if m[row][col].A == nil {
		return false
	}
	return true
}

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
	L1.OutlineNDraw()
	L1.OutlineSDraw()

	// cut off last, line element - it's doubly
	L1.OutlineN = L1.OutlineN[:len(L1.OutlineN)-1]

	// now combine the first elements from both outlines into one
	n1 := L1.OutlineN[0]
	s1 := L1.OutlineS[0]
	if n1.Direction == jsonDirection(UP) &&
		s1.Direction == jsonDirection(DOWN) &&
		s1.Col1 == n1.Col1 {
		combiLine := Line{}
		combiLine.Row1 = s1.Row2
		combiLine.Col1 = s1.Col2
		combiLine.Row2 = n1.Row2
		combiLine.Col2 = n1.Col2
		combiLine.Direction = jsonDirection(UP)
		combiLine.DrawRow = n1.Row2
		combiLine.DrawCol = n1.Col2
		combiLine.Vert = true
		combiLine.Length = combiLine.Row1 - combiLine.Row2
		L1.OutlineN[0] = combiLine    // replace
		L1.OutlineS = L1.OutlineS[1:] // remove
	}
	L1.CheckConcavityN()
	L1.CheckConcavityS()
}

func init() {
	layoutPipeline()
}
