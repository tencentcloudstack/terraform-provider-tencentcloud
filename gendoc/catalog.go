package main

// catalog.go — API catalog loading.
//
// 'make doc' fetches TencentCloud's public API catalog live on every run
// and rewrites the bundled snapshot, so the generated doc links always
// track the current cloud docs. When a fetch fails (e.g. offline), the
// embedded snapshot (gendoc/data/api-actions.json) is used as a fallback.
//
// Doc URLs are international-site only: the action catalog provides the
// (product, action) whitelist, and the rendered URL comes from the
// international docs catalogue (apiref_intl.go) keyed by the catalog
// entry's Chinese-site URL (its path carries the shared categoryId).
// Actions without an international page fall back to the international
// site's search URL — Chinese-site URLs are never rendered.

import (
	"embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

//go:embed data/api-actions.json
var apiActionsFS embed.FS

// cnDocCatRe extracts the catalog id from a Chinese-site doc URL. The
// international site shares the same product category ids, so the id
// keys the international URL lookup.
var cnDocCatRe = regexp.MustCompile(`cloud\.tencent\.com/document/api/(\d+)/`)

// catalog holds the resolved action index shared by all doc renders.
type catalog struct {
	actionIdx map[string]apiAction
}

// loadCatalog returns the API catalog: live data when the public endpoint
// is reachable (the snapshot is refreshed in the same pass), otherwise the
// embedded snapshot. International doc URLs are resolved in the same pass
// and stored on the catalog entries.
func loadCatalog() *catalog {
	actions, err := refreshCatalog()
	if err != nil {
		message("[WARN]live API catalog unavailable (%v); falling back to the embedded snapshot\n", err)
		actions, err = loadSnapshotCatalog()
		if err != nil {
			message("[WARN]bundled snapshot unavailable (%v); API Reference sections will be empty\n", err)
			return &catalog{actionIdx: map[string]apiAction{}}
		}
		return newCatalog(actions)
	}
	message("[SUCC.]live API catalog: %d actions (snapshot refreshed)\n", len(actions))
	return newCatalog(actions)
}

// newCatalog builds the action index and attaches international doc URLs.
// When the international catalogue cannot be fetched, every entry's URL
// stays empty and the renderer falls back to the international site's
// search URL per action — generation never fails on this.
func newCatalog(actions []apiAction) *catalog {
	intlURLs, err := fetchIntlActionURLs()
	if err != nil {
		message("[WARN]international doc catalogue unavailable (%v); actions will link to the search page\n", err)
	} else {
		matched := 0
		for i := range actions {
			a := &actions[i]
			if a.URL == "" {
				continue
			}
			m := cnDocCatRe.FindStringSubmatch(a.URL)
			if m == nil {
				continue
			}
			if u, ok := intlURLs[m[1]+"."+a.Action]; ok {
				a.URL = u
				matched++
			} else {
				a.URL = "" // no international page: fall back to search
			}
		}
		message("[SUCC.]international doc URLs resolved for %d/%d catalog actions\n", matched, len(actions))
	}
	return &catalog{actionIdx: buildActionIndex(actions)}
}

// loadSnapshotCatalogOnly returns a snapshot-backed catalog, never touching
// the network. Used by the offline probe mode.
func loadSnapshotCatalogOnly() *catalog {
	actions, err := loadSnapshotCatalog()
	if err != nil {
		return &catalog{actionIdx: map[string]apiAction{}}
	}
	return newCatalog(actions)
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
