package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func mangaread() {
	link := "https://mangaread.org/manga"
	response, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var j, totalPages, firstChapter, lastChapter int
	var url string
	var mangaImgSrc []string
	var fileName string

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter manga name: e.g. solo leveling")
	mangaName, _ := reader.ReadString('\n')
	dirName := strings.ToUpper(mangaName)

	dir := dirName
	os.Mkdir(dir, 0755)

	mangaName = trimMangaName(mangaName)

	fmt.Println("Enter the number of chapters you want to download e.g. 3..(if it's only one chapter you want, input 0):")
	fmt.Scanln(&lastChapter)
	fmt.Println("Enter the first chapter number for download e.g. 134 or 01 (if first episode starts with 01) or 1 (if first episode starts with 1, not 01)\n[To be sure, check mangaread.org/manga, search for your manga and check the number of the first episode.]:")
	variable, _ := reader.ReadString('\n')
	variable = strings.ReplaceAll(variable, "\n", "")

	firstChapter, _ = strconv.Atoi(variable)

	for i := firstChapter; i <= (firstChapter + lastChapter); i++ {
		j = 0
		if strings.HasPrefix(variable, "0") {
			url = link + "/" + mangaName + "/chapter-" + variable + "/"
		} else {
			url = link + "/" + mangaName + "/chapter-" + strconv.Itoa(i) + "/"
		}

		fmt.Println("URL:", url)
		// Getting the URL
		fmt.Println("Waiting for 10 seconds!!!")
		time.Sleep(10 * time.Second) // waiting for page to load depending on the internet speed
		response, err = http.Get(url)
		if err != nil {
			fmt.Println("Error getting the manga url:", err)
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
			fmt.Println("image SRC 1: ", imgSrc)

			if strings.Contains(imgSrc, "ads") || strings.Contains(imgSrc, "content/frontend") {
				document.Next()
			}
			// check if the link exists and the src of the link has some keywords in it

			if exists && (strings.Contains(imgSrc, "manga_") || strings.Contains(imgSrc, "WP-manga")) {
				fmt.Println("image SRC 2: ", imgSrc)
				mangaImgSrc = append(mangaImgSrc, imgSrc)

				if strings.Contains(imgSrc, ".jpeg") {
					fileName = "page_" + strconv.Itoa(j) + ".jpeg"
				} else if strings.Contains(imgSrc, ".jpg")  {
					fileName = "page_" + strconv.Itoa(j) + ".jpg"
				} else {
					fileName = "page_" + strconv.Itoa(j) + ".png"
				}

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
	return
}
