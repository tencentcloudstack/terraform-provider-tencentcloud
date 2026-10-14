package main

// apiref_en.go — English action descriptions sourced from the official
// international SDK (TencentCloud/tencentcloud-sdk-go-intl-en).
//
// The snapshot gendoc/data/api-actions-en.json maps
// (product, action) -> first prose line of the method's doc comment plus
// the official doc URL from the public catalog. Refresh both snapshots
// with 'make doc-refresh'. Doc rendering prefers the English description;
// entries without one degrade to the bare action name so that no Chinese
// ever reaches the generated docs.

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

// htmlTagRe matches inline HTML tags that leak from the SDK docs.
var htmlTagRe = regexp.MustCompile(`</?[a-zA-Z][^>]*>`)

// maxDescLen bounds a description to keep the generated doc readable.
const maxDescLen = 240

// cleanDesc strips inline HTML tags and collapses whitespace. Descriptions
// that grow beyond maxDescLen are cut at the last sentence end before the
// cap (falling back to a hard cut).
func cleanDesc(s string) string {
	if strings.Contains(s, "<") {
		s = htmlTagRe.ReplaceAllString(s, "")
	}
	s = strings.TrimSpace(s)
	if len(s) <= maxDescLen {
		return s
	}
	cut := s[:maxDescLen]
	if i := strings.LastIndexAny(cut, ".!?"); i > maxDescLen/2 {
		return strings.TrimSpace(cut[:i+1])
	}
	if i := strings.LastIndex(cut, " "); i > 0 {
		return strings.TrimSpace(cut[:i]) + " ..."
	}
	return cut + " ..."
}

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
		desc = cleanDesc(desc)
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
// action descriptions, merges doc URLs — preferring the international site,
// falling back to the Chinese site — and rewrites
// gendoc/data/api-actions-en.json atomically (the old snapshot survives any
// failure).
func refreshENData(rawActions []rawAction) error {
	// Prefer international URLs when the intl doc site lists the action;
	// fall back to the Chinese-site URL otherwise.
	urls := map[string]string{}
	for _, a := range rawActions {
		if a.ActionCNDocUr != "" {
			urls[strings.ToLower(a.Product)+"."+a.Action] = a.ActionCNDocUr
		}
	}
	intlURLs, err := refreshIntlURLs(rawActions)
	if err != nil {
		message("[WARN]international doc URLs unavailable (%v); using Chinese-site URLs only\n", err)
	} else {
		replaced := 0
		for key, u := range intlURLs {
			if _, cn := urls[key]; cn {
				replaced++
			}
			urls[key] = u
		}
		message("[SUCC.]international doc URLs: %d entries (%d replaced Chinese URLs)\n", len(intlURLs), replaced)
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
	outIdx := map[string]int{} // "<product>.<action>" -> index into out
	var out []enEntry
	for _, a := range rawActions {
		key := strings.ToLower(a.Product) + "." + a.Action
		if a.Action == "" {
			continue
		}
		if _, dup := outIdx[key]; dup {
			continue
		}
		outIdx[key] = len(out)
		out = append(out, enEntry{Product: a.Product, Action: a.Action, URL: urls[key]})
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
			if desc == "" {
				continue
			}
			key := product + "." + action
			if i, ok := outIdx[key]; ok {
				out[i].Desc = desc
				continue
			}
			// Not in the catalog: keep it anyway (catalog may lag the SDK).
			outIdx[key] = len(out)
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
