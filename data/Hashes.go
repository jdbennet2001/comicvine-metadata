package data

import (
	"comicvine-metadata/comicvine"
	"encoding/json"
	"fmt"
	"github.com/corona10/goimagehash"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

func HashData(rootDir string, issues []comicvine.Issue) map[int]string {

	fmt.Println("Checking covers and hash data.")

	hashes := loadHashes(rootDir)

	for ctr, issue := range issues {

		_, ok := hashes[issue.ID]
		if ok {
			continue // Already got the hash...
		}

		fmt.Println(".. ", strconv.Itoa(ctr), " hashing: ", strconv.Itoa(issue.ID), " : ", issue.Volume.Name+":", issue.IssueNumber)

		fileLocation := coverLocation(rootDir, issue)
		value, err := hash(fileLocation)
		if err != nil {
			fmt.Println("Hash error: ", err)
			continue
		}

		hashes[issue.ID] = value
		saveHashes(rootDir, hashes)

	}

	return hashes

}

// Return the location for a given issue
func coverLocation(root string, issue comicvine.Issue) string {

	var month = "0000-00"

	if len(issue.CoverDate) > 0 {
		month = issue.CoverDate[0:7]
	}

	cover := filepath.Join(root, "covers", month, strconv.Itoa(issue.ID)+".jpg")
	return cover

}

func loadHashes(rootDir string) map[int]string {

	cache := cacheLocation(rootDir)

	if _, e := os.Stat(cache); os.IsNotExist(e) {
		return make(map[int]string)
	}

	var data map[int]string
	file, _ := ioutil.ReadFile(cache)
	json.Unmarshal(file, &data)

	return data
}

// Return the hash value for a given file on disk
func hash(filepath string) (string, error) {

	file1, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file1.Close()

	img1, err := jpeg.Decode(file1)
	if err != nil {
		return "", err
	}

	hash1, _ := goimagehash.DifferenceHash(img1)
	value := strconv.FormatUint(hash1.GetHash(), 10)

	return value, nil
}

func saveHashes(rootDir string, hashData map[int]string) error {

	cache := cacheLocation(rootDir)

	file, _ := json.MarshalIndent(hashData, "", "\t")
	err := ioutil.WriteFile(cache, file, 0644)
	return err
}

func cacheLocation(rootDir string) string {
	cache := filepath.Join(rootDir, "covers", "hashes.json")
	return cache
}
