package main

import (
	"github.com/qiwenilli/goyacc"
)

func main() {

	strs := []string{
		"str eq '123行动\\''",
		"str eq '123行动\\''",
		"ab ne 0.123",
		"ab ne 1110.123",
		"ab ne -1110.123",
		"str ne -1123",
		"st ne 100",
	}
	for _, str := range strs {
		goyacc.Parse(str)
	}

}
