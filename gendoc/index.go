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

func replace(str, old, new string) string { return strings.Replace(str, old, new, -1) }
