package build

import "time"

var (
	commit  = ""
	date    = ""
	version = "dev"
	repoURL = ""
)

func init() {
	date, _ := time.Parse(time.RFC3339, date)

	Current = Build{
		Commit:     commit,
		Version:    version,
		Date:       date,
		RepoURL:    repoURL,
		CommitURL:  repoURL + "/tree/" + commit,
		LicenseURL: repoURL + "/blob/master/LICENSE",
		ReleaseURL: repoURL + "/releases/tag/" + version,
	}
	if repoURL == "" {
		Current.CommitURL = "#"
		Current.LicenseURL = "#"
		Current.ReleaseURL = "#"
	}
}

var Current Build

type Build struct {
	Commit     string    `json:"commit,omitempty"`
	Version    string    `json:"version,omitempty"`
	Date       time.Time `json:"date,omitempty"`
	RepoURL    string    `json:"repo_url,omitempty"`
	CommitURL  string    `json:"commit_url,omitempty"`
	LicenseURL string    `json:"license_url,omitempty"`
	ReleaseURL string    `json:"release_url,omitempty"`
}
