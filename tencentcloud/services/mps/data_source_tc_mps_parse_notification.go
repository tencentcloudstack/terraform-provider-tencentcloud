package mps

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMpsParseNotification() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsParseNotificationRead,
		Schema: map[string]*schema.Schema{
			"content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Event notification obtained from CMQ.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMpsParseNotificationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mps_parse_notification.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("content"); ok {
		paramMap["Content"] = helper.String(v.(string))
	}
	var (
		taskId string
		result *mps.ParseNotificationResponseParams
		e      error
	)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e = service.DescribeMpsParseNotificationByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		return nil
	})
	if err != nil {
		return err
	}
	if result != nil {
		taskId = *result.WorkflowTaskEvent.TaskId
	}

	d.SetId(taskId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
