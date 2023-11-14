/*
Provides a resource to create a dcdb activate_hour_instance

Example Usage

```hcl
resource "tencentcloud_dcdb_activate_hour_instance" "activate_hour_instance" {
  instance_ids =
}
```

Import

dcdb activate_hour_instance can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_activate_hour_instance.activate_hour_instance activate_hour_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"log"
)

func resourceTencentCloudDcdbActivateHourInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbActivateHourInstanceCreate,
		Read:   resourceTencentCloudDcdbActivateHourInstanceRead,
		Delete: resourceTencentCloudDcdbActivateHourInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of instance IDs in the format of dcdbt-ow728lmc, which can be obtained through the `DescribeDCDBInstances` API.",
			},
		},
	}
}

func resourceTencentCloudDcdbActivateHourInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_activate_hour_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewActiveHourDCDBInstanceRequest()
		response   = dcdb.NewActiveHourDCDBInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ActiveHourDCDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb activateHourInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudDcdbActivateHourInstanceRead(d, meta)
}

func resourceTencentCloudDcdbActivateHourInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_activate_hour_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbActivateHourInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_activate_hour_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
