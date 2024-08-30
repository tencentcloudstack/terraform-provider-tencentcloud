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

func ResourceTencentCloudIdentityCenterRoleConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterRoleConfigurationCreate,
		Read:   resourceTencentCloudIdentityCenterRoleConfigurationRead,
		Update: resourceTencentCloudIdentityCenterRoleConfigurationUpdate,
		Delete: resourceTencentCloudIdentityCenterRoleConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID.",
			},

			"role_configuration_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access configuration name, which contains up to 128 characters, including English letters, digits, and hyphens (-).",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Access configuration description, which contains up to 1024 characters.",
			},

			"session_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Session duration. It indicates the maximum session duration when CIC users use the access configuration to access the target account of the Tencent Cloud Organization. Unit: seconds. Value range: 900-43,200 (15 minutes to 12 hours). Default value: 3600 (1 hour).",
			},

			"relay_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Initial access page. It indicates the initial access page URL when CIC users use the access configuration to access the target account of the Tencent Cloud Organization. This page must be the Tencent Cloud console page. The default is null, which indicates navigating to the home page of the Tencent Cloud console.",
			},
			"role_configuration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Role configuration id.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Update time.",
			},
		},
	}
}

func resourceTencentCloudIdentityCenterRoleConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId              string
		roleConfigurationId string
	)
	var (
		request  = organization.NewCreateRoleConfigurationRequest()
		response = organization.NewCreateRoleConfigurationResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role_configuration_name"); ok {
		request.RoleConfigurationName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("session_duration"); ok {
		request.SessionDuration = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("relay_state"); ok {
		request.RelayState = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CreateRoleConfigurationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create identity center role configuration failed, reason:%+v", logId, err)
		return err
	}

	roleConfigurationId = *response.Response.RoleConfigurationInfo.RoleConfigurationId

	d.SetId(strings.Join([]string{zoneId, roleConfigurationId}, tccommon.FILED_SP))

	return resourceTencentCloudIdentityCenterRoleConfigurationRead(d, meta)
}

func resourceTencentCloudIdentityCenterRoleConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	roleConfigurationId := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeIdentityCenterRoleConfigurationById(ctx, zoneId, roleConfigurationId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_role_configuration` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.RoleConfigurationId != nil {
		_ = d.Set("role_configuration_id", respData.RoleConfigurationId)
	}

	if respData.RoleConfigurationName != nil {
		_ = d.Set("role_configuration_name", respData.RoleConfigurationName)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.SessionDuration != nil {
		_ = d.Set("session_duration", respData.SessionDuration)
	}

	if respData.RelayState != nil {
		_ = d.Set("relay_state", respData.RelayState)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	return nil
}

func resourceTencentCloudIdentityCenterRoleConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"zone_id", "role_configuration_name"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	roleConfigurationId := idSplit[1]

	needChange := false
	mutableArgs := []string{"description", "session_duration", "relay_state"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := organization.NewUpdateRoleConfigurationRequest()

		request.ZoneId = helper.String(zoneId)

		request.RoleConfigurationId = helper.String(roleConfigurationId)

		if v, ok := d.GetOk("description"); ok {
			request.NewDescription = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("session_duration"); ok {
			request.NewSessionDuration = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("relay_state"); ok {
			request.NewRelayState = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdateRoleConfigurationWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update identity center role configuration failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudIdentityCenterRoleConfigurationRead(d, meta)
}

func resourceTencentCloudIdentityCenterRoleConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_configuration.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	roleConfigurationId := idSplit[1]

	var (
		request  = organization.NewDeleteRoleConfigurationRequest()
		response = organization.NewDeleteRoleConfigurationResponse()
	)

	request.ZoneId = helper.String(zoneId)

	request.RoleConfigurationId = helper.String(roleConfigurationId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().DeleteRoleConfigurationWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete identity center role configuration failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
