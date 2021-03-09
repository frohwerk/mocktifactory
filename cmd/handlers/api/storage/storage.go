package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/frohwerk/mocktifactory/cmd/handlers/repository"
	"github.com/frohwerk/mocktifactory/cmd/ignores"
)

var (
	ignored = ignores.New(
		ignores.Equal("_maven.repositories"),
		ignores.Equal("_remote.repositories"),
		ignores.HasSuffix(".sha1"),
	)
)

func CreateHandler(r *repository.Repository) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		encoder := json.NewEncoder(rw)
		encoder.SetIndent("", "  ")

		path := strings.TrimPrefix(req.URL.Path, "/api/storage/libs-release-local")
		log.Printf("Request for path: %v\n", path)
		location := filepath.Join(r.Base(), path)

		if fsinfo, err := os.Stat(location); os.IsNotExist(err) {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusNotFound)
			if err := encoder.Encode(ErrorResponse{Errors: []Error{{Status: 404, Message: "File not found"}}}); err != nil {
				log.Printf("error writing response message: %v\n", err)
			}
		} else if fsinfo.IsDir() {
			if files, err := ioutil.ReadDir(location); err != nil {
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusInternalServerError)
				if err := encoder.Encode(ErrorResponse{Errors: []Error{{Status: 404, Message: "File not found"}}}); err != nil {
					log.Printf("error writing response message: %v\n", err)
				}
			} else {
				rw.Header().Set("Content-Type", "application/vnd.org.jfrog.artifactory.storage.FolderInfo+json")
				rw.WriteHeader(http.StatusOK)
				fi := folderInfo{
					Repo:         "libs-release-local",
					Path:         path,
					Created:      fsinfo.ModTime(),
					CreatedBy:    "jenkins-devops-user",
					LastModified: fsinfo.ModTime(),
					ModifiedBy:   "jenkins-devops-user",
					LastUpdated:  fsinfo.ModTime(),
					Children:     make([]child, 0),
					Uri:          fmt.Sprintf("http://%s%s", req.Host, req.URL),
				}
				for _, file := range files {
					if !ignored.Matches(file.Name()) {
						fi.Children = append(fi.Children, child{Uri: fmt.Sprintf("/%v", file.Name()), Folder: file.IsDir()})
					}
				}
				if err := encoder.Encode(fi); err != nil {
					log.Printf("error writing response message: %v\n", err)
				}
			}
		} else {
			rw.Header().Set("Content-Type", "application/vnd.org.jfrog.artifactory.storage.FileInfo+json")
			rw.WriteHeader(http.StatusOK)
			checksums := getChecksums(location)
			fi := fileInfo{
				Repo:              "libs-release-local",
				Path:              path,
				Created:           fsinfo.ModTime(),
				CreatedBy:         "jenkins-devops-user",
				LastModified:      fsinfo.ModTime(),
				ModifiedBy:        "jenkins-devops-user",
				LastUpdated:       fsinfo.ModTime(),
				Checksums:         checksums,
				OriginalChecksums: checksums,
				Uri:               fmt.Sprintf("http://%s%s", req.Host, req.URL),
			}
			if err := encoder.Encode(fi); err != nil {
				log.Printf("error writing response message: %v\n", err)
			}
		}
	}
}
