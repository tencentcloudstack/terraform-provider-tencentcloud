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

func ResourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPolicyAttachmentCreate,
		Read:   resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPolicyAttachmentRead,
		Delete: resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPolicyAttachmentDelete,
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
	}
}

func resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId              string
		roleConfigurationId string
		rolePolicyName      string
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

	if v, ok := d.GetOk("role_policy_name"); ok {
		rolePolicyName = v.(string)
		request.RolePolicyNames = []*string{helper.String(rolePolicyName)}
	}

	if v, ok := d.GetOk("role_policy_document"); ok {
		request.CustomPolicyDocument = helper.String(v.(string))
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

	d.SetId(strings.Join([]string{zoneId, roleConfigurationId, rolePolicyName}, tccommon.FILED_SP))

	return resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPolicyAttachmentRead(d, meta)
}

func resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment.read")()
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
	rolePolicyName := idSplit[2]

	_ = d.Set("zone_id", zoneId)

	_ = d.Set("role_configuration_id", roleConfigurationId)

	_ = d.Set("role_policy_name", rolePolicyName)

	var respData *organization.ListPermissionPoliciesInRoleConfigurationResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeIdentityCenterRoleConfigurationPermissionPolicyAttachmentById(ctx, zoneId, roleConfigurationId, "Custom")
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_role_configuration_permission_policy_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if respData.RolePolicies != nil {
		var rolePolicie *organization.RolePolicie
		for _, r := range respData.RolePolicies {
			if *r.RolePolicyName == rolePolicyName {
				rolePolicie = r
				break
			}
		}

		if rolePolicie == nil {
			log.Printf("[WARN]%s resource `identity_center_role_configuration_permission_policy_attachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
			return fmt.Errorf("RolePolicy %s is not exist", d.Id())
		}

		if rolePolicie.RolePolicyName != nil {
			_ = d.Set("role_policy_name", rolePolicie.RolePolicyName)
		}

		if rolePolicie.RolePolicyType != nil {
			_ = d.Set("role_policy_type", rolePolicie.RolePolicyType)
		}

		if rolePolicie.RolePolicyDocument != nil {
			_ = d.Set("role_policy_document", rolePolicie.RolePolicyDocument)
		}

		if rolePolicie.AddTime != nil {
			_ = d.Set("add_time", rolePolicie.AddTime)
		}

	}

	return nil
}

func resourceTencentCloudIdentityCenterRoleConfigurationPermissionCustomPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration_permission_custom_policy_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	roleConfigurationId := idSplit[1]
	rolePolicyName := idSplit[2]

	var (
		request  = organization.NewRemovePermissionPolicyFromRoleConfigurationRequest()
		response = organization.NewRemovePermissionPolicyFromRoleConfigurationResponse()
	)

	request.ZoneId = helper.String(zoneId)

	request.RoleConfigurationId = helper.String(roleConfigurationId)

	request.RolePolicyType = helper.String("Custom")

	request.RolePolicyId = helper.Int64(0)

	request.RolePolicyName = helper.String(rolePolicyName)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().RemovePermissionPolicyFromRoleConfigurationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete identity center role configuration permission policy attachment failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
