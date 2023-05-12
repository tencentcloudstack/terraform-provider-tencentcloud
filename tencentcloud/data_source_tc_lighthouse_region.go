/*
Use this data source to query detailed information of lighthouse region

Example Usage

```hcl
data "tencentcloud_lighthouse_region" "region" {
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseRegionRead,
		Schema: map[string]*schema.Schema{
			"region_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Region information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region name.",
						},
						"region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region description.",
						},
						"region_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region availability status.",
						},
						"is_china_mainland": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the region is in the Chinese mainland.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudLighthouseRegionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_region.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var regionSet []*lighthouse.RegionInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseRegionByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		regionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(regionSet))
	tmpList := make([]map[string]interface{}, 0, len(regionSet))

	if regionSet != nil {
		for _, regionInfo := range regionSet {
			regionInfoMap := map[string]interface{}{}

			if regionInfo.Region != nil {
				regionInfoMap["region"] = regionInfo.Region
			}

			if regionInfo.RegionName != nil {
				regionInfoMap["region_name"] = regionInfo.RegionName
			}

			if regionInfo.RegionState != nil {
				regionInfoMap["region_state"] = regionInfo.RegionState
			}

			if regionInfo.IsChinaMainland != nil {
				regionInfoMap["is_china_mainland"] = regionInfo.IsChinaMainland
			}

			ids = append(ids, *regionInfo.Region)
			tmpList = append(tmpList, regionInfoMap)
		}

		_ = d.Set("region_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
