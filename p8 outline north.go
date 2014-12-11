package main

func (l *Layout) OutlineDrawN() {

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
			line = l.completeAndAppend(line, prev, row, col)
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

		// pf("dir%v  rowcol: %v %v\n", direction, row, col)

		cntr++
		if col > l.East || row > l.CRow || cntr > 40 {
			line = l.completeAndAppend(line, prev, row, col)
			break
		}

		prev = direction

	}

}
