/*
Provides a resource to create a monitor grafana_sso_config

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_sso_config" "grafana_sso_config" {
  instance_id = "grafana-dp2hnnfa"
  enable_sso  = false
}
```

Import

monitor grafana_sso_config can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_sso_config.grafana_sso_config instance_id
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

func resourceTencentCloudMonitorGrafanaSsoConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaSsoConfigCreate,
		Read:   resourceTencentCloudMonitorGrafanaSsoConfigRead,
		Update: resourceTencentCloudMonitorGrafanaSsoConfigUpdate,
		Delete: resourceTencentCloudMonitorGrafanaSsoConfigDelete,
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

			"enable_sso": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable SSO: `true` for enabling; `false` for disabling.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaSsoConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_sso_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorGrafanaSsoConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorGrafanaSsoConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_sso_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	grafanaSsoConfig, err := service.DescribeMonitorGrafanaInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	if grafanaSsoConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaSsoConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if grafanaSsoConfig.EnableSSO != nil {
		_ = d.Set("enable_sso", grafanaSsoConfig.EnableSSO)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaSsoConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_sso_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := monitor.NewEnableGrafanaSSORequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOkExists("enable_sso"); ok {
		request.EnableSSO = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().EnableGrafanaSSO(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaSsoConfig failed, reason:%+v", logId, err)
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
			return resource.NonRetryableError(fmt.Errorf("grafanaInstance status is %v, update sso failed.", *instance.InstanceStatus))
		}
		return resource.RetryableError(fmt.Errorf("grafanaInstance status is %v, retry...", *instance.InstanceStatus))
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudMonitorGrafanaSsoConfigRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaSsoConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_sso_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
