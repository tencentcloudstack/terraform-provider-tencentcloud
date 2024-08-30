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

func ResourceTencentCloudIdentityCenterUserSyncProvisioning() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterUserSyncProvisioningCreate,
		Read:   resourceTencentCloudIdentityCenterUserSyncProvisioningRead,
		Update: resourceTencentCloudIdentityCenterUserSyncProvisioningUpdate,
		Delete: resourceTencentCloudIdentityCenterUserSyncProvisioningDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Space ID.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description.",
			},
			"principal_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity ID for the CAM user synchronization. Valid values:\nWhen the PrincipalType value is Group, it is the CIC user group ID (g-********).\nWhen the PrincipalType value is User, it is the CIC user ID (u-********).",
			},
			"principal_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Identity type for the CAM user synchronization. Valid values:\n\nUser: indicates that the identity for the CAM user synchronization is a CIC user.\nGroup: indicates that the identity for the CAM user synchronization is a CIC user group.",
			},
			"target_uin": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "UIN of the synchronized target account of the Tencent Cloud Organization.",
			},
			"duplication_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Conflict policy. It indicates the handling policy for existence of a user with the same username when CIC users are synchronized to CAM. Valid values: KeepBoth: Keep both, that is, add the _cic suffix to the CIC user's username and then try to create a CAM user with the username when CIC users are synchronized to CAM and a user with the same username already exists in CAM; TakeOver: Replace, that is, directly replace the existing CAM user with the synchronized CIC user when CIC users are synchronized to CAM and a user with the same username already exists in CAM.",
			},
			"deletion_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Deletion policy. It indicates the handling policy for CAM users already synchronized when the CAM user synchronization is deleted. Valid values: Delete: Delete the CAM users already synchronized from CIC to CAM when the CAM user synchronization is deleted; Keep: Keep the CAM users already synchronized from CIC to CAM when the CAM user synchronization is deleted.",
			},
			"target_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the synchronized target account of the Tencent Cloud Organization. ManagerUin: admin account; MemberUin: member account.",
			},
			"user_provisioning_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User provisioning id.",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Status of CAM user synchronization. Value:\n" +
					"	* Enabled: CAM user synchronization is enabled;\n" +
					"	* Disabled: CAM user synchronization is not enabled.",
			},
			"principal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identity name of the CAM user synchronization. Value: When PrincipalType is Group, the value is the CIC user group name; When PrincipalType takes the value to User, the value is the CIC user name.",
			},
			"target_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Group account The name of the target account..",
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

func resourceTencentCloudIdentityCenterUserSyncProvisioningCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_user_sync_provisioning.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId   string
		request  = organization.NewCreateUserSyncProvisioningRequest()
		response = organization.NewCreateUserSyncProvisioningResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	userSyncProvisioning := organization.UserSyncProvisioning{}
	if v, ok := d.GetOk("description"); ok {
		userSyncProvisioning.Description = helper.String(v.(string))
	}
	if v, ok := d.GetOk("principal_id"); ok {
		userSyncProvisioning.PrincipalId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("principal_type"); ok {
		userSyncProvisioning.PrincipalType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("target_uin"); ok {
		userSyncProvisioning.TargetUin = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("duplication_strategy"); ok {
		userSyncProvisioning.DuplicationStrategy = helper.String(v.(string))
	}
	if v, ok := d.GetOk("deletion_strategy"); ok {
		userSyncProvisioning.DeletionStrategy = helper.String(v.(string))
	}
	if v, ok := d.GetOk("target_type"); ok {
		userSyncProvisioning.TargetType = helper.String(v.(string))
	}

	request.UserSyncProvisionings = []*organization.UserSyncProvisioning{&userSyncProvisioning}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CreateUserSyncProvisioningWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create identity center user sync provisioning failed, reason:%+v", logId, err)
		return err
	}
	if len(response.Response.Tasks) > 0 {
		task := response.Response.Tasks[0]
		taskId := *task.TaskId
		userProvisioningId := *task.UserProvisioningId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			request := organization.NewGetProvisioningTaskStatusRequest()
			request.ZoneId = helper.String(zoneId)
			request.TaskId = helper.String(taskId)
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().GetProvisioningTaskStatus(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}
			if response.Response.TaskStatus != nil {
				status := *response.Response.TaskStatus.Status
				if status == "Failed" {
					return resource.NonRetryableError(fmt.Errorf("task status is %s", status))
				}
				if status != "Success" {
					return resource.RetryableError(fmt.Errorf("task status is %s", status))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create identity center user sync provisioning failed, reason:%+v", logId, err)
			return err
		}
		d.SetId(strings.Join([]string{zoneId, userProvisioningId}, tccommon.FILED_SP))

	}

	return resourceTencentCloudIdentityCenterUserSyncProvisioningRead(d, meta)
}

func resourceTencentCloudIdentityCenterUserSyncProvisioningRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_user_sync_provisioning.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	userProvisioningId := idSplit[1]
	respData, err := service.DescribeIdentityCenterUserSyncProvisioningById(ctx, zoneId, userProvisioningId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_user_sync_provisioning` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	_ = d.Set("zone_id", zoneId)
	_ = d.Set("user_provisioning_id", respData.UserProvisioningId)

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.PrincipalId != nil {
		_ = d.Set("principal_id", respData.PrincipalId)
	}

	if respData.PrincipalName != nil {
		_ = d.Set("principal_name", respData.PrincipalName)
	}

	if respData.PrincipalType != nil {
		_ = d.Set("principal_type", respData.PrincipalType)
	}

	if respData.TargetUin != nil {
		_ = d.Set("target_uin", respData.TargetUin)
	}

	if respData.TargetName != nil {
		_ = d.Set("target_name", respData.TargetName)
	}

	if respData.DuplicationStrategy != nil {
		_ = d.Set("duplication_strategy", respData.DuplicationStrategy)
	}

	if respData.DeletionStrategy != nil {
		_ = d.Set("deletion_strategy", respData.DeletionStrategy)
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.UpdateTime != nil {
		_ = d.Set("update_time", respData.UpdateTime)
	}

	if respData.TargetType != nil {
		_ = d.Set("target_type", respData.TargetType)
	}

	return nil
}

func resourceTencentCloudIdentityCenterUserSyncProvisioningUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_user_sync_provisioning.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"zone_id", "user_sync_provisionings"}
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
	userProvisioningId := idSplit[1]

	needChange := false
	mutableArgs := []string{"description", "duplication_stateful", "deletion_strategy"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := organization.NewUpdateUserSyncProvisioningRequest()

		request.ZoneId = helper.String(zoneId)

		request.UserProvisioningId = helper.String(userProvisioningId)

		if v, ok := d.GetOk("description"); ok {
			request.NewDescription = helper.String(v.(string))
		}

		if v, ok := d.GetOk("duplication_stateful"); ok {
			request.NewDuplicationStateful = helper.String(v.(string))
		}

		if v, ok := d.GetOk("deletion_strategy"); ok {
			request.NewDeletionStrategy = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdateUserSyncProvisioningWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update identity center user sync provisioning failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudIdentityCenterUserSyncProvisioningRead(d, meta)
}

func resourceTencentCloudIdentityCenterUserSyncProvisioningDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_user_sync_provisioning.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	userProvisioningId := idSplit[1]

	var (
		request  = organization.NewDeleteUserSyncProvisioningRequest()
		response = organization.NewDeleteUserSyncProvisioningResponse()
	)

	request.ZoneId = helper.String(zoneId)

	request.UserProvisioningId = helper.String(userProvisioningId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().DeleteUserSyncProvisioningWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete identity center user sync provisioning failed, reason:%+v", logId, err)
		return err
	}

	if response.Response != nil && response.Response.Tasks != nil && response.Response.Tasks.TaskId != nil {
		taskId := *response.Response.Tasks.TaskId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			request := organization.NewGetProvisioningTaskStatusRequest()
			request.ZoneId = helper.String(zoneId)
			request.TaskId = helper.String(taskId)
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().GetProvisioningTaskStatus(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			}
			if response.Response.TaskStatus != nil {
				status := *response.Response.TaskStatus.Status
				if status == "Failed" {
					return resource.NonRetryableError(fmt.Errorf("task status is %s", status))
				}
				if status != "Success" {
					return resource.RetryableError(fmt.Errorf("task status is %s", status))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s delete identity center user sync provisioning failed, reason:%+v", logId, err)
			return err
		}
	}

	return nil
}
