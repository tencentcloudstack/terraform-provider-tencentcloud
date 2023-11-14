/*
Provides a resource to create a ssm describe_async_request_info

Example Usage

```hcl
resource "tencentcloud_ssm_describe_async_request_info" "describe_async_request_info" {
  flow_i_d = 1
}
```

Import

ssm describe_async_request_info can be imported using the id, e.g.

```
terraform import tencentcloud_ssm_describe_async_request_info.describe_async_request_info describe_async_request_info_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSsmDescribeAsyncRequestInfo() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmDescribeAsyncRequestInfoCreate,
		Read:   resourceTencentCloudSsmDescribeAsyncRequestInfoRead,
		Delete: resourceTencentCloudSsmDescribeAsyncRequestInfoDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_i_d": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Asynchronous task ID number.",
			},
		},
	}
}

func resourceTencentCloudSsmDescribeAsyncRequestInfoCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_describe_async_request_info.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = ssm.NewDescribeAsyncRequestInfoRequest()
		response = ssm.NewDescribeAsyncRequestInfoResponse()
		flowID   int
	)
	if v, _ := d.GetOk("flow_i_d"); v != nil {
		flowID = v.(int64)
		request.FlowID = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSsmClient().DescribeAsyncRequestInfo(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssm describeAsyncRequestInfo failed, reason:%+v", logId, err)
		return err
	}

	flowID = *response.Response.FlowID
	d.SetId(helper.Int64ToStr(flowID))

	return resourceTencentCloudSsmDescribeAsyncRequestInfoRead(d, meta)
}

func resourceTencentCloudSsmDescribeAsyncRequestInfoRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_describe_async_request_info.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSsmDescribeAsyncRequestInfoDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_describe_async_request_info.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
