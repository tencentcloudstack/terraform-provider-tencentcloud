/*
Provides a resource to create a monitor dns_config

Example Usage

```hcl
resource "tencentcloud_monitor_dns_config" "dns_config" {
  instance_id = "grafana-12345678"
  name_servers =
}
```

Import

monitor dns_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_dns_config.dns_config dns_config_id
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

func resourceTencentCloudMonitorDnsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorDnsConfigCreate,
		Read:   resourceTencentCloudMonitorDnsConfigRead,
		Update: resourceTencentCloudMonitorDnsConfigUpdate,
		Delete: resourceTencentCloudMonitorDnsConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance ID.",
			},

			"name_servers": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "DNS nameserver list.",
			},
		},
	}
}

func resourceTencentCloudMonitorDnsConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_dns_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorDnsConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorDnsConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_dns_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	dnsConfigId := d.Id()

	dnsConfig, err := service.DescribeMonitorDnsConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if dnsConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorDnsConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if dnsConfig.InstanceId != nil {
		_ = d.Set("instance_id", dnsConfig.InstanceId)
	}

	if dnsConfig.NameServers != nil {
		_ = d.Set("name_servers", dnsConfig.NameServers)
	}

	return nil
}

func resourceTencentCloudMonitorDnsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_dns_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateDNSConfigRequest()

	dnsConfigId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "name_servers"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name_servers") {
		if v, ok := d.GetOk("name_servers"); ok {
			nameServersSet := v.(*schema.Set).List()
			for i := range nameServersSet {
				nameServers := nameServersSet[i].(string)
				request.NameServers = append(request.NameServers, &nameServers)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateDNSConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor dnsConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorDnsConfigRead(d, meta)
}

func resourceTencentCloudMonitorDnsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_dns_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
