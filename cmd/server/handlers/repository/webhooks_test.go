package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStuff(t *testing.T) {
	events := make([]*event, 0)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		require.Equal(t, r.Method, http.MethodPost)
		require.Equal(t, r.URL.Path, "/webhooks/artifactory")
		event := new(event)
		if err := json.NewDecoder(r.Body).Decode(event); err != nil {
			t.Fatalf("Error decoding incoming event: %v", err)
		}
		events = append(events, event)
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	t.Logf("Test server: %v", server.URL)

	webhooks := &Webhooks{
		Transport:  server.Client().Transport,
		Registered: []Webhook{{Url: fmt.Sprintf("%s/webhooks/artifactory", server.URL)}},
	}

	webhooks.ArtifactCreated("libs-release-local", "de/example/example-app/1.1", "example-app-1.1.jar", "TODO", 10)

	if len(events) != 1 {
		t.Fatalf("Expected exactly one event to be published, but it actually was %v", len(events))
	}

	event := events[0]
	require.Equal(t, event.Domain, "artifact", "Expecting Event.Domain to be 'artifiact', but it is '%v'", event.Domain)
	require.Equal(t, event.EventType, "deployed", "Expecting Event.EventType to be 'deployed', but it is '%v'", event.EventType)
	log.Printf("Type of Event.Data: %T", event.Data)
	data := event.Data.(map[string]interface{})
	require.Equal(t, data["name"], "example-app-1.1.jar")
}

type mockClient struct {
}

func (c *mockClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK}, nil
}
