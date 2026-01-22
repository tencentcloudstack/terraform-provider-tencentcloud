package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataTenantRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataTenantRolesRead,
		Schema: map[string]*schema.Schema{
			"role_display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Role Chinese display name fuzzy search, can only pass one value.",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Main account role list.",
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

func dataSourceTencentCloudWedataTenantRolesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_tenant_roles.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("role_display_name"); ok {
		paramMap["RoleDisplayName"] = helper.String(v.(string))
	}

	var respData []*wedatav20250806.SystemRole
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataTenantRolesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, data := range respData {
			dataMap := map[string]interface{}{}
			if data.RoleId != nil {
				dataMap["role_id"] = data.RoleId
			}

			if data.RoleName != nil {
				dataMap["role_name"] = data.RoleName
			}

			if data.RoleDisplayName != nil {
				dataMap["role_display_name"] = data.RoleDisplayName
			}

			if data.Description != nil {
				dataMap["description"] = data.Description
			}

			dataList = append(dataList, dataMap)
		}

		_ = d.Set("data", dataList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataList); e != nil {
			return e
		}
	}

	return nil
}
