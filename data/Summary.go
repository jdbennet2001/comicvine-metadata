package data

import (
	"comicvine-metadata/comicvine"
	"encoding/json"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

// Generate a JSON index of all issues available for processing
func SummaryRecords(issueData []comicvine.Issue, volumes map[int]comicvine.Volume, hashes map[int]string) []SummaryRecord {

	var summary []SummaryRecord

	// Remove publishers known to be bad
	skippedPublishers := map[string]bool{"Panini Comics": true,
		"ECC Ediciones":      true,
		"Norma Editorial":    true,
		"No Comprendo Press": true,
		"Timely":             true,
		"Planeta DeAgostini": true,
		"Toutain Editor":     true,
		"Dupuis":             true,
	}

	for _, issue := range issueData {

		volume := volumes[issue.Volume.ID]

		var record = SummaryRecord{
			ID:          issue.ID,
			Name:        normalize(issue.Name),
			IssueNumber: issue.IssueNumber,
			Image:       issue.Image.OriginalURL,
			CoverDate:   issue.CoverDate,
			IssueSource: issue.SiteDetailURL,
			VolumeName:  normalize(issue.Volume.Name),
			Publisher:   volume.Publisher.Name,
			VolumeCount: volume.CountOfIssues,
			Hash:        hashes[issue.ID],
			VolumeID:    strconv.Itoa(volume.Id),
			VolumeStart: volume.StartYear,
			Description: issue.Description,
		}

		// Drop the European publishers that push reprintss
		if skippedPublishers[record.Publisher] {
			continue
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
