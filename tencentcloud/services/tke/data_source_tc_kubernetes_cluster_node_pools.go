package tke

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKubernetesClusterNodePools() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesClusterNodePoolsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the cluster.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "NodePoolsName, Filter according to the node pool name, type: String, required: no. NodePoolsId, Filter according to the node pool ID, type: String, required: no. tags, Filter according to the label key value pairs, type: String, required: no. tag:tag-key, Filter according to the label key value pairs, type: String, required: no.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The attribute name, if there are multiple filters, the relationship between the filters is a logical AND relationship.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Attribute values, if there are multiple values in the same filter, the relationship between values under the same filter is a logical OR relationship.",
						},
					},
				},
			},

			"node_pool_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Node Pool List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_pool_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the node pool.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the node pool.",
						},
						"cluster_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the cluster.",
						},
						"life_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Life cycle state of the node pool, include: creating, normal, updating, deleting, deleted.",
						},
						"launch_configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of launch configuration.",
						},
						"autoscaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of autoscaling group.",
						},
						"labels": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Labels of the node pool.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name in the map table.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value in the map table.",
									},
								},
							},
						},
						"taints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Labels of the node pool.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key of taints mark.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Value of taints mark.",
									},
									"effect": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Effect of taints mark.",
									},
								},
							},
						},
						"node_count_summary": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node List.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"manually_added": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Manually managed nodes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"joining": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of nodes joining.",
												},
												"initializing": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of nodes in initialization.",
												},
												"normal": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Normal number of nodes.",
												},
												"total": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total number of nodes.",
												},
											},
										},
									},
									"autoscaling_added": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Automatically managed nodes.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"joining": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of nodes joining.",
												},
												"initializing": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of nodes in initialization.",
												},
												"normal": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Normal number of nodes.",
												},
												"total": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total number of nodes.",
												},
											},
										},
									},
								},
							},
						},
						"autoscaling_group_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status information.",
						},
						"max_nodes_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of nodes.",
						},
						"min_nodes_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum number of nodes.",
						},
						"desired_nodes_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Expected number of nodes.",
						},
						"node_pool_os": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node Pool OS Name.",
						},
						"os_customize_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Mirror version of container.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of image.",
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "When the cluster belongs to the node podCIDR size customization mode, the node pool needs to have the pod number attribute.",
						},
						"user_script": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User defined scripts.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Resource tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Label key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Label value.",
									},
								},
							},
						},
						"deletion_protection": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Remove protection switch.",
						},
						"extra_args": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kubelet": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Kubelet custom parameters.",
									},
								},
							},
						},
						"gpu_args": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "GPU driver related parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mig_enable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the MIG feature enabled.",
									},
									"driver": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "GPU driver version information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "GPU driver or CUDA version.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "GPU driver or CUDA name.",
												},
											},
										},
									},
									"cuda": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CUDA version information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "GPU driver or CUDA version.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "GPU driver or CUDA name.",
												},
											},
										},
									},
									"cudnn": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "CuDNN version information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Version of cuDNN.",
												},
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of cuDNN.",
												},
												"doc_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Doc name of cuDNN.",
												},
												"dev_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dev name of cuDNN.",
												},
											},
										},
									},
									"custom_driver": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Custom GPU driver information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Custom GPU driver address link.",
												},
											},
										},
									},
								},
							},
						},
						"docker_graph_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dockerd --graph specified value, default to /var/lib/docker.",
						},
						"data_disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Multi disk data disk mounting information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Cloud disk type.",
									},
									"file_system": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "File system(ext3/ext4/xfs).",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Cloud disk size(G).",
									},
									"auto_format_and_mount": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to automate the format disk and mount it.",
									},
									"mount_target": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mount directory.",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Mount device name or partition name.",
									},
								},
							},
						},
						"unschedulable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is it not schedulable.",
						},
						"pre_start_user_script": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User defined script, executed before User Script.",
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

func dataSourceTencentCloudKubernetesClusterNodePoolsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tke_cluster_node_pools.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tke.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tke.Filter{}
			filterMap := item.(map[string]interface{})

			if v1, ok1 := filterMap["name"]; ok1 {
				filter.Name = helper.String(v1.(string))
			}
			if v1, ok1 := filterMap["values"]; ok1 {
				valuesSet := v1.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var nodePoolSet []*tke.NodePool
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterNodePoolsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		nodePoolSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(nodePoolSet))
	tmpList := make([]map[string]interface{}, 0, len(nodePoolSet))

	if nodePoolSet != nil {
		for _, nodePool := range nodePoolSet {
			nodePoolMap := map[string]interface{}{}

			if nodePool.NodePoolId != nil {
				nodePoolMap["node_pool_id"] = nodePool.NodePoolId
			}

			if nodePool.Name != nil {
				nodePoolMap["name"] = nodePool.Name
			}

			if nodePool.ClusterInstanceId != nil {
				nodePoolMap["cluster_instance_id"] = nodePool.ClusterInstanceId
			}

			if nodePool.LifeState != nil {
				nodePoolMap["life_state"] = nodePool.LifeState
			}

			if nodePool.LaunchConfigurationId != nil {
				nodePoolMap["launch_configuration_id"] = nodePool.LaunchConfigurationId
			}

			if nodePool.AutoscalingGroupId != nil {
				nodePoolMap["autoscaling_group_id"] = nodePool.AutoscalingGroupId
			}

			if nodePool.Labels != nil {
				var labelsList []interface{}
				for _, labels := range nodePool.Labels {
					labelsMap := map[string]interface{}{}

					if labels.Name != nil {
						labelsMap["name"] = labels.Name
					}

					if labels.Value != nil {
						labelsMap["value"] = labels.Value
					}

					labelsList = append(labelsList, labelsMap)
				}

				nodePoolMap["labels"] = labelsList
			}

			if nodePool.Taints != nil {
				var taintsList []interface{}
				for _, taints := range nodePool.Taints {
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

				nodePoolMap["taints"] = taintsList
			}

			if nodePool.NodeCountSummary != nil {
				nodeCountSummaryMap := map[string]interface{}{}

				if nodePool.NodeCountSummary.ManuallyAdded != nil {
					manuallyAddedMap := map[string]interface{}{}

					if nodePool.NodeCountSummary.ManuallyAdded.Joining != nil {
						manuallyAddedMap["joining"] = nodePool.NodeCountSummary.ManuallyAdded.Joining
					}

					if nodePool.NodeCountSummary.ManuallyAdded.Initializing != nil {
						manuallyAddedMap["initializing"] = nodePool.NodeCountSummary.ManuallyAdded.Initializing
					}

					if nodePool.NodeCountSummary.ManuallyAdded.Normal != nil {
						manuallyAddedMap["normal"] = nodePool.NodeCountSummary.ManuallyAdded.Normal
					}

					if nodePool.NodeCountSummary.ManuallyAdded.Total != nil {
						manuallyAddedMap["total"] = nodePool.NodeCountSummary.ManuallyAdded.Total
					}

					nodeCountSummaryMap["manually_added"] = []interface{}{manuallyAddedMap}
				}

				if nodePool.NodeCountSummary.AutoscalingAdded != nil {
					autoscalingAddedMap := map[string]interface{}{}

					if nodePool.NodeCountSummary.AutoscalingAdded.Joining != nil {
						autoscalingAddedMap["joining"] = nodePool.NodeCountSummary.AutoscalingAdded.Joining
					}

					if nodePool.NodeCountSummary.AutoscalingAdded.Initializing != nil {
						autoscalingAddedMap["initializing"] = nodePool.NodeCountSummary.AutoscalingAdded.Initializing
					}

					if nodePool.NodeCountSummary.AutoscalingAdded.Normal != nil {
						autoscalingAddedMap["normal"] = nodePool.NodeCountSummary.AutoscalingAdded.Normal
					}

					if nodePool.NodeCountSummary.AutoscalingAdded.Total != nil {
						autoscalingAddedMap["total"] = nodePool.NodeCountSummary.AutoscalingAdded.Total
					}

					nodeCountSummaryMap["autoscaling_added"] = []interface{}{autoscalingAddedMap}
				}

				nodePoolMap["node_count_summary"] = []interface{}{nodeCountSummaryMap}
			}

			if nodePool.AutoscalingGroupStatus != nil {
				nodePoolMap["autoscaling_group_status"] = nodePool.AutoscalingGroupStatus
			}

			if nodePool.MaxNodesNum != nil {
				nodePoolMap["max_nodes_num"] = nodePool.MaxNodesNum
			}

			if nodePool.MinNodesNum != nil {
				nodePoolMap["min_nodes_num"] = nodePool.MinNodesNum
			}

			if nodePool.DesiredNodesNum != nil {
				nodePoolMap["desired_nodes_num"] = nodePool.DesiredNodesNum
			}

			if nodePool.NodePoolOs != nil {
				nodePoolMap["node_pool_os"] = nodePool.NodePoolOs
			}

			if nodePool.OsCustomizeType != nil {
				nodePoolMap["os_customize_type"] = nodePool.OsCustomizeType
			}

			if nodePool.ImageId != nil {
				nodePoolMap["image_id"] = nodePool.ImageId
			}

			if nodePool.DesiredPodNum != nil {
				nodePoolMap["desired_pod_num"] = nodePool.DesiredPodNum
			}

			if nodePool.UserScript != nil {
				nodePoolMap["user_script"] = nodePool.UserScript
			}

			if nodePool.Tags != nil {
				var tagsList []interface{}
				for _, tags := range nodePool.Tags {
					tagsMap := map[string]interface{}{}

					if tags.Key != nil {
						tagsMap["key"] = tags.Key
					}

					if tags.Value != nil {
						tagsMap["value"] = tags.Value
					}

					tagsList = append(tagsList, tagsMap)
				}

				nodePoolMap["tags"] = tagsList
			}

			if nodePool.DeletionProtection != nil {
				nodePoolMap["deletion_protection"] = nodePool.DeletionProtection
			}

			if nodePool.ExtraArgs != nil {
				extraArgsMap := map[string]interface{}{}

				if nodePool.ExtraArgs.Kubelet != nil {
					extraArgsMap["kubelet"] = nodePool.ExtraArgs.Kubelet
				}

				nodePoolMap["extra_args"] = []interface{}{extraArgsMap}
			}

			if nodePool.GPUArgs != nil {
				gPUArgsMap := map[string]interface{}{}

				if nodePool.GPUArgs.MIGEnable != nil {
					gPUArgsMap["mig_enable"] = nodePool.GPUArgs.MIGEnable
				}

				if nodePool.GPUArgs.Driver != nil {
					driverMap := map[string]interface{}{}

					if nodePool.GPUArgs.Driver.Version != nil {
						driverMap["version"] = nodePool.GPUArgs.Driver.Version
					}

					if nodePool.GPUArgs.Driver.Name != nil {
						driverMap["name"] = nodePool.GPUArgs.Driver.Name
					}

					gPUArgsMap["driver"] = []interface{}{driverMap}
				}

				if nodePool.GPUArgs.CUDA != nil {
					cUDAMap := map[string]interface{}{}

					if nodePool.GPUArgs.CUDA.Version != nil {
						cUDAMap["version"] = nodePool.GPUArgs.CUDA.Version
					}

					if nodePool.GPUArgs.CUDA.Name != nil {
						cUDAMap["name"] = nodePool.GPUArgs.CUDA.Name
					}

					gPUArgsMap["cuda"] = []interface{}{cUDAMap}
				}

				if nodePool.GPUArgs.CUDNN != nil {
					cUDNNMap := map[string]interface{}{}

					if nodePool.GPUArgs.CUDNN.Version != nil {
						cUDNNMap["version"] = nodePool.GPUArgs.CUDNN.Version
					}

					if nodePool.GPUArgs.CUDNN.Name != nil {
						cUDNNMap["name"] = nodePool.GPUArgs.CUDNN.Name
					}

					if nodePool.GPUArgs.CUDNN.DocName != nil {
						cUDNNMap["doc_name"] = nodePool.GPUArgs.CUDNN.DocName
					}

					if nodePool.GPUArgs.CUDNN.DevName != nil {
						cUDNNMap["dev_name"] = nodePool.GPUArgs.CUDNN.DevName
					}

					gPUArgsMap["cudnn"] = []interface{}{cUDNNMap}
				}

				if nodePool.GPUArgs.CustomDriver != nil {
					customDriverMap := map[string]interface{}{}

					if nodePool.GPUArgs.CustomDriver.Address != nil {
						customDriverMap["address"] = nodePool.GPUArgs.CustomDriver.Address
					}

					gPUArgsMap["custom_driver"] = []interface{}{customDriverMap}
				}

				nodePoolMap["gpu_args"] = []interface{}{gPUArgsMap}
			}

			if nodePool.DockerGraphPath != nil {
				nodePoolMap["docker_graph_path"] = nodePool.DockerGraphPath
			}

			if nodePool.DataDisks != nil {
				var dataDisksList []interface{}
				for _, dataDisks := range nodePool.DataDisks {
					dataDisksMap := map[string]interface{}{}

					if dataDisks.DiskType != nil {
						dataDisksMap["disk_type"] = dataDisks.DiskType
					}

					if dataDisks.FileSystem != nil {
						dataDisksMap["file_system"] = dataDisks.FileSystem
					}

					if dataDisks.DiskSize != nil {
						dataDisksMap["disk_size"] = dataDisks.DiskSize
					}

					if dataDisks.AutoFormatAndMount != nil {
						dataDisksMap["auto_format_and_mount"] = dataDisks.AutoFormatAndMount
					}

					if dataDisks.MountTarget != nil {
						dataDisksMap["mount_target"] = dataDisks.MountTarget
					}

					if dataDisks.DiskPartition != nil {
						dataDisksMap["disk_partition"] = dataDisks.DiskPartition
					}

					dataDisksList = append(dataDisksList, dataDisksMap)
				}

				nodePoolMap["data_disks"] = dataDisksList
			}

			if nodePool.Unschedulable != nil {
				nodePoolMap["unschedulable"] = nodePool.Unschedulable
			}

			if nodePool.PreStartUserScript != nil {
				nodePoolMap["pre_start_user_script"] = nodePool.PreStartUserScript
			}

			ids = append(ids, *nodePool.NodePoolId)
			tmpList = append(tmpList, nodePoolMap)
		}

		_ = d.Set("node_pool_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
