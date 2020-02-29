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
	var pathFind, fileOut, algosToCompare string
	flag.IntVar(&cells, "cells", 25, "The numbers of cell across and wide for the maze")
	flag.StringVar(&pathFind, "path-find", "", "The path finding algorithm to use available are [bfs, stack]")
	flag.StringVar(&algosToCompare, "compare-algos", "", "Comma separated list of algos to compare")
	flag.StringVar(&fileOut, "file-out", "", "Image file with the maze")
	flag.Parse()

	if pathFind != "" && algosToCompare != "" {
		fmt.Println("cannot provide both -path-find and -compare-algos")
		os.Exit(1)
	}

	m := maze.NewMaze(cells, cells)
	fmt.Println("Done creating maze; producing image")

	if fileOut == "" {
		fileOut = fmt.Sprintf("./out/maze-%dx%d-%d", cells, cells, time.Now().Unix())
	}

	if pathFind != "" {
		err := singlePathFind(pathFind, fileOut, m)
		if err != nil {
			fmt.Println("Failed:", err.Error())
			os.Exit(1)
		}

		return
	}

	if algosToCompare != "" {
		err := compareAlgos(algosToCompare, fileOut, m)
		if err != nil {
			fmt.Println("Failed:", err.Error())
			os.Exit(1)
		}

		return
	}

	err := m.Image(fileOut)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Image done")
}

func singlePathFind(algo, fileOut string, m *maze.Maze) error {
	var path []*maze.CellIndex
	var steps uint64
	var err error

	switch algo{
	case "bfs":
		path, steps, err = pathfinding.BFS(m)
	case "dfs":
		path, steps, err = pathfinding.DFS(m)
	case "astar":
		path, steps, err = pathfinding.Astart(m)
	default:
		return fmt.Errorf("unrecognized alogrithm: %s", algo)
	}

	if err != nil {
		return fmt.Errorf("%s failed to find path after %d steps", algo, steps)
	}

	err = m.ImageWithPath(path, fileOut, color.RGBA{
		R: 100,
		G: 0,
		B: 100,
		A: 255,
	})
	if err != nil {
		return fmt.Errorf("could not create image: %w", err.Error())
	}

	return nil
}

func compareAlgos(algosToCompare, fileOut string, m *maze.Maze) error {
	algos := strings.Split(algosToCompare, ",")
	paths := make([][]*maze.CellIndex, 0)
	colours := make([]color.Color, 0)
	for _, a := range algos {
		var out []*maze.CellIndex
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
		case "astar":
			out, steps, err = pathfinding.Astart(m)
			c = color.RGBA{
				R: 255,
				A: 200,
			}
		default:
			return fmt.Errorf("unknown algo `%s`", a)
		}

		if err != nil {
			return err
		}

		fmt.Printf("Path found using %s took %d steps path length %d\n", a, steps, len(out))
		paths = append(paths, out)
		colours = append(colours, c)
		m.UnVisitAll()
	}

	err := m.ImageWithMultiplePaths(paths, fileOut, colours)
	if err != nil {
		return fmt.Errorf("could not create image: %w", err)
	}

	return nil
}