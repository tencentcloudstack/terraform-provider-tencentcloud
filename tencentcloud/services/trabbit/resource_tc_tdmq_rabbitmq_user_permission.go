package trabbit

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTdmqRabbitmqUserPermission() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRabbitmqUserPermissionCreate,
		Read:   resourceTencentCloudTdmqRabbitmqUserPermissionRead,
		Update: resourceTencentCloudTdmqRabbitmqUserPermissionUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqUserPermissionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Cluster instance ID.",
			},
			"user": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Username.",
			},
			"virtual_host": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "VirtualHost name.",
			},
			"config_regexp": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Configure permission regexp, controls which resources can be declared.",
			},
			"write_regexp": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Write permission regexp, controls which resources can be written.",
			},
			"read_regexp": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Read permission regexp, controls which resources can be read.",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqUserPermissionCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_user_permission.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = tdmq.NewModifyRabbitMQPermissionRequest()
		instanceId  string
		user        string
		virtualHost string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("user"); ok {
		request.User = helper.String(v.(string))
		user = v.(string)
	}

	if v, ok := d.GetOk("virtual_host"); ok {
		request.VirtualHost = helper.String(v.(string))
		virtualHost = v.(string)
	}

	if v, ok := d.GetOk("config_regexp"); ok {
		request.ConfigRegexp = helper.String(v.(string))
	}

	if v, ok := d.GetOk("write_regexp"); ok {
		request.WriteRegexp = helper.String(v.(string))
	}

	if v, ok := d.GetOk("read_regexp"); ok {
		request.ReadRegexp = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQPermission(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqUserPermission failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{instanceId, user, virtualHost}, tccommon.FILED_SP))

	return resourceTencentCloudTdmqRabbitmqUserPermissionRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqUserPermissionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_user_permission.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	user := idSplit[1]
	virtualHost := idSplit[2]

	rabbitmqPermission, err := service.DescribeTdmqRabbitmqPermissionById(ctx, instanceId, user, virtualHost)
	if err != nil {
		return err
	}

	if rabbitmqPermission == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRabbitmqUserPermission` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rabbitmqPermission.InstanceId != nil {
		_ = d.Set("instance_id", rabbitmqPermission.InstanceId)
	}

	if rabbitmqPermission.User != nil {
		_ = d.Set("user", rabbitmqPermission.User)
	}

	if rabbitmqPermission.VirtualHost != nil {
		_ = d.Set("virtual_host", rabbitmqPermission.VirtualHost)
	}

	if rabbitmqPermission.ConfigRegexp != nil {
		_ = d.Set("config_regexp", rabbitmqPermission.ConfigRegexp)
	}

	if rabbitmqPermission.WriteRegexp != nil {
		_ = d.Set("write_regexp", rabbitmqPermission.WriteRegexp)
	}

	if rabbitmqPermission.ReadRegexp != nil {
		_ = d.Set("read_regexp", rabbitmqPermission.ReadRegexp)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqUserPermissionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_user_permission.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = tdmq.NewModifyRabbitMQPermissionRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	user := idSplit[1]
	virtualHost := idSplit[2]

	if d.HasChange("config_regexp") || d.HasChange("write_regexp") || d.HasChange("read_regexp") {
		request.InstanceId = &instanceId
		request.User = &user
		request.VirtualHost = &virtualHost

		if v, ok := d.GetOk("config_regexp"); ok {
			request.ConfigRegexp = helper.String(v.(string))
		}

		if v, ok := d.GetOk("write_regexp"); ok {
			request.WriteRegexp = helper.String(v.(string))
		}

		if v, ok := d.GetOk("read_regexp"); ok {
			request.ReadRegexp = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQPermission(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update tdmq rabbitmqUserPermission failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTdmqRabbitmqUserPermissionRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqUserPermissionDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_user_permission.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	user := idSplit[1]
	virtualHost := idSplit[2]

	if err := service.DeleteTdmqRabbitmqPermissionById(ctx, instanceId, user, virtualHost); err != nil {
		return err
	}

	return nil
}
