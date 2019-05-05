package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	cloud "github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"text/template"
)

const (
	cloudMark      = "tencentcloud"
	cloudTitle     = "TencentCloud"
	cloudMarkShort = "tc"
	docRoot        = "../website/docs"
)

func main() {
	provider := cloud.Provider()
	vProvider := runtime.FuncForPC(reflect.ValueOf(cloud.Provider).Pointer())

	fname, _ := vProvider.FileLine(0)
	fpath := filepath.Dir(fname)
	log.Printf("generating doc from: %s\n", fpath)

	// DataSourcesMap
	for k, v := range provider.DataSourcesMap {
		genDoc("data_source", fpath, k, v)
	}

	// ResourcesMap
	for k, v := range provider.ResourcesMap {
		genDoc("resource", fpath, k, v)
	}
}

// genDoc generating doc for resource
func genDoc(dtype, fpath, name string, resource *schema.Resource) {
	data := map[string]string{
		"name":              name,
		"dtype":             strings.Replace(dtype, "_", "", -1),
		"resource":          name[len(cloudMark)+1:],
		"cloud_mark":        cloudMark,
		"cloud_title":       cloudTitle,
		"example":           "",
		"description":       "",
		"description_short": "",
	}

	fname := fmt.Sprintf("%s_%s_%s.go", dtype, cloudMarkShort, data["resource"])
	log.Printf("get description from file: %s\n", fname)

	description, err := getFileDescription(fmt.Sprintf("%s/%s", fpath, fname))
	if err != nil {
		log.Printf("[SKIP]get description failed, skip: %s", err)
		return
	}

	description = strings.TrimSpace(description)
	if description == "" {
		log.Printf("[SKIP]description empty, skip: %s\n", fname)
		return
	}

	pos := strings.Index(description, "\nExample Usage\n")
	if pos != -1 {
		data["example"] = strings.TrimSpace(description[pos+15:])
		description = strings.TrimSpace(description[:pos])
	} else {
		log.Printf("[SKIP]example usage missing, skip: %s\n", fname)
		return
	}

	data["description"] = description
	pos = strings.Index(description, "\n\n")
	if pos != -1 {
		data["description_short"] = strings.TrimSpace(description[:pos])
	} else {
		data["description_short"] = description
	}

	requiredArgs := []string{}
	optionalArgs := []string{}
	attributes := []string{}

	for k, v := range resource.Schema {
		if v.Description == "" {
			continue
		}
		if v.Required {
			opt := "Required"
			if v.ForceNew {
				opt += ", ForceNew"
			}
			requiredArgs = append(requiredArgs, fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description))
		} else if v.Optional {
			opt := "Optional"
			if v.ForceNew {
				opt += ", ForceNew"
			}
			optionalArgs = append(optionalArgs, fmt.Sprintf("* `%s` - (%s) %s", k, opt, v.Description))
		} else if v.Computed {
			if v.Type == schema.TypeList {
				listAttributes := []string{}
				for kk, vv := range v.Elem.(*schema.Resource).Schema {
					if vv.Description == "" {
						continue
					}
					if v.Computed {
						listAttributes = append(listAttributes, fmt.Sprintf("  * `%s` - %s", kk, vv.Description))
					}
				}
				slistAttributes := ""
				sort.Strings(listAttributes)
				if len(listAttributes) > 0 {
					slistAttributes = "\n" + strings.Join(listAttributes, "\n")
				}
				attributes = append(attributes, fmt.Sprintf("* `%s` - %s%s", k, v.Description, slistAttributes))
			} else {
				attributes = append(attributes, fmt.Sprintf("* `%s` - %s", k, v.Description))
			}
		}
	}

	sort.Strings(requiredArgs)
	sort.Strings(optionalArgs)
	sort.Strings(attributes)

	requiredArgs = append(requiredArgs, optionalArgs...)
	data["arguments"] = strings.Join(requiredArgs, "\n")
	data["attributes"] = strings.Join(attributes, "\n")

	fname = fmt.Sprintf("%s/%s/%s.html.markdown", docRoot, dtype[0:1], data["resource"])
	fd, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("[FAIL]open file %s failed: %s", fname, err)
		return
	}

	defer fd.Close()
	t := template.Must(template.New("t").Parse(docTPL))
	err = t.Execute(fd, data)
	if err != nil {
		log.Printf("[FAIL]write file %s failed: %s", fname, err)
		return
	}

	log.Printf("[SUCC]write doc to file success: %s", fname)
}

// getFileDescription get description from go file
func getFileDescription(fname string) (string, error) {
	fset := token.NewFileSet()

	parsedAst, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	return parsedAst.Doc.Text(), nil
}
