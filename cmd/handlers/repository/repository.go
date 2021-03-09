package repository

import "net/http"

const defaultPath = "F:/data/.m2/repository"

type Repository struct {
	Path     string
	Webhooks *Webhooks
}

func (r *Repository) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.download(rw, req)
	case http.MethodPut:
		r.upload(rw, req)
	default:
		rw.Header().Set("Allow", "GET, PUT")
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *Repository) Base() string {
	if r.Path == "" {
		return defaultPath
	}
	return r.Path
}
