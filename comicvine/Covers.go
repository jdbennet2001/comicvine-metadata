package comicvine

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func UpdateCovers(rootDir string, issues []Issue) {

	fmt.Println("Checking ", len(issues), " covers.")

	for ctr, issue := range issues {

		if ctr%1000 == 0 {
			fmt.Println(".. ", strconv.Itoa(ctr))
		}

		fileLocation := coverLocation(rootDir, issue)
		if _, e := os.Stat(fileLocation); os.IsNotExist(e) {
			fmt.Println(".. retrieve image ", ctr, " - ", issue.Volume.Name, " #"+issue.IssueNumber+" -> ", fileLocation)
			downloadCover(fileLocation, issue.Image.OriginalURL)
		} else {

		}

	}

}

// Return the location for a given issue
func coverLocation(root string, issue Issue) string {

	var month = "0000-00"

	if len(issue.CoverDate) > 0 {
		month = issue.CoverDate[0:7]
	}

	cover := filepath.Join(root, "covers", month, strconv.Itoa(issue.ID)+".jpg")
	return cover

}

func downloadCover(fileName string, URL string) error {

	//  Make sure the directory exists
	parent := filepath.Dir(fileName)
	err := os.MkdirAll(parent, 0755)

	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	//time.Sleep(2 * time.Second)

	return nil
}
