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

func ResourceTencentCloudTdmqRabbitmqVirtualHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTdmqRabbitmqVirtualHostCreate,
		Read:   resourceTencentCloudTdmqRabbitmqVirtualHostRead,
		Update: resourceTencentCloudTdmqRabbitmqVirtualHostUpdate,
		Delete: resourceTencentCloudTdmqRabbitmqVirtualHostDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster instance ID.",
			},
			"virtual_host": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "vhost name.",
			},
			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "describe.",
			},
			"trace_flag": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Message track switch, true is on, false is off, default is off.",
			},
		},
	}
}

func resourceTencentCloudTdmqRabbitmqVirtualHostCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_virtual_host.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		request     = tdmq.NewCreateRabbitMQVirtualHostRequest()
		response    = tdmq.NewCreateRabbitMQVirtualHostResponse()
		instanceId  string
		virtualHost string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("virtual_host"); ok {
		request.VirtualHost = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("trace_flag"); ok {
		request.TraceFlag = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().CreateRabbitMQVirtualHost(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmq rabbitmqVirtualHost failed, reason:%+v", logId, err)
		return err
	}

	virtualHost = *response.Response.VirtualHost
	d.SetId(strings.Join([]string{instanceId, virtualHost}, tccommon.FILED_SP))

	return resourceTencentCloudTdmqRabbitmqVirtualHostRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqVirtualHostRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_virtual_host.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}
	instanceId := idSplit[0]
	virtualHost := idSplit[1]

	rabbitmqVirtualHost, err := service.DescribeTdmqRabbitmqVirtualHostById(ctx, instanceId, virtualHost)
	if err != nil {
		return err
	}

	if rabbitmqVirtualHost == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TdmqRabbitmqVirtualHost` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if rabbitmqVirtualHost.InstanceId != nil {
		_ = d.Set("instance_id", rabbitmqVirtualHost.InstanceId)
	}

	if rabbitmqVirtualHost.VirtualHost != nil {
		_ = d.Set("virtual_host", rabbitmqVirtualHost.VirtualHost)
	}

	if rabbitmqVirtualHost.Description != nil {
		_ = d.Set("description", rabbitmqVirtualHost.Description)
	}

	return nil
}

func resourceTencentCloudTdmqRabbitmqVirtualHostUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_virtual_host.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = tdmq.NewModifyRabbitMQVirtualHostRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	virtualHost := idSplit[1]

	immutableArgs := []string{"instance_id", "virtual_host"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("description") || d.HasChange("trace_flag") {
		request.InstanceId = &instanceId
		request.VirtualHost = &virtualHost

		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("trace_flag"); ok {
			request.TraceFlag = helper.Bool(v.(bool))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTdmqClient().ModifyRabbitMQVirtualHost(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update tdmq rabbitmqVirtualHost failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTdmqRabbitmqVirtualHostRead(d, meta)
}

func resourceTencentCloudTdmqRabbitmqVirtualHostDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tdmq_rabbitmq_virtual_host.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", idSplit)
	}

	instanceId := idSplit[0]
	virtualHost := idSplit[1]

	if err := service.DeleteTdmqRabbitmqVirtualHostById(ctx, instanceId, virtualHost); err != nil {
		return err
	}

	return nil
}
