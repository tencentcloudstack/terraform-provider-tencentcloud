package tat

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
)

func ResourceTencentCloudTatInvokerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTatInvokerConfigCreate,
		Read:   resourceTencentCloudTatInvokerConfigRead,
		Update: resourceTencentCloudTatInvokerConfigUpdate,
		Delete: resourceTencentCloudTatInvokerConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"invoker_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the invoker to be enabled.",
			},

			"invoker_status": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"on", "off"}),
				Description:  "Invoker on and off state, Values: `on`, `off`.",
			},
		},
	}
}

func resourceTencentCloudTatInvokerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invoker_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		invokerId string
	)
	if v, ok := d.GetOk("invoker_id"); ok {
		invokerId = v.(string)
	}

	d.SetId(invokerId)

	return resourceTencentCloudTatInvokerConfigUpdate(d, meta)
}

func resourceTencentCloudTatInvokerConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invoker_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	invokerId := d.Id()

	invokerConfig, err := service.DescribeTatInvokerConfigById(ctx, invokerId)
	if err != nil {
		return err
	}

	if invokerConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TatInvokerConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if invokerConfig.InvokerId != nil {
		_ = d.Set("invoker_id", invokerConfig.InvokerId)
	}

	if invokerConfig.Enable != nil && *invokerConfig.Enable {
		_ = d.Set("invoker_status", "on")
	} else {
		_ = d.Set("invoker_status", "off")
	}

	return nil
}

func resourceTencentCloudTatInvokerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invoker_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		disableInvokerRequest = tat.NewDisableInvokerRequest()
		enableInvokerRequest  = tat.NewEnableInvokerRequest()
		invokerId             = d.Id()
		err                   error
	)

	if v, ok := d.GetOk("invoker_status"); ok {
		status := v.(string)
		if status == "on" {
			enableInvokerRequest.InvokerId = &invokerId
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTatClient().EnableInvoker(enableInvokerRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableInvokerRequest.GetAction(), enableInvokerRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
		} else {
			disableInvokerRequest.InvokerId = &invokerId
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTatClient().DisableInvoker(disableInvokerRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableInvokerRequest.GetAction(), disableInvokerRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
		}
	}

	if err != nil {
		log.Printf("[CRITAL]%s update tat invokerConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTatInvokerConfigRead(d, meta)
}

func resourceTencentCloudTatInvokerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tat_invoker_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
