/*
Provides a resource to create a mariadb isolate_instance

Example Usage

```hcl
resource "tencentcloud_mariadb_isolate_instance" "isolate_instance" {
  instance_ids =
}
```

Import

mariadb isolate_instance can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_isolate_instance.isolate_instance isolate_instance_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"log"
)

func resourceTencentCloudMariadbIsolateInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbIsolateInstanceCreate,
		Read:   resourceTencentCloudMariadbIsolateInstanceRead,
		Delete: resourceTencentCloudMariadbIsolateInstanceDelete,
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
				Description: "Instance ID in the format of `tdsql-dasjkhd`. It is the same as the instance ID displayed on the TencentDB console. You can use the `DescribeDBInstances` API to query the ID, which is the value of the `InstanceId` output parameter.",
			},
		},
	}
}

func resourceTencentCloudMariadbIsolateInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_isolate_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewIsolateDBInstanceRequest()
		response   = mariadb.NewIsolateDBInstanceResponse()
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().IsolateDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb isolateInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudMariadbIsolateInstanceRead(d, meta)
}

func resourceTencentCloudMariadbIsolateInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_isolate_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbIsolateInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_isolate_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
