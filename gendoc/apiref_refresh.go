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
	actionsURL  = "https://main.qcloudimg.com/qcloud/main/scripts/release/common/api-all-actions.json"
	productsURL = "https://main.qcloudimg.com/qcloud/main/scripts/release/common/api-all-products.json"
)

// rawAction mirrors one entry of api-all-actions.json.
type rawAction struct {
	Product       string `json:"product"`
	Action        string `json:"action"`
	ActionName    string `json:"actionName"`
	ActionCNDocUr string `json:"actionCNDocUrl"`
}

// rawProduct mirrors one entry of api-all-products.json.
type rawProduct struct {
	Product             string `json:"product"`
	ProductName         string `json:"productName"`
	ProductCNDocIntroUr string `json:"productCNDocIntroUrl"`
}

// refreshData downloads both catalogs and rewrites the slim snapshots under
// gendoc/data/.
func refreshData() error {
	actions, err := fetchAndSlimActions()
	if err != nil {
		return fmt.Errorf("actions: %w", err)
	}
	products, err := fetchAndSlimProducts()
	if err != nil {
		return fmt.Errorf("products: %w", err)
	}

	if err := writeJSON("data/api-actions.json", actions); err != nil {
		return err
	}
	if err := writeJSON("data/api-products.json", products); err != nil {
		return err
	}
	message("[SUCC.]refreshed data/api-actions.json (%d actions) and data/api-products.json (%d products)\n",
		len(actions), len(products))
	return nil
}

func fetchAndSlimActions() ([]apiAction, error) {
	var raw []rawAction
	if err := fetchJSON(actionsURL, &raw); err != nil {
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

func fetchAndSlimProducts() ([]map[string]string, error) {
	var raw []rawProduct
	if err := fetchJSON(productsURL, &raw); err != nil {
		return nil, err
	}
	out := make([]map[string]string, 0, len(raw))
	for _, p := range raw {
		if p.Product == "" {
			continue
		}
		out = append(out, map[string]string{
			"product": p.Product,
			"name":    p.ProductName,
			"url":     p.ProductCNDocIntroUr,
		})
	}
	return out, nil
}

// fetchJSON downloads url (following redirects, transparent gzip) and
// unmarshals it into v. The response is size-capped for safety.
func fetchJSON(url string, v interface{}) error {
	client := &http.Client{Timeout: 2 * time.Minute}
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
