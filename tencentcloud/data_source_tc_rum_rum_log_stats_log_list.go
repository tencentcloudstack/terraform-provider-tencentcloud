/*
Use this data source to query detailed information of rum rum_log_stats_log_list

Example Usage

```hcl
data "tencentcloud_rum_rum_log_stats_log_list" "rum_log_stats_log_list" {
  start_time = 1625444040
  query = "id:123 AND type:&quot;log&quot;"
  end_time = 1625454840
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

func dataSourceTencentCloudRumRumLogStatsLogList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumRumLogStatsLogListRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start time but is represented using a timestamp in seconds.",
			},

			"query": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Log Query syntax statement.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End time but is represented using a timestamp in seconds.",
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

func dataSourceTencentCloudRumRumLogStatsLogListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_rum_log_stats_log_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
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
		result, e := service.DescribeRumRumLogStatsLogListByFilter(ctx, paramMap)
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
