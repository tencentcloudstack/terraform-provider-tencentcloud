/*
Provides a resource to create a monitor grafana_whitelist_config

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_whitelist_config" "grafana_whitelist_config" {
  instance_id = "grafana-dp2hnnfa"
  whitelist   = ["10.1.1.1", "10.1.1.2", "10.1.1.3"]
}
```

Import

monitor grafana_whitelist_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_whitelist_config.grafana_whitelist_config instance_id
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorGrafanaWhitelistConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaWhitelistConfigCreate,
		Read:   resourceTencentCloudMonitorGrafanaWhitelistConfigRead,
		Update: resourceTencentCloudMonitorGrafanaWhitelistConfigUpdate,
		Delete: resourceTencentCloudMonitorGrafanaWhitelistConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance ID.",
			},

			"whitelist": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The addresses in the whitelist.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaWhitelistConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_whitelist_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorGrafanaWhitelistConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorGrafanaWhitelistConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_whitelist_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	grafanaWhitelistConfig, err := service.DescribeMonitorGrafanaWhitelistConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if grafanaWhitelistConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaWhitelistConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if grafanaWhitelistConfig.WhiteList != nil {
		whiteList := grafanaWhitelistConfig.WhiteList
		if len(whiteList) == 1 && *whiteList[0] == "" {
			return nil
		}
		if len(whiteList) == 1 && strings.Contains(*whiteList[0], "\n") {
			_ = d.Set("whitelist", strings.Split(*whiteList[0], "\n"))
		}
		if len(whiteList) > 1 {
			_ = d.Set("whitelist", whiteList)
		}
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaWhitelistConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_whitelist_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateGrafanaWhiteListRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("whitelist"); ok {
		whitelistSet := v.(*schema.Set).List()
		for i := range whitelistSet {
			whitelist := whitelistSet[i].(string)
			request.Whitelist = append(request.Whitelist, &whitelist)
		}
	}

	if len(request.Whitelist) < 1 {
		request.Whitelist = append(request.Whitelist, helper.String(""))
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
		log.Printf("[CRITAL]%s update monitor grafanaWhitelistConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorGrafanaWhitelistConfigRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaWhitelistConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_whitelist_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
