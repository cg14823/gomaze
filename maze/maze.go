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
	Col int
	Row int
}

func (c *CellIndex) Equal(other *CellIndex) bool {
	return c.Row == other.Row && c.Col == other.Col
}

func (c *CellIndex) GetID(cols int) int {
	return (c.Row * cols) + c.Col
}

type Cell struct {
	Top    bool
	Bottom bool
	Left   bool
	Right  bool

	In      bool
	Visited bool
	End     bool
	Start   bool
}

func (c *Cell) Blocked() bool {
	return !c.In
}

type Maze struct {
	Rows  int
	Cols  int
	Start CellIndex
	End   CellIndex
	Cells [][]Cell
}

func (m *Maze) AsciiDraw() {
	for i := 0; i < m.Cols; i++ {
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
	// 0 <= Row < height
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
		m.join(fCell.Row, fCell.Col)
		delete(frontierSet, fCell)
		addToSet(frontierSet, m.getFrontierCell(fCell.Row, fCell.Col))
	}

	// choose an arbitrary Start and an End
	startRow := rand.Intn(rows)
	var startCol int
	if startRow == 0 || startRow == rows-1 {
		startCol = rand.Intn(cols)
	} else if rand.Intn(100) < 50 {
		startCol = cols - 1
	}

	start := CellIndex{
		Col: startCol,
		Row: startRow,
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
		Col: endCol,
		Row: endRow,
	}

	m.Start = start
	m.End = end
	m.Cells[start.Row][start.Col].Start = true
	m.Cells[end.Row][end.Col].End = true
}

func (m *Maze) VisitCell(row, col int) {
	m.Cells[row][col].Visited = true
}

func (m *Maze) UnVisitAll() {
	for r := range m.Cells {
		for c := range m.Cells[r] {
			m.Cells[r][c].Visited = false
		}
	}
}

func (m *Maze) join(row, col int) {
	indexes := []CellIndex{
		{
			Row: row + 1,
			Col: col,
		},
		{
			Row: row - 1,
			Col: col,
		},
		{
			Col: col - 1,
			Row: row,
		},
		{
			Col: col + 1,
			Row: row,
		},
	}

	possible := make([]int, 0)
	for i, neighbour := range indexes {
		// outside maze
		if neighbour.Col < 0 || neighbour.Col > m.Cols-1 || neighbour.Row < 0 || neighbour.Row > m.Rows-1 {
			continue
		}

		if !m.Cells[neighbour.Row][neighbour.Col].Blocked() {
			possible = append(possible, i)
		}
	}

	if len(possible) == 0 {
		return
	}

	chosenOne := possible[rand.Intn(len(possible))]
	switch chosenOne {
	case 0:
		m.Cells[indexes[0].Row][indexes[0].Col].Top = true
		m.Cells[row][col].Bottom = true
	case 1:
		m.Cells[indexes[1].Row][indexes[1].Col].Bottom = true
		m.Cells[row][col].Top = true
	case 2:
		m.Cells[indexes[2].Row][indexes[2].Col].Right = true
		m.Cells[row][col].Left = true
	case 3:
		m.Cells[indexes[3].Row][indexes[3].Col].Left = true
		m.Cells[row][col].Right = true
	}

	m.Cells[row][col].In = true
}

func (m *Maze) getFrontierCell(row, col int) []CellIndex {
	frontiers := make([]CellIndex, 0)
	// Top, Bottom, Left, Right
	indexes := []CellIndex{
		{
			Row: row + 1,
			Col: col,
		},
		{
			Row: row - 1,
			Col: col,
		},
		{
			Col: col - 1,
			Row: row,
		},
		{
			Col: col + 1,
			Row: row,
		},
	}

	for _, neighbour := range indexes {
		// outside maze
		if neighbour.Col < 0 || neighbour.Col > m.Cols-1 || neighbour.Row < 0 || neighbour.Row > m.Rows-1 {
			continue
		}

		if m.Cells[neighbour.Row][neighbour.Col].Blocked() {
			frontiers = append(frontiers, neighbour)
		}
	}

	return frontiers
}

func (m *Maze) Image(outImage string) error {
	cellWidth, cellHeight, wallWidth, xOffset, yOffset, margin := getMeasurements(m.Cols, m.Rows)
	xDimension := m.Cols*(cellWidth+wallWidth) + margin*2
	yDimensions := m.Rows*(cellHeight+wallWidth) + margin*2

	img := generateEmptyImage(xDimension, yDimensions)

	m.drawMap(img, cellWidth, cellHeight, wallWidth, xOffset, yOffset, margin)

	return saveImage(outImage, img)
}

func (m *Maze) drawMap(img *image.RGBA, cellWidth, cellHeight, wallWidth, xOffset, yOffset, margin int) {
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

				paintExits(img, xOffset, yOffset, x, y, m.Cols, m.Rows, cellWidth, cellHeight, cellColor)
			} else if c.End {
				cellColor = color.Color(color.RGBA{
					R: 0,
					G: 255,
					B: 0,
					A: 255,
				})

				paintExits(img, xOffset, yOffset, x, y, m.Cols, m.Rows, cellWidth, cellHeight, cellColor)
			}

			paintCell(img, xOffset, yOffset, cellWidth, cellHeight, cellColor)
			removeWall(img, xOffset, yOffset, cellWidth, cellHeight, wallWidth, &c)
			xOffset += cellWidth + wallWidth
		}

		yOffset += cellHeight + wallWidth
		xOffset = margin
	}
}

func saveImage(outImage string, img *image.RGBA) error {
	f, err := os.Create(outImage)
	if err != nil {
		return fmt.Errorf("could not create image: %w", err)
	}

	err = png.Encode(f, img)
	if err != nil {
		return fmt.Errorf("could not create image: %w", err)
	}

	return f.Close()
}

func generateEmptyImage(xSize, ySize int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, xSize, ySize))

	for r := 0; r < ySize; r++ {
		for c := 0; c < xSize; c++ {
			img.Set(c, r, color.Black)
		}
	}

	return img
}

// getMeasurements returns (scaler, cellWidth, cellHeight, wallWidth, xOffset, yOffset, margin)
func getMeasurements(cols, rows int) (int, int, int, int, int, int) {
	var scaler int
	if rows > cols {
		scaler = rows / 50
	} else {
		scaler = cols / 50
	}

	cellWidth, cellHeight := 2+scaler, 2+scaler
	wallWidth := 1 + scaler
	xOffset, yOffset := 5+scaler, 5+scaler
	margin := 5 + scaler
	return cellWidth, cellHeight, wallWidth, xOffset, yOffset, margin
}

func (m *Maze) ImageWithPath(path []*CellIndex, outImage string, cellColor color.Color) error {
	cellWidth, cellHeight, wallWidth, xOffset, yOffset, margin := getMeasurements(m.Cols, m.Rows)
	xDimension := m.Cols*(cellWidth+wallWidth) + margin*2
	yDimensions := m.Rows*(cellHeight+wallWidth) + margin*2

	img := generateEmptyImage(xDimension, yDimensions)

	m.drawMap(img, cellWidth, cellHeight, wallWidth, xOffset, yOffset, margin)
	for _, c := range path {
		paintCell(img, margin+c.Col*(cellWidth+wallWidth), margin+c.Row*(cellHeight+wallWidth), cellWidth,
			cellHeight, cellColor)
	}

	return saveImage(outImage, img)
}

func (m *Maze) ImageWithMultiplePaths(paths [][]*CellIndex, outImage string, cellColors []color.Color) error {
	cellWidth, cellHeight, wallWidth, xOffset, yOffset, margin := getMeasurements(m.Cols, m.Rows)
	xDimension := m.Cols*(cellWidth+wallWidth) + margin*2
	yDimensions := m.Rows*(cellHeight+wallWidth) + margin*2

	img := generateEmptyImage(xDimension, yDimensions)

	m.drawMap(img, cellWidth, cellHeight, wallWidth, xOffset, yOffset, margin)
	for i, path := range paths {
		for _, c := range path {
			paintCell(img, margin+c.Col*(cellWidth+wallWidth), margin+c.Row*(cellHeight+wallWidth), cellWidth,
				cellHeight, cellColors[i])
		}
	}

	return saveImage(outImage, img)
}

func paintCell(img *image.RGBA, x, y, width, height int, c color.Color) {
	for xi := 0; xi < width; xi++ {
		for yi := 0; yi < height; yi++ {
			img.Set(xi+x, yi+y, c)
		}
	}
}

func paintExits(img *image.RGBA, x, y, cx, cy, cols, rows, cellWidth, cellHeight int, c color.Color) {
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
		Rows: rows,
		Cols: cols,
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
