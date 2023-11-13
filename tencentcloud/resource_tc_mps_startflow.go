/*
Provides a resource to create a mps startflow

Example Usage

```hcl
resource "tencentcloud_mps_startflow" "startflow" {
  flow_id = "your-flow-id"
}
```

Import

mps startflow can be imported using the id, e.g.

```
terraform import tencentcloud_mps_startflow.startflow startflow_id
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

func resourceTencentCloudMpsStartflow() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsStartflowCreate,
		Read:   resourceTencentCloudMpsStartflowRead,
		Delete: resourceTencentCloudMpsStartflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Flow Id.",
			},
		},
	}
}

func resourceTencentCloudMpsStartflowCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_startflow.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewStartStreamLinkFlowRequest()
		response = mps.NewStartStreamLinkFlowResponse()
		flowId   string
	)
	if v, ok := d.GetOk("flow_id"); ok {
		flowId = v.(string)
		request.FlowId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().StartStreamLinkFlow(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps startflow failed, reason:%+v", logId, err)
		return err
	}

	flowId = *response.Response.FlowId
	d.SetId(flowId)

	return resourceTencentCloudMpsStartflowRead(d, meta)
}

func resourceTencentCloudMpsStartflowRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_startflow.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsStartflowDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_startflow.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
