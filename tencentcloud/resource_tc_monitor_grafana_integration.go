/*
Provides a resource to create a monitor grafana_integration

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_integration" "grafana_integration" {
  instance_id = &lt;nil&gt;
    kind = &lt;nil&gt;
  content = &lt;nil&gt;
}
```

Import

monitor grafana_integration can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_integration.grafana_integration grafana_integration_id
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
	"strings"
)

func resourceTencentCloudMonitorGrafanaIntegration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaIntegrationCreate,
		Read:   resourceTencentCloudMonitorGrafanaIntegrationRead,
		Update: resourceTencentCloudMonitorGrafanaIntegrationUpdate,
		Delete: resourceTencentCloudMonitorGrafanaIntegrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance id.",
			},

			"integration_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Integration id.",
			},

			"kind": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Integration json schema kind.",
			},

			"content": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Generated json string of given integration json schema.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaIntegrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_integration.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = monitor.NewCreateGrafanaIntegrationRequest()
		response      = monitor.NewCreateGrafanaIntegrationResponse()
		integrationId string
		kind          string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kind"); ok {
		kind = v.(string)
		request.Kind = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreateGrafanaIntegration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create monitor grafanaIntegration failed, reason:%+v", logId, err)
		return err
	}

	integrationId = *response.Response.IntegrationId
	d.SetId(strings.Join([]string{integrationId, kind}, FILED_SP))

	return resourceTencentCloudMonitorGrafanaIntegrationRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_integration.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	integrationId := idSplit[0]
	kind := idSplit[1]

	grafanaIntegration, err := service.DescribeMonitorGrafanaIntegrationById(ctx, integrationId, kind)
	if err != nil {
		return err
	}

	if grafanaIntegration == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaIntegration` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if grafanaIntegration.InstanceId != nil {
		_ = d.Set("instance_id", grafanaIntegration.InstanceId)
	}

	if grafanaIntegration.IntegrationId != nil {
		_ = d.Set("integration_id", grafanaIntegration.IntegrationId)
	}

	if grafanaIntegration.Kind != nil {
		_ = d.Set("kind", grafanaIntegration.Kind)
	}

	if grafanaIntegration.Content != nil {
		_ = d.Set("content", grafanaIntegration.Content)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_integration.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateGrafanaIntegrationRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	integrationId := idSplit[0]
	kind := idSplit[1]

	request.IntegrationId = &integrationId
	request.Kind = &kind

	immutableArgs := []string{"instance_id", "integration_id", "kind", "content"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateGrafanaIntegration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaIntegration failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorGrafanaIntegrationRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_integration.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	integrationId := idSplit[0]
	kind := idSplit[1]

	if err := service.DeleteMonitorGrafanaIntegrationById(ctx, integrationId, kind); err != nil {
		return err
	}

	return nil
}
