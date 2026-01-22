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

func ResourceTencentCloudKubernetesServerlessNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesServerlessNodePoolCreate,
		Read:   resourceTencentCloudKubernetesServerlessNodePoolRead,
		Update: resourceTencentCloudKubernetesServerlessNodePoolUpdate,
		Delete: resourceTencentCloudKubernetesServerlessNodePoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "cluster id of serverless node pool.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "serverless node pool name.",
			},

			"serverless_nodes": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "node list of serverless node pool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"display_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "display name of serverless node.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "subnet id of serverless node.",
						},
					},
				},
			},

			"security_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "security groups of serverless node pool.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "labels of serverless node.",
			},

			"taints": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "taints of serverless node.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key of the taint. The taint key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the taint.",
						},
						"effect": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Effect of the taint. Valid values are: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.",
						},
					},
				},
			},

			"life_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "life state of serverless node pool.",
			},
		},
	}
}

func resourceTencentCloudKubernetesServerlessNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_serverless_node_pool.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = tkev20180525.NewCreateClusterVirtualNodePoolRequest()
		response   = tkev20180525.NewCreateClusterVirtualNodePoolResponse()
		clusterId  string
		nodePoolId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(clusterId)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("taints"); ok {
		for _, item := range v.([]interface{}) {
			taintsMap := item.(map[string]interface{})
			taint := tkev20180525.Taint{}
			if v, ok := taintsMap["key"]; ok {
				taint.Key = helper.String(v.(string))
			}

			if v, ok := taintsMap["value"]; ok {
				taint.Value = helper.String(v.(string))
			}

			if v, ok := taintsMap["effect"]; ok {
				taint.Effect = helper.String(v.(string))
			}

			request.Taints = append(request.Taints, &taint)
		}
	}

	if v, ok := d.GetOk("serverless_nodes"); ok {
		for _, item := range v.([]interface{}) {
			virtualNodesMap := item.(map[string]interface{})
			virtualNodeSpec := tkev20180525.VirtualNodeSpec{}
			if v, ok := virtualNodesMap["display_name"]; ok {
				virtualNodeSpec.DisplayName = helper.String(v.(string))
			}

			if v, ok := virtualNodesMap["subnet_id"]; ok {
				virtualNodeSpec.SubnetId = helper.String(v.(string))
			}

			request.VirtualNodes = append(request.VirtualNodes, &virtualNodeSpec)
		}
	}

	if err := resourceTencentCloudKubernetesServerlessNodePoolCreatePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().CreateClusterVirtualNodePoolWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create kubernetes serverless node pool failed, Response is nil"))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes serverless node pool failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.NodePoolId == nil {
		return fmt.Errorf("NodePoolId is nil")
	}

	nodePoolId = *response.Response.NodePoolId
	d.SetId(strings.Join([]string{clusterId, nodePoolId}, tccommon.FILED_SP))

	// wait
	waitRequest := tkev20180525.NewDescribeClusterVirtualNodePoolsRequest()
	waitRequest.ClusterId = &clusterId
	err = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DescribeClusterVirtualNodePoolsWithContext(ctx, waitRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe cluster virtual nodepools failed, Response is nil."))
		}

		if result.Response.NodePoolSet == nil || len(result.Response.NodePoolSet) < 1 {
			return resource.NonRetryableError(fmt.Errorf("NodePoolSet is nil."))
		}

		var hasNpId bool
		for _, item := range result.Response.NodePoolSet {
			if item != nil && item.NodePoolId != nil && *item.NodePoolId == nodePoolId {
				if item.LifeState != nil && *item.LifeState == "normal" {
					return nil
				}

				hasNpId = true
			}
		}

		if !hasNpId {
			return resource.NonRetryableError(fmt.Errorf("NodePoolId %s is not found.", nodePoolId))
		}

		return resource.RetryableError(fmt.Errorf("serverless node pool %s is creating...", nodePoolId))
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudKubernetesServerlessNodePoolRead(d, meta)
}

func resourceTencentCloudKubernetesServerlessNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_serverless_node_pool.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	nodePoolId := idSplit[1]

	_ = d.Set("cluster_id", clusterId)

	var respData *tkev20180525.VirtualNodePool
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesServerlessNodePoolById(ctx, clusterId, nodePoolId)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe cluster virtual nodepools failed, Response is nil."))
		}

		if err := resourceTencentCloudKubernetesServerlessNodePoolReadRequestOnSuccess0(ctx, result); err != nil {
			return err
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read kubernetes serverless node pool failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_serverless_node_pool` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.LifeState != nil {
		_ = d.Set("life_state", respData.LifeState)
	}

	taintsList := make([]map[string]interface{}, 0, len(respData.Taints))
	if respData.Taints != nil {
		for _, taints := range respData.Taints {
			taintsMap := map[string]interface{}{}
			if taints.Key != nil {
				taintsMap["key"] = taints.Key
			}

			if taints.Value != nil {
				taintsMap["value"] = taints.Value
			}

			if taints.Effect != nil {
				taintsMap["effect"] = taints.Effect
			}

			taintsList = append(taintsList, taintsMap)
		}

		_ = d.Set("taints", taintsList)
	}

	if err := resourceTencentCloudKubernetesServerlessNodePoolReadPostHandleResponse0(ctx, respData); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudKubernetesServerlessNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_serverless_node_pool.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	nodePoolId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "taints"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := tkev20180525.NewModifyClusterVirtualNodePoolRequest()
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("taints"); ok {
			for _, item := range v.([]interface{}) {
				taintsMap := item.(map[string]interface{})
				taint := tkev20180525.Taint{}
				if v, ok := taintsMap["key"]; ok {
					taint.Key = helper.String(v.(string))
				}

				if v, ok := taintsMap["value"]; ok {
					taint.Value = helper.String(v.(string))
				}

				if v, ok := taintsMap["effect"]; ok {
					taint.Effect = helper.String(v.(string))
				}

				request.Taints = append(request.Taints, &taint)
			}
		}

		if err := resourceTencentCloudKubernetesServerlessNodePoolUpdatePostFillRequest0(ctx, request); err != nil {
			return err
		}

		request.ClusterId = helper.String(clusterId)
		request.NodePoolId = helper.String(nodePoolId)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterVirtualNodePoolWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes serverless node pool failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudKubernetesServerlessNodePoolRead(d, meta)
}

func resourceTencentCloudKubernetesServerlessNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_serverless_node_pool.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tkev20180525.NewDeleteClusterVirtualNodePoolRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	clusterId := idSplit[0]
	nodePoolId := idSplit[1]

	request.ClusterId = helper.String(clusterId)
	request.NodePoolIds = []*string{helper.String(nodePoolId)}
	request.Force = helper.Bool(true)
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DeleteClusterVirtualNodePoolWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete kubernetes serverless node pool failed, reason:%+v", logId, err)
		return err
	}

	// wait
	waitRequest := tkev20180525.NewDescribeClusterVirtualNodePoolsRequest()
	waitRequest.ClusterId = &clusterId
	err = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DescribeClusterVirtualNodePoolsWithContext(ctx, waitRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.NodePoolSet == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe cluster virtual nodepools failed, Response is nil."))
		}

		var hasNpId bool
		for _, item := range result.Response.NodePoolSet {
			if item.NodePoolId != nil && *item.NodePoolId == nodePoolId {
				hasNpId = true
			}
		}

		if !hasNpId {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("serverless node pool %s is deleting...", nodePoolId))
	})

	if err != nil {
		return err
	}

	return nil
}
