package tmp

import (
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpExporterIntegrationV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpExporterIntegrationV2Create,
		Read:   resourceTencentCloudMonitorTmpExporterIntegrationV2Read,
		Update: resourceTencentCloudMonitorTmpExporterIntegrationV2Update,
		Delete: resourceTencentCloudMonitorTmpExporterIntegrationV2Delete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance ID.",
			},

			"kind": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Integration config. For more details, please refer to [Cloud Monitoring](https://www.tencentcloud.com/document/product/248/63002?lang=en&pg=).",
			},

			"kube_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2, 3}),
				Description:  "Integration config. 1 - TKE; 2 - EKS; 3 - MEKS.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cluster ID.",
			},

			"disable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Integration is disabled when the value is true. Default is false.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpExporterIntegrationV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_exporter_integration_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		request      = monitor.NewCreateExporterIntegrationRequest()
		response     = monitor.NewCreateExporterIntegrationResponse()
		instanceId   string
		kubeType     int
		clusterId    string
		kind         string
		hasKubeType  bool
		hasClusterId bool
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("kind"); ok {
		request.Kind = helper.String(v.(string))
		kind = v.(string)
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("kube_type"); ok {
		request.KubeType = helper.IntInt64(v.(int))
		kubeType = v.(int)
		hasKubeType = true
	}

	if v, ok := d.GetOk("cluster_id"); ok && v.(string) != "" {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
		hasClusterId = true
	}

	if hasKubeType != hasClusterId {
		return fmt.Errorf("`kube_type` and `cluster_id` can only be set together or not set at all.")
	}

	initStatus := monitor.NewDescribePrometheusInstanceInitStatusRequest()
	initStatus.InstanceId = &instanceId
	err := resource.Retry(8*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		results, errRet := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().DescribePrometheusInstanceInitStatus(initStatus)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}

		if results == nil || results.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("prometheusInstanceInit results is nil, operate failed"))
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
			iniRequest.InstanceId = &instanceId
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().RunPrometheusInstance(iniRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, iniRequest.GetAction(), iniRequest.ToJsonString(), result.ToJsonString())
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

	if hasKubeType || hasClusterId {
		d.SetId(strings.Join([]string{tmpExporterIntegrationId, instanceId, strconv.Itoa(kubeType), clusterId, kind}, tccommon.FILED_SP))
	} else {
		d.SetId(strings.Join([]string{tmpExporterIntegrationId, instanceId, kind}, tccommon.FILED_SP))
	}

	// wait
	waitReq := monitor.NewDescribeExporterIntegrationsRequest()
	waitReq.InstanceId = request.InstanceId
	waitReq.KubeType = request.KubeType
	waitReq.ClusterId = request.ClusterId
	waitReq.Kind = request.Kind
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().DescribeExporterIntegrations(waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.IntegrationSet == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeExporterIntegrations is nil"))
		}

		integration := result.Response.IntegrationSet[0]
		if integration.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Status is nil"))
		}

		if *integration.Status == 2 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Create exporter integration status is %d", *integration.Status))
	})

	if err != nil {
		return err
	}

	// set status
	if v, ok := d.GetOkExists("disable"); ok {
		if v.(bool) {
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

			if v, ok := d.GetOkExists("kube_type"); ok {
				request.KubeType = helper.IntInt64(v.(int))
			}

			if v, ok := d.GetOk("cluster_id"); ok && v.(string) != "" {
				request.ClusterId = helper.String(v.(string))
			}

			if v, ok := d.GetOkExists("disable"); ok {
				request.Disable = helper.Bool(v.(bool))
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

			// wait
			waitReq := monitor.NewDescribeExporterIntegrationsRequest()
			waitReq.InstanceId = request.InstanceId
			waitReq.KubeType = request.KubeType
			waitReq.ClusterId = request.ClusterId
			waitReq.Kind = request.Kind
			err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().DescribeExporterIntegrations(waitReq)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
				}

				if result == nil || result.Response == nil || result.Response.IntegrationSet == nil {
					return resource.NonRetryableError(fmt.Errorf("DescribeExporterIntegrations is nil"))
				}

				integration := result.Response.IntegrationSet[0]
				if integration.Status == nil {
					return resource.NonRetryableError(fmt.Errorf("Status is nil"))
				}

				if *integration.Status == 5 {
					return nil
				}

				return resource.RetryableError(fmt.Errorf("update exporter integration status is %d", *integration.Status))
			})

			if err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudMonitorTmpExporterIntegrationV2Read(d, meta)
}

func resourceTencentCloudMonitorTmpExporterIntegrationV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_exporter_integration_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                    = tccommon.GetLogId(tccommon.ContextNil)
		ctx                      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service                  = svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		tmpExporterIntegrationId = d.Id()
	)

	tmpExporterIntegration, err := service.DescribeMonitorTmpExporterIntegration(ctx, tmpExporterIntegrationId)
	if err != nil {
		return err
	}

	if tmpExporterIntegration == nil {
		log.Printf("[WARN]%s resource `tencentcloud_monitor_tmp_exporter_integration_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if tmpExporterIntegration.Kind != nil {
		_ = d.Set("kind", tmpExporterIntegration.Kind)
	}

	if tmpExporterIntegration.Content != nil {
		_ = d.Set("content", tmpExporterIntegration.Content)
	}

	if tmpExporterIntegration.Status != nil {
		switch {
		case *tmpExporterIntegration.Status == 2:
			_ = d.Set("disable", false)
		case *tmpExporterIntegration.Status == 5:
			_ = d.Set("disable", true)
		}
	}

	return nil
}

func resourceTencentCloudMonitorTmpExporterIntegrationV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_exporter_integration_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = monitor.NewUpdateExporterIntegrationRequest()
	)

	immutableArgs := []string{"instance_id", "kind", "kube_type", "cluster_id"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed.", v)
		}
	}

	var disableFlag bool
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

	if v, ok := d.GetOk("cluster_id"); ok && v.(string) != "" {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("disable"); ok {
		request.Disable = helper.Bool(v.(bool))
		disableFlag = v.(bool)
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

	// wait
	waitReq := monitor.NewDescribeExporterIntegrationsRequest()
	waitReq.InstanceId = request.InstanceId
	waitReq.KubeType = request.KubeType
	waitReq.ClusterId = request.ClusterId
	waitReq.Kind = request.Kind
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().DescribeExporterIntegrations(waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.IntegrationSet == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeExporterIntegrations is nil"))
		}

		integration := result.Response.IntegrationSet[0]
		if integration.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Status is nil"))
		}

		if disableFlag && *integration.Status == 5 {
			return nil
		} else if !disableFlag && *integration.Status == 2 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("update exporter integration status is %d", *integration.Status))
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudMonitorTmpExporterIntegrationV2Read(d, meta)
}

func resourceTencentCloudMonitorTmpExporterIntegrationV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_exporter_integration_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                    = tccommon.GetLogId(tccommon.ContextNil)
		ctx                      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service                  = svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		tmpExporterIntegrationId = d.Id()
	)

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
