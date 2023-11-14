/*
Provides a resource to create a cdb restart_d_b_instances

Example Usage

```hcl
resource "tencentcloud_cdb_restart_d_b_instances" "restart_d_b_instances" {
  instance_ids =
}
```

Import

cdb restart_d_b_instances can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_restart_d_b_instances.restart_d_b_instances restart_d_b_instances_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"log"
	"time"
)

func resourceTencentCloudCdbRestartDBInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbRestartDBInstancesCreate,
		Read:   resourceTencentCloudCdbRestartDBInstancesRead,
		Delete: resourceTencentCloudCdbRestartDBInstancesDelete,
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
				Description: "An array of instance IDs in the format: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.",
			},
		},
	}
}

func resourceTencentCloudCdbRestartDBInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_restart_d_b_instances.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = cdb.NewRestartDBInstancesRequest()
		response = cdb.NewRestartDBInstancesResponse()
		idsHash  string
	)
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().RestartDBInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cdb restartDBInstances failed, reason:%+v", logId, err)
		return err
	}

	idsHash = *response.Response.IdsHash
	d.SetId(idsHash)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{""}, 1*readRetryTimeout, time.Second, service.CdbRestartDBInstancesStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbRestartDBInstancesRead(d, meta)
}

func resourceTencentCloudCdbRestartDBInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_restart_d_b_instances.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCdbRestartDBInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_restart_d_b_instances.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
