/*
Use this data source to query detailed information of rum log_list

Example Usage

```hcl
data "tencentcloud_rum_log_list" "log_list" {
  order_by   = "desc"
  start_time = 1625444040000
  query      = "id:123 AND type:\"log\""
  end_time   = 1625454840000
  project_id = 1
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
				Description: "Sorting method. `desc`:Descending order; `asc`: Ascending order.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time but is represented using a timestamp in milliseconds.",
			},

			"query": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Log Query syntax statement.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time but is represented using a timestamp in milliseconds.",
			},

			"project_id": {
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

	var (
		startTime string
		endTime   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		startTime = v.(string)
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("query"); ok {
		paramMap["Query"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		endTime = v.(string)
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeRumLogListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	var ids string
	if result != nil {
		ids = *result
		_ = d.Set("result", result)
	}

	d.SetId(helper.DataResourceIdsHash([]string{startTime, endTime, ids}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
