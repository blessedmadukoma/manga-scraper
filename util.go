package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// createDir creates a directory named manga to house all downloadable manga
func createDir() string {
	dir := "Manga"
	err := os.Mkdir(dir, 0755)
	if err != nil {
		if err.Error() != "mkdir Manga: file exists" {
			log.Fatal("Error creating manga directory:", err)
		}
	}
	return dir
}

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
