package main

/*


*/

import (
	"io"
	"math/rand"
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

	seeds := 3

	if AmorphsRandom[1].Nrows < 2 {
		AmorphsRandom[1].Nrows = 2 + rand.Intn(4)
	}
	if AmorphsRandom[1].Ncols < 2 {
		AmorphsRandom[1].Ncols = 2 + rand.Intn(4)
	}

	// for i := 0; i < seeds; i++ {
	// 	AmorphsRandom[i].Nrows++
	// 	AmorphsRandom[i].Ncols++
	// }

	// if AmorphsRandom[1].Nrows < AmorphsRandom[0].Nrows {
	// 	AmorphsRandom[1].Nrows = 1 + rand.Intn(5)
	// }

	// if AmorphsRandom[1].Ncols < 2 {
	// 	AmorphsRandom[1].Ncols = 2 + rand.Intn(5)
	// }

	// if AmorphsRandom[2].Nrows < AmorphsRandom[3].Nrows {
	// 	AmorphsRandom[2].Nrows = AmorphsRandom[3].Nrows + rand.Intn(4)
	// }

	lastRowPos := l.CRow
	lastColPos := l.CCol - AmorphsRandom[0].Ncols - AmorphsRandom[1].Ncols/2

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

const (
	_ = iota // ignore first value by assigning to blank identifier
	UP
	DOWN
	LEFT
	RIGHT
)

type Line struct {
	Row1, Col1, Row2, Col2 int
	DrawRow, DrawCol       int  // left-top for HTML rendering
	Vert                   bool // vertical or horizonal - for HTML rendering
	Length                 int
}

func (l *Layout) OutlineDraw() {

	// init
	var (
		direction, prev int
		line            = Line{}
	)

	// init position
	row := l.CRow
	col := l.West
	l.Outline = make([]Line, 0)

	// init line
	line.Row1 = row
	line.Col1 = col
	pf("outls rowcol: %v %v\n", row, col)

	// init direction
	if l.M.Filled(row-1, col) {
		direction, prev = UP, UP
		row--
	} else {
		direction, prev = RIGHT, RIGHT
		col++
	}
	pf("dir%v  rowcol: %v %v\n", direction, row, col)

	cntr := 0
	for {

		if direction == UP {
			// pf("%b --", !l.M.Filled(row-1, col))
			if !l.M.Filled(row-1, col) { // checking northeast
				direction = RIGHT
			}
		} else if direction == DOWN {
			if l.M.Filled(row, col) { // checking southheast
				direction = RIGHT
			}
		} else if direction == RIGHT {
			if l.M.Filled(row-1, col) { // checking northeast
				direction = UP
			} else if !l.M.Filled(row, col) { // checking southheast
				direction = DOWN
			}
		}

		if direction != prev {
			line = l.completeAndAppend(line, row, col)
		}

		if direction == UP {
			row--
		}
		if direction == DOWN {
			row++
		}
		if direction == RIGHT {
			col++
		}

		prev = direction

		pf("dir%v  rowcol: %v %v\n", direction, row, col)

		cntr++
		if col > l.East || row > l.CRow || cntr > 40 {
			line = l.completeAndAppend(line, row, col)
			break
		}
	}

}

func (l *Layout) completeAndAppend(line Line, row, col int) (newLine Line) {

	// complete line and append it
	line.Row2 = row
	line.Col2 = col

	line.Length = util.Abs(line.Col2 - line.Col1)
	if line.Col2 == line.Col1 {
		line.Vert = true
		line.Length = util.Abs(line.Row2 - line.Row1)
	}

	line.DrawRow = line.Row1
	if line.Row2 < line.Row1 {
		line.DrawRow = line.Row2
	}

	line.DrawCol = line.Col1
	if line.Col2 < line.Col1 {
		line.DrawCol = line.Col2
	}

	l.Outline = append(l.Outline, line)

	newLine.Row1 = row
	newLine.Col1 = col

	pf("  line closed %v\n", line)
	return
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
	L1.OutlineDraw()
}
