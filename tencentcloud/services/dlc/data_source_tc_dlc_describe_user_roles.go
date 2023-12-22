package dlc

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDescribeUserRoles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeUserRolesRead,
		Schema: map[string]*schema.Schema{
			"fuzzy": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "List according to ARN blur.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The return results are sorted according to this field.",
			},

			"sorting": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Positive or inverted, such as DESC.",
			},

			"user_roles": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "User role information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Character ID.",
						},
						"app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User app ID.",
						},
						"uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User ID.",
						},
						"arn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role Permissions.",
						},
						"modify_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Recently modify the time stamp.",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Character description information.",
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role NameNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creator UinNote: This field may return NULL, indicating that the valid value cannot be obtained.",
						},
						"cos_permission_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "COS authorization path listNote: This field may return NULL, indicating that the valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cos_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "COS pathNote: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
									"permissions": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Permissions [Read, WRITE]Note: This field may return NULL, indicating that the valid value cannot be obtained.",
									},
								},
							},
						},
						"permission_json": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CAM strategy jsonNote: This field may return NULL, indicating that the valid value cannot be obtained.",
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

func dataSourceTencentCloudDlcDescribeUserRolesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_describe_user_roles.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("fuzzy"); ok {
		paramMap["Fuzzy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sorting"); ok {
		paramMap["Sorting"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var userRoles []*dlc.UserRole

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDescribeUserRolesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		userRoles = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(userRoles))
	tmpList := make([]map[string]interface{}, 0, len(userRoles))

	if userRoles != nil {
		for _, userRole := range userRoles {
			userRoleMap := map[string]interface{}{}

			if userRole.RoleId != nil {
				userRoleMap["role_id"] = userRole.RoleId
			}

			if userRole.AppId != nil {
				userRoleMap["app_id"] = userRole.AppId
			}

			if userRole.Uin != nil {
				userRoleMap["uin"] = userRole.Uin
			}

			if userRole.Arn != nil {
				userRoleMap["arn"] = userRole.Arn
			}

			if userRole.ModifyTime != nil {
				userRoleMap["modify_time"] = userRole.ModifyTime
			}

			if userRole.Desc != nil {
				userRoleMap["desc"] = userRole.Desc
			}

			if userRole.RoleName != nil {
				userRoleMap["role_name"] = userRole.RoleName
			}

			if userRole.Creator != nil {
				userRoleMap["creator"] = userRole.Creator
			}

			if userRole.CosPermissionList != nil {
				var cosPermissionListList []interface{}
				for _, cosPermissionList := range userRole.CosPermissionList {
					cosPermissionListMap := map[string]interface{}{}

					if cosPermissionList.CosPath != nil {
						cosPermissionListMap["cos_path"] = cosPermissionList.CosPath
					}

					if cosPermissionList.Permissions != nil {
						cosPermissionListMap["permissions"] = cosPermissionList.Permissions
					}

					cosPermissionListList = append(cosPermissionListList, cosPermissionListMap)
				}

				userRoleMap["cos_permission_list"] = cosPermissionListList
			}

			if userRole.PermissionJson != nil {
				userRoleMap["permission_json"] = userRole.PermissionJson
			}

			ids = append(ids, helper.Int64ToStr(*userRole.RoleId))
			tmpList = append(tmpList, userRoleMap)
		}

		_ = d.Set("user_roles", tmpList)
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
