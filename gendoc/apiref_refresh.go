package main

// apiref_refresh.go — refreshes the bundled API catalog snapshots from
// TencentCloud's public endpoints. Run with:
//
//	go run ./gendoc -refresh-data
//
// or via `make doc-refresh-data`. The download is ~5 MB; only four fields
// (product/action/name/url) are kept so the snapshot stays small.

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Public catalog endpoints. The raw JSON is large (~6 MB) and carries many
// extra fields; the refresh keeps only what the doc renderer needs.
const (
	actionsURL = "https://main.qcloudimg.com/qcloud/main/scripts/release/common/api-all-actions.json"
)

// rawAction mirrors one entry of api-all-actions.json.
type rawAction struct {
	Product       string `json:"product"`
	Action        string `json:"action"`
	ActionName    string `json:"actionName"`
	ActionCNDocUr string `json:"actionCNDocUrl"`
}

// refreshData downloads both catalogs and rewrites the slim snapshots under
// gendoc/data/.
// refreshData downloads and rewrites the slim catalog snapshots. The raw
// action catalog is also returned so the English refresh can merge doc
// URLs without a second download.
func refreshData() ([]rawAction, error) {
	actions, err := fetchAndSlimActions()
	if err != nil {
		return nil, fmt.Errorf("actions: %w", err)
	}
	if err := writeJSON("data/api-actions.json", actions); err != nil {
		return nil, err
	}
	message("[SUCC.]refreshed data/api-actions.json (%d actions)\n", len(actions))

	var raw []rawAction
	if err := fetchJSON(nil, actionsURL, &raw); err != nil {
		// Non-fatal: the English refresh tolerates a missing catalog.
		message("[WARN]raw catalog re-fetch for URL merge failed: %v\n", err)
		return nil, nil
	}
	return raw, nil
}

func fetchAndSlimActions() ([]apiAction, error) {
	return fetchAndSlimActionsWith(http.DefaultClient)
}

// fetchAndSlimActionsWith downloads and slims the action catalog using the
// given HTTP client.
func fetchAndSlimActionsWith(httpClient *http.Client) ([]apiAction, error) {
	var raw []rawAction
	if err := fetchJSON(httpClient, actionsURL, &raw); err != nil {
		return nil, err
	}
	out := make([]apiAction, 0, len(raw))
	for _, a := range raw {
		if a.Product == "" || a.Action == "" {
			continue
		}
		out = append(out, apiAction{
			Product: a.Product,
			Action:  a.Action,
			Name:    a.ActionName,
			URL:     a.ActionCNDocUr,
		})
	}
	return out, nil
}

// fetchJSON downloads url (following redirects, transparent gzip) and
// unmarshals it into v. The response is size-capped for safety.
func fetchJSON(httpClient *http.Client, url string, v interface{}) error {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 2 * time.Minute}
	}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download %s: %w", url, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: status %s", url, resp.Status)
	}

	var reader io.Reader = resp.Body
	if strings.EqualFold(resp.Header.Get("Content-Encoding"), "gzip") {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return fmt.Errorf("gzip %s: %w", url, err)
		}
		defer func() { _ = gz.Close() }()
		reader = gz
	}

	// 64 MiB cap: the raw catalog is ~6 MB today.
	raw, err := io.ReadAll(io.LimitReader(reader, 64<<20))
	if err != nil {
		return fmt.Errorf("read %s: %w", url, err)
	}
	if err := json.Unmarshal(raw, v); err != nil {
		return fmt.Errorf("parse %s: %w", url, err)
	}
	return nil
}

// writeJSON marshals v compactly and writes it to
// path relative to the gendoc directory, preserving the old file on failure.
func writeJSON(path string, v interface{}) error {
	raw, err := json.Marshal(v)
	if err != nil {
		return err
	}
	final := filepath.Join(".", path)
	tmp := final + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, final)
}
