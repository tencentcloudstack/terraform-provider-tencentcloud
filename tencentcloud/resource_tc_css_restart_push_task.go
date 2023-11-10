/*
Provides a resource to create a css restart_push_task

Example Usage

```hcl
resource "tencentcloud_css_restart_push_task" "restart_push_task" {
  task_id = ""
  operator = ""
}
```
*/
package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssRestartPushTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssRestartPushTaskCreate,
		Read:   resourceTencentCloudCssRestartPushTaskRead,
		Delete: resourceTencentCloudCssRestartPushTaskDelete,
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

func resourceTencentCloudCssRestartPushTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_restart_push_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = css.NewRestartLivePullStreamTaskRequest()
		taskId  string
	)
	if v, ok := d.GetOk("task_id"); ok {
		taskId = v.(string)
		request.TaskId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("operator"); ok {
		request.Operator = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().RestartLivePullStreamTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate css restartPushTask failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(taskId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"active"}, 3*readRetryTimeout, time.Second, service.CssRestartPushTaskStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCssRestartPushTaskRead(d, meta)
}

func resourceTencentCloudCssRestartPushTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_restart_push_task.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCssRestartPushTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_restart_push_task.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
