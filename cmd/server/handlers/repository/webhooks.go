package repository

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Webhook struct {
	Url string
}

type Webhooks struct {
	Transport  http.RoundTripper
	Registered []Webhook
}

type event struct {
	Domain    string      `json:"domain,omitempty"`
	EventType string      `json:"event_type,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type artifactEvent struct {
	Name   string `json:"name,omitempty"`
	Path   string `json:"path,omitempty"`
	Repo   string `json:"repo_key,omitempty"`
	Sha256 string `json:"sha256,omitempty"`
	Size   uint   `json:"size,omitempty"`

	SourceRepoPath string `json:"source_repo_path,omitempty"`
	TargetRepoPath string `json:"target_repo_path,omitempty"`

	PropertyKey    string   `json:"property_key,omitempty"`
	PropertyValues []string `json:"property_values,omitempty"`
}

/*

{
  "repo_key":"sample_repo",
  "path":"sample_path_dir/sample_artifact",
  "name":"sample_artifact",
  "sha256":"ec1be623d148ed220f70f4f6125dc738b1d301a85b75e87c5b554fa3bb1b4141",
  "size":17848
}

*/

func (w *Webhooks) ArtifactCreated(repo, path, name, digest string, size uint) {
	w.publish(event{
		Domain:    "artifact",
		EventType: "deployed",
		Data: artifactEvent{
			Repo:   repo,
			Path:   path,
			Name:   name,
			Sha256: digest,
			Size:   size,
		},
	})
}

func (w *Webhooks) publish(a interface{}) {
	buf, err := json.Marshal(a)
	if err != nil {
		log.Printf("Error encoding event: %v\n", err)
		return
	}
	for _, r := range w.Registered {
		req, err := http.NewRequest(http.MethodPost, r.Url, bytes.NewReader(buf))
		if err != nil {
			log.Printf("Failed to create new http request: %v\n", err)
			continue
		}

		resp, err := w.transport().RoundTrip(req)
		if err != nil {
			log.Printf("URL: %s", req.URL)
			log.Printf("Failed to submit http request: %v", err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Printf("Unexpected http status %v in response: %v", resp.StatusCode, resp.Status)
		}
	}
}

func (w *Webhooks) transport() http.RoundTripper {
	if w.Transport != nil {
		return w.Transport
	}
	return http.DefaultTransport
}
