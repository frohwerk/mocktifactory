package repository

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type createdResponse struct {
	repo   string
	path   string
	name   string
	sha256 string
	size   uint
}

func (r *Repository) upload(rw http.ResponseWriter, req *http.Request) {
	if a, err := r.handleUpload(rw, req); err != nil {
		rw.WriteHeader(err.Status)
		if err := json.NewEncoder(rw).Encode(err); err != nil {
			log.Printf("failed to send json response: %v", err)
		}
	} else {
		go r.Webhooks.ArtifactCreated(a.repo, a.path, a.name, a.sha256, a.size)
		rw.WriteHeader(http.StatusCreated)
	}
}

func (r *Repository) handleUpload(rw http.ResponseWriter, req *http.Request) (*createdResponse, *errorResponse) {
	if _, err := os.Stat(r.Base()); os.IsNotExist(err) {
		log.Printf("Base directory %s does not exist", r.Base())
		return nil, &errorResponse{http.StatusInternalServerError, err.Error()}
	}

	path := strings.TrimPrefix(req.URL.Path, "/libs-release-local")
	log.Printf("Base directory: %s", path)
	file := filepath.Join(r.Base(), path)
	log.Printf("File location: %s", file)
	dir := filepath.Dir(file)
	log.Printf("File directory: %s", dir)

	switch _, err := os.Stat(dir); {
	case os.IsNotExist(err):
		if err := os.MkdirAll(dir, 0775); err != nil {
			log.Printf("Failed to create directory %s", dir)
			return nil, &errorResponse{http.StatusInternalServerError, err.Error()}
		}
	case err != nil:
		return nil, &errorResponse{http.StatusInternalServerError, err.Error()}
	}

	f, err := os.Create(file)
	if err, ok := err.(*os.PathError); ok {
		if err, ok := err.Err.(syscall.Errno); ok && err == 0x7b {
			return nil, &errorResponse{http.StatusBadRequest, fmt.Sprintf("Invalid file name or path: %v", path)}
		}
	} else if err != nil {
		log.Printf("Failed to create file %s", file)
		return nil, &errorResponse{http.StatusInternalServerError, err.Error()}
	}
	defer f.Close()

	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, &errorResponse{http.StatusInternalServerError, err.Error()}
	}
	req.Body.Close()

	size, err := io.Copy(f, bytes.NewReader(buf))
	if err != nil {
		return nil, &errorResponse{http.StatusInternalServerError, err.Error()}
	}

	sha256 := sha256.Sum256(buf)
	return &createdResponse{
		repo:   "libs-release-local",
		path:   path,
		name:   path[strings.LastIndex(path, "/")+1:],
		sha256: hex.EncodeToString(sha256[:]),
		size:   uint(size),
	}, nil
}
