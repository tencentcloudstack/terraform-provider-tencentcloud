/*
Provides a resource to create a live pullpush_task

Example Usage

```hcl
resource "tencentcloud_live_pullpush_task" "pullpush_task" {
  task_id = ""
  operator = ""
}
```

Import

live pullpush_task can be imported using the id, e.g.

```
terraform import tencentcloud_live_pullpush_task.pullpush_task pullpush_task_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLivePullpush_task() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLivePullpush_taskCreate,
		Read:   resourceTencentCloudLivePullpush_taskRead,
		Delete: resourceTencentCloudLivePullpush_taskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"task_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Task Id.",
			},

			"operator": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Task operator.",
			},
		},
	}
}

func resourceTencentCloudLivePullpush_taskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_pullpush_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = live.NewRestartLivePullStreamTaskRequest()
		response = live.NewRestartLivePullStreamTaskResponse()
		taskId   string
	)
	if v, ok := d.GetOk("task_id"); ok {
		taskId = v.(string)
		request.TaskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("operator"); ok {
		request.Operator = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().RestartLivePullStreamTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate live pullpush_task failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudLivePullpush_taskRead(d, meta)
}

func resourceTencentCloudLivePullpush_taskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_pullpush_task.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudLivePullpush_taskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_pullpush_task.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
