package main

import (
	"flag"
	"fmt"
)

func main() {
	var cells int
	flag.IntVar(&cells, "cells", 25, "The numbers of cell across and wide for the maze")
	flag.Parse()

	maze := maze.NewMaze(cells, cells)
	fmt.Println("Done creating maze; producing image")
	//maze.AsciiDraw()
	maze.Image()
	fmt.Println("Image done")
}
