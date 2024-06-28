package emr

import (
	"context"
	innerErr "errors"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svccdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cdb"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
)

func ResourceTencentCloudEmrCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudEmrClusterCreate,
		Read:   resourceTencentCloudEmrClusterRead,
		Delete: resourceTencentCloudEmrClusterDelete,
		Update: resourceTencentCloudEmrClusterUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"display_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Deprecated:  "It will be deprecated in later versions.",
				Description: "Display strategy of EMR instance.",
			},
			"product_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				Description: "Product ID. Different products ID represents different EMR product versions. Value range:\n" +
					"- 16: represents EMR-V2.3.0\n" +
					"- 20: indicates EMR-V2.5.0\n" +
					"- 25: represents EMR-V3.1.0\n" +
					"- 27: represents KAFKA-V1.0.0\n" +
					"- 30: indicates EMR-V2.6.0\n" +
					"- 33: represents EMR-V3.2.1\n" +
					"- 34: stands for EMR-V3.3.0\n" +
					"- 36: represents STARROCKS-V1.0.0\n" +
					"- 37: indicates EMR-V3.4.0\n" +
					"- 38: represents EMR-V2.7.0\n" +
					"- 39: stands for STARROCKS-V1.1.0\n" +
					"- 41: represents DRUID-V1.1.0.",
			},
			"vpc_settings": {
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    true,
				Description: "The private net config of EMR instance.",
			},
			"softwares": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The softwares of a EMR instance.",
			},
			"resource_spec": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"master_resource_spec": buildResourceSpecSchema(),
						"core_resource_spec":   buildResourceSpecSchema(),
						"task_resource_spec":   buildResourceSpecSchema(),
						"master_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The number of master node.",
						},
						"core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The number of core node.",
						},
						"task_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The number of core node.",
						},
						"common_resource_spec": buildResourceSpecSchema(),
						"common_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							ForceNew:    true,
							Description: "The number of common node.",
						},
					},
				},
				Description: "Resource specification of EMR instance.",
			},
			"support_ha": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
				Description:  "The flag whether the instance support high availability.(0=>not support, 1=>support).",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateStringLengthInRange(6, 36),
				Description:  "Name of the instance, which can contain 6 to 36 English letters, Chinese characters, digits, dashes(-), or underscores(_).",
			},
			"pay_mode": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 1),
				Description:  "The pay mode of instance. 0 represent POSTPAID_BY_HOUR, 1 represent PREPAID.",
			},
			"placement": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"placement", "placement_info"},
				Deprecated:   "It will be deprecated in later versions. Use `placement_info` instead.",
				Description:  "The location of the instance.",
			},
			"placement_info": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"placement", "placement_info"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Zone.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "Project id.",
						},
					},
				},
				Description: "The location of the instance.",
			},
			"time_span": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The length of time the instance was purchased. Use with TimeUnit.When TimeUnit is s, the parameter can only be filled in at 3600, representing a metered instance.\nWhen TimeUnit is m, the number filled in by this parameter indicates the length of purchase of the monthly instance of the package year, such as 1 for one month of purchase.",
			},
			"time_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unit of time in which the instance was purchased. When PayMode is 0, TimeUnit can only take values of s(second). When PayMode is 1, TimeUnit can only take the value m(month).",
			},
			"login_settings": {
				Type:      schema.TypeMap,
				Optional:  true,
				Sensitive: true,
				Description: "Instance login settings. There are two optional fields:" +
					"- password: Instance login password: 8-16 characters, including uppercase letters, lowercase letters, numbers and special characters. Special symbols only support! @% ^ *. The first bit of the password cannot be a special character;" +
					"- public_key_id: Public key id. After the key is associated, the instance can be accessed through the corresponding private key.",
			},
			"extend_fs_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access the external file system.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Created EMR instance id.",
			},
			"need_master_wan": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      EMR_MASTER_WAN_TYPE_NEED_MASTER_WAN,
				ValidateFunc: tccommon.ValidateAllowedStringValue(EMR_MASTER_WAN_TYPES),
				Description: `Whether to enable the cluster Master node public network. Value range:
				- NEED_MASTER_WAN: Indicates that the cluster Master node public network is enabled.
				- NOT_NEED_MASTER_WAN: Indicates that it is not turned on.
				By default, the cluster Master node internet is enabled.`,
			},
			"sg_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the security group to which the instance belongs, in the form of sg-xxxxxxxx.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudEmrClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster.update")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	immutableFields := []string{"placement", "placement_info", "display_strategy", "login_settings", "resource_spec.0.master_count", "resource_spec.0.task_count", "resource_spec.0.core_count"}
	for _, f := range immutableFields {
		if d.HasChange(f) {
			return fmt.Errorf("cannot update argument `%s`", f)
		}
	}

	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	instanceId := d.Id()
	timeUnit, hasTimeUnit := d.GetOkExists("time_unit")
	timeSpan, hasTimeSpan := d.GetOkExists("time_span")
	payMode, hasPayMode := d.GetOkExists("pay_mode")
	if !hasTimeUnit || !hasTimeSpan || !hasPayMode {
		return innerErr.New("Time_unit, time_span or pay_mode must be set.")
	}
	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		err := emrService.ModifyResourcesTags(ctx, meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region, instanceId, oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		if err != nil {
			return err
		}
	}

	hasChange := false
	request := emr.NewScaleOutInstanceRequest()
	request.TimeUnit = common.StringPtr(timeUnit.(string))
	request.TimeSpan = common.Uint64Ptr((uint64)(timeSpan.(int)))
	request.PayMode = common.Uint64Ptr((uint64)(payMode.(int)))
	request.InstanceId = common.StringPtr(instanceId)

	tmpResourceSpec := d.Get("resource_spec").([]interface{})
	resourceSpec := tmpResourceSpec[0].(map[string]interface{})

	if d.HasChange("resource_spec.0.master_count") {
		request.MasterCount = common.Uint64Ptr((uint64)(resourceSpec["master_count"].(int)))
		hasChange = true
	}
	if d.HasChange("resource_spec.0.task_count") {
		request.TaskCount = common.Uint64Ptr((uint64)(resourceSpec["task_count"].(int)))
		hasChange = true
	}
	if d.HasChange("resource_spec.0.core_count") {
		request.CoreCount = common.Uint64Ptr((uint64)(resourceSpec["core_count"].(int)))
		hasChange = true
	}
	if d.HasChange("extend_fs_field") {
		return innerErr.New("extend_fs_field not support update.")
	}
	if !hasChange {
		return nil
	}
	_, err := emrService.UpdateInstance(ctx, request)
	if err != nil {
		return err
	}
	err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}

		if len(clusters) > 0 {
			status := *(clusters[0].Status)
			if status != EmrInternetStatusCreated {
				return resource.RetryableError(
					fmt.Errorf("%v create cluster endpoint  status still is %v", instanceId, status))
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return resourceTencentCloudEmrClusterRead(d, meta)
}

func resourceTencentCloudEmrClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster.create")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	immutableFields := []string{"time_unit", "time_span", "login_settings"}
	for _, f := range immutableFields {
		if _, ok := d.GetOkExists(f); !ok {
			return fmt.Errorf("Argument `%s` must be set", f)
		}
	}

	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	instanceId, err := emrService.CreateInstance(ctx, d)
	if err != nil {
		return err
	}
	d.SetId(instanceId)
	_ = d.Set("instance_id", instanceId)
	var displayStrategy string
	if v, ok := d.GetOk("display_strategy"); ok {
		displayStrategy = v.(string)
	} else {
		displayStrategy = "clusterList"
	}
	err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstancesById(ctx, instanceId, displayStrategy)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}

		if len(clusters) > 0 {
			status := *(clusters[0].Status)
			if status != EmrInternetStatusCreated {
				return resource.RetryableError(
					fmt.Errorf("%v create cluster endpoint  status still is %v", instanceId, status))
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudEmrClusterRead(d, meta)
}

func resourceTencentCloudEmrClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_emr_cluster.delete")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	instanceId := d.Id()
	clusters, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)
	if len(clusters) == 0 {
		return innerErr.New("Not find clusters.")
	}
	metaDB := clusters[0].MetaDb
	if err != nil {
		return err
	}
	if err = emrService.DeleteInstance(ctx, d); err != nil {
		return err
	}
	err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		clusters, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
			if e.GetCode() == "UnauthorizedOperation" {
				return nil
			}
		}

		if len(clusters) > 0 {
			status := *(clusters[0].Status)
			if status != EmrInternetStatusDeleted {
				return resource.RetryableError(
					fmt.Errorf("%v create cluster endpoint  status still is %v", instanceId, status))
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if metaDB != nil && *metaDB != "" {
		// remove metadb
		mysqlService := svccdb.NewMysqlService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err := mysqlService.OfflineIsolatedInstances(ctx, *metaDB)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func resourceTencentCloudEmrClusterRead(d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	emrService := EMRService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	instanceId := d.Id()
	var instance *emr.ClusterInstancesInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, err := emrService.DescribeInstancesById(ctx, instanceId, DisplayStrategyIsclusterList)

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
				return nil
			}
		}

		if err != nil {
			return resource.RetryableError(err)
		}

		if len(result) > 0 {
			instance = result[0]
		}

		return nil
	})

	if err != nil {
		return err
	}

	_ = d.Set("instance_id", instanceId)
	if instance != nil {
		_ = d.Set("product_id", instance.ProductId)
		_ = d.Set("vpc_settings", map[string]interface{}{
			"vpc_id":    *instance.UniqVpcId,
			"subnet_id": *instance.UniqSubnetId,
		})
		if instance.Config != nil {
			if instance.Config.SoftInfo != nil {
				_ = d.Set("softwares", helper.PStrings(instance.Config.SoftInfo))
			}

			if instance.Config.SupportHA != nil {
				if *instance.Config.SupportHA {
					_ = d.Set("support_ha", 1)
				} else {
					_ = d.Set("support_ha", 0)
				}
			}

			if instance.Config.SecurityGroup != nil {
				_ = d.Set("sg_id", instance.Config.SecurityGroup)
			}
			resourceSpec := make(map[string]interface{})

			var masterCount int64
			if instance.Config.MasterNodeSize != nil {
				masterCount = *instance.Config.MasterNodeSize
				resourceSpec["master_count"] = masterCount
			}
			if masterCount != 0 && instance.Config.MasterResource != nil {
				masterResource := instance.Config.MasterResource
				masterResourceSpec := make(map[string]interface{})
				if masterResource.MemSize != nil {
					masterResourceSpec["mem_size"] = *masterResource.MemSize
				}
				if masterResource.Cpu != nil {
					masterResourceSpec["cpu"] = *masterResource.Cpu
				}
				if masterResource.DiskSize != nil {
					masterResourceSpec["disk_size"] = *masterResource.DiskSize
				}
				if masterResource.DiskType != nil {
					masterResourceSpec["disk_type"] = *masterResource.DiskType
				}
				if masterResource.Spec != nil {
					masterResourceSpec["spec"] = *masterResource.Spec
				}
				if masterResource.StorageType != nil {
					masterResourceSpec["storage_type"] = *masterResource.StorageType
				}
				if masterResource.RootSize != nil {
					masterResourceSpec["root_size"] = *masterResource.RootSize
				}
				resourceSpec["master_resource_spec"] = []interface{}{masterResourceSpec}
			}

			var coreCount int64
			if instance.Config.CoreNodeSize != nil {
				coreCount = *instance.Config.CoreNodeSize
				resourceSpec["core_count"] = coreCount
			}
			if coreCount != 0 && instance.Config.CoreResource != nil {
				coreResource := instance.Config.CoreResource
				coreResourceSpec := make(map[string]interface{})
				if coreResource.MemSize != nil {
					coreResourceSpec["mem_size"] = *coreResource.MemSize
				}
				if coreResource.Cpu != nil {
					coreResourceSpec["cpu"] = *coreResource.Cpu
				}
				if coreResource.DiskSize != nil {
					coreResourceSpec["disk_size"] = *coreResource.DiskSize
				}
				if coreResource.DiskType != nil {
					coreResourceSpec["disk_type"] = *coreResource.DiskType
				}
				if coreResource.Spec != nil {
					coreResourceSpec["spec"] = *coreResource.Spec
				}
				if coreResource.StorageType != nil {
					coreResourceSpec["storage_type"] = *coreResource.StorageType
				}
				if coreResource.RootSize != nil {
					coreResourceSpec["root_size"] = *coreResource.RootSize
				}
				resourceSpec["core_resource_spec"] = []interface{}{coreResourceSpec}
			}

			var taskCount int64
			if instance.Config.TaskNodeSize != nil {
				taskCount = *instance.Config.TaskNodeSize
				resourceSpec["task_count"] = taskCount
			}
			if taskCount != 0 && instance.Config.TaskResource != nil {
				taskResource := instance.Config.TaskResource
				taskResourceSpec := make(map[string]interface{})
				if taskResource.MemSize != nil {
					taskResourceSpec["mem_size"] = *taskResource.MemSize
				}
				if taskResource.Cpu != nil {
					taskResourceSpec["cpu"] = *taskResource.Cpu
				}
				if taskResource.DiskSize != nil {
					taskResourceSpec["disk_size"] = *taskResource.DiskSize
				}
				if taskResource.DiskType != nil {
					taskResourceSpec["disk_type"] = *taskResource.DiskType
				}
				if taskResource.Spec != nil {
					taskResourceSpec["spec"] = *taskResource.Spec
				}
				if taskResource.StorageType != nil {
					taskResourceSpec["storage_type"] = *taskResource.StorageType
				}
				if taskResource.RootSize != nil {
					taskResourceSpec["root_size"] = *taskResource.RootSize
				}
				resourceSpec["task_resource_spec"] = []interface{}{taskResourceSpec}
			}

			var commonCount int64
			if instance.Config.ComNodeSize != nil {
				commonCount = *instance.Config.ComNodeSize
				resourceSpec["common_count"] = commonCount
			}
			if commonCount != 0 && instance.Config.ComResource != nil {
				comResource := instance.Config.ComResource
				comResourceSpec := make(map[string]interface{})
				if comResource.MemSize != nil {
					comResourceSpec["mem_size"] = *comResource.MemSize
				}
				if comResource.Cpu != nil {
					comResourceSpec["cpu"] = *comResource.Cpu
				}
				if comResource.DiskSize != nil {
					comResourceSpec["disk_size"] = *comResource.DiskSize
				}
				if comResource.DiskType != nil {
					comResourceSpec["disk_type"] = *comResource.DiskType
				}
				if comResource.Spec != nil {
					comResourceSpec["spec"] = *comResource.Spec
				}
				if comResource.StorageType != nil {
					comResourceSpec["storage_type"] = *comResource.StorageType
				}
				if comResource.RootSize != nil {
					comResourceSpec["root_size"] = *comResource.RootSize
				}
				resourceSpec["common_resource_spec"] = []interface{}{comResourceSpec}
			}

			_ = d.Set("resource_spec", []interface{}{resourceSpec})
		}

		_ = d.Set("instance_name", instance.ClusterName)
		_ = d.Set("pay_mode", instance.ChargeType)
		placement := map[string]interface{}{
			"zone":       *instance.Zone,
			"project_id": *instance.ProjectId,
		}
		_ = d.Set("placement", map[string]interface{}{
			"zone": *instance.Zone,
		})
		_ = d.Set("placement_info", []interface{}{placement})
		if instance.MasterIp != nil {
			_ = d.Set("need_master_wan", "NEED_MASTER_WAN")
		} else {
			_ = d.Set("need_master_wan", "NOT_NEED_MASTER_WAN")
		}
	}

	tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
	tags, err := tagService.DescribeResourceTags(ctx, "emr", "emr-instance", region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)
	return nil
}
