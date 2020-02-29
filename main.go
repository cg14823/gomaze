package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/cg14823/gomaze/pathfinding"

	"github.com/cg14823/gomaze/maze"
)

func main() {
	var cells int
	var pathFind, fileOut string
	flag.IntVar(&cells, "cells", 25, "The numbers of cell across and wide for the maze")
	flag.StringVar(&pathFind, "path-find", "", "The path finding algorithm to use available are [bfs]")
	flag.StringVar(&fileOut, "file-out", "", "Image file with the maze")
	flag.Parse()

	m := maze.NewMaze(cells, cells)
	fmt.Println("Done creating maze; producing image")
	//maze.AsciiDraw()

	if fileOut == "" {
		fileOut = fmt.Sprintf("./out/maze-%dx%d-%d", cells, cells, time.Now().Unix())
	}

	var err error
	switch pathFind {
	case "bfs":
		out, steps, err := pathfinding.BFS(m)
		if err != nil {
			fmt.Printf("Could not find path afet %d steps\n", steps)
			os.Exit(1)
		}

		path := pathfinding.SearchCellToSlice(out)
		fmt.Printf("Path found took %d steps path lenght %d\n", steps, len(path))

		err = m.ImageWithPath(path, fileOut, color.RGBA{
			R: 100,
			G: 0,
			B: 100,
			A: 255,
		})
	default:
		err = m.Image(fileOut)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Image done")
}
