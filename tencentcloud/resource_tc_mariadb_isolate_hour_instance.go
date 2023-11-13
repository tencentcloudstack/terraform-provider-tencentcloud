/*
Provides a resource to create a mariadb isolate_hour_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_isolate_hour_instance" "isolate_hour_instance" {
  instance_ids =
}
```

Import

mariadb isolate_hour_instance can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_isolate_hour_instance.isolate_hour_instance isolate_hour_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"log"
)

func resourceTencentCloudMariadbIsolateHourInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbIsolateHourInstanceCreate,
		Read:   resourceTencentCloudMariadbIsolateHourInstanceRead,
		Delete: resourceTencentCloudMariadbIsolateHourInstanceDelete,
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

func resourceTencentCloudMariadbIsolateHourInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_isolate_hour_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewIsolateHourDBInstanceRequest()
		response   = mariadb.NewIsolateHourDBInstanceResponse()
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().IsolateHourDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb isolateHourInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudMariadbIsolateHourInstanceRead(d, meta)
}

func resourceTencentCloudMariadbIsolateHourInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_isolate_hour_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbIsolateHourInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_isolate_hour_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
