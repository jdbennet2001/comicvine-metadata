package main

import (
	"comicvine-metadata/comicvine"
	"comicvine-metadata/data"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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

	// Download new issues
	//comicvine.UpdateIssueData(*dir, *apiKey)
	issues := data.IssueData(*dir)

	//Download new volumes
	comicvine.UpdateVolumeData(*dir, *apiKey)
	volumes := data.VolumeData(*dir)
	comicvine.UpdateMissingVolumes(*dir, issues, volumes, *apiKey)

	//comicvine.UpdateCovers(*dir, issues)

	hashes := data.HashData(*dir, issues)

	fmt.Println("Done!, " + strconv.Itoa(len(hashes)) + " hashes")

}
