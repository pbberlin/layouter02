package main

func (l *Layout) OutlineNDraw() {

	l.OutlineN = make([]Line, 0)

	row := l.CRow
	col := l.West
	dir := 0
	prev := 0
	line := Line{Row1: row, Col1: col}

	for {

		northEast := l.M.Filled(row-1, col)
		southEast := l.M.Filled(row, col)
		chCol := 0
		chRow := 0

		switch {
		case northEast:
			chRow = -1
			dir = UP
		case southEast && !northEast:
			chCol = 1
			dir = RIGHT
		case !northEast && !southEast:
			chRow = 1
			dir = DOWN
		}

		if prev != dir {
			line.Row2 = row
			line.Col2 = col
			line = l.completeAndAppend(true, line, prev, row, col)
		}

		row += chRow
		col += chCol
		prev = dir
		// pf("n: dir%8s  rowcol: %v %v\n", jsonDirection(dir), row, col)

		// final condition
		emptySouthWest := col == l.East && !l.M.Filled(row, col-1)
		if emptySouthWest {
			line = l.completeAndAppend(true, line, prev, row, col)
			break
		}
	}

}
