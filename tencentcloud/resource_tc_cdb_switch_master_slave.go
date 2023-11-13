/*
Provides a resource to create a cdb switch_master_slave

Example Usage

```hcl
resource "tencentcloud_cdb_switch_master_slave" "switch_master_slave" {
  instance_id = ""
  dst_slave = ""
  force_switch =
  wait_switch =
}
```

Import

cdb switch_master_slave can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_switch_master_slave.switch_master_slave switch_master_slave_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCdbSwitchMasterSlave() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbSwitchMasterSlaveCreate,
		Read:   resourceTencentCloudCdbSwitchMasterSlaveRead,
		Delete: resourceTencentCloudCdbSwitchMasterSlaveDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"dst_slave": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target instance. Possible values: `first` - first standby; `second` - second standby. The default value is `first`, and only multi-AZ instances support setting it to `second`.",
			},

			"force_switch": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to force switch. Default is False. Note that if you set the mandatory switch to True, there is a risk of data loss on the instance, so use it with caution.",
			},

			"wait_switch": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to switch within the time window. The default is False, i.e. do not switch within the time window. Note that if the ForceSwitch parameter is set to True, this parameter will not take effect.",
			},
		},
	}
}

func resourceTencentCloudCdbSwitchMasterSlaveCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_switch_master_slave.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewSwitchDBInstanceMasterSlaveRequest()
		response   = cdb.NewSwitchDBInstanceMasterSlaveResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dst_slave"); ok {
		request.DstSlave = helper.String(v.(string))
	}

	if v, _ := d.GetOk("force_switch"); v != nil {
		request.ForceSwitch = helper.Bool(v.(bool))
	}

	if v, _ := d.GetOk("wait_switch"); v != nil {
		request.WaitSwitch = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().SwitchDBInstanceMasterSlave(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cdb switchMasterSlave failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbSwitchMasterSlaveStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbSwitchMasterSlaveRead(d, meta)
}

func resourceTencentCloudCdbSwitchMasterSlaveRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_switch_master_slave.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCdbSwitchMasterSlaveDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_switch_master_slave.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
