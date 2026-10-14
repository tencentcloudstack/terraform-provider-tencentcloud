package main

// apiref_intl.go — English-site doc URLs from the tencentcloud.com
// documentation site.
//
// The international site assigns its own document IDs, so Chinese-site URLs
// cannot be rewritten by pattern. The site exposes two internal endpoints
// (they require browser-like headers):
//
//	POST /document/api/mid/getCategoryByLang {"lang":"en"}
//	  -> top-level product catalog (categoryId, title)
//	POST /document/api/mid/getCatalogue     {"lang":"en","categoryId":N}
//	  -> full doc tree for one product; each leaf page title is the action
//	     name and its id forms /document/product/<categoryId>/<id>
//
// refreshIntlURLs() walks the products referenced by the raw catalog and
// builds (product, action) -> international URL, stored in the en snapshot
// during 'make doc-refresh'. Actions the international site has not
// translated fall back to the Chinese-site URL at snapshot-merge time.

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// cnDocCatRe extracts the catalog id from a Chinese-site doc URL.
var cnDocCatRe = regexp.MustCompile(`cloud\.tencent\.com/document/api/(\d+)/`)

const (
	intlBase   = "https://www.tencentcloud.com"
	intlCatURL = intlBase + "/document/api/mid/getCatalogue"
)

// intlNode is one node of the getCatalogue doc tree.
type intlNode struct {
	Title    string     `json:"title"`
	Type     string     `json:"type"`
	ID       int64      `json:"id"`
	Children []intlNode `json:"children"`
}

// fetchCatalogue fetches and decodes the doc tree of one categoryId via
// curl: the doc site fingerprints TLS and serves a JS challenge to Go's
// HTTP stack, while curl is accepted (curl must be on PATH).
func fetchCatalogue(catID string) ([]intlNode, error) {
	payload := fmt.Sprintf(`{"lang":"en","categoryId":%q}`, catID)
	cmd := exec.Command("curl", "-s", "--max-time", "30",
		"-X", "POST", intlCatURL,
		"-H", "Content-Type: application/json",
		"-H", "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36",
		"-d", payload)
	raw, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("curl: %w", err)
	}
	var envelope struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Title string     `json:"title"`
			List  []intlNode `json:"list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("response is not JSON (challenge page?): %w", err)
	}
	if envelope.Code != 0 {
		return nil, fmt.Errorf("code %d: %s", envelope.Code, envelope.Msg)
	}
	return envelope.Data.List, nil
}

// refreshIntlURLs walks the international doc trees for the products
// present in rawActions and returns (product, action) -> international URL.
func refreshIntlURLs(rawActions []rawAction) (map[string]string, error) {
	// Collect the catalog ids (== intl categoryId) per product from the
	// Chinese-site URLs; the two sites share the product category ids.
	prodCats := map[string]map[string]bool{}
	for _, a := range rawActions {
		if a.ActionCNDocUr == "" || a.Product == "" {
			continue
		}
		m := cnDocCatRe.FindStringSubmatch(a.ActionCNDocUr)
		if m == nil {
			continue
		}
		p := strings.ToLower(a.Product)
		if prodCats[p] == nil {
			prodCats[p] = map[string]bool{}
		}
		prodCats[p][m[1]] = true
	}
	if len(prodCats) == 0 {
		return nil, fmt.Errorf("no products with catalog ids found")
	}

	// Walk one tree per categoryId; entries are keyed
	// "<product>.<action>" so actions reused across products keep their
	// product-specific page.
	intlURLs := map[string]string{}
	fetched := 0
	for prod, cats := range prodCats {
		for catID := range cats {
			nodes, err := fetchCatalogue(catID)
			if err != nil {
				message("[WARN]intl catalogue for categoryId %s failed: %v\n", catID, err)
				continue
			}
			fetched++
			for action, docID := range collectActionPages(nodes) {
				key := prod + "." + action
				if _, exists := intlURLs[key]; !exists {
					intlURLs[key] = fmt.Sprintf("%s/document/product/%s/%s", intlBase, catID, docID)
				}
			}
		}
	}
	if fetched == 0 {
		return nil, fmt.Errorf("no international catalogues could be fetched")
	}
	return intlURLs, nil
}

// isActionTitle reports whether a tree page title reads like an API action
// name (CamelCase, no spaces) rather than a prose page ("Introduction").
func isActionTitle(title string) bool {
	if len(title) < 4 || title[0] < 'A' || title[0] > 'Z' {
		return false
	}
	if strings.ContainsAny(title, " _") {
		return false
	}
	// Reject all-caps prose-like titles and pure words such as "Overview".
	hasLower := false
	for i := 1; i < len(title); i++ {
		if title[i] >= 'a' && title[i] <= 'z' {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return false
	}
	for _, w := range []string{"Overview", "Introduction", "History", "Glossary", "Billing", "Strengths"} {
		if title == w {
			return false
		}
	}
	return true
}

// collectActionPages flattens the doc tree into action -> doc id.
func collectActionPages(nodes []intlNode) map[string]string {
	out := map[string]string{}
	var walk func(list []intlNode, inAPI bool)
	walk = func(list []intlNode, inAPI bool) {
		for _, n := range list {
			isAPIDir := n.Type == "directory" && strings.Contains(n.Title, "API")
			next := inAPI || isAPIDir
			if n.Type == "page" && next && isActionTitle(n.Title) {
				if _, dup := out[n.Title]; !dup {
					out[n.Title] = fmt.Sprintf("%d", n.ID)
				}
			}
			if len(n.Children) > 0 {
				walk(n.Children, next)
			}
		}
	}
	walk(nodes, false)
	return out
}
