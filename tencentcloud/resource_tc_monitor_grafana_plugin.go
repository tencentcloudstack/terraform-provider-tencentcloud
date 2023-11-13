/*
Provides a resource to create a monitor grafana_plugin

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_plugin" "grafana_plugin" {
  instance_id = &lt;nil&gt;
  plugin_id = &lt;nil&gt;
  version = &lt;nil&gt;
}
```

Import

monitor grafana_plugin can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_plugin.grafana_plugin grafana_plugin_id
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

func resourceTencentCloudMonitorGrafanaPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaPluginCreate,
		Read:   resourceTencentCloudMonitorGrafanaPluginRead,
		Update: resourceTencentCloudMonitorGrafanaPluginUpdate,
		Delete: resourceTencentCloudMonitorGrafanaPluginDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance id.",
			},

			"plugin_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Plugin id.",
			},

			"version": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Plugin version.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaPluginCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_plugin.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewInstallPluginsRequest()
		response = monitor.NewInstallPluginsResponse()
		pluginId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("plugin_id"); ok {
		pluginId = v.(string)
		request.PluginId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("version"); ok {
		request.Version = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().InstallPlugins(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create monitor grafanaPlugin failed, reason:%+v", logId, err)
		return err
	}

	pluginId = *response.Response.PluginId
	d.SetId(pluginId)

	return resourceTencentCloudMonitorGrafanaPluginRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaPluginRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_plugin.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	grafanaPluginId := d.Id()

	grafanaPlugin, err := service.DescribeMonitorGrafanaPluginById(ctx, pluginId)
	if err != nil {
		return err
	}

	if grafanaPlugin == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaPlugin` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if grafanaPlugin.InstanceId != nil {
		_ = d.Set("instance_id", grafanaPlugin.InstanceId)
	}

	if grafanaPlugin.PluginId != nil {
		_ = d.Set("plugin_id", grafanaPlugin.PluginId)
	}

	if grafanaPlugin.Version != nil {
		_ = d.Set("version", grafanaPlugin.Version)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaPluginUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_plugin.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		uninstallGrafanaPluginsRequest  = monitor.NewUninstallGrafanaPluginsRequest()
		uninstallGrafanaPluginsResponse = monitor.NewUninstallGrafanaPluginsResponse()
	)

	grafanaPluginId := d.Id()

	request.PluginId = &pluginId

	immutableArgs := []string{"instance_id", "plugin_id", "version"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UninstallGrafanaPlugins(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaPlugin failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorGrafanaPluginRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaPluginDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_plugin.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	grafanaPluginId := d.Id()

	if err := service.DeleteMonitorGrafanaPluginById(ctx, pluginId); err != nil {
		return err
	}

	return nil
}
