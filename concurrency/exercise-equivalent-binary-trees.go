package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		Walk(t1, ch1)
		close(ch1)
	}()
	go func() {
		Walk(t2, ch2)
		close(ch2)
	}()

	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2
		if ok1 != ok2 || v1 != v2 {
			return false
		}
		if !ok1 {
			break
		}
	}
	return true
}

func main() {
	ch := make(chan int)
	go func() {
		Walk(tree.New(1), ch)
		close(ch)
	}()
	fmt.Print("Values tree1: ")
	for v := range ch {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// test
	fmt.Println("Same(tree1, tree1) =", Same(tree.New(1), tree.New(1)))
	fmt.Println("Same(tree1, tree2) =", Same(tree.New(1), tree.New(2)))
}
