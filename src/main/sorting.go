package main

import (
	"fmt"
	"math/rand"
)

func quicksort(a []int) []int {
	if len(a) <= 1 {
		return a
	}
	left, right := 0, len(a)-1
	pivotIndex := rand.Int() % len(a)
	a[pivotIndex], a[right] = a[right], a[pivotIndex]
	for i := range a {
		if a[i] < a[right] {
			a[i], a[left] = a[left], a[i]
			left++
		}
	}
	a[left], a[right] = a[right], a[left]
	quicksort(a[:left])
	quicksort(a[left+1:])
	return a
}

func add(x, y int) int {
	return x + y
}

func fizzbuzz() int {
	for i := 0; i < 100; i++ {
		if i%15 == 0 {
			fmt.Println(`fizzbuzz`)
		} else if i%3 == 0 {
			fmt.Println(`fizz`)
		} else if i%5 == 0 {
			fmt.Println(`buzz`)
		} else {
			fmt.Println(i)
		}
	}
	return 0
}
