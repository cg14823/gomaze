package pathfinding

import (
	"fmt"

	"github.com/cg14823/gomaze/maze"
)

type SearchCell struct {
	Parent *SearchCell

	cell  *maze.Cell
	index *maze.CellIndex
}

func BFS(m *maze.Maze) (*SearchCell, uint64, error) {
	m.VisitCell(m.Start.Row, m.Start.Col)
	start := SearchCell{
		index: &m.Start,
		cell:  &m.Cells[m.Start.Row][m.Start.Col],
	}

	var steps uint64
	queue := make([]SearchCell, 0)
	queue = append(queue, start)
	for len(queue) > 0 {
		steps++
		current := queue[0]
		queue = queue[1:]

		if current.index.Equal(&m.End) {
			return &current, steps, nil
		}

		// Get all adjacent edges
		connected := getConnectedUnvisitedCells(current.index, m)
		for _, c := range connected {
			c.Parent = &current
			m.VisitCell(c.index.Row, c.index.Col)
			queue = append(queue, c)
		}
	}

	return nil, steps, fmt.Errorf("no path could be found")
}

func getConnectedUnvisitedCells(current *maze.CellIndex, m *maze.Maze) []SearchCell {
	searchCells := make([]SearchCell, 0)
	cell := m.Cells[current.Row][current.Col]

	if cell.Top && !m.Cells[current.Row-1][current.Col].Visited {
		searchCells = append(searchCells, SearchCell{
			cell: &m.Cells[current.Row-1][current.Col],
			index: &maze.CellIndex{
				Col: current.Col,
				Row: current.Row - 1,
			},
		})
	}

	if cell.Bottom && !m.Cells[current.Row+1][current.Col].Visited {
		searchCells = append(searchCells, SearchCell{
			cell: &m.Cells[current.Row+1][current.Col],
			index: &maze.CellIndex{
				Col: current.Col,
				Row: current.Row + 1,
			},
		})
	}

	if cell.Left && !m.Cells[current.Row][current.Col-1].Visited {
		searchCells = append(searchCells, SearchCell{
			cell: &m.Cells[current.Row][current.Col-1],
			index: &maze.CellIndex{
				Col: current.Col - 1,
				Row: current.Row,
			},
		})
	}

	if cell.Right && !m.Cells[current.Row][current.Col+1].Visited {
		searchCells = append(searchCells, SearchCell{
			cell: &m.Cells[current.Row][current.Col+1],
			index: &maze.CellIndex{
				Col: current.Col + 1,
				Row: current.Row,
			},
		})
	}

	return searchCells
}

func SearchCellToSlice(s *SearchCell) []*maze.CellIndex {
	path := make([]*maze.CellIndex, 0)
	path = append(path, s.index)
	current := s.Parent
	for current != nil {
		path = append(path, current.index)
		current = current.Parent
	}

	return path
}
