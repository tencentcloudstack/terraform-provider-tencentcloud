package wedata

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataCodeMaxPermission() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataCodeMaxPermissionRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique ID of authorization resource, folder ID or file ID.",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User's recursive maximum permission type for CodeStudio files/folders.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Authorization permission type (CAN_VIEW/CAN_RUN/CAN_EDIT/CAN_MANAGE).",
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

func dataSourceTencentCloudWedataCodeMaxPermissionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_code_max_permission.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(nil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		projectId  string
		resourceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("resource_id"); ok {
		paramMap["ResourceId"] = helper.String(v.(string))
		resourceId = v.(string)
	}

	var respData *wedatav20250806.CodeStudioMaxPermission
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataCodeMaxPermissionByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataMap := map[string]interface{}{}
	if respData != nil {
		if respData.PermissionType != nil {
			dataMap["permission_type"] = respData.PermissionType
		}

		_ = d.Set("data", []interface{}{dataMap})
	}

	d.SetId(strings.Join([]string{projectId, resourceId}, tccommon.FILED_SP))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataMap); e != nil {
			return e
		}
	}

	return nil
}
