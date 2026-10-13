package main

// apiref_probe.go — PoC verification entry. Run with:
//
//	go run . -probe=vpc
//
// It prints the lifecycle -> cloud API mapping for a single resource
// without touching the normal doc generation pipeline.

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cloud "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
)

// probeAllResource prints the lifecycle -> cloud API mapping for one resource.
func probeAllResource(repoRoot string) {
	analyzer, err := newAPIAnalyzer(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init analyzer: %v\n", err)
		os.Exit(1)
	}

	var names []string
	for name := range cloud.Provider().ResourcesMap {
		names = append(names, name)
	}
	sort.Strings(names)

	covered, empty := 0, 0
	for _, name := range names {
		res := cloud.Provider().ResourcesMap[name]
		refsByVerb := collectResourceAPIRefs(analyzer, repoRoot, name, res)
		total := 0
		for _, list := range refsByVerb {
			total += len(list)
		}
		if total > 0 {
			covered++
		} else {
			empty++
			fmt.Printf("NO-API %s\n", name)
		}
	}
	fmt.Fprintf(os.Stderr, "\nresources: %d, with APIs: %d (%.1f%%), no APIs found: %d\n",
		len(names), covered, float64(covered)*100/float64(len(names)), empty)
}

// collectResourceAPIRefs traces one resource's lifecycle functions.
func collectResourceAPIRefs(analyzer *apiAnalyzer, repoRoot, resName string, res *schema.Resource) map[string][]apiRef {
	resFile := findResourceFile(repoRoot, resName)
	if resFile == "" {
		return nil
	}
	refsByVerb := map[string][]apiRef{}
	for verb, fnPtr := range map[string]func(*schema.ResourceData, interface{}) error{
		"Create": res.Create, "Read": res.Read, "Update": res.Update, "Delete": res.Delete,
	} {
		if fnPtr == nil {
			continue
		}
		fd := findFuncDeclByPointer(analyzer, resFile, fnPtr)
		if fd == nil {
			continue
		}
		refsByVerb[verb] = analyzer.analyzeLifecycle(fd)
	}
	// dedup per verb
	out := map[string][]apiRef{}
	for verb, list := range refsByVerb {
		seen := map[apiRef]bool{}
		var uniq []apiRef
		for _, r := range list {
			if !seen[r] {
				seen[r] = true
				uniq = append(uniq, r)
			}
		}
		out[verb] = uniq
	}
	return out
}

// probeResource prints the lifecycle -> cloud API mapping for one resource.
func probeResource(repoRoot, shortName string) {
	resName := "tencentcloud_" + shortName
	res, ok := cloud.Provider().ResourcesMap[resName]
	if !ok {
		fmt.Fprintf(os.Stderr, "resource %s not found in ResourcesMap\n", resName)
		os.Exit(1)
	}
	resFile := findResourceFile(repoRoot, resName)
	if resFile == "" {
		fmt.Fprintf(os.Stderr, "resource file for %s not found\n", resName)
		os.Exit(1)
	}
	fmt.Printf("resource: %s\nfile: %s\n", resName, relToRepo(repoRoot, resFile))

	analyzer, err := newAPIAnalyzer(repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init analyzer: %v\n", err)
		os.Exit(1)
	}

	lifecycleFns := map[string]func(*schema.ResourceData, interface{}) error{
		"Create": res.Create, "Read": res.Read, "Update": res.Update, "Delete": res.Delete,
	}
	refsByVerb := map[string][]apiRef{}
	for verb, fnPtr := range lifecycleFns {
		if fnPtr == nil {
			continue
		}
		fd := findFuncDeclByPointer(analyzer, resFile, fnPtr)
		if fd == nil {
			fmt.Printf("[%s] function not located\n", verb)
			continue
		}
		refs := analyzer.analyzeLifecycle(fd)
		refsByVerb[verb] = refs
		fmt.Printf("\n[%s] %s:\n", verb, fd.Name.Name)
		for _, r := range refs {
			fmt.Printf("    %s (%s) %s\n", r.Action, r.Product, r.URL)
		}
	}

	fmt.Println("\n=== rendered section ===")
	fmt.Println(renderAPIRefs(dedupVerbRefs(refsByVerb)))
}



func findResourceFile(repoRoot, resName string) string {
	short := strings.TrimPrefix(resName, "tencentcloud_")
	candidates := []string{
		fmt.Sprintf("resource_tc_%s.go", short),
		// Some resources share an implementation file with an
		// "_operation" suffixed sibling (e.g. teo_purge_task ->
		// teo_purge_task_operation), or a "_config" suffixed one
		// (e.g. teo_prefetch_origin_limit ->
		// teo_prefetch_origin_limit_config).
		fmt.Sprintf("resource_tc_%s_operation.go", short),
		fmt.Sprintf("resource_tc_%s_config.go", short),
	}
	var found string
	filepath.Walk(filepath.Join(repoRoot, "tencentcloud", "services"), func(path string, fi os.FileInfo, err error) error {
		if err != nil || found != "" || fi.IsDir() {
			return nil
		}
		for _, c := range candidates {
			if filepath.Base(path) == c {
				found = path
			}
		}
		return nil
	})
	return found
}

func relToRepo(repoRoot, path string) string {
	rel, err := filepath.Rel(repoRoot, path)
	if err != nil {
		return path
	}
	return rel
}

// runtimeFuncName extracts the bare function name from a function pointer.
func runtimeFuncName(fn interface{}) string {
	v := reflect.ValueOf(fn)
	if !v.IsValid() || v.IsNil() {
		return ""
	}
	full := runtime.FuncForPC(v.Pointer()).Name()
	if i := strings.LastIndex(full, "."); i >= 0 {
		return full[i+1:]
	}
	return full
}
