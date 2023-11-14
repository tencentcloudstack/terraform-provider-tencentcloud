/*
Provides a resource to create a monitor grafana_env

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_env" "grafana_env" {
  instance_id = "grafana-12345678"
  envs = ""
}
```

Import

monitor grafana_env can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_env.grafana_env grafana_env_id
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

func resourceTencentCloudMonitorGrafanaEnv() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaEnvCreate,
		Read:   resourceTencentCloudMonitorGrafanaEnvRead,
		Update: resourceTencentCloudMonitorGrafanaEnvUpdate,
		Delete: resourceTencentCloudMonitorGrafanaEnvDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance ID.",
			},

			"envs": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment variables.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaEnvCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_env.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorGrafanaEnvUpdate(d, meta)
}

func resourceTencentCloudMonitorGrafanaEnvRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_env.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	grafanaEnvId := d.Id()

	grafanaEnv, err := service.DescribeMonitorGrafanaEnvById(ctx, instanceId)
	if err != nil {
		return err
	}

	if grafanaEnv == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaEnv` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if grafanaEnv.InstanceId != nil {
		_ = d.Set("instance_id", grafanaEnv.InstanceId)
	}

	if grafanaEnv.Envs != nil {
		_ = d.Set("envs", grafanaEnv.Envs)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaEnvUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_env.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateGrafanaEnvironmentsRequest()

	grafanaEnvId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "envs"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("envs") {
		if v, ok := d.GetOk("envs"); ok {
			request.Envs = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateGrafanaEnvironments(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaEnv failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorGrafanaEnvRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaEnvDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_env.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
