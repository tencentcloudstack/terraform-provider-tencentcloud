/*
Use this data source to query detailed information of tdmq vip_instance

Example Usage

```hcl
data "tencentcloud_tdmq_vip_instance" "vip_instance" {
  cluster_id = &lt;nil&gt;
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdmqVipInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqVipInstanceRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"cluster_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster Name.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation time, in milliseconds.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster description informationNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"public_end_point": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public network access address.",
						},
						"vpc_end_point": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC access address.",
						},
						"support_namespace_endpoint": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether namespace access points are supportedNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"vpcs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "VPC informationNote: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "VPC ID.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subnet Id.",
									},
								},
							},
						},
						"is_vip": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it is a dedicated instanceNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"rocket_m_q_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Rocketmq cluster identificationNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing status, 1 means normal, 2 means stopped, 3 means destroyedNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"isolate_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Overdue suspension time, in millisecondsNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"http_public_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP protocol public network access addressNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"http_vpc_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP protocol VPC access addressNote: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"instance_config": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_tps_per_namespace": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Single namespace TPS upper limit.",
						},
						"max_namespace_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of namespaces.",
						},
						"used_namespace_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of used namespaces.",
						},
						"max_topic_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of topics.",
						},
						"used_topic_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of topics used.",
						},
						"max_group_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of groups.",
						},
						"used_group_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of used groups.",
						},
						"config_display": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster type.",
						},
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of cluster nodes.",
						},
						"node_distribution": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node distribution.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone.",
									},
									"zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone id.",
									},
									"node_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of nodes.",
									},
								},
							},
						},
						"topic_distribution": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Topic distribution.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Topic type.",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of topics.",
									},
								},
							},
						},
						"max_queues_per_topic": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of queues per topicNote: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTdmqVipInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_vip_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var clusterInfo []*tdmq.RocketMQClusterInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqVipInstanceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		clusterInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(clusterInfo))
	if clusterInfo != nil {
		rocketMQClusterInfoMap := map[string]interface{}{}

		if clusterInfo.ClusterId != nil {
			rocketMQClusterInfoMap["cluster_id"] = clusterInfo.ClusterId
		}

		if clusterInfo.ClusterName != nil {
			rocketMQClusterInfoMap["cluster_name"] = clusterInfo.ClusterName
		}

		if clusterInfo.Region != nil {
			rocketMQClusterInfoMap["region"] = clusterInfo.Region
		}

		if clusterInfo.CreateTime != nil {
			rocketMQClusterInfoMap["create_time"] = clusterInfo.CreateTime
		}

		if clusterInfo.Remark != nil {
			rocketMQClusterInfoMap["remark"] = clusterInfo.Remark
		}

		if clusterInfo.PublicEndPoint != nil {
			rocketMQClusterInfoMap["public_end_point"] = clusterInfo.PublicEndPoint
		}

		if clusterInfo.VpcEndPoint != nil {
			rocketMQClusterInfoMap["vpc_end_point"] = clusterInfo.VpcEndPoint
		}

		if clusterInfo.SupportNamespaceEndpoint != nil {
			rocketMQClusterInfoMap["support_namespace_endpoint"] = clusterInfo.SupportNamespaceEndpoint
		}

		if clusterInfo.Vpcs != nil {
			vpcsList := []interface{}{}
			for _, vpcs := range clusterInfo.Vpcs {
				vpcsMap := map[string]interface{}{}

				if vpcs.VpcId != nil {
					vpcsMap["vpc_id"] = vpcs.VpcId
				}

				if vpcs.SubnetId != nil {
					vpcsMap["subnet_id"] = vpcs.SubnetId
				}

				vpcsList = append(vpcsList, vpcsMap)
			}

			rocketMQClusterInfoMap["vpcs"] = []interface{}{vpcsList}
		}

		if clusterInfo.IsVip != nil {
			rocketMQClusterInfoMap["is_vip"] = clusterInfo.IsVip
		}

		if clusterInfo.RocketMQFlag != nil {
			rocketMQClusterInfoMap["rocket_m_q_flag"] = clusterInfo.RocketMQFlag
		}

		if clusterInfo.Status != nil {
			rocketMQClusterInfoMap["status"] = clusterInfo.Status
		}

		if clusterInfo.IsolateTime != nil {
			rocketMQClusterInfoMap["isolate_time"] = clusterInfo.IsolateTime
		}

		if clusterInfo.HttpPublicEndpoint != nil {
			rocketMQClusterInfoMap["http_public_endpoint"] = clusterInfo.HttpPublicEndpoint
		}

		if clusterInfo.HttpVpcEndpoint != nil {
			rocketMQClusterInfoMap["http_vpc_endpoint"] = clusterInfo.HttpVpcEndpoint
		}

		ids = append(ids, *clusterInfo.ClusterId)
		_ = d.Set("cluster_info", rocketMQClusterInfoMap)
	}

	if instanceConfig != nil {
		rocketMQInstanceConfigMap := map[string]interface{}{}

		if instanceConfig.MaxTpsPerNamespace != nil {
			rocketMQInstanceConfigMap["max_tps_per_namespace"] = instanceConfig.MaxTpsPerNamespace
		}

		if instanceConfig.MaxNamespaceNum != nil {
			rocketMQInstanceConfigMap["max_namespace_num"] = instanceConfig.MaxNamespaceNum
		}

		if instanceConfig.UsedNamespaceNum != nil {
			rocketMQInstanceConfigMap["used_namespace_num"] = instanceConfig.UsedNamespaceNum
		}

		if instanceConfig.MaxTopicNum != nil {
			rocketMQInstanceConfigMap["max_topic_num"] = instanceConfig.MaxTopicNum
		}

		if instanceConfig.UsedTopicNum != nil {
			rocketMQInstanceConfigMap["used_topic_num"] = instanceConfig.UsedTopicNum
		}

		if instanceConfig.MaxGroupNum != nil {
			rocketMQInstanceConfigMap["max_group_num"] = instanceConfig.MaxGroupNum
		}

		if instanceConfig.UsedGroupNum != nil {
			rocketMQInstanceConfigMap["used_group_num"] = instanceConfig.UsedGroupNum
		}

		if instanceConfig.ConfigDisplay != nil {
			rocketMQInstanceConfigMap["config_display"] = instanceConfig.ConfigDisplay
		}

		if instanceConfig.NodeCount != nil {
			rocketMQInstanceConfigMap["node_count"] = instanceConfig.NodeCount
		}

		if instanceConfig.NodeDistribution != nil {
			nodeDistributionList := []interface{}{}
			for _, nodeDistribution := range instanceConfig.NodeDistribution {
				nodeDistributionMap := map[string]interface{}{}

				if nodeDistribution.ZoneName != nil {
					nodeDistributionMap["zone_name"] = nodeDistribution.ZoneName
				}

				if nodeDistribution.ZoneId != nil {
					nodeDistributionMap["zone_id"] = nodeDistribution.ZoneId
				}

				if nodeDistribution.NodeCount != nil {
					nodeDistributionMap["node_count"] = nodeDistribution.NodeCount
				}

				nodeDistributionList = append(nodeDistributionList, nodeDistributionMap)
			}

			rocketMQInstanceConfigMap["node_distribution"] = []interface{}{nodeDistributionList}
		}

		if instanceConfig.TopicDistribution != nil {
			topicDistributionList := []interface{}{}
			for _, topicDistribution := range instanceConfig.TopicDistribution {
				topicDistributionMap := map[string]interface{}{}

				if topicDistribution.TopicType != nil {
					topicDistributionMap["topic_type"] = topicDistribution.TopicType
				}

				if topicDistribution.Count != nil {
					topicDistributionMap["count"] = topicDistribution.Count
				}

				topicDistributionList = append(topicDistributionList, topicDistributionMap)
			}

			rocketMQInstanceConfigMap["topic_distribution"] = []interface{}{topicDistributionList}
		}

		if instanceConfig.MaxQueuesPerTopic != nil {
			rocketMQInstanceConfigMap["max_queues_per_topic"] = instanceConfig.MaxQueuesPerTopic
		}

		ids = append(ids, *instanceConfig.ClusterId)
		_ = d.Set("instance_config", rocketMQInstanceConfigMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), rocketMQClusterInfoMap); e != nil {
			return e
		}
	}
	return nil
}
