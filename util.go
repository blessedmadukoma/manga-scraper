package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// convertToInt converts string variable to int
func convertToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Error converting %s to int: %s \n", s, err)
	}
	return i
}

// isContainString checks if the image cotains certain URLs which makes it downloadable
func isContainString(s string) bool {
	return strings.Contains(s, "heaven") || strings.Contains(s, "fun") || strings.Contains(s, "manga") || (strings.Contains(s, "mytoon.net/images")) || (strings.Contains(s, "mytoon.net/img"))
}

// isExist cheks if the file already exists in the directory
func isExist(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
