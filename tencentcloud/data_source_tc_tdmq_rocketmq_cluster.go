/*
Use this data source to query detailed information of tdmqRocketmq cluster

Example Usage

```hcl
data "tencentcloud_tdmq_rocketmq_cluster" "example" {
  name_keyword = tencentcloud_tdmq_rocketmq_cluster.example.cluster_name
}

resource "tencentcloud_tdmq_rocketmq_cluster" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rocketmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTdmqRocketmqCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRocketmqClusterRead,
		Schema: map[string]*schema.Schema{
			"id_keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search by cluster ID.",
			},

			"name_keyword": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Search by cluster name.",
			},

			"cluster_id_list": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Filter by cluster ID.",
			},

			"cluster_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cluster information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Basic cluster information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"remark": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster description (up to 128 characters).",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster ID.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster name.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region information.",
									},
									"create_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Creation time in milliseconds.",
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
										Description: "Whether the namespace access point is supported.",
									},
									"vpcs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Vpc list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vpc_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Vpc ID.",
												},
												"subnet_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Subnet ID.",
												},
											},
										},
									},
									"is_vip": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether it is an exclusive instance.",
									},
									"rocketmq_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Rocketmq cluster identification.",
									},
								},
							},
						},
						"config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cluster configuration information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_tps_per_namespace": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum TPS per namespace.",
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
										Description: "Number of used topics.",
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
									"max_retention_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum message retention period in milliseconds.",
									},
									"max_latency_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum message delay in millisecond.",
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster status. `0`: Creating; `1`: Normal; `2`: Terminating; `3`: Deleted; `4`: Isolated; `5`: Creation failed; `6`: Deletion failed.",
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

func dataSourceTencentCloudRocketmqClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rocketmq_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("id_keyword"); ok {
		paramMap["id_keyword"] = v.(string)
	}

	if v, ok := d.GetOk("name_keyword"); ok {
		paramMap["name_keyword"] = v.(string)
	}

	if v, ok := d.GetOk("cluster_id_list"); ok {
		clusterIds := v.(*schema.Set).List()
		clusterIdList := make([]string, 0)
		for _, i := range clusterIds {
			clusterId := i.(string)
			clusterIdList = append(clusterIdList, clusterId)
		}
		paramMap["cluster_id_list"] = clusterIdList
	}

	service := TdmqRocketmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var clusterList []*rocketmq.RocketMQClusterDetail
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		results, e := service.DescribeRocketmqClusterByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		clusterList = results
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read Rocketmq clusterList failed, reason:%+v", logId, err)
		return err
	}
	clusterListList := []interface{}{}
	ids := make([]string, 0)

	for _, clusterList := range clusterList {
		clusterListMap := map[string]interface{}{}
		if clusterList.Info != nil {
			infoMap := map[string]interface{}{}
			infoMap["cluster_id"] = clusterList.Info.ClusterId
			ids = append(ids, *clusterList.Info.ClusterId)
			if clusterList.Info.Remark != nil {
				infoMap["remark"] = clusterList.Info.Remark
			}
			if clusterList.Info.ClusterName != nil {
				infoMap["cluster_name"] = clusterList.Info.ClusterName
			}
			if clusterList.Info.Region != nil {
				infoMap["region"] = clusterList.Info.Region
			}
			if clusterList.Info.CreateTime != nil {
				infoMap["create_time"] = clusterList.Info.CreateTime
			}
			if clusterList.Info.PublicEndPoint != nil {
				infoMap["public_end_point"] = clusterList.Info.PublicEndPoint
			}
			if clusterList.Info.VpcEndPoint != nil {
				infoMap["vpc_end_point"] = clusterList.Info.VpcEndPoint
			}
			if clusterList.Info.SupportNamespaceEndpoint != nil {
				infoMap["support_namespace_endpoint"] = clusterList.Info.SupportNamespaceEndpoint
			}
			if clusterList.Info.Vpcs != nil {
				vpcsList := []interface{}{}
				for _, vpcs := range clusterList.Info.Vpcs {
					vpcsMap := map[string]interface{}{}
					if vpcs.VpcId != nil {
						vpcsMap["vpc_id"] = vpcs.VpcId
					}
					if vpcs.SubnetId != nil {
						vpcsMap["subnet_id"] = vpcs.SubnetId
					}

					vpcsList = append(vpcsList, vpcsMap)
				}
				infoMap["vpcs"] = vpcsList
			}
			if clusterList.Info.IsVip != nil {
				infoMap["is_vip"] = clusterList.Info.IsVip
			}
			if clusterList.Info.RocketMQFlag != nil {
				infoMap["rocketmq_flag"] = clusterList.Info.RocketMQFlag
			}

			clusterListMap["info"] = []interface{}{infoMap}
		}
		if clusterList.Config != nil {
			configMap := map[string]interface{}{}
			if clusterList.Config.MaxTpsPerNamespace != nil {
				configMap["max_tps_per_namespace"] = clusterList.Config.MaxTpsPerNamespace
			}
			if clusterList.Config.MaxNamespaceNum != nil {
				configMap["max_namespace_num"] = clusterList.Config.MaxNamespaceNum
			}
			if clusterList.Config.UsedNamespaceNum != nil {
				configMap["used_namespace_num"] = clusterList.Config.UsedNamespaceNum
			}
			if clusterList.Config.MaxTopicNum != nil {
				configMap["max_topic_num"] = clusterList.Config.MaxTopicNum
			}
			if clusterList.Config.UsedTopicNum != nil {
				configMap["used_topic_num"] = clusterList.Config.UsedTopicNum
			}
			if clusterList.Config.MaxGroupNum != nil {
				configMap["max_group_num"] = clusterList.Config.MaxGroupNum
			}
			if clusterList.Config.UsedGroupNum != nil {
				configMap["used_group_num"] = clusterList.Config.UsedGroupNum
			}
			if clusterList.Config.MaxRetentionTime != nil {
				configMap["max_retention_time"] = clusterList.Config.MaxRetentionTime
			}
			if clusterList.Config.MaxLatencyTime != nil {
				configMap["max_latency_time"] = clusterList.Config.MaxLatencyTime
			}

			clusterListMap["config"] = []interface{}{configMap}
		}
		if clusterList.Status != nil {
			clusterListMap["status"] = clusterList.Status
		}

		clusterListList = append(clusterListList, clusterListMap)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("cluster_list", clusterListList)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), clusterListList); e != nil {
			return e
		}
	}

	return nil
}
