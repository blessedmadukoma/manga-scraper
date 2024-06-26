package main

import (
	"bufio"
	"fmt"
	"strings"
)

// getMangaName retrieves the manga name from user input
func getMangaName(rd *bufio.Reader) (string, string) {
	fmt.Println("Enter manga name: e.g. solo leveling")
	mangaName, _ := rd.ReadString('\n')
	dirName := strings.Title(mangaName)

	dirName = strings.ReplaceAll(dirName, "?", "")

	mangaName = strings.ReplaceAll(mangaName, " ", "+")
	mangaName = strings.ReplaceAll(mangaName, "\n", "")

	return mangaName, dirName
}

// Trim the manga name
func trimMangaName(mangaName string) string {
	mangaName = strings.ToLower(mangaName)
	mangaName = strings.ReplaceAll(mangaName, " ", "-")
	mangaName = strings.ReplaceAll(mangaName, "+", "-")
	mangaName = strings.ReplaceAll(mangaName, "'", "")
	mangaName = strings.ReplaceAll(mangaName, "~", "-")
	mangaName = strings.ReplaceAll(mangaName, " ", "-")
	mangaName = strings.ReplaceAll(mangaName, ".", "-")
	mangaName = strings.ReplaceAll(mangaName, ", ", "-")
	mangaName = strings.ReplaceAll(mangaName, ",", "-")
	mangaName = strings.ReplaceAll(mangaName, "?", "-")
	mangaName = strings.ReplaceAll(mangaName, " ? ", "-")
	mangaName = strings.ReplaceAll(mangaName, ". ", "-")
	mangaName = strings.ReplaceAll(mangaName, ":", "-")
	mangaName = strings.ReplaceAll(mangaName, ": ", "-")
	mangaName = strings.ReplaceAll(mangaName, " - ", "-")
	mangaName = strings.ReplaceAll(mangaName, "--", "-")
	mangaName = strings.ReplaceAll(mangaName, "\n", "")
	mangaName = strings.ReplaceAll(mangaName, "’", "")

	return mangaName
}
