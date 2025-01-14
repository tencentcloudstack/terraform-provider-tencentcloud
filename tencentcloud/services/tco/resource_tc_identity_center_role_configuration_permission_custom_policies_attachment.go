package tco

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachmentCreate,
		Read:   resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachmentRead,
		Delete: resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Space ID.",
			},

			"role_configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Permission configuration ID.",
			},

			"policies": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_policy_document": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Role policy document.",
						},

						"role_policy_name": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Role policy name.",
						},
						"role_policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role policy type.",
						},

						"add_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Role policy add time.",
						},
					},
				},
				Description: "Policies.",
			},
		},
	}
}

func resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId              string
		roleConfigurationId string
		rolePolicyNames     []string
	)
	var (
		request  = organization.NewAddPermissionPolicyToRoleConfigurationRequest()
		response = organization.NewAddPermissionPolicyToRoleConfigurationResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	if v, ok := d.GetOk("role_configuration_id"); ok {
		roleConfigurationId = v.(string)
		request.RoleConfigurationId = helper.String(roleConfigurationId)
	}

	if v, ok := d.GetOk("policies"); ok {
		policies := v.(*schema.Set).List()
		for _, poilcy := range policies {
			policyMap := poilcy.(map[string]interface{})
			rolePolicyName := policyMap["role_policy_name"]
			rolePolicyDocument := policyMap["role_policy_document"]
			rolePolicyNames = append(rolePolicyNames, rolePolicyName.(string))
			request.RolePolicyNames = append(request.RolePolicyNames, helper.String(rolePolicyName.(string)))
			request.CustomPolicyDocuments = append(request.CustomPolicyDocuments, helper.String(rolePolicyDocument.(string)))
		}
	}

	request.RolePolicyType = helper.String("Custom")

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddPermissionPolicyToRoleConfigurationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create identity center role configuration permission policy attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response

	rolePolicyNameStr := strings.Join(rolePolicyNames, tccommon.COMMA_SP)
	d.SetId(strings.Join([]string{zoneId, roleConfigurationId, rolePolicyNameStr}, tccommon.FILED_SP))

	return resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachmentRead(d, meta)
}

func resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	roleConfigurationId := idSplit[1]
	rolePolicyNames := strings.Split(idSplit[2], tccommon.COMMA_SP)
	rolePolicyNameSet := make(map[string]struct{})

	for _, rolePolicyName := range rolePolicyNames {
		rolePolicyNameSet[rolePolicyName] = struct{}{}
	}

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("role_configuration_id", roleConfigurationId)

	respData, err := service.DescribeIdentityCenterRoleConfigurationPermissionPolicyAttachmentById(ctx, zoneId, roleConfigurationId, "Custom")
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_role_configuration_permission_policy_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.RolePolicies != nil {
		policies := make([]interface{}, 0)
		for _, r := range respData.RolePolicies {
			if _, ok := rolePolicyNameSet[*r.RolePolicyName]; ok {
				policyMap := make(map[string]interface{})

				if r.RolePolicyName != nil {
					policyMap["role_policy_name"] = *r.RolePolicyName
				}

				if r.RolePolicyType != nil {
					policyMap["role_policy_type"] = *r.RolePolicyType
				}

				if r.RolePolicyDocument != nil {
					policyMap["role_policy_document"] = *r.RolePolicyDocument
				}

				if r.AddTime != nil {
					policyMap["add_time"] = *r.AddTime
				}

				policies = append(policies, policyMap)
			}
		}
		_ = d.Set("policies", policies)
	}

	return nil
}

func resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPoliciesAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration_permission_custom_policies_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	roleConfigurationId := idSplit[1]
	rolePolicyNames := strings.Split(idSplit[2], tccommon.COMMA_SP)

	var (
		request = organization.NewRemovePermissionPolicyFromRoleConfigurationRequest()
	)

	request.ZoneId = helper.String(zoneId)

	request.RoleConfigurationId = helper.String(roleConfigurationId)

	request.RolePolicyType = helper.String("Custom")

	request.RolePolicyId = helper.Int64(0)

	for _, rolePolicyName := range rolePolicyNames {
		request.RolePolicyName = helper.String(rolePolicyName)

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().RemovePermissionPolicyFromRoleConfigurationWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s delete identity center role configuration permission policy attachment failed, reason:%+v", logId, err)
			return err
		}
	}

	return nil
}
