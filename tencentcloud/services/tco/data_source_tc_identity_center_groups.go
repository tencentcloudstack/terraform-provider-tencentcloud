package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIdentityCenterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIdentityCenterGroupsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID.",
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter criterion. Format: <Attribute> <Operator> <Value>, case-insensitive. Currently, <Attribute> supports only GroupName, and <Operator> supports only eq (Equals) and sw (Start With). For example, Filter = \"GroupName sw test\" indicates querying all user groups with names starting with test; Filter = \"GroupName eq testgroup\" indicates querying the user group with the name testgroup.",
			},

			"group_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User group type. Manual: manually created; Synchronized: externally imported.",
			},

			"filter_users": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Filtered user. IsSelected=1 will be returned for the user group associated with this user.",
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

			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User group list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User group name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User group description.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Creation time of the user group.",
						},
						"group_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User group type. Manual: manually created; Synchronized: externally imported.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Modification time of the user group.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User group ID.",
						},
						"member_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of group members.",
						},
						"is_selected": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If the input parameter FilterUsers is provided, return true when the user is in the user group; otherwise, return false.",
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

func dataSourceTencentCloudIdentityCenterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_identity_center_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		paramMap["Filter"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_type"); ok {
		paramMap["GroupType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter_users"); ok {
		filterUsersList := []*string{}
		filterUsersSet := v.(*schema.Set).List()
		for i := range filterUsersSet {
			filterUsers := filterUsersSet[i].(string)
			filterUsersList = append(filterUsersList, helper.String(filterUsers))
		}
		paramMap["FilterUsers"] = filterUsersList
	}

	if v, ok := d.GetOk("sort_field"); ok {
		paramMap["SortField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_type"); ok {
		paramMap["SortType"] = helper.String(v.(string))
	}

	var groups []*organization.GroupInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIdentityCenterGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		groups = result
		return nil
	})
	if err != nil {
		return err
	}

	groupsList := make([]map[string]interface{}, 0, len(groups))
	ids := make([]string, 0, len(groups))
	for _, group := range groups {
		groupsMap := map[string]interface{}{}

		if group.GroupName != nil {
			groupsMap["group_name"] = group.GroupName
		}

		if group.Description != nil {
			groupsMap["description"] = group.Description
		}

		if group.CreateTime != nil {
			groupsMap["create_time"] = group.CreateTime
		}

		if group.GroupType != nil {
			groupsMap["group_type"] = group.GroupType
		}

		if group.UpdateTime != nil {
			groupsMap["update_time"] = group.UpdateTime
		}

		if group.GroupId != nil {
			groupsMap["group_id"] = group.GroupId
			ids = append(ids, *group.GroupId)
		}

		if group.MemberCount != nil {
			groupsMap["member_count"] = group.MemberCount
		}

		if group.IsSelected != nil {
			groupsMap["is_selected"] = group.IsSelected
		}

		groupsList = append(groupsList, groupsMap)
	}

	_ = d.Set("groups", groupsList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), groupsList); e != nil {
			return e
		}
	}

	return nil
}
