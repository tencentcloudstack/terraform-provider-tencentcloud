/*
Provides a resource to create a monitor tmpExporterIntegration

Example Usage

```hcl
resource "tencentcloud_monitor_tmp_exporter_integration" "tmpExporterIntegration" {
  instance_id = "prom-dko9d0nu"
  kind = "blackbox-exporter"
  content = "{\"name\":\"test\",\"kind\":\"blackbox-exporter\",\"spec\":{\"instanceSpec\":{\"module\":\"http_get\",\"urls\":[\"xx\"]}}}"
  kube_type = 1
  cluster_id = "cls-bmuaukfu"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpExporterIntegration() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpExporterIntegrationRead,
		Create: resourceTencentCloudMonitorTmpExporterIntegrationCreate,
		Update: resourceTencentCloudMonitorTmpExporterIntegrationUpdate,
		Delete: resourceTencentCloudMonitorTmpExporterIntegrationDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},

			"kind": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Integration config.",
			},

			"kube_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Integration config.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
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
		instanceId string
		kubeType   int
		clusterId  string
		kind       string
	)

	var (
		request  = monitor.NewCreateExporterIntegrationRequest()
		response *monitor.CreateExporterIntegrationResponse
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}

	if v, ok := d.GetOk("kind"); ok {
		kind = v.(string)
		request.Kind = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kube_type"); ok {
		kubeType = v.(int)
		request.KubeType = helper.IntInt64(kubeType)
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(clusterId)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreateExporterIntegration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpExporterIntegration failed, reason:%+v", logId, err)
		return err
	}

	tmpExporterIntegrationId := *response.Response.Names[0]

	d.SetId(strings.Join([]string{tmpExporterIntegrationId, instanceId, strconv.Itoa(kubeType), clusterId, kind}, FILED_SP))

	return resourceTencentCloudMonitorTmpExporterIntegrationRead(d, meta)
}

func resourceTencentCloudMonitorTmpExporterIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmpExporterIntegration.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	tmpExporterIntegrationId := d.Id()

	tmpExporterIntegration, err := service.DescribeMonitorTmpExporterIntegration(ctx, tmpExporterIntegrationId)

	if err != nil {
		return err
	}

	if tmpExporterIntegration == nil {
		d.SetId("")
		return fmt.Errorf("resource `tmpExporterIntegration` %s does not exist", tmpExporterIntegrationId)
	}

	if tmpExporterIntegration.Kind != nil {
		_ = d.Set("kind", tmpExporterIntegration.Kind)
	}

	if tmpExporterIntegration.Content != nil {
		_ = d.Set("content", tmpExporterIntegration.Content)
	}

	return nil
}

func resourceTencentCloudMonitorTmpExporterIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_exporter_integration.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateExporterIntegrationRequest()

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kind"); ok {
		request.Kind = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kube_type"); ok {
		request.KubeType = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateExporterIntegration(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
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

	if err := service.DeleteMonitorTmpExporterIntegrationById(ctx, tmpExporterIntegrationId); err != nil {
		return err
	}

	err := resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		tmpExporterIntegration, errRet := service.DescribeMonitorTmpExporterIntegration(ctx, tmpExporterIntegrationId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if tmpExporterIntegration == nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("exporter integration status is %v, retry...", *tmpExporterIntegration.Status))
	})
	if err != nil {
		return err
	}

	return nil
}
