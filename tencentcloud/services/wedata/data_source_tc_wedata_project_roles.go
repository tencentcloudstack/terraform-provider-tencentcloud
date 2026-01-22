package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataProjectRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataProjectRolesRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID.",
			},

			"role_display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Role Chinese display name fuzzy search, can only pass one value.",
			},

			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Role information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role ID.",
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role name.",
						},
						"role_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role display name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
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

func dataSourceTencentCloudWedataProjectRolesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_project_roles.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(nil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		projectId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("project_id"); ok {
		paramMap["ProjectId"] = helper.String(v.(string))
		projectId = v.(string)
	}

	if v, ok := d.GetOk("role_display_name"); ok {
		paramMap["RoleDisplayName"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.SystemRole
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataProjectRolesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	itemsList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, items := range respData {
			itemsMap := map[string]interface{}{}
			if items.RoleId != nil {
				itemsMap["role_id"] = items.RoleId
			}

			if items.RoleName != nil {
				itemsMap["role_name"] = items.RoleName
			}

			if items.RoleDisplayName != nil {
				itemsMap["role_display_name"] = items.RoleDisplayName
			}

			if items.Description != nil {
				itemsMap["description"] = items.Description
			}

			itemsList = append(itemsList, itemsMap)
		}

		_ = d.Set("items", itemsList)
	}

	d.SetId(projectId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), itemsList); e != nil {
			return e
		}
	}

	return nil
}
