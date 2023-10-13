/*
Provides a resource to create a monitor grafana_dns_config

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_dns_config" "grafana_dns_config" {
  instance_id  = "grafana-dp2hnnfa"
  name_servers = ["10.1.2.1", "10.1.2.2", "10.1.2.3"]
}
```

Import

monitor grafana_dns_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_dns_config.grafana_dns_config instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorGrafanaDnsConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaDnsConfigCreate,
		Read:   resourceTencentCloudMonitorGrafanaDnsConfigRead,
		Update: resourceTencentCloudMonitorGrafanaDnsConfigUpdate,
		Delete: resourceTencentCloudMonitorGrafanaDnsConfigDelete,
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

			"name_servers": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "DNS nameserver list.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaDnsConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_dns_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorGrafanaDnsConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorGrafanaDnsConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_dns_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()
	grafanaDnsConfig, err := service.DescribeMonitorGrafanaDnsConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if grafanaDnsConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaDnsConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if grafanaDnsConfig.NameServers != nil {
		_ = d.Set("name_servers", grafanaDnsConfig.NameServers)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaDnsConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_dns_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := monitor.NewUpdateDNSConfigRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("name_servers"); ok {
		nameServersSet := v.(*schema.Set).List()
		for i := range nameServersSet {
			nameServers := nameServersSet[i].(string)
			request.NameServers = append(request.NameServers, &nameServers)
		}
	}

	if len(request.NameServers) < 1 {
		request.NameServers = append(request.NameServers, helper.String(""))
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
		log.Printf("[CRITAL]%s update monitor grafanaDnsConfig failed, reason:%+v", logId, err)
		return err
	}

	time.Sleep(3 * time.Second)
	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(1*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeMonitorGrafanaInstance(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *instance.InstanceStatus == 2 {
			return nil
		}
		if *instance.InstanceStatus == 3 {
			return resource.NonRetryableError(fmt.Errorf("grafanaInstance status is %v, update dns config failed.", *instance.InstanceStatus))
		}
		return resource.RetryableError(fmt.Errorf("grafanaInstance status is %v, retry...", *instance.InstanceStatus))
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudMonitorGrafanaDnsConfigRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaDnsConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_dns_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
