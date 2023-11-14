/*
Provides a resource to create a mariadb activate_hour_db_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_activate_hour_db_instance" "activate_hour_db_instance" {
  instance_ids =
}
```

Import

mariadb activate_hour_db_instance can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_activate_hour_db_instance.activate_hour_db_instance activate_hour_db_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"log"
)

func resourceTencentCloudMariadbActivateHourDbInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbActivateHourDbInstanceCreate,
		Read:   resourceTencentCloudMariadbActivateHourDbInstanceRead,
		Delete: resourceTencentCloudMariadbActivateHourDbInstanceDelete,
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
				Description: "Instance ID list.",
			},
		},
	}
}

func resourceTencentCloudMariadbActivateHourDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_activate_hour_db_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewActivateHourDBInstanceRequest()
		response   = mariadb.NewActivateHourDBInstanceResponse()
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

	instanceId = *response.Response.InstanceId
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
