package comicvine

// ComicVine metadata about a given issue
type Issue struct {
	APIDetailURL    string `json:"api_detail_url"`
	CoverDate       string `json:"cover_date"`
	DateAdded       string `json:"date_added"`
	DateLastUpdated string `json:"date_last_updated"`
	Description     string `json:"description"`
	ID              int    `json:"id"`
	Image           struct {
		IconURL        string `json:"icon_url"`
		MediumURL      string `json:"medium_url"`
		ScreenURL      string `json:"screen_url"`
		ScreenLargeURL string `json:"screen_large_url"`
		SmallURL       string `json:"small_url"`
		SuperURL       string `json:"super_url"`
		ThumbURL       string `json:"thumb_url"`
		TinyURL        string `json:"tiny_url"`
		OriginalURL    string `json:"original_url"`
		ImageTags      string `json:"image_tags"`
	} `json:"image"`
	IssueNumber   string `json:"issue_number"`
	Name          string `json:"name"`
	SiteDetailURL string `json:"site_detail_url"`
	StoreDate     string `json:"store_date"`
	Volume        struct {
		APIDetailURL  string `json:"api_detail_url"`
		ID            int    `json:"id"`
		Name          string `json:"name"`
		SiteDetailURL string `json:"site_detail_url"`
	} `json:"volume"`
}

// ComicVine metadata about a given volume
type Volume struct {
	ApiDetailUrl    string  `json:"api_detail_url"`
	CountOfIssues   int     `json:"count_of_issues"`
	DateAdded       string  `json:"date_added"`
	DateLastUpdated string  `json:"date_last_updated"`
	Deck            *string `json:"deck"`
	Description     string  `json:"description"`
	FirstIssue      struct {
		ApiDetailUrl string `json:"api_detail_url"`
		Id           int    `json:"id"`
		Name         string `json:"name"`
		IssueNumber  string `json:"issue_number"`
	} `json:"first_issue"`
	Id    int `json:"id"`
	Image struct {
		IconUrl        string `json:"icon_url"`
		MediumUrl      string `json:"medium_url"`
		ScreenUrl      string `json:"screen_url"`
		ScreenLargeUrl string `json:"screen_large_url"`
		SmallUrl       string `json:"small_url"`
		SuperUrl       string `json:"super_url"`
		ThumbUrl       string `json:"thumb_url"`
		TinyUrl        string `json:"tiny_url"`
		OriginalUrl    string `json:"original_url"`
		ImageTags      string `json:"image_tags"`
	} `json:"image"`
	LastIssue struct {
		ApiDetailUrl string `json:"api_detail_url"`
		Id           int    `json:"id"`
		Name         string `json:"name"`
		IssueNumber  string `json:"issue_number"`
	} `json:"last_issue"`
	Name      string `json:"name"`
	Publisher struct {
		ApiDetailUrl string `json:"api_detail_url"`
		Id           int    `json:"id"`
		Name         string `json:"name"`
	} `json:"publisher"`
	SiteDetailUrl string `json:"site_detail_url"`
	StartYear     string `json:"start_year"`
}
