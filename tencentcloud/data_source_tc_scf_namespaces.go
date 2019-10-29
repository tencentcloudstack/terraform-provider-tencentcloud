/*
Use this data source to query SCF namespaces.

Example Usage

```hcl
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cos_bucket_name   = "scf-code-1234567890"
  cos_object_name   = "code.zip"
  cos_bucket_region = "ap-guangzhou"
}

data "tencentcloud_scf_functions" "foo" {
  name = "${tencentcloud_scf_function.foo.name}"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func dataSourceTencentCloudScfNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfNamespacesRead,
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the SCF namespace to be queried.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the SCF namespace to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"namespaces": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of namespace. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SCF namespace.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the SCF namespace.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the SCF namespace.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modify time of the SCF namespace.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the SCF namespace.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudScfNamespacesRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_namespaces.read")
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := ScfService{client: m.(*TencentCloudClient).apiV3Conn}

	var (
		namespace *string
		desc      *string
	)

	if raw, ok := d.GetOk("namespace"); ok {
		namespace = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("description"); ok {
		desc = stringToPointer(raw.(string))
	}

	nss, err := service.DescribeNamespaces(ctx)
	if err != nil {
		log.Printf("[CRITAL]%s read namespace list failed: %+v", logId, err)
		return err
	}

	namespaces := make([]map[string]*string, 0, len(nss))
	ids := make([]string, 0, len(nss))

	for _, ns := range nss {
		if namespace != nil && !strings.Contains(*ns.Name, *namespace) {
			continue
		}
		if desc != nil && !strings.Contains(*ns.Description, *desc) {
			continue
		}

		ids = append(ids, *ns.Name)

		namespaces = append(namespaces, map[string]*string{
			"namespace":   ns.Name,
			"description": ns.Description,
			"create_time": ns.AddTime,
			"modify_time": ns.ModTime,
			"type":        ns.Type,
		})
	}

	d.Set("namespaces", namespaces)
	d.SetId(dataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), namespaces); err != nil {
			err = errors.WithStack(err)
			log.Printf("[CRITAL]%s output file[%s] fail, reason: %+v", logId, output.(string), err)
			return err
		}
	}

	return nil
}
