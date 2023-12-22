package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcCheckDataEngineImageCanBeRollback() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcCheckDataEngineImageCanBeRollbackRead,
		Schema: map[string]*schema.Schema{
			"data_engine_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine unique id.",
			},

			"to_record_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Log record id after rollback.",
			},

			"from_record_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Log record id before rollback.",
			},

			"is_rollback": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Is it possible to roll back.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDlcCheckDataEngineImageCanBeRollbackRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_check_data_engine_image_can_be_rollback.read")()
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
	response := &dlc.CheckDataEngineImageCanBeRollbackResponseParams{}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcCheckDataEngineImageCanBeRollbackByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}
	var data = make(map[string]interface{}, 0)
	if response.ToRecordId != nil {
		_ = d.Set("to_record_id", response.ToRecordId)
		data["to_record_id"] = response.ToRecordId
	}

	if response.FromRecordId != nil {
		_ = d.Set("from_record_id", response.FromRecordId)
		data["from_record_id"] = response.FromRecordId
	}

	if response.IsRollback != nil {
		_ = d.Set("is_rollback", response.IsRollback)
		data["is_rollback"] = response.IsRollback
	}

	d.SetId(dataEngineId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), data); e != nil {
			return e
		}
	}
	return nil
}
