/*
Provides a resource to create a mps parse_notify_operation

Example Usage

```hcl
resource "tencentcloud_mps_parse_notify_operation" "parse_notify_operation" {
  content = ""
}
```

Import

mps parse_notify_operation can be imported using the id, e.g.

```
terraform import tencentcloud_mps_parse_notify_operation.parse_notify_operation parse_notify_operation_id
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsParseNotifyOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsParseNotifyOperationCreate,
		Read:   resourceTencentCloudMpsParseNotifyOperationRead,
		Delete: resourceTencentCloudMpsParseNotifyOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"content": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Event notification obtained from CMQ.",
			},
		},
	}
}

func resourceTencentCloudMpsParseNotifyOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_parse_notify_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewParseNotificationRequest()
		response = mps.NewParseNotificationResponse()
		taskId   string
	)
	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ParseNotification(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps parseNotifyOperation failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.WorkflowTaskEvent.TaskId
	d.SetId(taskId)

	return resourceTencentCloudMpsParseNotifyOperationRead(d, meta)
}

func resourceTencentCloudMpsParseNotifyOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_parse_notify_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsParseNotifyOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_parse_notify_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
