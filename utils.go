package jscoreWorker

import "math/rand"

var uniqueIntGenerated = make(map[int32]bool)

func uniqueInt() int {
	for {
		i := rand.Int31()
		if !uniqueIntGenerated[i] {
			uniqueIntGenerated[i] = true
			return int(i)
		}
	}
}
