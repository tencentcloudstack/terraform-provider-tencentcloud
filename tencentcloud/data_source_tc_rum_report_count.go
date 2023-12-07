package tencentcloud

import (
	"context"
	"strconv"

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

			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"report_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Report type, empty is meaning all type count. `log`:log report count, `pv`:pv report count, `event`:event report count, `speed`:speed report count, `performance`:performance report count, `custom`:custom report count, `webvitals`:webvitals report count, `miniProgramData`:miniProgramData report count.",
			},

			"instance_id": {
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

	var (
		startTime int
		endTime   int
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("start_time"); v != nil {
		startTime = v.(int)
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		endTime = v.(int)
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("report_type"); ok {
		paramMap["ReportType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceID"] = helper.String(v.(string))
	}

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeRumReportCountByFilter(ctx, paramMap)
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

	d.SetId(helper.DataResourceIdsHash([]string{strconv.Itoa(startTime), strconv.Itoa(endTime), ids}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
