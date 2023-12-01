/*
Provides a resource to create a tmp tke cluster agent

Example Usage

```hcl
variable "default_instance_type" {
  default = "SA1.MEDIUM2"
}

variable "availability_zone_first" {
  default = "ap-guangzhou-3"
}

variable "availability_zone_second" {
  default = "ap-guangzhou-4"
}

variable "example_cluster_cidr" {
  default = "10.31.0.0/16"
}

locals {
  first_vpc_id     = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.vpc_id
  first_subnet_id  = data.tencentcloud_vpc_subnets.vpc_one.instance_list.0.subnet_id
  second_vpc_id    = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.vpc_id
  second_subnet_id = data.tencentcloud_vpc_subnets.vpc_two.instance_list.0.subnet_id
  sg_id            = tencentcloud_security_group.sg.id
  image_id         = data.tencentcloud_images.default.image_id
}

data "tencentcloud_vpc_subnets" "vpc_one" {
  is_default        = true
  availability_zone = var.availability_zone_first
}

data "tencentcloud_vpc_subnets" "vpc_two" {
  is_default        = true
  availability_zone = var.availability_zone_second
}

resource "tencentcloud_security_group" "sg" {
  name = "tf-example-sg"
}

resource "tencentcloud_security_group_lite_rule" "sg_rule" {
  security_group_id = tencentcloud_security_group.sg.id

  ingress = [
    "ACCEPT#10.0.0.0/16#ALL#ALL",
    "ACCEPT#172.16.0.0/22#ALL#ALL",
    "DROP#0.0.0.0/0#ALL#ALL",
  ]

  egress = [
    "ACCEPT#172.16.0.0/22#ALL#ALL",
  ]
}

data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

resource "tencentcloud_kubernetes_cluster" "example" {
  vpc_id                          = local.first_vpc_id
  cluster_cidr                    = var.example_cluster_cidr
  cluster_max_pod_num             = 32
  cluster_name                    = "tf_example_cluster"
  cluster_desc                    = "example for tke cluster"
  cluster_max_service_num         = 32
  cluster_internet                = false
  cluster_internet_security_group = local.sg_id
  cluster_version                 = "1.22.5"
  cluster_deploy_type             = "MANAGED_CLUSTER"

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_first
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.first_subnet_id
    img_id                     = local.image_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    password = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  worker_config {
    count                      = 1
    availability_zone          = var.availability_zone_second
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = local.second_subnet_id

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    # key_ids                   = ["skey-11112222"]
    cam_role_name = "CVM_QcsRole"
    password      = "ZZXXccvv1212" // Optional, should be set if key_ids not set.
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }
}


# create monitor
variable "zone" {
  default = "ap-guangzhou"
}

variable "cluster_type" {
  default = "tke"
}

resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = local.first_vpc_id
  subnet_id           = local.first_subnet_id
  data_retention_time = 30
  zone                = var.availability_zone_second
  tags = {
    "createdBy" = "terraform"
  }
}

# tmp tke bind
resource "tencentcloud_monitor_tmp_tke_cluster_agent" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id

  agents {
    region          = var.zone
    cluster_type    = var.cluster_type
    cluster_id      = tencentcloud_kubernetes_cluster.example.id
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
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

	request := monitor.NewCreatePrometheusClusterAgentRequest()

	instanceId := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	clusterId := ""
	clusterType := ""
	if dMap, ok := helper.InterfacesHeadMap(d, "agents"); ok {
		prometheusClusterAgent := monitor.PrometheusClusterAgentBasic{}
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
			var clusterAgentPodConfig *monitor.PrometheusClusterAgentPodConfig
			if len(v.([]interface{})) > 0 {
				podConfig := v.([]interface{})[0].(map[string]interface{})

				if vv, ok := podConfig["host_net"]; ok {
					clusterAgentPodConfig.HostNet = helper.Bool(vv.(bool))
				}
				if vv, ok := podConfig["node_selector"]; ok {
					labelsList := vv.([]interface{})
					nodeSelectorKV := make([]*monitor.Label, 0, len(labelsList))
					for _, labels := range labelsList {
						label := labels.(map[string]interface{})
						var kv monitor.Label
						kv.Name = helper.String(label["name"].(string))
						kv.Value = helper.String(label["value"].(string))
						nodeSelectorKV = append(nodeSelectorKV, &kv)
					}
					clusterAgentPodConfig.NodeSelector = nodeSelectorKV
				}
				if vv, ok := podConfig["tolerations"]; ok {
					tolerationList := vv.([]interface{})
					tolerations := make([]*monitor.Toleration, 0, len(tolerationList))
					for _, t := range tolerationList {
						tolerationMap := t.(map[string]interface{})
						var toleration monitor.Toleration
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
			externalKV := make([]*monitor.Label, 0, len(labelsList))
			for _, labels := range labelsList {
				label := labels.(map[string]interface{})
				var kv monitor.Label
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
		var prometheusClusterAgents []*monitor.PrometheusClusterAgentBasic
		prometheusClusterAgents = append(prometheusClusterAgents, &prometheusClusterAgent)
		request.Agents = prometheusClusterAgents
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusClusterAgent(request)
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

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		clusterAgent, errRet := service.DescribeTmpTkeClusterAgentsById(ctx, instanceId, clusterId, clusterType)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if *clusterAgent.Status == "normal" {
			return nil
		}
		// if *clusterAgent.Status == "abnormal" {
		// 	return resource.NonRetryableError(fmt.Errorf("cluster agent status is %v, operate failed.", *clusterAgent.Status))
		// }
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

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

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

	var agents []map[string]interface{}
	agent := make(map[string]interface{})
	agent["cluster_id"] = clusterAgent.ClusterId
	agent["cluster_type"] = clusterAgent.ClusterType
	agent["status"] = clusterAgent.Status
	agent["cluster_name"] = clusterAgent.ClusterName
	agent["region"] = clusterAgent.Region
	//if clusterAgent.ExternalLabels != nil {
	//	clusterAgentList := clusterAgent.ExternalLabels
	//	result := make([]map[string]interface{}, 0, len(clusterAgentList))
	//	for _, v := range clusterAgentList {
	//		mapping := map[string]interface{}{
	//			"name":  v.Name,
	//			"value": v.Value,
	//		}
	//		result = append(result, mapping)
	//	}
	//	agent["external_labels"] = result
	//}
	agents = append(agents, agent)
	_ = d.Set("agents", agents)

	return nil
}

func resourceTencentCloudMonitorTmpTkeClusterAgentUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_cluster_agent.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewModifyPrometheusAgentExternalLabelsRequest()
		response *monitor.ModifyPrometheusAgentExternalLabelsResponse
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
				externalKV := make([]*monitor.Label, 0, len(labelsList))
				for _, labels := range labelsList {
					label := labels.(map[string]interface{})
					var kv monitor.Label
					kv.Name = helper.String(label["name"].(string))
					kv.Value = helper.String(label["value"].(string))
					externalKV = append(externalKV, &kv)
				}
				request.ExternalLabels = externalKV
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().ModifyPrometheusAgentExternalLabels(request)
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
	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

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
