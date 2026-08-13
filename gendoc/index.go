package main

import (
	"bufio"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Product struct {
	Name        string
	DataSources []string
	Resources   []string
	// Extra collections introduced when the framework stack was merged
	// into the same provider.md. They are populated by GetIndex when a
	// product node carries the matching sub-section header (Function /
	// Ephemeral Resource / List Resource / Action) and consumed by both
	// idxTPL and the framework rendering pipeline in framework.go.
	Functions  []string
	Ephemerals []string
	Lists      []string
	Actions    []string
}

// sectionHeaders enumerates every legal sub-section header that may
// appear under a product node in tencentcloud/provider.md. Adding a
// new framework reference type only requires adding a new entry here
// plus a switch arm in GetIndex (and a sidebar block in template.go).
var sectionHeaders = map[string]struct{}{
	"Data Source":        {},
	"Resource":           {},
	"Function":           {},
	"Ephemeral Resource": {},
	"List Resource":      {},
	"Action":             {},
}

// GetIndex parses the unified `Resources List` section that lives at
// the bottom of tencentcloud/provider.md. It understands six sub-section
// headers: "Data Source", "Resource", "Function", "Ephemeral Resource",
// "List Resource" and "Action". The historical "Provider Data Sources"
// pseudo-product (which lists data sources without a header) is still
// accepted unchanged.
func GetIndex(doc string) ([]Product, error) {
	scanner := bufio.NewScanner(strings.NewReader(doc))

	var (
		prods       []Product
		prod        Product
		currentType string
	)

	flush := func() {
		sort.Strings(prod.DataSources)
		sort.Strings(prod.Resources)
		sort.Strings(prod.Functions)
		sort.Strings(prod.Ephemerals)
		sort.Strings(prod.Lists)
		sort.Strings(prod.Actions)
		prods = append(prods, prod)
	}

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if text == "" {
			continue
		}

		if !strings.HasPrefix(text, cloudMark) {
			if _, isHeader := sectionHeaders[text]; isHeader {
				if currentType == text {
					return nil, fmt.Errorf("duplicate %s in product %s", currentType, prod.Name)
				}
				currentType = text
				continue
			}

			if prod.Name == "" {
				prod.Name = text
			} else {
				flush()
				prod = Product{Name: text}
				currentType = ""
			}
			continue
		}

		switch currentType {
		case "Data Source":
			prod.DataSources = append(prod.DataSources, text)
		case "Resource":
			prod.Resources = append(prod.Resources, text)
		case "Function":
			prod.Functions = append(prod.Functions, text)
		case "Ephemeral Resource":
			prod.Ephemerals = append(prod.Ephemerals, text)
		case "List Resource":
			prod.Lists = append(prod.Lists, text)
		case "Action":
			prod.Actions = append(prod.Actions, text)
		default:
			// Backward-compatibility: the original "Provider Data
			// Sources" pseudo-product lists data sources directly
			// without a `Data Source` header.
			if prod.Name == "Provider Data Sources" {
				prod.DataSources = append(prod.DataSources, text)
				continue
			}
			return nil, errors.New("no section header before " + text)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(prods) == 0 && prod.Name == "" {
		return nil, errors.New("no product found")
	}

	if prod.Name != "" {
		flush()
	}

	sort.Slice(prods, func(i, j int) bool {
		// make sure Provider Data Sources at first
		if prods[i].Name == "Provider Data Sources" {
			return true
		}
		if prods[j].Name == "Provider Data Sources" {
			return false
		}
		return prods[i].Name < prods[j].Name
	})

	return prods, nil
}

func replace(str, old, new string) string { return strings.Replace(str, old, new, -1) }
