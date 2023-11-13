/*
Use this data source to query detailed information of lighthouse zone

Example Usage

```hcl
data "tencentcloud_lighthouse_zone" "zone" {
  order_field = "ZONE"
  order = "ASC"
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

func dataSourceTencentCloudLighthouseZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseZoneRead,
		Schema: map[string]*schema.Schema{
			"order_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sorting field. Valid values:&amp;amp;lt;li&amp;amp;gt;ZONE：Sort by the availability zone.&amp;amp;lt;li&amp;amp;gt;INSTANCE_DISPLAY_LABEL：Sort by visibility labels (HIDDEN, NORMAL and SELECTED). Default: [&amp;amp;#39;HIDDEN&amp;amp;#39;, &amp;amp;#39;NORMAL&amp;amp;#39;, &amp;amp;#39;SELECTED&amp;amp;#39;].",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Specifies how availability zones are listed. Valid values:&amp;amp;lt;li&amp;amp;gt;ASC：Ascending sort. &amp;amp;lt;li&amp;amp;gt;DESC：Descending sort.The default value is ASC.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudLighthouseZoneRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_zone.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("order_field"); ok {
		paramMap["OrderField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var zoneInfoSet []*lighthouse.ZoneInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseZoneByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		zoneInfoSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(zoneInfoSet))
	tmpList := make([]map[string]interface{}, 0, len(zoneInfoSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
