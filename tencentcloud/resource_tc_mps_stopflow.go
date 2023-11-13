/*
Provides a resource to create a mps stopflow

Example Usage

```hcl
resource "tencentcloud_mps_stopflow" "stopflow" {
  flow_id = "your flow id"
}
```

Import

mps stopflow can be imported using the id, e.g.

```
terraform import tencentcloud_mps_stopflow.stopflow stopflow_id
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

func resourceTencentCloudMpsStopflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsStopflowCreate,
		Read:   resourceTencentCloudMpsStopflowRead,
		Delete: resourceTencentCloudMpsStopflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Flow Id。.",
			},
		},
	}
}

func resourceTencentCloudMpsStopflowCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_stopflow.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewStopStreamLinkFlowRequest()
		response = mps.NewStopStreamLinkFlowResponse()
		flowId   string
	)
	if v, ok := d.GetOk("flow_id"); ok {
		flowId = v.(string)
		request.FlowId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().StopStreamLinkFlow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps stopflow failed, reason:%+v", logId, err)
		return err
	}

	flowId = *response.Response.FlowId
	d.SetId(flowId)

	return resourceTencentCloudMpsStopflowRead(d, meta)
}

func resourceTencentCloudMpsStopflowRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_stopflow.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsStopflowDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_stopflow.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
