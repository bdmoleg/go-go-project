package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bdmoleg/go-go-project/configs"
	"github.com/stretchr/testify/assert"
)

var discardLogger = slog.New(slog.NewJSONHandler(io.Discard, nil))

func TestGetAsanaProjects(t *testing.T) {
	fmt.Println("here is fmt.Println(`fmt msg`)")
	configs.AsanaToken = "123456"
	t.Cleanup(func() {
		configs.AsanaToken = ""
	})

	httpTestServer := createTestMockedHttpServer()
	defer httpTestServer.Close()

	fullPath, err := url.JoinPath(httpTestServer.URL, "projects")
	assert.NoError(t, err)
	asanaExtractor := NewAsanaExtractor(httpTestServer.Client(), discardLogger)
	projects, err := asanaExtractor.GetAsanaProjects(fullPath)

	assert.NoError(t, err)
	assert.Len(t, projects, 4)
	assert.Equal(t, 1, projects[0].Gid)
	assert.Equal(t, "projname1", projects[0].Name)
	assert.Equal(t, "resource-type-1", projects[0].Gid)
}

func TestGetAsanaUsers(t *testing.T) {
	configs.AsanaToken = "123456"
	t.Cleanup(func() {
		configs.AsanaToken = ""
	})
	httpTestServer := createTestMockedHttpServer()
	defer httpTestServer.Close()

	// try projects endpoint

}

func TestAbc(t *testing.T) {

}

func createTestMockedHttpServer() *httptest.Server {
	httpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectsData := AsanaProjectsResponse{
			AsanaProjects: []AsanaProject{
				{Gid: "1", Name: "projname1", ResourceType: "resource-type-1"},
				{Gid: "2", Name: "projname2", ResourceType: "resource-type-2"},
				{Gid: "3", Name: "projname3", ResourceType: "resource-type-3"},
				{Gid: "4", Name: "projname4", ResourceType: "resource-type-4"},
			},
		}
		usersData := AsanaUsersResponse{
			AsanaUsers: []AsanaUser{
				{Gid: "1", Name: "username1", ResourceType: "resource-type-1"},
				{Gid: "100", Name: "username2", ResourceType: "resource-type-2"},
				{Gid: "200", Name: "username3", ResourceType: "resource-type-3"},
			},
		}

		if r.URL.Path == "/projects" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(projectsData)
			return
		} else if r.URL.Path == "/users" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(usersData)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "error"}`))
	}))

	return httpServer
}
