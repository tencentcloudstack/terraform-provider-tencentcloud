/*
Provides a resource to create a tat invoker_config

Example Usage

```hcl
resource "tencentcloud_tat_invoker_config" "invoker_config" {
  invoker_id = "ivk-cas4upyf"
  invoker_status = "on"
}
```

Import

tat invoker_config can be imported using the id, e.g.

```
terraform import tencentcloud_tat_invoker_config.invoker_config invoker_config_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
)

func resourceTencentCloudTatInvokerConfig() *schema.Resource {
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
				ValidateFunc: validateAllowedStringValue([]string{"on", "off"}),
				Description:  "Invoker on and off state, Values: `on`, `off`.",
			},
		},
	}
}

func resourceTencentCloudTatInvokerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tat_invoker_config.create")()
	defer inconsistentCheck(d, meta)()

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
	defer logElapsed("resource.tencentcloud_tat_invoker_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TatService{client: meta.(*TencentCloudClient).apiV3Conn}

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
	defer logElapsed("resource.tencentcloud_tat_invoker_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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
			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseTatClient().EnableInvoker(enableInvokerRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableInvokerRequest.GetAction(), enableInvokerRequest.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
		} else {
			disableInvokerRequest.InvokerId = &invokerId
			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseTatClient().DisableInvoker(disableInvokerRequest)
				if e != nil {
					return retryError(e)
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
	defer logElapsed("resource.tencentcloud_tat_invoker_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
