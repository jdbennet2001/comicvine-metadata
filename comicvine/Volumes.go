package comicvine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func UpdateVolumeData(rootDir string, issues []Issue, apiKey string) map[int]Volume {

	volumes := loadVolumes(rootDir)

	for _, issue := range issues {
		_, ok := volumes[issue.Volume.ID]
		if ok {
			continue
		}

		volume := getVolumeData(issue, apiKey)
		volumes[volume.Id] = volume
	}

	saveVolumes(rootDir, volumes)
	return volumes

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

// Return information for a given volume
func getVolumeData(issue Issue, apiKey string) Volume {

	url := issue.Volume.APIDetailURL
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

func loadVolumes(rootDir string) map[int]Volume {

	location := filepath.Join(rootDir, "volumes.json")

	if _, e := os.Stat(location); os.IsNotExist(e) {
		return make(map[int]Volume)
	}

	var data map[int]Volume
	file, _ := ioutil.ReadFile(location)
	json.Unmarshal(file, &data)

	return data
}

func saveVolumes(rootDir string, volumes map[int]Volume) error {

	location := filepath.Join(rootDir, "volumes.json")

	file, _ := json.MarshalIndent(volumes, "", "\t")
	err := ioutil.WriteFile(location, file, 0644)
	return err
}
