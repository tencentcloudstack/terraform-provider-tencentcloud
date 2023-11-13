/*
Provides a resource to create a dcdb switch_d_b_instance_h_a

Example Usage

```hcl
resource "tencentcloud_dcdb_switch_d_b_instance_h_a" "switch_d_b_instance_h_a" {
  instance_id = ""
  zone = ""
}
```

Import

dcdb switch_d_b_instance_h_a can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_switch_d_b_instance_h_a.switch_d_b_instance_h_a switch_d_b_instance_h_a_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDcdbSwitchDBInstanceHA() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbSwitchDBInstanceHACreate,
		Read:   resourceTencentCloudDcdbSwitchDBInstanceHARead,
		Delete: resourceTencentCloudDcdbSwitchDBInstanceHADelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of tdsqlshard-ow728lmc.",
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

func resourceTencentCloudDcdbSwitchDBInstanceHACreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_switch_d_b_instance_h_a.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewSwitchDBInstanceHARequest()
		response   = dcdb.NewSwitchDBInstanceHAResponse()
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().SwitchDBInstanceHA(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb switchDBInstanceHA failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudDcdbSwitchDBInstanceHARead(d, meta)
}

func resourceTencentCloudDcdbSwitchDBInstanceHARead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_switch_d_b_instance_h_a.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbSwitchDBInstanceHADelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_switch_d_b_instance_h_a.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
