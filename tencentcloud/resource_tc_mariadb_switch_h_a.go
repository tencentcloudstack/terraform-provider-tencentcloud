/*
Provides a resource to create a mariadb switch_h_a

Example Usage

```hcl
resource "tencentcloud_mariadb_switch_h_a" "switch_h_a" {
  instance_id = ""
  zone = ""
}
```

Import

mariadb switch_h_a can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_switch_h_a.switch_h_a switch_h_a_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMariadbSwitchHA() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbSwitchHACreate,
		Read:   resourceTencentCloudMariadbSwitchHARead,
		Delete: resourceTencentCloudMariadbSwitchHADelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of tdsql-ow728lmc.",
			},

			"zone": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target AZ. The node with the lowest delay in the target AZ will be automatically promoted to primary node.",
			},
		},
	}
}

func resourceTencentCloudMariadbSwitchHACreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_switch_h_a.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewSwitchDBInstanceHARequest()
		response   = mariadb.NewSwitchDBInstanceHAResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok {
		request.Zone = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().SwitchDBInstanceHA(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb switchHA failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudMariadbSwitchHARead(d, meta)
}

func resourceTencentCloudMariadbSwitchHARead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_switch_h_a.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbSwitchHADelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_switch_h_a.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
