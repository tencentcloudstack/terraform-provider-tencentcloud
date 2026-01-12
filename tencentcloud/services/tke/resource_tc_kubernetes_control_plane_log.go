package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudKubernetesControlPlaneLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesControlPlaneLogCreate,
		Read:   resourceTencentCloudKubernetesControlPlaneLogRead,
		Update: resourceTencentCloudKubernetesControlPlaneLogUpdate,
		Delete: resourceTencentCloudKubernetesControlPlaneLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID.",
			},

			"cluster_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster type. currently only support tke.",
			},

			"components": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Component name list. currently supports cluster-autoscaler, kapenter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Component name.",
						},
						"log_level": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Log level. for components that support dynamic adjustment, you can specify this parameter when enabling logs.",
						},
						"log_set_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Logset ID. if not specified, auto-create.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Log topic ID. if not specified, auto-create.",
						},
						"topic_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "topic region. this parameter enables cross-region shipping of logs.",
						},
					},
				},
			},

			"delete_log_set_and_topic": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to simultaneously delete the log set and topic. If the log set and topic are used by other collection rules, they will not be deleted. Default is false.",
			},
		},
	}
}

func resourceTencentCloudKubernetesControlPlaneLogCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_control_plane_log.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = tkev20180525.NewEnableControlPlaneLogsRequest()
		clusterId     string
		clusterType   string
		componentName string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = helper.String(v.(string))
		clusterType = v.(string)
	}

	if v, ok := d.GetOk("components"); ok {
		for _, item := range v.([]interface{}) {
			componentsMap := item.(map[string]interface{})
			componentLogConfig := tkev20180525.ComponentLogConfig{}
			if v, ok := componentsMap["name"].(string); ok && v != "" {
				componentLogConfig.Name = helper.String(v)
				componentName = v
			}

			if v, ok := componentsMap["log_level"].(int); ok {
				componentLogConfig.LogLevel = helper.IntInt64(v)
			}

			if v, ok := componentsMap["log_set_id"].(string); ok && v != "" {
				componentLogConfig.LogSetId = helper.String(v)
			}

			if v, ok := componentsMap["topic_id"].(string); ok && v != "" {
				componentLogConfig.TopicId = helper.String(v)
			}

			if v, ok := componentsMap["topic_region"].(string); ok && v != "" {
				componentLogConfig.TopicRegion = helper.String(v)
			}

			request.Components = append(request.Components, &componentLogConfig)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().EnableControlPlaneLogsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create kubernetes control plane log failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(strings.Join([]string{clusterId, clusterType, componentName}, tccommon.FILED_SP))

	// wait
	waitReq := tkev20180525.NewDescribeControlPlaneLogsRequest()
	waitReq.ClusterId = request.ClusterId
	waitReq.ClusterType = request.ClusterType

	waitErr := resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DescribeControlPlaneLogsWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Details == nil || len(result.Response.Details) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe kubernetes control plane logs failed, Response is nil."))
		}

		for _, item := range result.Response.Details {
			if item != nil && item.Name != nil && *item.Name == componentName {
				return nil
			}
		}

		return resource.RetryableError(fmt.Errorf("wait for kubernetes control plane log creating..."))
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s create kubernetes control plane log failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return resourceTencentCloudKubernetesControlPlaneLogRead(d, meta)
}

func resourceTencentCloudKubernetesControlPlaneLogRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_control_plane_log.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	clusterType := idSplit[1]
	componentName := idSplit[2]

	respData, err := service.DescribeKubernetesControlPlaneLogById(ctx, clusterId, clusterType, componentName)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_control_plane_log` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("cluster_id", clusterId)
	_ = d.Set("cluster_type", clusterType)

	dMap := map[string]interface{}{}
	if respData.Name != nil {
		dMap["name"] = respData.Name
	}

	if respData.LogLevel != nil {
		dMap["log_level"] = respData.LogLevel
	}

	if respData.LogSetId != nil {
		dMap["log_set_id"] = respData.LogSetId
	}

	if respData.TopicId != nil {
		dMap["topic_id"] = respData.TopicId
	}

	if respData.TopicRegion != nil {
		dMap["topic_region"] = respData.TopicRegion
	}

	_ = d.Set("components", []interface{}{dMap})

	return nil
}

func resourceTencentCloudKubernetesControlPlaneLogUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_control_plane_log.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	// if d.HasChange("delete_log_set_and_topic") {
	// 	if v, ok := d.GetOkExists("delete_log_set_and_topic"); ok {
	// 		_ = d.Set("delete_log_set_and_topic", v.(bool))
	// 	}
	// }

	return resourceTencentCloudKubernetesControlPlaneLogRead(d, meta)
}

func resourceTencentCloudKubernetesControlPlaneLogDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_control_plane_log.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tkev20180525.NewDisableControlPlaneLogsRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	clusterType := idSplit[1]
	componentName := idSplit[2]

	if v, ok := d.GetOkExists("delete_log_set_and_topic"); ok {
		request.DeleteLogSetAndTopic = helper.Bool(v.(bool))
	}

	request.ClusterId = &clusterId
	request.ClusterType = &clusterType
	request.ComponentNames = append(request.ComponentNames, &componentName)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DisableControlPlaneLogsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete kubernetes control plane log failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitReq := tkev20180525.NewDescribeControlPlaneLogsRequest()
	waitReq.ClusterId = request.ClusterId
	waitReq.ClusterType = request.ClusterType

	waitErr := resource.Retry(3*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DescribeControlPlaneLogsWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe kubernetes control plane logs failed, Response is nil."))
		}

		if result.Response.Details == nil || len(result.Response.Details) == 0 {
			return nil
		}

		var has bool
		for _, item := range result.Response.Details {
			if item != nil && item.Name != nil && *item.Name == componentName {
				has = true
			}
		}

		if has {
			resource.RetryableError(fmt.Errorf("wait for kubernetes control plane log deleteing..."))
		}

		return nil
	})

	if waitErr != nil {
		log.Printf("[CRITAL]%s delete kubernetes control plane log failed, reason:%+v", logId, waitErr)
		return waitErr
	}

	return nil
}
