/*
Provides a resource to create a monitor tmp_grafana_config

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_grafana_config" "tmp_grafana_config" {
    config      = jsonencode(
        {
            server = {
                http_port           = 8080
                root_url            = "https://cloud-grafana.woa.com/grafana-ffrdnrfa/"
                serve_from_sub_path = true
            }
        }
    )
    instance_id = "grafana-29phe08q"
}
```

Import

monitor tmp_grafana_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_grafana_config.tmp_grafana_config tmp_grafana_config_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpGrafanaConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpGrafanaConfigCreate,
		Read:   resourceTencentCloudMonitorTmpGrafanaConfigRead,
		Update: resourceTencentCloudMonitorTmpGrafanaConfigUpdate,
		Delete: resourceTencentCloudMonitorTmpGrafanaConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"config": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "JSON encoded string.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpGrafanaConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_grafana_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorTmpGrafanaConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorTmpGrafanaConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_grafana_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	tmpGrafanaConfig, err := service.DescribeMonitorTmpGrafanaConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if tmpGrafanaConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpGrafanaConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if tmpGrafanaConfig.Config != nil {
		_ = d.Set("config", tmpGrafanaConfig.Config)
	}

	return nil
}

func resourceTencentCloudMonitorTmpGrafanaConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_grafana_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateGrafanaConfigRequest()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("config"); ok {
		request.Config = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateGrafanaConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor tmpGrafanaConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorTmpGrafanaConfigRead(d, meta)
}

func resourceTencentCloudMonitorTmpGrafanaConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_grafana_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
