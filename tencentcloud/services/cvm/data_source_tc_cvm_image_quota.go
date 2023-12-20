package cvm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmImageQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmImageQuotaRead,
		Schema: map[string]*schema.Schema{
			"image_num_quota": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The image quota of an account.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCvmImageQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_image_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var imageNumQuota int64
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmImageQuotaByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		imageNumQuota = result
		return nil
	})
	if err != nil {
		return err
	}

	_ = d.Set("image_num_quota", imageNumQuota)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), map[string]interface{}{
			"image_num_quota": imageNumQuota,
		}); e != nil {
			return e
		}
	}
	return nil
}
