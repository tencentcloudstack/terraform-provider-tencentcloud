package trocket

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmqRocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdmqRocketmqEnvironmentRole() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTdmqRocketmqEnvironmentRoleRead,
		Create: resourceTencentCloudTdmqRocketmqEnvironmentRoleCreate,
		Update: resourceTencentCloudTdmqRocketmqEnvironmentRoleUpdate,
		Delete: resourceTencentCloudTdmqRocketmqEnvironmentRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Environment (namespace) name.",
			},

			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Role Name.",
			},

			"permissions": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Permissions, which is a non-empty string array of `produce` and `consume` at the most.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID (required).",
			},
		},
	}
}

func resourceTencentCloudTdmqRocketmqEnvironmentRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmqRocketmq_environment_role.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request         = tdmqRocketmq.NewCreateEnvironmentRoleRequest()
		clusterId       string
		roleName        string
		environmentName string
	)

	if v, ok := d.GetOk("environment_name"); ok {
		environmentName = v.(string)
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role_name"); ok {
		roleName = v.(string)
		request.RoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("permissions"); ok {
		permissionsSet := v.(*schema.Set).List()
		for i := range permissionsSet {
			permissions := permissionsSet[i].(string)
			request.Permissions = append(request.Permissions, &permissions)
		}
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateEnvironmentRole(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq environmentRole failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId + tccommon.FILED_SP + roleName + tccommon.FILED_SP + environmentName)
	return resourceTencentCloudTdmqRocketmqEnvironmentRoleRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqEnvironmentRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmqRocketmq_environment_role.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TdmqRocketmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]
	environmentName := idSplit[2]

	environmentRoles, err := service.DescribeTdmqRocketmqEnvironmentRole(ctx, clusterId, roleName, environmentName)

	if err != nil {
		return err
	}

	if len(environmentRoles) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `environmentRole` %s does not exist", roleName)
	}
	environmentRole := environmentRoles[0]
	_ = d.Set("environment_name", environmentRole.EnvironmentId)
	_ = d.Set("role_name", environmentRole.RoleName)
	permissions := make([]string, 0)
	for _, i := range environmentRole.Permissions {
		permissions = append(permissions, *i)
	}
	_ = d.Set("permissions", permissions)
	_ = d.Set("cluster_id", clusterId)

	return nil
}

func resourceTencentCloudTdmqRocketmqEnvironmentRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmqRocketmq_environment_role.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := tdmqRocketmq.NewModifyEnvironmentRoleRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]
	environmentId := idSplit[2]

	request.ClusterId = &clusterId
	request.RoleName = &roleName
	request.EnvironmentId = &environmentId

	if d.HasChange("permissions") {
		if v, ok := d.GetOk("permissions"); ok {
			permissionsSet := v.(*schema.Set).List()
			for i := range permissionsSet {
				permissions := permissionsSet[i].(string)
				request.Permissions = append(request.Permissions, &permissions)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyEnvironmentRole(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmqRocketmq environmentRole failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTdmqRocketmqEnvironmentRoleRead(d, meta)
}

func resourceTencentCloudTdmqRocketmqEnvironmentRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmqRocketmq_environment_role.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TdmqRocketmqService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	roleName := idSplit[1]
	environmentId := idSplit[2]

	if err := service.DeleteTdmqRocketmqEnvironmentRoleById(ctx, clusterId, roleName, environmentId); err != nil {
		return err
	}

	return nil
}
