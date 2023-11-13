/*
Provides a resource to create a autoscaling complete_lifecycle

Example Usage

```hcl
resource "tencentcloud_autoscaling_complete_lifecycle" "complete_lifecycle" {
  lifecycle_hook_id = "ash-xxxxxxxx"
  lifecycle_action_result = "CONTINUE"
  instance_id = "ins-xxxxxxxx"
  lifecycle_action_token = &lt;nil&gt;
}
```

Import

autoscaling complete_lifecycle can be imported using the id, e.g.

```
terraform import tencentcloud_autoscaling_complete_lifecycle.complete_lifecycle complete_lifecycle_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	autoscaling "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudAutoscalingCompleteLifecycle() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAutoscalingCompleteLifecycleCreate,
		Read:   resourceTencentCloudAutoscalingCompleteLifecycleRead,
		Delete: resourceTencentCloudAutoscalingCompleteLifecycleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"lifecycle_hook_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Lifecycle hook ID.",
			},

			"lifecycle_action_result": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Result of the lifecycle action. Value range: CONTINUE, ABANDON.",
			},

			"instance_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID. Either InstanceId or LifecycleActionToken must be specified.",
			},

			"lifecycle_action_token": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Either InstanceId or LifecycleActionToken must be specified.",
			},
		},
	}
}

func resourceTencentCloudAutoscalingCompleteLifecycleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_complete_lifecycle.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = autoscaling.NewCompleteLifecycleActionRequest()
		response        = autoscaling.NewCompleteLifecycleActionResponse()
		lifecycleHookId string
	)
	if v, ok := d.GetOk("lifecycle_hook_id"); ok {
		lifecycleHookId = v.(string)
		request.LifecycleHookId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("lifecycle_action_result"); ok {
		request.LifecycleActionResult = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("lifecycle_action_token"); ok {
		request.LifecycleActionToken = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseAutoscalingClient().CompleteLifecycleAction(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate autoscaling completeLifecycle failed, reason:%+v", logId, err)
		return err
	}

	lifecycleHookId = *response.Response.LifecycleHookId
	d.SetId(lifecycleHookId)

	return resourceTencentCloudAutoscalingCompleteLifecycleRead(d, meta)
}

func resourceTencentCloudAutoscalingCompleteLifecycleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_complete_lifecycle.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAutoscalingCompleteLifecycleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_autoscaling_complete_lifecycle.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
