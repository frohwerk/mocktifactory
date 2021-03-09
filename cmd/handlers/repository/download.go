package repository

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (r *Repository) download(rw http.ResponseWriter, req *http.Request) {
	if err := r.handleDownload(rw, req); err != nil {
		rw.Header().Set("Content-Type", "text/plain")
		rw.WriteHeader(err.Status)
		if _, err := rw.Write([]byte(err.Message)); err != nil {
			log.Printf("error writing response message: %v\n", err)
		}
	}
}

func (r *Repository) handleDownload(rw http.ResponseWriter, req *http.Request) *errorResponse {
	path := strings.TrimPrefix(req.URL.Path, "/libs-release-local")
	location := filepath.Join(r.Base(), path)

	fsinfo, err := os.Stat(location)
	switch {
	case os.IsNotExist(err):
		return &errorResponse{Status: http.StatusNotFound, Message: fmt.Sprintf("%v does not exist\n", req.URL.Path)}
	case fsinfo.IsDir():
		return &errorResponse{Status: http.StatusBadRequest, Message: fmt.Sprintf("%v is a directory\n", req.URL.Path)}
	}

	file, err := os.Open(location)
	if err != nil {
		return &errorResponse{Status: http.StatusInternalServerError, Message: fmt.Sprintf("error opening file '%v'\n", location)}
	}

	rw.Header().Set("Content-Type", "application/java-archive")
	rw.Header().Set("Content-Length", fmt.Sprint(fsinfo.Size()))
	rw.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"; filename*=''%v`, fsinfo.Name(), fsinfo.Name()))
	rw.Header().Set("Last-Modified", fsinfo.ModTime().Format(time.RFC1123))
	rw.Header().Set("ETag", etag(location))
	rw.WriteHeader(http.StatusOK)
	for remaining := fsinfo.Size(); remaining > 0; {
		written, err := io.CopyN(rw, file, 8192)
		if err != nil && err != io.EOF {
			log.Printf("error writing response message: %v\n", err)
		} else {
			remaining -= written
		}
	}
	return nil
}

func etag(filename string) string {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("error reading file '%v' for etag calculation: %v\n", filename, err)
		return fmt.Sprint(time.Now().Unix())
	}
	sha1 := sha1.Sum(buf)
	return hex.EncodeToString(sha1[:])
}
