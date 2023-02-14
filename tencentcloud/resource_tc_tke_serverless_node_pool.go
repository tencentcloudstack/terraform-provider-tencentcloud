/*
Provide a resource to create serverless node pool of cluster.

Example Usage
```
resource "tencentcloud_kubernetes_serverless_node_pool" "example_serverless_node_pool" {
	cluster_id = tencentcloud_kubernetes_cluster.example.id
	name = "example_node_pool"
	serverless_nodes = {
		display_name = "serverless_node1"
		subnet_id = "subnet-xxx"
	}
	severless_nodes = {
		display_name = "serverless_node2"
		subnet_id = "subnet-xxx"
	}
	security_group_ids = ["sg-xxx"]
	labels = {
		"example1": "test1",
		"example2": "test2",
	}
}
```

Import

serverless node pool can be imported, e.g.

```
$ terraform import tencentcloud_kubernetes_serverless_node_pool.test cls-xxx#np-xxx
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTkeServerLessNodePool() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Read:   resourceTkeServerlessNodePoolRead,
		Create: resourceTkeServerlessNodePoolCreate,
		Update: resourceTkeServerlessNodePoolUpdate,
		Delete: resourceTkeServerlessNodePoolDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "cluster id of serverless node pool.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "serverless node pool name.",
			},
			"serverless_nodes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
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
				Description: "node list of serverless node pool.",
			},
			"security_group_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				ForceNew:    true,
				Description: "security groups of serverless node pool.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "labels of serverless node.",
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
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
				Description: "taints of serverless node.",
			},
			"life_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "life state of serverless node pool.",
			},
		},
	}
}

func resourceTkeServerlessNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_eks_cluster.read")()

	var (
		items = strings.Split(d.Id(), FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	var (
		nodePool *tke.VirtualNodePool
		has      bool
	)

	outErr := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var err error
		nodePool, has, err = service.DescribeServerlessNodePoolByClusterIdAndNodePoolId(ctx, clusterId, nodePoolId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("serverless node pool %s not exists", d.Id()))
		}
		if shouldServerlessNodePoolRetryReading(*nodePool.LifeState) {
			return resource.RetryableError(fmt.Errorf("serverless node pool %s is now %s, retrying", d.Id(), *nodePool.LifeState))
		}
		return nil
	})

	if outErr != nil {
		return outErr
	}

	if !has {
		d.SetId("")
		return nil
	}

	if err := setDataFromDescribeVirtualNodePoolResponse(clusterId, nodePool, d); err != nil {
		return err
	}

	return nil
}

func resourceTkeServerlessNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_serverless_node_pool.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}

	request := genCreateClusterVirtualNodePoolReq(d)

	nodePoolId, err := service.CreateClusterVirtualNodePool(ctx, request)

	if err != nil {
		return err
	}

	clusterId := *request.ClusterId
	d.SetId(clusterId + FILED_SP + nodePoolId)

	return resourceTkeServerlessNodePoolRead(d, meta)
}
func resourceTkeServerlessNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	// currently only name, labels and taints can be modified
	defer logElapsed("resource.tencentcloud_kubernetes_serverless_node_pool.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		items = strings.Split(d.Id(), FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}

	request := tke.NewModifyClusterVirtualNodePoolRequest()
	request.ClusterId = common.StringPtr(clusterId)
	request.NodePoolId = &nodePoolId

	if d.HasChange("labels") {
		request.Labels = GetOptimizedTkeLabels(d, "labels")
	}

	if d.HasChange("taints") {
		// if taints is empty, need to recreate this resource. But tf need to inform user at applying...
		request.Taints = GetOptimizedTkeTaints(d, "taints")
	}
	if d.HasChange("name") {
		request.Name = common.StringPtr(d.Get("name").(string))
	}

	if err := service.ModifyClusterVirtualNodePool(ctx, request); err != nil {
		return err
	}

	return resourceTkeServerlessNodePoolRead(d, meta)
}
func resourceTkeServerlessNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_serverless_node_pool.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		items = strings.Split(d.Id(), FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	client := meta.(*TencentCloudClient).apiV3Conn
	service := TkeService{client: client}

	request := tke.NewDeleteClusterVirtualNodePoolRequest()
	request.NodePoolIds = []*string{&nodePoolId}
	request.ClusterId = common.StringPtr(clusterId)
	request.Force = common.BoolPtr(true)

	if err := service.DeleteClusterVirtualNodePool(ctx, request); err != nil {
		return err
	}

	return nil
}

func genCreateClusterVirtualNodePoolReq(d *schema.ResourceData) *tke.CreateClusterVirtualNodePoolRequest {
	var (
		clusterId        = d.Get("cluster_id").(string)
		name             = d.Get("name").(string)
		serverlessNodes  = d.Get("serverless_nodes").([]interface{})
		securityGroupIds = d.Get("security_group_ids").([]interface{})
	)

	virtualNodes := make([]*tke.VirtualNodeSpec, 0)
	for _, node := range serverlessNodes {
		nodeItem := node.(map[string]interface{})
		virtualNodes = append(virtualNodes, &tke.VirtualNodeSpec{
			DisplayName: common.StringPtr(nodeItem["display_name"].(string)),
			SubnetId:    common.StringPtr(nodeItem["subnet_id"].(string)),
		})
	}
	sgIds := make([]string, len(securityGroupIds))
	for i := 0; i < len(securityGroupIds); i++ {
		sgIds[i] = securityGroupIds[i].(string)
	}

	request := tke.NewCreateClusterVirtualNodePoolRequest()
	request.ClusterId = common.StringPtr(clusterId)
	request.Name = common.StringPtr(name)
	request.VirtualNodes = virtualNodes
	request.SecurityGroupIds = common.StringPtrs(sgIds)
	request.Labels = GetTkeLabels(d, "labels")
	request.Taints = GetTkeTaints(d, "taints")

	return request
}

func setDataFromDescribeVirtualNodePoolResponse(clusterId string, res *tke.VirtualNodePool, d *schema.ResourceData) error {
	d.SetId(clusterId + FILED_SP + *res.NodePoolId)
	_ = d.Set("name", res.Name)
	_ = d.Set("life_state", res.LifeState)
	labels := make(map[string]interface{})
	taints := make([]map[string]interface{}, 0)
	for i := 0; i < len(res.Labels); i++ {
		if res.Labels != nil && res.Labels[i].Name != nil && res.Labels[i].Value != nil {
			labels[*res.Labels[i].Name] = *res.Labels[i].Value
		}
	}
	for i := 0; i < len(res.Taints); i++ {
		if res.Taints != nil && res.Taints[i].Key != nil && res.Taints[i].Value != nil && res.Taints[i].Effect != nil {
			taint := map[string]interface{}{
				"key":    *res.Taints[i].Key,
				"value":  *res.Taints[i].Value,
				"effect": *res.Taints[i].Effect,
			}
			taints = append(taints, taint)
		}
	}
	_ = d.Set("labels", labels)
	_ = d.Set("taints", taints)

	return nil
}

func shouldServerlessNodePoolRetryReading(state string) bool {
	return state != "normal"
}

func GetOptimizedTkeLabels(d *schema.ResourceData, k string) []*tke.Label {
	labels := make([]*tke.Label, 0)
	if raw, ok := d.GetOk(k); ok {
		for k, v := range raw.(map[string]interface{}) {
			labels = append(labels, &tke.Label{Name: helper.String(k), Value: common.StringPtr(v.(string))})
		}
	}
	return labels
}

func GetOptimizedTkeTaints(d *schema.ResourceData, k string) []*tke.Taint {
	taints := make([]*tke.Taint, 0)
	if raw, ok := d.GetOk(k); ok {
		for _, v := range raw.([]interface{}) {
			vv := v.(map[string]interface{})
			taints = append(taints, &tke.Taint{Key: helper.String(vv["key"].(string)), Value: common.StringPtr(vv["value"].(string)), Effect: helper.String(vv["effect"].(string))})
		}
	}
	return taints
}
