/*
Provides a resource to create a mps parse_live_stream_process_notify_operation

Example Usage

```hcl
resource "tencentcloud_mps_parse_live_stream_process_notify_operation" "operation" {
  content = "{\"EventType\":\"WorkflowTask\", xxx}"
}
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

func resourceTencentCloudMpsParseLiveStreamProcessNotifyOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsParseLiveStreamProcessNotifyOperationCreate,
		Read:   resourceTencentCloudMpsParseLiveStreamProcessNotifyOperationRead,
		Delete: resourceTencentCloudMpsParseLiveStreamProcessNotifyOperationDelete,
		Schema: map[string]*schema.Schema{
			"content": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Live stream event notification obtained from CMQ.",
			},
		},
	}
}

func resourceTencentCloudMpsParseLiveStreamProcessNotifyOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_parse_live_stream_process_notify_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewParseLiveStreamProcessNotificationRequest()
		response = mps.NewParseLiveStreamProcessNotificationResponse()
		taskId   string
	)
	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ParseLiveStreamProcessNotification(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps parseLiveStreamProcessNotifyOperation failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudMpsParseLiveStreamProcessNotifyOperationRead(d, meta)
}

func resourceTencentCloudMpsParseLiveStreamProcessNotifyOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_parse_live_stream_process_notify_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsParseLiveStreamProcessNotifyOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_parse_live_stream_process_notify_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
