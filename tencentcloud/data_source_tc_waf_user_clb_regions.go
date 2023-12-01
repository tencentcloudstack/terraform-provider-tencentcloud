/*
Use this data source to query detailed information of waf user_clb_regions

Example Usage

```hcl
data "tencentcloud_waf_user_clb_regions" "example" {}
```
*/
package tencentcloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
)

func dataSourceTencentCloudWafUserClbRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafUserClbRegionsRead,
		Schema: map[string]*schema.Schema{
			"data": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Region list(ap-xxx format).",
			},
			"rich_datas": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Detail info for region.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region ID.",
						},
						"text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Chinese description for region.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "English description for region.",
						},
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region code.",
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

func dataSourceTencentCloudWafUserClbRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_user_clb_regions.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId          = getLogId(contextNil)
		ctx            = context.WithValue(context.TODO(), logIdKey, logId)
		service        = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
		userClbRegions *waf.DescribeUserClbWafRegionsResponseParams
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafUserClbRegionsByFilter(ctx)
		if e != nil {
			return retryError(e)
		}

		userClbRegions = result
		return nil
	})

	if err != nil {
		return err
	}

	if userClbRegions.Data != nil {
		_ = d.Set("data", userClbRegions.Data)
	}

	if userClbRegions.RichDatas != nil {
		tmpList := make([]map[string]interface{}, 0, len(userClbRegions.RichDatas))
		for _, clbWafRegionItem := range userClbRegions.RichDatas {
			clbWafRegionItemMap := map[string]interface{}{}

			if clbWafRegionItem.Id != nil {
				clbWafRegionItemMap["id"] = clbWafRegionItem.Id
			}

			if clbWafRegionItem.Text != nil {
				clbWafRegionItemMap["text"] = clbWafRegionItem.Text
			}

			if clbWafRegionItem.Value != nil {
				clbWafRegionItemMap["value"] = clbWafRegionItem.Value
			}

			if clbWafRegionItem.Code != nil {
				clbWafRegionItemMap["code"] = clbWafRegionItem.Code
			}

			tmpList = append(tmpList, clbWafRegionItemMap)
		}

		_ = d.Set("rich_datas", tmpList)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
