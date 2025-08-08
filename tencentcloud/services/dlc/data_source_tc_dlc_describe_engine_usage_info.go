package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDescribeEngineUsageInfo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeEngineUsageInfoRead,
		Schema: map[string]*schema.Schema{
			"data_engine_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The data engine ID.",
			},

			"used": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The used cluster spec.",
			},

			"available": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The available cluster spec.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDlcDescribeEngineUsageInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_engine_usage_info.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var dataEngineId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("data_engine_id"); ok {
		dataEngineId = v.(string)
		paramMap["DataEngineId"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var describeInfo *dlc.DescribeEngineUsageInfoResponseParams
	tmp := make(map[string]interface{}, 0)

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeEngineUsageInfoByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		describeInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	if describeInfo.Used != nil {
		_ = d.Set("used", describeInfo.Used)
		tmp["used"] = describeInfo.Used
	}

	if describeInfo.Available != nil {
		_ = d.Set("available", describeInfo.Available)
		tmp["available"] = describeInfo.Available
	}

	d.SetId(dataEngineId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmp); e != nil {
			return e
		}
	}
	return nil
}
