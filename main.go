package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
)

// Main function
func main() {
	link := "https://ww4.beetoon.net"
	response, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var j, totalPages, firstChapter, lastChapter int
	var url string
	var mangaImgSrc []string

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter manga name: e.g. solo leveling")
	mangaName, _ := reader.ReadString('\n')
	dirName := strings.ToUpper(mangaName)

	dir := dirName
	os.Mkdir(dir, 0755)

	mangaName = trimMangaName(mangaName)

	fmt.Println("Enter the number of chapters you want to download e.g. 3..(if it's only one chapter you want, input 0):")
	fmt.Scanln(&lastChapter)
	fmt.Println("Enter the first chapter number for download e.g. 134 or 01 (if first episode starts with 01) or 1 (if first episode starts with 1, not 01)\n[To be sure, check https://ww4.beetoon.net, search for your manga and check the number of the first episode.]:")
	variable, _ := reader.ReadString('\n')
	variable = strings.ReplaceAll(variable, "\n", "")

	firstChapter, _ = strconv.Atoi(variable)

	for i := firstChapter; i <= (firstChapter + lastChapter); i++ {
		j = 0
		if strings.HasPrefix(variable, "0") {
			url = link + "/" + mangaName + "-chap-" + variable + "/"
		} else {
			url = link + "/" + mangaName + "-chap-" + strconv.Itoa(i) + "/"
		}

		fmt.Printf("\nLoading URL %s!\n", url)
		time.Sleep(5 * time.Second) // waiting for page to load depending on the internet speed

		response, err = http.Get(url)
		if err != nil {
			fmt.Println("Error getting manga url:", err)
			return
		}

		// Change to the manga directory
		os.Chdir(dir)

		fmt.Println("Chapter", strconv.Itoa(i), "starting download!")

		// Make a directory with a chapter subdirectory, 0755 is the permision
		chapter := "chapter_" + strconv.Itoa(i)
		os.Mkdir(chapter, 0755)

		chapter = strings.ReplaceAll(chapter, "\n", "")
		chapter = strings.ReplaceAll(chapter, " ", "")
		dirChapter := chapter

		err := os.Chdir(dirChapter)
		if err != nil {
			fmt.Println("\nError changing chapter directory:", err)
		}

		document, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal("Error loading HTTP response body. ", err)
		}

		// select all the image tags
		document.Find("img").Each(func(index int, element *goquery.Selection) {
			// select all image tags with src attribute
			imgSrc, exists := element.Attr("src")

			if strings.Contains(imgSrc, "ads") || strings.Contains(imgSrc, "content/frontend") {
				document.Next()
			}

			// check if the link exists and has a contains the URLs
			if exists && containString(imgSrc) {
				mangaImgSrc = append(mangaImgSrc, imgSrc)

				fileName := "page_" + strconv.Itoa(j) + ".jpg"
				if fileExists(fileName) {
					j++
					fmt.Printf("%s already exists, skipping...\n", fileName)
					goto DOWNLOADCODE
				}
				fmt.Println("Waiting for 2 seconds!")
				time.Sleep(2 * time.Second)
				fmt.Printf("%s with URL %s download started!\n", fileName, imgSrc)

				err = downloadFile(imgSrc, fileName)
				if err != nil {
					log.Fatalf("Error downloading %s: %s\n", fileName, err)
				}
				fmt.Printf("=> %s download finished!\n", fileName)
				j++
			DOWNLOADCODE:
			}
		})
		fmt.Print("\n")

		err = os.Chdir("../../" + dir)
		if err != nil {
			log.Fatal("Error changing to original directory: ", err)
		}

		totalPages = j
		fmt.Println("Chapter ", strconv.Itoa(i), " downloaded with ", strconv.Itoa(totalPages), "pages downloaded")
	}
	fmt.Println("All Downloads Completed!!")
}

// WriteCounter counts the total kb downloaded
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress prints the progress
func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download. We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(wc.Total))
}

// containString checks if the image cotains certain URLs which makes it downloadable
func containString(imgSrc string) bool {
	return strings.Contains(imgSrc, "heaven") || strings.Contains(imgSrc, "fun") || strings.Contains(imgSrc, "manga") || (strings.Contains(imgSrc, "mytoon.net/images")) || (strings.Contains(imgSrc, "mytoon.net/img"))
}

// Download the file
func downloadFile(URL, fileName string) error {
	out, err := os.Create(fileName + ".tmp")
	if err != nil {
		return err
	}

	// remove whitespaces e.g. \n
	URL = strings.TrimSpace(URL)

	response, err := http.Get(URL)
	if err != nil {
		log.Println("Error getting URL:", err)
		out.Close()
		return err
	}
	defer response.Body.Close()

	// Create our progress reported and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(response.Body, counter)); err != nil {
		out.Close()
		return err
	}

	fmt.Print("\n")
	out.Close()

	if err = os.Rename(fileName+".tmp", fileName); err != nil {
		return err
	}
	return nil
}

// Function to check if the file already exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Trim the manga name
func trimMangaName(mangaName string) string {
	mangaName = strings.ToLower(mangaName)
	mangaName = strings.ReplaceAll(mangaName, " ", "-")
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

	return mangaName
}
