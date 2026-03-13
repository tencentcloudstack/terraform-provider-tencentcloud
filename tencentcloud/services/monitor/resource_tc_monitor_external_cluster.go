package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorExternalCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorExternalClusterCreate,
		Read:   resourceTencentCloudMonitorExternalClusterRead,
		Delete: resourceTencentCloudMonitorExternalClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"cluster_region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster region.",
			},

			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Cluster name.",
			},

			"external_labels": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "External labels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Label name.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Label value.",
						},
					},
				},
			},

			"enable_external": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to enable public network access.",
			},

			// computed
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster ID.",
			},

			"cluster_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster type, returned by API.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Agent status.",
			},
		},
	}
}

func resourceTencentCloudMonitorExternalClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_external_cluster.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = monitor.NewCreateExternalClusterRequest()
		response   = monitor.NewCreateExternalClusterResponse()
		instanceId string
		clusterId  string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("cluster_region"); ok {
		request.ClusterRegion = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("external_labels"); ok {
		for _, item := range v.([]interface{}) {
			labelMap := item.(map[string]interface{})
			label := monitor.Label{}
			if v, ok := labelMap["name"].(string); ok {
				label.Name = helper.String(v)
			}

			if v, ok := labelMap["value"].(string); ok && v != "" {
				label.Value = helper.String(v)
			}

			request.ExternalLabels = append(request.ExternalLabels, &label)
		}
	}

	if v, ok := d.GetOkExists("enable_external"); ok {
		request.EnableExternal = helper.Bool(v.(bool))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreateExternalClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create monitor external cluster failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create monitor external cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ClusterId == nil {
		return fmt.Errorf("ClusterId is nil.")
	}

	clusterId = *response.Response.ClusterId
	d.SetId(strings.Join([]string{instanceId, clusterId}, tccommon.FILED_SP))

	// Wait for the cluster to be ready
	service := NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	err := resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		cluster, e := service.DescribeMonitorExternalClusterById(ctx, instanceId, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if cluster == nil {
			return resource.NonRetryableError(fmt.Errorf("cluster not found after creation"))
		}
		if cluster.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("cluster status is nil"))
		}

		status := *cluster.Status
		if status == "waiting" || status == "normal" {
			log.Printf("[DEBUG]%s cluster status is %s, creation completed\n", logId, status)
			return nil
		}

		log.Printf("[DEBUG]%s cluster status is %s, waiting for completion\n", logId, status)
		return resource.RetryableError(fmt.Errorf("cluster status is %s, expected 'waiting' or 'normal'", status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s wait for monitor external cluster ready failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorExternalClusterRead(d, meta)
}

func resourceTencentCloudMonitorExternalClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_external_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	clusterId := idSplit[1]

	respData, err := service.DescribeMonitorExternalClusterById(ctx, instanceId, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_monitor_external_cluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if respData.ClusterName != nil {
		_ = d.Set("cluster_name", respData.ClusterName)
	}

	if respData.Region != nil {
		_ = d.Set("cluster_region", respData.Region)
	}

	if len(respData.ExternalLabels) > 0 {
		externalLabelsList := make([]map[string]interface{}, 0, len(respData.ExternalLabels))
		for _, label := range respData.ExternalLabels {
			labelMap := make(map[string]interface{})
			if label.Name != nil {
				labelMap["name"] = *label.Name
			}
			if label.Value != nil {
				labelMap["value"] = *label.Value
			}
			externalLabelsList = append(externalLabelsList, labelMap)
		}
		_ = d.Set("external_labels", externalLabelsList)
	}

	if respData.EnableExternal != nil {
		_ = d.Set("enable_external", respData.EnableExternal)
	}

	if respData.ClusterId != nil {
		_ = d.Set("cluster_id", respData.ClusterId)
	}

	if respData.ClusterType != nil {
		_ = d.Set("cluster_type", respData.ClusterType)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	return nil
}

func resourceTencentCloudMonitorExternalClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_external_cluster.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = monitor.NewDeletePrometheusClusterAgentRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	clusterId := idSplit[1]

	request.InstanceId = &instanceId

	// Get cluster_type from state
	clusterType, ok := d.GetOk("cluster_type")
	if !ok || clusterType.(string) == "" {
		return fmt.Errorf("cluster_type not found in state, cannot delete cluster")
	}

	agentInfo := monitor.PrometheusAgentInfo{
		ClusterId:   &clusterId,
		ClusterType: helper.String(clusterType.(string)),
	}
	request.Agents = append(request.Agents, &agentInfo)

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().DeletePrometheusClusterAgentWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete monitor external cluster failed, Response is nil."))
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete monitor external cluster failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Wait for the cluster to be deleted (max 15 minutes)
	service := NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	err := resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		cluster, e := service.DescribeMonitorExternalClusterById(ctx, instanceId, clusterId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		// If cluster is not found (nil), deletion is complete
		if cluster == nil {
			log.Printf("[DEBUG]%s cluster not found, deletion completed\n", logId)
			return nil
		}

		log.Printf("[DEBUG]%s cluster still exists, waiting for deletion\n", logId)
		return resource.RetryableError(fmt.Errorf("cluster still exists, waiting for deletion"))
	})

	if err != nil {
		log.Printf("[CRITAL]%s wait for monitor external cluster deletion failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
