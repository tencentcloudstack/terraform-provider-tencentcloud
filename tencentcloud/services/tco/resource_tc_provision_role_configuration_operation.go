package tco

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudProvisionRoleConfigurationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudProvisionRoleConfigurationOperationCreate,
		Read:   resourceTencentCloudProvisionRoleConfigurationOperationRead,
		Delete: resourceTencentCloudProvisionRoleConfigurationOperationDelete,
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

			"target_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of the synchronized target account of the Tencent Cloud Organization. ManagerUin: admin account; MemberUin: member account.",
			},

			"target_uin": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "UIN of the target account of the Tencent Cloud Organization.",
			},
		},
	}
}

func resourceTencentCloudProvisionRoleConfigurationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_provision_role_configuration_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId               = tccommon.GetLogId(tccommon.ContextNil)
		ctx                 = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request             = organization.NewProvisionRoleConfigurationRequest()
		response            = organization.NewProvisionRoleConfigurationResponse()
		zoneId              string
		roleConfigurationId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	if v, ok := d.GetOk("role_configuration_id"); ok {
		roleConfigurationId = v.(string)
		request.RoleConfigurationId = helper.String(roleConfigurationId)
	}

	if v, ok := d.GetOk("target_type"); ok {
		request.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("target_uin"); ok {
		request.TargetUin = helper.IntInt64(v.(int))
	}

	// skipped indicates that the provisioning is not needed (already in the
	// desired state), in which case the operation ends successfully.
	var skipped bool
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().ProvisionRoleConfigurationWithContext(ctx, request)
		if e != nil {
			// The role configuration is already deployed to the target and does
			// not need re-deployment; treat it as success and end the operation.
			var sdkErr *sdkErrors.TencentCloudSDKError
			if errors.As(e, &sdkErr) &&
				sdkErr.GetCode() == "FailedOperation.RoleConfigurationProvisioningStatusError" &&
				sdkErr.GetMessage() == "role configuration provisioning not need re-deploy" {
				skipped = true
				return nil
			}

			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create provision role configuration operation failed, response is nil"))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create provision role configuration operation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{zoneId, roleConfigurationId}, tccommon.FILED_SP))

	if skipped {
		return nil
	}

	if response.Response.Task == nil || response.Response.Task.TaskId == nil {
		return fmt.Errorf("TaskId is nil")
	}

	if _, err := (&resource.StateChangeConf{
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{},
		Refresh:    resourceProvisionRoleConfigurationOperationCreateStateRefreshFunc_0_0(ctx, zoneId, *response.Response.Task.TaskId),
		Target:     []string{"Success"},
		Timeout:    600 * time.Second,
	}).WaitForStateContext(ctx); err != nil {
		return err
	}

	return resourceTencentCloudProvisionRoleConfigurationOperationRead(d, meta)
}

func resourceTencentCloudProvisionRoleConfigurationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_provision_role_configuration_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudProvisionRoleConfigurationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_provision_role_configuration_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceProvisionRoleConfigurationOperationCreateStateRefreshFunc_0_0(ctx context.Context, zoneId string, taskId string) resource.StateRefreshFunc {
	var req *organization.GetTaskStatusRequest
	return func() (interface{}, string, error) {
		meta := tccommon.ProviderMetaFromContext(ctx)
		if meta == nil {
			return nil, "", fmt.Errorf("resource data can not be nil")
		}
		if req == nil {
			d := tccommon.ResourceDataFromContext(ctx)
			if d == nil {
				return nil, "", fmt.Errorf("resource data can not be nil")
			}
			_ = d
			req = organization.NewGetTaskStatusRequest()
			req.ZoneId = helper.String(zoneId)
			req.TaskId = helper.String(taskId)

		}
		resp, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().GetTaskStatusWithContext(ctx, req)
		if err != nil {
			return nil, "", err
		}
		if resp == nil || resp.Response == nil {
			return nil, "", nil
		}
		state := fmt.Sprintf("%v", *resp.Response.TaskStatus.Status)
		return resp.Response, state, nil
	}
}
