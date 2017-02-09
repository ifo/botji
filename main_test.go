package main

import "testing"

func Test_parseEmoji(t *testing.T) {
	emoji = getEmojiSet("emoji.txt")

	cases := map[int]struct {
		In  string
		Out Set
	}{
		1: {In: "airplane_arrival", Out: Set{"airplane_arrival": struct{}{}}},
		2: {In: "some Alien Monster", Out: Set{"alien_monster": struct{}{}}},
		3: {In: "some Alien_Monster", Out: Set{"alien_monster": struct{}{}}},
	}

	for n, c := range cases {
		set := parseEmoji(c.In)
		if len(set) != len(c.Out) {
			t.Fatalf("Actual: %+v, Expected: %+v, Case: %d", set, c.Out, n)
		}
		for k := range set {
			if _, exist := c.Out[k]; !exist {
				t.Fatalf("Got but did not have: %s, Case: %d", k, n)
			}
		}
	}
}
