package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIdentityCenterUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIdentityCenterUsersRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID.",
			},

			"user_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User status: Enabled, Disabled.",
			},

			"user_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User type. Manual: manually created; Synchronized: externally imported.",
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter criterion, which currently only supports username, email address, userId, and description.",
			},

			"filter_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filtered user group. IsSelected=1 will be returned for the sub-user associated with this user group.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"sort_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sorting field, which currently only supports CreateTime. The default is the CreateTime field.",
			},

			"sort_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sorting type. Desc: descending order; Asc: ascending order. It should be set along with SortField.",
			},

			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Queried username.",
						},
						"first_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "First name of the user.",
						},
						"last_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last name of the user.",
						},
						"display_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Display name of the user.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User description.",
						},
						"email": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Email address of the user, which must be unique within the directory.",
						},
						"user_status": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User status. Valid values: Enabled, Disabled.",
						},
						"user_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User type. Manual: manually created; Synchronized: externally imported.",
						},
						"user_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User ID.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Creation time of the user.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Modification time of the user.",
						},
						"is_selected": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether selected.",
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

func dataSourceTencentCloudIdentityCenterUsersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_identity_center_users.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_status"); ok {
		paramMap["UserStatus"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_type"); ok {
		paramMap["UserType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		paramMap["Filter"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_groups"); ok {
		filterGroupsList := []*string{}
		filterGroupsSet := v.(*schema.Set).List()
		for i := range filterGroupsSet {
			filterGroups := filterGroupsSet[i].(string)
			filterGroupsList = append(filterGroupsList, helper.String(filterGroups))
		}
		paramMap["FilterGroups"] = filterGroupsList
	}

	if v, ok := d.GetOk("sort_field"); ok {
		paramMap["SortField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_type"); ok {
		paramMap["SortType"] = helper.String(v.(string))
	}

	var users []*organization.UserInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIdentityCenterUsersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		users = result
		return nil
	})
	if err != nil {
		return err
	}

	usersList := make([]map[string]interface{}, 0, len(users))
	ids := make([]string, 0, len(users))
	for _, user := range users {
		usersMap := map[string]interface{}{}

		if user.UserName != nil {
			usersMap["user_name"] = user.UserName
		}

		if user.FirstName != nil {
			usersMap["first_name"] = user.FirstName
		}

		if user.LastName != nil {
			usersMap["last_name"] = user.LastName
		}

		if user.DisplayName != nil {
			usersMap["display_name"] = user.DisplayName
		}

		if user.Description != nil {
			usersMap["description"] = user.Description
		}

		if user.Email != nil {
			usersMap["email"] = user.Email
		}

		if user.UserStatus != nil {
			usersMap["user_status"] = user.UserStatus
		}

		if user.UserType != nil {
			usersMap["user_type"] = user.UserType
		}

		if user.UserId != nil {
			usersMap["user_id"] = user.UserId
			ids = append(ids, *user.UserId)
		}

		if user.CreateTime != nil {
			usersMap["create_time"] = user.CreateTime
		}

		if user.UpdateTime != nil {
			usersMap["update_time"] = user.UpdateTime
		}

		if user.IsSelected != nil {
			usersMap["is_selected"] = user.IsSelected
		}

		usersList = append(usersList, usersMap)

		_ = d.Set("users", usersList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), usersList); e != nil {
			return e
		}
	}

	return nil
}
