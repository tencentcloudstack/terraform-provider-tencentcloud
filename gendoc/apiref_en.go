package main

// apiref_en.go — English action descriptions sourced from the official
// international SDK (TencentCloud/tencentcloud-sdk-go-intl-en).
//
// The snapshot gendoc/data/api-actions-en.json maps
// (product, action) -> first prose line of the method's doc comment.
// Refresh it with:
//
//	make doc-refresh-en
//
// which downloads the SDK master tarball (~6 MB), re-extracts the
// descriptions and atomically rewrites the snapshot. Doc rendering prefers
// the English description; entries without one degrade to the bare action
// name so that no Chinese ever reaches the generated docs.

import (
	"archive/tar"
	"compress/gzip"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

//go:embed data/api-actions-en.json
var apiActionsEnFS embed.FS

const intlSDKTarballURL = "https://github.com/TencentCloud/tencentcloud-sdk-go-intl-en/archive/refs/heads/master.tar.gz"

// enEntry is one record of gendoc/data/api-actions-en.json.
type enEntry struct {
	Product string `json:"product"`
	Action  string `json:"action"`
	Desc    string `json:"desc"`
	URL     string `json:"url,omitempty"`
}

// blockRe matches a doc-comment block starting with the action-name line:
//
//	// ActionName
//	// ...comment lines...
var blockRe = regexp.MustCompile(`(?m)^// ([A-Z][A-Za-z0-9]*)$\n(//[^\n]*\n)+`)

// ecRe matches SDK error-code comment lines: `NAME = "Namespace.Name"`.
var ecRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*\s*=\s*"`)

// isProseDesc reports whether s reads like a prose description line
// (uppercase-initial sentence) rather than an error-code entry.
func isProseDesc(s string) bool {
	if len(s) < 15 || s[0] < 'A' || s[0] > 'Z' {
		return false
	}
	if strings.Contains(s, "= \"") {
		return false
	}
	return len(strings.Fields(s)) >= 4
}

// loadENIndex loads the embedded English descriptions into a lookup index
// keyed by "<lowercase product>.<action>".
func loadENIndex() map[string]string {
	raw, err := apiActionsEnFS.ReadFile("data/api-actions-en.json")
	if err != nil {
		return map[string]string{}
	}
	var entries []enEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return map[string]string{}
	}
	idx := make(map[string]string, len(entries))
	for _, e := range entries {
		idx[strings.ToLower(e.Product)+"."+e.Action] = e.Desc
	}
	return idx
}

// loadENURLIndex loads the embedded doc URLs keyed the same way. Empty
// entries are omitted so the renderer falls back to the URL-less form.
func loadENURLIndex() map[string]string {
	raw, err := apiActionsEnFS.ReadFile("data/api-actions-en.json")
	if err != nil {
		return map[string]string{}
	}
	var entries []enEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return map[string]string{}
	}
	idx := make(map[string]string, len(entries))
	for _, e := range entries {
		if e.URL != "" {
			idx[strings.ToLower(e.Product)+"."+e.Action] = e.URL
		}
	}
	return idx
}

// extractENFromClient extracts (action, desc) pairs from one intl SDK
// client.go. Each API method carries a doc comment block: the first line is
// the action name, the following lines hold the description, then the
// error-code list. Both the plain and the WithContext method repeat the
// same comment; the first occurrence wins.
func extractENFromClient(src string) map[string]string {
	out := map[string]string{}
	locs := blockRe.FindAllStringSubmatchIndex(src, -1)
	for _, loc := range locs {
		action := src[loc[2]:loc[3]]
		// The description is the comment line right after the action-name
		// line (loc[0] is the block start).
		line1End := strings.Index(src[loc[0]:], "\n")
		rest := src[loc[0]+line1End+1 : loc[1]]
		lines := strings.Split(strings.TrimSuffix(rest, "\n"), "\n")
		desc := ""
		for _, ln := range lines {
			t := strings.TrimSpace(strings.TrimPrefix(ln, "//"))
			if t == "" {
				continue
			}
			if strings.HasPrefix(t, "error code that may be returned") {
				break
			}
			if ecRe.MatchString(t) || !isProseDesc(t) {
				continue
			}
			desc = t
			break
		}
		if desc == "" {
			continue
		}
		// Anchor to the method declaration that follows this comment.
		// The error-code list can be long, so scan a generous window.
		tail := src[loc[1]:minEN(len(src), loc[1]+16384)]
		end := strings.Index(tail, "\nfunc ")
		if end == -1 || end > 16384 {
			continue
		}
		if !strings.Contains(tail[:end], "func (c *Client) "+action+"(") {
			continue
		}
		if _, dup := out[action]; !dup {
			out[action] = desc
		}
	}
	return out
}

func minEN(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// refreshENData downloads the intl SDK master tarball, extracts the English
// action descriptions, merges the official doc URLs from the public catalog
// and rewrites gendoc/data/api-actions-en.json atomically (the old snapshot
// survives any failure).
func refreshENData() error {
	// Doc URLs live in the public catalog; fetch it best-effort. Without
	// it the snapshot still carries descriptions, just without URLs.
	urls := map[string]string{}
	var rawActions []rawAction
	if err := fetchJSON(nil, actionsURL, &rawActions); err == nil {
		for _, a := range rawActions {
			if a.ActionCNDocUr != "" {
				urls[strings.ToLower(a.Product)+"."+a.Action] = a.ActionCNDocUr
			}
		}
	} else {
		message("[WARN]catalog fetch for doc URLs failed (%v); writing descriptions without URLs\n", err)
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(intlSDKTarballURL)
	if err != nil {
		return fmt.Errorf("download %s: %w", intlSDKTarballURL, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: status %s", intlSDKTarballURL, resp.Status)
	}

	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("gzip: %w", err)
	}
	tr := tar.NewReader(gz)

	// Base the snapshot on the full public catalog so that entries without
	// an English description still carry their doc URL; the intl SDK only
	// contributes descriptions.
	seen := map[string]bool{}
	var out []enEntry
	for _, a := range rawActions {
		key := strings.ToLower(a.Product) + "." + a.Action
		if seen[key] || a.Action == "" {
			continue
		}
		seen[key] = true
		out = append(out, enEntry{Product: a.Product, Action: a.Action, URL: a.ActionCNDocUr})
	}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar: %w", err)
		}
		if !strings.HasSuffix(filepath.Base(hdr.Name), "client.go") ||
			strings.Contains(hdr.Name, "/examples/") {
			continue
		}
		parts := strings.Split(hdr.Name, "/")
		pi := -1
		for i, p := range parts {
			if p == "tencentcloud" {
				pi = i
			}
		}
		if pi == -1 || pi+2 >= len(parts) {
			continue
		}
		product := strings.ToLower(parts[pi+1])
		src, err := io.ReadAll(io.LimitReader(tr, 32<<20))
		if err != nil {
			continue
		}
		for action, desc := range extractENFromClient(string(src)) {
			key := product + "." + action
			if seen[key] {
				// Catalog-known action: attach the extracted description.
				for i := range out {
					if strings.ToLower(out[i].Product)+"."+out[i].Action == key {
						out[i].Desc = desc
						break
					}
				}
				continue
			}
			seen[key] = true
			out = append(out, enEntry{Product: product, Action: action, Desc: desc})
		}
	}

	if len(out) == 0 {
		return fmt.Errorf("no catalog entries extracted; refusing to overwrite the snapshot")
	}

	var b strings.Builder
	b.WriteString("[")
	for i, e := range out {
		if i > 0 {
			b.WriteString(",")
		}
		if e.URL != "" {
			fmt.Fprintf(&b, `{"product":%q,"action":%q,"desc":%q,"url":%q}`, e.Product, e.Action, e.Desc, e.URL)
		} else {
			fmt.Fprintf(&b, `{"product":%q,"action":%q,"desc":%q}`, e.Product, e.Action, e.Desc)
		}
	}
	b.WriteString("]")

	const path = "data/api-actions-en.json"
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(b.String()), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		return err
	}
	withDesc := 0
	for _, e := range out {
		if e.Desc != "" {
			withDesc++
		}
	}
	message("[SUCC.]refreshed %s (%d entries, %d with descriptions)\n", path, len(out), withDesc)
	return nil
}
