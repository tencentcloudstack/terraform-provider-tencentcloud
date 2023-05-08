/*
Use this data source to query detailed information of tsf cluster

Example Usage

```hcl
data "tencentcloud_tsf_cluster" "cluster" {
  cluster_id_list = ["cluster-vwgj5e6y"]
  cluster_type = "V"
  # search_word = ""
  disable_program_auth_check = true
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTsfCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTsfClusterRead,
		Schema: map[string]*schema.Schema{
			"cluster_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Cluster ID list to be queried, if not filled in or passed, all content will be queried.",
			},

			"cluster_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The type of cluster to be queried, if left blank or not passed, all content will be queried. C: container, V: virtual machine.",
			},

			"search_word": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter by keywords for Cluster Id or name.",
			},

			"disable_program_auth_check": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to disable dataset authentication.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "TSF cluster pagination object. Note: This field may return null, indicating no valid value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of items. Note: This field may return null, indicating that no valid value was found.",
						},
						"content": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cluster list. Note: This field may return null, indicating no valid values.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster ID. Note: This field may return null, indicating no valid value.",
									},
									"cluster_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster name. Note: This field may return null, indicating no valid value.",
									},
									"cluster_desc": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster description. Note: This field may return null, indicating no valid value.",
									},
									"cluster_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster type. Note: This field may return null, indicating no valid value.",
									},
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Private network ID of the cluster. Note: This field may return null, indicating no valid value.",
									},
									"cluster_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "cluster status. Note: This field may return null, indicating no valid value.",
									},
									"cluster_cidr": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "cluster CIDR. Note: This field may return null, indicating no valid value.",
									},
									"cluster_total_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Total CPU of the cluster, unit: cores. Note: This field may return null, indicating that no valid value was found.",
									},
									"cluster_total_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Total memory of the cluster, unit: G. Note: This field may return null, indicating that no valid value is obtained.",
									},
									"cluster_used_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Used CPU of the cluster, in cores. This field may return null, indicating no valid value.",
									},
									"cluster_used_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Cluster used memory in GB. This field may return null, indicating no valid value.",
									},
									"instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cluster instance number. This field may return null, indicating no valid value.",
									},
									"run_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cluster running instance number. This field may return null, indicating no valid value.",
									},
									"normal_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cluster normal instance number. This field may return null, indicating no valid value.",
									},
									"delete_flag": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Deletion tag: true means it can be deleted, false means it cannot be deleted. Note: This field may return null, indicating no valid value.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "CreationTime. Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"update_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "last update time.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"tsf_region_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "region ID of TSF.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"tsf_region_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "region name of TSF.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"tsf_zone_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone Id of TSF.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"tsf_zone_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Zone name of TSF.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"delete_flag_reason": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Reason why the cluster cannot be deleted.  Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"cluster_limit_cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Maximum CPU limit of the cluster, in cores. This field may return null, indicating that no valid value was found.",
									},
									"cluster_limit_mem": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Cluster maximum memory limit in GB. This field may return null, indicating that no valid value was found.",
									},
									"run_service_instance_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of available service instances in the cluster. Note: This field may return null, indicating no valid value.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cluster subnet ID. Note: This field may return null, indicating no valid values.",
									},
									"operation_info": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Control information returned to the frontend. This field may return null, indicating no valid value.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"init": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Control information of the initialization button returned to the front end. Note: This field may return null, indicating no valid value.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disabled_reason": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Reason for not displaying. Note: This field may return null, indicating no valid value.",
															},
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "The availability of the button (whether it is clickable) may return null indicating that the information is not available.",
															},
															"supported": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether to display the button. Note: This field may return null, indicating that no valid value was found.",
															},
														},
													},
												},
												"add_instance": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Add instance button control information, Note: This field may return null, indicating that no valid value is obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disabled_reason": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The reason why this button is not displayed, may return null if not applicable.",
															},
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the button is clickable. Note: This field may return null, indicating that no valid value is obtained.",
															},
															"supported": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the button is clickable. Note: This field may return null, indicating that no valid value was found.",
															},
														},
													},
												},
												"destroy": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Control information for destroying machine, may return null if no valid value is obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disabled_reason": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The reason why this button is not displayed, may return null if not applicable.",
															},
															"enabled": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the button is clickable. Note: This field may return null, indicating that no valid value is obtained.",
															},
															"supported": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether the button is clickable. Note: This field may return null, indicating that no valid value was found.",
															},
														},
													},
												},
											},
										},
									},
									"cluster_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The cluster version, may return null if not applicable.",
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

func dataSourceTencentCloudTsfClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tsf_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id_list"); ok {
		clusterIdListSet := v.(*schema.Set).List()
		paramMap["ClusterIdList"] = helper.InterfacesStringsPoint(clusterIdListSet)
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		paramMap["ClusterType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_word"); ok {
		paramMap["SearchWord"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("disable_program_auth_check"); v != nil {
		paramMap["DisableProgramAuthCheck"] = helper.Bool(v.(bool))
	}

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var cluster *tsf.TsfPageCluster
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTsfClusterByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		cluster = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(cluster.Content))
	tsfPageClusterMap := map[string]interface{}{}
	if cluster != nil {
		if cluster.TotalCount != nil {
			tsfPageClusterMap["total_count"] = cluster.TotalCount
		}

		if cluster.Content != nil {
			contentList := []interface{}{}
			for _, content := range cluster.Content {
				contentMap := map[string]interface{}{}

				if content.ClusterId != nil {
					contentMap["cluster_id"] = content.ClusterId
				}

				if content.ClusterName != nil {
					contentMap["cluster_name"] = content.ClusterName
				}

				if content.ClusterDesc != nil {
					contentMap["cluster_desc"] = content.ClusterDesc
				}

				if content.ClusterType != nil {
					contentMap["cluster_type"] = content.ClusterType
				}

				if content.VpcId != nil {
					contentMap["vpc_id"] = content.VpcId
				}

				if content.ClusterStatus != nil {
					contentMap["cluster_status"] = content.ClusterStatus
				}

				if content.ClusterCIDR != nil {
					contentMap["cluster_cidr"] = content.ClusterCIDR
				}

				if content.ClusterTotalCpu != nil {
					contentMap["cluster_total_cpu"] = content.ClusterTotalCpu
				}

				if content.ClusterTotalMem != nil {
					contentMap["cluster_total_mem"] = content.ClusterTotalMem
				}

				if content.ClusterUsedCpu != nil {
					contentMap["cluster_used_cpu"] = content.ClusterUsedCpu
				}

				if content.ClusterUsedMem != nil {
					contentMap["cluster_used_mem"] = content.ClusterUsedMem
				}

				if content.InstanceCount != nil {
					contentMap["instance_count"] = content.InstanceCount
				}

				if content.RunInstanceCount != nil {
					contentMap["run_instance_count"] = content.RunInstanceCount
				}

				if content.NormalInstanceCount != nil {
					contentMap["normal_instance_count"] = content.NormalInstanceCount
				}

				if content.DeleteFlag != nil {
					contentMap["delete_flag"] = content.DeleteFlag
				}

				if content.CreateTime != nil {
					contentMap["create_time"] = content.CreateTime
				}

				if content.UpdateTime != nil {
					contentMap["update_time"] = content.UpdateTime
				}

				if content.TsfRegionId != nil {
					contentMap["tsf_region_id"] = content.TsfRegionId
				}

				if content.TsfRegionName != nil {
					contentMap["tsf_region_name"] = content.TsfRegionName
				}

				if content.TsfZoneId != nil {
					contentMap["tsf_zone_id"] = content.TsfZoneId
				}

				if content.TsfZoneName != nil {
					contentMap["tsf_zone_name"] = content.TsfZoneName
				}

				if content.DeleteFlagReason != nil {
					contentMap["delete_flag_reason"] = content.DeleteFlagReason
				}

				if content.ClusterLimitCpu != nil {
					contentMap["cluster_limit_cpu"] = content.ClusterLimitCpu
				}

				if content.ClusterLimitMem != nil {
					contentMap["cluster_limit_mem"] = content.ClusterLimitMem
				}

				if content.RunServiceInstanceCount != nil {
					contentMap["run_service_instance_count"] = content.RunServiceInstanceCount
				}

				if content.SubnetId != nil {
					contentMap["subnet_id"] = content.SubnetId
				}

				if content.OperationInfo != nil {
					operationInfoMap := map[string]interface{}{}

					if content.OperationInfo.Init != nil {
						initMap := map[string]interface{}{}

						if content.OperationInfo.Init.DisabledReason != nil {
							initMap["disabled_reason"] = content.OperationInfo.Init.DisabledReason
						}

						if content.OperationInfo.Init.Enabled != nil {
							initMap["enabled"] = content.OperationInfo.Init.Enabled
						}

						if content.OperationInfo.Init.Supported != nil {
							initMap["supported"] = content.OperationInfo.Init.Supported
						}

						operationInfoMap["init"] = []interface{}{initMap}
					}

					if content.OperationInfo.AddInstance != nil {
						addInstanceMap := map[string]interface{}{}

						if content.OperationInfo.AddInstance.DisabledReason != nil {
							addInstanceMap["disabled_reason"] = content.OperationInfo.AddInstance.DisabledReason
						}

						if content.OperationInfo.AddInstance.Enabled != nil {
							addInstanceMap["enabled"] = content.OperationInfo.AddInstance.Enabled
						}

						if content.OperationInfo.AddInstance.Supported != nil {
							addInstanceMap["supported"] = content.OperationInfo.AddInstance.Supported
						}

						operationInfoMap["add_instance"] = []interface{}{addInstanceMap}
					}

					if content.OperationInfo.Destroy != nil {
						destroyMap := map[string]interface{}{}

						if content.OperationInfo.Destroy.DisabledReason != nil {
							destroyMap["disabled_reason"] = content.OperationInfo.Destroy.DisabledReason
						}

						if content.OperationInfo.Destroy.Enabled != nil {
							destroyMap["enabled"] = content.OperationInfo.Destroy.Enabled
						}

						if content.OperationInfo.Destroy.Supported != nil {
							destroyMap["supported"] = content.OperationInfo.Destroy.Supported
						}

						operationInfoMap["destroy"] = []interface{}{destroyMap}
					}

					contentMap["operation_info"] = []interface{}{operationInfoMap}
				}

				if content.ClusterVersion != nil {
					contentMap["cluster_version"] = content.ClusterVersion
				}

				contentList = append(contentList, contentMap)
				ids = append(ids, *content.ClusterId)
			}

			tsfPageClusterMap["content"] = contentList
		}

		_ = d.Set("result", []interface{}{tsfPageClusterMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tsfPageClusterMap); e != nil {
			return e
		}
	}
	return nil
}
