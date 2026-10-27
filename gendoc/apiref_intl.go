package main

// apiref_intl.go — international-site doc URLs from www.tencentcloud.com.
//
// Two sources, both from the international site only:
//
//  1. The API docs landing page (GET /document/api) is a Next.js SSG
//     page whose full product catalogue tree is embedded in the HTML as
//     __NEXT_DATA__ → props.pageProps.data: one root directory node per
//     product, whose leaf pages are the API references. A page's
//     (categoryId, id) pair forms the doc URL
//     /document/api/<categoryId>/<pageId>, and its title is the action
//     name (e.g. DescribeInstances). The root directory node's id is the
//     product's categoryId, and getCategoryByLang maps that id to the
//     product slug (cvm, teo, csip, ...) used by static analysis.
//
//  2. POST /document/api/mid/getCategoryByLang {"lang":"en"} returns the
//     product category tree (id -> slug) used to key the catalogue tree
//     by product slug.
//
// The site is fronted by EdgeOne bot management: without the challenge
// cookies the first response is a ~1 KB obfuscated JS challenge instead
// of the page. The challenge currently sets two constant cookies, so
// sending them up front passes the check (see edgeOneCookie). If the
// constants ever become dynamic, refresh them from a browser session.
//
// The fetches shell out to curl: the doc site fingerprints TLS and
// serves the EdgeOne JS challenge to Go's HTTP stack, while curl is
// accepted (curl must be on PATH).

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const (
	intlDocURL  = "https://www.tencentcloud.com/document/api"
	intlCatURL  = "https://www.tencentcloud.com/document/api/mid/getCategoryByLang"
	intlBaseURL = "https://www.tencentcloud.com"
)

// EdgeOne bot-management challenge cookies (currently hard-coded
// constants inside the challenge JS). Sending them up front skips the
// JS challenge.
const edgeOneCookie = "__tst_status=1942791768#; EO_Bot_Ssid=1254555648"

// intlUserAgent is the browser UA accepted by the doc site.
const intlUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"

// intlNode is one node of the __NEXT_DATA__ catalogue tree. The id
// fields are numeric in deep nodes but string-typed on some top-level
// product roots (e.g. "id": "213"), so both are accepted.
type intlNode struct {
	ID         json.Number `json:"id"`
	CategoryID json.Number `json:"categoryId"`
	Title      string      `json:"title"`
	Type       string      `json:"type"`
	Link       string      `json:"link"`
	Children   []intlNode  `json:"children"`
}

// nodeID returns the int64 value of a node id field, 0 when absent.
func nodeID(n json.Number) int64 {
	v, err := n.Int64()
	if err != nil {
		return 0
	}
	return v
}

// nextDataRe extracts the __NEXT_DATA__ JSON payload from the page HTML.
var nextDataRe = regexp.MustCompile(`(?s)<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)

// intlLinkRe parses a node's relative link "/document/api/<cat>/<page>".
var intlLinkRe = regexp.MustCompile(`^/document/api/(\d+)/(\d+)$`)

// fetchIntlActionURLs downloads the international API docs landing page
// and flattens its catalogue tree into (productCategoryId, action) ->
// absolute doc URL plus per-product action sets (keyed by categoryId,
// for cross-referencing with the slug mapping from fetchIntlProductSlugs).
func fetchIntlActionURLs() (map[string]string, map[int64]map[string]bool, error) {
	body := fetchIntlPage(intlDocURL)
	if !strings.Contains(body, "__NEXT_DATA__") {
		// ~1 KB obfuscated JS: the EdgeOne challenge. The cookie constants
		// above no longer work and must be refreshed.
		return nil, nil, fmt.Errorf("EdgeOne challenge page served (%d bytes); refresh edgeOneCookie", len(body))
	}
	m := nextDataRe.FindSubmatch([]byte(body))
	if m == nil {
		return nil, nil, fmt.Errorf("__NEXT_DATA__ script not found in %s HTML", intlDocURL)
	}
	var payload struct {
		Props struct {
			PageProps struct {
				Data []intlNode `json:"data"`
			} `json:"pageProps"`
		} `json:"props"`
	}
	if err := json.Unmarshal(m[1], &payload); err != nil {
		return nil, nil, fmt.Errorf("parse __NEXT_DATA__: %w", err)
	}
	if len(payload.Props.PageProps.Data) == 0 {
		return nil, nil, fmt.Errorf("__NEXT_DATA__ catalogue tree is empty")
	}

	urls := map[string]string{}
	// productPages maps a top-level product categoryId to the set of page
	// titles found under it.
	productPages := map[int64]map[string]bool{}
	var walk func(nodes []intlNode, prod int64)
	walk = func(nodes []intlNode, prod int64) {
		for _, n := range nodes {
			prodHere := prod
			if n.Type == "directory" && prodHere == 0 {
				if id := nodeID(n.ID); id != 0 {
					prodHere = id
				}
			}
			if n.Type == "page" {
				catID, pageID := nodeID(n.CategoryID), nodeID(n.ID)
				if (catID == 0 || pageID == 0) && n.Link != "" {
					if lm := intlLinkRe.FindStringSubmatch(n.Link); lm != nil {
						catID = parseInt64(lm[1])
						pageID = parseInt64(lm[2])
					}
				}
				if catID != 0 && pageID != 0 {
					key := fmt.Sprintf("%d.%s", catID, n.Title)
					if _, dup := urls[key]; !dup {
						urls[key] = fmt.Sprintf("%s/document/api/%d/%d", intlBaseURL, catID, pageID)
					}
					if prodHere != 0 {
						if productPages[prodHere] == nil {
							productPages[prodHere] = map[string]bool{}
						}
						productPages[prodHere][n.Title] = true
					}
				}
			}
			walk(n.Children, prodHere)
		}
	}
	walk(payload.Props.PageProps.Data, 0)
	if len(urls) == 0 {
		return nil, nil, fmt.Errorf("no action pages found in the __NEXT_DATA__ catalogue tree")
	}
	return urls, productPages, nil
}

// fetchIntlPage fetches one doc-site URL via curl with the EdgeOne
// cookies and a browser UA.
func fetchIntlPage(url string) string {
	cmd := exec.Command("curl", "-s", "--max-time", "120",
		url,
		"-H", "Cookie: "+edgeOneCookie,
		"-H", "User-Agent: "+intlUserAgent)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

// intlCatNode is one node of the getCategoryByLang product tree.
type intlCatNode struct {
	ID       int64         `json:"id"`
	Slug     string        `json:"slug"`
	Title    string        `json:"title"`
	Children []intlCatNode `json:"children"`
}

// fetchIntlProductSlugs returns categoryId -> lowercase product slug
// from the international site's product category tree.
func fetchIntlProductSlugs() (map[int64]string, error) {
	body := fetchIntlPostJSON(intlCatURL, `{"lang":"en"}`)
	if body == "" {
		return nil, fmt.Errorf("empty response from %s", intlCatURL)
	}
	var envelope struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data []intlCatNode
	}
	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		return nil, fmt.Errorf("parse %s response: %w", intlCatURL, err)
	}
	if envelope.Code != 0 {
		return nil, fmt.Errorf("%s: code %d: %s", intlCatURL, envelope.Code, envelope.Msg)
	}
	slugs := map[int64]string{}
	var walk func(nodes []intlCatNode)
	walk = func(nodes []intlCatNode) {
		for _, n := range nodes {
			if n.ID != 0 && n.Slug != "" {
				slugs[n.ID] = strings.ToLower(n.Slug)
			}
			walk(n.Children)
		}
	}
	walk(envelope.Data)
	if len(slugs) == 0 {
		return nil, fmt.Errorf("no product slugs in %s response", intlCatURL)
	}
	return slugs, nil
}

// fetchIntlPostJSON POSTs a JSON payload to the doc site's internal
// gateway via curl.
func fetchIntlPostJSON(url, payload string) string {
	cmd := exec.Command("curl", "-s", "--max-time", "60",
		"-X", "POST", url,
		"-H", "Content-Type: application/json",
		"-H", "Cookie: "+edgeOneCookie,
		"-H", "User-Agent: "+intlUserAgent,
		"-d", payload)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func parseInt64(s string) int64 {
	var v int64
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0
		}
		v = v*10 + int64(s[i]-'0')
	}
	return v
}
