/*
Provides a resource to create a dcdb cancel_dcn_job_operation

Example Usage

```hcl
resource "tencentcloud_dcdb_cancel_dcn_job_operation" "cancel_dcn_job_operation" {
  instance_id = ""
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcdbCancelDcnJobOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbCancelDcnJobOperationCreate,
		Read:   resourceTencentCloudDcdbCancelDcnJobOperationRead,
		Delete: resourceTencentCloudDcdbCancelDcnJobOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},
		},
	}
}

func resourceTencentCloudDcdbCancelDcnJobOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_cancel_dcn_job_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewCancelDcnJobRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId =  helper.String(v.(string))
		instanceId=v.(string)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().CancelDcnJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb cancelDcnJobOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudDcdbCancelDcnJobOperationRead(d, meta)
}

func resourceTencentCloudDcdbCancelDcnJobOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_cancel_dcn_job_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbCancelDcnJobOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_cancel_dcn_job_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
