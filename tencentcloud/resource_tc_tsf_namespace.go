/*
Provides a resource to create a tsf namespace

Example Usage

```hcl
resource "tencentcloud_tsf_namespace" "namespace" {
  namespace_name = ""
  cluster_id = ""
  namespace_desc = ""
  namespace_resource_type = ""
  namespace_type = ""
  namespace_id = ""
  is_ha_enable = ""
  program_id = ""
    program_id_list =
              }
```

Import

tsf namespace can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_namespace.namespace namespace_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTsfNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfNamespaceCreate,
		Read:   resourceTencentCloudTsfNamespaceRead,
		Update: resourceTencentCloudTsfNamespaceUpdate,
		Delete: resourceTencentCloudTsfNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"namespace_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace name.",
			},

			"cluster_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"namespace_desc": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Namespace description.",
			},

			"namespace_resource_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Namespace resource type (default is DEF).",
			},

			"namespace_type": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Whether it is a global namespace (the default is DEF, which means a common namespace; GLOBAL means a global namespace).",
			},

			"namespace_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Namespace ID.",
			},

			"is_ha_enable": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Whether to enable high availability.",
			},

			"program_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the dataset to be bound.",
			},

			"kube_inject_enable": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "KubeInjectEnable value.",
			},

			"program_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program id list.",
			},

			"namespace_code": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Namespace encoding.",
			},

			"is_default": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Default namespace.",
			},

			"namespace_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Namespace status.",
			},

			"delete_flag": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Delete ID.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creation time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"cluster_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster array, only carrying basic information such as cluster ID, cluster name, and cluster type.",
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
							Description: "Cluster name.",
						},
						"cluster_desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster description.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster type.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the private network to which the cluster belongs.",
						},
						"cluster_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster status.",
						},
						"cluster_c_i_d_r": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster CIDR.",
						},
						"cluster_total_cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total CPU of the cluster, unit: core.",
						},
						"cluster_total_mem": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total memory of the cluster, unit: G.",
						},
						"cluster_used_cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The CPU used by the cluster, unit: core.",
						},
						"cluster_used_mem": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The cluster has used memory, unit: G.",
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of cluster machine instances.",
						},
						"run_instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of machine instances available in the cluster.",
						},
						"normal_instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of machine instances in the normal state of the cluster.",
						},
						"delete_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Delete flag: true: can be deleted; false: can not be deleted.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"tsf_region_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the TSF region to which the cluster belongs.",
						},
						"tsf_region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the TSF region to which the cluster belongs.",
						},
						"tsf_zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the TSF availability zone to which the cluster belongs.",
						},
						"tsf_zone_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the TSF availability zone to which the cluster belongs.",
						},
						"delete_flag_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason why the cluster cannot be deleted.",
						},
						"cluster_limit_cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The maximum CPU limit of the cluster, unit: core.",
						},
						"cluster_limit_mem": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The maximum memory limit of the cluster, unit: G.",
						},
						"run_service_instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of service instances available in the cluster.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet to which the cluster belongs.",
						},
						"operation_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Control information returned to the front end.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"init": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Initialize the control information of the button.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disabled_reason": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The reason for not displaying.",
												},
												"enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the button is clickable.",
												},
												"supported": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to display the button.",
												},
											},
										},
									},
									"add_instance": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Add the control information of the instance button.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disabled_reason": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The reason for not displaying.",
												},
												"enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the button is clickable.",
												},
												"supported": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to display the button.",
												},
											},
										},
									},
									"destroy": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Destroy the control information of the machine.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disabled_reason": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The reason for not displaying.",
												},
												"enabled": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether the button is clickable.",
												},
												"supported": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to display the button.",
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
							Description: "Cluster version.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTsfNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_namespace.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = tsf.NewCreateNamespaceRequest()
		response    = tsf.NewCreateNamespaceResponse()
		namespaceId string
	)
	if v, ok := d.GetOk("namespace_name"); ok {
		request.NamespaceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_id"); ok {
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_desc"); ok {
		request.NamespaceDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_resource_type"); ok {
		request.NamespaceResourceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_type"); ok {
		request.NamespaceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		namespaceId = v.(string)
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("is_ha_enable"); ok {
		request.IsHaEnable = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_id"); ok {
		request.ProgramId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateNamespace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf namespace failed, reason:%+v", logId, err)
		return err
	}

	namespaceId = *response.Response.namespaceId
	d.SetId(namespaceId)

	return resourceTencentCloudTsfNamespaceRead(d, meta)
}

func resourceTencentCloudTsfNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_namespace.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	namespaceId := d.Id()

	namespace, err := service.DescribeTsfNamespaceById(ctx, namespaceId)
	if err != nil {
		return err
	}

	if namespace == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfNamespace` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if namespace.NamespaceName != nil {
		_ = d.Set("namespace_name", namespace.NamespaceName)
	}

	if namespace.ClusterId != nil {
		_ = d.Set("cluster_id", namespace.ClusterId)
	}

	if namespace.NamespaceDesc != nil {
		_ = d.Set("namespace_desc", namespace.NamespaceDesc)
	}

	if namespace.NamespaceResourceType != nil {
		_ = d.Set("namespace_resource_type", namespace.NamespaceResourceType)
	}

	if namespace.NamespaceType != nil {
		_ = d.Set("namespace_type", namespace.NamespaceType)
	}

	if namespace.NamespaceId != nil {
		_ = d.Set("namespace_id", namespace.NamespaceId)
	}

	if namespace.IsHaEnable != nil {
		_ = d.Set("is_ha_enable", namespace.IsHaEnable)
	}

	if namespace.ProgramId != nil {
		_ = d.Set("program_id", namespace.ProgramId)
	}

	if namespace.KubeInjectEnable != nil {
		_ = d.Set("kube_inject_enable", namespace.KubeInjectEnable)
	}

	if namespace.ProgramIdList != nil {
		_ = d.Set("program_id_list", namespace.ProgramIdList)
	}

	if namespace.NamespaceCode != nil {
		_ = d.Set("namespace_code", namespace.NamespaceCode)
	}

	if namespace.IsDefault != nil {
		_ = d.Set("is_default", namespace.IsDefault)
	}

	if namespace.NamespaceStatus != nil {
		_ = d.Set("namespace_status", namespace.NamespaceStatus)
	}

	if namespace.DeleteFlag != nil {
		_ = d.Set("delete_flag", namespace.DeleteFlag)
	}

	if namespace.CreateTime != nil {
		_ = d.Set("create_time", namespace.CreateTime)
	}

	if namespace.UpdateTime != nil {
		_ = d.Set("update_time", namespace.UpdateTime)
	}

	if namespace.ClusterList != nil {
		clusterListList := []interface{}{}
		for _, clusterList := range namespace.ClusterList {
			clusterListMap := map[string]interface{}{}

			if namespace.ClusterList.ClusterId != nil {
				clusterListMap["cluster_id"] = namespace.ClusterList.ClusterId
			}

			if namespace.ClusterList.ClusterName != nil {
				clusterListMap["cluster_name"] = namespace.ClusterList.ClusterName
			}

			if namespace.ClusterList.ClusterDesc != nil {
				clusterListMap["cluster_desc"] = namespace.ClusterList.ClusterDesc
			}

			if namespace.ClusterList.ClusterType != nil {
				clusterListMap["cluster_type"] = namespace.ClusterList.ClusterType
			}

			if namespace.ClusterList.VpcId != nil {
				clusterListMap["vpc_id"] = namespace.ClusterList.VpcId
			}

			if namespace.ClusterList.ClusterStatus != nil {
				clusterListMap["cluster_status"] = namespace.ClusterList.ClusterStatus
			}

			if namespace.ClusterList.ClusterCIDR != nil {
				clusterListMap["cluster_c_i_d_r"] = namespace.ClusterList.ClusterCIDR
			}

			if namespace.ClusterList.ClusterTotalCpu != nil {
				clusterListMap["cluster_total_cpu"] = namespace.ClusterList.ClusterTotalCpu
			}

			if namespace.ClusterList.ClusterTotalMem != nil {
				clusterListMap["cluster_total_mem"] = namespace.ClusterList.ClusterTotalMem
			}

			if namespace.ClusterList.ClusterUsedCpu != nil {
				clusterListMap["cluster_used_cpu"] = namespace.ClusterList.ClusterUsedCpu
			}

			if namespace.ClusterList.ClusterUsedMem != nil {
				clusterListMap["cluster_used_mem"] = namespace.ClusterList.ClusterUsedMem
			}

			if namespace.ClusterList.InstanceCount != nil {
				clusterListMap["instance_count"] = namespace.ClusterList.InstanceCount
			}

			if namespace.ClusterList.RunInstanceCount != nil {
				clusterListMap["run_instance_count"] = namespace.ClusterList.RunInstanceCount
			}

			if namespace.ClusterList.NormalInstanceCount != nil {
				clusterListMap["normal_instance_count"] = namespace.ClusterList.NormalInstanceCount
			}

			if namespace.ClusterList.DeleteFlag != nil {
				clusterListMap["delete_flag"] = namespace.ClusterList.DeleteFlag
			}

			if namespace.ClusterList.CreateTime != nil {
				clusterListMap["create_time"] = namespace.ClusterList.CreateTime
			}

			if namespace.ClusterList.UpdateTime != nil {
				clusterListMap["update_time"] = namespace.ClusterList.UpdateTime
			}

			if namespace.ClusterList.TsfRegionId != nil {
				clusterListMap["tsf_region_id"] = namespace.ClusterList.TsfRegionId
			}

			if namespace.ClusterList.TsfRegionName != nil {
				clusterListMap["tsf_region_name"] = namespace.ClusterList.TsfRegionName
			}

			if namespace.ClusterList.TsfZoneId != nil {
				clusterListMap["tsf_zone_id"] = namespace.ClusterList.TsfZoneId
			}

			if namespace.ClusterList.TsfZoneName != nil {
				clusterListMap["tsf_zone_name"] = namespace.ClusterList.TsfZoneName
			}

			if namespace.ClusterList.DeleteFlagReason != nil {
				clusterListMap["delete_flag_reason"] = namespace.ClusterList.DeleteFlagReason
			}

			if namespace.ClusterList.ClusterLimitCpu != nil {
				clusterListMap["cluster_limit_cpu"] = namespace.ClusterList.ClusterLimitCpu
			}

			if namespace.ClusterList.ClusterLimitMem != nil {
				clusterListMap["cluster_limit_mem"] = namespace.ClusterList.ClusterLimitMem
			}

			if namespace.ClusterList.RunServiceInstanceCount != nil {
				clusterListMap["run_service_instance_count"] = namespace.ClusterList.RunServiceInstanceCount
			}

			if namespace.ClusterList.SubnetId != nil {
				clusterListMap["subnet_id"] = namespace.ClusterList.SubnetId
			}

			if namespace.ClusterList.OperationInfo != nil {
				operationInfoMap := map[string]interface{}{}

				if namespace.ClusterList.OperationInfo.Init != nil {
					initMap := map[string]interface{}{}

					if namespace.ClusterList.OperationInfo.Init.DisabledReason != nil {
						initMap["disabled_reason"] = namespace.ClusterList.OperationInfo.Init.DisabledReason
					}

					if namespace.ClusterList.OperationInfo.Init.Enabled != nil {
						initMap["enabled"] = namespace.ClusterList.OperationInfo.Init.Enabled
					}

					if namespace.ClusterList.OperationInfo.Init.Supported != nil {
						initMap["supported"] = namespace.ClusterList.OperationInfo.Init.Supported
					}

					operationInfoMap["init"] = []interface{}{initMap}
				}

				if namespace.ClusterList.OperationInfo.AddInstance != nil {
					addInstanceMap := map[string]interface{}{}

					if namespace.ClusterList.OperationInfo.AddInstance.DisabledReason != nil {
						addInstanceMap["disabled_reason"] = namespace.ClusterList.OperationInfo.AddInstance.DisabledReason
					}

					if namespace.ClusterList.OperationInfo.AddInstance.Enabled != nil {
						addInstanceMap["enabled"] = namespace.ClusterList.OperationInfo.AddInstance.Enabled
					}

					if namespace.ClusterList.OperationInfo.AddInstance.Supported != nil {
						addInstanceMap["supported"] = namespace.ClusterList.OperationInfo.AddInstance.Supported
					}

					operationInfoMap["add_instance"] = []interface{}{addInstanceMap}
				}

				if namespace.ClusterList.OperationInfo.Destroy != nil {
					destroyMap := map[string]interface{}{}

					if namespace.ClusterList.OperationInfo.Destroy.DisabledReason != nil {
						destroyMap["disabled_reason"] = namespace.ClusterList.OperationInfo.Destroy.DisabledReason
					}

					if namespace.ClusterList.OperationInfo.Destroy.Enabled != nil {
						destroyMap["enabled"] = namespace.ClusterList.OperationInfo.Destroy.Enabled
					}

					if namespace.ClusterList.OperationInfo.Destroy.Supported != nil {
						destroyMap["supported"] = namespace.ClusterList.OperationInfo.Destroy.Supported
					}

					operationInfoMap["destroy"] = []interface{}{destroyMap}
				}

				clusterListMap["operation_info"] = []interface{}{operationInfoMap}
			}

			if namespace.ClusterList.ClusterVersion != nil {
				clusterListMap["cluster_version"] = namespace.ClusterList.ClusterVersion
			}

			clusterListList = append(clusterListList, clusterListMap)
		}

		_ = d.Set("cluster_list", clusterListList)

	}

	return nil
}

func resourceTencentCloudTsfNamespaceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_namespace.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyNamespaceRequest()

	namespaceId := d.Id()

	request.NamespaceId = &namespaceId

	immutableArgs := []string{"namespace_name", "cluster_id", "namespace_desc", "namespace_resource_type", "namespace_type", "namespace_id", "is_ha_enable", "program_id", "kube_inject_enable", "program_id_list", "namespace_code", "is_default", "namespace_status", "delete_flag", "create_time", "update_time", "cluster_list"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("namespace_name") {
		if v, ok := d.GetOk("namespace_name"); ok {
			request.NamespaceName = helper.String(v.(string))
		}
	}

	if d.HasChange("namespace_desc") {
		if v, ok := d.GetOk("namespace_desc"); ok {
			request.NamespaceDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("namespace_id") {
		if v, ok := d.GetOk("namespace_id"); ok {
			request.NamespaceId = helper.String(v.(string))
		}
	}

	if d.HasChange("is_ha_enable") {
		if v, ok := d.GetOk("is_ha_enable"); ok {
			request.IsHaEnable = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyNamespace(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf namespace failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudTsfNamespaceRead(d, meta)
}

func resourceTencentCloudTsfNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_namespace.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	namespaceId := d.Id()

	if err := service.DeleteTsfNamespaceById(ctx, namespaceId); err != nil {
		return err
	}

	return nil
}
