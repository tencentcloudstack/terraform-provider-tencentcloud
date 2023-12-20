package cvm

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmImageSharePermission() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmImageSharePermissionRead,
		Schema: map[string]*schema.Schema{
			"image_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of the image to be shared.",
			},

			"share_permission_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Information on image sharing.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when an image was shared.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the account with which the image is shared.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCvmImageSharePermissionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_image_share_permission.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("image_id"); ok {
		paramMap["ImageId"] = helper.String(v.(string))
	}

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var sharePermissionSet []*cvm.SharePermission

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCvmImageSharePermissionByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		sharePermissionSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(sharePermissionSet))
	tmpList := make([]map[string]interface{}, 0, len(sharePermissionSet))

	if sharePermissionSet != nil {
		for _, sharePermission := range sharePermissionSet {
			sharePermissionMap := map[string]interface{}{}

			if sharePermission.CreatedTime != nil {
				sharePermissionMap["created_time"] = sharePermission.CreatedTime
			}

			if sharePermission.AccountId != nil {
				sharePermissionMap["account_id"] = sharePermission.AccountId
			}

			ids = append(ids, *sharePermission.AccountId)
			tmpList = append(tmpList, sharePermissionMap)
		}

		_ = d.Set("share_permission_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
