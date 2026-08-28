package upstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Extension is a Marketplace extension summary.
type Extension struct {
	Publisher    string `json:"publisher"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	LatestVersion string `json:"latestVersion"`
}

// ExtensionVersion is one version of an extension.
type ExtensionVersion struct {
	Version       string    `json:"version"`
	PublishTime   time.Time `json:"publishTime"`
	EnginesVSCode string    `json:"enginesVscode"`
	DownloadURL   string    `json:"downloadUrl"`
}

type extensionClient struct {
	hc *http.Client
}

func newExtensionClient() *extensionClient {
	return &extensionClient{hc: &http.Client{Timeout: 10 * time.Second}}
}

// SearchExtensions hits the Marketplace gallery search API.
// Sorting is left to the upstream default (no flags sent).
func SearchExtensions(ctx context.Context, query string, page, pageSize int) ([]Extension, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}
	u := fmt.Sprintf("https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery?page=%d&pageSize=%d", page, pageSize)
	body := map[string]any{
		"filters": []map[string]any{
			{"criteria": []map[string]any{
				{"filterType": 10, "value": query},
			}, "pageNumber": page, "pageSize": pageSize},
		},
		"assetTypes": []string{},
		"flags":      0x1 | 0x4, // IncludeLatestVersionOnly + IncludeFiles
	}
	c := newExtensionClient()
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, u, jsonReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json;api-version=7.1-preview.1")
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, 0, fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	var out struct {
		Results []struct {
			Extensions []struct {
				Publisher       struct{ PublisherName string `json:"publisherName"` } `json:"publisher"`
				ExtensionName   string `json:"extensionName"`
				DisplayName     string `json:"displayName"`
				Versions        []struct {
					Version string `json:"version"`
				} `json:"versions"`
			} `json:"extensions"`
			ResultMetadata []struct {
				MetadataType string `json:"metadataType"`
				TotalCount   int    `json:"totalCount"`
			} `json:"resultMetadata"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, 0, err
	}
	var extensions []Extension
	if len(out.Results) > 0 {
		for _, e := range out.Results[0].Extensions {
			latest := ""
			if len(e.Versions) > 0 {
				latest = e.Versions[0].Version
			}
			extensions = append(extensions, Extension{
				Publisher:     strings.ToLower(e.Publisher.PublisherName),
				Name:          e.ExtensionName,
				DisplayName:   e.DisplayName,
				LatestVersion: latest,
			})
		}
	}
	total := 0
	for _, r := range out.Results {
		for _, m := range r.ResultMetadata {
			if m.MetadataType == "ResultCount" {
				total = m.TotalCount
			}
		}
	}
	return extensions, total, nil
}

// FetchExtensionVersions returns all versions for a single extension.
func FetchExtensionVersions(ctx context.Context, publisher, name string) ([]ExtensionVersion, error) {
	u := fmt.Sprintf("https://marketplace.visualstudio.com/_apis/public/gallery/publishers/%s/vsextensions/%s",
		url.PathEscape(publisher), url.PathEscape(name))
	c := newExtensionClient()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	req.Header.Set("Accept", "application/json;api-version=7.1-preview.1")
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, ErrExtensionNotFound
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("upstream status %d", resp.StatusCode)
	}
	var payload struct {
		Versions []struct {
			Version     string    `json:"version"`
			LastUpdated time.Time `json:"lastUpdated"`
			Engines     []struct {
				Runtime string `json:"runtime"`
			} `json:"engines"`
			Files []struct {
				AssetType string `json:"assetType"`
				Source    string `json:"source"`
			} `json:"files"`
		} `json:"versions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	out := make([]ExtensionVersion, 0, len(payload.Versions))
	for _, v := range payload.Versions {
		engines := ""
		for _, e := range v.Engines {
			if strings.HasPrefix(e.Runtime, "vscode") {
				engines = e.Runtime
				break
			}
		}
		download := ""
		for _, f := range v.Files {
			if f.AssetType == "Microsoft.VisualStudio.Services.VSIXPackage" {
				download = f.Source
				break
			}
		}
		out = append(out, ExtensionVersion{
			Version:       v.Version,
			PublishTime:   v.LastUpdated,
			EnginesVSCode: engines,
			DownloadURL:   download,
		})
	}
	return out, nil
}

// jsonReader is a tiny helper that returns an io.Reader from any JSON-encodable value.
type jsonBodyReader struct {
	v any
}

func (j jsonBodyReader) Read(p []byte) (int, error) {
	if j.v == nil {
		return 0, io.EOF
	}
	data, err := json.Marshal(j.v)
	if err != nil {
		return 0, err
	}
	j.v = nil
	n := copy(p, data)
	if n < len(data) {
		return n, nil
	}
	return n, io.EOF
}

func jsonReader(v any) *jsonBodyReader { return &jsonBodyReader{v: v} }

// ErrExtensionNotFound signals a missing extension.
var ErrExtensionNotFound = fmt.Errorf("extension not found")