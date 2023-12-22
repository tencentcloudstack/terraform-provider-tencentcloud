package rum

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRumLogStatsLogList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumLogStatsLogListRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time but is represented using a timestamp in seconds.",
			},

			"query": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Log Query syntax statement.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time but is represented using a timestamp in seconds.",
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

func dataSourceTencentCloudRumLogStatsLogListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_rum_log_stats_log_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		startTime string
		endTime   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("start_time"); ok {
		startTime = v.(string)
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		endTime = v.(string)
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("query"); ok {
		paramMap["Query"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	service := RumService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result *string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeRumLogStatsLogListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
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
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
