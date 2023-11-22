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

func VolumeData(rootDir string) []comicvine.Volume {

	directory := filepath.Join(rootDir, "volumes")

	var results []comicvine.Volume

	// Get a list of all covers
	volumeFiles := fsutils.Walk(directory, ".json")
	print(len(volumeFiles), " volume data sets")

	for _, volumeFile := range volumeFiles {
		volumes := loadVolumes(volumeFile)
		fmt.Println("Loading volume data from: ", filepath.Base(volumeFile), ".. ", len(volumes), " issues")
		results = append(results, volumes...)
	}

	return results
}

func loadVolumes(filename string) []comicvine.Volume {

	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var users []comicvine.Volume

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &users)

	return users

}
