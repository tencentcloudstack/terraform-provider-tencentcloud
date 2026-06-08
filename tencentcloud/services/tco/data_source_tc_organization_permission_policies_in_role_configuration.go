package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOrganizationPermissionPoliciesInRoleConfiguration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOrganizationPermissionPoliciesInRoleConfigurationRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID.",
			},

			"role_configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Role configuration ID.",
			},

			"role_policy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Permission policy type. Valid values: `System`: System policy, reuses CAM system policies. `Custom`: Custom policy, written according to CAM permission policy syntax and structure.",
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search by policy name.",
			},

			"total_counts": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of permission policies.",
			},

			"role_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Permission policy list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_policy_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Policy ID.",
						},
						"role_policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Permission policy name.",
						},
						"role_policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Permission policy type.",
						},
						"role_policy_document": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Custom policy content. Only returned for custom policies.",
						},
						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time when the permission policy was added to the role configuration.",
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

func dataSourceTencentCloudOrganizationPermissionPoliciesInRoleConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_organization_permission_policies_in_role_configuration.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		zoneId              string
		roleConfigurationId string
	)

	zoneId = d.Get("zone_id").(string)
	roleConfigurationId = d.Get("role_configuration_id").(string)

	request := organization.NewListPermissionPoliciesInRoleConfigurationRequest()
	request.ZoneId = helper.String(zoneId)
	request.RoleConfigurationId = helper.String(roleConfigurationId)

	if v, ok := d.GetOk("role_policy_type"); ok {
		request.RolePolicyType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		request.Filter = helper.String(v.(string))
	}

	var response *organization.ListPermissionPoliciesInRoleConfigurationResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient()
		result, e := client.ListPermissionPoliciesInRoleConfigurationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		response = result
		return nil
	})

	if err != nil {
		return err
	}

	if response == nil || response.Response == nil {
		d.SetId("")
		return nil
	}

	if response.Response.TotalCounts != nil {
		_ = d.Set("total_counts", *response.Response.TotalCounts)
	}

	rolePoliciesList := make([]map[string]interface{}, 0)
	if response.Response.RolePolicies != nil {
		for _, rolePolicy := range response.Response.RolePolicies {
			rolePolicyMap := map[string]interface{}{}

			if rolePolicy.RolePolicyId != nil {
				rolePolicyMap["role_policy_id"] = *rolePolicy.RolePolicyId
			}

			if rolePolicy.RolePolicyName != nil {
				rolePolicyMap["role_policy_name"] = *rolePolicy.RolePolicyName
			}

			if rolePolicy.RolePolicyType != nil {
				rolePolicyMap["role_policy_type"] = *rolePolicy.RolePolicyType
			}

			if rolePolicy.RolePolicyDocument != nil {
				rolePolicyMap["role_policy_document"] = *rolePolicy.RolePolicyDocument
			}

			if rolePolicy.AddTime != nil {
				rolePolicyMap["add_time"] = *rolePolicy.AddTime
			}

			rolePoliciesList = append(rolePoliciesList, rolePolicyMap)
		}
	}

	_ = d.Set("role_policies", rolePoliciesList)

	d.SetId(zoneId + tccommon.FILED_SP + roleConfigurationId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), rolePoliciesList); e != nil {
			return e
		}
	}

	return nil
}
