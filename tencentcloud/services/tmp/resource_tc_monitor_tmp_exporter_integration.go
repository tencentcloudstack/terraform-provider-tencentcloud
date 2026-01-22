package tmp

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpExporterIntegration() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.81.182. Please use 'tencentcloud_monitor_tmp_exporter_integration_v2' instead.",
		Read:               resourceTencentCloudMonitorTmpExporterIntegrationRead,
		Create:             resourceTencentCloudMonitorTmpExporterIntegrationCreate,
		Update:             resourceTencentCloudMonitorTmpExporterIntegrationUpdate,
		Delete:             resourceTencentCloudMonitorTmpExporterIntegrationDelete,
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
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_exporter_integration.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	initStatus := monitor.NewDescribePrometheusInstanceInitStatusRequest()
	initStatus.InstanceId = request.InstanceId
	err := resource.Retry(8*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, errRet := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().DescribePrometheusInstanceInitStatus(initStatus)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		status := results.Response.Status
		if status == nil {
			return resource.NonRetryableError(fmt.Errorf("prometheusInstanceInit status is nil, operate failed"))
		}
		if *status == "running" {
			return nil
		}
		if *status == "uninitialized" {
			iniRequest := monitor.NewRunPrometheusInstanceRequest()
			iniRequest.InstanceId = request.InstanceId
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().RunPrometheusInstance(iniRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if err != nil {
				return resource.RetryableError(fmt.Errorf("prometheusInstanceInit error %v, operate failed", err))
			}
			return resource.RetryableError(fmt.Errorf("prometheusInstance initializing, retry..."))
		}
		return resource.RetryableError(fmt.Errorf("prometheusInstanceInit status is %v, retry...", *status))
	})
	if err != nil {
		return err
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreateExporterIntegration(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Names == nil {
			return resource.NonRetryableError(fmt.Errorf("Create monitor tmpExporterIntegration failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpExporterIntegration failed, reason:%+v", logId, err)
		return err
	}

	if len(response.Response.Names) < 1 {
		return fmt.Errorf("Names is nil.")
	}

	tmpExporterIntegrationId := *response.Response.Names[0]

	d.SetId(strings.Join([]string{tmpExporterIntegrationId, instanceId, strconv.Itoa(kubeType), clusterId, kind}, tccommon.FILED_SP))

	return resourceTencentCloudMonitorTmpExporterIntegrationRead(d, meta)
}

func resourceTencentCloudMonitorTmpExporterIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_exporter_integration.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	tmpExporterIntegrationId := d.Id()

	tmpExporterIntegration, err := service.DescribeMonitorTmpExporterIntegration(ctx, tmpExporterIntegrationId)

	if err != nil {
		return err
	}

	if tmpExporterIntegration == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `tmpExporterIntegration` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
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
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_exporter_integration.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdateExporterIntegration(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_exporter_integration.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	tmpExporterIntegrationId := d.Id()

	if err := service.DeleteMonitorTmpExporterIntegrationById(ctx, tmpExporterIntegrationId); err != nil {
		return err
	}

	err := resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		tmpExporterIntegration, errRet := service.DescribeMonitorTmpExporterIntegration(ctx, tmpExporterIntegrationId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
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
