package main

import (
	"math/rand"
)

func sampleArray(arr []string) string {
	// tempting to use unsafe here, but I know that this will be
	// used on a windows machine, and not too sure how that will
	// behave there.

	index := rand.Intn(len(arr))
	return arr[index]
}
