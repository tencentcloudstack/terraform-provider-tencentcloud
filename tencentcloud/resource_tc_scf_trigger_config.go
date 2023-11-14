/*
Provides a resource to create a scf trigger_config

Example Usage

```hcl
resource "tencentcloud_scf_trigger_config" "trigger_config" {
  enable = ""
  function_name = ""
  trigger_name = ""
  type = ""
  qualifier = ""
  namespace = ""
  trigger_desc = ""
}
```

Import

scf trigger_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_trigger_config.trigger_config trigger_config_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudScfTriggerConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfTriggerConfigCreate,
		Read:   resourceTencentCloudScfTriggerConfigRead,
		Delete: resourceTencentCloudScfTriggerConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"enable": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Initial status of the trigger. Values: `OPEN` (enabled); `CLOSE` disabled).",
			},

			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"trigger_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Trigger name.",
			},

			"type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Trigger Type.",
			},

			"qualifier": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function version. It defaults to `$LATEST`. It’s recommended to use `[$DEFAULT](https://intl.cloud.tencent.com/document/product/583/36149?from_cn_redirect=1#.E9.BB.98.E8.AE.A4.E5.88.AB.E5.90.8D)` for canary release.",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function namespace.",
			},

			"trigger_desc": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "To update a COS trigger, this field is required. It stores the data {event:cos:ObjectCreated:*} in the JSON format. The data content of this field is in the same format as that of SetTrigger. This field is optional if a scheduled trigger or CMQ trigger is to be deleted.",
			},
		},
	}
}

func resourceTencentCloudScfTriggerConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_trigger_config.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = scf.NewUpdateTriggerStatusRequest()
		response     = scf.NewUpdateTriggerStatusResponse()
		functionName string
		namespace    string
		triggerName  string
	)
	if v, ok := d.GetOk("enable"); ok {
		request.Enable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trigger_name"); ok {
		triggerName = v.(string)
		request.TriggerName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		request.Qualifier = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		namespace = v.(string)
		request.Namespace = helper.String(v.(string))
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
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate scf triggerConfig failed, reason:%+v", logId, err)
		return err
	}

	functionName = *response.Response.FunctionName
	d.SetId(strings.Join([]string{functionName, namespace, triggerName}, FILED_SP))

	return resourceTencentCloudScfTriggerConfigRead(d, meta)
}

func resourceTencentCloudScfTriggerConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_trigger_config.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudScfTriggerConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_trigger_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
