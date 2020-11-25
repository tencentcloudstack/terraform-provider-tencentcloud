/*
Use this data source to query detailed information of TCR namespaces.

Example Usage

```hcl
data "tencentcloud_tcr_namespaces" "name" {
  instance_id 			= "cls-satg5125"
  namespace_name       = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTCRNamespaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTCRNamespacesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the instance that the namespace belongs to.",
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the TCR namespace to query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"namespace_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information list of the dedicated TCR namespaces.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of TCR namespace.",
						},
						//Computed values
						"is_public": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicate that the namespace is public or not.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTCRNamespacesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcr_namespaces.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var name, instanceId string
	instanceId = d.Get("instance_id").(string)
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	tcrService := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}
	var outErr, inErr error
	namespaces, outErr := tcrService.DescribeTCRNameSpaces(ctx, instanceId, name)
	if outErr != nil {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			namespaces, inErr = tcrService.DescribeTCRNameSpaces(ctx, instanceId, name)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
	}
	if outErr != nil {
		return outErr
	}

	ids := make([]string, 0, len(namespaces))
	namespaceList := make([]map[string]interface{}, 0, len(namespaces))
	for _, namespace := range namespaces {
		mapping := map[string]interface{}{
			"name":      namespace.Name,
			"is_public": namespace.Public,
		}

		namespaceList = append(namespaceList, mapping)
		ids = append(ids, instanceId+FILED_SP+*namespace.Name)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("namespace_list", namespaceList); e != nil {
		log.Printf("[CRITAL]%s provider set TCR namespace list fail, reason:%s\n", logId, e)
		return e
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), namespaceList); e != nil {
			return e
		}
	}

	return nil

}
