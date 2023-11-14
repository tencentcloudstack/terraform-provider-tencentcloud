/*
Provides a resource to create a monitor grafana_whitelist

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_whitelist" "grafana_whitelist" {
  instance_id = "grafana-abcdefgh"
  whitelist =
}
```

Import

monitor grafana_whitelist can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_whitelist.grafana_whitelist grafana_whitelist_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"log"
)

func resourceTencentCloudMonitorGrafanaWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaWhitelistCreate,
		Read:   resourceTencentCloudMonitorGrafanaWhitelistRead,
		Update: resourceTencentCloudMonitorGrafanaWhitelistUpdate,
		Delete: resourceTencentCloudMonitorGrafanaWhitelistDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance ID.",
			},

			"whitelist": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The addresses in the whitelist.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_whitelist.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorGrafanaWhitelistUpdate(d, meta)
}

func resourceTencentCloudMonitorGrafanaWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_whitelist.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	grafanaWhitelistId := d.Id()

	grafanaWhitelist, err := service.DescribeMonitorGrafanaWhitelistById(ctx, instanceId)
	if err != nil {
		return err
	}

	if grafanaWhitelist == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaWhitelist` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if grafanaWhitelist.InstanceId != nil {
		_ = d.Set("instance_id", grafanaWhitelist.InstanceId)
	}

	if grafanaWhitelist.Whitelist != nil {
		_ = d.Set("whitelist", grafanaWhitelist.Whitelist)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_whitelist.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateGrafanaWhiteListRequest()

	grafanaWhitelistId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "whitelist"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("whitelist") {
		if v, ok := d.GetOk("whitelist"); ok {
			whitelistSet := v.(*schema.Set).List()
			for i := range whitelistSet {
				whitelist := whitelistSet[i].(string)
				request.Whitelist = append(request.Whitelist, &whitelist)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateGrafanaWhiteList(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaWhitelist failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorGrafanaWhitelistRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_whitelist.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
