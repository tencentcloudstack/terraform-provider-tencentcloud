/*
Provides a resource to create a monitor grafana_version_upgrade

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_version_upgrade" "grafana_version_upgrade" {
  instance_id = ""
  alias = ""
}
```

Import

monitor grafana_version_upgrade can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_version_upgrade.grafana_version_upgrade grafana_version_upgrade_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMonitorGrafanaVersionUpgrade() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaVersionUpgradeCreate,
		Read:   resourceTencentCloudMonitorGrafanaVersionUpgradeRead,
		Update: resourceTencentCloudMonitorGrafanaVersionUpgradeUpdate,
		Delete: resourceTencentCloudMonitorGrafanaVersionUpgradeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance ID.",
			},

			"alias": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Version alias.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaVersionUpgradeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_version_upgrade.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorGrafanaVersionUpgradeUpdate(d, meta)
}

func resourceTencentCloudMonitorGrafanaVersionUpgradeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_version_upgrade.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	grafanaVersionUpgradeId := d.Id()

	grafanaVersionUpgrade, err := service.DescribeMonitorGrafanaVersionUpgradeById(ctx, instanceId)
	if err != nil {
		return err
	}

	if grafanaVersionUpgrade == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaVersionUpgrade` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if grafanaVersionUpgrade.InstanceId != nil {
		_ = d.Set("instance_id", grafanaVersionUpgrade.InstanceId)
	}

	if grafanaVersionUpgrade.Alias != nil {
		_ = d.Set("alias", grafanaVersionUpgrade.Alias)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaVersionUpgradeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_version_upgrade.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpgradeGrafanaInstanceRequest()

	grafanaVersionUpgradeId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "alias"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("alias") {
		if v, ok := d.GetOk("alias"); ok {
			request.Alias = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpgradeGrafanaInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaVersionUpgrade failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorGrafanaVersionUpgradeRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaVersionUpgradeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_version_upgrade.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
