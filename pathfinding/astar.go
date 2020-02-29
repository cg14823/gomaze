package pathfinding

import (
	"fmt"

	"github.com/cg14823/gomaze/maze"
)

type AStartSearchCell struct {
	Parent *AStartSearchCell
	f      uint64
	h      uint64
	g      uint64

	cell  *maze.Cell
	index *maze.CellIndex
}

func Astart(m *maze.Maze) ([]*maze.CellIndex, uint64, error) {
	out, steps, err := astart(m)
	if err != nil {
		return nil, steps, err
	}

	path := constructPath(out)
	return path, steps, nil
}

func astart(m *maze.Maze) (*AStartSearchCell, uint64, error) {
	openSet := make([]*AStartSearchCell, 0)
	inOpenList := make(map[maze.CellIndex]struct{})
	h := heuristic(&m.End, &m.Start)

	start := AStartSearchCell{
		Parent: nil,
		f:      h,
		h:      h,
		cell:   &m.Cells[m.Start.Row][m.Start.Col],
		index:  &m.Start,
	}

	var step uint64
	openSet = append(openSet, &start)
	inOpenList[*start.index] = struct{}{}
	for len(openSet) > 0 {
		step++
		current := openSet[0]
		openSet = openSet[1:]

		fmt.Printf("Step: %d, openSet size: %d current: %+v\n", step, len(openSet), current.index)

		delete(inOpenList, *current.index)
		m.VisitCell(current.index.Row, current.index.Col)
		if current.index.Equal(&m.End) {
			return current, step, nil
		}

		adjacent := getAdjacent(current.index, &m.End, m, current.g)
		for _, c := range adjacent {
			if m.Cells[c.index.Row][c.index.Col].Visited {
				continue
			}

			_, ok := inOpenList[*c.index]
			if ok {
				ix := findInSorted(openSet, c.index)
				if ix == -1 {
					fmt.Println("This should not happen")
					continue
				}

				if openSet[ix].g < c.g {
					continue
				}

				// remove that entry from the openSet and replace it by the new one
				openSet = append(openSet[:ix], openSet[ix+1:]...)
			}

			c.Parent = current
			openSet = sortedInsert(openSet, c)
			inOpenList[*c.index] = struct{}{}
		}
	}

	return nil, step, fmt.Errorf("could not find path")
}

// heuristic is based on the straight line distance from the cell to the goal
func heuristic(goal *maze.CellIndex, current *maze.CellIndex) uint64 {
	return uint64((goal.Row-current.Row)*(goal.Row-current.Row) + (goal.Col-current.Col)*(goal.Col-current.Col))
}

func findInSorted(open []*AStartSearchCell, cellIndex *maze.CellIndex) int {
	for i, c := range open {
		if c.index.Equal(cellIndex) {
			return i
		}
	}

	return -1
}

func sortedInsert(open []*AStartSearchCell, insert *AStartSearchCell) []*AStartSearchCell {
	var insertIx int
	for i, c := range open {
		insertIx = i
		if c.f > insert.f {
			break
		}
	}

	switch insertIx {
	case 0:
		open = append([]*AStartSearchCell{insert}, open...)
	case len(open) - 1:
		open = append(open, insert)
	default:
		open = append(append(open[:insertIx], insert), open[insertIx:]...)
	}

	return open
}

func getAdjacent(current, goal *maze.CellIndex, m *maze.Maze, g uint64) []*AStartSearchCell {
	searchCells := make([]*AStartSearchCell, 0)
	cell := m.Cells[current.Row][current.Col]

	if cell.Top {
		cell := &AStartSearchCell{
			cell: &m.Cells[current.Row-1][current.Col],
			index: &maze.CellIndex{
				Col: current.Col,
				Row: current.Row - 1,
			},
		}

		h := heuristic(goal, cell.index)
		cell.h = h
		cell.g = g + 1
		cell.f = g + h + 1

		searchCells = append(searchCells, cell)
	}

	if cell.Bottom {
		cell := &AStartSearchCell{
			cell: &m.Cells[current.Row+1][current.Col],
			index: &maze.CellIndex{
				Col: current.Col,
				Row: current.Row + 1,
			},
		}

		h := heuristic(goal, cell.index)
		cell.h = h
		cell.g = g + 1
		cell.f = g + h + 1

		searchCells = append(searchCells, cell)
	}

	if cell.Left {
		cell := &AStartSearchCell{
			cell: &m.Cells[current.Row][current.Col-1],
			index: &maze.CellIndex{
				Col: current.Col - 1,
				Row: current.Row,
			},
		}

		h := heuristic(goal, cell.index)
		cell.h = h
		cell.g = g + 1
		cell.f = g + h + 1

		searchCells = append(searchCells, cell)
	}

	if cell.Right {
		cell := &AStartSearchCell{
			cell: &m.Cells[current.Row][current.Col+1],
			index: &maze.CellIndex{
				Col: current.Col + 1,
				Row: current.Row,
			},
		}

		h := heuristic(goal, cell.index)
		cell.h = h
		cell.g = g + 1
		cell.f = g + h + 1

		searchCells = append(searchCells, cell)
	}

	return searchCells
}

func constructPath(s *AStartSearchCell) []*maze.CellIndex {
	path := make([]*maze.CellIndex, 0)
	path = append(path, s.index)
	current := s.Parent
	for current != nil {
		path = append(path, current.index)
		current = current.Parent
	}

	return path
}
