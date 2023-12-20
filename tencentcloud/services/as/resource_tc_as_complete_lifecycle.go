package as

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsCompleteLifecycle() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsCompleteLifecycleCreate,
		Read:   resourceTencentCloudAsCompleteLifecycleRead,
		Delete: resourceTencentCloudAsCompleteLifecycleDelete,
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
				Description: "Result of the lifecycle action. Value range: `CONTINUE`, `ABANDON`.",
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

func resourceTencentCloudAsCompleteLifecycleCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_complete_lifecycle.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request         = as.NewCompleteLifecycleActionRequest()
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CompleteLifecycleAction(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate as completeLifecycle failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(lifecycleHookId)

	return resourceTencentCloudAsCompleteLifecycleRead(d, meta)
}

func resourceTencentCloudAsCompleteLifecycleRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_complete_lifecycle.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudAsCompleteLifecycleDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_complete_lifecycle.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
