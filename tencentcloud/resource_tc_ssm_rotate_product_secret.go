/*
Provides a resource to create a ssm rotate_product_secret

Example Usage

```hcl
resource "tencentcloud_ssm_rotate_product_secret" "example" {
  secret_name = "tf_example"
}
```
*/
package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssm/v20190923"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSsmRotateProductSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSsmRotateProductSecretCreate,
		Read:   resourceTencentCloudSsmRotateProductSecretRead,
		Delete: resourceTencentCloudSsmRotateProductSecretDelete,

		Schema: map[string]*schema.Schema{
			"secret_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Secret name.",
			},
		},
	}
}

func resourceTencentCloudSsmRotateProductSecretCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_rotate_product_secret.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId        = getLogId(contextNil)
		request      = ssm.NewRotateProductSecretRequest()
		response     = ssm.NewRotateProductSecretResponse()
		asyncRequest = ssm.NewDescribeAsyncRequestInfoRequest()
		secretName   string
		flowId       int64
	)

	if v, ok := d.GetOk("secret_name"); ok {
		request.SecretName = helper.String(v.(string))
		secretName = v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSsmClient().RotateProductSecret(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate ssm rotateProductSecret failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(secretName)

	// wait
	flowId = *response.Response.FlowID
	asyncRequest.FlowID = &flowId
	err = resource.Retry(writeRetryTimeout*3, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSsmClient().DescribeAsyncRequestInfo(asyncRequest)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if *result.Response.TaskStatus == 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("operate ssm describeAsyncRequestInfo status is %d", *result.Response.TaskStatus))
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate ssm describeAsyncRequestInfo failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSsmRotateProductSecretRead(d, meta)
}

func resourceTencentCloudSsmRotateProductSecretRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_rotate_product_secret.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSsmRotateProductSecretDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssm_rotate_product_secret.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
