package main

// apiref_intl.go — international-site doc URLs from www.tencentcloud.com.
//
// The API docs landing page (GET /document/api) is a Next.js SSG page
// whose full product catalogue tree is embedded in the HTML as
// __NEXT_DATA__ → props.pageProps.data: one root directory node per
// product, whose leaf pages are the API references. A page's
// (categoryId, id) pair forms the doc URL
// /document/api/<categoryId>/<pageId>, and its title is the action name
// (e.g. DescribeInstances).
//
// The site is fronted by EdgeOne bot management: without the challenge
// cookies the first response is a ~1 KB obfuscated JS challenge instead
// of the page. The challenge currently sets two constant cookies, so
// sending them up front passes the check (see edgeOneCookie). If the
// constants ever become dynamic, refresh them from a browser session.
//
// fetchIntlActionURLs() returns (categoryId.Action -> absolute URL).
// Actions the international site has not translated fall back to the
// site's search URL at resolve time — Chinese-site URLs are never
// rendered.

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

const intlDocURL = "https://www.tencentcloud.com/document/api"

// EdgeOne bot-management challenge cookies (currently hard-coded
// constants inside the challenge JS). Sending them up front skips the
// JS challenge.
const edgeOneCookie = "__tst_status=1942791768#; EO_Bot_Ssid=1254555648"

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
// and flattens its catalogue tree into "categoryId.Action" -> absolute
// doc URL. Returns an error when the challenge page is served or the
// payload shape changed, so callers can degrade to search URLs.
//
// The fetch shells out to curl: the doc site fingerprints TLS and serves
// the EdgeOne JS challenge to Go's HTTP stack, while curl is accepted
// (curl must be on PATH).
func fetchIntlActionURLs() (map[string]string, error) {
	cmd := exec.Command("curl", "-s", "--max-time", "120",
		intlDocURL,
		"-H", "Cookie: "+edgeOneCookie,
		"-H", "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36")
	body, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("curl %s: %w", intlDocURL, err)
	}
	if !strings.Contains(string(body), "__NEXT_DATA__") {
		// ~1 KB obfuscated JS: the EdgeOne challenge. The cookie constants
		// above no longer work and must be refreshed.
		return nil, fmt.Errorf("EdgeOne challenge page served (%d bytes); refresh edgeOneCookie", len(body))
	}
	m := nextDataRe.FindSubmatch(body)
	if m == nil {
		return nil, fmt.Errorf("__NEXT_DATA__ script not found in %s HTML", intlDocURL)
	}
	var payload struct {
		Props struct {
			PageProps struct {
				Data []intlNode `json:"data"`
			} `json:"pageProps"`
		} `json:"props"`
	}
	if err := json.Unmarshal(m[1], &payload); err != nil {
		return nil, fmt.Errorf("parse __NEXT_DATA__: %w", err)
	}
	if len(payload.Props.PageProps.Data) == 0 {
		return nil, fmt.Errorf("__NEXT_DATA__ catalogue tree is empty")
	}

	urls := map[string]string{}
	var walk func(nodes []intlNode)
	walk = func(nodes []intlNode) {
		for _, n := range nodes {
			if n.Type == "page" && isActionTitle(n.Title) {
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
						urls[key] = fmt.Sprintf("https://www.tencentcloud.com/document/api/%d/%d", catID, pageID)
					}
				}
			}
			walk(n.Children)
		}
	}
	walk(payload.Props.PageProps.Data)
	if len(urls) == 0 {
		return nil, fmt.Errorf("no action pages found in the __NEXT_DATA__ catalogue tree")
	}
	return urls, nil
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
