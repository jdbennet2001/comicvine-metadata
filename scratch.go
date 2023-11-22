package main

import (
	"bufio"
	"comic-match/archive"
	"comic-match/fsutils"
	"comic-match/stringUtils"
	"errors"


	"os/exec"

	"comic-match/catalog"
	"comic-match/metadata"
	"encoding/json"
	"fmt"

	"os"
	"path/filepath"
	"strings"
)


const cvDrive = "/Volumes/Seagate/comicVine"
const catalogDir = "/Volumes/Storage/comics"


const idLength = 5
const tpbPages = 100		// Dividing line between trade and single issues

var nas catalog.Catalog
var dataStore metadata.ComicVineDB

var coverSummaries map[string]metadata.SummaryRecord = make(map[string]metadata.SummaryRecord)
var stringSummaries map[string]metadata.SummaryRecord = make(map[string]metadata.SummaryRecord)


func main() {


	// Get classification records
	dataStore = metadata.NewDataStore(cvDrive)

	// What are we processing?
	nas = catalog.NewCatalog(catalogDir)
	pendingComics := nas.OnDeck()

	// Batch all ML classification logic to reduce user interaction time
	for ctr, pendingComic := range pendingComics{
		fmt.Println(".. preprocessing ", ctr, " : ", pendingComic)
		coverSummaries[pendingComic], _ = dataStore.ClosestCover(pendingComic)
		stringSummaries[pendingComic] = dataStore.ClosestString(pendingComic)
	}

	fmt.Println(".. initialization done. ")

	// File comics that don't need human intervention
	for _, pendingComic := range pendingComics {

		pages, _ := archive.Pages(pendingComic)
		if len(pages) == 0{
			continue
		}

		// Import trades to known series
		series, _ := nas.Series(pendingComic)
		if len(pages) >= tpbPages && len(series) > 0 {
			fmt.Println(".. importing trade: ", filepath.Base(pendingComic))
			nas.ImportComic(pendingComic, metadata.SummaryRecord{})
			continue
		}

		summary, _ := coverSummaries[pendingComic]


		// Import high confidence matches
		if stringUtils.StringMatch(pendingComic, summary.VolumeName) && strings.Contains(pendingComic, summary.IssueNumber) {
			fmt.Println(filepath.Base(pendingComic), " - ", len(pages), " : ", series)
			fmt.Println(prettyPrint(summary))
			nas.ImportComic(pendingComic, summary)
			continue
		}
	}

	fmt.Println(".. automatic processing done done. ")

	for _, pendingComic := range pendingComics {

		if _, err := os.Stat(pendingComic); errors.Is(err, os.ErrNotExist) {
			continue
		}

		pages, _ := archive.Pages(pendingComic)
		if len(pages) == 0{
			continue
		}

		//Use string matching for the comic..
		summary := coverSummaries[pendingComic]
		promptIssue(pendingComic, summary)

	}

	fsutils.RemoveEmptyDirs(catalogDir)
}

func clear(){
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}

func promptIssue( pendingComic string, summary metadata.SummaryRecord){

	clear()

	pages, _ := archive.Pages(pendingComic)

	fmt.Println(filepath.Base(pendingComic), " - ", len(pages))
	fmt.Println(prettyPrint(summary))

	// Sanity check, blank record
	if summary.VolumeName == "" || summary.CoverDate == "" {
		return
	}

	prompt := StringPrompt("(i)mport, (s)kip, (m)atch title or enter new id")
	if len(prompt) >= idLength {
		issueSummary, _ := dataStore.Issue(prompt)
		promptIssue(pendingComic, issueSummary)
	}else if prompt == "m"{
		issueSummary := stringSummaries[pendingComic]
		promptIssue(pendingComic, issueSummary)
	}else if prompt == "i"{
		nas.ImportComic(pendingComic, summary)
	}

	return

}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}