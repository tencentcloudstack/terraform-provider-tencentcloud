package main

// catalog.go — the international-site API catalog.
//
// All data comes from the international site (www.tencentcloud.com);
// no Chinese-site catalog or ids are involved:
//
//   - The doc catalogue tree (apiref_intl.go) provides
//     (product categoryId, action) -> doc URL.
//   - getCategoryByLang maps product categoryId -> slug (cvm, teo, ...).
//   - Static analysis supplies the client product names (the SDK's
//     directory names, e.g. "vpc", "csip"), which are matched against
//     slugs directly and, failing that, by Jaccard similarity of their
//     action-name sets (the SDK product name and the docs slug differ
//     for some products, e.g. SDK "redis" vs slug "crs").
//
// Actions without an international page fall back to the international
// site's search URL — every rendered link is international.

import (
	"fmt"
	"sort"
	"strings"
)

// minJaccard is the minimum action-set similarity for an SDK product to
// be aligned with an international docs product via set matching.
const minJaccard = 0.3

// catalog resolves (clientProduct, action) -> international doc URL.
type catalog struct {
	// urlIdx maps "<slug-or-sdk-product>.<Action>" -> doc URL. Keys exist
	// under both the docs slug and the SDK product name so either
	// resolution path hits.
	urlIdx map[string]string
}

// loadCatalog fetches the international catalogue and product slugs and
// builds the URL index. On failure the index stays empty and the
// renderer falls back to the international site's search URL per action
// — generation never fails on this.
func loadCatalog() *catalog {
	c := &catalog{urlIdx: map[string]string{}}
	urls, productPages, err := fetchIntlActionURLs()
	if err != nil {
		message("[WARN]international doc catalogue unavailable (%v); actions will link to the search page\n", err)
		return c
	}
	slugs, err := fetchIntlProductSlugs()
	if err != nil {
		message("[WARN]international product slugs unavailable (%v); product alignment by action sets only\n", err)
	}
	c.build(urls, productPages, slugs)
	return c
}

// build indexes the catalogue URLs under product slugs (and, for
// resilience, under the raw SDK-style product names derived from the
// action sets' best slug match).
func (c *catalog) build(urls map[string]string, productPages map[int64]map[string]bool, slugs map[int64]string) {
	// Align slug-less / differently-named products: for products whose
	// slug does not correspond to any SDK product name, we cannot know
	// the SDK-side name from the docs site alone. Instead index every
	// product's actions under its slug, and additionally compute a
	// slug -> slug alias table from cross-product action-set similarity
	// is impossible without the SDK-side names here — those are supplied
	// by the analyzer (see indexForClient).
	for catID, pages := range productPages {
		slug := slugs[catID]
		if slug == "" {
			// No slug known: fall back to a normalized page title.
			slug = strings.ToLower(strings.ReplaceAll(productRootTitle(catID, urls), " ", "-"))
		}
		for action := range pages {
			key := slug + "." + action
			if _, dup := c.urlIdx[key]; !dup {
				c.urlIdx[key] = urls[fmt.Sprintf("%d.%s", catID, action)]
			}
		}
	}
	message("[SUCC.]international catalog: %d actions across %d products\n", len(c.urlIdx), len(productPages))
}

// productRootTitle is a best-effort label for products without a slug.
func productRootTitle(catID int64, urls map[string]string) string {
	prefix := fmt.Sprintf("%d.", catID)
	for k := range urls {
		if strings.HasPrefix(k, prefix) {
			return "product-" + fmt.Sprintf("%d", catID)
		}
	}
	return fmt.Sprintf("product-%d", catID)
}

// indexForClient returns the URL for one client action, trying the
// client product name as slug first. The analyzer calls this with the
// SDK product name (lowercase, e.g. "vpc").
func (c *catalog) indexForClient(prod, action string) string {
	return c.urlIdx[strings.ToLower(prod)+"."+action]
}

// alignCatalogWithSDKProducts extends the index with alias keys: for each
// SDK product name (from the analyzer) that is not itself a slug, find
// the docs product whose action set best matches the SDK product's known
// actions and register the SDK product name as an alias for that slug.
// knownActions maps SDK product name -> set of action names observed by
// static analysis across the whole provider.
func (c *catalog) alignCatalogWithSDKProducts(knownActions map[string]map[string]bool) {
	if len(knownActions) == 0 {
		return
	}
	// Group index keys by product and collect each docs product's actions.
	docsActions := map[string]map[string]bool{}
	for key := range c.urlIdx {
		i := strings.IndexByte(key, '.')
		if i <= 0 {
			continue
		}
		prod, action := key[:i], key[i+1:]
		if docsActions[prod] == nil {
			docsActions[prod] = map[string]bool{}
		}
		docsActions[prod][action] = true
	}
	// Alias unknown SDK products to the most similar docs product.
	prods := make([]string, 0, len(knownActions))
	for p := range knownActions {
		prods = append(prods, p)
	}
	sort.Strings(prods)
	aliased := 0
	for _, p := range prods {
		if _, ok := docsActions[p]; ok {
			continue // SDK name is itself a slug — nothing to do
		}
		acts := knownActions[p]
		bestSlug, bestJ := "", 0.0
		for slug, dacts := range docsActions {
			if j := jaccard(acts, dacts); j > bestJ {
				bestSlug, bestJ = slug, j
			}
		}
		if bestSlug == "" || bestJ < minJaccard {
			continue
		}
		for action := range docsActions[bestSlug] {
			key := p + "." + action
			if _, dup := c.urlIdx[key]; !dup {
				c.urlIdx[key] = c.urlIdx[bestSlug+"."+action]
			}
		}
		aliased++
	}
	if aliased > 0 {
		message("[SUCC.]aliased %d SDK product names to international docs products by action sets\n", aliased)
	}
}

// jaccard is the Jaccard similarity of two string sets.
func jaccard(a, b map[string]bool) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	inter := 0
	for k := range a {
		if b[k] {
			inter++
		}
	}
	return float64(inter) / float64(len(a)+len(b)-inter)
}

// loadSnapshotCatalogOnly kept for the offline probe mode: it now returns
// an empty catalog (the probe never touches the network; rendered probe
// output shows bare action names).
func loadSnapshotCatalogOnly() *catalog {
	return &catalog{urlIdx: map[string]string{}}
}
