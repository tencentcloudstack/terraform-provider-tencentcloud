/*
Provides a resource to create a monitor tmp_exporter_integration

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_exporter_integration" "tmp_exporter_integration" {
  instance_id = "prom-dko9d0nu"
  kind = "blackbox-exporter"
  content = "blackbox-exporter"
  kube_type = 1
  cluster_id = "job_name: demo-config"
}
```

Import

monitor tmp_exporter_integration can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_tmp_exporter_integration.tmp_exporter_integration tmp_exporter_integration_id
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

func resourceTencentCloudMonitorTmpExporterIntegration() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpExporterIntegrationCreate,
		Read:   resourceTencentCloudMonitorTmpExporterIntegrationRead,
		Update: resourceTencentCloudMonitorTmpExporterIntegrationUpdate,
		Delete: resourceTencentCloudMonitorTmpExporterIntegrationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"kind": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Type.",
			},

			"content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Integration config.",
			},

			"kube_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Integration config.",
			},

			"cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpExporterIntegrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_exporter_integration.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewCreateExporterIntegrationRequest()
		response = monitor.NewCreateExporterIntegrationResponse()
		jobId    string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kind"); ok {
		request.Kind = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("kube_type"); ok {
		request.KubeType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreateExporterIntegration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpExporterIntegration failed, reason:%+v", logId, err)
		return err
	}

	jobId = *response.Response.JobId
	d.SetId(jobId)

	return resourceTencentCloudMonitorTmpExporterIntegrationRead(d, meta)
}

func resourceTencentCloudMonitorTmpExporterIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_exporter_integration.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	tmpExporterIntegrationId := d.Id()

	tmpExporterIntegration, err := service.DescribeMonitorTmpExporterIntegrationById(ctx, jobId)
	if err != nil {
		return err
	}

	if tmpExporterIntegration == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpExporterIntegration` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if tmpExporterIntegration.InstanceId != nil {
		_ = d.Set("instance_id", tmpExporterIntegration.InstanceId)
	}

	if tmpExporterIntegration.Kind != nil {
		_ = d.Set("kind", tmpExporterIntegration.Kind)
	}

	if tmpExporterIntegration.Content != nil {
		_ = d.Set("content", tmpExporterIntegration.Content)
	}

	if tmpExporterIntegration.KubeType != nil {
		_ = d.Set("kube_type", tmpExporterIntegration.KubeType)
	}

	if tmpExporterIntegration.ClusterId != nil {
		_ = d.Set("cluster_id", tmpExporterIntegration.ClusterId)
	}

	return nil
}

func resourceTencentCloudMonitorTmpExporterIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_exporter_integration.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateExporterIntegrationRequest()

	tmpExporterIntegrationId := d.Id()

	request.JobId = &jobId

	immutableArgs := []string{"instance_id", "kind", "content", "kube_type", "cluster_id"}

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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateExporterIntegration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor tmpExporterIntegration failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorTmpExporterIntegrationRead(d, meta)
}

func resourceTencentCloudMonitorTmpExporterIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_exporter_integration.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	tmpExporterIntegrationId := d.Id()

	if err := service.DeleteMonitorTmpExporterIntegrationById(ctx, jobId); err != nil {
		return err
	}

	return nil
}
