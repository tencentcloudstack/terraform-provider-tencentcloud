/*
Provides a resource to create a dcdb activate_hour_instance_operation

Example Usage

```hcl
resource "tencentcloud_dcdb_activate_hour_instance_operation" "activate_hour_instance_operation" {
  instance_id = local.dcdb_id
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

func resourceTencentCloudDcdbActivateHourInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbActivateHourInstanceOperationCreate,
		Read:   resourceTencentCloudDcdbActivateHourInstanceOperationRead,
		Delete: resourceTencentCloudDcdbActivateHourInstanceOperationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "instance ID in the format of dcdbt-ow728lmc, which can be obtained through the `DescribeDCDBInstances` API.",
			},
		},
	}
}

func resourceTencentCloudDcdbActivateHourInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_activate_hour_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewActiveHourDCDBInstanceRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceIds = []*string{helper.String(instanceId)}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ActiveHourDCDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb activateHourInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudDcdbActivateHourInstanceOperationRead(d, meta)
}

func resourceTencentCloudDcdbActivateHourInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_activate_hour_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbActivateHourInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_activate_hour_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
