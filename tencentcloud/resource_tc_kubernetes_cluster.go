/*
Provide a resource to create a kubernetes cluster.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "vpc" {
  default = "vpc-dk8zmwuf"
}

variable "subnet" {
  default = "subnet-pqfek0t8"
}

variable "default_instance_type" {
  default = "S1.SMALL1"
}

#examples for MANAGED_CLUSTER cluster
resource "tencentcloud_kubernetes_cluster" "managed_cluster" {
  vpc_id                  = var.vpc
  cluster_cidr            = "10.31.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

  worker_config {
    count                      = 2
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = var.subnet

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2"
  }

  cluster_deploy_type = "MANAGED_CLUSTER"
}

#examples for INDEPENDENT_CLUSTER cluster
resource "tencentcloud_kubernetes_cluster" "independing_cluster" {
  vpc_id                  = var.vpc
  cluster_cidr            = "10.1.0.0/16"
  cluster_max_pod_num     = 32
  cluster_name            = "test"
  cluster_desc            = "test cluster desc"
  cluster_max_service_num = 32

  master_config {
    count                      = 3
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = var.subnet

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "MMMZZXXccvv1212"
  }

  worker_config {
    count                      = 2
    availability_zone          = var.availability_zone
    instance_type              = var.default_instance_type
    system_disk_type           = "CLOUD_SSD"
    system_disk_size           = 60
    internet_charge_type       = "TRAFFIC_POSTPAID_BY_HOUR"
    internet_max_bandwidth_out = 100
    public_ip_assigned         = true
    subnet_id                  = var.subnet

    data_disk {
      disk_type = "CLOUD_PREMIUM"
      disk_size = 50
    }

    enhanced_security_service = false
    enhanced_monitor_service  = false
    user_data                 = "dGVzdA=="
    password                  = "ZZXXccvv1212"
  }

  labels = {
    "test1" = "test1",
    "test2" = "test2"
  }

  cluster_deploy_type = "INDEPENDENT_CLUSTER"
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func tkeCvmState() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the cvm.",
		},
		"instance_role": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Role of the cvm.",
		},
		"instance_state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "State of the cvm.",
		},
		"failed_reason": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Information of the cvm when it is failed.",
		},
	}
}

func tkeSecurityInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "User name of account.",
		},
		"password": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Password of account.",
		},
		"certification_authority": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The certificate used for access.",
		},
		"cluster_external_endpoint": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "External network address to access.",
		},
		"domain": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Domain name for access.",
		},
		"pgw_endpoint": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The Intranet address used for access.",
		},
		"security_policy": {
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Access policy.",
		},
	}
}

func TkeCvmCreateInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"count": {
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Default:     1,
			Description: "Number of cvm.",
		},
		"availability_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Indicates which availability zone will be used.",
		},
		"instance_name": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Default:     "sub machine of tke",
			Description: "Name of the CVMs.",
		},
		"instance_type": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "Specified types of CVM instance.",
		},
		// payment
		"instance_charge_type": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      CVM_CHARGE_TYPE_POSTPAID,
			ValidateFunc: validateAllowedStringValue(TKE_INSTANCE_CHARGE_TYPE),
			Description:  "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`, The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`, `PREPAID` instance will not terminated after cluster deleted, and may not allow to delete before expired.",
		},
		"instance_charge_type_prepaid_period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      1,
			ValidateFunc: validateAllowedIntValue(CVM_PREPAID_PERIOD),
			Description:  "The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.",
		},
		"instance_charge_type_prepaid_renew_flag": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      CVM_PREPAID_RENEW_FLAG_NOTIFY_AND_MANUAL_RENEW,
			ValidateFunc: validateAllowedStringValue(CVM_PREPAID_RENEW_FLAG),
			Description:  "When enabled, the CVM instance will be renew automatically when it reach the end of the prepaid tenancy. Valid values are `NOTIFY_AND_AUTO_RENEW`, `NOTIFY_AND_MANUAL_RENEW` and `DISABLE_NOTIFY_AND_MANUAL_RENEW`. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validateStringLengthInRange(4, 100),
			Description:  "Private network ID.",
		},
		"system_disk_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
			ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
			Description:  "Type of a CVM disk, and available values include CLOUD_PREMIUM and CLOUD_SSD. Default is CLOUD_PREMIUM.",
		},
		"system_disk_size": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			Default:      50,
			ValidateFunc: validateIntegerInRange(50, 500),
			Description:  "Volume of system disk in GB. Default is 50.",
		},
		"data_disk": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			MaxItems:    11,
			Description: "Configurations of data disk.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"disk_type": {
						Type:         schema.TypeString,
						ForceNew:     true,
						Optional:     true,
						Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
						ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
						Description:  "Types of disk, available values: CLOUD_PREMIUM and CLOUD_SSD.",
					},
					"disk_size": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Optional:    true,
						Default:     0,
						Description: "Volume of disk in GB. Default is 0.",
					},
					"snapshot_id": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "Data disk snapshot ID.",
					},
				},
			},
		},
		"internet_charge_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID_BY_HOUR,
			ValidateFunc: validateAllowedStringValue(INTERNET_CHARGE_ALLOW_TYPE),
			Description:  "Charge types for network traffic. Available values include TRAFFIC_POSTPAID_BY_HOUR.",
		},
		"internet_max_bandwidth_out": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Default:     0,
			Description: "Max bandwidth of Internet access in Mbps. Default is 0.",
		},
		"public_ip_assigned": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Description: "Specify whether to assign an Internet IP address.",
		},
		"password": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validateAsConfigPassword,
			Description:  "Password to access, should be set if `key_ids` not set.",
		},
		"key_ids": {
			MaxItems:    1,
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "ID list of keys, should be set if `password` not set.",
		},
		"security_group_ids": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Security groups to which a CVM instance belongs.",
		},
		"enhanced_security_service": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "To specify whether to enable cloud security service. Default is TRUE.",
		},
		"enhanced_monitor_service": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
		},
		"user_data": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "ase64-encoded User Data text, the length limit is 16KB.",
		},
	}
}

func resourceTencentCloudTkeCluster() *schema.Resource {
	schemaBody := map[string]*schema.Schema{
		"cluster_name": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Name of the cluster.",
		},
		"cluster_desc": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Description of the cluster.",
		},
		"cluster_os": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      TKE_CLUSTER_OS_UBUNTU16,
			ValidateFunc: validateAllowedStringValue(TKE_CLUSTER_OS),
			Description: "Operating system of the cluster, the available values include: '" + strings.Join(TKE_CLUSTER_OS, "','") +
				"'. Default is '" + TKE_CLUSTER_OS_UBUNTU16 + "'.",
		},
		"cluster_os_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      TKE_CLUSTER_OS_TYPE_GENERAL,
			ValidateFunc: validateAllowedStringValue(TKE_CLUSTER_OS_TYPES),
			Description: "Image type of the cluster os, the available values include: '" + strings.Join(TKE_CLUSTER_OS_TYPES, "','") +
				"'. Default is '" + TKE_CLUSTER_OS_TYPE_GENERAL + "'. 'DOCKER_CUSTOMIZE' means 'TKE-Optimized'. " +
				"Only 'centos7.6x86_64' or 'ubuntu18.04.1 LTSx86_64' support 'DOCKER_CUSTOMIZE' now.",
		},
		"container_runtime": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      TKE_RUNTIME_DOCKER,
			ValidateFunc: validateAllowedStringValue(TKE_RUNTIMES),
			Description:  "Runtime type of the cluster, the available values include: 'docker' and 'containerd'. Default is 'docker'.",
		},
		"cluster_deploy_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      TKE_DEPLOY_TYPE_MANAGED,
			ValidateFunc: validateAllowedStringValue(TKE_DEPLOY_TYPES),
			Description:  "Deployment type of the cluster, the available values include: 'MANAGED_CLUSTER' and 'INDEPENDENT_CLUSTER', Default is 'MANAGED_CLUSTER'.",
		},
		"cluster_version": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Default:     "1.10.5",
			Description: "Version of the cluster, Default is '1.10.5'.",
		},
		"cluster_ipvs": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "Indicates whether ipvs is enabled. Default is true.",
		},
		"vpc_id": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validateStringLengthInRange(4, 100),
			Description:  "Vpc Id of the cluster.",
		},
		"cluster_internet": {
			Type:        schema.TypeBool,
			Default:     false,
			Optional:    true,
			Description: "Open internet access or not.",
		},
		"cluster_intranet": {
			Type:        schema.TypeBool,
			Default:     false,
			Optional:    true,
			Description: "Open intranet access or not.",
		},
		"managed_cluster_internet_security_policies": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Description: "Security policies for managed cluster internet, like:'192.168.1.0/24' or '113.116.51.27', '0.0.0.0/0' means all." +
				" This field can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true." +
				" `managed_cluster_internet_security_policies` can not delete or empty once be set.",
		},
		"cluster_intranet_subnet_id": {
			Type:     schema.TypeString,
			Optional: true,
			Description: "Subnet id who can access this independent cluster, this field must and can only set  when `cluster_intranet` is true." +
				" `cluster_intranet_subnet_id` can not modify once be set.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Description: "Project ID, default value is 0.",
		},
		"cluster_cidr": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.",
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				value := v.(string)
				_, ipnet, err := net.ParseCIDR(value)
				if err != nil {
					errors = append(errors, fmt.Errorf("%q must contain a valid CIDR, got error parsing: %s", k, err))
					return
				}
				if ipnet == nil || value != ipnet.String() {
					errors = append(errors, fmt.Errorf("%q must contain a valid network CIDR, expected %q, got %q", k, ipnet, value))
					return
				}
				if !strings.Contains(value, "/") {
					errors = append(errors, fmt.Errorf("%q must be a network segment", k))
					return
				}
				if !strings.HasPrefix(value, "10.") && !strings.HasPrefix(value, "192.168.") && !strings.HasPrefix(value, "172.") {
					errors = append(errors, fmt.Errorf("%q must in 10. | 192.168. | 172.[16-31]", k))
					return
				}

				if strings.HasPrefix(value, "172.") {
					nextNo := strings.Split(value, ".")[1]
					no, _ := strconv.ParseInt(nextNo, 10, 64)
					if no < 16 || no > 31 {
						errors = append(errors, fmt.Errorf("%q must in 10. | 192.168. | 172.[16-31]", k))
						return
					}
				}
				return
			},
		},
		"ignore_cluster_cidr_conflict": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     false,
			Description: "Indicates whether to ignore the cluster cidr conflict error. Default is false.",
		},
		"cluster_max_pod_num": {
			Type:     schema.TypeInt,
			ForceNew: true,
			Optional: true,
			Default:  256,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				value := v.(int)
				if value%16 != 0 {
					errors = append(errors, fmt.Errorf(
						"%q  has to be a multiple of 16 ", k))
				}
				if value < 32 {
					errors = append(errors, fmt.Errorf(
						"%q cannot be lower than %d: %d", k, 32, value))
				}
				return
			},
			Description: "The maximum number of Pods per node in the cluster. Default is 256. Must be a multiple of 16 and large than 32.",
		},
		"cluster_max_service_num": {
			Type:     schema.TypeInt,
			ForceNew: true,
			Optional: true,
			Default:  256,
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				value := v.(int)
				if value%16 != 0 {
					errors = append(errors, fmt.Errorf(
						"%q  has to be a multiple of 16 ", k))
				}
				return
			},
			Description: "The maximum number of services in the cluster. Default is 256. Must be a multiple of 16.",
		},
		"master_config": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: TkeCvmCreateInfo(),
			},
			Description: "Deploy the machine configuration information of the 'MASTER_ETCD' service, and create <=7 units for common users.",
		},
		"worker_config": {
			Type:     schema.TypeList,
			ForceNew: true,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: TkeCvmCreateInfo(),
			},
			Description: "Deploy the machine configuration information of the 'WORKER' service, and create <=20 units for common users. The other 'WORK' service are added by 'tencentcloud_kubernetes_worker'.",
		},
		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "The tags of the cluster.",
		},

		// Computed values
		"cluster_node_num": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Number of nodes in the cluster.",
		},
		"worker_instances_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: tkeCvmState(),
			},
			Description: "An information list of cvm within the 'WORKER' clusters. Each element contains the following attributes:",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Description: "Labels of tke cluster.",
		},
	}

	for k, v := range tkeSecurityInfo() {
		schemaBody[k] = v
	}

	return &schema.Resource{
		Create: resourceTencentCloudTkeClusterCreate,
		Read:   resourceTencentCloudTkeClusterRead,
		Update: resourceTencentCloudTkeClusterUpdate,
		Delete: resourceTencentCloudTkeClusterDelete,
		Schema: schemaBody,
	}
}

func tkeGetCvmRunInstancesPara(dMap map[string]interface{}, meta interface{},
	vpcId string, projectId int64) (cvmJson string, count int64, errRet error) {

	request := cvm.NewRunInstancesRequest()

	var place cvm.Placement
	request.Placement = &place

	place.ProjectId = &projectId

	configRegion := meta.(*TencentCloudClient).apiV3Conn.Region
	if v, ok := dMap["availability_zone"]; ok {
		if !strings.Contains(v.(string), configRegion) {
			errRet = fmt.Errorf("availability_zone[%s] should in [%s]", v.(string), configRegion)
			return
		}
		place.Zone = helper.String(v.(string))
	}

	if v, ok := dMap["instance_type"]; ok {
		request.InstanceType = helper.String(v.(string))
	} else {
		errRet = fmt.Errorf("instance_type must be set.")
		return
	}

	subnetId := ""

	if v, ok := dMap["subnet_id"]; ok {
		subnetId = v.(string)
	}

	if (vpcId == "" && subnetId != "") ||
		(vpcId != "" && subnetId == "") {
		errRet = fmt.Errorf("Parameters cvm.`subnet_id` and cluster.`vpc_id` are both set or neither")
		return
	}

	if vpcId != "" {
		request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{
			VpcId:    &vpcId,
			SubnetId: &subnetId,
		}
	}

	if v, ok := dMap["system_disk_type"]; ok {
		if request.SystemDisk == nil {
			request.SystemDisk = &cvm.SystemDisk{}
		}
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := dMap["system_disk_size"]; ok {
		if request.SystemDisk == nil {
			request.SystemDisk = &cvm.SystemDisk{}
		}
		request.SystemDisk.DiskSize = helper.Int64(int64(v.(int)))

	}

	if v, ok := dMap["data_disk"]; ok {

		dataDisks := v.([]interface{})
		request.DataDisks = make([]*cvm.DataDisk, 0, len(dataDisks))

		for _, d := range dataDisks {

			var (
				value      = d.(map[string]interface{})
				diskType   = value["disk_type"].(string)
				diskSize   = int64(value["disk_size"].(int))
				snapshotId = value["snapshot_id"].(string)
				dataDisk   = cvm.DataDisk{
					DiskType: &diskType,
					DiskSize: &diskSize,
				}
			)
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	if v, ok := dMap["internet_charge_type"]; ok {

		if request.InternetAccessible == nil {
			request.InternetAccessible = &cvm.InternetAccessible{}
		}
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}

	if v, ok := dMap["internet_max_bandwidth_out"]; ok {
		if request.InternetAccessible == nil {
			request.InternetAccessible = &cvm.InternetAccessible{}
		}
		request.InternetAccessible.InternetMaxBandwidthOut = helper.Int64(int64(v.(int)))
	}

	if v, ok := dMap["public_ip_assigned"]; ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	if v, ok := dMap["password"]; ok {
		if request.LoginSettings == nil {
			request.LoginSettings = &cvm.LoginSettings{}
		}

		if v.(string) != "" {
			request.LoginSettings.Password = helper.String(v.(string))
		}
	}

	if v, ok := dMap["instance_name"]; ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := dMap["key_ids"]; ok {
		if request.LoginSettings == nil {
			request.LoginSettings = &cvm.LoginSettings{}
		}
		keyIds := v.([]interface{})

		if len(keyIds) != 0 {
			request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
			for i := range keyIds {
				keyId := keyIds[i].(string)
				request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
			}
		}
	}

	if request.LoginSettings.Password == nil && len(request.LoginSettings.KeyIds) == 0 {
		errRet = fmt.Errorf("Parameters cvm.`key_ids` and cluster.`password` should be set one")
		return
	}

	if request.LoginSettings.Password != nil && len(request.LoginSettings.KeyIds) != 0 {
		errRet = fmt.Errorf("Parameters cvm.`key_ids` and cluster.`password` can only be supported one")
		return
	}

	if v, ok := dMap["security_group_ids"]; ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
		}
	}

	if v, ok := dMap["enhanced_security_service"]; ok {

		if request.EnhancedService == nil {
			request.EnhancedService = &cvm.EnhancedService{}
		}

		securityService := v.(bool)
		request.EnhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := dMap["enhanced_monitor_service"]; ok {
		if request.EnhancedService == nil {
			request.EnhancedService = &cvm.EnhancedService{}
		}
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}
	if v, ok := dMap["user_data"]; ok {
		request.UserData = helper.String(v.(string))
	}
	if v, ok := dMap["instance_charge_type"]; ok {
		instanceChargeType := v.(string)
		request.InstanceChargeType = &instanceChargeType
		if instanceChargeType == CVM_CHARGE_TYPE_PREPAID {
			request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{}
			if period, ok := dMap["instance_charge_type_prepaid_period"]; ok {
				periodInt64 := int64(period.(int))
				request.InstanceChargePrepaid.Period = &periodInt64
			} else {
				errRet = fmt.Errorf("instance charge type prepaid period can not be empty when charge type is %s",
					instanceChargeType)
				return
			}
			if renewFlag, ok := dMap["instance_charge_type_prepaid_renew_flag"]; ok {
				request.InstanceChargePrepaid.RenewFlag = helper.String(renewFlag.(string))
			}
		}
	}
	if v, ok := dMap["count"]; ok {
		count = int64(v.(int))
	} else {
		count = 1
	}
	request.InstanceCount = &count

	cvmJson = request.ToJsonString()

	cvmJson = strings.Replace(cvmJson, `"Password":"",`, "", -1)

	return
}

func resourceTencentCloudTkeClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		basic            ClusterBasicSetting
		advanced         ClusterAdvancedSettings
		cvms             RunInstancesForNode
		iAdvanced        InstanceAdvancedSettings
		cidrSet          ClusterCidrSettings
		securityPolicies []string
		clusterInternet  = d.Get("cluster_internet").(bool)
		clusterIntranet  = d.Get("cluster_intranet").(bool)
		intranetSubnetId = d.Get("cluster_intranet_subnet_id").(string)
	)

	if temp, ok := d.GetOkExists("managed_cluster_internet_security_policies"); ok {
		securityPolicies = helper.InterfacesStrings(temp.([]interface{}))
	}
	clusterDeployType := d.Get("cluster_deploy_type").(string)

	if clusterIntranet && intranetSubnetId == "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` must set when `cluster_intranet` is true")
	}
	if !clusterIntranet && intranetSubnetId != "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` can only set when `cluster_intranet` is true")
	}

	if clusterDeployType == TKE_DEPLOY_TYPE_INDEPENDENT {
		if len(securityPolicies) != 0 {
			return fmt.Errorf("`managed_cluster_internet_security_policies` can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true.")
		}
	}

	if clusterDeployType == TKE_DEPLOY_TYPE_MANAGED {
		if !clusterInternet && len(securityPolicies) != 0 {
			return fmt.Errorf("`managed_cluster_internet_security_policies` can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true.")
		}
	}

	vpcId := d.Get("vpc_id").(string)
	if vpcId != "" {
		basic.VpcId = vpcId
	}

	basic.ProjectId = int64(d.Get("project_id").(int))

	cluster_os := d.Get("cluster_os").(string)

	basic.ClusterOs = tkeClusterOsMap[cluster_os]
	basic.ClusterOsType = d.Get("cluster_os_type").(string)

	if basic.ClusterOsType == TKE_CLUSTER_OS_TYPE_DOCKER_CUSTOMIZE {
		if cluster_os != TKE_CLUSTER_OS_UBUNTU18 && cluster_os != TKE_CLUSTER_OS_CENTOS76 {
			return fmt.Errorf("Only 'centos7.6x86_64' or 'ubuntu18.04.1 LTSx86_64' support 'DOCKER_CUSTOMIZE' now, can not be " + basic.ClusterOs)
		}
	}

	basic.ClusterVersion = d.Get("cluster_version").(string)
	if v, ok := d.GetOk("cluster_name"); ok {
		basic.ClusterName = v.(string)
	}
	if v, ok := d.GetOk("cluster_desc"); ok {
		basic.ClusterDescription = v.(string)
	}

	advanced.ContainerRuntime = d.Get("container_runtime").(string)
	advanced.Ipvs = d.Get("cluster_ipvs").(bool)

	cidrSet.ClusterCidr = d.Get("cluster_cidr").(string)
	cidrSet.IgnoreClusterCidrConflict = d.Get("ignore_cluster_cidr_conflict").(bool)
	cidrSet.MaxClusterServiceNum = int64(d.Get("cluster_max_service_num").(int))
	cidrSet.MaxNodePodNum = int64(d.Get("cluster_max_pod_num").(int))

	items := strings.Split(cidrSet.ClusterCidr, "/")
	if len(items) != 2 {
		return fmt.Errorf("`cluster_cidr` must be network segment ")
	}

	bitNumber, err := strconv.ParseInt(items[1], 10, 64)

	if err != nil {
		return fmt.Errorf("`cluster_cidr` must be network segment ")
	}

	if math.Pow(2, float64(32-bitNumber)) <= float64(cidrSet.MaxNodePodNum) {
		return fmt.Errorf("`cluster_cidr` Network segment range is too small, can not cover cluster_max_service_num")
	}

	if masters, ok := d.GetOk("master_config"); ok {
		if clusterDeployType == TKE_DEPLOY_TYPE_MANAGED {
			return fmt.Errorf("if `cluster_deploy_type` is `MANAGED_CLUSTER` , You don't need define the master yourself")
		}
		var masterCount int64 = 0
		masterList := masters.([]interface{})
		for index := range masterList {
			master := masterList[index].(map[string]interface{})
			paraJson, count, err := tkeGetCvmRunInstancesPara(master, meta, vpcId, basic.ProjectId)
			if err != nil {
				return err
			}

			cvms.Master = append(cvms.Master, paraJson)
			masterCount += count
		}
		if masterCount < 3 {
			return fmt.Errorf("if `cluster_deploy_type` is `TKE_DEPLOY_TYPE_INDEPENDENT` len(master_config) should >=3")
		}
	} else {
		if clusterDeployType == TKE_DEPLOY_TYPE_INDEPENDENT {
			return fmt.Errorf("if `cluster_deploy_type` is `TKE_DEPLOY_TYPE_INDEPENDENT` , You need define the master yourself")
		}

	}

	if workers, ok := d.GetOk("worker_config"); ok {
		workerList := workers.([]interface{})
		for index := range workerList {
			worker := workerList[index].(map[string]interface{})
			paraJson, _, err := tkeGetCvmRunInstancesPara(worker, meta, vpcId, basic.ProjectId)

			if err != nil {
				return err
			}
			cvms.Work = append(cvms.Work, paraJson)
		}
	}

	tags := helper.GetTags(d, "tags")

	labels := make([]*tke.Label, 0)
	if v, ok := d.GetOk("labels"); ok {
		vlabels := v.(map[string]interface{})

		for key, value := range vlabels {
			keyTmp, valueTmp := key, value

			valueResult, ok := valueTmp.(string)
			if !ok {
				continue
			}

			labels = append(labels, &tke.Label{Name: &keyTmp, Value: &valueResult})
		}

		iAdvanced.Labels = labels
	}

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateCluster(ctx, basic, advanced, cvms, iAdvanced, cidrSet, tags)
	if err != nil {
		return err
	}

	d.SetId(id)

	_, _, err = service.DescribeClusterInstances(ctx, d.Id())

	if err != nil {
		// create often cost more than 20 Minutes.
		err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
			_, _, err = service.DescribeClusterInstances(ctx, d.Id())

			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}

			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
	}

	if err != nil {
		return err
	}

	//intranet
	if clusterIntranet {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := service.CreateClusterEndpoint(ctx, id, intranetSubnetId, false)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			status, message, inErr := service.DescribeClusterEndpointStatus(ctx, id)
			if inErr != nil {
				return retryError(inErr)
			}
			if status == TkeInternetStatusCreating {
				return resource.RetryableError(
					fmt.Errorf("%s create intranet cluster endpoint status still is %s", id, status))
			}
			if status == TkeInternetStatusNotfound || status == TkeInternetStatusCreated {
				return nil
			}
			return resource.NonRetryableError(
				fmt.Errorf("%s create intranet cluster endpoint error ,status is %s,message is %s", id, status, message))
		})
		if err != nil {
			return err
		}
	}

	//TKE_DEPLOY_TYPE_MANAGED Open the internet
	if clusterDeployType == TKE_DEPLOY_TYPE_MANAGED && clusterInternet {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := service.CreateClusterEndpointVip(ctx, id, securityPolicies)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			status, message, inErr := service.DescribeClusterEndpointVipStatus(ctx, id)
			if inErr != nil {
				return retryError(inErr)
			}
			if status == TkeInternetStatusCreating {
				return resource.RetryableError(
					fmt.Errorf("%s create cluster endpoint vip status still is %s", id, status))
			}
			if status == TkeInternetStatusNotfound || status == TkeInternetStatusCreated {
				return nil
			}
			return resource.NonRetryableError(
				fmt.Errorf("%s create cluster endpoint vip error ,status is %s,message is %s", id, status, message))
		})
		if err != nil {
			return err
		}
	}

	//TKE_DEPLOY_TYPE_INDEPENDENT Open the internet
	if clusterDeployType == TKE_DEPLOY_TYPE_INDEPENDENT && clusterInternet {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := service.CreateClusterEndpoint(ctx, id, "", true)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			status, message, inErr := service.DescribeClusterEndpointStatus(ctx, id)
			if inErr != nil {
				return retryError(inErr)
			}
			if status == TkeInternetStatusCreating {
				return resource.RetryableError(
					fmt.Errorf("%s create cluster internet endpoint status still is %s", id, status))
			}
			if status == TkeInternetStatusNotfound || status == TkeInternetStatusCreated {
				return nil
			}
			return resource.NonRetryableError(
				fmt.Errorf("%s create cluster internet endpoint error ,status is %s,message is %s", id, status, message))
		})
		if err != nil {
			return err
		}
	}

	if err = resourceTencentCloudTkeClusterRead(d, meta); err != nil {
		log.Printf("[WARN]%s resource.kubernetes_cluster.read after create fail , %s", logId, err.Error())
		return err
	}
	return nil
}

func resourceTencentCloudTkeClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	info, has, err := service.DescribeCluster(ctx, d.Id())
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			info, has, err = service.DescribeCluster(ctx, d.Id())
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

	_ = d.Set("cluster_name", info.ClusterName)
	_ = d.Set("cluster_desc", info.ClusterDescription)
	_ = d.Set("cluster_os", tkeToShowClusterOs(info.ClusterOs))
	_ = d.Set("cluster_deploy_type", info.DeployType)
	_ = d.Set("cluster_version", info.ClusterVersion)
	_ = d.Set("cluster_ipvs", info.Ipvs)
	_ = d.Set("vpc_id", info.VpcId)
	_ = d.Set("project_id", info.ProjectId)
	_ = d.Set("cluster_cidr", info.ClusterCidr)
	_ = d.Set("ignore_cluster_cidr_conflict", info.IgnoreClusterCidrConflict)
	_ = d.Set("cluster_max_pod_num", info.MaxNodePodNum)
	_ = d.Set("cluster_max_service_num", info.MaxClusterServiceNum)
	_ = d.Set("cluster_node_num", info.ClusterNodeNum)
	_ = d.Set("tags", info.Tags)

	_, workers, err := service.DescribeClusterInstances(ctx, d.Id())
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			_, workers, err = service.DescribeClusterInstances(ctx, d.Id())

			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	workerInstancesList := make([]map[string]interface{}, 0, len(workers))
	for _, worker := range workers {
		tempMap := make(map[string]interface{})
		tempMap["instance_id"] = worker.InstanceId
		tempMap["instance_role"] = worker.InstanceRole
		tempMap["instance_state"] = worker.InstanceState
		tempMap["failed_reason"] = worker.FailedReason
		workerInstancesList = append(workerInstancesList, tempMap)
	}

	_ = d.Set("worker_instances_list", workerInstancesList)

	securityRet, err := service.DescribeClusterSecurity(ctx, d.Id())

	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			securityRet, err = service.DescribeClusterSecurity(ctx, d.Id())
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	var emptyStrFunc = func(ptr *string) string {
		if ptr == nil {
			return ""
		} else {
			return *ptr
		}
	}

	policies := make([]string, 0, len(securityRet.Response.SecurityPolicy))
	for _, v := range securityRet.Response.SecurityPolicy {
		policies = append(policies, *v)
	}

	_ = d.Set("user_name", emptyStrFunc(securityRet.Response.UserName))
	_ = d.Set("password", emptyStrFunc(securityRet.Response.Password))
	_ = d.Set("certification_authority", emptyStrFunc(securityRet.Response.CertificationAuthority))
	_ = d.Set("cluster_external_endpoint", emptyStrFunc(securityRet.Response.ClusterExternalEndpoint))
	_ = d.Set("domain", emptyStrFunc(securityRet.Response.Domain))
	_ = d.Set("pgw_endpoint", emptyStrFunc(securityRet.Response.PgwEndpoint))
	_ = d.Set("security_policy", policies)

	if emptyStrFunc(securityRet.Response.ClusterExternalEndpoint) == "" {
		_ = d.Set("cluster_internet", false)
	} else {
		_ = d.Set("cluster_internet", true)
	}

	if emptyStrFunc(securityRet.Response.PgwEndpoint) == "" {
		_ = d.Set("cluster_intranet", false)
	} else {
		_ = d.Set("cluster_intranet", true)
	}

	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		_, clusterAsGroupSet, err := service.DescribeClusterAsGroups(ctx, d.Id())

		if err != nil {
			return retryError(err)
		}

		if len(clusterAsGroupSet) == 0 {
			d.Set("labels", "")
			return nil
		}

		var labelsMap = map[string]string{}
		for _, value := range clusterAsGroupSet {
			labels := value.Labels

			if len(labels) == 0 {
				d.Set("labels", "")
				return nil
			}

			for _, v := range labels {
				labelsMap[*v.Name] = *v.Value
			}

		}

		marshal, _ := json.Marshal(labelsMap)
		d.Set("labels", string(marshal))

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudTkeClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	client := meta.(*TencentCloudClient).apiV3Conn
	service := TagService{client: client}
	tkeService := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := client.Region
	d.Partial(true)

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("ccs", "cluster", region, id)
		if err := service.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	var (
		securityPolicies []string
		clusterInternet  = d.Get("cluster_internet").(bool)
		clusterIntranet  = d.Get("cluster_intranet").(bool)
		intranetSubnetId = d.Get("cluster_intranet_subnet_id").(string)
	)

	if temp, ok := d.GetOkExists("managed_cluster_internet_security_policies"); ok {
		securityPolicies = helper.InterfacesStrings(temp.([]interface{}))
	}
	clusterDeployType := d.Get("cluster_deploy_type").(string)

	if d.HasChange("cluster_intranet_subnet_id") {
		oldKey, newKey := d.GetChange("cluster_intranet_subnet_id")
		if (oldKey.(string) != "" && newKey.(string) == "") || (oldKey.(string) != "" && newKey.(string) != "") {
			return fmt.Errorf("`cluster_intranet_subnet_id` can not modify once be set")
		}
	}
	if clusterIntranet && intranetSubnetId == "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` must set when `cluster_intranet` is true")
	}

	if clusterDeployType == TKE_DEPLOY_TYPE_INDEPENDENT {
		if len(securityPolicies) != 0 {
			return fmt.Errorf("`managed_cluster_internet_security_policies` can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' ")
		}
	}

	if d.HasChange("cluster_intranet") {
		//open intranet
		if clusterIntranet {
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				inErr := tkeService.CreateClusterEndpoint(ctx, id, intranetSubnetId, false)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if err != nil {
				return err
			}
			err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
				status, message, inErr := tkeService.DescribeClusterEndpointStatus(ctx, id)
				if inErr != nil {
					return retryError(inErr)
				}
				if status == TkeInternetStatusCreating {
					return resource.RetryableError(
						fmt.Errorf("%s create intranet cluster endpoint status still is %s", id, status))
				}
				if status == TkeInternetStatusNotfound || status == TkeInternetStatusCreated {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("%s create intranet cluster endpoint error ,status is %s,message is %s", id, status, message))
			})
			if err != nil {
				return err
			}
			//close
		} else {
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				inErr := tkeService.DeleteClusterEndpoint(ctx, id, false)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if err != nil {
				return err
			}
			err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
				status, message, inErr := tkeService.DescribeClusterEndpointStatus(ctx, id)
				if inErr != nil {
					return retryError(inErr)
				}
				if status == TkeInternetStatusDeleting {
					return resource.RetryableError(
						fmt.Errorf("%s close cluster internet endpoint status still is %s", id, status))
				}
				if status == TkeInternetStatusNotfound || status == TkeInternetStatusDeleted || status == TkeInternetStatusCreated {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("%s close cluster internet endpoint error ,status is %s,message is %s", id, status, message))
			})
			if err != nil {
				return err
			}
		}

		d.SetPartial("cluster_intranet")
	}

	if d.HasChange("cluster_internet") {

		//TKE_DEPLOY_TYPE_INDEPENDENT open internet
		if clusterDeployType == TKE_DEPLOY_TYPE_INDEPENDENT && clusterInternet {
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				inErr := tkeService.CreateClusterEndpoint(ctx, id, "", true)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if err != nil {
				return err
			}
			err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
				status, message, inErr := tkeService.DescribeClusterEndpointStatus(ctx, id)
				if inErr != nil {
					return retryError(inErr)
				}
				if status == TkeInternetStatusCreating {
					return resource.RetryableError(
						fmt.Errorf("%s create cluster internet endpoint status still is %s", id, status))
				}
				if status == TkeInternetStatusNotfound || status == TkeInternetStatusCreated {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("%s create cluster internet endpoint error ,status is %s,message is %s", id, status, message))
			})
			if err != nil {
				return err
			}
		}

		//TKE_DEPLOY_TYPE_INDEPENDENT close internet
		if clusterDeployType == TKE_DEPLOY_TYPE_INDEPENDENT && !clusterInternet {
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				inErr := tkeService.DeleteClusterEndpoint(ctx, id, true)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if err != nil {
				return err
			}
			err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
				status, message, inErr := tkeService.DescribeClusterEndpointStatus(ctx, id)
				if inErr != nil {
					return retryError(inErr)
				}
				if status == TkeInternetStatusDeleting {
					return resource.RetryableError(
						fmt.Errorf("%s close cluster internet endpoint status still is %s", id, status))
				}
				if status == TkeInternetStatusNotfound || status == TkeInternetStatusDeleted || status == TkeInternetStatusCreated {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("%s close cluster internet endpoint error ,status is %s,message is %s", id, status, message))
			})
			if err != nil {
				return err
			}
		}

		//TKE_DEPLOY_TYPE_MANAGED open internet
		if clusterDeployType == TKE_DEPLOY_TYPE_MANAGED && clusterInternet {
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				inErr := tkeService.CreateClusterEndpointVip(ctx, id, securityPolicies)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if err != nil {
				return err
			}
			err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
				status, message, inErr := tkeService.DescribeClusterEndpointVipStatus(ctx, id)
				if inErr != nil {
					return retryError(inErr)
				}
				if status == TkeInternetStatusCreating {
					return resource.RetryableError(
						fmt.Errorf("%s create cluster endpoint vip status still is %s", id, status))
				}
				if status == TkeInternetStatusNotfound || status == TkeInternetStatusCreated {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("%s create cluster endpoint vip error ,status is %s,message is %s", id, status, message))
			})
			if err != nil {
				return err
			}
		}

		//TKE_DEPLOY_TYPE_MANAGED close internet
		if clusterDeployType == TKE_DEPLOY_TYPE_MANAGED && !clusterInternet {
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				inErr := tkeService.DeleteClusterEndpointVip(ctx, id)
				if inErr != nil {
					return retryError(inErr)
				}
				return nil
			})
			if err != nil {
				return err
			}
			err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
				status, message, inErr := tkeService.DescribeClusterEndpointVipStatus(ctx, id)
				if inErr != nil {
					return retryError(inErr)
				}
				if status == TkeInternetStatusDeleting {
					return resource.RetryableError(
						fmt.Errorf("%s close cluster internet endpoint status still is %s", id, status))
				}
				if status == TkeInternetStatusNotfound || status == TkeInternetStatusDeleted || status == TkeInternetStatusCreated {
					return nil
				}
				return resource.NonRetryableError(
					fmt.Errorf("%s close cluster internet endpoint error ,status is %s,message is %s", id, status, message))
			})
			if err != nil {
				return err
			}
		}
		d.SetPartial("cluster_internet")
	}

	if clusterInternet {
		if !d.HasChange("cluster_internet") && d.HasChange("managed_cluster_internet_security_policies") {
			if len(securityPolicies) == 0 {
				return fmt.Errorf("`managed_cluster_internet_security_policies` can not delete or empty once be setted")
			}
			if err := tkeService.ModifyClusterEndpointSP(ctx, id, securityPolicies); err != nil {
				return err
			}
			d.SetPartial("managed_cluster_internet_security_policies")
		}
	} else {
		d.SetPartial("managed_cluster_internet_security_policies")
	}
	d.Partial(false)
	if err := resourceTencentCloudTkeClusterRead(d, meta); err != nil {
		log.Printf("[WARN]%s resource.kubernetes_cluster.read after update fail , %s", logId, err.Error())
	}

	return nil
}

func resourceTencentCloudTkeClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		err := service.DeleteCluster(ctx, d.Id())

		if e, ok := err.(*errors.TencentCloudSDKError); ok {
			if e.GetCode() == "InternalError.ClusterNotFound" {
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
	_, _, err = service.DescribeClusterInstances(ctx, d.Id())

	if err != nil {
		err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
			_, _, err = service.DescribeClusterInstances(ctx, d.Id())
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InvalidParameter.ClusterNotFound" {
					return nil
				}
			}
			if err != nil {
				return retryError(err, InternalError)
			}
			return nil
		})
	}
	return err

}
