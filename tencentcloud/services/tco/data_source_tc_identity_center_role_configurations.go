package tco

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIdentityCenterRoleConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIdentityCenterRoleConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID.",
			},

			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter criteria, which are case insensitive. Currently, only RoleConfigurationName is supported and only eq (Equals) and sw (Start With) are supported. Example: Filter = \"RoleConfigurationName, only sw test\" means querying all permission configurations starting with test. Filter = \"RoleConfigurationName, only eq TestRoleConfiguration\" means querying the permission configuration named TestRoleConfiguration.",
			},

			"filter_targets": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Check whether the member account has been configured with permissions. If configured, return IsSelected: true; otherwise, return false.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"principal_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "UserId of the authorized user or GroupId of the authorized user group, which must be set together with the input parameter FilterTargets.",
			},

			"role_configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Permission configuration list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_configuration_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Permission configuration ID.",
						},
						"role_configuration_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Permission configuration name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Permission configuration description.",
						},
						"session_duration": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Session duration. It indicates the maximum session duration when CIC users use the access configuration to access member accounts.\nUnit: seconds.",
						},
						"relay_state": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Initial access page. It indicates the initial access page URL when CIC users use the access configuration to access member accounts.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Creation time of the permission configuration.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Update time of the permission configuration.",
						},
						"is_selected": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If the input parameter FilterTargets is provided, check whether the member account has been configured with permissions. If configured, return true; otherwise, return false.",
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

func dataSourceTencentCloudIdentityCenterRoleConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_identity_center_role_configurations.read")()
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

	if v, ok := d.GetOk("filter_targets"); ok {
		filterTargetsList := []*int64{}
		filterTargetsSet := v.(*schema.Set).List()
		for i := range filterTargetsSet {
			filterTargets := filterTargetsSet[i].(int)
			filterTargetsList = append(filterTargetsList, helper.IntInt64(filterTargets))
		}
		paramMap["FilterTargets"] = filterTargetsList
	}

	if v, ok := d.GetOk("principal_id"); ok {
		paramMap["PrincipalId"] = helper.String(v.(string))
	}

	var roleConfigurations []*organization.RoleConfiguration
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIdentityCenterRoleConfigurationsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		roleConfigurations = result
		return nil
	})
	if err != nil {
		return err
	}

	roleConfigurationsList := make([]map[string]interface{}, 0, len(roleConfigurations))
	ids := make([]string, 0, len(roleConfigurations))
	for _, roleConfiguration := range roleConfigurations {
		roleConfigurationsMap := map[string]interface{}{}

		if roleConfiguration.RoleConfigurationId != nil {
			roleConfigurationsMap["role_configuration_id"] = roleConfiguration.RoleConfigurationId
			ids = append(ids, *roleConfiguration.RoleConfigurationId)
		}

		if roleConfiguration.RoleConfigurationName != nil {
			roleConfigurationsMap["role_configuration_name"] = roleConfiguration.RoleConfigurationName
		}

		if roleConfiguration.Description != nil {
			roleConfigurationsMap["description"] = roleConfiguration.Description
		}

		if roleConfiguration.SessionDuration != nil {
			roleConfigurationsMap["session_duration"] = roleConfiguration.SessionDuration
		}

		if roleConfiguration.RelayState != nil {
			roleConfigurationsMap["relay_state"] = roleConfiguration.RelayState
		}

		if roleConfiguration.CreateTime != nil {
			roleConfigurationsMap["create_time"] = roleConfiguration.CreateTime
		}

		if roleConfiguration.UpdateTime != nil {
			roleConfigurationsMap["update_time"] = roleConfiguration.UpdateTime
		}

		if roleConfiguration.IsSelected != nil {
			roleConfigurationsMap["is_selected"] = roleConfiguration.IsSelected
		}

		roleConfigurationsList = append(roleConfigurationsList, roleConfigurationsMap)
	}

	_ = d.Set("role_configurations", roleConfigurationsList)

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), roleConfigurationsList); e != nil {
			return e
		}
	}

	return nil
}
