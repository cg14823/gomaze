package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"github.com/cg14823/gomaze/pathfinding"

	"github.com/cg14823/gomaze/maze"
)

func main() {
	var cells int
	var pathFind, fileOut, compareAlgos string
	flag.IntVar(&cells, "cells", 25, "The numbers of cell across and wide for the maze")
	flag.StringVar(&pathFind, "path-find", "", "The path finding algorithm to use available are [bfs, stack]")
	flag.StringVar(&compareAlgos, "compare-algos", "", "Comma separated list of algos to compare")
	flag.StringVar(&fileOut, "file-out", "", "Image file with the maze")
	flag.Parse()

	if pathFind != "" && compareAlgos != "" {
		fmt.Println("cannot provide both -path-find and -compare-algos")
		os.Exit(1)
	}

	m := maze.NewMaze(cells, cells)
	fmt.Println("Done creating maze; producing image")
	//maze.AsciiDraw()

	if fileOut == "" {
		fileOut = fmt.Sprintf("./out/maze-%dx%d-%d", cells, cells, time.Now().Unix())
	}

	if pathFind != "" {
		var out *pathfinding.SearchCell
		var steps uint64
		var err error

		switch pathFind {
		case "bfs":
			out, steps, err = pathfinding.BFS(m)
		case "dfs":
			out, steps, err = pathfinding.DFS(m)
		default:
			fmt.Println("Unrecognized path finding algo:", pathFind)
			os.Exit(1)
		}

		if err != nil {
			fmt.Printf("Could not find path afet %d steps\n", steps)
			os.Exit(1)
		}

		path := pathfinding.SearchCellToSlice(out)
		fmt.Printf("Path found took %d steps path length %d\n", steps, len(path))

		err = m.ImageWithPath(path, fileOut, color.RGBA{
			R: 100,
			G: 0,
			B: 100,
			A: 255,
		})

		if err != nil {
			fmt.Println("Could not create image:", err.Error())
			os.Exit(1)
		}

		fmt.Println("Image done")
		os.Exit(0)
	}

	if compareAlgos != "" {
		algos := strings.Split(compareAlgos, ",")
		paths := make([][]*maze.CellIndex, 0)
		colours := make([]color.Color, 0)
		for _, a := range algos {
			var out *pathfinding.SearchCell
			var steps uint64
			var err error
			var c color.Color

			switch strings.TrimSpace(a) {
			case "dfs":
				out, steps, err = pathfinding.BFS(m)
				c = color.RGBA{
					B: 250,
					A: 150,
				}
			case "bfs":
				out, steps, err = pathfinding.DFS(m)
				c = color.RGBA{
					G: 255,
					A: 130,
				}
			default:
				fmt.Println("Unknown algo:", a)
				os.Exit(1)
			}

			if err != nil {
				fmt.Printf("Could not find path afet %d steps\n", steps)
				os.Exit(1)
			}

			path := pathfinding.SearchCellToSlice(out)
			fmt.Printf("Path found using %s took %d steps path length %d\n", a, steps, len(path))
			paths = append(paths, path)
			colours = append(colours, c)
			m.UnVisitAll()
		}

		err := m.ImageWithMultiplePaths(paths, fileOut, colours)
		if err != nil {
			fmt.Println("Could not create image:", err.Error())
			os.Exit(1)
		}

		fmt.Println("Image done")
		os.Exit(0)
	}

	err := m.Image(fileOut)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Image done")
}
