package data

import (
	"comicvine-metadata/comicvine"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Generate a JSON index of all issues available for processing
func getSummaryRecords(root string, issueData []comicvine.Issue, volumes []comicvine.Volume) []SummaryRecord {

	var summary []SummaryRecord

	// Remove publishers known to be bad
	//skippedPublishers := map[string]bool{"Panini Comics": true,
	//	"ECC Ediciones":      true,
	//	"Norma Editorial":    true,
	//	"No Comprendo Press": true,
	//	"Timely":             true,
	//	"Planeta DeAgostini": true,
	//	"Toutain Editor":     true,
	//	"Dupuis":             true,

	for ctr, issue := range issueData {

		volume, err := getVolume(issue.Volume.ID, volumes)

		if err != nil {
			fmt.Println("Skipping ", issue.Name, " : ", issue.IssueNumber+", no volume.")
			continue
		}

		if ctr%100 == 0 {
			fmt.Println("..", ctr, " covers processed.")
		}

		var record = SummaryRecord{
			ID:          issue.ID,
			Name:        issue.Name,
			IssueNumber: issue.IssueNumber,
			Image:       issue.Image.OriginalURL,
			CoverDate:   issue.CoverDate,
			IssueSource: issue.SiteDetailURL,
			VolumeName:  issue.Volume.Name,
			Publisher:   volume.Publisher.Name,
			VolumeCount: volume.CountOfIssues,
			//Hash:        ,
			VolumeID:    strconv.Itoa(volume.Id),
			VolumeStart: volume.StartYear,
		}

		summary = append(summary, record)
	}

	return summary

}

func getVolume(volumeId int, volumes []comicvine.Volume) (comicvine.Volume, error) {
	for _, volume := range volumes {
		if volume.Id == volumeId {
			return volume, nil
		}
	}

	return comicvine.Volume{}, errors.New("Missing volume: " + strconv.Itoa(volumeId))
}

// Write summary data to disk
func saveClassificationMetaData(indexLocation string, records []SummaryRecord) error {
	jsonData, _ := json.MarshalIndent(records, "", "\t")
	err := ioutil.WriteFile(indexLocation, jsonData, 0644)
	return err
}

func normalize(str string) string {
	// Clean up known special characters
	str = strings.ReplaceAll(str, "?", "")
	str = strings.ReplaceAll(str, "/", "-")
	str = strings.ReplaceAll(str, ":", " ")
	str = strings.ReplaceAll(str, "\"", "'")
	str = strings.ReplaceAll(str, "|", "-")
	str = strings.ReplaceAll(str, "...", " ")

	return str
}
