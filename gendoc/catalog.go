package main

// catalog.go — API catalog loading with online-first, snapshot-fallback
// semantics.
//
// make doc fetches the live public catalog on every run. When the download
// fails or times out, the bundled snapshot (go:embed) is used instead so
// that doc generation never hard-fails on network availability; the stale
// sections are simply rendered from the snapshot and a warning is printed.

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//go:embed data/api-actions.json
var apiActionsFS embed.FS

// catalog holds the resolved action index shared by all doc renders.
type catalog struct {
	actionIdx map[string]apiAction
	source    string // "live" or "snapshot"
}

// loadCatalog returns the embedded snapshot catalog. Doc generation is
// fully offline: action URLs and English descriptions both come from the
// bundled snapshots, which are refreshed explicitly via
// 'make doc-refresh-data' / 'make doc-refresh-en'.
func loadCatalog(*http.Client) *catalog {
	idx, err := loadSnapshotCatalog()
	if err != nil {
		message("[WARN]bundled snapshot unavailable (%v); API Reference sections will be empty\n", err)
		return &catalog{actionIdx: map[string]apiAction{}, source: "none"}
	}
	return &catalog{actionIdx: idx, source: "snapshot"}
}

// fetchActionCatalog downloads the live catalog and builds the lookup index.
func fetchActionCatalog(httpClient *http.Client) (map[string]apiAction, error) {
	actions, err := fetchAndSlimActionsWith(httpClient)
	if err != nil {
		return nil, err
	}
	if len(actions) == 0 {
		return nil, fmt.Errorf("live catalog is empty")
	}
	return buildActionIndex(actions), nil
}

// loadSnapshotCatalogOnly returns a snapshot-backed catalog, never touching
// the network. Used by the offline probe mode.
func loadSnapshotCatalogOnly() *catalog {
	idx, err := loadSnapshotCatalog()
	if err != nil {
		return &catalog{actionIdx: map[string]apiAction{}, source: "none"}
	}
	return &catalog{actionIdx: idx, source: "snapshot"}
}

// loadSnapshotCatalog builds the lookup index from the embedded snapshot.
func loadSnapshotCatalog() (map[string]apiAction, error) {
	raw, err := apiActionsFS.ReadFile("data/api-actions.json")
	if err != nil {
		return nil, fmt.Errorf("read embedded api-actions.json: %w", err)
	}
	var actions []apiAction
	if err := json.Unmarshal(raw, &actions); err != nil {
		return nil, fmt.Errorf("parse api-actions.json: %w", err)
	}
	return buildActionIndex(actions), nil
}

func buildActionIndex(actions []apiAction) map[string]apiAction {
	idx := make(map[string]apiAction, len(actions))
	for _, a := range actions {
		idx[strings.ToLower(a.Product)+"."+a.Action] = a
	}
	return idx
}

// liveHTTPClient is the client used for the per-run live fetch. The timeout
// is deliberately short: doc generation should not hang on a slow network,
// and the snapshot fallback covers the failure case.
func liveHTTPClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}
