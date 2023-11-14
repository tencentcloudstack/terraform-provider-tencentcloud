/*
Use this data source to query detailed information of tdmq pro_instance_detail

Example Usage

```hcl
data "tencentcloud_tdmq_pro_instance_detail" "pro_instance_detail" {
  cluster_id = ""
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
				Description: "Cluster infomration.",
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
							Description: "Cluster version.",
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
							Description: "Access address.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance id.",
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
							Description: "Peak tps.",
						},
						"max_band_width": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Peak bandwidth. Unit: mbps.",
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

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	service := TdmqService{client: meta.(*TencentCloudClient).apiV3Conn}

	var clusterInfo []*tdmq.PulsarProClusterInfo

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

	ids := make([]string, 0, len(clusterInfo))
	if clusterInfo != nil {
		pulsarProClusterInfoMap := map[string]interface{}{}

		if clusterInfo.ClusterId != nil {
			pulsarProClusterInfoMap["cluster_id"] = clusterInfo.ClusterId
		}

		if clusterInfo.ClusterName != nil {
			pulsarProClusterInfoMap["cluster_name"] = clusterInfo.ClusterName
		}

		if clusterInfo.Remark != nil {
			pulsarProClusterInfoMap["remark"] = clusterInfo.Remark
		}

		if clusterInfo.CreateTime != nil {
			pulsarProClusterInfoMap["create_time"] = clusterInfo.CreateTime
		}

		if clusterInfo.Status != nil {
			pulsarProClusterInfoMap["status"] = clusterInfo.Status
		}

		if clusterInfo.Version != nil {
			pulsarProClusterInfoMap["version"] = clusterInfo.Version
		}

		if clusterInfo.NodeDistribution != nil {
			nodeDistributionList := []interface{}{}
			for _, nodeDistribution := range clusterInfo.NodeDistribution {
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

			pulsarProClusterInfoMap["node_distribution"] = []interface{}{nodeDistributionList}
		}

		if clusterInfo.MaxStorage != nil {
			pulsarProClusterInfoMap["max_storage"] = clusterInfo.MaxStorage
		}

		if clusterInfo.CanEditRoute != nil {
			pulsarProClusterInfoMap["can_edit_route"] = clusterInfo.CanEditRoute
		}

		ids = append(ids, *clusterInfo.ClusterId)
		_ = d.Set("cluster_info", pulsarProClusterInfoMap)
	}

	if networkAccessPointInfos != nil {
		for _, pulsarNetworkAccessPointInfo := range networkAccessPointInfos {
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

			ids = append(ids, *pulsarNetworkAccessPointInfo.ClusterId)
			tmpList = append(tmpList, pulsarNetworkAccessPointInfoMap)
		}

		_ = d.Set("network_access_point_infos", tmpList)
	}

	if clusterSpecInfo != nil {
		pulsarProClusterSpecInfoMap := map[string]interface{}{}

		if clusterSpecInfo.SpecName != nil {
			pulsarProClusterSpecInfoMap["spec_name"] = clusterSpecInfo.SpecName
		}

		if clusterSpecInfo.MaxTps != nil {
			pulsarProClusterSpecInfoMap["max_tps"] = clusterSpecInfo.MaxTps
		}

		if clusterSpecInfo.MaxBandWidth != nil {
			pulsarProClusterSpecInfoMap["max_band_width"] = clusterSpecInfo.MaxBandWidth
		}

		if clusterSpecInfo.MaxNamespaces != nil {
			pulsarProClusterSpecInfoMap["max_namespaces"] = clusterSpecInfo.MaxNamespaces
		}

		if clusterSpecInfo.MaxTopics != nil {
			pulsarProClusterSpecInfoMap["max_topics"] = clusterSpecInfo.MaxTopics
		}

		if clusterSpecInfo.ScalableTps != nil {
			pulsarProClusterSpecInfoMap["scalable_tps"] = clusterSpecInfo.ScalableTps
		}

		ids = append(ids, *clusterSpecInfo.ClusterId)
		_ = d.Set("cluster_spec_info", pulsarProClusterSpecInfoMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), pulsarProClusterInfoMap); e != nil {
			return e
		}
	}
	return nil
}
