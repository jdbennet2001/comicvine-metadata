package main

import (
	"comicvine-metadata/comicvine"
	"comicvine-metadata/data"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	// Define flags for dir and apiKey
	dir := flag.String("dir", "", "Directory argument")
	apiKey := flag.String("apiKey", "", "API Key argument")

	// Parse the command-line arguments
	flag.Parse()

	// Check if the required arguments are provided
	if *dir == "" || *apiKey == "" {
		fmt.Println("Usage: go run comicvine-metadata.go -dir <directory> -apiKey <apiKey>")
		return
	}

	// Sanity check, just exit if the backing DB doesn't exist
	if _, e := os.Stat(*dir); os.IsNotExist(e) {
		log.Fatal("Metadata dir: ", *dir, " not present")
	}

	// Access the values of dir and apiKey
	fmt.Println("Directory:", *dir)
	fmt.Println("API Key:", *apiKey)

	//// Download new issues
	issues := comicvine.UpdateIssueData(*dir, *apiKey)

	// Use issue data to download volume / cover / hashes ...
	comicvine.UpdateCovers(*dir, issues)
	volumes := comicvine.UpdateVolumeData(*dir, issues, *apiKey)
	hashes := comicvine.UpdateHashData(*dir, issues)
	//
	//// Summarize the data
	summary := data.SummaryRecords(issues, volumes, hashes)

	// And store
	file, _ := json.MarshalIndent(summary, "", "\t")
	err := ioutil.WriteFile("index.json", file, 0644)
	if err == nil {
		msg := fmt.Sprintf("Done! %d records available.", len(summary))
		fmt.Println(msg)
	} else {
		fmt.Println("Error serializing data: ", err)
	}

	fmt.Println(len(issues), ".. done")

}
