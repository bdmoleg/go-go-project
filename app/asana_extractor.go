package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bdmoleg/go-go-project/configs"
)

var RateLimitedErr = errors.New("rate limited error")

const (
	AsanaApiProjectsPath = "projects"
	AsanaApiUsersPath    = "users"
)

type AsanaExtractor struct {
	client        *http.Client
	logger        *slog.Logger
	mu            sync.Mutex
	isRateLimited bool
	retryAfter    int
}

func NewAsanaExtractor(client *http.Client, logger *slog.Logger) *AsanaExtractor {
	return &AsanaExtractor{
		client: client,
		logger: logger,
	}
}

func (a *AsanaExtractor) Extract(ctx context.Context, extractInterval time.Duration) {
	fullProjectsEndpoint, err := url.JoinPath(configs.AsanaAPIEndpoint, AsanaApiProjectsPath)
	if err != nil {
		a.logger.Error("failed to construct full endpoint API for projects", slog.String("error", err.Error()))
		panic(err.Error())
	}

	fullUsersEndpoint, err := url.JoinPath(configs.AsanaAPIEndpoint, AsanaApiUsersPath)

	a.logger.Info("full enpoint projects path", slog.String("endpoint", fullProjectsEndpoint))
	a.logger.Info("full enpoint users path", slog.String("endpoint", fullUsersEndpoint))

	if err != nil {
		a.logger.Error("failed to construct full endpoint API for users", slog.String("error", err.Error()))
		panic(err.Error())
	}

	// there are 2 requests to Asana API made per tick
	ticker := time.NewTicker(extractInterval)
	for {

		select {
		case <-ticker.C:
			asanaProjects, err := a.GetAsanaProjects(fullProjectsEndpoint)
			if err != nil {
				return
			}

			asanaUsers, err := a.GetAsanaUsers(fullUsersEndpoint)
			if err != nil {
				return
			}

			a.logger.Info("projects extracted", "projects", asanaProjects)
			a.logger.Info("users extracted", "users", asanaUsers)
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				a.logger.Error("context was not finished by cancel", slog.String("error", ctx.Err().Error()))
			}
			a.logger.Info("Asana Extractor finished execution")
			return
		}
	}
}

func (a *AsanaExtractor) GetAsanaProjects(endpoint string) ([]AsanaProject, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.isRateLimited {
		return nil, RateLimitedErr
	}

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		a.logger.Error("failed to construct HTTP request", slog.String("error", err.Error()))
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", configs.AsanaToken))

	response, err := a.client.Do(request)
	if err != nil {
		a.logger.Error("failed to get response from the endpoint", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		errClose := response.Body.Close()
		if errClose != nil {
			a.logger.Error("failed to close response body", slog.String("error", errClose.Error()))
		}
	}()

	responsePayload, err := io.ReadAll(response.Body)
	if err != nil {
		a.logger.Error("failed to read response body", slog.String("error", err.Error()))
		return nil, err
	}

	asanaProjects := &AsanaProjectsResponse{}
	err = json.Unmarshal(responsePayload, asanaProjects)
	if err != nil {
		a.logger.Error("failed to unmarshal json", slog.String("error", err.Error()))
		return nil, err
	}

	return asanaProjects.AsanaProjects, nil
}

func (a *AsanaExtractor) GetAsanaUsers(endpoint string) ([]AsanaUser, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.isRateLimited {
		return nil, RateLimitedErr
	}

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		a.logger.Error("failed to construct HTTP request", slog.String("error", err.Error()))
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", configs.AsanaToken))

	response, err := a.client.Do(request)
	if err != nil {
		a.logger.Error("failed to get response from the endpoint", slog.String("error", err.Error()))
		return nil, err
	}
	defer func() {
		errClose := response.Body.Close()
		if errClose != nil {
			a.logger.Error("failed to close response body", slog.String("error", errClose.Error()))
		}
	}()
	if response.StatusCode == 429 {
		a.isRateLimited = true
		a.logger.Info("got a rate limiting")
		return nil, RateLimitedErr
	}

	retryAfterHeader := response.Header.Get("Retry-After")
	a.logger.Info("rate-limiting data", "retry-after", retryAfterHeader, "statusCode", response.StatusCode)

	responsePayload, err := io.ReadAll(response.Body)
	if err != nil {
		a.logger.Error("failed to read response body", slog.String("error", err.Error()))
		return nil, err
	}

	asanaProjects := &AsanaUsersResponse{}
	err = json.Unmarshal(responsePayload, asanaProjects)
	if err != nil {
		a.logger.Error("failed to unmarshal json", slog.String("error", err.Error()))
		return nil, err
	}

	return asanaProjects.AsanaUsers, nil
}

func (a *AsanaExtractor) ThreadsTestRateLimiting(wg *sync.WaitGroup, numOf int) {
	fullProjectsEndpoint, _ := url.JoinPath(configs.AsanaAPIEndpoint, AsanaApiProjectsPath)
	for i := 0; i < numOf; i++ {
		go func() {
			defer wg.Done()
			a.GetAsanaUsers(fullProjectsEndpoint)
		}()
	}
}

// RateLimiterWatcher update rate limiting boolean
func (a *AsanaExtractor) RateLimiterWatcher() {

}
