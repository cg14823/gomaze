package pathfinding

import (
	"fmt"

	"github.com/cg14823/gomaze/maze"
)

func DFS(m *maze.Maze) (*SearchCell, uint64, error) {
	m.VisitCell(m.Start.Row, m.Start.Col)
	start := SearchCell{
		index: &m.Start,
		cell:  &m.Cells[m.Start.Row][m.Start.Col],
	}

	var steps uint64
	stack := make([]SearchCell, 0)
	stack = append(stack, start)
	for len(stack) > 0 {
		steps++
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if current.index.Equal(&m.End) {
			return &current, steps, nil
		}

		connected := getConnectedUnvisitedCells(current.index, m)
		for _, c := range connected {
			c.Parent = &current
			m.VisitCell(c.index.Row, c.index.Col)
			stack = append(stack, c)
		}
	}

	return nil, steps, fmt.Errorf("could not find path")
}
