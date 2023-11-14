/*
Provides a resource to create a mps parse_live_stream_process_notification

Example Usage

```hcl
resource "tencentcloud_mps_parse_live_stream_process_notification" "parse_live_stream_process_notification" {
  content = &lt;nil&gt;
}
```

Import

mps parse_live_stream_process_notification can be imported using the id, e.g.

```
terraform import tencentcloud_mps_parse_live_stream_process_notification.parse_live_stream_process_notification parse_live_stream_process_notification_id
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

func resourceTencentCloudMpsParseLiveStreamProcessNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsParseLiveStreamProcessNotificationCreate,
		Read:   resourceTencentCloudMpsParseLiveStreamProcessNotificationRead,
		Delete: resourceTencentCloudMpsParseLiveStreamProcessNotificationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"content": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The live stream event notification content obtained from CMQ.",
			},
		},
	}
}

func resourceTencentCloudMpsParseLiveStreamProcessNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_parse_live_stream_process_notification.create")()
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
		log.Printf("[CRITAL]%s operate mps parseLiveStreamProcessNotification failed, reason:%+v", logId, err)
		return err
	}

	taskId = *response.Response.TaskId
	d.SetId(taskId)

	return resourceTencentCloudMpsParseLiveStreamProcessNotificationRead(d, meta)
}

func resourceTencentCloudMpsParseLiveStreamProcessNotificationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_parse_live_stream_process_notification.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsParseLiveStreamProcessNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_parse_live_stream_process_notification.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
