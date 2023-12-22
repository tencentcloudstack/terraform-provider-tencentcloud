package oceanus

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOceanusCheckSavepoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusCheckSavepointRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job id.",
			},
			"serial_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Snapshot resource ID.",
			},
			"record_type": {
				Required:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateAllowedIntValue(RECORD_TYPE),
				Description:  "Snapshot type. 1:savepoint; 2:checkpoint; 3:cancelWithSavepoint.",
			},
			"savepoint_path": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Snapshot path, currently only supports COS path.",
			},
			"work_space_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Workspace ID.",
			},
			"savepoint_status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "1=available, 2=unavailable.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudOceanusCheckSavepointRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_oceanus_check_savepoint.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		checkSavepoint *oceanus.CheckSavepointResponseParams
		serialId       string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("serial_id"); ok {
		paramMap["SerialId"] = helper.String(v.(string))
		serialId = v.(string)
	}

	if v, ok := d.GetOkExists("record_type"); ok {
		paramMap["RecordType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("savepoint_path"); ok {
		paramMap["SavepointPath"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusCheckSavepointByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		checkSavepoint = result
		return nil
	})

	if err != nil {
		return err
	}

	if checkSavepoint.SerialId != nil {
		_ = d.Set("serial_id", checkSavepoint.SerialId)
	}

	if checkSavepoint.SavepointStatus != nil {
		_ = d.Set("savepoint_status", checkSavepoint.SavepointStatus)
	}

	d.SetId(serialId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
