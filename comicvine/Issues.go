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

// Pull all data before Jan 1st, 1938
const cutoffYear = 1938

const RATE_LIMIT = 10

type IssueResponse struct {
	Results              []Issue
	Limit                int    `json:"limit"`
	Version              string `json:"version"`
	StatusCode           int    `json:"status_code"`
	Error                string `json:"error"`
	NumberOfPageResults  int    `json:"number_of_page_results"`
	NumberOfTotalResults int    `json:"number_of_total_results"`
	Offset               int    `json:"offset"`
}

func UpdateIssueData(rootDir string, apiKey string) []Issue {

	months := months(cutoffYear)

	// Update any missing notes
	for _, month := range months {
		updateIssues(rootDir, month, apiKey)
	}

	// Update this month -> next year data
	updateExtra(rootDir, apiKey)

	return issueRecords(rootDir)

}

// Download issues for this / future months. Data can be tossed away as the updates march on...
func updateExtra(root string, apiKey string) error {
	t := time.Now()
	start := fmt.Sprintf("%d-%02d", t.Year(), int(t.Month()))
	end := fmt.Sprintf("%d-%02d", 2032, int(t.Month()))
	dateRange := fmt.Sprintf("%s-01|%s-31", start, end)
	data := issues(dateRange, apiKey)

	location := filepath.Join(root, "issues", "extra.json")
	fmt.Println(".. updating issue data for ", len(data), " issues")
	file, _ := json.MarshalIndent(data, "", "\t")
	err := ioutil.WriteFile(location, file, 0644)

	return err

}

// Download metadata for comics, between last month and
func updateIssues(root string, month string, apiKey string) error {

	location := filepath.Join(root, "issues", month+".json")
	if _, err := os.Stat(location); err == nil {
		fmt.Println(".. skipping ", location)
		return nil
	}

	fmt.Println(".. downloading issue data for ", month)
	dateRange := fmt.Sprintf("%s-01|%s-31", month, month)
	data := issues(dateRange, apiKey)
	file, _ := json.MarshalIndent(data, "", "\t")
	_ = ioutil.WriteFile(location, file, 0644)

	return nil
}

// Returns all issues for a given month
func issues(dateRange string, apiKey string) []Issue {

	var offset = 0

	var results []Issue

	var totalResults = int(math.Pow(2, 32))

	for len(results) < totalResults {

		url := fmt.Sprintf("https://comicvine.gamespot.com/api/issues/?api_key=%s&filter=cover_date:%s&format=json&offset=%d", apiKey, dateRange, offset)

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("User-Agent", "request")
		req.Header.Add("Accept-Encoding", "")
		req.Header.Add("Content-Type", "application/json")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var result IssueResponse

		if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
			fmt.Println("Can not unmarshal JSON", err)
			inputData := string(body)
			panic("Cannot unmarshal JSON data: " + inputData)
		}

		arr := result.Results
		totalResults = result.NumberOfTotalResults

		results = append(results, arr...)
		offset = offset + 100

		fmt.Println("Issue query: ", dateRange, ".. ", len(results), " of ", totalResults)

		time.Sleep(RATE_LIMIT * time.Second)

	}

	return results
}

// Returns all issues
func issueRecords(rootDir string) []Issue {

	var issues []Issue

	directory := filepath.Join(rootDir, "issues")
	files, _ := ioutil.ReadDir(directory) //read the files from the directory
	for _, file := range files {
		fileLocation := filepath.Join(directory, file.Name())
		data := loadIssueFile(fileLocation)
		issues = append(issues, data...)
	}

	return issues

}

func loadIssueFile(filename string) []Issue {

	fmt.Println(".. loading ", filename)

	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var results []Issue
	json.Unmarshal(byteValue, &results)
	return results
}
