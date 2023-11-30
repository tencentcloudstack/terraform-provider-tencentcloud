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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

var importClsFlag = false

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
		"lan_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "LAN IP of the cvm.",
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
			Description:  "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`, `PREPAID` instance will not terminated after cluster deleted, and may not allow to delete before expired.",
		},
		"instance_charge_type_prepaid_period": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      1,
			ValidateFunc: validateAllowedIntValue(CVM_PREPAID_PERIOD),
			Description:  "The tenancy (time unit is month) of the prepaid instance. NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
		},
		"instance_charge_type_prepaid_renew_flag": {
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validateAllowedStringValue(CVM_PREPAID_RENEW_FLAG),
			Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
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
			Description:  "System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `CLOUD_BASIC`, `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.",
		},
		"system_disk_size": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			Default:      50,
			ValidateFunc: validateIntegerInRange(20, 1024),
			Description:  "Volume of system disk in GB. Default is `50`.",
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
						Description:  "Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD` and `CLOUD_HSSD` and `CLOUD_TSSD`.",
					},
					"disk_size": {
						Type:        schema.TypeInt,
						ForceNew:    true,
						Optional:    true,
						Default:     0,
						Description: "Volume of disk in GB. Default is `0`.",
					},
					"snapshot_id": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "Data disk snapshot ID.",
					},
					"encrypt": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Indicates whether to encrypt data disk, default `false`.",
					},
					"kms_key_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "ID of the custom CMK in the format of UUID or `kms-abcd1234`. This parameter is used to encrypt cloud disks.",
					},
					"file_system": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "File system, e.g. `ext3/ext4/xfs`.",
					},
					"auto_format_and_mount": {
						Type:        schema.TypeBool,
						ForceNew:    true,
						Optional:    true,
						Default:     false,
						Description: "Indicate whether to auto format and mount or not. Default is `false`.",
					},
					"mount_target": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "Mount target.",
					},
					"disk_partition": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						Description: "The name of the device or partition to mount.",
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
			Description:  "Charge types for network traffic. Available values include `TRAFFIC_POSTPAID_BY_HOUR`.",
		},
		"internet_max_bandwidth_out": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Max bandwidth of Internet access in Mbps. Default is 0.",
		},
		"bandwidth_package_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.",
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
		"cam_role_name": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "CAM role name authorized to access.",
		},
		"hostname": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
			Description: "The host name of the attached instance. " +
				"Dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. " +
				"Windows example: The length of the name character is [2, 15], letters (capitalization is not restricted), numbers and dashes (-) are allowed, dots (.) are not supported, and not all numbers are allowed. " +
				"Examples of other types (Linux, etc.): The character length is [2, 60], and multiple dots are allowed. There is a segment between the dots. Each segment allows letters (with no limitation on capitalization), numbers and dashes (-).",
		},
		"disaster_recover_group_ids": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			MaxItems:    1,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Disaster recover groups to which a CVM instance belongs. Only support maximum 1.",
		},
		"img_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateImageID,
			Description:  "The valid image id, format of img-xxx.",
		},
		// InstanceAdvancedSettingsOverrides
		"desired_pod_num": {
			Type:     schema.TypeInt,
			ForceNew: true,
			Optional: true,
			Default:  DefaultDesiredPodNum,
			Description: "Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, " +
				"and it override `[globe_]desired_pod_num` for current node. Either all the fields `desired_pod_num` or none.",
		},
		"hpc_cluster_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Id of cvm hpc cluster.",
		},
	}
}

func TkeExistCvmCreateInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"node_role": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validateAllowedStringValue([]string{TKE_ROLE_WORKER, TKE_ROLE_MASTER_ETCD}),
			Description:  "Role of existed node. value:MASTER_ETCD or WORKER.",
		},
		"instances_para": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"instance_ids": {
						Type:        schema.TypeList,
						ForceNew:    true,
						Required:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Description: "Cluster IDs.",
					},
				},
			},
			Description: "Reinstallation parameters of an existing instance.",
		},
		"desired_pod_numbers": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Description: "Custom mode cluster, you can specify the number of pods for each node. corresponding to the existed_instances_para.instance_ids parameter.",
		},
	}
}

func TkeNodePoolGlobalConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"is_scale_in_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Indicates whether to enable scale-in.",
		},
		"expander": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Indicates which scale-out method will be used when there are multiple scaling groups. Valid values: `random` - select a random scaling group, `most-pods` - select the scaling group that can schedule the most pods, `least-waste` - select the scaling group that can ensure the fewest remaining resources after Pod scheduling.",
		},
		"max_concurrent_scale_in": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Max concurrent scale-in volume.",
		},
		"scale_in_delay": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Number of minutes after cluster scale-out when the system starts judging whether to perform scale-in.",
		},
		"scale_in_unneeded_time": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Number of consecutive minutes of idleness after which the node is subject to scale-in.",
		},
		"scale_in_utilization_threshold": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Percentage of node resource usage below which the node is considered to be idle.",
		},
		"ignore_daemon_sets_utilization": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Whether to ignore DaemonSet pods by default when calculating resource usage.",
		},
		"skip_nodes_with_local_storage": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "During scale-in, ignore nodes with local storage pods.",
		},
		"skip_nodes_with_system_pods": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "During scale-in, ignore nodes with pods in the kube-system namespace that are not managed by DaemonSet.",
		},
	}
}

func resourceTencentCloudTkeCluster() *schema.Resource {
	schemaBody := map[string]*schema.Schema{
		"cluster_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the cluster.",
		},
		"cluster_desc": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the cluster.",
		},
		"cluster_os": {
			Type:     schema.TypeString,
			ForceNew: true,
			Optional: true,
			Default:  TKE_CLUSTER_OS_LINUX24,
			Description: "Operating system of the cluster, the available values include: '" + strings.Join(TKE_CLUSTER_OS, "','") +
				"'. Default is '" + TKE_CLUSTER_OS_LINUX24 + "'.",
		},
		"cluster_subnet_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Subnet ID of the cluster, such as: subnet-b3p7d7q5.",
		},

		"cluster_os_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      TKE_CLUSTER_OS_TYPE_GENERAL,
			ValidateFunc: validateAllowedStringValue(TKE_CLUSTER_OS_TYPES),
			Description: "Image type of the cluster os, the available values include: '" + strings.Join(TKE_CLUSTER_OS_TYPES, "','") +
				"'. Default is '" + TKE_CLUSTER_OS_TYPE_GENERAL + "'.",
		},
		"container_runtime": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      TKE_RUNTIME_DOCKER,
			ValidateFunc: validateAllowedStringValue(TKE_RUNTIMES),
			Description: "Runtime type of the cluster, the available values include: 'docker' and 'containerd'." +
				"The Kubernetes v1.24 has removed dockershim, so please use containerd in v1.24 or higher." +
				"Default is 'docker'.",
		},
		"cluster_deploy_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      TKE_DEPLOY_TYPE_MANAGED,
			ValidateFunc: validateAllowedStringValue(TKE_DEPLOY_TYPES),
			Description:  "Deployment type of the cluster, the available values include: 'MANAGED_CLUSTER' and 'INDEPENDENT_CLUSTER'. Default is 'MANAGED_CLUSTER'.",
		},
		"cluster_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Version of the cluster. Use `tencentcloud_kubernetes_available_cluster_versions` to get the upgradable cluster version.",
		},
		"upgrade_instances_follow_cluster": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicates whether upgrade all instances when cluster_version change. Default is false.",
		},
		"cluster_ipvs": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     true,
			Description: "Indicates whether `ipvs` is enabled. Default is true. False means `iptables` is enabled.",
		},
		"cluster_as_enabled": {
			Type:        schema.TypeBool,
			Computed:    true,
			Deprecated:  "This argument is deprecated because the TKE auto-scaling group was no longer available.",
			Description: "Indicates whether to enable cluster node auto scaling. Default is false.",
		},
		"cluster_level": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Specify cluster level, valid for managed cluster, use data source `tencentcloud_kubernetes_cluster_levels` to query available levels. Available value examples `L5`, `L20`, `L50`, `L100`, etc.",
		},
		"auto_upgrade_cluster_level": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether the cluster level auto upgraded, valid for managed cluster.",
		},
		"acquire_cluster_admin_role": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If set to true, it will acquire the ClusterRole tke:admin. NOTE: this arguments cannot revoke to `false` after acquired.",
		},
		"node_pool_global_config": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: TkeNodePoolGlobalConfig(),
			},
			Description: "Global config effective for all node pools.",
		},
		"cluster_extra_args": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"kube_apiserver": {
						Type:        schema.TypeList,
						ForceNew:    true,
						Optional:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Description: "The customized parameters for kube-apiserver.",
					},
					"kube_controller_manager": {
						Type:        schema.TypeList,
						ForceNew:    true,
						Optional:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Description: "The customized parameters for kube-controller-manager.",
					},
					"kube_scheduler": {
						Type:        schema.TypeList,
						ForceNew:    true,
						Optional:    true,
						Elem:        &schema.Schema{Type: schema.TypeString},
						Description: "The customized parameters for kube-scheduler.",
					},
				},
			},
			Description: "Customized parameters for master component,such as kube-apiserver, kube-controller-manager, kube-scheduler.",
		},
		"node_name_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      "lan-ip",
			Description:  "Node name type of Cluster, the available values include: 'lan-ip' and 'hostname', Default is 'lan-ip'.",
			ValidateFunc: validateAllowedStringValue(TKE_CLUSTER_NODE_NAME_TYPE),
		},
		"network_type": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			Default:      "GR",
			ValidateFunc: validateAllowedStringValue(TKE_CLUSTER_NETWORK_TYPE),
			Description:  "Cluster network type, the available values include: 'GR' and 'VPC-CNI' and 'CiliumOverlay'. Default is GR.",
		},
		"enable_customized_pod_cidr": {
			Type: schema.TypeBool,
			//ForceNew:    true,
			Optional:    true,
			Default:     false,
			Description: "Whether to enable the custom mode of node podCIDR size. Default is false.",
		},
		"base_pod_num": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Description: "The number of basic pods. valid when enable_customized_pod_cidr=true.",
		},
		"is_non_static_ip_mode": {
			Type:        schema.TypeBool,
			ForceNew:    true,
			Optional:    true,
			Default:     false,
			Description: "Indicates whether non-static ip mode is enabled. Default is false.",
		},
		"deletion_protection": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicates whether cluster deletion protection is enabled. Default is false.",
		},
		"kube_proxy_mode": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
			Description: "Cluster kube-proxy mode, the available values include: 'kube-proxy-bpf'. Default is not set." +
				"When set to kube-proxy-bpf, cluster version greater than 1.14 and with Tencent Linux 2.4 is required.",
		},
		"vpc_cni_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validateAllowedStringValue([]string{"tke-route-eni", "tke-direct-eni"}),
			Description: "Distinguish between shared network card multi-IP mode and independent network card mode. " +
				"Fill in `tke-route-eni` for shared network card multi-IP mode and `tke-direct-eni` for independent network card mode. " +
				"The default is shared network card mode. When it is necessary to turn off the vpc-cni container network capability, both `eni_subnet_ids` and `vpc_cni_type` must be set to empty.",
		},
		"vpc_id": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validateStringLengthInRange(4, 100),
			Description:  "Vpc Id of the cluster.",
		},
		"cluster_internet": {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
			Description: "Open internet access or not." +
				" If this field is set 'true', the field below `worker_config` must be set." +
				" Because only cluster with node is allowed enable access endpoint. You may open it through `tencentcloud_kubernetes_cluster_endpoint`.",
		},
		"cluster_internet_domain": {
			Type:     schema.TypeString,
			Optional: true,
			Description: "Domain name for cluster Kube-apiserver internet access." +
				" Be careful if you modify value of this parameter, the cluster_external_endpoint value may be changed automatically too.",
		},
		"cluster_intranet": {
			Type:     schema.TypeBool,
			Default:  false,
			Optional: true,
			Description: "Open intranet access or not." +
				" If this field is set 'true', the field below `worker_config` must be set." +
				" Because only cluster with node is allowed enable access endpoint. You may open it through `tencentcloud_kubernetes_cluster_endpoint`.",
		},
		"cluster_intranet_domain": {
			Type:     schema.TypeString,
			Optional: true,
			Description: "Domain name for cluster Kube-apiserver intranet access." +
				" Be careful if you modify value of this parameter, the pgw_endpoint value may be changed automatically too.",
		},
		"cluster_internet_security_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify security group, NOTE: This argument must not be empty if cluster internet enabled.",
		},
		"managed_cluster_internet_security_policies": {
			Type:       schema.TypeList,
			Optional:   true,
			Elem:       &schema.Schema{Type: schema.TypeString},
			Deprecated: "this argument was deprecated, use `cluster_internet_security_group` instead.",
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
			Optional:    true,
			Description: "Project ID, default value is 0.",
		},
		"cluster_cidr": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.",
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				value := v.(string)
				if value == "" {
					return
				}
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
				if !strings.HasPrefix(value, "9.") && !strings.HasPrefix(value, "10.") && !strings.HasPrefix(value, "192.168.") && !strings.HasPrefix(value, "172.") {
					errors = append(errors, fmt.Errorf("%q must in 9. | 10. | 192.168. | 172.[16-31]", k))
					return
				}

				if strings.HasPrefix(value, "172.") {
					nextNo := strings.Split(value, ".")[1]
					no, _ := strconv.ParseInt(nextNo, 10, 64)
					if no < 16 || no > 31 {
						errors = append(errors, fmt.Errorf("%q must in 9.0 | 10. | 192.168. | 172.[16-31]", k))
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
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Default:     256,
			Description: "The maximum number of Pods per node in the cluster. Default is 256. The minimum value is 4. When its power unequal to 2, it will round upward to the closest power of 2.",
		},
		"cluster_max_service_num": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Default:     256,
			Description: "The maximum number of services in the cluster. Default is 256. The range is from 32 to 32768. When its power unequal to 2, it will round upward to the closest power of 2.",
		},
		"service_cidr": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "A network address block of the service. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.",
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				value := v.(string)
				if value == "" {
					return
				}
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
				if !strings.HasPrefix(value, "9.") && !strings.HasPrefix(value, "10.") && !strings.HasPrefix(value, "192.168.") && !strings.HasPrefix(value, "172.") {
					errors = append(errors, fmt.Errorf("%q must in 9. | 10. | 192.168. | 172.[16-31]", k))
					return
				}

				if strings.HasPrefix(value, "172.") {
					nextNo := strings.Split(value, ".")[1]
					no, _ := strconv.ParseInt(nextNo, 10, 64)
					if no < 16 || no > 31 {
						errors = append(errors, fmt.Errorf("%q must in 9. | 10. | 192.168. | 172.[16-31]", k))
						return
					}
				}
				return
			},
		},
		"eni_subnet_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Description: "Subnet Ids for cluster with VPC-CNI network mode." +
				" This field can only set when field `network_type` is 'VPC-CNI'." +
				" `eni_subnet_ids` can not empty once be set.",
		},
		"claim_expired_seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
			Description: "Claim expired seconds to recycle ENI." +
				" This field can only set when field `network_type` is 'VPC-CNI'." +
				" `claim_expired_seconds` must greater or equal than 300 and less than 15768000.",
			ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
				value := v.(int)
				if value < 300 || value > 15768000 {
					errors = append(errors, fmt.Errorf("%q must greater or equal than 300 and less than 15768000", k))
					return
				}
				return
			},
		},
		"master_config": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: TkeCvmCreateInfo(),
			},
			Description: "Deploy the machine configuration information of the 'MASTER_ETCD' service, and create <=7 units for common users.",
		},
		"worker_config": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: TkeCvmCreateInfo(),
			},
			Description: "Deploy the machine configuration information of the 'WORKER' service, and create <=20 units for common users. The other 'WORK' service are added by 'tencentcloud_kubernetes_worker'.",
		},
		"exist_instance": {
			Type:     schema.TypeList,
			ForceNew: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: TkeExistCvmCreateInfo(),
			},
			Description: "create tke cluster by existed instances.",
		},
		"auth_options": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"use_tke_default": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "If set to `true`, the issuer and jwks_uri will be generated automatically by tke, please do not set issuer and jwks_uri, and they will be ignored.",
					},
					"jwks_uri": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specify service-account-jwks-uri. If use_tke_default is set to `true`, please do not set this field, it will be ignored anyway.",
					},
					"issuer": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specify service-account-issuer. If use_tke_default is set to `true`, please do not set this field, it will be ignored anyway.",
					},
					"auto_create_discovery_anonymous_auth": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "If set to `true`, the rbac rule will be created automatically which allow anonymous user to access '/.well-known/openid-configuration' and '/openid/v1/jwks'.",
					},
				},
			},
			Description: "Specify cluster authentication configuration. Only available for managed cluster and `cluster_version` >= 1.20.",
		},
		"extension_addon": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Information of the add-on to be installed.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Add-on name.",
					},
					"param": {
						Type:             schema.TypeString,
						Required:         true,
						DiffSuppressFunc: helper.DiffSupressJSON,
						Description:      "Parameter of the add-on resource object in JSON string format, please check the example at the top of page for reference.",
					},
				},
			},
		},
		"log_agent": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specify cluster log agent config.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Whether the log agent enabled.",
					},
					"kubelet_root_dir": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Kubelet root directory as the literal.",
					},
				},
			},
		},
		"event_persistence": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specify cluster Event Persistence config. NOTE: Please make sure your TKE CamRole have permission to access CLS service.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Specify weather the Event Persistence enabled.",
					},
					"log_set_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specify id of existing CLS log set, or auto create a new set by leave it empty.",
					},
					"topic_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specify id of existing CLS log topic, or auto create a new topic by leave it empty.",
					},
					"delete_event_log_and_topic": {
						Type:     schema.TypeBool,
						Optional: true,
						Description: "when you want to close the cluster event persistence or delete the cluster, you can use this parameter to determine " +
							"whether the event persistence log set and topic created by default will be deleted.",
					},
				},
			},
		},
		"cluster_audit": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Specify Cluster Audit config. NOTE: Please make sure your TKE CamRole have permission to access CLS service.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: "Specify weather the Cluster Audit enabled. NOTE: Enable Cluster Audit will also auto install Log Agent.",
					},
					"log_set_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specify id of existing CLS log set, or auto create a new set by leave it empty.",
					},
					"topic_id": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Specify id of existing CLS log topic, or auto create a new topic by leave it empty.",
					},
					"delete_audit_log_and_topic": {
						Type:     schema.TypeBool,
						Optional: true,
						Description: "when you want to close the cluster audit log or delete the cluster, you can use " +
							"this parameter to determine whether the audit log set and topic created by default will" +
							" be deleted.",
					},
				},
			},
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
		//advanced instance setting
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Description: "Labels of tke cluster nodes.",
		},
		"unschedulable": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if new == "0" && old == "" {
					return true
				} else {
					return old == new
				}
			},
			Default:     0,
			Description: "Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.",
		},
		"mount_target": {
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Description: "Mount target. Default is not mounting.",
		},
		"globe_desired_pod_num": {
			Type:        schema.TypeInt,
			ForceNew:    true,
			Optional:    true,
			Description: "Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, and it takes effect for all nodes.",
		},
		"docker_graph_path": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "/var/lib/docker",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				if new == "/var/lib/docker" && old == "" || old == "/var/lib/docker" && new == "" {
					return true
				} else {
					return old == new
				}
			},
			Description: "Docker graph path. Default is `/var/lib/docker`.",
		},
		"extra_args": {
			Type:        schema.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Description: "Custom parameter information related to the node.",
		},
		"runtime_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Container Runtime version.",
		},

		"kube_config": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Kubernetes config.",
		},
		"kube_config_intranet": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Kubernetes config of private network.",
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
		Importer: &schema.ResourceImporter{
			//State: schema.ImportStatePassthrough,
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				importClsFlag = true
				err := resourceTencentCloudTkeClusterRead(d, m)
				if err != nil {
					return nil, fmt.Errorf("failed to import resource")
				}

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: schemaBody,
	}
}

func tkeGetCvmRunInstancesPara(dMap map[string]interface{}, meta interface{},
	vpcId string, projectId int64) (cvmJson string, count int64, errRet error) {

	request := cvm.NewRunInstancesRequest()

	var place cvm.Placement
	request.Placement = &place

	place.ProjectId = &projectId

	if v, ok := dMap["availability_zone"]; ok {
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

	if v, ok := dMap["cam_role_name"]; ok {
		request.CamRoleName = helper.String(v.(string))
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
				encrypt    = value["encrypt"].(bool)
				kmsKeyId   = value["kms_key_id"].(string)
				dataDisk   = cvm.DataDisk{
					DiskType: &diskType,
				}
			)
			if diskSize > 0 {
				dataDisk.DiskSize = &diskSize
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			if encrypt {
				dataDisk.Encrypt = &encrypt
			}
			if kmsKeyId != "" {
				dataDisk.KmsKeyId = &kmsKeyId
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

	if v, ok := dMap["bandwidth_package_id"]; ok {
		if v.(string) != "" {
			request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
		}
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

	if v, ok := dMap["disaster_recover_group_ids"]; ok {
		disasterGroups := v.([]interface{})
		request.DisasterRecoverGroupIds = make([]*string, 0, len(disasterGroups))
		for i := range disasterGroups {
			disasterGroup := disasterGroups[i].(string)
			request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &disasterGroup)
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

	if v, ok := dMap["hostname"]; ok {
		hostname := v.(string)
		if hostname != "" {
			request.HostName = &hostname
		}
	}

	if v, ok := dMap["img_id"]; ok && v.(string) != "" {
		request.ImageId = helper.String(v.(string))
	}

	if v, ok := dMap["hpc_cluster_id"]; ok && v.(string) != "" {
		request.HpcClusterId = helper.String(v.(string))
	}

	cvmJson = request.ToJsonString()

	cvmJson = strings.Replace(cvmJson, `"Password":"",`, "", -1)

	return
}

func tkeGetCvmExistInstancesPara(dMap map[string]interface{}) (tke.ExistedInstancesForNode, error) {

	inst := tke.ExistedInstancesForNode{}

	if temp, ok := dMap["instances_para"]; ok {
		paras := temp.([]interface{})
		if len(paras) > 0 {
			paraMap := paras[0].(map[string]interface{})
			instanceIds := paraMap["instance_ids"].([]interface{})
			inst.ExistedInstancesPara = &tke.ExistedInstancesPara{}
			inst.ExistedInstancesPara.InstanceIds = make([]*string, 0)
			for _, v := range instanceIds {
				inst.ExistedInstancesPara.InstanceIds = append(inst.ExistedInstancesPara.InstanceIds, helper.String(v.(string)))
			}
		}
	}
	if temp, ok := dMap["desired_pod_numbers"]; ok {
		inst.DesiredPodNumbers = make([]*int64, 0)
		podNums := temp.([]interface{})
		for _, v := range podNums {
			inst.DesiredPodNumbers = append(inst.DesiredPodNumbers, helper.Int64(int64(v.(int))))
		}
	}
	if temp, ok := dMap["node_role"]; ok {
		nodeRole := temp.(string)
		inst.NodeRole = &nodeRole
	}

	return inst, nil
}

func tkeGetNodePoolGlobalConfig(d *schema.ResourceData) *tke.ModifyClusterAsGroupOptionAttributeRequest {
	request := tke.NewModifyClusterAsGroupOptionAttributeRequest()
	request.ClusterId = helper.String(d.Id())

	clusterAsGroupOption := &tke.ClusterAsGroupOption{}
	if v, ok := d.GetOkExists("node_pool_global_config.0.is_scale_in_enabled"); ok {
		clusterAsGroupOption.IsScaleDownEnabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.expander"); ok {
		clusterAsGroupOption.Expander = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.max_concurrent_scale_in"); ok {
		clusterAsGroupOption.MaxEmptyBulkDelete = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.scale_in_delay"); ok {
		clusterAsGroupOption.ScaleDownDelay = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.scale_in_unneeded_time"); ok {
		clusterAsGroupOption.ScaleDownUnneededTime = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.scale_in_utilization_threshold"); ok {
		clusterAsGroupOption.ScaleDownUtilizationThreshold = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.ignore_daemon_sets_utilization"); ok {
		clusterAsGroupOption.IgnoreDaemonSetsUtilization = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.skip_nodes_with_local_storage"); ok {
		clusterAsGroupOption.SkipNodesWithLocalStorage = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("node_pool_global_config.0.skip_nodes_with_system_pods"); ok {
		clusterAsGroupOption.SkipNodesWithSystemPods = helper.Bool(v.(bool))
	}

	request.ClusterAsGroupOption = clusterAsGroupOption
	return request
}

func tkeGetAuthOptions(d *schema.ResourceData) *tke.ModifyClusterAuthenticationOptionsRequest {
	raw, ok := d.GetOk("auth_options")
	options := raw.([]interface{})

	request := tke.NewModifyClusterAuthenticationOptionsRequest()
	request.ClusterId = helper.String(d.Id())
	request.ServiceAccounts = &tke.ServiceAccountAuthenticationOptions{
		AutoCreateDiscoveryAnonymousAuth: helper.Bool(false),
	}

	if !ok || len(options) == 0 {
		request.ServiceAccounts.JWKSURI = helper.String("")
		return request
	}

	option := options[0].(map[string]interface{})

	if v, ok := option["auto_create_discovery_anonymous_auth"]; ok {
		request.ServiceAccounts.AutoCreateDiscoveryAnonymousAuth = helper.Bool(v.(bool))
	}

	if v, ok := option["use_tke_default"]; ok && v.(bool) {
		request.ServiceAccounts.UseTKEDefault = helper.Bool(true)
	} else {
		if v, ok := option["issuer"]; ok {
			request.ServiceAccounts.Issuer = helper.String(v.(string))
		}

		if v, ok := option["jwks_uri"]; ok {
			request.ServiceAccounts.JWKSURI = helper.String(v.(string))
		}
	}

	return request
}

// upgradeClusterInstances upgrade instances, upgrade type try seq:major, hot.
func upgradeClusterInstances(tkeService TkeService, ctx context.Context, id string) error {
	// get all available instances for upgrade
	upgradeType := "major"
	instanceIds, err := tkeService.CheckInstancesUpgradeAble(ctx, id, upgradeType)
	if err != nil {
		return err
	}
	if len(instanceIds) == 0 {
		upgradeType = "hot"
		instanceIds, err = tkeService.CheckInstancesUpgradeAble(ctx, id, upgradeType)
		if err != nil {
			return err
		}
	}
	log.Println("instancesIds for upgrade:", instanceIds)
	instNum := len(instanceIds)
	if instNum == 0 {
		return nil
	}

	// upgrade instances
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		inErr := tkeService.UpgradeClusterInstances(ctx, id, upgradeType, instanceIds)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// check update status: upgrade instance one by one, so timeout depend on instance number.
	timeout := readRetryTimeout * time.Duration(instNum)
	err = resource.Retry(timeout, func() *resource.RetryError {
		done, inErr := tkeService.GetUpgradeInstanceResult(ctx, id)
		if inErr != nil {
			return retryError(inErr)
		}
		if done {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("cluster %s, retry...", id))
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudTkeClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_kubernetes_cluster.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		basic                        ClusterBasicSetting
		advanced                     ClusterAdvancedSettings
		cvms                         RunInstancesForNode
		iAdvanced                    InstanceAdvancedSettings
		iDiskMountSettings           []*tke.InstanceDataDiskMountSetting
		cidrSet                      ClusterCidrSettings
		extensionAddons              []*tke.ExtensionAddon
		clusterInternet              = d.Get("cluster_internet").(bool)
		clusterIntranet              = d.Get("cluster_intranet").(bool)
		intranetSubnetId             = d.Get("cluster_intranet_subnet_id").(string)
		clusterInternetSecurityGroup = d.Get("cluster_internet_security_group").(string)
		clusterInternetDomain        = d.Get("cluster_internet_domain").(string)
		clusterIntranetDomain        = d.Get("cluster_intranet_domain").(string)
	)

	clusterDeployType := d.Get("cluster_deploy_type").(string)

	if clusterIntranet && intranetSubnetId == "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` must set when `cluster_intranet` is true")
	}
	if !clusterIntranet && intranetSubnetId != "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` can only set when `cluster_intranet` is true")
	}

	vpcId := d.Get("vpc_id").(string)
	if vpcId != "" {
		basic.VpcId = vpcId
	}

	basic.ProjectId = int64(d.Get("project_id").(int))

	cluster_os := d.Get("cluster_os").(string)

	if v, ok := tkeClusterOsMap[cluster_os]; ok {
		basic.ClusterOs = v
	} else {
		basic.ClusterOs = cluster_os
	}

	if tkeClusterOsMap[cluster_os] != "" {
		basic.ClusterOs = tkeClusterOsMap[cluster_os]
	} else {
		basic.ClusterOs = cluster_os
	}

	basic.ClusterOsType = d.Get("cluster_os_type").(string)
	basic.SubnetId = d.Get("cluster_subnet_id").(string)
	basic.ClusterVersion = d.Get("cluster_version").(string)
	if v, ok := d.GetOk("cluster_name"); ok {
		basic.ClusterName = v.(string)
	}
	if v, ok := d.GetOk("cluster_desc"); ok {
		basic.ClusterDescription = v.(string)
	}

	if v, ok := d.GetOk("cluster_level"); ok {
		basic.ClusterLevel = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_upgrade_cluster_level"); ok {
		basic.AutoUpgradeClusterLevel = helper.Bool(v.(bool))
	}

	advanced.Ipvs = d.Get("cluster_ipvs").(bool)
	advanced.AsEnabled = d.Get("cluster_as_enabled").(bool)
	advanced.ContainerRuntime = d.Get("container_runtime").(string)
	advanced.NodeNameType = d.Get("node_name_type").(string)
	advanced.NetworkType = d.Get("network_type").(string)
	advanced.IsNonStaticIpMode = d.Get("is_non_static_ip_mode").(bool)
	advanced.DeletionProtection = d.Get("deletion_protection").(bool)
	advanced.KubeProxyMode = d.Get("kube_proxy_mode").(string)
	advanced.EnableCustomizedPodCIDR = d.Get("enable_customized_pod_cidr").(bool)

	if advanced.NetworkType == TKE_CLUSTER_NETWORK_TYPE_VPC_CNI {
		if v, ok := d.GetOk("vpc_cni_type"); ok {
			advanced.VpcCniType = v.(string)
		} else {
			advanced.VpcCniType = "tke-route-eni"
		}
	}
	if v, ok := d.GetOk("base_pod_num"); ok {
		advanced.BasePodNumber = int64(v.(int))
	}

	if extraArgs, ok := d.GetOk("cluster_extra_args"); ok {
		extraArgList := extraArgs.([]interface{})
		for index := range extraArgList {
			extraArg := extraArgList[index].(map[string]interface{})
			if apiserverArgs, exist := extraArg["kube_apiserver"]; exist {
				args := apiserverArgs.([]interface{})
				for index := range args {
					advanced.ExtraArgs.KubeAPIServer = append(advanced.ExtraArgs.KubeAPIServer, args[index].(string))
				}
			}
			if cmArgs, exist := extraArg["kube_controller_manager"]; exist {
				args := cmArgs.([]interface{})
				for index := range args {
					advanced.ExtraArgs.KubeControllerManager = append(advanced.ExtraArgs.KubeControllerManager, args[index].(string))
				}
			}
			if schedulerArgs, exist := extraArg["kube_scheduler"]; exist {
				args := schedulerArgs.([]interface{})
				for index := range args {
					advanced.ExtraArgs.KubeScheduler = append(advanced.ExtraArgs.KubeScheduler, args[index].(string))
				}
			}
		}
	}
	cidrSet.ClusterCidr = d.Get("cluster_cidr").(string)
	cidrSet.IgnoreClusterCidrConflict = d.Get("ignore_cluster_cidr_conflict").(bool)
	cidrSet.MaxClusterServiceNum = int64(d.Get("cluster_max_service_num").(int))
	cidrSet.MaxNodePodNum = int64(d.Get("cluster_max_pod_num").(int))
	cidrSet.ServiceCIDR = d.Get("service_cidr").(string)

	if ClaimExpiredSeconds, ok := d.GetOk("claim_expired_seconds"); ok {
		cidrSet.ClaimExpiredSeconds = int64(ClaimExpiredSeconds.(int))
	} else {
		cidrSet.ClaimExpiredSeconds = int64(300)

		if err := d.Set("claim_expired_seconds", 300); err != nil {
			return fmt.Errorf("error setting claim_expired_seconds: %s", err)
		}
	}

	if advanced.NetworkType == TKE_CLUSTER_NETWORK_TYPE_VPC_CNI {
		// VPC-CNI cluster need to set eni subnet and service cidr.
		eniSubnetIdList := d.Get("eni_subnet_ids").([]interface{})
		for index := range eniSubnetIdList {
			subnetId := eniSubnetIdList[index].(string)
			cidrSet.EniSubnetIds = append(cidrSet.EniSubnetIds, subnetId)
		}
		if cidrSet.ServiceCIDR == "" || len(cidrSet.EniSubnetIds) == 0 {
			return fmt.Errorf("`service_cidr` must be set and `eni_subnet_ids` must be set when cluster `network_type` is VPC-CNI.")
		}
	} else {
		// GR cluster
		if cidrSet.ClusterCidr == "" {
			return fmt.Errorf("`cluster_cidr` must be set when cluster `network_type` is GR")
		}
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

		if advanced.NetworkType == TKE_CLUSTER_NETWORK_TYPE_CILIUM_OVERLAY && d.Get("cluster_subnet_id").(string) == "" {
			return fmt.Errorf("`cluster_subnet_id` must be set ")
		}
	}

	if version, ok := d.GetOk("runtime_version"); ok {
		advanced.RuntimeVersion = version.(string)
	}

	overrideSettings := OverrideSettings{
		Master: make([]tke.InstanceAdvancedSettings, 0),
		Work:   make([]tke.InstanceAdvancedSettings, 0),
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

			if v, ok := master["desired_pod_num"]; ok {
				dpNum := int64(v.(int))
				if dpNum != DefaultDesiredPodNum {
					overrideSettings.Master = append(overrideSettings.Master, tke.InstanceAdvancedSettings{DesiredPodNumber: helper.Int64(dpNum)})
				}
			}
		}
		if masterCount < 3 {
			return fmt.Errorf("if `cluster_deploy_type` is `TKE_DEPLOY_TYPE_INDEPENDENT` len(master_config) should >=3")
		}
	} else if clusterDeployType == TKE_DEPLOY_TYPE_INDEPENDENT {
		return fmt.Errorf("if `cluster_deploy_type` is `TKE_DEPLOY_TYPE_INDEPENDENT` , You need define the master yourself")
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

			if v, ok := worker["desired_pod_num"]; ok {
				dpNum := int64(v.(int))
				if dpNum != DefaultDesiredPodNum {
					overrideSettings.Work = append(overrideSettings.Work, tke.InstanceAdvancedSettings{DesiredPodNumber: helper.Int64(dpNum)})
				}
			}

			if v, ok := worker["data_disk"]; ok {
				var (
					instanceType = worker["instance_type"].(string)
					zone         = worker["availability_zone"].(string)
				)
				iDiskMountSetting := &tke.InstanceDataDiskMountSetting{
					InstanceType: &instanceType,
					Zone:         &zone,
				}

				diskList := v.([]interface{})
				for _, d := range diskList {
					var (
						disk               = d.(map[string]interface{})
						diskType           = disk["disk_type"].(string)
						diskSize           = int64(disk["disk_size"].(int))
						fileSystem         = disk["file_system"].(string)
						autoFormatAndMount = disk["auto_format_and_mount"].(bool)
						mountTarget        = disk["mount_target"].(string)
						diskPartition      = disk["disk_partition"].(string)
					)

					dataDisk := &tke.DataDisk{
						DiskType:           &diskType,
						DiskSize:           &diskSize,
						AutoFormatAndMount: &autoFormatAndMount,
					}

					if fileSystem != "" {
						dataDisk.FileSystem = &fileSystem
					}

					if mountTarget != "" {
						dataDisk.MountTarget = &mountTarget
					}

					if diskPartition != "" {
						dataDisk.DiskPartition = &diskPartition
					}

					iDiskMountSetting.DataDisks = append(iDiskMountSetting.DataDisks, dataDisk)
				}

				iDiskMountSettings = append(iDiskMountSettings, iDiskMountSetting)
			}
		}
	}

	tags := helper.GetTags(d, "tags")

	iAdvanced.Labels = GetTkeLabels(d, "labels")

	if temp, ok := d.GetOk("extra_args"); ok {
		extraArgs := helper.InterfacesStrings(temp.([]interface{}))
		for i := range extraArgs {
			iAdvanced.ExtraArgs.Kubelet = append(iAdvanced.ExtraArgs.Kubelet, &extraArgs[i])
		}
	}
	if temp, ok := d.GetOk("unschedulable"); ok {
		iAdvanced.Unschedulable = int64(temp.(int))
	}
	if temp, ok := d.GetOk("docker_graph_path"); ok {
		iAdvanced.DockerGraphPath = temp.(string)
	} else {
		iAdvanced.DockerGraphPath = "/var/lib/docker"
	}
	if temp, ok := d.GetOk("mount_target"); ok {
		iAdvanced.MountTarget = temp.(string)
	}
	if temp, ok := d.GetOk("globe_desired_pod_num"); ok {
		iAdvanced.DesiredPodNum = int64(temp.(int))
	}

	// ExistedInstancesForNode
	existInstances := make([]*tke.ExistedInstancesForNode, 0)
	if instances, ok := d.GetOk("exist_instance"); ok {
		instanceList := instances.([]interface{})
		for index := range instanceList {
			instance := instanceList[index].(map[string]interface{})
			existedInstance, _ := tkeGetCvmExistInstancesPara(instance)
			existInstances = append(existInstances, &existedInstance)
		}
	}

	// RunInstancesForNode（master_config+worker_config) 和 ExistedInstancesForNode 不能同时存在
	if len(cvms.Master)+len(cvms.Work) > 0 && len(existInstances) > 0 {
		return fmt.Errorf("master_config+worker_config and exist_instance can not exist at the same time")
	}

	if v, ok := d.GetOk("extension_addon"); ok {
		for _, i := range v.([]interface{}) {
			dMap := i.(map[string]interface{})
			name := dMap["name"].(string)
			param := dMap["param"].(string)
			addon := &tke.ExtensionAddon{
				AddonName:  helper.String(name),
				AddonParam: helper.String(param),
			}
			extensionAddons = append(extensionAddons, addon)
		}
	}

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}
	id, err := service.CreateCluster(ctx, basic, advanced, cvms, iAdvanced, cidrSet, tags, existInstances, &overrideSettings, iDiskMountSettings, extensionAddons)
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

	err = service.CheckOneOfClusterNodeReady(ctx, d.Id(), clusterInternet || clusterIntranet)

	if err != nil {
		return err
	}

	//intranet
	if clusterIntranet {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := service.CreateClusterEndpoint(ctx, id, intranetSubnetId, clusterInternetSecurityGroup, false, clusterIntranetDomain, "")
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			status, message, inErr := service.DescribeClusterEndpointStatus(ctx, id, false)
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

	if clusterInternet {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := service.CreateClusterEndpoint(ctx, id, "", clusterInternetSecurityGroup, true, clusterInternetDomain, "")
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			status, message, inErr := service.DescribeClusterEndpointStatus(ctx, id, true)
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

	//Modify node pool global config
	if _, ok := d.GetOk("node_pool_global_config"); ok {
		request := tkeGetNodePoolGlobalConfig(d)
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := service.ModifyClusterNodePoolGlobalConfig(ctx, request)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if v, ok := d.GetOk("acquire_cluster_admin_role"); ok && v.(bool) {
		err := service.AcquireClusterAdminRole(ctx, id)
		if err != nil {
			return err
		}
	}

	if _, ok := d.GetOk("auth_options"); ok {
		request := tkeGetAuthOptions(d)
		if err := service.ModifyClusterAuthenticationOptions(ctx, request); err != nil {
			return err
		}
	}

	if v, ok := helper.InterfacesHeadMap(d, "log_agent"); ok {
		enabled := v["enabled"].(bool)
		rootDir := v["kubelet_root_dir"].(string)

		if enabled {
			err := service.SwitchLogAgent(ctx, id, rootDir, enabled)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := helper.InterfacesHeadMap(d, "event_persistence"); ok {
		enabled := v["enabled"].(bool)
		logSetId := v["log_set_id"].(string)
		topicId := v["topic_id"].(string)
		if enabled {
			err := service.SwitchEventPersistence(ctx, id, logSetId, topicId, enabled, false)
			if err != nil {
				return err
			}
		}
	}

	if v, ok := helper.InterfacesHeadMap(d, "cluster_audit"); ok {
		enabled := v["enabled"].(bool)
		logSetId := v["log_set_id"].(string)
		topicId := v["topic_id"].(string)
		if enabled {
			err := service.SwitchClusterAudit(ctx, id, logSetId, topicId, enabled, false)
			if err != nil {
				return err
			}
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
	cvmService := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

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
	// 兼容旧的 cluster_os 的 key, 由于 cluster_os有默认值，所以不大可能为空
	oldOs := d.Get("cluster_os").(string)
	newOs := tkeToShowClusterOs(info.ClusterOs)

	if (oldOs == TkeClusterOsCentOS76 && newOs == TKE_CLUSTER_OS_CENTOS76) ||
		(oldOs == TkeClusterOsUbuntu18 && newOs == TKE_CLUSTER_OS_UBUNTU18) {
		newOs = oldOs
	}

	_ = d.Set("cluster_name", info.ClusterName)
	_ = d.Set("cluster_desc", info.ClusterDescription)
	_ = d.Set("cluster_os", newOs)
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
	_ = d.Set("deletion_protection", info.DeletionProtection)
	_ = d.Set("cluster_level", info.ClusterLevel)
	_ = d.Set("vpc_cni_type", info.VpcCniType)
	_ = d.Set("eni_subnet_ids", info.EniSubnetIds)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(info.Property), &data)
	if err != nil {
		return fmt.Errorf("error:%v", err)
	}

	if importClsFlag && info.DeployType == TKE_DEPLOY_TYPE_INDEPENDENT {
		var masters []InstanceInfo
		var errRet error
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			masters, _, errRet = service.DescribeClusterInstancesByRole(ctx, d.Id(), "MASTER_OR_ETCD")
			if e, ok := errRet.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() == "InternalError.ClusterNotFound" {
					return nil
				}
			}
			if errRet != nil {
				return resource.RetryableError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		var instances []*cvm.Instance
		instanceIds := make([]*string, 0)
		for _, instance := range masters {
			instanceIds = append(instanceIds, helper.String(instance.InstanceId))
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			instances, errRet = cvmService.DescribeInstanceByFilter(ctx, instanceIds, nil)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			return nil
		})
		if err != nil {
			return err
		}

		instanceList := make([]interface{}, 0, len(instances))
		for _, instance := range instances {
			mapping := map[string]interface{}{
				"count":                               1,
				"instance_charge_type_prepaid_period": 1,
				"instance_type":                       helper.PString(instance.InstanceType),
				"subnet_id":                           helper.PString(instance.VirtualPrivateCloud.SubnetId),
				"availability_zone":                   helper.PString(instance.Placement.Zone),
				"instance_name":                       helper.PString(instance.InstanceName),
				"instance_charge_type":                helper.PString(instance.InstanceChargeType),
				"system_disk_type":                    helper.PString(instance.SystemDisk.DiskType),
				"system_disk_size":                    helper.PInt64(instance.SystemDisk.DiskSize),
				"internet_charge_type":                helper.PString(instance.InternetAccessible.InternetChargeType),
				"bandwidth_package_id":                helper.PString(instance.InternetAccessible.BandwidthPackageId),
				"internet_max_bandwidth_out":          helper.PInt64(instance.InternetAccessible.InternetMaxBandwidthOut),
				"security_group_ids":                  helper.StringsInterfaces(instance.SecurityGroupIds),
				"img_id":                              helper.PString(instance.ImageId),
			}

			if instance.RenewFlag != nil && helper.PString(instance.InstanceChargeType) == "PREPAID" {
				mapping["instance_charge_type_prepaid_renew_flag"] = helper.PString(instance.RenewFlag)
			} else {
				mapping["instance_charge_type_prepaid_renew_flag"] = ""
			}
			if helper.PInt64(instance.InternetAccessible.InternetMaxBandwidthOut) > 0 {
				mapping["public_ip_assigned"] = true
			}

			if instance.CamRoleName != nil {
				mapping["cam_role_name"] = instance.CamRoleName
			}
			if instance.LoginSettings != nil {
				if instance.LoginSettings.KeyIds != nil && len(instance.LoginSettings.KeyIds) > 0 {
					mapping["key_ids"] = helper.StringsInterfaces(instance.LoginSettings.KeyIds)
				}
				if instance.LoginSettings.Password != nil {
					mapping["password"] = helper.PString(instance.LoginSettings.Password)
				}
			}
			if instance.DisasterRecoverGroupId != nil && helper.PString(instance.DisasterRecoverGroupId) != "" {
				mapping["disaster_recover_group_ids"] = []string{helper.PString(instance.DisasterRecoverGroupId)}
			}
			if instance.HpcClusterId != nil {
				mapping["hpc_cluster_id"] = helper.PString(instance.HpcClusterId)
			}

			dataDisks := make([]interface{}, 0, len(instance.DataDisks))
			for _, v := range instance.DataDisks {
				dataDisk := map[string]interface{}{
					"disk_type":   helper.PString(v.DiskType),
					"disk_size":   helper.PInt64(v.DiskSize),
					"snapshot_id": helper.PString(v.DiskId),
					"encrypt":     helper.PBool(v.Encrypt),
					"kms_key_id":  helper.PString(v.KmsKeyId),
				}
				dataDisks = append(dataDisks, dataDisk)
			}

			mapping["data_disk"] = dataDisks
			instanceList = append(instanceList, mapping)
		}

		_ = d.Set("master_config", instanceList)
	}

	if importClsFlag {
		networkType, _ := data["NetworkType"].(string)
		_ = d.Set("network_type", networkType)

		nodeNameType, _ := data["NodeNameType"].(string)
		_ = d.Set("node_name_type", nodeNameType)

		enableCustomizedPodCIDR, _ := data["EnableCustomizedPodCIDR"].(bool)
		_ = d.Set("enable_customized_pod_cidr", enableCustomizedPodCIDR)

		basePodNumber, _ := data["BasePodNumber"].(int)
		_ = d.Set("base_pod_num", basePodNumber)

		isNonStaticIpMode, _ := data["IsNonStaticIpMode"].(bool)
		_ = d.Set("is_non_static_ip_mode", isNonStaticIpMode)

		_ = d.Set("runtime_version", info.RuntimeVersion)
		_ = d.Set("cluster_os_type", info.OsCustomizeType)
		_ = d.Set("container_runtime", info.ContainerRuntime)
		_ = d.Set("kube_proxy_mode", info.KubeProxyMode)
		_ = d.Set("service_cidr", info.ServiceCIDR)
		_ = d.Set("upgrade_instances_follow_cluster", false)

		switchSet, err := service.DescribeLogSwitches(ctx, d.Id())
		if err != nil {
			return err
		}
		logAgents := make([]map[string]interface{}, 0)
		events := make([]map[string]interface{}, 0)
		audits := make([]map[string]interface{}, 0)
		for _, switchItem := range switchSet {
			if switchItem.Log != nil && switchItem.Log.Enable != nil && helper.PBool(switchItem.Log.Enable) {
				logAgent := map[string]interface{}{
					"enabled": helper.PBool(switchItem.Log.Enable),
				}
				logAgents = append(logAgents, logAgent)
			}
			if switchItem.Event != nil && switchItem.Event.Enable != nil && helper.PBool(switchItem.Event.Enable) {
				event := map[string]interface{}{
					"enabled":    helper.PBool(switchItem.Event.Enable),
					"log_set_id": helper.PString(switchItem.Event.LogsetId),
					"topic_id":   helper.PString(switchItem.Event.TopicId),
				}
				events = append(events, event)
			}
			if switchItem.Audit != nil && switchItem.Audit.Enable != nil && helper.PBool(switchItem.Audit.Enable) {
				audit := map[string]interface{}{
					"enabled":    helper.PBool(switchItem.Audit.Enable),
					"log_set_id": helper.PString(switchItem.Audit.LogsetId),
					"topic_id":   helper.PString(switchItem.Audit.TopicId),
				}
				audits = append(audits, audit)
			}
		}
		if len(logAgents) > 0 {
			_ = d.Set("log_agent", logAgents)
		}
		if len(events) > 0 {
			_ = d.Set("event_persistence", events)
		}
		if len(audits) > 0 {
			_ = d.Set("cluster_audit", audits)
		}

		applist, err := service.DescribeExtensionAddonList(ctx, d.Id())
		if err != nil {
			return err
		}
		addons := make([]map[string]interface{}, 0)
		for _, item := range applist.Items {
			if item.Status.Phase == "Succeeded" && item.Labels["application.tkestack.io/type"] == "internal-addon" {
				addonParam := AddonRequestBody{
					Kind: helper.String("App"),
					Spec: &AddonSpec{
						Chart: &AddonSpecChart{
							ChartName:    item.Spec.Chart.ChartName,
							ChartVersion: item.Spec.Chart.ChartVersion,
						},
						Values: &AddonSpecValues{
							Values:        item.Spec.Values.Values,
							RawValues:     item.Spec.Values.RawValues,
							RawValuesType: item.Spec.Values.RawValuesType,
						},
					},
				}
				result, err := json.Marshal(addonParam)
				if err != nil {
					return err
				}

				addon := map[string]interface{}{
					"name":  item.Name,
					"param": string(result),
				}
				addons = append(addons, addon)
			}
		}
		if len(addons) > 0 {
			_ = d.Set("extension_addon", addons)
		}

		resp, err := service.DescribeClusterExtraArgs(ctx, d.Id())
		if err != nil {
			return err
		}
		fmt.Println(&resp)
		flag := false
		extraArgs := make(map[string]interface{}, 0)
		if len(resp.KubeAPIServer) > 0 {
			flag = true
			extraArgs["kube_apiserver"] = resp.KubeAPIServer
		}
		if len(resp.KubeControllerManager) > 0 {
			flag = true
			extraArgs["kube_controller_manager"] = resp.KubeControllerManager
		}
		if len(resp.KubeScheduler) > 0 {
			flag = true
			extraArgs["kube_scheduler"] = resp.KubeScheduler
		}

		if flag {
			_ = d.Set("cluster_extra_args", []map[string]interface{}{extraArgs})
		}

		if networkType == TKE_CLUSTER_NETWORK_TYPE_CILIUM_OVERLAY {
			resp, err := service.DescribeExternalNodeSupportConfig(ctx, d.Id())
			if err != nil {
				return err
			}
			_ = d.Set("cluster_subnet_id", resp.SubnetId)
		}
		if networkType == TKE_CLUSTER_NETWORK_TYPE_VPC_CNI {
			resp, err := service.DescribeIPAMD(ctx, d.Id())
			if err != nil {
				return err
			}
			_ = d.Set("eni_subnet_ids", helper.PStrings(resp.SubnetIds))

			duration, err := time.ParseDuration(helper.PString(resp.ClaimExpiredDuration))
			if err != nil {
				return err
			}
			seconds := int(duration.Seconds())
			if seconds > 0 {
				_ = d.Set("claim_expired_seconds", seconds)
			}
		}

		if info.DeployType == TKE_DEPLOY_TYPE_MANAGED {
			options, state, _, err := service.DescribeClusterAuthenticationOptions(ctx, d.Id())
			if err != nil {
				return err
			}
			if state == "Success" {
				authOptions := make(map[string]interface{}, 0)
				if helper.PBool(options.UseTKEDefault) {
					authOptions["use_tke_default"] = helper.PBool(options.UseTKEDefault)
				} else {
					authOptions["jwks_uri"] = helper.PString(options.JWKSURI)
					authOptions["issuer"] = helper.PString(options.Issuer)
				}
				authOptions["auto_create_discovery_anonymous_auth"] = helper.PBool(options.AutoCreateDiscoveryAnonymousAuth)
				_ = d.Set("auth_options", []map[string]interface{}{authOptions})
			}
		}
	}

	if _, ok := d.GetOkExists("auto_upgrade_cluster_level"); ok {
		_ = d.Set("auto_upgrade_cluster_level", info.AutoUpgradeClusterLevel)
	} else if importClsFlag {
		_ = d.Set("auto_upgrade_cluster_level", info.AutoUpgradeClusterLevel)
		importClsFlag = false
	}

	err = checkClusterEndpointStatus(ctx, &service, d, false)
	if err != nil {
		return fmt.Errorf("get internet failed, %s", err.Error())
	}

	err = checkClusterEndpointStatus(ctx, &service, d, true)
	if err != nil {
		return fmt.Errorf("get intranet failed, %s\n", err.Error())
	}

	_, workers, err := service.DescribeClusterInstances(ctx, d.Id())
	if err != nil {
		err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
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
		return fmt.Errorf("get worker instances failed, %s", err.Error())
	}

	workerInstancesList := make([]map[string]interface{}, 0, len(workers))
	for _, worker := range workers {
		tempMap := make(map[string]interface{})
		tempMap["instance_id"] = worker.InstanceId
		tempMap["instance_role"] = worker.InstanceRole
		tempMap["instance_state"] = worker.InstanceState
		tempMap["failed_reason"] = worker.FailedReason
		tempMap["lan_ip"] = worker.LanIp
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

	//if v, ok := d.GetOk("worker_config"); ok && len(v.([]interface{})) > 0 {
	//	if emptyStrFunc(securityRet.Response.ClusterExternalEndpoint) == "" {
	//		_ = d.Set("cluster_internet", false)
	//	} else {
	//		_ = d.Set("cluster_internet", true)
	//	}
	//
	//	if emptyStrFunc(securityRet.Response.PgwEndpoint) == "" {
	//		_ = d.Set("cluster_intranet", false)
	//	} else {
	//		_ = d.Set("cluster_intranet", true)
	//	}
	//}

	var globalConfig *tke.ClusterAsGroupOption
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		globalConfig, err = service.DescribeClusterNodePoolGlobalConfig(ctx, d.Id())
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
	if err != nil {
		return err
	}

	if globalConfig != nil {
		temp := make(map[string]interface{})
		temp["is_scale_in_enabled"] = globalConfig.IsScaleDownEnabled
		temp["expander"] = globalConfig.Expander
		temp["max_concurrent_scale_in"] = globalConfig.MaxEmptyBulkDelete
		temp["scale_in_delay"] = globalConfig.ScaleDownDelay
		temp["scale_in_unneeded_time"] = globalConfig.ScaleDownUnneededTime
		temp["scale_in_utilization_threshold"] = globalConfig.ScaleDownUtilizationThreshold
		temp["ignore_daemon_sets_utilization"] = globalConfig.IgnoreDaemonSetsUtilization
		temp["skip_nodes_with_local_storage"] = globalConfig.SkipNodesWithLocalStorage
		temp["skip_nodes_with_system_pods"] = globalConfig.SkipNodesWithSystemPods

		_ = d.Set("node_pool_global_config", []map[string]interface{}{temp})
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

	if d.HasChange("cluster_subnet_id") {
		return fmt.Errorf("argument cluster_subnet_id cannot be changed")
	}

	if d.HasChange("tags") {
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))

		resourceName := BuildTagResourceName("ccs", "cluster", region, id)
		if err := service.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}

	}

	var (
		clusterInternet              = d.Get("cluster_internet").(bool)
		clusterIntranet              = d.Get("cluster_intranet").(bool)
		intranetSubnetId             = d.Get("cluster_intranet_subnet_id").(string)
		clusterInternetSecurityGroup = d.Get("cluster_internet_security_group").(string)
		clusterInternetDomain        = d.Get("cluster_internet_domain").(string)
		clusterIntranetDomain        = d.Get("cluster_intranet_domain").(string)
	)

	if clusterIntranet && intranetSubnetId == "" {
		return fmt.Errorf("`cluster_intranet_subnet_id` must set when `cluster_intranet` is true")
	}

	if d.HasChange("cluster_intranet_subnet_id") && !d.HasChange("cluster_intranet") {
		return fmt.Errorf("`cluster_intranet_subnet_id` must modified with `cluster_intranet`")
	}

	if d.HasChange("cluster_internet_security_group") && !d.HasChange("cluster_internet") {
		if clusterInternet {
			err := tkeService.ModifyClusterEndpointSG(ctx, id, clusterInternetSecurityGroup)
			if err != nil {
				return err
			}
		}
	}

	if d.HasChange("cluster_intranet") {
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTRANET, clusterIntranet, clusterInternetSecurityGroup, intranetSubnetId, clusterIntranetDomain); err != nil {
			return err
		}

	}

	if d.HasChange("cluster_internet") {
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTERNET, clusterInternet, clusterInternetSecurityGroup, "", clusterInternetDomain); err != nil {
			return err
		}
	}

	// situation when only domain changed
	if !d.HasChange("cluster_intranet") && clusterIntranet && d.HasChange("cluster_intranet_domain") {
		// recreate the cluster intranet endpoint using new domain
		// first close
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTRANET, TKE_CLUSTER_CLOSE_ACCESS, clusterInternetSecurityGroup, intranetSubnetId, clusterIntranetDomain); err != nil {
			return err
		}
		// then reopen
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTRANET, TKE_CLUSTER_OPEN_ACCESS, clusterInternetSecurityGroup, intranetSubnetId, clusterIntranetDomain); err != nil {
			return err
		}
	}
	if !d.HasChange("cluster_internet") && clusterInternet && d.HasChange("cluster_internet_domain") {
		// recreate the cluster internet endpoint using new domain
		// first close
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTERNET, TKE_CLUSTER_CLOSE_ACCESS, clusterInternetSecurityGroup, "", clusterInternetDomain); err != nil {
			return err
		}
		// then reopen
		if err := ModifyClusterInternetOrIntranetAccess(ctx, d, &tkeService, TKE_CLUSTER_INTERNET, TKE_CLUSTER_OPEN_ACCESS, clusterInternetSecurityGroup, "", clusterInternetDomain); err != nil {
			return err
		}
	}

	if d.HasChange("project_id") || d.HasChange("cluster_name") || d.HasChange("cluster_desc") || d.HasChange("cluster_level") || d.HasChange("auto_upgrade_cluster_level") {
		projectId := int64(d.Get("project_id").(int))
		clusterName := d.Get("cluster_name").(string)
		clusterDesc := d.Get("cluster_desc").(string)
		clusterLevel := d.Get("cluster_level").(string)
		autoUpgradeClusterLevel := d.Get("auto_upgrade_cluster_level").(bool)

		ins, _, err := tkeService.DescribeCluster(ctx, id)
		if err != nil {
			return err
		}

		//ignore same cluster level if same
		if *ins.ClusterLevel == clusterLevel {
			clusterLevel = ""
		}

		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			err := tkeService.ModifyClusterAttribute(ctx, id, projectId, clusterName, clusterDesc, clusterLevel, autoUpgradeClusterLevel)
			if err != nil {
				// create and update immediately may cause cluster level syntax error, this error can wait until cluster level state normal
				return retryError(err, tke.INTERNALERROR_UNEXPECTEDINTERNAL, tke.RESOURCEUNAVAILABLE)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	//update VPC-CNI container network capability
	if !d.HasChange("eni_subnet_ids") && (d.HasChange("vpc_cni_type") || d.HasChange("claim_expired_seconds")) {
		err := fmt.Errorf("changing only `vpc_cni_type` or `claim_expired_seconds` is not supported, when turning on or off the vpc-cni container network capability, `eni_subnet_ids` must be changed")
		return err
	}
	if d.HasChange("eni_subnet_ids") {
		eniSubnetIdList := d.Get("eni_subnet_ids").([]interface{})
		if len(eniSubnetIdList) == 0 {
			err := tkeService.DisableVpcCniNetworkType(ctx, id)
			if err != nil {
				return err
			}
			time.Sleep(3 * time.Second)
			err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
				ipamdResp, inErr := tkeService.DescribeIPAMD(ctx, id)
				enableIPAMD := *ipamdResp.EnableIPAMD
				disableVpcCniMode := *ipamdResp.DisableVpcCniMode
				phase := *ipamdResp.Phase
				if inErr != nil {
					return resource.NonRetryableError(inErr)
				}
				if !enableIPAMD || (disableVpcCniMode && phase != "upgrading") {
					return nil
				}
				return resource.RetryableError(fmt.Errorf("%s close vpc cni network type task is in progress and waiting to be completed", id))
			})
			if err != nil {
				return err
			}
		} else {
			info, _, err := tkeService.DescribeCluster(ctx, id)
			if err != nil {
				err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
					newInfo, _, inErr := tkeService.DescribeCluster(ctx, id)
					if inErr != nil {
						return retryError(inErr)
					}
					info = newInfo
					return nil
				})
				if err != nil {
					return err
				}
			}
			oldSubnets := info.EniSubnetIds
			var subnets []string
			for index := range eniSubnetIdList {
				subnetId := eniSubnetIdList[index].(string)
				subnets = append(subnets, subnetId)
			}
			if len(oldSubnets) > 0 {
				exist, addSubnets := helper.CheckElementsExist(oldSubnets, subnets)
				if !exist {
					err = fmt.Errorf("the `eni_subnet_ids` parameter does not allow modification of existing subnet ID data %v. "+
						"if you want to modify the existing subnet ID, please first set eni_subnet_ids to empty to turn off the VPC-CNI network capability, "+
						"and then fill in the latest subnet ID", oldSubnets)
					return err
				}
				if d.HasChange("vpc_cni_type") || d.HasChange("claim_expired_seconds") {
					err = fmt.Errorf("modifying `vpc_cni_type` and `claim_expired_seconds` is not supported when adding a cluster subnet")
					return err
				}
				if len(addSubnets) > 0 {
					vpcId := d.Get("vpc_id").(string)
					err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
						inErr := tkeService.AddVpcCniSubnets(ctx, id, addSubnets, vpcId)
						if inErr != nil {
							return resource.NonRetryableError(inErr)
						}
						return nil
					})
					if err != nil {
						return err
					}
				}
			} else {
				var vpcCniType string
				if v, ok := d.GetOk("vpc_cni_type"); ok {
					vpcCniType = v.(string)
				} else {
					vpcCniType = "tke-route-eni"
				}
				enableStaticIp := !d.Get("is_non_static_ip_mode").(bool)
				expiredSeconds := uint64(d.Get("claim_expired_seconds").(int))

				err = tkeService.EnableVpcCniNetworkType(ctx, id, vpcCniType, enableStaticIp, subnets, expiredSeconds)
				if err != nil {
					return err
				}
				time.Sleep(3 * time.Second)
				err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
					ipamdResp, inErr := tkeService.DescribeIPAMD(ctx, id)
					disableVpcCniMode := *ipamdResp.DisableVpcCniMode
					phase := *ipamdResp.Phase
					if inErr != nil {
						return resource.NonRetryableError(inErr)
					}
					if !disableVpcCniMode && phase == "running" {
						return nil
					}
					if !disableVpcCniMode && phase == "initializing" {
						return resource.RetryableError(fmt.Errorf("%s enable vpc cni network type task is in progress and waiting to be completed", id))
					}
					return resource.NonRetryableError(fmt.Errorf("%s enable vpc cni network type task disableVpcCniMode is %v and phase is %s,we won't wait for it finish", id, disableVpcCniMode, phase))
				})
				if err != nil {
					return err
				}
			}
		}
	}

	//upgrade k8s cluster version
	if d.HasChange("cluster_version") {
		newVersion := d.Get("cluster_version").(string)
		isOk, err := tkeService.CheckClusterVersion(ctx, id, newVersion)
		if err != nil {
			return err
		}
		if !isOk {
			return fmt.Errorf("version %s is unsupported", newVersion)
		}
		extraArgs, ok := d.GetOk("cluster_extra_args")
		if !ok {
			extraArgs = nil
		}
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := tkeService.ModifyClusterVersion(ctx, id, newVersion, extraArgs)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}
		//check status
		err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
			ins, has, inErr := tkeService.DescribeCluster(ctx, id)
			if inErr != nil {
				return retryError(inErr)
			}
			if !has {
				return resource.NonRetryableError(fmt.Errorf("Cluster %s is not exist", id))
			}
			if ins.ClusterStatus == "Running" {
				return nil
			} else {
				return resource.RetryableError(fmt.Errorf("cluster %s status %s, retry...", id, ins.ClusterStatus))
			}
		})
		if err != nil {
			return err
		}

		// upgrade instances version
		upgrade := false
		if v, ok := d.GetOk("upgrade_instances_follow_cluster"); ok {
			upgrade = v.(bool)
		}
		if upgrade {
			err := upgradeClusterInstances(tkeService, ctx, id)
			if err != nil {
				return err
			}
		}
	}

	// update node pool global config
	if d.HasChange("node_pool_global_config") {
		request := tkeGetNodePoolGlobalConfig(d)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := tkeService.ModifyClusterNodePoolGlobalConfig(ctx, request)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if err != nil {
			return err
		}

	}

	if d.HasChange("auth_options") {
		request := tkeGetAuthOptions(d)
		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			inErr := tkeService.ModifyClusterAuthenticationOptions(ctx, request)
			if inErr != nil {
				return retryError(inErr, tke.RESOURCEUNAVAILABLE_CLUSTERSTATE)
			}
			return nil
		})
		if err != nil {
			return err
		}
		_, _, err = tkeService.WaitForAuthenticationOptionsUpdateSuccess(ctx, id)
		if err != nil {
			return err
		}
	}

	if d.HasChange("deletion_protection") {
		enable := d.Get("deletion_protection").(bool)
		if err := tkeService.ModifyDeletionProtection(ctx, id, enable); err != nil {
			return err
		}

	}

	if d.HasChange("acquire_cluster_admin_role") {
		o, n := d.GetChange("acquire_cluster_admin_role")
		if o.(bool) && !n.(bool) {
			return fmt.Errorf("argument `acquire_cluster_admin_role` cannot set to false")
		}
		err := tkeService.AcquireClusterAdminRole(ctx, id)
		if err != nil {
			return err
		}
	}

	if d.HasChange("log_agent") {
		v, ok := helper.InterfacesHeadMap(d, "log_agent")
		enabled := false
		rootDir := ""
		if ok {
			rootDir = v["kubelet_root_dir"].(string)
			enabled = v["enabled"].(bool)
		}
		err := tkeService.SwitchLogAgent(ctx, id, rootDir, enabled)
		if err != nil {
			return err
		}
	}

	if d.HasChange("event_persistence") {
		v, ok := helper.InterfacesHeadMap(d, "event_persistence")
		enabled := false
		logSetId := ""
		topicId := ""
		deleteEventLog := false
		if ok {
			enabled = v["enabled"].(bool)
			logSetId = v["log_set_id"].(string)
			topicId = v["topic_id"].(string)
			deleteEventLog = v["delete_event_log_and_topic"].(bool)
		}

		err := tkeService.SwitchEventPersistence(ctx, id, logSetId, topicId, enabled, deleteEventLog)
		if err != nil {
			return err
		}
	}

	if d.HasChange("cluster_audit") {
		v, ok := helper.InterfacesHeadMap(d, "cluster_audit")
		enabled := false
		logSetId := ""
		topicId := ""
		deleteAuditLog := false
		if ok {
			enabled = v["enabled"].(bool)
			logSetId = v["log_set_id"].(string)
			topicId = v["topic_id"].(string)
			deleteAuditLog = v["delete_audit_log_and_topic"].(bool)
		}

		err := tkeService.SwitchClusterAudit(ctx, id, logSetId, topicId, enabled, deleteAuditLog)
		if err != nil {
			return err
		}
	}

	if d.HasChange("extension_addon") {
		o, n := d.GetChange("extension_addon")
		adds, removes, changes := resourceTkeGetAddonsDiffs(o.([]interface{}), n.([]interface{}))
		updates := append(adds, changes...)
		for i := range updates {
			var err error
			addon := updates[i].(map[string]interface{})
			param := addon["param"].(string)
			name, err := tkeService.GetAddonNameFromJson(param)
			if err != nil {
				return err
			}
			_, has, _ := tkeService.PollingAddonsPhase(ctx, id, name, nil)
			if has {
				err = tkeService.UpdateExtensionAddon(ctx, id, name, param)
			} else {
				err = tkeService.CreateExtensionAddon(ctx, id, param)
			}
			if err != nil {
				return err
			}
			_, _, err = tkeService.PollingAddonsPhase(ctx, id, name, nil)
			if err != nil {
				return err
			}
		}

		for i := range removes {
			addon := removes[i].(map[string]interface{})
			param := addon["param"].(string)
			name, err := tkeService.GetAddonNameFromJson(param)
			if err != nil {
				return err
			}
			_, has, _ := tkeService.PollingAddonsPhase(ctx, id, name, nil)
			if !has {
				continue
			}
			err = tkeService.DeleteExtensionAddon(ctx, id, name)
			if err != nil {
				return err
			}
			_, has, _ = tkeService.PollingAddonsPhase(ctx, id, name, nil)
			if has {
				return fmt.Errorf("addon %s still exists", name)
			}
		}

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
	deleteEventLogSetAndTopic := false
	enableEventLog := false
	deleteAuditLogSetAndTopic := false
	if v, ok := helper.InterfacesHeadMap(d, "event_persistence"); ok {
		deleteEventLogSetAndTopic = v["delete_event_log_and_topic"].(bool)
		// get cluster current enabled status
		enableEventLog = v["enabled"].(bool)
	}

	if v, ok := helper.InterfacesHeadMap(d, "cluster_audit"); ok {
		deleteAuditLogSetAndTopic = v["delete_audit_log_and_topic"].(bool)
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		if deleteEventLogSetAndTopic && enableEventLog {
			err := service.SwitchEventPersistence(ctx, d.Id(), "", "", false, true)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() != "FailedOperation.ClusterNotFound" {
					return retryError(err, InternalError)
				}
			} else if err != nil {
				return retryError(err, InternalError)
			}
		}
		if deleteAuditLogSetAndTopic {
			err := service.SwitchClusterAudit(ctx, d.Id(), "", "", false, true)
			if e, ok := err.(*errors.TencentCloudSDKError); ok {
				if e.GetCode() != "ResourceNotFound.ClusterNotFound" {
					return retryError(err, InternalError)
				}
			} else if err != nil {
				return retryError(err, InternalError)
			}
		}
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

func resourceTkeGetAddonsDiffs(o, n []interface{}) (adds, removes, changes []interface{}) {
	indexByName := func(i interface{}) int {
		v := i.(map[string]interface{})
		return helper.HashString(v["name"].(string))
	}
	indexAll := func(i interface{}) int {
		v := i.(map[string]interface{})
		name := v["name"].(string)
		param := v["param"].(string)
		return helper.HashString(fmt.Sprintf("%s#%s", name, param))
	}

	os := schema.NewSet(indexByName, o)
	ns := schema.NewSet(indexByName, n)

	adds = ns.Difference(os).List()
	removes = os.Difference(ns).List()

	fullIndexedKeeps := schema.NewSet(indexAll, ns.Intersection(os).List())
	fullIndexedOlds := schema.NewSet(indexAll, o)

	changes = fullIndexedKeeps.Difference(fullIndexedOlds).List()
	return
}

func checkClusterEndpointStatus(ctx context.Context, service *TkeService, d *schema.ResourceData, isInternet bool) (err error) {
	var status, config string
	var response tke.DescribeClusterEndpointsResponseParams
	var isOpened bool
	var errRet error
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		status, _, errRet = service.DescribeClusterEndpointStatus(ctx, d.Id(), isInternet)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if status == TkeInternetStatusCreating || status == TkeInternetStatusDeleting {
			return resource.RetryableError(
				fmt.Errorf("%s create cluster internet endpoint status still is %s", d.Id(), status))
		}
		return nil
	})
	if err != nil {
		return err
	}
	if status == TkeInternetStatusNotfound || status == TkeInternetStatusDeleted {
		isOpened = false
	}
	if status == TkeInternetStatusCreated {
		isOpened = true
	}
	if isInternet {
		_ = d.Set("cluster_internet", isOpened)
	} else {
		_ = d.Set("cluster_intranet", isOpened)
	}

	if isOpened {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			config, errRet = service.DescribeClusterConfig(ctx, d.Id(), isInternet)
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, errRet = service.DescribeClusterEndpoints(ctx, d.Id())
			if errRet != nil {
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if isInternet {
			_ = d.Set("kube_config", config)
			_ = d.Set("cluster_internet_domain", helper.PString(response.ClusterExternalDomain))
			_ = d.Set("cluster_internet_security_group", helper.PString(response.SecurityGroup))
		} else {
			_ = d.Set("kube_config_intranet", config)
			_ = d.Set("cluster_intranet_domain", helper.PString(response.ClusterIntranetDomain))
			_ = d.Set("cluster_intranet_subnet_id", helper.PString(response.ClusterIntranetSubnetId))
		}

	} else {
		if isInternet {
			_ = d.Set("kube_config", "")
		} else {
			_ = d.Set("kube_config_intranet", "")
		}
	}
	return nil
}
