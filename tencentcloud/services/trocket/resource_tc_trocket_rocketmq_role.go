package trocket

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	trocket "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/trocket/v20230308"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTrocketRocketmqRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTrocketRocketmqRoleCreate,
		Read:   resourceTencentCloudTrocketRocketmqRoleRead,
		Update: resourceTencentCloudTrocketRocketmqRoleUpdate,
		Delete: resourceTencentCloudTrocketRocketmqRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "ID of instance.",
			},

			"role": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Name of role.",
			},

			"remark": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "remark.",
			},

			"perm_write": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable production permission.",
			},

			"perm_read": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable consumption permission.",
			},

			"access_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Access key.",
			},

			"secret_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Secret key.",
			},

			"created_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Created time.",
			},

			"modified_time": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Modified time.",
			},
		},
	}
}

func resourceTencentCloudTrocketRocketmqRoleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_role.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = trocket.NewCreateRoleRequest()
		response   = trocket.NewCreateRoleResponse()
		instanceId string
		role       string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("role"); ok {
		request.Role = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("perm_write"); ok {
		request.PermWrite = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("perm_read"); ok {
		request.PermRead = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTrocketClient().CreateRole(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create trocket rocketmqRole failed, reason:%+v", logId, err)
		return err
	}

	role = *response.Response.Role
	d.SetId(instanceId + tccommon.FILED_SP + role)

	return resourceTencentCloudTrocketRocketmqRoleRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqRoleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_role.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	role := idSplit[1]

	rocketmqRole, err := service.DescribeTrocketRocketmqRoleById(ctx, instanceId, role)
	if err != nil {
		return err
	}

	if rocketmqRole == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TrocketRocketmqRole` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("role", role)

	if rocketmqRole.Remark != nil {
		_ = d.Set("remark", rocketmqRole.Remark)
	}

	if rocketmqRole.PermWrite != nil {
		_ = d.Set("perm_write", rocketmqRole.PermWrite)
	}

	if rocketmqRole.PermRead != nil {
		_ = d.Set("perm_read", rocketmqRole.PermRead)
	}

	if rocketmqRole.AccessKey != nil {
		_ = d.Set("access_key", rocketmqRole.AccessKey)
	}

	if rocketmqRole.SecretKey != nil {
		_ = d.Set("secret_key", rocketmqRole.SecretKey)
	}

	if rocketmqRole.CreatedTime != nil {
		_ = d.Set("created_time", rocketmqRole.CreatedTime)
	}

	if rocketmqRole.ModifiedTime != nil {
		_ = d.Set("modified_time", rocketmqRole.ModifiedTime)
	}

	return nil
}

func resourceTencentCloudTrocketRocketmqRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_role.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := trocket.NewModifyRoleRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	role := idSplit[1]

	request.InstanceId = &instanceId
	request.Role = &role

	mutableArgs := []string{"remark", "perm_write", "perm_read"}
	needChange := false

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
		}
	}

	if needChange {
		if v, ok := d.GetOk("remark"); ok {
			request.Remark = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("perm_write"); ok {
			request.PermWrite = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("perm_read"); ok {
			request.PermRead = helper.Bool(v.(bool))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTrocketClient().ModifyRole(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update trocket rocketmqRole failed, reason:%+v", logId, err)
			return err
		}
	}
	return resourceTencentCloudTrocketRocketmqRoleRead(d, meta)
}

func resourceTencentCloudTrocketRocketmqRoleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_trocket_rocketmq_role.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TrocketService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	role := idSplit[1]

	if err := service.DeleteTrocketRocketmqRoleById(ctx, instanceId, role); err != nil {
		return err
	}

	return nil
}
