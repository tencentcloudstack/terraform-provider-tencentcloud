/*
Provides a resource to create a mps execute_function

Example Usage

```hcl
resource "tencentcloud_mps_execute_function" "execute_function" {
  function_name = ""
  function_arg = ""
}
```

Import

mps execute_function can be imported using the id, e.g.

```
terraform import tencentcloud_mps_execute_function.execute_function execute_function_id
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

func resourceTencentCloudMpsExecuteFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsExecuteFunctionCreate,
		Read:   resourceTencentCloudMpsExecuteFunctionRead,
		Delete: resourceTencentCloudMpsExecuteFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Name of called backend API.",
			},

			"function_arg": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "API parameter. Parameter format will depend on the actual function definition.",
			},
		},
	}
}

func resourceTencentCloudMpsExecuteFunctionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_execute_function.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = mps.NewExecuteFunctionRequest()
		response     = mps.NewExecuteFunctionResponse()
		functionName string
	)
	if v, ok := d.GetOk("function_name"); ok {
		functionName = v.(string)
		request.FunctionName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_arg"); ok {
		request.FunctionArg = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ExecuteFunction(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mps executeFunction failed, reason:%+v", logId, err)
		return err
	}

	functionName = *response.Response.FunctionName
	d.SetId(functionName)

	return resourceTencentCloudMpsExecuteFunctionRead(d, meta)
}

func resourceTencentCloudMpsExecuteFunctionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_execute_function.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsExecuteFunctionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_execute_function.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
