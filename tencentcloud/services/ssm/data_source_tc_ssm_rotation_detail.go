package ssm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudSsmRotationDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmRotationDetailRead,
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Secret name.",
			},
			"enable_rotation": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to allow rotation.",
			},
			"frequency": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The rotation frequency, in days, defaults to 1 day.",
			},
			"latest_rotate_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Time of last rotation.",
			},
			"next_rotate_begin_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The time to start the next rotation.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSsmRotationDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ssm_rotation_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = SsmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		rotationDetail *ssm.DescribeRotationDetailResponseParams
		secretName     string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("secret_name"); ok {
		paramMap["SecretName"] = helper.String(v.(string))
		secretName = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmRotationDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		rotationDetail = result
		return nil
	})

	if err != nil {
		return err
	}

	if rotationDetail.EnableRotation != nil {
		_ = d.Set("enable_rotation", rotationDetail.EnableRotation)
	}

	if rotationDetail.Frequency != nil {
		_ = d.Set("frequency", rotationDetail.Frequency)
	}

	if rotationDetail.LatestRotateTime != nil {
		_ = d.Set("latest_rotate_time", rotationDetail.LatestRotateTime)
	}

	if rotationDetail.NextRotateBeginTime != nil {
		_ = d.Set("next_rotate_begin_time", rotationDetail.NextRotateBeginTime)
	}

	d.SetId(secretName)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
