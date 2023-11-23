package data

// Classification information about a given issue (issue + volume + hash data)
type SummaryRecord struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IssueNumber string `json:"issue_number"`
	Image       string `json:"image"`
	CoverDate   string `json:"cover_date"`
	IssueSource string `json:"issue_source"`
	VolumeName  string `json:"volume_name"`
	Publisher   string `json:"publisher"`
	VolumeCount int    `json:"volume_count"`
	Hash        string `json:"hash"`
	VolumeID    string `json:"volume_id"`
	VolumeStart string `json:"volume_start"`
	Description string `json:"description"`
}
