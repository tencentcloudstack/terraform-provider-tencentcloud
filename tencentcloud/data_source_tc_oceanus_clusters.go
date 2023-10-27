/*
Use this data source to query detailed information of oceanus clusters

Example Usage

Query all clusters

```hcl
data "tencentcloud_oceanus_clusters" "example" {}
```

Query the specified cluster

```hcl
data "tencentcloud_oceanus_clusters" "example" {
  cluster_ids = ["cluster-5c42n3a5"]
  order_type  = 1
  filters {
    name   = "name"
    values = ["tf_example"]
  }
  work_space_id = "space-2idq8wbr"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudOceanusClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_ids": {
				Optional:    true,
				Type:        schema.TypeSet,
				MaxItems:    100,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Query one or more clusters by their ID. The maximum number of clusters that can be queried at once is 100.",
			},
			"order_type": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateAllowedIntValue(CLUSTER_ORDER_TYPE),
				Description:  "The sorting rule of the cluster information results. Possible values are 1 (sort by time in descending order), 2 (sort by time in ascending order), and 3 (sort by status).",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The filtering rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The field to be filtered.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "The filtering values of the field.",
						},
					},
				},
			},
			"work_space_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Workspace SerialId.",
			},
			"cluster_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the cluster.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the cluster.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The region where the cluster is located.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The user AppID.",
						},
						"owner_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The main account UIN.",
						},
						"creator_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creator UIN.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the cluster. Possible values are 1 (uninitialized), 3 (initializing), and 2 (running).",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A description of the cluster.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the cluster was created.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time of the last operation on the cluster.",
						},
						"cu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of CUs.",
						},
						"cu_mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The memory specification of the CU.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The availability zone.",
						},
						"status_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status description.",
						},
						"ccns": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The network.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the VPC.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the subnet.",
									},
									"ccn_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the Cloud Connect Network (CCN), such as ccn-rahigzjd.",
									},
								},
							},
						},
						"net_environment_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The network.",
						},
						"free_cu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of free CUs.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The tags bound to the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The tag key.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"tag_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The tag value.Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"isolated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the cluster was isolated. If the cluster has not been isolated, this field will be -.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expiration time of the cluster. If the cluster does not have an expiration time, this field will be -.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"seconds_until_expiry": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The number of seconds until the cluster expires. If the cluster does not have an expiration time, this field will be -.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The auto-renewal flag. 0 indicates the default state (the user has not set it, which is the initial state; if the user has enabled the prepaid non-stop privilege, the cluster will be automatically renewed), 1 indicates automatic renewal, and 2 indicates no automatic renewal (set by the user).Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"default_cos_bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default COS bucket of the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"cls_log_set": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CLS logset of the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"cls_topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CLS topic ID of the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"cls_log_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the CLS logset of the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"cls_topic_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the CLS topic of the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"version": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The version information of the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flink": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Flink version of the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"supported_flink": {
										Type:        schema.TypeSet,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: "The Flink versions supported by the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"free_cu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The number of free CUs at the granularity level.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"default_log_collect_conf": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default log collection configuration of the cluster.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"customized_dns_enabled": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Value: 0 - not set, 1 - set, 2 - not allowed to set.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"correlations": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Space information.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_group_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cluster ID.",
									},
									"cluster_group_serial_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster SerialId.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster name.",
									},
									"work_space_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace SerialId.",
									},
									"work_space_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace name.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Binding status. 2 - bound, 1 - unbound.",
									},
									"project_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Project ID.",
									},
									"project_id_str": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Project ID in string format.Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"running_cu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Running CU.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0 - postpaid, 1 - prepaid.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"is_need_manage_node": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Front-end distinguishes whether the cluster needs 2CU logic, because historical clusters do not need to be changed. Default is 1. All new clusters need to be changed.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"cluster_sessions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Session cluster information.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"arch_generation": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "V3 version = 2.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"cluster_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0: TKE, 1: EKS.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"orders": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Order information.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 - create, 2 - renew, 3 - scale.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"auto_renew_flag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "1 - auto-renewal.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"operate_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "UIN of the operator.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"compute_cu": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of CUs in the final cluster.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"order_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The time of the order.Note: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"sql_gateways": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Gateway information.Note: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"serial_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Unique identifier.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"flink_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flink kernel version.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Status. 1 - stopped, 2 - starting, 3 - started, 4 - start failed, 5 - stopping.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"creator_uin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creator.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"resource_refs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Reference resources.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"workspace_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique identifier of the space.",
												},
												"resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unique identifier of the resource.",
												},
												"version": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Version number.",
												},
												"type": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Reference type. 0: user resource.Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"cu_spec": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "CU specification.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Creation time.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Update time.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"properties": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Configuration parameters.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Key of the system configuration.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value of the system configuration.",
												},
											},
										},
									},
								},
							},
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

func dataSourceTencentCloudOceanusClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_oceanus_clusters.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = OceanusService{client: meta.(*TencentCloudClient).apiV3Conn}
		clusterSet []*oceanus.Cluster
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_ids"); ok {
		clusterIdsSet := v.(*schema.Set).List()
		paramMap["ClusterIds"] = helper.InterfacesStringsPoint(clusterIdsSet)
	}

	if v, ok := d.GetOkExists("order_type"); ok {
		paramMap["OrderType"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*oceanus.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := oceanus.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusClustersByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		clusterSet = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(clusterSet))
	tmpList := make([]map[string]interface{}, 0, len(clusterSet))

	if clusterSet != nil {
		for _, cluster := range clusterSet {
			clusterMap := map[string]interface{}{}

			if cluster.ClusterId != nil {
				clusterMap["cluster_id"] = cluster.ClusterId
			}

			if cluster.Name != nil {
				clusterMap["name"] = cluster.Name
			}

			if cluster.Region != nil {
				clusterMap["region"] = cluster.Region
			}

			if cluster.AppId != nil {
				clusterMap["app_id"] = cluster.AppId
			}

			if cluster.OwnerUin != nil {
				clusterMap["owner_uin"] = cluster.OwnerUin
			}

			if cluster.CreatorUin != nil {
				clusterMap["creator_uin"] = cluster.CreatorUin
			}

			if cluster.Status != nil {
				clusterMap["status"] = cluster.Status
			}

			if cluster.Remark != nil {
				clusterMap["remark"] = cluster.Remark
			}

			if cluster.CreateTime != nil {
				clusterMap["create_time"] = cluster.CreateTime
			}

			if cluster.UpdateTime != nil {
				clusterMap["update_time"] = cluster.UpdateTime
			}

			if cluster.CuNum != nil {
				clusterMap["cu_num"] = cluster.CuNum
			}

			if cluster.CuMem != nil {
				clusterMap["cu_mem"] = cluster.CuMem
			}

			if cluster.Zone != nil {
				clusterMap["zone"] = cluster.Zone
			}

			if cluster.StatusDesc != nil {
				clusterMap["status_desc"] = cluster.StatusDesc
			}

			if cluster.CCNs != nil {
				CCNsList := []interface{}{}
				for _, cCNs := range cluster.CCNs {
					cCNsMap := map[string]interface{}{}

					if cCNs.VpcId != nil {
						cCNsMap["vpc_id"] = cCNs.VpcId
					}

					if cCNs.SubnetId != nil {
						cCNsMap["subnet_id"] = cCNs.SubnetId
					}

					if cCNs.CcnId != nil {
						cCNsMap["ccn_id"] = cCNs.CcnId
					}

					CCNsList = append(CCNsList, cCNsMap)
				}

				clusterMap["ccns"] = CCNsList
			}

			if cluster.NetEnvironmentType != nil {
				clusterMap["net_environment_type"] = cluster.NetEnvironmentType
			}

			if cluster.FreeCuNum != nil {
				clusterMap["free_cu_num"] = cluster.FreeCuNum
			}

			if cluster.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range cluster.Tags {
					tagsMap := map[string]interface{}{}

					if tags.TagKey != nil {
						tagsMap["tag_key"] = tags.TagKey
					}

					if tags.TagValue != nil {
						tagsMap["tag_value"] = tags.TagValue
					}

					tagsList = append(tagsList, tagsMap)
				}

				clusterMap["tags"] = tagsList
			}

			if cluster.IsolatedTime != nil {
				clusterMap["isolated_time"] = cluster.IsolatedTime
			}

			if cluster.ExpireTime != nil {
				clusterMap["expire_time"] = cluster.ExpireTime
			}

			if cluster.SecondsUntilExpiry != nil {
				clusterMap["seconds_until_expiry"] = cluster.SecondsUntilExpiry
			}

			if cluster.AutoRenewFlag != nil {
				clusterMap["auto_renew_flag"] = cluster.AutoRenewFlag
			}

			if cluster.DefaultCOSBucket != nil {
				clusterMap["default_cos_bucket"] = cluster.DefaultCOSBucket
			}

			if cluster.CLSLogSet != nil {
				clusterMap["cls_log_set"] = cluster.CLSLogSet
			}

			if cluster.CLSTopicId != nil {
				clusterMap["cls_topic_id"] = cluster.CLSTopicId
			}

			if cluster.CLSLogName != nil {
				clusterMap["cls_log_name"] = cluster.CLSLogName
			}

			if cluster.CLSTopicName != nil {
				clusterMap["cls_topic_name"] = cluster.CLSTopicName
			}

			if cluster.Version != nil {
				versionMap := map[string]interface{}{}

				if cluster.Version.Flink != nil {
					versionMap["flink"] = cluster.Version.Flink
				}

				if cluster.Version.SupportedFlink != nil {
					versionMap["supported_flink"] = cluster.Version.SupportedFlink
				}

				clusterMap["version"] = []interface{}{versionMap}
			}

			if cluster.FreeCu != nil {
				clusterMap["free_cu"] = cluster.FreeCu
			}

			if cluster.DefaultLogCollectConf != nil {
				clusterMap["default_log_collect_conf"] = cluster.DefaultLogCollectConf
			}

			if cluster.CustomizedDNSEnabled != nil {
				clusterMap["customized_dns_enabled"] = cluster.CustomizedDNSEnabled
			}

			if cluster.Correlations != nil {
				correlationsList := []interface{}{}
				for _, correlations := range cluster.Correlations {
					correlationsMap := map[string]interface{}{}

					if correlations.ClusterGroupId != nil {
						correlationsMap["cluster_group_id"] = correlations.ClusterGroupId
					}

					if correlations.ClusterGroupSerialId != nil {
						correlationsMap["cluster_group_serial_id"] = correlations.ClusterGroupSerialId
					}

					if correlations.ClusterName != nil {
						correlationsMap["cluster_name"] = correlations.ClusterName
					}

					if correlations.WorkSpaceId != nil {
						correlationsMap["work_space_id"] = correlations.WorkSpaceId
					}

					if correlations.WorkSpaceName != nil {
						correlationsMap["work_space_name"] = correlations.WorkSpaceName
					}

					if correlations.Status != nil {
						correlationsMap["status"] = correlations.Status
					}

					if correlations.ProjectId != nil {
						correlationsMap["project_id"] = correlations.ProjectId
					}

					if correlations.ProjectIdStr != nil {
						correlationsMap["project_id_str"] = correlations.ProjectIdStr
					}

					correlationsList = append(correlationsList, correlationsMap)
				}

				clusterMap["correlations"] = correlationsList
			}

			if cluster.RunningCu != nil {
				clusterMap["running_cu"] = cluster.RunningCu
			}

			if cluster.PayMode != nil {
				clusterMap["pay_mode"] = cluster.PayMode
			}

			if cluster.IsNeedManageNode != nil {
				clusterMap["is_need_manage_node"] = cluster.IsNeedManageNode
			}

			if cluster.ClusterSessions != nil {
				tmpList = make([]map[string]interface{}, 0, len(cluster.ClusterSessions))
				//for _, item := range cluster.ClusterSessions {
				//	sessionMap := map[string]interface{}{}
				//	if item != nil {
				//
				//	}
				//
				//	tmpList = append(tmpList, sessionMap)
				//}

				clusterMap["cluster_sessions"] = tmpList
			}

			if cluster.ArchGeneration != nil {
				clusterMap["arch_generation"] = cluster.ArchGeneration
			}

			if cluster.ClusterType != nil {
				clusterMap["cluster_type"] = cluster.ClusterType
			}

			if cluster.Orders != nil {
				ordersList := []interface{}{}
				for _, orders := range cluster.Orders {
					ordersMap := map[string]interface{}{}

					if orders.Type != nil {
						ordersMap["type"] = orders.Type
					}

					if orders.AutoRenewFlag != nil {
						ordersMap["auto_renew_flag"] = orders.AutoRenewFlag
					}

					if orders.OperateUin != nil {
						ordersMap["operate_uin"] = orders.OperateUin
					}

					if orders.ComputeCu != nil {
						ordersMap["compute_cu"] = orders.ComputeCu
					}

					if orders.OrderTime != nil {
						ordersMap["order_time"] = orders.OrderTime
					}

					ordersList = append(ordersList, ordersMap)
				}

				clusterMap["orders"] = ordersList
			}

			if cluster.SqlGateways != nil {
				sqlGatewaysList := []interface{}{}
				for _, sqlGateways := range cluster.SqlGateways {
					sqlGatewaysMap := map[string]interface{}{}

					if sqlGateways.SerialId != nil {
						sqlGatewaysMap["serial_id"] = sqlGateways.SerialId
					}

					if sqlGateways.FlinkVersion != nil {
						sqlGatewaysMap["flink_version"] = sqlGateways.FlinkVersion
					}

					if sqlGateways.Status != nil {
						sqlGatewaysMap["status"] = sqlGateways.Status
					}

					if sqlGateways.CreatorUin != nil {
						sqlGatewaysMap["creator_uin"] = sqlGateways.CreatorUin
					}

					if sqlGateways.ResourceRefs != nil {
						resourceRefsList := []interface{}{}
						for _, resourceRefs := range sqlGateways.ResourceRefs {
							resourceRefsMap := map[string]interface{}{}

							if resourceRefs.WorkspaceId != nil {
								resourceRefsMap["workspace_id"] = resourceRefs.WorkspaceId
							}

							if resourceRefs.ResourceId != nil {
								resourceRefsMap["resource_id"] = resourceRefs.ResourceId
							}

							if resourceRefs.Version != nil {
								resourceRefsMap["version"] = resourceRefs.Version
							}

							if resourceRefs.Type != nil {
								resourceRefsMap["type"] = resourceRefs.Type
							}

							resourceRefsList = append(resourceRefsList, resourceRefsMap)
						}

						sqlGatewaysMap["resource_refs"] = resourceRefsList
					}

					if sqlGateways.CuSpec != nil {
						sqlGatewaysMap["cu_spec"] = sqlGateways.CuSpec
					}

					if sqlGateways.CreateTime != nil {
						sqlGatewaysMap["create_time"] = sqlGateways.CreateTime
					}

					if sqlGateways.UpdateTime != nil {
						sqlGatewaysMap["update_time"] = sqlGateways.UpdateTime
					}

					if sqlGateways.Properties != nil {
						propertiesList := []interface{}{}
						for _, properties := range sqlGateways.Properties {
							propertiesMap := map[string]interface{}{}

							if properties.Key != nil {
								propertiesMap["key"] = properties.Key
							}

							if properties.Value != nil {
								propertiesMap["value"] = properties.Value
							}

							propertiesList = append(propertiesList, propertiesMap)
						}

						sqlGatewaysMap["properties"] = propertiesList
					}

					sqlGatewaysList = append(sqlGatewaysList, sqlGatewaysMap)
				}

				clusterMap["sql_gateways"] = sqlGatewaysList
			}

			ids = append(ids, *cluster.ClusterId)
			tmpList = append(tmpList, clusterMap)
		}

		_ = d.Set("cluster_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
