/*
Use this data source to query detailed information of organization members

Example Usage

```hcl
data "tencentcloud_organization_members" "members" {}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudOrganizationMembers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationMembersRead,
		Schema: map[string]*schema.Schema{
			"lang": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Valid values: `en` (Tencent Cloud International); `zh` (Tencent Cloud).",
			},

			"search_key": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Search by member name or ID.",
			},

			"auth_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Entity name.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Abbreviation of the trusted service, which is required during querying the trusted service admin.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Member list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"member_uin": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Member UINNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member nameNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"member_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member type. Valid values: `Invite` (invited); `Create` (created).Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"org_policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Relationship policy typeNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"org_policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Relationship policy nameNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"org_permission": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Relationship policy permissionNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Permission ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Permission name.",
									},
								},
							},
						},
						"node_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node IDNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node nameNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "RemarksNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timeNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update timeNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"is_allow_quit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the member is allowed to leave. Valid values: `Allow`, `Denied`.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"pay_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payer UINNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"pay_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Payer nameNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"org_identity": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Management identityNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"identity_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Identity ID.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"identity_alias_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Identity name.Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"bind_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Security information binding status. Valid values: `Unbound`, `Valid`, `Success`, `Failed`.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"permission_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Member permission status. Valid values: `Confirmed`, `UnConfirmed`.Note: This field may return null, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudOrganizationMembersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_organization_members.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("lang"); ok {
		paramMap["Lang"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("auth_name"); ok {
		paramMap["AuthName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*organization.OrgMember

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOrganizationMembersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, orgMember := range items {
			orgMemberMap := map[string]interface{}{}

			if orgMember.MemberUin != nil {
				orgMemberMap["member_uin"] = orgMember.MemberUin
			}

			if orgMember.Name != nil {
				orgMemberMap["name"] = orgMember.Name
			}

			if orgMember.MemberType != nil {
				orgMemberMap["member_type"] = orgMember.MemberType
			}

			if orgMember.OrgPolicyType != nil {
				orgMemberMap["org_policy_type"] = orgMember.OrgPolicyType
			}

			if orgMember.OrgPolicyName != nil {
				orgMemberMap["org_policy_name"] = orgMember.OrgPolicyName
			}

			if orgMember.OrgPermission != nil {
				orgPermissionList := []interface{}{}
				for _, orgPermission := range orgMember.OrgPermission {
					orgPermissionMap := map[string]interface{}{}

					if orgPermission.Id != nil {
						orgPermissionMap["id"] = orgPermission.Id
					}

					if orgPermission.Name != nil {
						orgPermissionMap["name"] = orgPermission.Name
					}

					orgPermissionList = append(orgPermissionList, orgPermissionMap)
				}

				orgMemberMap["org_permission"] = orgPermissionList
			}

			if orgMember.NodeId != nil {
				orgMemberMap["node_id"] = orgMember.NodeId
			}

			if orgMember.NodeName != nil {
				orgMemberMap["node_name"] = orgMember.NodeName
			}

			if orgMember.Remark != nil {
				orgMemberMap["remark"] = orgMember.Remark
			}

			if orgMember.CreateTime != nil {
				orgMemberMap["create_time"] = orgMember.CreateTime
			}

			if orgMember.UpdateTime != nil {
				orgMemberMap["update_time"] = orgMember.UpdateTime
			}

			if orgMember.IsAllowQuit != nil {
				orgMemberMap["is_allow_quit"] = orgMember.IsAllowQuit
			}

			if orgMember.PayUin != nil {
				orgMemberMap["pay_uin"] = orgMember.PayUin
			}

			if orgMember.PayName != nil {
				orgMemberMap["pay_name"] = orgMember.PayName
			}

			if orgMember.OrgIdentity != nil {
				orgIdentityList := []interface{}{}
				for _, orgIdentity := range orgMember.OrgIdentity {
					orgIdentityMap := map[string]interface{}{}

					if orgIdentity.IdentityId != nil {
						orgIdentityMap["identity_id"] = orgIdentity.IdentityId
					}

					if orgIdentity.IdentityAliasName != nil {
						orgIdentityMap["identity_alias_name"] = orgIdentity.IdentityAliasName
					}

					orgIdentityList = append(orgIdentityList, orgIdentityMap)
				}

				orgMemberMap["org_identity"] = orgIdentityList
			}

			if orgMember.BindStatus != nil {
				orgMemberMap["bind_status"] = orgMember.BindStatus
			}

			if orgMember.PermissionStatus != nil {
				orgMemberMap["permission_status"] = orgMember.PermissionStatus
			}

			ids = append(ids, *orgMember.Name)
			tmpList = append(tmpList, orgMemberMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
