/*
Provides a resource to create a cdb switch_for_upgrade

Example Usage

```hcl
resource "tencentcloud_cdb_switch_for_upgrade" "switch_for_upgrade" {
  instance_id = ""
}
```

Import

cdb switch_for_upgrade can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_switch_for_upgrade.switch_for_upgrade switch_for_upgrade_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCdbSwitchForUpgrade() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbSwitchForUpgradeCreate,
		Read:   resourceTencentCloudCdbSwitchForUpgradeRead,
		Delete: resourceTencentCloudCdbSwitchForUpgradeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed on the TencentDB Console page.",
			},
		},
	}
}

func resourceTencentCloudCdbSwitchForUpgradeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_switch_for_upgrade.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewSwitchForUpgradeRequest()
		response   = cdb.NewSwitchForUpgradeResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().SwitchForUpgrade(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cdb switchForUpgrade failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{0}, 1*readRetryTimeout, time.Second, service.CdbSwitchForUpgradeStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbSwitchForUpgradeRead(d, meta)
}

func resourceTencentCloudCdbSwitchForUpgradeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_switch_for_upgrade.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	switchForUpgradeId := d.Id()

	switchForUpgrade, err := service.DescribeCdbSwitchForUpgradeById(ctx, instanceId)
	if err != nil {
		return err
	}

	if switchForUpgrade == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CdbSwitchForUpgrade` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if switchForUpgrade.InstanceId != nil {
		_ = d.Set("instance_id", switchForUpgrade.InstanceId)
	}

	return nil
}

func resourceTencentCloudCdbSwitchForUpgradeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_switch_for_upgrade.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	switchForUpgradeId := d.Id()

	if err := service.DeleteCdbSwitchForUpgradeById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
