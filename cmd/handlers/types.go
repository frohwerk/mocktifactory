package handlers

import "time"

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}

type folderInfo struct {
	Repo         string    `json:"repo"`
	Path         string    `json:"path"`
	Created      time.Time `json:"created"`
	CreatedBy    string    `json:"createdBy"`
	LastModified time.Time `json:"lastModified"`
	ModifiedBy   string    `json:"modifiedBy"`
	LastUpdated  time.Time `json:"lastUpdated"`
	Children     []child   `json:"children"`
	Uri          string    `json:"uri"`
}

type fileInfo struct {
	Repo              string    `json:"repo"`
	Path              string    `json:"path"`
	Created           time.Time `json:"created"`
	CreatedBy         string    `json:"createdBy"`
	LastModified      time.Time `json:"lastModified"`
	ModifiedBy        string    `json:"modifiedBy"`
	LastUpdated       time.Time `json:"lastUpdated"`
	Checksums         checksums `json:"checksums"`
	OriginalChecksums checksums `json:"originalChecksums"`
	Uri               string    `json:"uri"`
}

type checksums struct {
	Sha1   string `json:"sha1"`
	Md5    string `json:"md5"`
	Sha256 string `json:"sh2561"`
}

type child struct {
	Uri    string `json:"uri"`
	Folder bool   `json:"folder"`
}
