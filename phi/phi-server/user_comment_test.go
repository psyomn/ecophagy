package main

import (
	"testing"
)

func TestUserCommentHasCorrectFormat(t *testing.T) {
	phi := userComment{
		Phi: phi{
			Username:  "the-username",
			Timestamp: 123,
			Tags:      []string{"one", "two", "three"},
		},
	}

	bytes := phi.toJSON()

	if len(bytes) == 0 {
		log.Fail("should generate some json")
	}
}
