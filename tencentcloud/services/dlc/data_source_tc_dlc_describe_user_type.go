package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDescribeUserType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeUserTypeRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User id (uin), if left blank, it defaults to the caller's sub-uin.",
			},

			"user_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "User type, only support: ADMIN: ddministrator/COMMON: ordinary user.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDlcDescribeUserTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_user_type.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	var userId string
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		paramMap["UserId"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var userType *string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeUserTypeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		userType = result
		return nil
	})
	if err != nil {
		return err
	}

	if userType != nil {
		_ = d.Set("user_type", userType)
	}

	d.SetId(userId + tccommon.FILED_SP + *userType)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), userType); e != nil {
			return e
		}
	}
	return nil
}
