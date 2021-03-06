package main

import "github.com/pbberlin/tools/util"

const (
	_ = iota // ignore first value by assigning to blank identifier
	UP
	DOWN
	LEFT
	RIGHT
)

type jsonDirection int

type Line struct {
	Row1, Col1, Row2, Col2 int
	DrawRow, DrawCol       int  // left-top for HTML rendering
	Vert                   bool // vertical or horizonal - for HTML rendering
	Length                 int
	Direction              jsonDirection // UP, DOWN, ...
}

// overriding normal json marshal => writing the string instead of the int
func (d jsonDirection) MarshalJSON() ([]byte, error) {

	s := ""
	switch {
	case d == 1:
		s = "up"
	case d == 2:
		s = "down"
	case d == 3:
		s = "left"
	case d == 4:
		s = "right"
	}

	return []byte(`"` + s + `"`), nil
}
func (d jsonDirection) String() string {
	b, _ := d.MarshalJSON()
	return string(b)
}
func (l *Layout) completeAndAppend(NorthOrSouth bool, line Line, direction, row, col int) (newLine Line) {

	// complete old line
	line.Row2 = row
	line.Col2 = col

	// start init new line
	newLine.Row1 = row
	newLine.Col1 = col

	// mere dots originate in init loop of drawOutlineX() - skip those:
	// if line.Col1 == line.Col2 && line.Row1 == line.Row2 {
	// 	return
	// }

	// init loop => skip
	if direction == 0 {
		return
	}

	line.Direction = jsonDirection(direction)

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

	if NorthOrSouth {
		l.OutlineN = append(l.OutlineN, line)
	} else {
		l.OutlineS = append(l.OutlineS, line)
	}

	// pf("  line closed %v\n", line)
	return
}
