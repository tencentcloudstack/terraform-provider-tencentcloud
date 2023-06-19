/*
Use this data source to query detailed information of tdmq pro_instance_detail

Example Usage

```hcl
data "tencentcloud_tdmq_pro_instance_detail" "pro_instance_detail" {
  cluster_id = "pulsar-9n95ax58b9vn"
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

func dataSourceTencentCloudTdmqProInstanceDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqProInstanceDetailRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster Id.",
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
							Description: "Cluster Id.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster name.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Descriptive information.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster status, 0: creating, 1: normal, 2: isolated.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cluster version.",
						},
						"node_distribution": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node distributionNote: This field may return null, indicating that no valid value can be obtained.",
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
										Description: "Availability zone ID.",
									},
									"node_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of nodes.",
									},
								},
							},
						},
						"max_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum storage capacity, unit: MB.",
						},
						"can_edit_route": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Can the route be modifiedNote: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},
			"network_access_point_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster network access point informationNote: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the vpc, the supporting network and the access point of the public network, this field is emptyNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet id, support network and public network access point, this field is emptyNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "access address.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"route_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Access point type: 0: support network access point 1: VPC access point 2: public network access point.",
						},
					},
				},
			},
			"cluster_spec_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster specification informationNote: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster specification name.",
						},
						"max_tps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "peak tps.",
						},
						"max_band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "peak bandwidth. Unit: mbps.",
						},
						"max_namespaces": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of namespaces.",
						},
						"max_topics": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of topic partitions.",
						},
						"scalable_tps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Elastic TPS outside specificationNote: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudTdmqProInstanceDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tdmq_pro_instance_detail.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId       = getLogId(contextNil)
		ctx         = context.WithValue(context.TODO(), logIdKey, logId)
		service     = TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}
		clusterInfo *tdmq.DescribePulsarProInstanceDetailResponseParams
		clusterId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqProInstanceDetailByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		clusterInfo = result
		return nil
	})

	if err != nil {
		return err
	}

	if clusterInfo != nil {
		pulsarProClusterInfoMap := map[string]interface{}{}
		tmpList := []interface{}{}

		if clusterInfo.ClusterInfo.ClusterId != nil {
			pulsarProClusterInfoMap["cluster_id"] = clusterInfo.ClusterInfo.ClusterId
		}

		if clusterInfo.ClusterInfo.ClusterName != nil {
			pulsarProClusterInfoMap["cluster_name"] = clusterInfo.ClusterInfo.ClusterName
		}

		if clusterInfo.ClusterInfo.Remark != nil {
			pulsarProClusterInfoMap["remark"] = clusterInfo.ClusterInfo.Remark
		}

		if clusterInfo.ClusterInfo.CreateTime != nil {
			pulsarProClusterInfoMap["create_time"] = clusterInfo.ClusterInfo.CreateTime
		}

		if clusterInfo.ClusterInfo.Status != nil {
			pulsarProClusterInfoMap["status"] = clusterInfo.ClusterInfo.Status
		}

		if clusterInfo.ClusterInfo.Version != nil {
			pulsarProClusterInfoMap["version"] = clusterInfo.ClusterInfo.Version
		}

		if clusterInfo.ClusterInfo.NodeDistribution != nil {
			nodeDistributionList := []interface{}{}
			for _, nodeDistribution := range clusterInfo.ClusterInfo.NodeDistribution {
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

			pulsarProClusterInfoMap["node_distribution"] = nodeDistributionList
		}

		if clusterInfo.ClusterInfo.MaxStorage != nil {
			pulsarProClusterInfoMap["max_storage"] = clusterInfo.ClusterInfo.MaxStorage
		}

		if clusterInfo.ClusterInfo.CanEditRoute != nil {
			pulsarProClusterInfoMap["can_edit_route"] = clusterInfo.ClusterInfo.CanEditRoute
		}

		tmpList = append(tmpList, pulsarProClusterInfoMap)
		_ = d.Set("cluster_info", tmpList)
	}

	if clusterInfo.NetworkAccessPointInfos != nil {
		tmpList := []interface{}{}
		for _, pulsarNetworkAccessPointInfo := range clusterInfo.NetworkAccessPointInfos {
			pulsarNetworkAccessPointInfoMap := map[string]interface{}{}

			if pulsarNetworkAccessPointInfo.VpcId != nil {
				pulsarNetworkAccessPointInfoMap["vpc_id"] = pulsarNetworkAccessPointInfo.VpcId
			}

			if pulsarNetworkAccessPointInfo.SubnetId != nil {
				pulsarNetworkAccessPointInfoMap["subnet_id"] = pulsarNetworkAccessPointInfo.SubnetId
			}

			if pulsarNetworkAccessPointInfo.Endpoint != nil {
				pulsarNetworkAccessPointInfoMap["endpoint"] = pulsarNetworkAccessPointInfo.Endpoint
			}

			if pulsarNetworkAccessPointInfo.InstanceId != nil {
				pulsarNetworkAccessPointInfoMap["instance_id"] = pulsarNetworkAccessPointInfo.InstanceId
			}

			if pulsarNetworkAccessPointInfo.RouteType != nil {
				pulsarNetworkAccessPointInfoMap["route_type"] = pulsarNetworkAccessPointInfo.RouteType
			}

			tmpList = append(tmpList, pulsarNetworkAccessPointInfoMap)
		}

		_ = d.Set("network_access_point_infos", tmpList)
	}

	if clusterInfo.ClusterSpecInfo != nil {
		pulsarProClusterSpecInfoMap := map[string]interface{}{}
		tmpList := []interface{}{}

		if clusterInfo.ClusterSpecInfo.SpecName != nil {
			pulsarProClusterSpecInfoMap["spec_name"] = clusterInfo.ClusterSpecInfo.SpecName
		}

		if clusterInfo.ClusterSpecInfo.MaxTps != nil {
			pulsarProClusterSpecInfoMap["max_tps"] = clusterInfo.ClusterSpecInfo.MaxTps
		}

		if clusterInfo.ClusterSpecInfo.MaxBandWidth != nil {
			pulsarProClusterSpecInfoMap["max_band_width"] = clusterInfo.ClusterSpecInfo.MaxBandWidth
		}

		if clusterInfo.ClusterSpecInfo.MaxNamespaces != nil {
			pulsarProClusterSpecInfoMap["max_namespaces"] = clusterInfo.ClusterSpecInfo.MaxNamespaces
		}

		if clusterInfo.ClusterSpecInfo.MaxTopics != nil {
			pulsarProClusterSpecInfoMap["max_topics"] = clusterInfo.ClusterSpecInfo.MaxTopics
		}

		if clusterInfo.ClusterSpecInfo.ScalableTps != nil {
			pulsarProClusterSpecInfoMap["scalable_tps"] = clusterInfo.ClusterSpecInfo.ScalableTps
		}

		tmpList = append(tmpList, pulsarProClusterSpecInfoMap)
		_ = d.Set("cluster_spec_info", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
