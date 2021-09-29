/*
Provide a resource to create an auto scaling group for kubernetes cluster.

~> **NOTE:**  We recommend the usage of one cluster with essential worker config + node pool to manage cluster and nodes. Its a more flexible way than manage worker config with tencentcloud_kubernetes_cluster, tencentcloud_kubernetes_scale_worker or exist node management of `tencentcloud_kubernetes_attachment`. Cause some unchangeable parameters of `worker_config` may cause the whole cluster resource `force new`.

Example Usage

```hcl

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "cluster_cidr" {
  default = "172.31.0.0/16"
}

data "tencentcloud_vpc_subnets" "vpc" {
    is_default        = true
    availability_zone = var.availability_zone
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

//this is the cluster with empty worker config
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  cluster_cidr            = var.cluster_cidr
  cluster_max_pod_num     = 32
  cluster_name            = "tf-tke-unit-test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32
  cluster_version         = "1.18.4"
  cluster_deploy_type = "MANAGED_CLUSTER"
}

//this is one example of managing node using node pool
resource "tencentcloud_kubernetes_node_pool" "mynodepool" {
  name = "mynodepool"
  cluster_id = tencentcloud_kubernetes_cluster.managed_cluster.id
  max_size = 6
  min_size = 1
  vpc_id               = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_ids           = [data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id]
  retry_policy         = "INCREMENTAL_INTERVALS"
  desired_capacity     = 4
  enable_auto_scale    = true

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    security_group_ids = ["sg-24vswocp"]

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 10
    public_ip_assigned         = true
    password                   = "test123#"
    enhanced_security_service  = false
    enhanced_monitor_service   = false

  }

  labels = {
    "test1" = "test1",
    "test2" = "test2",
  }

  taints {
	key = "test_taint"
    value = "taint_value"
    effect = "PreferNoSchedule"
  }

  taints {
	key = "test_taint2"
    value = "taint_value2"
    effect = "PreferNoSchedule"
  }

  node_config {
      extra_args = [
 	"root-dir=/var/lib/kubelet"
  ]
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func getNodePoolInstanceTypes(d *schema.ResourceData) []*string {
	configParas := d.Get("auto_scaling_config").([]interface{})
	dMap := configParas[0].(map[string]interface{})
	instanceType, _ := dMap["instance_type"]
	currInsType := instanceType.(string)
	v, ok := dMap["backup_instance_types"]
	backupInstanceTypes := v.([]interface{})
	instanceTypes := make([]*string, 0)
	if !ok || len(backupInstanceTypes) == 0 {
		instanceTypes = append(instanceTypes, &currInsType)
		return instanceTypes
	}
	headType := backupInstanceTypes[0].(string)
	if headType != currInsType {
		instanceTypes = append(instanceTypes, &currInsType)
	}
	for i := range backupInstanceTypes {
		insType := backupInstanceTypes[i].(string)
		instanceTypes = append(instanceTypes, &insType)
	}

	return instanceTypes
}

func composedKubernetesAsScalingConfigPara() map[string]*schema.Schema {
	needSchema := map[string]*schema.Schema{
		"instance_type": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Specified types of CVM instance.",
		},
		"backup_instance_types": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Backup CVM instance types if specified instance type sold out or mismatch.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"system_disk_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
			ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
			Description:  "Type of a CVM disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`. Default is `CLOUD_PREMIUM`.",
		},
		"system_disk_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      50,
			ValidateFunc: validateIntegerInRange(50, 500),
			Description:  "Volume of system disk in GB. Default is `50`.",
		},
		"data_disk": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    11,
			Description: "Configurations of data disk.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
						ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
						Description:  "Types of disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`.",
					},
					"disk_size": {
						Type:        schema.TypeInt,
						Optional:    true,
						ForceNew:    true,
						Default:     0,
						Description: "Volume of disk in GB. Default is `0`.",
					},
					"snapshot_id": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						Description: "Data disk snapshot ID.",
					},
				},
			},
		},
		"internet_charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID_BY_HOUR,
			ValidateFunc: validateAllowedStringValue(INTERNET_CHARGE_ALLOW_TYPE),
			Description:  "Charge types for network traffic. Valid value: `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
		},
		"internet_max_bandwidth_out": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Max bandwidth of Internet access in Mbps. Default is `0`.",
		},
		"bandwidth_package_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.",
		},
		"public_ip_assigned": {
			Type:        schema.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Description: "Specify whether to assign an Internet IP address.",
		},
		"password": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			ForceNew:      true,
			ValidateFunc:  validateAsConfigPassword,
			ConflictsWith: []string{"auto_scaling_config.0.key_ids"},
			Description:   "Password to access.",
		},
		"key_ids": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      true,
			Elem:          &schema.Schema{Type: schema.TypeString},
			ConflictsWith: []string{"auto_scaling_config.0.password"},
			Description:   "ID list of keys.",
		},
		"security_group_ids": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Security groups to which a CVM instance belongs.",
		},
		"enhanced_security_service": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
			Description: "To specify whether to enable cloud security service. Default is TRUE.",
		},
		"enhanced_monitor_service": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			ForceNew:    true,
			Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
		},
	}

	return needSchema
}

func ResourceTencentCloudKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesNodePoolCreate,
		Read:   resourceKubernetesNodePoolRead,
		Delete: resourceKubernetesNodePoolDelete,
		Update: resourceKubernetesNodePoolUpdate,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "ID of the cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the node pool. The name does not exceed 25 characters, and only supports Chinese, English, numbers, underscores, separators (`-`) and decimal points.",
			},
			"max_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 2000),
				Description:  "Maximum number of node.",
			},
			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateIntegerInRange(0, 2000),
				Description:  "Minimum number of node.",
			},
			"desired_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateIntegerInRange(0, 2000),
				Description:  "Desired capacity ot the node. If `enable_auto_scale` is set `true`, this will be a computed parameter.",
			},
			"enable_auto_scale": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate whether to enable auto scaling or not.",
			},
			"retry_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Available values for retry policies include `IMMEDIATE_RETRY` and `INCREMENTAL_INTERVALS`.",
				Default:     SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
				ValidateFunc: validateAllowedStringValue([]string{SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY,
					SCALING_GROUP_RETRY_POLICY_INCREMENTAL_INTERVALS}),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of VPC network.",
			},
			"subnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "ID list of subnet, and for VPC it is required.",
			},
			"scaling_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Auto scaling mode. Valid values are `CLASSIC_SCALING`(scaling by create/destroy instances), " +
					"`WAKE_UP_STOPPED_SCALING`(Boot priority for expansion. When expanding the capacity, the shutdown operation is given priority to the shutdown of the instance." +
					" If the number of instances is still lower than the expected number of instances after the startup, the instance will be created, and the method of destroying the instance will still be used for shrinking)" +
					".",
			},
			"node_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: TkeInstanceAdvancedSetting(),
				},
				Description: "Node config.",
			},
			"auto_scaling_config": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: composedKubernetesAsScalingConfigPara(),
				},
				Description: "Auto scaling config parameters.",
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Labels of kubernetes node pool created nodes. The label key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').",
			},
			"unschedulable": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.",
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Key of the taint. The taint key name does not exceed 63 characters, only supports English, numbers,'/','-', and does not allow beginning with ('/').",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the taint.",
						},
						"effect": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Effect of the taint. Valid values are: `NoSchedule`, `PreferNoSchedule`, `NoExecute`.",
						},
					},
				},
				Description: "Taints of kubernetes node pool created nodes.",
			},
			"delete_keep_instance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate to keep the CVM instance when delete the node pool. Default is `true`.",
			},
			"node_os": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tlinux2.4x86_64",
				Description: "Operating system of the cluster, the available values include: `tlinux2.4x86_64`, `ubuntu18.04.1x86_64`, `ubuntu16.04.1 LTSx86_64`, `centos7.6.0_x64` and `centos7.2x86_64`. Default is 'tlinux2.4x86_64'. This parameter will only affect new nodes, not including the existing nodes.",
			},
			"node_os_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "GENERAL",
				Description: "The image version of the node. Valida values are `DOCKER_CUSTOMIZE` and `GENERAL`. Default is `GENERAL`. This parameter will only affect new nodes, not including the existing nodes.",
			},
			// asg pass through arguments
			"scaling_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of relative scaling group.",
			},
			"zones": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of auto scaling group available zones, for Basic network it is required.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"scaling_group_project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Project ID the scaling group belongs to.",
			},
			"default_cooldown": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Seconds of scaling group cool down. Default value is `300`.",
			},
			"termination_policies": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Policy of scaling group termination. Available values: `[\"OLDEST_INSTANCE\"]`, `[\"NEWEST_INSTANCE\"]`.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			//computed
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the node pool.",
			},
			"node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total node count.",
			},
			"launch_config_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The launch config ID.",
			},
			"auto_scaling_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The auto scaling group ID.",
			},
		},
		//compare to console, miss cam_role and running_version and lock_initial_node and security_proof
	}
}

//this function composes every single parameter to a as scale parameter with json string format
func composeParameterToAsScalingGroupParaSerial(d *schema.ResourceData) (string, error) {
	var (
		result string
		errRet error
	)

	request := as.NewCreateAutoScalingGroupRequest()

	//this is an empty string
	request.MaxSize = helper.IntUint64(d.Get("max_size").(int))
	request.MinSize = helper.IntUint64(d.Get("min_size").(int))

	if *request.MinSize > *request.MaxSize {
		return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
	}

	request.VpcId = helper.String(d.Get("vpc_id").(string))

	if v, ok := d.GetOk("desired_capacity"); ok {
		request.DesiredCapacity = helper.IntUint64(v.(int))
		if *request.DesiredCapacity > *request.MaxSize ||
			*request.DesiredCapacity < *request.MinSize {
			return "", fmt.Errorf("constraints `min_size <= desired_capacity <= max_size` must be established,")
		}

	}

	if v, ok := d.GetOk("retry_policy"); ok {
		request.RetryPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_ids"); ok {
		subnetIds := v.([]interface{})
		request.SubnetIds = make([]*string, 0, len(subnetIds))
		for i := range subnetIds {
			subnetId := subnetIds[i].(string)
			request.SubnetIds = append(request.SubnetIds, &subnetId)
		}
	}

	if v, ok := d.GetOk("scaling_mode"); ok {
		request.ServiceSettings = &as.ServiceSettings{ScalingMode: helper.String(v.(string))}
	}

	result = request.ToJsonString()

	return result, errRet
}

//this function is similar to kubernetesAsScalingConfigParaSerial, but less parameter
func composedKubernetesAsScalingConfigParaSerial(dMap map[string]interface{}, meta interface{}) (string, error) {
	var (
		result string
		errRet error
	)

	request := as.NewCreateLaunchConfigurationRequest()

	instanceType := dMap["instance_type"].(string)
	request.InstanceType = &instanceType

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := dMap["system_disk_type"]; ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := dMap["data_disk"]; ok {
		dataDisks := v.([]interface{})
		request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["disk_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			snapshotId := value["snapshot_id"].(string)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
				DiskSize: &diskSize,
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	request.InternetAccessible = &as.InternetAccessible{}
	if v, ok := dMap["internet_charge_type"]; ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := dMap["bandwidth_package_id"]; ok {
		if v.(string) != "" {
			request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
		}
	}
	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	request.LoginSettings = &as.LoginSettings{}

	if v, ok := dMap["password"]; ok {
		request.LoginSettings.Password = helper.String(v.(string))
	}
	if v, ok := dMap["key_ids"]; ok {
		keyIds := v.([]interface{})
		request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
		}
	}

	if request.LoginSettings.Password != nil && *request.LoginSettings.Password == "" {
		request.LoginSettings.Password = nil
	}

	if request.LoginSettings.Password == nil && len(request.LoginSettings.KeyIds) == 0 {
		errRet = fmt.Errorf("Parameters `key_ids` and `password` should be set one")
		return result, errRet
	}

	if request.LoginSettings.Password != nil && len(request.LoginSettings.KeyIds) != 0 {
		errRet = fmt.Errorf("Parameters `key_ids` and `password` can only be supported one")
		return result, errRet
	}

	if v, ok := dMap["security_group_ids"]; ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
		}
	}

	request.EnhancedService = &as.EnhancedService{}

	if v, ok := dMap["enhanced_security_service"]; ok {
		securityService := v.(bool)
		request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := dMap["enhanced_monitor_service"]; ok {
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &as.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}

	chargeType := INSTANCE_CHARGE_TYPE_POSTPAID
	request.InstanceChargeType = &chargeType

	result = request.ToJsonString()
	return result, errRet
}

func resourceKubernetesNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_node_pool.read")()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		asService = AsService{client: meta.(*TencentCloudClient).apiV3Conn}
		items     = strings.Split(d.Id(), FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	_, has, err := service.DescribeCluster(ctx, clusterId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = service.DescribeCluster(ctx, clusterId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if !has {
		d.SetId("")
		return nil
	}

	//Describe Node Pool
	nodePool, has, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, has, err = service.DescribeNodePool(ctx, clusterId, nodePoolId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if !has {
		d.SetId("")
		return nil
	}

	//set not force new parameters
	d.Set("max_size", nodePool.MaxNodesNum)
	d.Set("min_size", nodePool.MinNodesNum)
	d.Set("desired_capacity", nodePool.DesiredNodesNum)
	d.Set("name", nodePool.Name)
	d.Set("status", nodePool.LifeState)
	d.Set("node_count", nodePool.NodeCountSummary)
	d.Set("auto_scaling_group_id", nodePool.AutoscalingGroupId)
	d.Set("launch_config_id", nodePool.LaunchConfigurationId)
	d.Set("enable_auto_scale", *nodePool.AutoscalingGroupStatus == "enabled")
	d.Set("node_os", *nodePool.NodePoolOs)
	d.Set("node_system_type", *nodePool.OsCustomizeType)

	//set composed struct
	lables := make(map[string]interface{}, len(nodePool.Labels))
	for _, v := range nodePool.Labels {
		lables[*v.Name] = *v.Value
	}
	d.Set("labels", lables)

	// Relative scaling group status
	asg, hasAsg, err := asService.DescribeAutoScalingGroupById(ctx, *nodePool.AutoscalingGroupId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			asg, hasAsg, err = asService.DescribeAutoScalingGroupById(ctx, *nodePool.AutoscalingGroupId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}

	if err != nil {
		return nil
	}

	if hasAsg >= 1 {
		_ = d.Set("scaling_group_name", asg.AutoScalingGroupName)
		_ = d.Set("zones", asg.ZoneSet)
		_ = d.Set("scaling_group_project_id", asg.ProjectId)
		_ = d.Set("default_cooldown", asg.DefaultCooldown)
		_ = d.Set("termination_policies", asg.TerminationPolicySet)
	}

	taints := make([]map[string]interface{}, len(nodePool.Taints))
	for i, v := range nodePool.Taints {
		taint := map[string]interface{}{
			"key":    *v.Key,
			"value":  *v.Value,
			"effect": *v.Effect,
		}
		taints[i] = taint
	}
	d.Set("taints", taints)

	return nil
}

func resourceKubernetesNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_node_pool.create")()
	var (
		logId           = getLogId(contextNil)
		ctx             = context.WithValue(context.TODO(), logIdKey, logId)
		clusterId       = d.Get("cluster_id").(string)
		nodeConfig      = d.Get("node_config").([]interface{})
		enableAutoScale = d.Get("enable_auto_scale").(bool)
		configParas     = d.Get("auto_scaling_config").([]interface{})
		name            = d.Get("name").(string)
		iAdvanced       tke.InstanceAdvancedSettings
	)
	if len(configParas) != 1 {
		return fmt.Errorf("need only one auto_scaling_config")
	}

	if len(nodeConfig) > 1 {
		return fmt.Errorf("need only one node_config")
	}

	groupParaStr, err := composeParameterToAsScalingGroupParaSerial(d)
	if err != nil {
		return err
	}

	configParaStr, err := composedKubernetesAsScalingConfigParaSerial(configParas[0].(map[string]interface{}), meta)
	if err != nil {
		return err
	}

	labels := GetTkeLabels(d, "labels")
	taints := GetTkeTaints(d, "taints")

	//compose InstanceAdvancedSettings
	if workConfig, ok := d.GetOk("node_config"); ok {
		workConfigList := workConfig.([]interface{})
		if len(workConfigList) == 1 {
			workConfigPara := workConfigList[0].(map[string]interface{})
			setting := tkeGetInstanceAdvancedPara(workConfigPara, meta)
			iAdvanced = setting
		}
	}

	if temp, ok := d.GetOk("extra_args"); ok {
		extraArgs := helper.InterfacesStrings(temp.([]interface{}))
		for _, extraArg := range extraArgs {
			iAdvanced.ExtraArgs.Kubelet = append(iAdvanced.ExtraArgs.Kubelet, &extraArg)
		}
	}
	if temp, ok := d.GetOk("unschedulable"); ok {
		iAdvanced.Unschedulable = helper.Int64(int64(temp.(int)))
	}

	nodeOs := d.Get("node_os").(string)
	nodeOsType := d.Get("node_os_type").(string)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	nodePoolId, err := service.CreateClusterNodePool(ctx, clusterId, name, groupParaStr, configParaStr, enableAutoScale, nodeOs, nodeOsType, labels, taints, iAdvanced)
	if err != nil {
		return err
	}

	d.SetId(clusterId + FILED_SP + nodePoolId)

	// wait for status ok
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		nodePool, _, errRet := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if nodePool != nil && *nodePool.LifeState == "normal" {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("node pool status is %s, retry...", *nodePool.LifeState))
	})
	if err != nil {
		return err
	}

	instanceTypes := getNodePoolInstanceTypes(d)

	if len(instanceTypes) != 0 {
		err := service.ModifyClusterNodePoolInstanceTypes(ctx, clusterId, nodePoolId, instanceTypes)
		if err != nil {
			return err
		}
	}

	//modify os, instanceTypes and image
	err = resourceKubernetesNodePoolUpdate(d, meta)
	if err != nil {
		return err
	}

	return nil
}

func resourceKubernetesNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_node_pool.update")()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		asService = AsService{client: meta.(*TencentCloudClient).apiV3Conn}
		items     = strings.Split(d.Id(), FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	d.Partial(true)

	if d.HasChange("min_size") || d.HasChange("max_size") || d.HasChange("name") || d.HasChange("labels") || d.HasChange("taints") || d.HasChange("enable_auto_scale") || d.HasChange("node_os_type") || d.HasChange("node_os") {
		maxSize := int64(d.Get("max_size").(int))
		minSize := int64(d.Get("min_size").(int))
		enableAutoScale := d.Get("enable_auto_scale").(bool)
		name := d.Get("name").(string)
		nodeOs := d.Get("node_os").(string)
		nodeOsType := d.Get("node_os_type").(string)
		labels := GetTkeLabels(d, "labels")
		taints := GetTkeTaints(d, "taints")
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePool(ctx, clusterId, nodePoolId, name, enableAutoScale, minSize, maxSize, nodeOs, nodeOsType, labels, taints)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("min_size")
		d.SetPartial("max_size")
		d.SetPartial("name")
		d.SetPartial("enable_auto_scale")
		d.SetPartial("node_os")
		d.SetPartial("node_os_type")
		d.SetPartial("labels")
		d.SetPartial("taints")
	}

	if d.HasChange("scaling_group_name") ||
		d.HasChange("zones") ||
		d.HasChange("scaling_group_project_id") ||
		d.HasChange("default_cooldown") ||
		d.HasChange("termination_policies") {

		nodePool, _, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if err != nil {
			return err
		}

		var (
			scalingGroupId    = *nodePool.AutoscalingGroupId
			name              = d.Get("scaling_group_name").(string)
			projectId         = d.Get("scaling_group_project_id").(int)
			defaultCooldown   = d.Get("default_cooldown").(int)
			zones             []*string
			terminationPolicies []*string
		)

		if v, ok := d.GetOk("zones"); ok {
			for _, zone := range v.([]interface{}) {
				zones = append(zones, helper.String(zone.(string)))
			}
		}

		if v, ok := d.GetOk("termination_policies"); ok {
			for _, policy := range v.([]interface{}) {
				terminationPolicies = append(terminationPolicies, helper.String(policy.(string)))
			}
		}

		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := asService.ModifyScalingGroup(ctx, scalingGroupId, name, projectId, defaultCooldown, zones, terminationPolicies)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})

		if err != nil {
			return err
		}
		d.SetPartial("scaling_group_name")
		d.SetPartial("zones")
		d.SetPartial("scaling_group_project_id")
		d.SetPartial("default_cooldown")
		d.SetPartial("termination_policies")
	}

	if d.HasChange("desired_capacity") {
		desiredCapacity := int64(d.Get("desired_capacity").(int))
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePoolDesiredCapacity(ctx, clusterId, nodePoolId, desiredCapacity)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.SetPartial("desired_capacity")
	}

	if d.HasChange("auto_scaling_config.0.backup_instance_types") {
		instanceTypes := getNodePoolInstanceTypes(d)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := service.ModifyClusterNodePoolInstanceTypes(ctx, clusterId, nodePoolId, instanceTypes)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
		d.Set("auto_scaling_config.0.backup_instance_types", instanceTypes)
	}
	d.Partial(false)

	return resourceKubernetesNodePoolRead(d, meta)
}

func resourceKubernetesNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_node_pool.delete")()

	var (
		logId              = getLogId(contextNil)
		ctx                = context.WithValue(context.TODO(), logIdKey, logId)
		service            = TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
		items              = strings.Split(d.Id(), FILED_SP)
		deleteKeepInstance = d.Get("delete_keep_instance").(bool)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	//delete as group
	hasDelete := false
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err := service.DeleteClusterNodePool(ctx, clusterId, nodePoolId, deleteKeepInstance)

		if sdkErr, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			if sdkErr.Code == "InternalError.Param" && strings.Contains(sdkErr.Message, "Not Found") {
				hasDelete = true
				return nil
			}
		}
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if hasDelete {
		return nil
	}

	// wait for delete ok
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		nodePool, has, errRet := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if has {
			resource.RetryableError(fmt.Errorf("node pool %s still alive, status %s", nodePoolId, *nodePool.LifeState))
		}
		return nil
	})

	return err
}
