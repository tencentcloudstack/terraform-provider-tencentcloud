package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTsfCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfClusterCreate,
		Read:   resourceTencentCloudTsfClusterRead,
		Update: resourceTencentCloudTsfClusterUpdate,
		Delete: resourceTencentCloudTsfClusterDelete,
		// Importer: &schema.ResourceImporter{
		// 	State: schema.ImportStatePassthrough,
		// },
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"cluster_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster name.",
			},

			"cluster_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster type.",
			},

			"vpc_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Vpc id.",
			},

			"cluster_cidr": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "CIDR assigned to cluster containers and service IP.",
			},

			"cluster_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "cluster notes.",
			},

			"tsf_region_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The TSF region to which the cluster belongs.",
			},

			"tsf_zone_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The TSF availability zone to which the cluster belongs.",
			},

			"subnet_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Subnet id.",
			},

			"cluster_version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "cluster version.",
			},

			"max_node_pod_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of Pods on each Node in the cluster. The value ranges from 4 to 256. When the value is not a power of 2, the nearest power of 2 will be taken up.",
			},

			"max_cluster_service_num": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The maximum number of services in the cluster. The value ranges from 32 to 32768. If it is not a power of 2, the nearest power of 2 will be taken up.",
			},

			"program_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The dataset ID to be bound.",
			},

			"kubernete_api_server": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "api address.",
			},

			"kubernete_native_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "`K`:kubeconfig, `S`:service account.",
			},

			"kubernete_native_secret": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "native secret.",
			},

			"program_id_list": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Program id list.",
			},

			"cluster_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "cluster status.",
			},

			"cluster_total_cpu": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "The total CPU of the cluster, unit: core.",
			},

			"cluster_total_mem": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "The total memory of the cluster, unit: G.",
			},

			"cluster_used_cpu": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "CPU used by the cluster, unit: core.",
			},

			"cluster_used_mem": {
				Computed:    true,
				Type:        schema.TypeFloat,
				Description: "The memory used by the cluster, unit: G.",
			},

			"instance_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of cluster machine instances.",
			},

			"run_instance_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of machine instances running in the cluster.",
			},

			"normal_instance_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The number of machine instances in the normal state of the cluster.",
			},

			"delete_flag": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Delete flag: `true`: can be deleted; `false`: can not be deleted.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Create time.",
			},

			"update_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Update time.",
			},

			"tsf_region_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Name of the TSF region to which the cluster belongs.",
			},

			"tsf_zone_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The name of the TSF availability zone to which the cluster belongs.",
			},

			"delete_flag_reason": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Reasons why clusters cannot be deleted.",
			},

			"cluster_limit_cpu": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cluster remaining cpu limit.",
			},

			"cluster_limit_mem": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Cluster remaining memory limit.",
			},

			"run_service_instance_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of running service instances.",
			},

			"operation_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Control information for buttons on the front end.",
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
										Description: "Reason for not showing.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the button clickable.",
									},
									"supported": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether to show the button.",
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
										Description: "Reason for not showing.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the button clickable.",
									},
									"supported": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether to show the button.",
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
										Description: "Reason for not showing.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the button clickable.",
									},
									"supported": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "whether to show the button.",
									},
								},
							},
						},
					},
				},
			},

			"group_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total number of deployment groups.",
			},

			"run_group_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of Deployment Groups in progress.",
			},

			"stop_group_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Number of deployment groups in stop.",
			},

			"abnormal_group_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Abnormal number of deployment groups.",
			},

			"cluster_remark_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "cluster remark name.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTsfClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_cluster.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = tsf.NewCreateClusterRequest()
		response  = tsf.NewCreateClusterResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_name"); ok {
		request.ClusterName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_type"); ok {
		request.ClusterType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_cidr"); ok {
		request.ClusterCIDR = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_desc"); ok {
		request.ClusterDesc = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tsf_region_id"); ok {
		request.TsfRegionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tsf_zone_id"); ok {
		request.TsfZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cluster_version"); ok {
		request.ClusterVersion = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("max_node_pod_num"); ok {
		request.MaxNodePodNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_cluster_service_num"); ok {
		request.MaxClusterServiceNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("program_id"); ok {
		request.ProgramId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kubernete_api_server"); ok {
		request.KuberneteApiServer = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kubernete_native_type"); ok {
		request.KuberneteNativeType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kubernete_native_secret"); ok {
		request.KuberneteNativeSecret = helper.String(v.(string))
	}

	if v, ok := d.GetOk("program_id_list"); ok {
		programIdListSet := v.(*schema.Set).List()
		for i := range programIdListSet {
			programIdList := programIdListSet[i].(string)
			request.ProgramIdList = append(request.ProgramIdList, &programIdList)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().CreateCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf cluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.Result
	d.SetId(clusterId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	conf := BuildStateChangeConf(
		[]string{"Creating"},
		[]string{"Running"},
		8*readRetryTimeout,
		time.Second,
		service.TsfClusterStateRefreshFunc(d.Id(), []string{}),
	)
	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:cluster/%s", region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfClusterRead(d, meta)
}

func resourceTencentCloudTsfClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}

	clusterId := d.Id()

	cluster, err := service.DescribeTsfClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	if cluster == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfCluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if cluster.ClusterId != nil {
		_ = d.Set("cluster_id", cluster.ClusterId)
	}

	if cluster.ClusterName != nil {
		_ = d.Set("cluster_name", cluster.ClusterName)
	}

	if cluster.ClusterType != nil {
		_ = d.Set("cluster_type", cluster.ClusterType)
	}

	if cluster.VpcId != nil {
		_ = d.Set("vpc_id", cluster.VpcId)
	}

	if cluster.ClusterCIDR != nil {
		_ = d.Set("cluster_cidr", cluster.ClusterCIDR)
	}

	if cluster.ClusterDesc != nil {
		_ = d.Set("cluster_desc", cluster.ClusterDesc)
	}

	if cluster.TsfRegionId != nil {
		_ = d.Set("tsf_region_id", cluster.TsfRegionId)
	}

	if cluster.TsfZoneId != nil {
		_ = d.Set("tsf_zone_id", cluster.TsfZoneId)
	}

	if cluster.SubnetId != nil {
		_ = d.Set("subnet_id", cluster.SubnetId)
	}

	if cluster.ClusterVersion != nil {
		_ = d.Set("cluster_version", cluster.ClusterVersion)
	}

	// if cluster.MaxNodePodNum != nil {
	// 	_ = d.Set("max_node_pod_num", cluster.MaxNodePodNum)
	// }

	// if cluster.MaxClusterServiceNum != nil {
	// 	_ = d.Set("max_cluster_service_num", cluster.MaxClusterServiceNum)
	// }

	// if cluster.ProgramId != nil {
	// 	_ = d.Set("program_id", cluster.ProgramId)
	// }

	if cluster.KuberneteApiServer != nil {
		_ = d.Set("kubernete_api_server", cluster.KuberneteApiServer)
	}

	if cluster.KuberneteNativeType != nil {
		_ = d.Set("kubernete_native_type", cluster.KuberneteNativeType)
	}

	if cluster.KuberneteNativeSecret != nil {
		_ = d.Set("kubernete_native_secret", cluster.KuberneteNativeSecret)
	}

	// if cluster.ProgramIdList != nil {
	// 	_ = d.Set("program_id_list", cluster.ProgramIdList)
	// }

	if cluster.ClusterStatus != nil {
		_ = d.Set("cluster_status", cluster.ClusterStatus)
	}

	if cluster.ClusterTotalCpu != nil {
		_ = d.Set("cluster_total_cpu", cluster.ClusterTotalCpu)
	}

	if cluster.ClusterTotalMem != nil {
		_ = d.Set("cluster_total_mem", cluster.ClusterTotalMem)
	}

	if cluster.ClusterUsedCpu != nil {
		_ = d.Set("cluster_used_cpu", cluster.ClusterUsedCpu)
	}

	if cluster.ClusterUsedMem != nil {
		_ = d.Set("cluster_used_mem", cluster.ClusterUsedMem)
	}

	if cluster.InstanceCount != nil {
		_ = d.Set("instance_count", cluster.InstanceCount)
	}

	if cluster.RunInstanceCount != nil {
		_ = d.Set("run_instance_count", cluster.RunInstanceCount)
	}

	if cluster.NormalInstanceCount != nil {
		_ = d.Set("normal_instance_count", cluster.NormalInstanceCount)
	}

	if cluster.DeleteFlag != nil {
		_ = d.Set("delete_flag", cluster.DeleteFlag)
	}

	if cluster.CreateTime != nil {
		_ = d.Set("create_time", cluster.CreateTime)
	}

	if cluster.UpdateTime != nil {
		_ = d.Set("update_time", cluster.UpdateTime)
	}

	if cluster.TsfRegionName != nil {
		_ = d.Set("tsf_region_name", cluster.TsfRegionName)
	}

	if cluster.TsfZoneName != nil {
		_ = d.Set("tsf_zone_name", cluster.TsfZoneName)
	}

	if cluster.DeleteFlagReason != nil {
		_ = d.Set("delete_flag_reason", cluster.DeleteFlagReason)
	}

	if cluster.ClusterLimitCpu != nil {
		_ = d.Set("cluster_limit_cpu", cluster.ClusterLimitCpu)
	}

	if cluster.ClusterLimitMem != nil {
		_ = d.Set("cluster_limit_mem", cluster.ClusterLimitMem)
	}

	if cluster.RunServiceInstanceCount != nil {
		_ = d.Set("run_service_instance_count", cluster.RunServiceInstanceCount)
	}

	if cluster.OperationInfo != nil {
		operationInfoMap := map[string]interface{}{}

		if cluster.OperationInfo.Init != nil {
			initMap := map[string]interface{}{}

			if cluster.OperationInfo.Init.DisabledReason != nil {
				initMap["disabled_reason"] = cluster.OperationInfo.Init.DisabledReason
			}

			if cluster.OperationInfo.Init.Enabled != nil {
				initMap["enabled"] = cluster.OperationInfo.Init.Enabled
			}

			if cluster.OperationInfo.Init.Supported != nil {
				initMap["supported"] = cluster.OperationInfo.Init.Supported
			}

			operationInfoMap["init"] = []interface{}{initMap}
		}

		if cluster.OperationInfo.AddInstance != nil {
			addInstanceMap := map[string]interface{}{}

			if cluster.OperationInfo.AddInstance.DisabledReason != nil {
				addInstanceMap["disabled_reason"] = cluster.OperationInfo.AddInstance.DisabledReason
			}

			if cluster.OperationInfo.AddInstance.Enabled != nil {
				addInstanceMap["enabled"] = cluster.OperationInfo.AddInstance.Enabled
			}

			if cluster.OperationInfo.AddInstance.Supported != nil {
				addInstanceMap["supported"] = cluster.OperationInfo.AddInstance.Supported
			}

			operationInfoMap["add_instance"] = []interface{}{addInstanceMap}
		}

		if cluster.OperationInfo.Destroy != nil {
			destroyMap := map[string]interface{}{}

			if cluster.OperationInfo.Destroy.DisabledReason != nil {
				destroyMap["disabled_reason"] = cluster.OperationInfo.Destroy.DisabledReason
			}

			if cluster.OperationInfo.Destroy.Enabled != nil {
				destroyMap["enabled"] = cluster.OperationInfo.Destroy.Enabled
			}

			if cluster.OperationInfo.Destroy.Supported != nil {
				destroyMap["supported"] = cluster.OperationInfo.Destroy.Supported
			}

			operationInfoMap["destroy"] = []interface{}{destroyMap}
		}

		_ = d.Set("operation_info", []interface{}{operationInfoMap})
	}

	if cluster.GroupCount != nil {
		_ = d.Set("group_count", cluster.GroupCount)
	}

	if cluster.RunGroupCount != nil {
		_ = d.Set("run_group_count", cluster.RunGroupCount)
	}

	if cluster.StopGroupCount != nil {
		_ = d.Set("stop_group_count", cluster.StopGroupCount)
	}

	if cluster.AbnormalGroupCount != nil {
		_ = d.Set("abnormal_group_count", cluster.AbnormalGroupCount)
	}

	if cluster.ClusterRemarkName != nil {
		_ = d.Set("cluster_remark_name", cluster.ClusterRemarkName)
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "cluster", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_cluster.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tsf.NewModifyClusterRequest()

	clusterId := d.Id()

	request.ClusterId = &clusterId

	immutableArgs := []string{"cluster_id", "cluster_type", "vpc_id", "cluster_cidr", "tsf_region_id", "tsf_zone_id", "subnet_id", "cluster_version", "max_node_pod_num", "max_cluster_service_num", "program_id", "kubernete_api_server", "kubernete_native_type", "kubernete_native_secret", "program_id_list", "cluster_status", "cluster_total_cpu", "cluster_total_mem", "cluster_used_cpu", "cluster_used_mem", "instance_count", "run_instance_count", "normal_instance_count", "delete_flag", "create_time", "update_time", "tsf_region_name", "tsf_zone_name", "delete_flag_reason", "cluster_limit_cpu", "cluster_limit_mem", "run_service_instance_count", "operation_info", "group_count", "run_group_count", "stop_group_count", "abnormal_group_count"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("cluster_name") {
		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}
	}

	if d.HasChange("cluster_desc") {
		if v, ok := d.GetOk("cluster_desc"); ok {
			request.ClusterDesc = helper.String(v.(string))
		}
	}

	if d.HasChange("cluster_remark_name") {
		if v, ok := d.GetOk("cluster_remark_name"); ok {
			request.ClusterRemarkName = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTsfClient().ModifyCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf cluster failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := BuildTagResourceName("tsf", "cluster", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfClusterRead(d, meta)
}

func resourceTencentCloudTsfClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tsf_cluster.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TsfService{client: meta.(*TencentCloudClient).apiV3Conn}
	clusterId := d.Id()

	if err := service.DeleteTsfClusterById(ctx, clusterId); err != nil {
		return err
	}

	return nil
}
