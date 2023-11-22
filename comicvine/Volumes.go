package comicvine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type VolumeResponse struct {
	Error                string `json:"error"`
	Limit                int    `json:"limit"`
	Offset               int    `json:"offset"`
	NumberOfPageResults  int    `json:"number_of_page_results"`
	NumberOfTotalResults int    `json:"number_of_total_results"`
	StatusCode           int    `json:"status_code"`
	Results              []Volume
	Version              string `json:"version"`
}

type SingleVolumeResponse struct {
	Error                string `json:"error"`
	Limit                int    `json:"limit"`
	Offset               int    `json:"offset"`
	NumberOfPageResults  int    `json:"number_of_page_results"`
	NumberOfTotalResults int    `json:"number_of_total_results"`
	StatusCode           int    `json:"status_code"`
	Results              Volume
	Version              string `json:"version"`
}

func UpdateVolumeData(rootDir string, apiKey string) {

	months := months(cutoffYear)

	for _, month := range months {
		updateVolumes(rootDir, month, apiKey)
	}

}

func updateVolumes(root string, month string, apiKey string) error {

	location := filepath.Join(root, "volumes", month+".json")
	if _, err := os.Stat(location); err == nil {
		return nil
	}

	fmt.Println(".. downloading volume data for ", month)
	dateRange := fmt.Sprintf("%s-01|%s-31", month, month)
	data := volumes(dateRange, apiKey)
	file, _ := json.MarshalIndent(data, "", "\t")
	_ = ioutil.WriteFile(location, file, 0644)

	return nil
}

// Returns all volumes for a given month
func volumes(dateRange string, apiKey string) []Volume {

	var offset = 0

	var results []Volume

	var totalResults = int(math.Pow(2, 32))

	for len(results) < totalResults {

		url := fmt.Sprintf("https://comicvine.gamespot.com/api/volumes/?api_key=%s&filter=date_added:%s&format=json&offset=%d", apiKey, dateRange, offset)

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("User-Agent", "request")
		req.Header.Add("Accept-Encoding", "")
		req.Header.Add("Content-Type", "application/json")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var result VolumeResponse

		if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON", err)
			inputData := string(body)
			panic("Error unmarshalling JSON data for volume data, " + inputData)
		}

		arr := result.Results
		totalResults = result.NumberOfTotalResults

		results = append(results, arr...)
		offset = offset + 100

		fmt.Println("Volume query: ", dateRange, ".. ", len(results), " of ", totalResults)

		time.Sleep(RATE_LIMIT * time.Second)

	}

	return results
}

// ---------------------------------------------------------------------------------------------------- //

func UpdateMissingVolumes(rootDir string, issues []Issue, volumes []Volume, apiKey string) []Volume {

	vmap := volumeMap(volumes)

	var entries []Volume

	baseName := "extra-" + time.Now().Format("2006-01-02 15:04:05") + ".json"
	fullName := filepath.Join(rootDir, "volumes", baseName)

	for _, issue := range issues {

		volumeID := issue.Volume.ID
		if _, ok := vmap[volumeID]; ok {
			continue
		}

		volume := getVolumeData(issue.Volume.APIDetailURL, apiKey)
		vmap[volumeID] = volume
		entries = append(entries, volume)

		// Dump Hash data to disk
		jsonData, _ := json.MarshalIndent(entries, "", "    ")
		err := ioutil.WriteFile(fullName, jsonData, 0644)
		if err != nil {
			log.Fatal(err)
		}

	}

	return entries

}

// Build a lookup table of volumes id / volume data for better performance
func volumeMap(volumes []Volume) map[int]Volume {

	var vmap map[int]Volume = make(map[int]Volume)

	for _, volume := range volumes {
		vmap[volume.Id] = volume
	}

	return vmap
}

// Return information for a given volume
func getVolumeData(url string, apiKey string) Volume {

	authURL := fmt.Sprintf("%s?api_key=%s&format=json", url, apiKey)

	req, _ := http.NewRequest("GET", authURL, nil)
	req.Header.Add("User-Agent", "request")
	req.Header.Add("Accept-Encoding", "")
	req.Header.Add("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var result SingleVolumeResponse

	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON", err)
		inputData := string(body)
		panic("Error unmarshaling JSON data for volume data, " + inputData)
	}

	volume := result.Results

	fmt.Println(".. updating volume information for ", volume.Name, " (", volume.Publisher.Name, "), published in ", volume.StartYear, " (", volume.DateAdded, ")")
	time.Sleep(RATE_LIMIT * time.Second)

	return volume
}
