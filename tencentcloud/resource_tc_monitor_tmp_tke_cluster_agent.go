/*
Provides a resource to create a tmp tke cluster agent

Example Usage

```hcl

resource "tencentcloud_monitor_tmp_tke_cluster_agent" "tmpClusterAgent" {
  instance_id = "prom-xxx"

  agents {
    region          = "ap-xxx"
    cluster_type    = "eks"
    cluster_id      = "cls-xxx"
    enable_external = false
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpTkeClusterAgent() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpTkeClusterAgentRead,
		Create: resourceTencentCloudMonitorTmpTkeClusterAgentCreate,
		Update: resourceTencentCloudMonitorTmpTkeClusterAgentUpdate,
		Delete: resourceTencentCloudMonitorTmpTkeClusterAgentDelete,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance Id.",
			},

			"agents": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "agent list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Limitation of region.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of cluster.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "An id identify the cluster, like `cls-xxxxxx`.",
						},
						"enable_external": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether to enable the public network CLB.",
						},
						"in_cluster_pod_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Pod configuration for components deployed in the cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_net": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether to use HostNetWork.",
									},
									"node_selector": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specify the pod to run the node.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The pod configuration name of the component deployed in the cluster.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Pod configuration values for components deployed in the cluster.",
												},
											},
										},
									},
									"tolerations": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Tolerate Stain.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The taint key to which the tolerance applies.",
												},
												"operator": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "key-value relationship.",
												},
												"effect": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "blemish effect to match.",
												},
											},
										},
									},
								},
							},
						},
						"external_labels": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "All metrics collected by the cluster will carry these labels.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Indicator name.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Index value.",
									},
								},
							},
						},
						"not_install_basic_scrape": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to install the default collection configuration.",
						},
						"not_scrape": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to collect indicators, true means drop all indicators, false means collect default indicators.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the name of the cluster.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "agent state, `normal`, `abnormal`.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMonitorTmpTkeClusterAgentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_cluster_agent.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tke.NewCreatePrometheusClusterAgentRequest()

	instanceId := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	clusterId := ""
	clusterType := ""
	if dMap, ok := helper.InterfacesHeadMap(d, "agents"); ok {
		prometheusClusterAgent := tke.PrometheusClusterAgentBasic{}
		if v, ok := dMap["region"]; ok {
			prometheusClusterAgent.Region = helper.String(v.(string))
		}
		if v, ok := dMap["cluster_type"]; ok {
			clusterType = v.(string)
			prometheusClusterAgent.ClusterType = helper.String(v.(string))
		}
		if v, ok := dMap["cluster_id"]; ok {
			clusterId = v.(string)
			prometheusClusterAgent.ClusterId = helper.String(v.(string))
		}
		if v, ok := dMap["enable_external"]; ok {
			prometheusClusterAgent.EnableExternal = helper.Bool(v.(bool))
		}
		if v, ok := dMap["in_cluster_pod_config"]; ok {
			var clusterAgentPodConfig *tke.PrometheusClusterAgentPodConfig
			if len(v.([]interface{})) > 0 {
				podConfig := v.([]interface{})[0].(map[string]interface{})

				if vv, ok := podConfig["host_net"]; ok {
					clusterAgentPodConfig.HostNet = helper.Bool(vv.(bool))
				}
				if vv, ok := podConfig["node_selector"]; ok {
					labelsList := vv.([]interface{})
					nodeSelectorKV := make([]*tke.Label, 0, len(labelsList))
					for _, labels := range labelsList {
						label := labels.(map[string]interface{})
						var kv tke.Label
						kv.Name = helper.String(label["name"].(string))
						kv.Value = helper.String(label["value"].(string))
						nodeSelectorKV = append(nodeSelectorKV, &kv)
					}
					clusterAgentPodConfig.NodeSelector = nodeSelectorKV
				}
				if vv, ok := podConfig["tolerations"]; ok {
					tolerationList := vv.([]interface{})
					tolerations := make([]*tke.Toleration, 0, len(tolerationList))
					for _, t := range tolerationList {
						tolerationMap := t.(map[string]interface{})
						var toleration tke.Toleration
						toleration.Key = helper.String(tolerationMap["key"].(string))
						toleration.Operator = helper.String(tolerationMap["operator"].(string))
						toleration.Effect = helper.String(tolerationMap["effect"].(string))
						tolerations = append(tolerations, &toleration)
					}
					clusterAgentPodConfig.Tolerations = tolerations
				}
			}
			prometheusClusterAgent.InClusterPodConfig = clusterAgentPodConfig
		}
		if v, ok := dMap["external_labels"]; ok {
			labelsList := v.([]interface{})
			externalKV := make([]*tke.Label, 0, len(labelsList))
			for _, labels := range labelsList {
				label := labels.(map[string]interface{})
				var kv tke.Label
				kv.Name = helper.String(label["name"].(string))
				kv.Value = helper.String(label["value"].(string))
				externalKV = append(externalKV, &kv)
			}
			prometheusClusterAgent.ExternalLabels = externalKV
		}
		if v, ok := dMap["not_install_basic_scrape"]; ok {
			prometheusClusterAgent.NotInstallBasicScrape = helper.Bool(v.(bool))
		}
		if v, ok := dMap["not_scrape"]; ok {
			prometheusClusterAgent.NotScrape = helper.Bool(v.(bool))
		}
		var prometheusClusterAgents []*tke.PrometheusClusterAgentBasic
		prometheusClusterAgents = append(prometheusClusterAgents, &prometheusClusterAgent)
		request.Agents = prometheusClusterAgents
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().CreatePrometheusClusterAgent(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tke cluster agent failed, reason:%+v", logId, err)
		return err
	}

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		clusterAgent, errRet := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *clusterAgent.Status == "normal" {
			return nil
		}
		if *clusterAgent.Status == "abnormal" {
			return resource.NonRetryableError(fmt.Errorf("cluster agent status is %v, operate failed.", *clusterAgent.Status))
		}
		return resource.RetryableError(fmt.Errorf("cluster agent status is %v, retry...", *clusterAgent.Status))
	})
	if err != nil {
		return err
	}

	d.SetId(strings.Join([]string{instanceId, clusterId, clusterType}, FILED_SP))

	return resourceTencentCloudMonitorTmpTkeClusterAgentRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeClusterAgentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_cluster_agent.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	clusterId := ids[1]
	clusterType := ids[2]

	clusterAgent, err := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)

	if err != nil {
		return err
	}

	if clusterAgent == nil {
		d.SetId("")
		return fmt.Errorf("resource `global_notification` %s does not exist", instanceId)
	}

	_ = d.Set("cluster_id", clusterAgent.ClusterId)
	_ = d.Set("cluster_type", clusterAgent.ClusterType)
	_ = d.Set("status", clusterAgent.Status)
	_ = d.Set("cluster_name", clusterAgent.ClusterName)
	if clusterAgent.ExternalLabels != nil {
		clusterAgentList := clusterAgent.ExternalLabels
		result := make([]map[string]interface{}, 0, len(clusterAgentList))
		for _, v := range clusterAgentList {
			mapping := map[string]interface{}{
				"name":  v.Name,
				"value": v.Value,
			}
			result = append(result, mapping)
		}
		_ = d.Set("external_labels", result)
	}

	return nil
}

func resourceTencentCloudMonitorTmpTkeClusterAgentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_cluster_agent.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tke.NewModifyPrometheusAgentExternalLabelsRequest()
		response *tke.ModifyPrometheusAgentExternalLabelsResponse
	)

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	clusterId := ids[1]
	request.InstanceId = &instanceId
	request.ClusterId = &clusterId

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("cluster_id") {
		return fmt.Errorf("`cluster_id` do not support change now.")
	}

	if d.HasChange("agents") {
		if dMap, ok := helper.InterfacesHeadMap(d, "agents"); ok {
			if v, ok := dMap["external_labels"]; ok {
				labelsList := v.([]interface{})
				externalKV := make([]*tke.Label, 0, len(labelsList))
				for _, labels := range labelsList {
					label := labels.(map[string]interface{})
					var kv tke.Label
					kv.Name = helper.String(label["name"].(string))
					kv.Value = helper.String(label["value"].(string))
					externalKV = append(externalKV, &kv)
				}
				request.ExternalLabels = externalKV
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().ModifyPrometheusAgentExternalLabels(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if err != nil {
		return err
	}

	return resourceTencentCloudTkeTmpAlertPolicyRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeClusterAgentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_cluster_agent.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	clusterId := ids[1]
	clusterType := ids[2]

	if err := service.DeletePrometheusClusterAgent(ctx, instanceId, clusterId, clusterType); err != nil {
		return err
	}

	err := resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		clusterAgent, errRet := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if clusterAgent == nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cluster agent status is %v, retry...", *clusterAgent.Status))
	})
	if err != nil {
		return err
	}

	return nil
}
