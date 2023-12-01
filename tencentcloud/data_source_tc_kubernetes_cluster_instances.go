/*
Use this data source to query detailed information of kubernetes cluster_instances

Example Usage

```hcl
data "tencentcloud_kubernetes_cluster_instances" "cluster_instances" {
  cluster_id    = "cls-ely08ic4"
  instance_ids  = ["ins-kqmx8dm2"]
  instance_role = "WORKER"
  filters {
    name   = "nodepool-id"
    values = ["np-p4e6whqu"]
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kubernetes "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKubernetesClusterInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKubernetesClusterInstancesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the cluster.",
			},

			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of node instance IDs to be obtained. If it is empty, it means pulling all node instances in the cluster.",
			},

			"instance_role": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Node role, MASTER, WORKER, ETCD, MASTER_ETCD,ALL, default is WORKER.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "List of filter conditions. The optional values of Name are `nodepool-id` and `nodepool-instance-type`. Name is `nodepool-id`, which means filtering machines based on node pool id, and Value is the specific node pool id. Name is `nodepool-instance-type`, which indicates how the node is added to the node pool. Value is MANUALLY_ADDED (manually added to the node pool), AUTOSCALING_ADDED (joined by scaling group expansion method), ALL (manually join the node pool and join the node pool through scaling group expansion).",
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

			"instance_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of instances in the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node role, MASTER, WORKER, ETCD, MASTER_ETCD,ALL, default is WORKER.",
						},
						"failed_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reasons for instance exception (or being initialized).",
						},
						"instance_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the instance (running, initializing, failed).",
						},
						"drain_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the instance is blocked.",
						},
						"instance_advanced_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"desired_pod_number": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "When the node belongs to the podCIDR size customization mode, you can specify the upper limit of the number of pods running on the node.",
									},
									"gpu_args": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "GPU driver related parameters, obtain related GPU parameters: https://cloud.tencent.com/document/api/213/15715.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mig_enable": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to enable MIG features.",
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
																Description: "The name of the GPU driver or CUDA.",
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
																Description: "The name of the GPU driver or CUDA.",
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
																Description: "CuDNN name.",
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
									"pre_start_user_script": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded user script, executed before initializing the node, currently only effective for adding existing nodes.",
									},
									"taints": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Node taint.",
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
									"mount_target": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data disk mount point, the data disk is not mounted by default. Formatted ext3, ext4, xfs file system data disks will be mounted directly. Other file systems or unformatted data disks will be automatically formatted as ext4 (tlinux system formatted as xfs) and mounted. Please pay attention to backing up the data. This setting does not take effect for cloud hosts that have no data disks or multiple data disks.",
									},
									"docker_graph_path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Dockerd --graph specifies the value, the default is /var/lib/docker.",
									},
									"user_script": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Base64 encoded userscript.",
									},
									"unschedulable": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Set whether the added node participates in scheduling. The default value is 0, which means participating in scheduling; non-0 means not participating in scheduling. After the node initialization is completed, you can execute kubectl uncordon nodename to join the node in scheduling.",
									},
									"labels": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Node Label array.",
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
													Description: "Value in map table.",
												},
											},
										},
									},
									"data_disks": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Multi-disk data disk mounting information.",
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
													Description: "File system (ext3/ext4/xfs).",
												},
												"disk_size": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Cloud disk size (G).",
												},
												"auto_format_and_mount": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to automatically format the disk and mount it.",
												},
												"mount_target": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Mount directory.",
												},
												"disk_partition": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Mount device name or partition name, required when and only when adding an existing node.",
												},
											},
										},
									},
									"extra_args": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Node-related custom parameter information.",
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
								},
							},
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Add time.",
						},
						"lan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node intranet IP.",
						},
						"node_pool_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource pool ID.",
						},
						"autoscaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Auto scaling group ID.",
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

func dataSourceTencentCloudKubernetesClusterInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kubernetes_cluster_instances.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("instance_role"); ok {
		paramMap["InstanceRole"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*kubernetes.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := kubernetes.Filter{}
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

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	var instanceSet []*kubernetes.Instance
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterInstancesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceSet))
	tmpList := make([]map[string]interface{}, 0, len(instanceSet))

	if instanceSet != nil {
		for _, instance := range instanceSet {
			instanceMap := map[string]interface{}{}

			if instance.InstanceId != nil {
				instanceMap["instance_id"] = instance.InstanceId
			}

			if instance.InstanceRole != nil {
				instanceMap["instance_role"] = instance.InstanceRole
			}

			if instance.FailedReason != nil {
				instanceMap["failed_reason"] = instance.FailedReason
			}

			if instance.InstanceState != nil {
				instanceMap["instance_state"] = instance.InstanceState
			}

			if instance.DrainStatus != nil {
				instanceMap["drain_status"] = instance.DrainStatus
			}

			if instance.InstanceAdvancedSettings != nil {
				instanceAdvancedSettingsMap := map[string]interface{}{}

				if instance.InstanceAdvancedSettings.DesiredPodNumber != nil {
					instanceAdvancedSettingsMap["desired_pod_number"] = instance.InstanceAdvancedSettings.DesiredPodNumber
				}

				if instance.InstanceAdvancedSettings.GPUArgs != nil {
					gPUArgsMap := map[string]interface{}{}

					if instance.InstanceAdvancedSettings.GPUArgs.MIGEnable != nil {
						gPUArgsMap["mig_enable"] = instance.InstanceAdvancedSettings.GPUArgs.MIGEnable
					}

					if instance.InstanceAdvancedSettings.GPUArgs.Driver != nil {
						driverMap := map[string]interface{}{}

						if instance.InstanceAdvancedSettings.GPUArgs.Driver.Version != nil {
							driverMap["version"] = instance.InstanceAdvancedSettings.GPUArgs.Driver.Version
						}

						if instance.InstanceAdvancedSettings.GPUArgs.Driver.Name != nil {
							driverMap["name"] = instance.InstanceAdvancedSettings.GPUArgs.Driver.Name
						}

						gPUArgsMap["driver"] = []interface{}{driverMap}
					}

					if instance.InstanceAdvancedSettings.GPUArgs.CUDA != nil {
						cUDAMap := map[string]interface{}{}

						if instance.InstanceAdvancedSettings.GPUArgs.CUDA.Version != nil {
							cUDAMap["version"] = instance.InstanceAdvancedSettings.GPUArgs.CUDA.Version
						}

						if instance.InstanceAdvancedSettings.GPUArgs.CUDA.Name != nil {
							cUDAMap["name"] = instance.InstanceAdvancedSettings.GPUArgs.CUDA.Name
						}

						gPUArgsMap["cuda"] = []interface{}{cUDAMap}
					}

					if instance.InstanceAdvancedSettings.GPUArgs.CUDNN != nil {
						cUDNNMap := map[string]interface{}{}

						if instance.InstanceAdvancedSettings.GPUArgs.CUDNN.Version != nil {
							cUDNNMap["version"] = instance.InstanceAdvancedSettings.GPUArgs.CUDNN.Version
						}

						if instance.InstanceAdvancedSettings.GPUArgs.CUDNN.Name != nil {
							cUDNNMap["name"] = instance.InstanceAdvancedSettings.GPUArgs.CUDNN.Name
						}

						if instance.InstanceAdvancedSettings.GPUArgs.CUDNN.DocName != nil {
							cUDNNMap["doc_name"] = instance.InstanceAdvancedSettings.GPUArgs.CUDNN.DocName
						}

						if instance.InstanceAdvancedSettings.GPUArgs.CUDNN.DevName != nil {
							cUDNNMap["dev_name"] = instance.InstanceAdvancedSettings.GPUArgs.CUDNN.DevName
						}

						gPUArgsMap["cudnn"] = []interface{}{cUDNNMap}
					}

					if instance.InstanceAdvancedSettings.GPUArgs.CustomDriver != nil {
						customDriverMap := map[string]interface{}{}

						if instance.InstanceAdvancedSettings.GPUArgs.CustomDriver.Address != nil {
							customDriverMap["address"] = instance.InstanceAdvancedSettings.GPUArgs.CustomDriver.Address
						}

						gPUArgsMap["custom_driver"] = []interface{}{customDriverMap}
					}

					instanceAdvancedSettingsMap["gpu_args"] = []interface{}{gPUArgsMap}
				}

				if instance.InstanceAdvancedSettings.PreStartUserScript != nil {
					instanceAdvancedSettingsMap["pre_start_user_script"] = instance.InstanceAdvancedSettings.PreStartUserScript
				}

				if instance.InstanceAdvancedSettings.Taints != nil {
					var taintsList []interface{}
					for _, taints := range instance.InstanceAdvancedSettings.Taints {
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

					instanceAdvancedSettingsMap["taints"] = taintsList
				}

				if instance.InstanceAdvancedSettings.MountTarget != nil {
					instanceAdvancedSettingsMap["mount_target"] = instance.InstanceAdvancedSettings.MountTarget
				}

				if instance.InstanceAdvancedSettings.DockerGraphPath != nil {
					instanceAdvancedSettingsMap["docker_graph_path"] = instance.InstanceAdvancedSettings.DockerGraphPath
				}

				if instance.InstanceAdvancedSettings.UserScript != nil {
					instanceAdvancedSettingsMap["user_script"] = instance.InstanceAdvancedSettings.UserScript
				}

				if instance.InstanceAdvancedSettings.Unschedulable != nil {
					instanceAdvancedSettingsMap["unschedulable"] = instance.InstanceAdvancedSettings.Unschedulable
				}

				if instance.InstanceAdvancedSettings.Labels != nil {
					var labelsList []interface{}
					for _, labels := range instance.InstanceAdvancedSettings.Labels {
						labelsMap := map[string]interface{}{}

						if labels.Name != nil {
							labelsMap["name"] = labels.Name
						}

						if labels.Value != nil {
							labelsMap["value"] = labels.Value
						}

						labelsList = append(labelsList, labelsMap)
					}

					instanceAdvancedSettingsMap["labels"] = labelsList
				}

				if instance.InstanceAdvancedSettings.DataDisks != nil {
					var dataDisksList []interface{}
					for _, dataDisks := range instance.InstanceAdvancedSettings.DataDisks {
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

					instanceAdvancedSettingsMap["data_disks"] = dataDisksList
				}

				if instance.InstanceAdvancedSettings.ExtraArgs != nil {
					extraArgsMap := map[string]interface{}{}

					if instance.InstanceAdvancedSettings.ExtraArgs.Kubelet != nil {
						extraArgsMap["kubelet"] = instance.InstanceAdvancedSettings.ExtraArgs.Kubelet
					}

					instanceAdvancedSettingsMap["extra_args"] = []interface{}{extraArgsMap}
				}

				instanceMap["instance_advanced_settings"] = []interface{}{instanceAdvancedSettingsMap}
			}

			if instance.CreatedTime != nil {
				instanceMap["created_time"] = instance.CreatedTime
			}

			if instance.LanIP != nil {
				instanceMap["lan_ip"] = instance.LanIP
			}

			if instance.NodePoolId != nil {
				instanceMap["node_pool_id"] = instance.NodePoolId
			}

			if instance.AutoscalingGroupId != nil {
				instanceMap["autoscaling_group_id"] = instance.AutoscalingGroupId
			}

			ids = append(ids, *instance.InstanceId)
			tmpList = append(tmpList, instanceMap)
		}

		_ = d.Set("instance_set", tmpList)
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
