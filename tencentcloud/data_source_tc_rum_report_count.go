/*
Use this data source to query detailed information of rum report_count

Example Usage

```hcl
data "tencentcloud_rum_report_count" "report_count" {
  start_time = 1625444040
  end_time = 1625454840
  i_d = 1
  report_type = "log"
  instance_i_d = "rum-xxx"
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

func dataSourceTencentCloudRumReportCount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumReportCountRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start time but is represented using a timestamp in seconds.",
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

			"report_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Report type, empty is meaning all type count. `log`:&amp;amp;#39;log report count&amp;amp;#39;, `pv`:&amp;amp;#39;pv report count&amp;amp;#39;, `event`:&amp;amp;#39;event report count&amp;amp;#39;, `speed`:&amp;amp;#39;speed report count&amp;amp;#39;, `performance`:&amp;amp;#39;performance report count&amp;amp;#39;, `custom`:&amp;amp;#39;custom report count&amp;amp;#39;, `webvitals`:&amp;amp;#39;webvitals report count&amp;amp;#39;, `miniProgramData`:&amp;amp;#39;miniProgramData report count&amp;amp;#39; .",
			},

			"instance_i_d": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
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

func dataSourceTencentCloudRumReportCountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_report_count.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("i_d"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("report_type"); ok {
		paramMap["ReportType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_i_d"); ok {
		paramMap["InstanceID"] = helper.String(v.(string))
	}

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRumReportCountByFilter(ctx, paramMap)
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
