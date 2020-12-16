package handlers

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

	"github.com/frohwerk/mocktifactory/cmd/constants"
)

func DownloadHandler(rw http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/libs-release-local")
	log.Printf("Request for path: %v\n", path)
	location := filepath.Join(constants.Base, path)

	if fsinfo, err := os.Stat(location); os.IsNotExist(err) {
		rw.WriteHeader(http.StatusNotFound)
		rw.Header().Set("Content-Type", "text/plain")
		if _, err := rw.Write(([]byte(fmt.Sprintf("%v does not exist\n", r.URL.Path)))); err != nil {
			log.Printf("error writing response message: %v\n", err)
		}
	} else if fsinfo.IsDir() {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Set("Content-Type", "text/plain")
		if _, err := rw.Write(([]byte(fmt.Sprintf("%v is a directory\n", r.URL.Path)))); err != nil {
			log.Printf("error writing response message: %v\n", err)
		}
	} else if file, err := os.Open(location); err != nil {
		log.Printf("error opening file '%v'\n", location)
		rw.WriteHeader(http.StatusInternalServerError)
		if _, err := rw.Write(([]byte{})); err != nil {
			log.Printf("error writing response message: %v\n", err)
		}
	} else {
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
	}
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
