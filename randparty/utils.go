package main

// TODO: consider moving this to root common, and having common
// capabilities for random sampling of arrays.

import (
	mrand "math/rand"
	"time"
)

func initializeRandEngine() {
	// we don't care for strength because this app just wants to shuffle
	// things, not be crypto secure.
	//nolint
	mrand.New(mrand.NewSource(time.Now().Unix()))
}

func sampleArray(arr []string) string {
	// tempting to use unsafe here, but I know that this will be
	// used on a windows machine, and not too sure how that will
	// behave there.

	// This complains of not very strong crypto rand, but we really
	// don't care for an application like this -- we just want random
	// words from a list, and that's good enough. One could argue that
	// maybe the sequences that are going to show up will not be as
	// unique, but again, I don't think it's worth investing on this
	// subject...

	//nolint
	index := mrand.Intn(len(arr))
	return arr[index]
}
