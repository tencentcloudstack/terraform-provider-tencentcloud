/*
Use this data source to query detailed information of rum log_list

Example Usage

```hcl
data "tencentcloud_rum_log_list" "log_list" {
  order_by = "desc"
  start_time = 1625444040000
  query = "id:123 AND type:&quot;log&quot;"
  end_time = 1625454840000
  i_d = 1
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRumLogList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumLogListRead,
		Schema: map[string]*schema.Schema{
			"order_by": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Sorting method. `desc`:&amp;amp;#39;Descending order&amp;amp;#39;; `asc`: &amp;amp;#39;Ascending order&amp;amp;#39;.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start time but is represented using a timestamp in milliseconds.",
			},

			"query": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Log Query syntax statement.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End time but is represented using a timestamp in milliseconds.",
			},

			"i_d": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Return value.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudRumLogListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_log_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("query"); ok {
		paramMap["Query"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("i_d"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRumLogListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		_ = d.Set("result", result)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
