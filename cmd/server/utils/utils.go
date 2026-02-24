package utils

import (
	"strings"
	"math/rand"
	"fmt"
)

func GenerateFileName() string {
	name := strings.Builder{}

	for range 6 {
		randIndex := rand.Intn(10)
		fmt.Fprintf(&name, "%d", randIndex)
	}
	
	return name.String()
}

func FileNameMatches(s, substring string) bool {
	return strings.Contains(s, substring)
}
