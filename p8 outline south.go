package main

func (l *Layout) OutlineSDraw() {

	l.OutlineS = make([]Line, 0)

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
		case northEast && !southEast:
			chCol = 1
			dir = RIGHT
		case southEast:
			chRow = 1
			dir = DOWN
		case !northEast && !southEast:
			chRow = -1
			dir = UP
		}

		if prev != dir {
			// pf("  line finalized %v %v\n", prev, dir)
			line.Row2 = row
			line.Col2 = col
			line = l.completeAndAppend(false, line, prev, row, col)
		}

		row += chRow
		col += chCol
		prev = dir
		pf("dir%8s  rowcol: %v %v\n", jsonDirection(dir), row, col)

		// final condition
		emptyNorthWest := col == l.East && !l.M.Filled(row-1, col-1)
		if emptyNorthWest || row > l.South || row <= l.North || col > l.East {
			line = l.completeAndAppend(false, line, prev, row, col)
			break
		}
	}

}
