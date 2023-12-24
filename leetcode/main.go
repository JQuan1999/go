package main

import (
	"fmt"
)

type drawer interface {
	draw()
}

type Graph struct {
	name string
}

func (graph *Graph) draw() {
	fmt.Println("graph draw")
}

func (graph *Graph) Name() {
	fmt.Println("name: ", graph.name)
}

type Square struct {
	Graph
}

func (sq *Square) draw() {
	fmt.Printf("Square:%s draw\n", sq.name)
}

type Circle struct {
	Graph
}

func (circle *Circle) draw() {
	fmt.Printf("circle:%s draw\n", circle.name)
}

type Box struct {
	GraphList []drawer
}

func (box *Box) Show() {
	for _, graph := range box.GraphList {
		graph.draw()
	}
}

func (box *Box) Name() {
	for _, graph := range box.GraphList {
		if v, ok := graph.(*Graph); ok {
			v.Name()
		}
	}
}

func main() {

}
