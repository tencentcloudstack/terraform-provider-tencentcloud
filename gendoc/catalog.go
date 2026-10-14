package main

// catalog.go — API catalog loading from the bundled snapshots.
//
// Doc generation is fully offline: 'make doc' reads only the embedded
// snapshots. Snapshot freshness is maintained explicitly via
// 'make doc-refresh', which re-downloads the upstream catalogs.

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
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
// bundled snapshots, which are refreshed explicitly via 'make doc-refresh'.
func loadCatalog() *catalog {
	idx, err := loadSnapshotCatalog()
	if err != nil {
		message("[WARN]bundled snapshot unavailable (%v); API Reference sections will be empty\n", err)
		return &catalog{actionIdx: map[string]apiAction{}, source: "none"}
	}
	return &catalog{actionIdx: idx, source: "snapshot"}
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
