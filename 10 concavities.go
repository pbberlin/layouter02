package main

func (l *Layout) CheckConcavityN() {
	trendMarker := jsonDirection(0)

	concavities := make([]Line, 0) // we use a lines to store concavities into  - left-top-bottomright
	tmpC := Line{}

	for i := 0; i < len(l.OutlineN); i++ {
		lpLine := l.OutlineN[i]
		if lpLine.Direction == jsonDirection(DOWN) &&
			trendMarker != lpLine.Direction { // opening angle of cavity
			trendMarker = lpLine.Direction
			tmpC = Line{}
			tmpC.Row1 = lpLine.Row1
			tmpC.Col1 = lpLine.Col1
		} else if lpLine.Direction == jsonDirection(UP) && // first closing angle
			trendMarker == jsonDirection(DOWN) {
			trendMarker = jsonDirection(0) // reset
			tmpC.Row2 = lpLine.Row1
			tmpC.Col2 = lpLine.Col1
			concavities = append(concavities, tmpC)
		}
	}

	for i := 0; i < len(concavities); i++ {
		pf("n concavity #%v  %v %v %v %v \n", i, concavities[i].Row1, concavities[i].Col1, concavities[i].Row2, concavities[i].Col2)
	}
	l.ConcavitiesN = concavities

}

func (l *Layout) CheckConcavityS() {

	trendMarker := jsonDirection(0)

	concavities := make([]Line, 0) // we use a lines to store concavities into  - left-top-bottomright
	tmpC := Line{}

	for i := 0; i < len(l.OutlineS); i++ {
		lpLine := l.OutlineS[i]
		if lpLine.Direction == jsonDirection(UP) &&
			trendMarker != lpLine.Direction { // opening angle of cavity
			trendMarker = lpLine.Direction
			tmpC = Line{}
			tmpC.Row1 = lpLine.Row1
			tmpC.Col1 = lpLine.Col1
		} else if lpLine.Direction == jsonDirection(DOWN) && // first closing angle
			trendMarker == jsonDirection(UP) {
			trendMarker = jsonDirection(0) // reset
			tmpC.Row2 = lpLine.Row1
			tmpC.Col2 = lpLine.Col1
			concavities = append(concavities, tmpC)
		}
	}

	for i := 0; i < len(concavities); i++ {
		pf("s concavity #%v  %v %v %v %v \n", i, concavities[i].Row1, concavities[i].Col1, concavities[i].Row2, concavities[i].Col2)
	}
	l.ConcavitiesS = concavities
}
