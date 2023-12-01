/*
Provides a resource to create a monitor grafanaPlugin

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_grafana_instance" "foo" {
  instance_name         = "test-grafana"
  vpc_id                = tencentcloud_vpc.vpc.id
  subnet_ids            = [tencentcloud_subnet.subnet.id]
  grafana_init_password = "1234567890"
  enable_internet = false

  tags = {
    "createdBy" = "test"
  }
}

resource "tencentcloud_monitor_grafana_plugin" "grafanaPlugin" {
  instance_id = tencentcloud_monitor_grafana_instance.foo.id
  plugin_id   = "grafana-piechart-panel"
  version     = "1.6.2"
}

```
Import

monitor grafanaPlugin can be imported using the instance_id#plugin_id, e.g.
```
$ terraform import tencentcloud_monitor_grafana_plugin.grafanaPlugin grafana-50nj6v00#grafana-piechart-panel
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorGrafanaPlugin() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorGrafanaPluginRead,
		Create: resourceTencentCloudMonitorGrafanaPluginCreate,
		Update: resourceTencentCloudMonitorGrafanaPluginUpdate,
		Delete: resourceTencentCloudMonitorGrafanaPluginDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Grafana instance id.",
			},

			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Plugin id.",
			},

			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
		request = monitor.NewInstallPluginsRequest()
		//response   *monitor.InstallPluginsResponse
		pluginId     string
		instanceId   string
		descResquest = monitor.NewDescribeInstalledPluginsRequest()
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	var plugin monitor.GrafanaPlugin
	if v, ok := d.GetOk("plugin_id"); ok {
		pluginId = v.(string)
		plugin.PluginId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("version"); ok {
		plugin.Version = helper.String(v.(string))
	}
	request.Plugins = append(request.Plugins, &plugin)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().InstallPlugins(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		//response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor grafanaPlugin failed, reason:%+v", logId, err)
		return err
	}

	descResquest.PluginId = &pluginId
	descResquest.InstanceId = &instanceId
	outErr := resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().DescribeInstalledPlugins(descResquest)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err, InternalError)
		}
		if len(response.Response.PluginSet) < 1 {
			return resource.RetryableError(fmt.Errorf("Installing pluin %v, retry...", pluginId))
		}
		return nil
	})
	if outErr != nil {
		log.Printf("[CRITAL]%s Inquire monitor grafanaPlugin failed, reason:%+v", logId, outErr)
		return outErr
	}

	d.SetId(strings.Join([]string{instanceId, pluginId}, FILED_SP))
	return resourceTencentCloudMonitorGrafanaPluginRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaPluginRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_plugin.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	pluginId := idSplit[1]

	grafanaPlugin, err := service.DescribeMonitorGrafanaPlugin(ctx, instanceId, pluginId)

	if err != nil {
		return err
	}

	if grafanaPlugin == nil {
		d.SetId("")
		return fmt.Errorf("resource `grafanaPlugin` %s does not exist", pluginId)
	}

	_ = d.Set("instance_id", instanceId)

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

	request := monitor.NewUninstallGrafanaPluginsRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	pluginId := idSplit[1]

	request.InstanceId = &instanceId
	request.PluginIds = []*string{&pluginId}

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("plugin_id") {
		return fmt.Errorf("`plugin_id` do not support change now.")
	}

	err := resourceTencentCloudMonitorGrafanaPluginDelete(d, meta)
	if err != nil {
		return err
	}

	return resourceTencentCloudMonitorGrafanaPluginCreate(d, meta)
}

func resourceTencentCloudMonitorGrafanaPluginDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_plugin.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	pluginId := idSplit[1]

	if err := service.DeleteMonitorGrafanaPluginById(ctx, instanceId, pluginId); err != nil {
		return err
	}

	return nil
}
