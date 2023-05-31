/*
Provides a resource to create a scf trigger_config

Example Usage

```hcl
resource "tencentcloud_scf_trigger_config" "trigger_config" {
  enable        = "OPEN"
  function_name = "keep-1676351130"
  trigger_name  = "SCF-timer-1685540160"
  type          = "timer"
  qualifier     = "$DEFAULT"
  namespace     = "default"
}
```

Import

scf trigger_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_trigger_config.trigger_config functionName#namespace#triggerName
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudScfTriggerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfTriggerConfigCreate,
		Read:   resourceTencentCloudScfTriggerConfigRead,
		Update: resourceTencentCloudScfTriggerConfigUpdate,
		Delete: resourceTencentCloudScfTriggerConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"trigger_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Trigger name.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Trigger Type.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Initial status of the trigger. Values: `OPEN` (enabled); `CLOSE` disabled).",
			},

			"qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function version. It defaults to `$LATEST`. It's recommended to use `[$DEFAULT](https://intl.cloud.tencent.com/document/product/583/36149?from_cn_redirect=1#.E9.BB.98.E8.AE.A4.E5.88.AB.E5.90.8D)` for canary release.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "default",
				Description: "Function namespace.",
			},

			"trigger_desc": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "To update a COS trigger, this field is required. It stores the data {event:cos:ObjectCreated:*} in the JSON format. The data content of this field is in the same format as that of SetTrigger. This field is optional if a scheduled trigger or CMQ trigger is to be deleted.",
			},
		},
	}
}

func resourceTencentCloudScfTriggerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_trigger_config.create")()
	defer inconsistentCheck(d, meta)()

	var (
		functionName string
		triggerName  string
		namespace    string
	)
	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
	}

	if v, ok := d.GetOk("trigger_name"); ok {
		triggerName = v.(string)
	}

	d.SetId(functionName + FILED_SP + namespace + FILED_SP + triggerName)

	return resourceTencentCloudScfTriggerConfigUpdate(d, meta)
}

func resourceTencentCloudScfTriggerConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_trigger_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	namespace := idSplit[1]
	triggerName := idSplit[2]

	triggerConfig, err := service.DescribeScfTriggerConfigById(ctx, functionName, namespace, triggerName)
	if err != nil {
		return err
	}

	if triggerConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `ScfTriggerConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if triggerConfig.Enable != nil {
		if *triggerConfig.Enable == 1 {
			_ = d.Set("enable", "OPEN")
		} else {
			_ = d.Set("enable", "CLOSE")
		}
	}

	_ = d.Set("function_name", functionName)

	_ = d.Set("namespace", namespace)

	if triggerConfig.TriggerName != nil {
		_ = d.Set("trigger_name", triggerConfig.TriggerName)
	}

	if triggerConfig.Type != nil {
		_ = d.Set("type", triggerConfig.Type)
	}

	if triggerConfig.Qualifier != nil {
		_ = d.Set("qualifier", triggerConfig.Qualifier)
	}

	if triggerConfig.TriggerDesc != nil {
		_ = d.Set("trigger_desc", triggerConfig.TriggerDesc)
	}

	return nil
}

func resourceTencentCloudScfTriggerConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_trigger_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := scf.NewUpdateTriggerStatusRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	functionName := idSplit[0]
	namespace := idSplit[1]
	triggerName := idSplit[2]

	request.FunctionName = &functionName
	request.Namespace = &namespace
	request.TriggerName = &triggerName

	if v, ok := d.GetOk("enable"); ok {
		request.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		request.Qualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trigger_desc"); ok {
		request.TriggerDesc = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().UpdateTriggerStatus(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update scf triggerConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudScfTriggerConfigRead(d, meta)
}

func resourceTencentCloudScfTriggerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_trigger_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
