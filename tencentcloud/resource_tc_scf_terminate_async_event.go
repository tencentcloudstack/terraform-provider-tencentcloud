/*
Provides a resource to create a scf terminate_async_event

Example Usage

```hcl
resource "tencentcloud_scf_terminate_async_event" "terminate_async_event" {
  function_name = "test"
  invoke_request_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  namespace = "testNamespace"
  grace_shutdown = true
}
```

Import

scf terminate_async_event can be imported using the id, e.g.

```
terraform import tencentcloud_scf_terminate_async_event.terminate_async_event terminate_async_event_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudScfTerminateAsyncEvent() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudScfTerminateAsyncEventCreate,
		Read:   resourceTencentCloudScfTerminateAsyncEventRead,
		Delete: resourceTencentCloudScfTerminateAsyncEventDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"invoke_request_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Terminated invocation request ID.",
			},

			"namespace": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Namespace.",
			},

			"grace_shutdown": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable grace shutdown. If it’s true, a SIGTERM signal is sent to the specified request. See [Sending termination signal](https://www.tencentcloud.com/document/product/583/63969?from_cn_redirect=1#.E5.8F.91.E9.80.81.E7.BB.88.E6.AD.A2.E4.BF.A1.E5.8F.B7]. It’s set to false by default.",
			},
		},
	}
}

func resourceTencentCloudScfTerminateAsyncEventCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_terminate_async_event.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = scf.NewTerminateAsyncEventRequest()
		response        = scf.NewTerminateAsyncEventResponse()
		invokeRequestId string
	)
	if v, ok := d.GetOk("function_name"); ok {
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("invoke_request_id"); ok {
		invokeRequestId = v.(string)
		request.InvokeRequestId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
	}

	if v, _ := d.GetOk("grace_shutdown"); v != nil {
		request.GraceShutdown = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseScfClient().TerminateAsyncEvent(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate scf terminateAsyncEvent failed, reason:%+v", logId, err)
		return err
	}

	invokeRequestId = *response.Response.InvokeRequestId
	d.SetId(invokeRequestId)

	return resourceTencentCloudScfTerminateAsyncEventRead(d, meta)
}

func resourceTencentCloudScfTerminateAsyncEventRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_terminate_async_event.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudScfTerminateAsyncEventDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_scf_terminate_async_event.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
