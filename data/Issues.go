package data

import (
	"comicvine-metadata/comicvine"
	"comicvine-metadata/fsutils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func IssueData(rootDir string) []comicvine.Issue {

	directory := filepath.Join(rootDir, "issues")

	var results []comicvine.Issue

	// Get a list of all covers
	issueFiles := fsutils.Walk(directory, ".json")
	print(len(issueFiles), " covers")

	for _, issueFile := range issueFiles {
		issues := loadIssues(issueFile)
		fmt.Println("Loading issue data from: ", filepath.Base(issueFile), ".. ", len(issues), " issues")
		results = append(results, issues...)
	}

	return results

}

func loadIssues(filename string) []comicvine.Issue {

	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users []comicvine.Issue

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &users)

	return users

}
