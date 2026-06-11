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
	// Framework-only collections. They are populated by
	// mergeProductsForErb in main.go and consumed by idxTPL when a
	// product entry comes from the framework provider.md index.
	Functions  []string
	Ephemerals []string
	Lists      []string
	Actions    []string
}

func GetIndex(doc string) ([]Product, error) {
	scanner := bufio.NewScanner(strings.NewReader(doc))

	var (
		prods       []Product
		prod        Product
		currentType string
	)

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if text == "" {
			continue
		}

		if !strings.HasPrefix(text, cloudMark) {
			switch text {
			case "Data Source", "Resource":
				if currentType == text {
					return nil, fmt.Errorf("duplicate %s in product %s", currentType, prod.Name)
				}

				currentType = text

			default:
				if prod.Name == "" {
					prod.Name = text
				} else {
					sort.Slice(prod.DataSources, func(i, j int) bool {
						return prod.DataSources[i] < prod.DataSources[j]
					})

					sort.Slice(prod.Resources, func(i, j int) bool {
						return prod.Resources[i] < prod.Resources[j]
					})

					prods = append(prods, prod)
					prod = Product{Name: text}
					currentType = ""
				}
			}
		} else {
			switch currentType {
			case "Data Source":
				prod.DataSources = append(prod.DataSources, text)

			case "Resource":
				prod.Resources = append(prod.Resources, text)

			default:
				if prod.Name == "Provider Data Sources" {
					prod.DataSources = append(prod.DataSources, text)

					continue
				}

				return nil, errors.New("no Data Source or Resource")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(prods) == 0 {
		return nil, errors.New("no product found")
	}

	if len(prod.DataSources) == 0 && len(prod.Resources) == 0 {
		return nil, fmt.Errorf("product %s has no data source and resource", prod.Name)
	}

	prods = append(prods, prod)

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

// getFrameworkIndex parses the framework provider.md "Resources List"
// section. Compared to GetIndex it understands four extra section
// headers — "Function", "Ephemeral Resource", "List Resource",
// "Action" — and stores them on the frameworkProduct struct (defined
// alongside in framework.go). The on-disk syntax for each block is the
// same as the sdkv2 provider.md:
//
//	<Product Name>
//	<Section Header>
//	tencentcloud_xxx
//	tencentcloud_yyy
//
//	<Section Header>
//	tencentcloud_zzz
//
// A blank line is interpreted as a product separator (mirroring the
// sdkv2 parser).
func getFrameworkIndex(doc string) ([]frameworkProduct, error) {
	scanner := bufio.NewScanner(strings.NewReader(doc))

	var (
		prods       []frameworkProduct
		prod        frameworkProduct
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
			switch text {
			case "Data Source", "Resource", "Function", "Ephemeral Resource", "List Resource", "Action":
				currentType = text
			default:
				if prod.Name == "" {
					prod.Name = text
				} else {
					flush()
					prod = frameworkProduct{Name: text}
					currentType = ""
				}
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
			return nil, errors.New("framework index: missing type header before " + text)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if prod.Name != "" {
		flush()
	}
	if len(prods) == 0 {
		return nil, errors.New("framework index: no product found")
	}
	sort.Slice(prods, func(i, j int) bool {
		// Pin "Provider Meta" first so the sidebar order is deterministic.
		if prods[i].Name == "Provider Meta" {
			return true
		}
		if prods[j].Name == "Provider Meta" {
			return false
		}
		return prods[i].Name < prods[j].Name
	})
	return prods, nil
}

func replace(str, old, new string) string { return strings.Replace(str, old, new, -1) }
