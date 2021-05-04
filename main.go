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
	link := "https://ww9.koomanga.com"
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
	fmt.Println("Enter the first chapter number for download e.g. 134 or 01 (if first episode starts with 01) or 1 (if first episode starts with 1, not 01)\n[To be sure, check koomanga.com, search for your manga and check the number of the first episode.]:")
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

		// Getting the URL
		fmt.Println("\nLoading URL!!!")
		time.Sleep(3 * time.Second) // waiting for page to load depending on the internet speed
		response, err = http.Get(url)
		if err != nil {
			fmt.Println("Error when getting the new url for the mangas:", err)
			return
		}

		// Change to the manga directory
		os.Chdir(dir)

		fmt.Println("Chapter", strconv.Itoa(i), "starting download!")
		fmt.Println()

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

			// check if the link exists and has a "heaven" in it
			if exists && (strings.Contains(imgSrc, "heaven") || strings.Contains(imgSrc, "fun") || strings.Contains(imgSrc, "manga") || strings.Contains(imgSrc, "image")) {
				mangaImgSrc = append(mangaImgSrc, imgSrc)

				fileName := "page_" + strconv.Itoa(j) + ".jpg"
				if fileExists(fileName) {
					j++
					fmt.Println("Skipping ", fileName, " since it exists...")
					goto DOWNLOADCODE
				}
				fmt.Println("Waiting for 2 seconds!")
				time.Sleep(2 * time.Second)
				fmt.Println(fileName, "download started!")
				fmt.Println("Download URL:", imgSrc)

				err = downloadFile(imgSrc, fileName)
				if err != nil {
					log.Fatal(err)
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

// Download the file
func downloadFile(URL, fileName string) error {
	out, err := os.Create(fileName + ".tmp")
	if err != nil {
		return err
	}

	response, err := http.Get(URL)
	if err != nil {
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
	mangaName = strings.ReplaceAll(mangaName, ". ", "-")
	mangaName = strings.ReplaceAll(mangaName, ":", "-")
	mangaName = strings.ReplaceAll(mangaName, ": ", "-")
	mangaName = strings.ReplaceAll(mangaName, " - ", "-")
	mangaName = strings.ReplaceAll(mangaName, "--", "-")
	mangaName = strings.ReplaceAll(mangaName, "\n", "")

	return mangaName
}
