package main

// catalog.go — API catalog loading.
//
// 'make doc' fetches TencentCloud's public API catalog live on every run
// and rewrites the bundled snapshot, so the generated doc links always
// track the current cloud docs. When the fetch fails (e.g. offline), the
// embedded snapshot (gendoc/data/api-actions.json) is used as a fallback.

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
}

// loadCatalog returns the API catalog: live data when the public endpoint
// is reachable (the snapshot is refreshed in the same pass), otherwise the
// embedded snapshot.
func loadCatalog() *catalog {
	actions, err := refreshCatalog()
	if err != nil {
		message("[WARN]live API catalog unavailable (%v); falling back to the embedded snapshot\n", err)
		actions, err = loadSnapshotCatalog()
		if err != nil {
			message("[WARN]bundled snapshot unavailable (%v); API Reference sections will be empty\n", err)
			return &catalog{actionIdx: map[string]apiAction{}}
		}
		return &catalog{actionIdx: buildActionIndex(actions)}
	}
	message("[SUCC.]live API catalog: %d actions (snapshot refreshed)\n", len(actions))
	return &catalog{actionIdx: buildActionIndex(actions)}
}

// loadSnapshotCatalogOnly returns a snapshot-backed catalog, never touching
// the network. Used by the offline probe mode.
func loadSnapshotCatalogOnly() *catalog {
	actions, err := loadSnapshotCatalog()
	if err != nil {
		return &catalog{actionIdx: map[string]apiAction{}}
	}
	return &catalog{actionIdx: buildActionIndex(actions)}
}

// loadSnapshotCatalog builds the lookup index from the embedded snapshot.
func loadSnapshotCatalog() ([]apiAction, error) {
	raw, err := apiActionsFS.ReadFile("data/api-actions.json")
	if err != nil {
		return nil, fmt.Errorf("read embedded api-actions.json: %w", err)
	}
	var actions []apiAction
	if err := json.Unmarshal(raw, &actions); err != nil {
		return nil, fmt.Errorf("parse api-actions.json: %w", err)
	}
	return actions, nil
}

func buildActionIndex(actions []apiAction) map[string]apiAction {
	idx := make(map[string]apiAction, len(actions))
	for _, a := range actions {
		idx[strings.ToLower(a.Product)+"."+a.Action] = a
	}
	return idx
}
