package main

import (
	"io/ioutil"
	"strings"
)

type Set map[string]struct{}

func getEmojiSet(fileName string) Set {
	ebts, _ := ioutil.ReadFile(fileName)
	out := Set{}
	for _, e := range strings.Fields(string(ebts)) {
		out[e] = struct{}{}
	}

	return out
}

func (s *Set) Has(elem string) bool {
	_, ok := (*s)[elem]
	return ok
}

func (s *Set) Union(s2 Set) {
	for k := range s2 {
		(*s)[k] = struct{}{}
	}
}
