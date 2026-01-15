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

func DataSourceTencentCloudWedataWorkflowMaxPermission() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataWorkflowMaxPermissionRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"entity_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authorization entity ID.",
			},

			"entity_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Authorization entity type, folder/workflow.",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Current user's recursive maximum permission type for entity resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Authorization permission type (CAN_VIEW/CAN_RUN/CAN_EDIT/CAN_MANAGE, currently only supports CAN_MANAGE).",
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

func dataSourceTencentCloudWedataWorkflowMaxPermissionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_workflow_max_permission.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(nil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		projectId  string
		entityId   string
		entityType string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("entity_id"); ok {
		paramMap["EntityId"] = helper.String(v.(string))
		entityId = v.(string)
	}

	if v, ok := d.GetOk("entity_type"); ok {
		paramMap["EntityType"] = helper.String(v.(string))
		entityType = v.(string)
	}

	var respData *wedatav20250806.WorkflowMaxPermission
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataWorkflowMaxPermissionByFilter(ctx, paramMap)
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

	d.SetId(strings.Join([]string{projectId, entityId, entityType}, tccommon.FILED_SP))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataMap); e != nil {
			return e
		}
	}

	return nil
}
