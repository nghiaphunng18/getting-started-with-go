package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
    diff := 1.0

    for diff >= 1e-10 {
        next := z - (z*z-x)/(2*z)
        diff = math.Abs(next - z)
        z = next
    }

    return z
}

func main() {
	fmt.Println(Sqrt(2))
}
