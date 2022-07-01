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
  multi_zone_subnet_policy = "EQUALITY"

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

Using Spot CVM Instance
```hcl
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
  multi_zone_subnet_policy = "EQUALITY"

  auto_scaling_config {
    instance_type      = var.default_instance_type
    system_disk_type   = "CLOUD_PREMIUM"
    system_disk_size   = "50"
    security_group_ids = ["sg-24vswocp"]
	instance_charge_type = "SPOTPAID"
    spot_instance_type = "one-time"
    spot_max_price = "1000"

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

}

```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// merge `instance_type` to `backup_instance_types` as param `instance_types`
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
			Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
			ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
			Description:  "Type of a CVM disk. Valid value: `CLOUD_PREMIUM` and `CLOUD_SSD`. Default is `CLOUD_PREMIUM`.",
		},
		"system_disk_size": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      50,
			ValidateFunc: validateIntegerInRange(50, 500),
			Description:  "Volume of system disk in GB. Default is `50`.",
		},
		"data_disk": {
			Type:        schema.TypeList,
			Optional:    true,
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
					"delete_with_instance": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Indicates whether the disk remove after instance terminated.",
					},
				},
			},
		},
		// payment
		"instance_charge_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID`. The default is `POSTPAID_BY_HOUR`. NOTE: `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time.",
		},
		"instance_charge_type_prepaid_period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue(CVM_PREPAID_PERIOD),
			Description:  "The tenancy (in month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
		},
		"instance_charge_type_prepaid_renew_flag": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validateAllowedStringValue(CVM_PREPAID_RENEW_FLAG),
			Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
		},
		"spot_instance_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateAllowedStringValue([]string{"one-time"}),
			Description:  "Type of spot instance, only support `one-time` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.",
		},
		"spot_max_price": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringNumber,
			Description:  "Max price of a spot instance, is the format of decimal string, for example \"0.50\". Note: it only works when instance_charge_type is set to `SPOTPAID`.",
		},
		"internet_charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
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
		"cam_role_name": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Name of cam role.",
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
			"multi_zone_subnet_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{MultiZoneSubnetPolicyPriority,
					MultiZoneSubnetPolicyEquality}),
				Description: "Multi-availability zone/subnet policy. Valid values: PRIORITY and EQUALITY. Default value: PRIORITY.",
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
				Computed:    true,
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
				Computed:    true,
				Description: "Seconds of scaling group cool down. Default value is `300`.",
			},
			"termination_policies": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
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
			"autoscaling_added_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total of autoscaling added node.",
			},
			"manually_added_total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total of manually added node.",
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
		request.SubnetIds = helper.InterfacesStringsPoint(subnetIds)
	}

	if v, ok := d.GetOk("scaling_mode"); ok {
		request.ServiceSettings = &as.ServiceSettings{ScalingMode: helper.String(v.(string))}
	}

	if v, ok := d.GetOk("multi_zone_subnet_policy"); ok {
		request.MultiZoneSubnetPolicy = helper.String(v.(string))
	}

	result = request.ToJsonString()

	return result, errRet
}

//This function is used to specify tke as group launch config, similar to kubernetesAsScalingConfigParaSerial, but less parameter
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
			deleteWithInstance, dOk := value["delete_with_instance"].(bool)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
			}
			if diskSize > 0 {
				dataDisk.DiskSize = &diskSize
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			if dOk {
				dataDisk.DeleteWithInstance = &deleteWithInstance
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

	chargeType, ok := dMap["instance_charge_type"].(string)
	if !ok || chargeType == "" {
		chargeType = INSTANCE_CHARGE_TYPE_POSTPAID
	}

	if chargeType == INSTANCE_CHARGE_TYPE_SPOTPAID {
		spotMaxPrice := dMap["spot_max_price"].(string)
		spotInstanceType := dMap["spot_instance_type"].(string)
		request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
			MarketType: helper.String("spot"),
			SpotOptions: &as.SpotMarketOptions{
				MaxPrice:         &spotMaxPrice,
				SpotInstanceType: &spotInstanceType,
			},
		}
	}

	if chargeType == INSTANCE_CHARGE_TYPE_PREPAID {
		period := dMap["instance_charge_type_prepaid_period"].(int)
		renewFlag := dMap["instance_charge_type_prepaid_renew_flag"].(string)
		request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
			Period:    helper.IntInt64(period),
			RenewFlag: &renewFlag,
		}
	}

	request.InstanceChargeType = &chargeType

	if v, ok := dMap["cam_role_name"]; ok {
		request.CamRoleName = helper.String(v.(string))
	}
	result = request.ToJsonString()
	return result, errRet
}

func composeAsLaunchConfigModifyRequest(d *schema.ResourceData, launchConfigId string) *as.ModifyLaunchConfigurationAttributesRequest {
	launchConfigRaw := d.Get("auto_scaling_config").([]interface{})
	dMap := launchConfigRaw[0].(map[string]interface{})
	request := as.NewModifyLaunchConfigurationAttributesRequest()
	request.LaunchConfigurationId = &launchConfigId

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
			deleteWithInstance, dOk := value["delete_with_instance"].(bool)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
			}
			if diskSize > 0 {
				dataDisk.DiskSize = &diskSize
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			if dOk {
				dataDisk.DeleteWithInstance = &deleteWithInstance
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	} else {
		request.DataDisks = []*as.DataDisk{}
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

	if v, ok := dMap["security_group_ids"]; ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
		}
	}

	chargeType, ok := dMap["instance_charge_type"].(string)

	if !ok || chargeType == "" {
		chargeType = INSTANCE_CHARGE_TYPE_POSTPAID
	}

	if chargeType == INSTANCE_CHARGE_TYPE_SPOTPAID {
		spotMaxPrice := dMap["spot_max_price"].(string)
		spotInstanceType := dMap["spot_instance_type"].(string)
		request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
			MarketType: helper.String("spot"),
			SpotOptions: &as.SpotMarketOptions{
				MaxPrice:         &spotMaxPrice,
				SpotInstanceType: &spotInstanceType,
			},
		}
	}

	if chargeType == INSTANCE_CHARGE_TYPE_PREPAID {
		period := dMap["instance_charge_type_prepaid_period"].(int)
		renewFlag := dMap["instance_charge_type_prepaid_renew_flag"].(string)
		request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
			Period:    helper.IntInt64(period),
			RenewFlag: &renewFlag,
		}
	}

	request.InstanceChargeType = &chargeType

	return request
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

	_ = d.Set("cluster_id", clusterId)

	//Describe Node Pool
	var (
		nodePool *tke.NodePool
	)

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		nodePool, has, err = service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		status := *nodePool.AutoscalingGroupStatus

		if status == "enabling" || status == "disabling" {
			return resource.RetryableError(fmt.Errorf("node pool status is %s, retrying", status))
		}

		return nil
	})

	if err != nil {
		return err
	}

	if !has {
		d.SetId("")
		return nil
	}

	_ = d.Set("name", nodePool.Name)
	_ = d.Set("status", nodePool.LifeState)
	AutoscalingAddedTotal := *nodePool.NodeCountSummary.AutoscalingAdded.Total
	ManuallyAddedTotal := *nodePool.NodeCountSummary.ManuallyAdded.Total
	_ = d.Set("autoscaling_added_total", AutoscalingAddedTotal)
	_ = d.Set("manually_added_total", ManuallyAddedTotal)
	_ = d.Set("node_count", AutoscalingAddedTotal+ManuallyAddedTotal)
	_ = d.Set("auto_scaling_group_id", nodePool.AutoscalingGroupId)
	_ = d.Set("launch_config_id", nodePool.LaunchConfigurationId)
	//set not force new parameters
	if nodePool.MaxNodesNum != nil {
		_ = d.Set("max_size", nodePool.MaxNodesNum)
	}
	if nodePool.MinNodesNum != nil {
		_ = d.Set("min_size", nodePool.MinNodesNum)
	}
	if nodePool.DesiredNodesNum != nil {
		_ = d.Set("desired_capacity", nodePool.DesiredNodesNum)
	}
	if nodePool.AutoscalingGroupStatus != nil {
		_ = d.Set("enable_auto_scale", *nodePool.AutoscalingGroupStatus == "enabled")
	}
	if nodePool.NodePoolOs != nil {
		_ = d.Set("node_os", nodePool.NodePoolOs)
	}
	if nodePool.OsCustomizeType != nil {
		_ = d.Set("node_os_type", nodePool.OsCustomizeType)
	}

	//set composed struct
	lables := make(map[string]interface{}, len(nodePool.Labels))
	for _, v := range nodePool.Labels {
		lables[*v.Name] = *v.Value
	}
	_ = d.Set("labels", lables)

	// set launch config
	launchCfg, hasLC, err := asService.DescribeLaunchConfigurationById(ctx, *nodePool.LaunchConfigurationId)

	if hasLC > 0 {
		launchConfig := make(map[string]interface{})
		if launchCfg.InstanceTypes != nil {
			insTypes := launchCfg.InstanceTypes
			launchConfig["instance_type"] = insTypes[0]
			backupInsTypes := insTypes[1:]
			if len(backupInsTypes) > 0 {
				launchConfig["backup_instance_types"] = helper.StringsInterfaces(backupInsTypes)
			}
		} else {
			launchConfig["instance_type"] = launchCfg.InstanceType
		}
		if launchCfg.SystemDisk.DiskType != nil {
			launchConfig["system_disk_type"] = launchCfg.SystemDisk.DiskType
		}
		if launchCfg.SystemDisk.DiskSize != nil {
			launchConfig["system_disk_size"] = launchCfg.SystemDisk.DiskSize
		}
		if launchCfg.InternetAccessible.InternetChargeType != nil {
			launchConfig["internet_charge_type"] = launchCfg.InternetAccessible.InternetChargeType
		}
		if launchCfg.InternetAccessible.InternetMaxBandwidthOut != nil {
			launchConfig["internet_max_bandwidth_out"] = launchCfg.InternetAccessible.InternetMaxBandwidthOut
		}
		if launchCfg.InternetAccessible.BandwidthPackageId != nil {
			launchConfig["bandwidth_package_id"] = launchCfg.InternetAccessible.BandwidthPackageId
		}
		if launchCfg.InternetAccessible.PublicIpAssigned != nil {
			launchConfig["public_ip_assigned"] = launchCfg.InternetAccessible.PublicIpAssigned
		}
		if launchCfg.InstanceChargeType != nil {
			launchConfig["instance_charge_type"] = launchCfg.InstanceChargeType
			if *launchCfg.InstanceChargeType == INSTANCE_CHARGE_TYPE_SPOTPAID && launchCfg.InstanceMarketOptions != nil {
				launchConfig["spot_instance_type"] = launchCfg.InstanceMarketOptions.SpotOptions.SpotInstanceType
				launchConfig["spot_max_price"] = launchCfg.InstanceMarketOptions.SpotOptions.MaxPrice
			}
			if *launchCfg.InstanceChargeType == INSTANCE_CHARGE_TYPE_PREPAID && launchCfg.InstanceChargePrepaid != nil {
				launchConfig["instance_charge_type_prepaid_period"] = launchCfg.InstanceChargePrepaid.Period
				launchConfig["instance_charge_type_prepaid_renew_flag"] = launchCfg.InstanceChargePrepaid.RenewFlag
			}
		}
		if len(launchCfg.DataDisks) > 0 {
			dataDisks := make([]map[string]interface{}, 0, len(launchCfg.DataDisks))
			for i := range launchCfg.DataDisks {
				item := launchCfg.DataDisks[i]
				disk := make(map[string]interface{})
				disk["disk_type"] = *item.DiskType
				disk["disk_size"] = *item.DiskSize
				if item.SnapshotId != nil {
					disk["snapshot_id"] = *item.SnapshotId
				}
				if item.DeleteWithInstance != nil {
					disk["delete_with_instance"] = *item.DeleteWithInstance
				}
				dataDisks = append(dataDisks, disk)
			}
			launchConfig["data_disk"] = dataDisks
		}
		if launchCfg.LoginSettings != nil {
			launchConfig["key_ids"] = helper.StringsInterfaces(launchCfg.LoginSettings.KeyIds)
		}
		// keep existing password in new launchConfig object
		if v, ok := d.GetOk("auto_scaling_config.0.password"); ok {
			launchConfig["password"] = v.(string)
		}
		launchConfig["security_group_ids"] = helper.StringsInterfaces(launchCfg.SecurityGroupIds)

		enableSecurity := launchCfg.EnhancedService.SecurityService.Enabled
		enableMonitor := launchCfg.EnhancedService.MonitorService.Enabled
		// Only declared or diff from exist will set.
		if _, ok := d.GetOk("enhanced_security_service"); ok || enableSecurity != nil {
			launchConfig["enhanced_security_service"] = *enableSecurity
		}
		if _, ok := d.GetOk("enhanced_monitor_service"); ok || enableMonitor != nil {
			launchConfig["enhanced_monitor_service"] = *enableMonitor
		}
		if _, ok := d.GetOk("cam_role_name"); ok || launchCfg.CamRoleName != nil {
			launchConfig["cam_role_name"] = launchCfg.CamRoleName
		}
		asgConfig := make([]interface{}, 0, 1)
		asgConfig = append(asgConfig, launchConfig)
		if err := d.Set("auto_scaling_config", asgConfig); err != nil {
			return err
		}
	}

	// asg node unschedulable
	clusterAsg, err := service.DescribeClusterAsGroupsByGroupId(ctx, clusterId, *nodePool.AutoscalingGroupId)

	if err != nil {
		return err
	}

	unschedulable := 0
	if clusterAsg != nil {
		if clusterAsg.IsUnschedulable != nil && *clusterAsg.IsUnschedulable {
			unschedulable = 1
		}
	}
	_ = d.Set("unschedulable", unschedulable)

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

	if hasAsg > 0 {
		_ = d.Set("scaling_group_name", asg.AutoScalingGroupName)
		_ = d.Set("zones", asg.ZoneSet)
		_ = d.Set("scaling_group_project_id", asg.ProjectId)
		_ = d.Set("default_cooldown", asg.DefaultCooldown)
		_ = d.Set("termination_policies", helper.StringsInterfaces(asg.TerminationPolicySet))
		_ = d.Set("vpc_id", asg.VpcId)
		_ = d.Set("retry_policy", asg.RetryPolicy)
		_ = d.Set("subnet_ids", helper.StringsInterfaces(asg.SubnetIdSet))

		// If not check, the diff between computed and default empty value leads to force replacement
		if _, ok := d.GetOk("multi_zone_subnet_policy"); ok {
			_ = d.Set("multi_zone_subnet_policy", asg.MultiZoneSubnetPolicy)
		}
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
		enableAutoScale = d.Get("enable_auto_scale").(bool)
		configParas     = d.Get("auto_scaling_config").([]interface{})
		name            = d.Get("name").(string)
		iAdvanced       tke.InstanceAdvancedSettings
	)
	if len(configParas) != 1 {
		return fmt.Errorf("need only one auto_scaling_config")
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
	if workConfig, ok := helper.InterfacesHeadMap(d, "node_config"); ok {
		iAdvanced = tkeGetInstanceAdvancedPara(workConfig, meta)
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
		client    = meta.(*TencentCloudClient).apiV3Conn
		service   = TkeService{client: client}
		asService = AsService{client: client}
		items     = strings.Split(d.Id(), FILED_SP)
	)
	if len(items) != 2 {
		return fmt.Errorf("resource_tc_kubernetes_node_pool id  is broken")
	}
	clusterId := items[0]
	nodePoolId := items[1]

	d.Partial(true)

	// LaunchConfig
	if d.HasChange("auto_scaling_config") {
		nodePool, _, err := service.DescribeNodePool(ctx, clusterId, nodePoolId)
		if err != nil {
			return err
		}
		launchConfigId := *nodePool.LaunchConfigurationId
		request := composeAsLaunchConfigModifyRequest(d, launchConfigId)
		_, err = client.UseAsClient().ModifyLaunchConfigurationAttributes(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return err
		}
		d.SetPartial("auto_scaling_config")
	}

	// ModifyClusterNodePool
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

	// ModifyScalingGroup
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
			scalingGroupId      = *nodePool.AutoscalingGroupId
			name                = d.Get("scaling_group_name").(string)
			projectId           = d.Get("scaling_group_project_id").(int)
			defaultCooldown     = d.Get("default_cooldown").(int)
			zones               []*string
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
			errCode := errRet.(*sdkErrors.TencentCloudSDKError).Code
			if errCode == "InternalError.UnexpectedInternal" {
				return nil
			}
			return retryError(errRet, InternalError)
		}
		if has {
			return resource.RetryableError(fmt.Errorf("node pool %s still alive, status %s", nodePoolId, *nodePool.LifeState))
		}
		return nil
	})

	return err
}
