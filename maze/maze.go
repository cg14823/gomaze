package maze

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

type CellIndex struct {
	col int
	row int
}

type Cell struct {
	Top    bool
	Bottom bool
	Left   bool
	Right  bool

	In    bool
	End   bool
	Start bool
}

func (c *Cell) Blocked() bool {
	return !c.In
}

type Maze struct {
	rows  int
	cols  int
	start CellIndex
	end   CellIndex
	Cells [][]Cell
}

func (m *Maze) AsciiDraw() {
	for i := 0; i < m.cols; i++ {
		fmt.Printf("___")
	}

	fmt.Println()
	for _, r := range m.Cells {
		for _, c := range r {
			if c.Left {
				fmt.Printf("_")
			} else {
				fmt.Printf("|")
			}

			if c.Start {
				fmt.Printf("S")
			} else if c.End {
				fmt.Printf("E")
			} else if c.Bottom {
				fmt.Printf(" ")
			} else {
				fmt.Printf("_")
			}

			if c.Right {
				fmt.Printf("_")
			} else {
				fmt.Printf("|")
			}
		}
		fmt.Println()
	}
}

func (m *Maze) Create(rows, cols int) {
	rand.Seed(time.Now().Unix())
	// initialise grid
	m.Cells = make([][]Cell, rows)
	for r := 0; r < rows; r++ {
		m.Cells[r] = make([]Cell, cols)
	}

	// pick a random cell on the edge
	// 0 <= row < height
	row := rand.Intn(rows)
	col := rand.Intn(cols)

	m.Cells[row][col].In = true

	frontierSet := make(map[CellIndex]struct{})
	frontiers := m.getFrontierCell(row, col)
	addToSet(frontierSet, frontiers)

	var step int
	for len(frontierSet) != 0 {
		step++

		fCell := getRandomFrontierCell(frontierSet)
		m.join(fCell.row, fCell.col)
		delete(frontierSet, fCell)
		addToSet(frontierSet, m.getFrontierCell(fCell.row, fCell.col))
	}

	// choose an arbitrary start and an end
	startRow := rand.Intn(rows)
	var startCol int
	if startRow == 0 || startRow == rows-1 {
		startCol = rand.Intn(cols)
	} else if rand.Intn(100) < 50 {
		startCol = cols - 1
	}

	start := CellIndex{
		col: startRow,
		row: startCol,
	}

	endRow := startRow
	endCol := startCol
	for startRow == endRow && endCol == startCol {
		endRow := rand.Intn(rows)
		if endRow == 0 || endRow == rows-1 {
			endCol = rand.Intn(cols)
		} else if rand.Intn(100) < 50 {
			endCol = cols - 1
		} else {
			endCol = 0
		}
	}

	end := CellIndex{
		col: endCol,
		row: endRow,
	}

	m.start = start
	m.end = end
	m.Cells[start.row][start.col].Start = true
	m.Cells[end.row][end.col].End = true

}

func (m *Maze) join(row, col int) {
	indexes := []CellIndex{
		{
			row: row + 1,
			col: col,
		},
		{
			row: row - 1,
			col: col,
		},
		{
			col: col - 1,
			row: row,
		},
		{
			col: col + 1,
			row: row,
		},
	}

	possible := make([]int, 0)
	for i, neighbour := range indexes {
		// outside maze
		if neighbour.col < 0 || neighbour.col > m.cols-1 || neighbour.row < 0 || neighbour.row > m.rows-1 {
			continue
		}

		if !m.Cells[neighbour.row][neighbour.col].Blocked() {
			possible = append(possible, i)
		}
	}

	if len(possible) == 0 {
		return
	}

	chosenOne := possible[rand.Intn(len(possible))]
	switch chosenOne {
	case 0:
		m.Cells[indexes[0].row][indexes[0].col].Top = true
		m.Cells[row][col].Bottom = true
	case 1:
		m.Cells[indexes[1].row][indexes[1].col].Bottom = true
		m.Cells[row][col].Top = true
	case 2:
		m.Cells[indexes[2].row][indexes[2].col].Right = true
		m.Cells[row][col].Left = true
	case 3:
		m.Cells[indexes[3].row][indexes[3].col].Left = true
		m.Cells[row][col].Right = true
	}

	m.Cells[row][col].In = true
}

func (m *Maze) getFrontierCell(row, col int) []CellIndex {
	frontiers := make([]CellIndex, 0)
	// Top, Bottom, Left, Right
	indexes := []CellIndex{
		{
			row: row + 1,
			col: col,
		},
		{
			row: row - 1,
			col: col,
		},
		{
			col: col - 1,
			row: row,
		},
		{
			col: col + 1,
			row: row,
		},
	}

	for _, neighbour := range indexes {
		// outside maze
		if neighbour.col < 0 || neighbour.col > m.cols-1 || neighbour.row < 0 || neighbour.row > m.rows-1 {
			continue
		}

		if m.Cells[neighbour.row][neighbour.col].Blocked() {
			frontiers = append(frontiers, neighbour)
		}
	}

	return frontiers
}

func (m *Maze) Image() {
	var scaler int
	if m.rows > m.cols {
		scaler = m.rows / 50
	} else {
		scaler = m.cols / 50
	}

	cellWidth, cellHeight := 2+scaler, 2+scaler
	wallWidth := 1 + scaler

	img := image.NewRGBA(image.Rect(0, 0, m.cols*(cellWidth+wallWidth)+(5+scaler)*2, m.rows*(cellHeight+wallWidth)+(5+scaler)*2))
	xOffest, yOffset := 5+scaler, 5+scaler
	margin := 5 + scaler

	for r := 0; r < m.rows*(cellHeight+wallWidth)+(5+scaler)*2; r++ {
		for c := 0; c < m.cols*(cellWidth+wallWidth)+(5+scaler)*2; c++ {
			img.Set(c, r, color.Black)
		}
	}

	for y, r := range m.Cells {
		for x, c := range r {
			// draw main block
			cellColor := color.Color(color.White)
			if c.Start {
				cellColor = color.Color(color.RGBA{
					R: 255,
					G: 0,
					B: 0,
					A: 255,
				})

				paintExist(img, xOffest, yOffset, x, y, m.cols, m.rows, cellWidth, cellHeight, cellColor)
			} else if c.End {
				cellColor = color.Color(color.RGBA{
					R: 0,
					G: 255,
					B: 0,
					A: 255,
				})

				paintExist(img, xOffest, yOffset, x, y, m.cols, m.rows, cellWidth, cellHeight, cellColor)
			}

			paintCell(img, xOffest, yOffset, cellWidth, cellHeight, cellColor)
			removeWall(img, xOffest, yOffset, cellWidth, cellHeight, wallWidth, &c)
			xOffest += cellWidth + wallWidth
		}

		yOffset += cellHeight + wallWidth
		xOffest = margin
	}

	name := fmt.Sprintf("./out/maze-%dx%d-%d.png", m.rows, m.cols, time.Now().Unix())
	f, _ := os.Create(name)
	err := png.Encode(f, img)
	if err != nil {
		fmt.Println("Could not do image: ", err)
	}
}

func paintCell(img *image.RGBA, x, y, width, height int, c color.Color) {
	for xi := 0; xi < width; xi++ {
		for yi := 0; yi < height; yi++ {
			img.Set(xi+x, yi+y, c)
		}
	}
}

func paintExist(img *image.RGBA, x, y, cx, cy, cols, rows, cellWidth, cellHeight int, c color.Color) {
	if cy-1 < 0 {
		paintCell(img, x, y-cellHeight, cellWidth, cellHeight, c)
	} else if cy+1 >= rows {
		paintCell(img, x, y+cellHeight, cellWidth, cellHeight, c)
	} else if cx-1 < 0 {
		paintCell(img, x-cellWidth, y, cellWidth, cellHeight, c)
	} else if cx+1 >= cols {
		paintCell(img, x+cellWidth, y, cellWidth, cellHeight, c)
	}
}

func removeWall(img *image.RGBA, x, y, cellWidth, cellHeight, wallWidth int, c *Cell) {
	if c.Top {
		paintCell(img, x, y-wallWidth, cellWidth, wallWidth, color.White)
	}

	if c.Right {
		paintCell(img, x+cellWidth, y, wallWidth, cellHeight, color.White)
	}

	if c.Left {
		paintCell(img, x-wallWidth, y, wallWidth, cellHeight, color.White)
	}

	if c.Bottom {
		paintCell(img, x, y+cellHeight, wallWidth, cellHeight, color.White)
	}
}

func NewMaze(rows, cols int) *Maze {
	maze := &Maze{
		rows: rows,
		cols: cols,
	}

	maze.Create(rows, cols)
	return maze
}

func addToSet(set map[CellIndex]struct{}, sl []CellIndex) {
	for _, x := range sl {
		set[x] = struct{}{}
	}
}

func getRandomFrontierCell(set map[CellIndex]struct{}) CellIndex {
	n := len(set)
	chosenOne := rand.Intn(n)
	var count int
	for k := range set {
		if count == chosenOne {
			return k
		}

		count++
	}

	return CellIndex{}
}
