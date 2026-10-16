package main

// apiref_refresh.go — live refresh of the API catalog used by the doc
// generation pipeline. 'make doc' fetches TencentCloud's public API catalog
// on every run so the generated links always match the current state of the
// cloud docs; the embedded snapshot (gendoc/data/api-actions.json) is a
// fallback for offline runs and is rewritten whenever a fetch succeeds.

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

// Public catalog endpoint. The raw JSON is large (~6 MB) and carries many
// extra fields; only a few fields (product/action/name/url) are kept so
// the snapshot stays small.
const actionsURL = "https://main.qcloudimg.com/qcloud/main/scripts/release/common/api-all-actions.json"

// rawAction mirrors one entry of api-all-actions.json.
type rawAction struct {
	Product       string `json:"product"`
	Action        string `json:"action"`
	ActionName    string `json:"actionName"`
	ActionCNDocUr string `json:"actionCNDocUrl"`
}

// refreshCatalog fetches the live API catalog and rewrites the embedded
// snapshot. Returns the slim catalog for the current run. On failure the
// error is returned and callers fall back to the bundled snapshot.
func refreshCatalog() ([]apiAction, error) {
	actions, err := fetchAndSlimActions(http.DefaultClient)
	if err != nil {
		return nil, err
	}
	if err := writeJSON("data/api-actions.json", actions); err != nil {
		return nil, err
	}
	return actions, nil
}

// fetchAndSlimActions downloads and slims the action catalog using the
// given HTTP client.
func fetchAndSlimActions(httpClient *http.Client) ([]apiAction, error) {
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
