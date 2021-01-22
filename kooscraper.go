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

func main() {
	response, err := http.Get("https://ww6.koomanga.com")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var j, firstChapter, lastChapter int
	var url string
	var mangaImgSrc []string

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter manga name: e.g. solo leveling")
	mangaName, _ := reader.ReadString('\n')
	dirName := strings.ToUpper(mangaName)

	dir := dirName

	os.Mkdir(dir, 0755)

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

	fmt.Println("Enter the number of chapters you want to download e.g. 3 (if it's only one chapter you want or the first number of the chapter starts with 0 (i.e. 02), input 0):")
	fmt.Scanln(&lastChapter)
	fmt.Println("Enter the first chapter number for download e.g. 134 or 01 (if first episode starts with 01) or 1 (if first episode starts with 1 not 01)\n[To be sure, check koomanga.com, search for your manga, check the number of the first episode.]:")
	variable, _ := reader.ReadString('\n')
	variable = strings.ReplaceAll(variable, "\n", "")
	// fmt.Println("Variable:", variable)

	firstChapter, _ = strconv.Atoi(variable)

	for i := firstChapter; i <= (firstChapter + lastChapter); i++ {
		j = 1
		if strings.HasPrefix(variable, "0") {
			url = "https://ww6.koomanga.com/" + mangaName + "-chap-" + variable + "/"
		} else {
			url = "https://ww6.koomanga.com/" + mangaName + "-chap-" + strconv.Itoa(i) + "/"
		}

		// Getting the URL
		// fmt.Println("URL:", url)
		fmt.Println("Waiting for 6 seconds!")
		fmt.Println()
		time.Sleep(6 * time.Second) // waiting for page to load depending on the internet speed
		response, err = http.Get(url)
		if err != nil {
			fmt.Println("Error when getting the new url for the mangas:", err)
			return
		}

		err = os.Chdir(dir)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		pwd0, _ := os.Getwd()
		pwd0 = strings.ReplaceAll(pwd0, "\n", "")
		pwd0 = strings.ReplaceAll(pwd0, " ", "")
		// fmt.Println("First directory change pwd:", pwd0)

		fmt.Println("Chapter", strconv.Itoa(i), "starting download!")
		fmt.Println()

		// Make a directory with a chapter subdirectory, 0755 is the permision
		chapter := "chapter_" + strconv.Itoa(i)
		os.Mkdir(chapter, 0755)

		chapter = strings.ReplaceAll(chapter, "\n", "")
		chapter = strings.ReplaceAll(chapter, " ", "")
		dirChapter := chapter
		pwd, _ := os.Getwd()

		dirChapter = strings.ReplaceAll(dirChapter, "\n", "")
		dirChapter = strings.ReplaceAll(dirChapter, " ", "")
		pwd = strings.ReplaceAll(pwd, "\n", "")
		pwd = strings.ReplaceAll(pwd, " ", "")

		err := os.Chdir(dirChapter)
		pwd1, _ := os.Getwd()
		pwd1 = strings.ReplaceAll(pwd1, "\n", "")
		pwd1 = strings.ReplaceAll(pwd1, " ", "")
		if err != nil {
			fmt.Println("\nError:", err)
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

				// fmt.Println("j:", j, "\nURL:", imgSrc)
				fileName := "page_" + strconv.Itoa(j) + ".jpg"
				fmt.Println("Waiting for 3 seconds!")
				time.Sleep(3 * time.Second)
				fmt.Println(fileName, "download started!")
				err = downloadFile(imgSrc, fileName)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("=> %s download finished!\n", fileName)
				j++
			}
		})
		fmt.Println("Chapter", strconv.Itoa(i), "with", len(mangaImgSrc), "pages completely downloaded!!")

		fmt.Print("\n\n")

		err = os.Chdir("../../" + dir)
		if err != nil {
			fmt.Println(err)
		}
		pwd0, _ = os.Getwd()
		pwd0 = strings.ReplaceAll(pwd0, "\n", "")
		pwd0 = strings.ReplaceAll(pwd0, " ", "")

	}
	fmt.Println("All Downloads Completed!!")
}

// function returns an error or nil if no error
func downloadFile(URL, fileName string) error {
	out, err := os.Create(fileName + ".tmp")
	if err != nil {
		return err
	}

	// getting the url
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

	// Progress uses the same line, so print a new line
	fmt.Print("\n")

	// Close the file without defer so it happens before rename
	out.Close()

	if err = os.Rename(fileName+".tmp", fileName); err != nil {
		return err
	}
	return nil
}
