/*
Provides a resource to create a mps manage_task

Example Usage

```hcl
resource "tencentcloud_mps_manage_task" "manage_task" {
  operation_type = ""
  task_id = ""
}
```

Import

mps manage_task can be imported using the id, e.g.

```
terraform import tencentcloud_mps_manage_task.manage_task manage_task_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMpsManageTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsManageTaskCreate,
		Read:   resourceTencentCloudMpsManageTaskRead,
		Delete: resourceTencentCloudMpsManageTaskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"operation_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Operation type. Valid values:&amp;amp;lt;ul&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;Abort: task termination. Description:&amp;amp;lt;ul&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;If the [task type](https://intl.cloud.tencent.com/document/product/862/37614?from_cn_redirect=1#3.-.E8.BE.93.E5.87.BA.E5.8F.82.E6.95.B0) is live stream processing (`LiveStreamProcessTask`), tasks whose [task status](https://intl.cloud.tencent.com/document/product/862/37614?from_cn_redirect=1#3.-.E8.BE.93.E5.87.BA.E5.8F.82.E6.95.B0) is `WAITING` or `PROCESSING` can be terminated.&amp;amp;lt;/li&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;For other [task types](https://intl.cloud.tencent.com/document/product/862/37614?from_cn_redirect=1#3.-.E8.BE.93.E5.87.BA.E5.8F.82.E6.95.B0), only tasks whose [task status](https://intl.cloud.tencent.com/document/product/862/37614?from_cn_redirect=1#3.-.E8.BE.93.E5.87.BA.E5.8F.82.E6.95.B0) is `WAITING` can be terminated.&amp;amp;lt;/li&amp;amp;gt;&amp;amp;lt;/ul&amp;amp;gt;&amp;amp;lt;/li&amp;amp;gt;&amp;amp;lt;/ul&amp;amp;gt;.",
			},

			"task_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Video processing task ID.",
			},
		},
	}
}

func resourceTencentCloudMpsManageTaskCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_manage_task.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewManageTaskRequest()
		response = mps.NewManageTaskResponse()
		taskId   string
	)
	if v, ok := d.GetOk("operation_type"); ok {
		request.OperationType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_id"); ok {
		taskId = v.(string)
		request.TaskId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ManageTask(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps manageTask failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudMpsManageTaskRead(d, meta)
}

func resourceTencentCloudMpsManageTaskRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_manage_task.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsManageTaskDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_manage_task.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
