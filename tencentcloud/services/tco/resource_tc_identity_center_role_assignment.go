package tco

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIdentityCenterRoleAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIdentityCenterRoleAssignmentCreate,
		Read:   resourceTencentCloudIdentityCenterRoleAssignmentRead,
		Delete: resourceTencentCloudIdentityCenterRoleAssignmentDelete,
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
			"principal_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity ID for the CAM user synchronization. Valid values:\nWhen the PrincipalType value is Group, it is the CIC user group ID (g-********).\nWhen the PrincipalType value is User, it is the CIC user ID (u-********).",
			},
			"principal_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Identity type for the CAM user synchronization. Valid values:\n\nUser: indicates that the identity for the CAM user synchronization is a CIC user.\nGroup: indicates that the identity for the CAM user synchronization is a CIC user group.",
			},
			"target_uin": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "UIN of the synchronized target account of the Tencent Cloud Organization.",
			},
			"target_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of the synchronized target account of the Tencent Cloud Organization. ManagerUin: admin account; MemberUin: member account.",
			},
			"role_configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Permission configuration ID.",
			},
			"deprovision_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "None",
				ForceNew:    true,
				Description: "When you remove the last authorization configured with a certain privilege on a group account target account, whether to cancel the privilege configuration deployment at the same time. Value: DeprovisionForLastRoleAssignmentOnAccount: Remove privileges to configure deployment. None (default): Configure deployment without delegating privileges.",
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
			"role_configuration_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Role configuration name.",
			},
			"target_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Target name.",
			},
			"principal_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Principal name.",
			},
		},
	}
}

func resourceTencentCloudIdentityCenterRoleAssignmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_assignment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var (
		zoneId              string
		roleConfigurationId string
		targetType          string
		targetUin           int64
		principalType       string
		principalId         string
	)
	var (
		request  = organization.NewCreateRoleAssignmentRequest()
		response = organization.NewCreateRoleAssignmentResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(zoneId)
	}

	roleAssignmentInfo := organization.RoleAssignmentInfo{}
	if v, ok := d.GetOk("principal_id"); ok {
		principalId = v.(string)
		roleAssignmentInfo.PrincipalId = helper.String(principalId)
	}
	if v, ok := d.GetOk("principal_type"); ok {
		principalType = v.(string)
		roleAssignmentInfo.PrincipalType = helper.String(principalType)
	}
	if v, ok := d.GetOk("target_uin"); ok {
		targetUin = int64(v.(int))
		roleAssignmentInfo.TargetUin = helper.Int64(targetUin)
	}
	if v, ok := d.GetOk("target_type"); ok {
		targetType = v.(string)
		roleAssignmentInfo.TargetType = helper.String(targetType)
	}
	if v, ok := d.GetOk("role_configuration_id"); ok {
		roleConfigurationId = v.(string)
		roleAssignmentInfo.RoleConfigurationId = helper.String(roleConfigurationId)
	}
	request.RoleAssignmentInfo = []*organization.RoleAssignmentInfo{&roleAssignmentInfo}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CreateRoleAssignmentWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create identity center role assignment failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.Tasks) > 0 {
		task := response.Response.Tasks[0]
		taskId := *task.TaskId
		roleConfigurationId := *task.RoleConfigurationId
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"Success"}, 2*tccommon.ReadRetryTimeout, time.Second, service.AssignmentTaskStatusStateRefreshFunc(zoneId, taskId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}

		targetUinString := strconv.FormatInt(targetUin, 10)
		d.SetId(strings.Join([]string{zoneId, roleConfigurationId, targetType, targetUinString, principalType, principalId}, tccommon.FILED_SP))
	}

	return resourceTencentCloudIdentityCenterRoleAssignmentRead(d, meta)
}

func resourceTencentCloudIdentityCenterRoleAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_assignment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	respData, err := service.DescribeIdentityCenterRoleAssignmentById(ctx, d.Id())
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `identity_center_role_assignment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if len(respData.RoleAssignments) > 0 {
		roleAssignment := respData.RoleAssignments[0]
		if roleAssignment.RoleConfigurationId != nil {
			_ = d.Set("role_configuration_id", roleAssignment.RoleConfigurationId)
		}
		if roleAssignment.RoleConfigurationName != nil {
			_ = d.Set("role_configuration_name", roleAssignment.RoleConfigurationName)
		}
		if roleAssignment.TargetUin != nil {
			_ = d.Set("target_uin", roleAssignment.TargetUin)
		}
		if roleAssignment.TargetType != nil {
			_ = d.Set("target_type", roleAssignment.TargetType)
		}
		if roleAssignment.PrincipalId != nil {
			_ = d.Set("principal_id", roleAssignment.PrincipalId)
		}
		if roleAssignment.PrincipalType != nil {
			_ = d.Set("principal_type", roleAssignment.PrincipalType)
		}
		if roleAssignment.PrincipalName != nil {
			_ = d.Set("principal_name", roleAssignment.PrincipalName)
		}
		if roleAssignment.TargetName != nil {
			_ = d.Set("target_name", roleAssignment.TargetName)
		}
		if roleAssignment.CreateTime != nil {
			_ = d.Set("create_time", roleAssignment.CreateTime)
		}
		if roleAssignment.UpdateTime != nil {
			_ = d.Set("update_time", roleAssignment.UpdateTime)
		}

	} else {
		d.SetId("")
	}

	return nil
}

func resourceTencentCloudIdentityCenterRoleAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_identity_center_role_assignment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 6 {
		return fmt.Errorf("roleAssignmentId is broken,%s", d.Id())
	}

	zoneId := idSplit[0]
	roleConfigurationId := idSplit[1]
	targetType := idSplit[2]
	targetUinString := idSplit[3]
	principalType := idSplit[4]
	principalId := idSplit[5]

	var (
		deleteRoleAssignmentRequest        = organization.NewDeleteRoleAssignmentRequest()
		deleteRoleAssignmentResponse       = organization.NewDeleteRoleAssignmentResponse()
		dismantleRoleConfigurationRequest  = organization.NewDismantleRoleConfigurationRequest()
		dismantleRoleConfigurationResponse = organization.NewDismantleRoleConfigurationResponse()
	)
	deleteRoleAssignmentRequest.ZoneId = helper.String(zoneId)
	deleteRoleAssignmentRequest.RoleConfigurationId = helper.String(roleConfigurationId)
	deleteRoleAssignmentRequest.TargetType = helper.String(targetType)
	targetUin, err := strconv.ParseInt(targetUinString, 10, 64)
	if err != nil {
		return err
	}
	deleteRoleAssignmentRequest.TargetUin = helper.Int64(targetUin)
	deleteRoleAssignmentRequest.PrincipalType = helper.String(principalType)
	deleteRoleAssignmentRequest.PrincipalId = helper.String(principalId)
	if v, ok := d.GetOk("deprovision_strategy"); ok {
		deleteRoleAssignmentRequest.DeprovisionStrategy = helper.String(v.(string))
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().DeleteRoleAssignment(deleteRoleAssignmentRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, deleteRoleAssignmentRequest.GetAction(), deleteRoleAssignmentRequest.ToJsonString(), result.ToJsonString())
		}
		deleteRoleAssignmentResponse = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete identity center role assignment failed, reason:%+v", logId, err)
		return err
	}

	if deleteRoleAssignmentResponse.Response != nil && deleteRoleAssignmentResponse.Response.Task != nil && deleteRoleAssignmentResponse.Response.Task.TaskId != nil {
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"Success"}, 2*tccommon.ReadRetryTimeout, time.Second, service.AssignmentTaskStatusStateRefreshFunc(zoneId, *deleteRoleAssignmentResponse.Response.Task.TaskId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	dismantleRoleConfigurationRequest.RoleConfigurationId = helper.String(roleConfigurationId)
	dismantleRoleConfigurationRequest.ZoneId = helper.String(zoneId)
	dismantleRoleConfigurationRequest.TargetType = helper.String(targetType)
	dismantleRoleConfigurationRequest.TargetUin = helper.Int64(targetUin)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().DismantleRoleConfiguration(dismantleRoleConfigurationRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, dismantleRoleConfigurationRequest.GetAction(), dismantleRoleConfigurationRequest.ToJsonString(), result.ToJsonString())
		}
		dismantleRoleConfigurationResponse = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete identity center role assignment failed, reason:%+v", logId, err)
		return err
	}

	if dismantleRoleConfigurationResponse.Response != nil && dismantleRoleConfigurationResponse.Response.Task != nil && dismantleRoleConfigurationResponse.Response.Task.TaskId != nil {
		conf := tccommon.BuildStateChangeConf([]string{}, []string{"Success"}, 2*tccommon.ReadRetryTimeout, time.Second, service.AssignmentTaskStatusStateRefreshFunc(zoneId, *dismantleRoleConfigurationResponse.Response.Task.TaskId, []string{}))
		if _, e := conf.WaitForState(); e != nil {
			return e
		}
	}

	return nil
}
