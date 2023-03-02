package main

import (
	"fmt"
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
	return strings.Contains(s, "heaven") || (strings.Contains(s, "mytoon.net/images")) || (strings.Contains(s, "mytoon.net/img")) || (strings.Contains(s, "mytoon.net/cloud")) || (strings.Contains(s, "funmanga.com/uploads/chapter_files")) || (strings.Contains(s, "mytoon.net/uploads"))	
}

// isExist cheks if the file already exists in the directory
func isExist(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// isSpecialCharacter checks if the manga has special characters in the links such as mbx11, mbx12 e.g. Eleceed
func isContainsCharacter(url, mangaName string, i int) string {
	fmt.Println("Herrrrrrrerererere:", url)

	switch {
	// case strings.Contains(url, "eleceed"), strings.Contains(url, "chainsaw man"):
	case strings.Contains(url, "eleceed"):
		if strings.Contains(url, "200") {
			url = fmt.Sprintf("%s/mbx11-%s-chapter-%s/", link, mangaName, strconv.Itoa(i))
		}
		if strings.Contains(url, "214") {
			url = fmt.Sprintf("%s/mbx15-%s-chapter-%s/", link, mangaName, strconv.Itoa(i))
		}
		if strings.Contains(url, "202") || strings.Contains(url, "203") || strings.Contains(url, "204") || strings.Contains(url, "205") || strings.Contains(url, "206") || strings.Contains(url, "207") || strings.Contains(url, "208") || strings.Contains(url, "209") || strings.Contains(url, "210") {
			url = fmt.Sprintf("%s/mbx12-%s-chapter-%s/", link, mangaName, strconv.Itoa(i))
		}
		if strings.Contains(url, "211") || strings.Contains(url, "212") {
			url = fmt.Sprintf("%s/mbx14-%s-chapter-%s/", link, mangaName, strconv.Itoa(i))
		}

		fmt.Println("URL here:", url)
		break
	default:
		url = fmt.Sprintf("%s/%s-chapter-%s/", link, mangaName, strconv.Itoa(i))
		break
	}

	return url
}
