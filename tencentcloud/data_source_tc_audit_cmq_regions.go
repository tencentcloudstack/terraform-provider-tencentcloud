/*
Use this data source to query the region list supported by the audit cmq.

Example Usage

```hcl
data "tencentcloud_audit_cmq_region" "cmq_region" {
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
	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAuditCmqRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAuditCmqRegionsRead,

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
			"cmq_region_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available zones supported by cmq.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cmq_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cmq region.",
						},
						"cmq_region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cmq region description.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAuditCmqRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_audit_cmq_regions.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	auditService := AuditService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var regions []*audit.CmqRegionInfo
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		regions, errRet = auditService.DescribeAuditCmqRegions(ctx)
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
			"cmq_region":      region.CmqRegion,
			"cmq_region_name": region.CmqRegionName,
		}
		regionList = append(regionList, mapping)
		ids = append(ids, *region.CmqRegion)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("cmq_region_list", regionList)
	if err != nil {
		log.Printf("[CRITAL]%s audit cmq read regions list fail, reason:%s\n ", logId, err.Error())
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
