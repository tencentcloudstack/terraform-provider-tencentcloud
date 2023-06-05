/*
Provides a resource to create a mariadb activate_hour_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_activate_hour_db_instance" "activate_hour_db_instance" {
  instance_id = "tdsql-9vqvls95"
}
```
*/
package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
)

func resourceTencentCloudMariadbActivateHourDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbActivateHourDbInstanceCreate,
		Read:   resourceTencentCloudMariadbActivateHourDbInstanceRead,
		Delete: resourceTencentCloudMariadbActivateHourDbInstanceDelete,

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

func resourceTencentCloudMariadbActivateHourDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_activate_hour_db_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		request    = mariadb.NewActivateHourDBInstanceRequest()
		response   = mariadb.NewActivateHourDBInstanceResponse()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceIds = common.StringPtrs([]string{v.(string)})
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ActivateHourDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb activateHourDbInstance failed, reason:%+v", logId, err)
		return err
	}

	if response == nil {
		return fmt.Errorf("operate mariadb activateHourDbInstance not found")
	}

	instanceId = *response.Response.SuccessInstanceIds[0]

	d.SetId(instanceId)

	return resourceTencentCloudMariadbActivateHourDbInstanceRead(d, meta)
}

func resourceTencentCloudMariadbActivateHourDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_activate_hour_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbActivateHourDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_activate_hour_db_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
