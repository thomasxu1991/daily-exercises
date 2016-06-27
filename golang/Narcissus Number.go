package main

import (
	"fmt"
	"math"
	"time"
)

func check(num uint64) bool {
	tmp := num
	length := 0;
	slices := make([]uint64, 0, 0)
	for length = 0; num > 0; length++ {
		slices = append(slices, uint64(num % 10))
		num = num / 10
	}
	indexSum := uint64(0);
	for _, value := range slices {
		indexSum += uint64(math.Pow(float64(value), float64(length)))
	}
	if tmp == uint64(indexSum) {
		return true
	} else {
		return false
	}
}

func main() {
	for i := uint64(100); i < 8999999999999999999; i++ {
		if (check(i)) {
			fmt.Println(time.Now(),i)
		}
	}

}
