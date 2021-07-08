/*
Copyright 2018 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"testing"

	"github.com/psyomn/ecophagy/cynic/lib"
)

func TestSimpleBuilder(t *testing.T) {
	setup := func(eventCount, maxTime int) func(t *testing.T) {
		return func(t *testing.T) {
			var events []cynic.Event

			for i := 0; i < eventCount; i++ {
				event := cynic.EventNew(1)
				events = append(events, event)
			}

			builder := cynic.EventBuilderNew(events)
			builder.DistributeEvents(maxTime)

			session, ok := builder.Build()
			assert(t, ok)

			for _, el := range session.Events {
				assert(t, el.GetSecs() == (maxTime/eventCount))
			}
		}
	}

	type testCase struct {
		name     string
		serCount int
		maxTime  int
	}

	testCases := [...]testCase{
		{"maxtime 5, event count 5", 5, 5},
		{"maxtime 1000 event count 100", 100, 1000},
		{"maxtime 999 event count 100", 100, 999},
	}

	for _, c := range testCases {
		t.Run(c.name, setup(c.serCount, c.maxTime))
	}
}

func TestSimpleErrorCases(t *testing.T) {
	setup := func(eventCount, maxTime int) func(t *testing.T) {
		return func(t *testing.T) {
			var events []cynic.Event
			builder := cynic.EventBuilderNew(events)
			_, ok := builder.Build()
			assert(t, !ok)
		}
	}

	type testCase struct {
		name     string
		serCount int
		maxTime  int
	}

	tests := [...]testCase{
		{"maxtime -10 event count 1", 1, -10},
		{"maxtime 0 event count 10", 0, 10},
		{"maxtime 10 event count 11", 10, 11},
	}

	for _, c := range tests {
		t.Run(c.name, setup(c.serCount, c.maxTime))
	}
}
