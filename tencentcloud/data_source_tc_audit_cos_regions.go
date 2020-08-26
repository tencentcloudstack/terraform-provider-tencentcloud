/*
Use this data source to query the region list supported by the audit cos.

Example Usage

```hcl
data "tencentcloud_audit_cos_region" "cos_region" {
  website_type   = "zh"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	auditcos "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAuditCosRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAuditCosRegionsRead,

		Schema: map[string]*schema.Schema{
			"website_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "zh",
				Description: "Site type. zh means China region, en means international region.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"cos_region_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available zones supported by cos.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cos_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cos region.",
						},
						"cos_region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cos region description.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAuditCosRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_audit_cos_regions.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	auditCosService := AuditCosService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var regions []*auditcos.CosRegionInfo
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		regions, errRet = auditCosService.DescribeRegions(ctx)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	regionList := make([]map[string]interface{}, 0, len(regions))
	ids := make([]string, 0, len(regions))
	for _, region := range regions {
		mapping := map[string]interface{}{
			"cos_region":      region.CosRegion,
			"cos_region_name": region.CosRegionName,
		}
		regionList = append(regionList, mapping)
		ids = append(ids, *region.CosRegion)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("cos_region_list", regionList)
	if err != nil {
		log.Printf("[CRITAL]%s audit cos read regions list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), regionList); e != nil {
			return e
		}
	}
	return nil
}
