package tke

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

func ResourceTencentCloudKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesClusterCreate,
		Read:   resourceTencentCloudKubernetesClusterRead,
		Update: resourceTencentCloudKubernetesClusterUpdate,
		Delete: resourceTencentCloudKubernetesClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: customResourceImporter,
		},
		CustomizeDiff: customdiff.All(
			customizeDiffForContainerRuntimeDefault,
		),
		Schema: map[string]*schema.Schema{
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
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tlinux2.4x86_64",
				Description: "Cluster operating system, supports setting public images (the field passes the corresponding image Name) and custom images (the field passes the corresponding image ID). For details, please refer to: https://cloud.tencent.com/document/product/457/68289.",
			},

			"cluster_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Control Plane Subnet Information. This field is required only in the following scenarios: When the container network plugin is CiliumOverlay, TKE will obtain 2 IPs from this subnet to create an internal load balancer; When creating a managed cluster that supports CDC with the VPC-CNI network plugin, at least 12 IPs must be reserved.",
			},

			"cluster_os_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "GENERAL",
				Description:  "Image type of the cluster os, the available values include: 'GENERAL'. Default is 'GENERAL'.",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_CLUSTER_OS_TYPES),
			},

			"container_runtime": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "Runtime type of the cluster, the available values include: 'docker' and 'containerd'.The Kubernetes v1.24 has removed dockershim, so please use containerd in v1.24 or higher. The default value is `docker` for versions below v1.24 and `containerd` for versions above v1.24.",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_RUNTIMES),
			},

			"cluster_deploy_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "MANAGED_CLUSTER",
				Description:  "Deployment type of the cluster, the available values include: 'MANAGED_CLUSTER' and 'INDEPENDENT_CLUSTER'. Default is 'MANAGED_CLUSTER'.",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_DEPLOY_TYPES),
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
				Description: "Indicates whether upgrade all cluster instances. Default is false.",
			},

			"cluster_ipvs": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "Indicates whether `ipvs` is enabled. Default is true. False means `iptables` is enabled.",
			},

			"cluster_as_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether to enable cluster node auto scaling. Default is false.",
				Deprecated:  "This argument is deprecated because the TKE auto-scaling group was no longer available.",
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
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Global config effective for all node pools.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},

			"cluster_extra_args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Customized parameters for master component,such as kube-apiserver, kube-controller-manager, kube-scheduler.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kube_apiserver": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "The customized parameters for kube-apiserver.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"kube_controller_manager": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "The customized parameters for kube-controller-manager.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"kube_scheduler": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "The customized parameters for kube-scheduler.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"is_dual_stack": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "In the VPC-CNI mode of the cluster, the dual stack cluster status defaults to false, indicating a non dual stack cluster.",
			},

			"node_name_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "lan-ip",
				Description:  "Node name type of Cluster, the available values include: 'lan-ip' and 'hostname', Default is 'lan-ip'.",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_CLUSTER_NODE_NAME_TYPE),
			},

			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "GR",
				Description:  "Cluster network type, the available values include: 'GR' and 'VPC-CNI' and 'CiliumOverlay'. Default is GR.",
				ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_CLUSTER_NETWORK_TYPE),
			},

			"enable_customized_pod_cidr": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable the custom mode of node podCIDR size. Default is false.",
			},

			"base_pod_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The number of basic pods. valid when enable_customized_pod_cidr=true.",
			},

			"is_non_static_ip_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Indicates whether non-static ip mode is enabled. Default is false.",
			},

			"data_plane_v2": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Whether to enable DataPlaneV2 (replace kube-proxy with cilium). `data_plane_v2` and `cluster_ipvs` should not be set at the same time.",
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates whether cluster deletion protection is enabled. Default is false.",
			},

			"resource_delete_options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "The resource deletion policy when the cluster is deleted. Currently, CBS is supported (CBS is retained by default). Only valid when deleting cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource type, valid values are `CBS`, `CLB`, and `CVM`.",
						},
						"delete_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The deletion mode of CBS resources when the cluster is deleted, `terminate` (destroy), `retain` (retain). Other resources are deleted by default.",
						},
						"skip_deletion_protection": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to skip resources with deletion protection enabled, the default is false.",
						},
					},
				},
			},

			"kube_proxy_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Cluster kube-proxy mode, the available values include: 'kube-proxy-bpf'. Default is not set.When set to kube-proxy-bpf, cluster version greater than 1.14 and with Tencent Linux 2.4 is required.",
			},

			"vpc_cni_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Distinguish between shared network card multi-IP mode and independent network card mode. Fill in `tke-route-eni` for shared network card multi-IP mode and `tke-direct-eni` for independent network card mode. The default is shared network card mode. When it is necessary to turn off the vpc-cni container network capability, both `eni_subnet_ids` and `vpc_cni_type` must be set to empty.",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"tke-route-eni", "tke-direct-eni"}),
			},

			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Vpc Id of the cluster.",
				ValidateFunc: tccommon.ValidateStringLengthInRange(4, 100),
			},

			"cluster_internet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Open internet access or not. If this field is set 'true', the field below `worker_config` must be set. Because only cluster with node is allowed enable access endpoint. You may open it through `tencentcloud_kubernetes_cluster_endpoint`.",
			},

			"cluster_internet_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Domain name for cluster Kube-apiserver internet access. Be careful if you modify value of this parameter, the cluster_external_endpoint value may be changed automatically too.",
			},

			"cluster_intranet": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Open intranet access or not. If this field is set 'true', the field below `worker_config` must be set. Because only cluster with node is allowed enable access endpoint. You may open it through `tencentcloud_kubernetes_cluster_endpoint`.",
			},

			"cluster_intranet_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Domain name for cluster Kube-apiserver intranet access. Be careful if you modify value of this parameter, the pgw_endpoint value may be changed automatically too.",
			},

			"cluster_internet_security_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify security group, NOTE: This argument must not be empty if cluster internet enabled.",
			},

			"managed_cluster_internet_security_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Security policies for managed cluster internet, like:'192.168.1.0/24' or '113.116.51.27', '0.0.0.0/0' means all. This field can only set when field `cluster_deploy_type` is 'MANAGED_CLUSTER' and `cluster_internet` is true. `managed_cluster_internet_security_policies` can not delete or empty once be set.",
				Deprecated:  "this argument was deprecated, use `cluster_internet_security_group` instead.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"cluster_intranet_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Subnet id who can access this independent cluster, this field must and can only set  when `cluster_intranet` is true. `cluster_intranet_subnet_id` can not modify once be set.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project ID, default value is 0.",
			},

			"cluster_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A network address block of the cluster. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.",
				// ValidateFunc: clusterCidrValidateFunc,
			},

			"ignore_cluster_cidr_conflict": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Indicates whether to ignore the cluster cidr conflict error. Default is false.",
			},

			"ignore_service_cidr_conflict": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Indicates whether to ignore the service cidr conflict error. Only valid in `VPC-CNI` mode.",
			},

			"cluster_max_pod_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     256,
				Description: "The maximum number of Pods per node in the cluster. Default is 256. The minimum value is 4. When its power unequal to 2, it will round upward to the closest power of 2.",
			},

			"cluster_max_service_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     256,
				Description: "The maximum number of services in the cluster. Default is 256. The range is from 32 to 32768. When its power unequal to 2, it will round upward to the closest power of 2.",
			},

			"service_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A network address block of the service. Different from vpc cidr and cidr of other clusters within this vpc. Must be in  10./192.168/172.[16-31] segments.",
				// ValidateFunc: serviceCidrValidateFunc,
			},

			"eni_subnet_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Subnet Ids for cluster with VPC-CNI network mode. This field can only set when field `network_type` is 'VPC-CNI'. `eni_subnet_ids` can not empty once be set.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"claim_expired_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Claim expired seconds to recycle ENI. This field can only set when field `network_type` is 'VPC-CNI'. `claim_expired_seconds` must greater or equal than 300 and less than 15768000.",
				ValidateFunc: claimExpiredSecondsValidateFunc,
			},

			"master_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Deploy the machine configuration information of the 'MASTER_ETCD' service, and create <=7 units for common users.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Default:     1,
							Description: "Number of cvm.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Indicates which availability zone will be used.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Default:     "sub machine of tke",
							Description: "Name of the CVMs.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Specified types of CVM instance.",
						},
						"instance_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "POSTPAID_BY_HOUR",
							Description:  "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`, `PREPAID` instance will not terminated after cluster deleted, and may not allow to delete before expired.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_INSTANCE_CHARGE_TYPE),
						},
						"instance_charge_type_prepaid_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Default:      1,
							Description:  "The tenancy (time unit is month) of the prepaid instance. NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
							ValidateFunc: tccommon.ValidateAllowedIntValue(svccvm.CVM_PREPAID_PERIOD),
						},
						"instance_charge_type_prepaid_renew_flag": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svccvm.CVM_PREPAID_RENEW_FLAG),
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Description:  "Private network ID.",
							ValidateFunc: tccommon.ValidateStringLengthInRange(4, 100),
						},
						"system_disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "CLOUD_PREMIUM",
							Description:  "System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `CLOUD_BASIC`, `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
						},
						"system_disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Default:      50,
							Description:  "Volume of system disk in GB. Default is `50`.",
							ValidateFunc: tccommon.ValidateIntegerInRange(20, 1024),
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
										Default:      "CLOUD_PREMIUM",
										Description:  "Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD` and `CLOUD_HSSD` and `CLOUD_TSSD`.",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
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
										Optional:    true,
										ForceNew:    true,
										Description: "File system, e.g. `ext3/ext4/xfs`.",
									},
									"auto_format_and_mount": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Default:     false,
										Description: "Indicate whether to auto format and mount or not. Default is `false`.",
									},
									"mount_target": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Mount target.",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "The name of the device or partition to mount.",
									},
								},
							},
						},
						"internet_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "TRAFFIC_POSTPAID_BY_HOUR",
							Description:  "Charge types for network traffic. Available values include `TRAFFIC_POSTPAID_BY_HOUR`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.INTERNET_CHARGE_ALLOW_TYPE),
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
							Optional:    true,
							ForceNew:    true,
							Description: "Specify whether to assign an Internet IP address.",
						},
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Sensitive:    true,
							Description:  "Password to access, should be set if `key_ids` not set.",
							ValidateFunc: tccommon.ValidateAsConfigPassword,
						},
						"key_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "ID list of keys, should be set if `password` not set.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"security_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Security groups to which a CVM instance belongs.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enhanced_security_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Description: "To specify whether to enable cloud security service. Default is TRUE.",
						},
						"enhanced_monitor_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "ase64-encoded User Data text, the length limit is 16KB.",
						},
						"cam_role_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "CAM role name authorized to access.",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The host name of the attached instance. Dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. Windows example: The length of the name character is [2, 15], letters (capitalization is not restricted), numbers and dashes (-) are allowed, dots (.) are not supported, and not all numbers are allowed. Examples of other types (Linux, etc.): The character length is [2, 60], and multiple dots are allowed. There is a segment between the dots. Each segment allows letters (with no limitation on capitalization), numbers and dashes (-).",
						},
						"disaster_recover_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Disaster recover groups to which a CVM instance belongs. Only support maximum 1.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"img_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The valid image id, format of img-xxx. Note: `img_id` will be replaced with the image corresponding to TKE `cluster_os`.",
							ValidateFunc: tccommon.ValidateImageID,
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Default:     0,
							Description: "Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, and it override `[globe_]desired_pod_num` for current node. Either all the fields `desired_pod_num` or none.",
						},
						"hpc_cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Id of cvm hpc cluster.",
						},
					},
				},
			},

			"worker_config": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Deploy the machine configuration information of the 'WORKER' service, and create <=20 units for common users. The other 'WORK' service are added by 'tencentcloud_kubernetes_scale_worker'.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Default:     1,
							Description: "Number of cvm.",
						},
						"availability_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Indicates which availability zone will be used.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Default:     "sub machine of tke",
							Description: "Name of the CVMs.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Specified types of CVM instance.",
						},
						"instance_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "POSTPAID_BY_HOUR",
							Description:  "The charge type of instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`, `PREPAID` instance will not terminated after cluster deleted, and may not allow to delete before expired.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(TKE_INSTANCE_CHARGE_TYPE),
						},
						"instance_charge_type_prepaid_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Default:      1,
							Description:  "The tenancy (time unit is month) of the prepaid instance. NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
							ValidateFunc: tccommon.ValidateAllowedIntValue(svccvm.CVM_PREPAID_PERIOD),
						},
						"instance_charge_type_prepaid_renew_flag": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svccvm.CVM_PREPAID_RENEW_FLAG),
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							Description:  "Private network ID.",
							ValidateFunc: tccommon.ValidateStringLengthInRange(4, 100),
						},
						"system_disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "CLOUD_PREMIUM",
							Description:  "System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `CLOUD_BASIC`, `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
						},
						"system_disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							Default:      50,
							Description:  "Volume of system disk in GB. Default is `50`.",
							ValidateFunc: tccommon.ValidateIntegerInRange(20, 1024),
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
										Default:      "CLOUD_PREMIUM",
										Description:  "Types of disk, available values: `CLOUD_PREMIUM` and `CLOUD_SSD` and `CLOUD_HSSD` and `CLOUD_TSSD`.",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
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
										Optional:    true,
										ForceNew:    true,
										Description: "File system, e.g. `ext3/ext4/xfs`.",
									},
									"auto_format_and_mount": {
										Type:        schema.TypeBool,
										Optional:    true,
										ForceNew:    true,
										Default:     false,
										Description: "Indicate whether to auto format and mount or not. Default is `false`.",
									},
									"mount_target": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "Mount target.",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "The name of the device or partition to mount.",
									},
								},
							},
						},
						"internet_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "TRAFFIC_POSTPAID_BY_HOUR",
							Description:  "Charge types for network traffic. Available values include `TRAFFIC_POSTPAID_BY_HOUR`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.INTERNET_CHARGE_ALLOW_TYPE),
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
							Optional:    true,
							ForceNew:    true,
							Description: "Specify whether to assign an Internet IP address.",
						},
						"password": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Sensitive:    true,
							Description:  "Password to access, should be set if `key_ids` not set.",
							ValidateFunc: tccommon.ValidateAsConfigPassword,
						},
						"key_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "ID list of keys, should be set if `password` not set.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"security_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Security groups to which a CVM instance belongs.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"enhanced_security_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Description: "To specify whether to enable cloud security service. Default is TRUE.",
						},
						"enhanced_monitor_service": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "ase64-encoded User Data text, the length limit is 16KB.",
						},
						"cam_role_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "CAM role name authorized to access.",
						},
						"hostname": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The host name of the attached instance. Dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. Windows example: The length of the name character is [2, 15], letters (capitalization is not restricted), numbers and dashes (-) are allowed, dots (.) are not supported, and not all numbers are allowed. Examples of other types (Linux, etc.): The character length is [2, 60], and multiple dots are allowed. There is a segment between the dots. Each segment allows letters (with no limitation on capitalization), numbers and dashes (-).",
						},
						"disaster_recover_group_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Disaster recover groups to which a CVM instance belongs. Only support maximum 1.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"img_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "The valid image id, format of img-xxx. Note: `img_id` will be replaced with the image corresponding to TKE `cluster_os`.",
							ValidateFunc: tccommon.ValidateImageID,
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Default:     0,
							Description: "Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, and it override `[globe_]desired_pod_num` for current node. Either all the fields `desired_pod_num` or none.",
						},
						"hpc_cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Id of cvm hpc cluster.",
						},
					},
				},
			},

			"exist_instance": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create tke cluster by existed instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_role": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Role of existed node. Value: MASTER_ETCD or WORKER.",
						},
						"instances_para": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Reinstallation parameters of an existing instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_ids": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Cluster IDs.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"security_group_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Security groups to which a CVM instance belongs.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"password": {
										Type:         schema.TypeString,
										Optional:     true,
										Sensitive:    true,
										Description:  "Password to access, should be set if `key_ids` not set.",
										ValidateFunc: tccommon.ValidateAsConfigPassword,
									},
									"key_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "ID list of keys, should be set if `password` not set.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"enhanced_security_service": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "To specify whether to enable cloud security service. Default is TRUE.",
									},
									"enhanced_monitor_service": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
									},
									"master_config": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "Advanced Node Settings. commonly used to attach existing instances.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_target": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Mount target. Default is not mounting.",
												},
												"docker_graph_path": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Docker graph path. Default is `/var/lib/docker`.",
												},
												"user_script": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "User script encoded in base64, which will be executed after the k8s component runs. The user needs to ensure the script's reentrant and retry logic. The script and its generated log files can be viewed in the node path /data/ccs_userscript/. If the node needs to be initialized before joining the schedule, it can be used in conjunction with the `unschedulable` parameter. After the final initialization of the userScript is completed, add the command \"kubectl uncordon nodename --kubeconfig=/root/.kube/config\" to add the node to the schedule.",
												},
												"unschedulable": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Set whether the joined nodes participate in scheduling, with a default value of 0, indicating participation in scheduling; Non 0 means not participating in scheduling.",
												},
												"labels": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Node label list.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Name of map.",
															},
															"value": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Value of map.",
															},
														},
													},
												},
												"data_disk": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Configurations of data disk.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disk_type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Types of disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD` and `CLOUD_BSSD`.",
															},
															"file_system": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "File system, e.g. `ext3/ext4/xfs`.",
															},
															"disk_size": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Volume of disk in GB. Default is `0`.",
															},
															"auto_format_and_mount": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Indicate whether to auto format and mount or not. Default is `false`.",
															},
															"mount_target": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Mount target.",
															},
															"disk_partition": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The name of the device or partition to mount. NOTE: this argument doesn't support setting in node pool, or will leads to mount error.",
															},
														},
													},
												},
												"extra_args": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "Custom parameter information related to the node. This is a white-list parameter.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"kubelet": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Kubelet custom parameter. The parameter format is [\"k1=v1\", \"k1=v2\"].",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"desired_pod_number": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Indicate to set desired pod number in node. valid when the cluster is podCIDR.",
												},
												"gpu_args": {
													Type:        schema.TypeList,
													Optional:    true,
													MaxItems:    1,
													Description: "GPU driver parameters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mig_enable": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Whether to enable MIG.",
															},
															"driver": {
																Type:         schema.TypeMap,
																Optional:     true,
																Description:  "GPU driver version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.",
																ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
															},
															"cuda": {
																Type:         schema.TypeMap,
																Optional:     true,
																Description:  "CUDA  version. Format like: `{ version: String, name: String }`. `version`: Version of GPU driver or CUDA; `name`: Name of GPU driver or CUDA.",
																ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
															},
															"cudnn": {
																Type:         schema.TypeMap,
																Optional:     true,
																Description:  "cuDNN version. Format like: `{ version: String, name: String, doc_name: String, dev_name: String }`. `version`: cuDNN version; `name`: cuDNN name; `doc_name`: Doc name of cuDNN; `dev_name`: Dev name of cuDNN.",
																ValidateFunc: tccommon.ValidateTkeGpuDriverVersion,
															},
															"custom_driver": {
																Type:        schema.TypeMap,
																Optional:    true,
																Description: "Custom GPU driver. Format like: `{address: String}`. `address`: URL of custom GPU driver address.",
															},
														},
													},
												},
												"taints": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Node taint.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Key of the taint.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Value of the taint.",
															},
															"effect": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Effect of the taint.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"desired_pod_numbers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Custom mode cluster, you can specify the number of pods for each node. corresponding to the existed_instances_para.instance_ids parameter.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},

			"auth_options": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Specify cluster authentication configuration. Only available for managed cluster and `cluster_version` >= 1.20.",
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
			},

			"extension_addon": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Information of the add-on to be installed. It is recommended to use resource `tencentcloud_kubernetes_addon` management cluster addon.",
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
							Description:      "Parameter of the add-on resource object in JSON string format, please check the example at the top of page for reference.",
							DiffSuppressFunc: helper.DiffSupressJSON,
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
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "when you want to close the cluster event persistence or delete the cluster, you can use this parameter to determine whether the event persistence log set and topic created by default will be deleted.",
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
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "when you want to close the cluster audit log or delete the cluster, you can use this parameter to determine whether the audit log set and topic created by default will be deleted.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The tags of the cluster.",
			},

			"cluster_node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of nodes in the cluster.",
			},

			"worker_instances_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of cvm within the 'WORKER' clusters. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},

			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Labels of tke cluster nodes.",
			},

			"unschedulable": {
				Type:             schema.TypeInt,
				Optional:         true,
				ForceNew:         true,
				Default:          0,
				Description:      "Sets whether the joining node participates in the schedule. Default is '0'. Participate in scheduling.",
				DiffSuppressFunc: unschedulableDiffSuppressFunc,
			},

			"mount_target": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Mount target. Default is not mounting.",
			},

			"globe_desired_pod_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicate to set desired pod number in node. valid when enable_customized_pod_cidr=true, and it takes effect for all nodes.",
			},

			"docker_graph_path": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          "/var/lib/docker",
				Description:      "Docker graph path. Default is `/var/lib/docker`.",
				DiffSuppressFunc: dockerGraphPathDiffSuppressFunc,
			},

			"pre_start_user_script": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Base64-encoded user script, executed before initializing the node, currently only effective for adding existing nodes.",
			},

			"extra_args": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom parameter information related to the node.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"runtime_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Container Runtime version.",
			},

			"kube_config": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Kubernetes config.",
			},

			"kube_config_intranet": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Kubernetes config of private network.",
			},

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
				Description: "Access policy.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"cdc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CDC ID.",
			},

			"instance_delete_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The strategy for deleting cluster instances: terminate (destroy instances, only support pay as you go cloud host instances) retain (remove only, keep instances), Default is terminate.",
			},

			"disable_addons": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "To prevent the installation of a specific Addon component, enter the corresponding AddonName.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudKubernetesClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		clusterId string
	)
	var (
		request  = tkev20180525.NewCreateClusterRequest()
		response = tkev20180525.NewCreateClusterResponse()
	)

	if v, ok := d.GetOk("cdc_id"); ok {
		request.CdcId = helper.String(v.(string))
	}

	clusterCIDRSettings := tkev20180525.ClusterCIDRSettings{}
	if v, ok := d.GetOk("cluster_cidr"); ok {
		clusterCIDRSettings.ClusterCIDR = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("ignore_cluster_cidr_conflict"); ok {
		clusterCIDRSettings.IgnoreClusterCIDRConflict = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("ignore_service_cidr_conflict"); ok {
		clusterCIDRSettings.IgnoreServiceCIDRConflict = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("cluster_max_service_num"); ok {
		clusterCIDRSettings.MaxClusterServiceNum = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOkExists("cluster_max_pod_num"); ok {
		clusterCIDRSettings.MaxNodePodNum = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOk("service_cidr"); ok {
		clusterCIDRSettings.ServiceCIDR = helper.String(v.(string))
	}
	request.ClusterCIDRSettings = &clusterCIDRSettings

	clusterBasicSettings := tkev20180525.ClusterBasicSettings{}
	if v, ok := d.GetOk("cluster_version"); ok {
		clusterBasicSettings.ClusterVersion = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		clusterBasicSettings.ClusterName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_desc"); ok {
		clusterBasicSettings.ClusterDescription = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("project_id"); ok {
		clusterBasicSettings.ProjectId = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("cluster_os_type"); ok {
		clusterBasicSettings.OsCustomizeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_subnet_id"); ok {
		clusterBasicSettings.SubnetId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cluster_level"); ok {
		clusterBasicSettings.ClusterLevel = helper.String(v.(string))
	}
	autoUpgradeClusterLevel := tkev20180525.AutoUpgradeClusterLevel{}
	if v, ok := d.GetOkExists("auto_upgrade_cluster_level"); ok {
		autoUpgradeClusterLevel.IsAutoUpgrade = helper.Bool(v.(bool))
	}
	clusterBasicSettings.AutoUpgradeClusterLevel = &autoUpgradeClusterLevel
	request.ClusterBasicSettings = &clusterBasicSettings

	clusterAdvancedSettings := tkev20180525.ClusterAdvancedSettings{}
	if v, ok := d.GetOkExists("cluster_ipvs"); ok {
		clusterAdvancedSettings.IPVS = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("cluster_as_enabled"); ok {
		clusterAdvancedSettings.AsEnabled = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("container_runtime"); ok {
		clusterAdvancedSettings.ContainerRuntime = helper.String(v.(string))
	}
	if v, ok := d.GetOk("node_name_type"); ok {
		clusterAdvancedSettings.NodeNameType = helper.String(v.(string))
	}
	if extraArgsMap, ok := helper.InterfacesHeadMap(d, "cluster_extra_args"); ok {
		clusterExtraArgs := tkev20180525.ClusterExtraArgs{}
		if v, ok := extraArgsMap["kube_apiserver"]; ok {
			kubeAPIServerSet := v.([]interface{})
			for i := range kubeAPIServerSet {
				if kubeAPIServer, ok := kubeAPIServerSet[i].(string); ok && kubeAPIServer != "" {
					clusterExtraArgs.KubeAPIServer = append(clusterExtraArgs.KubeAPIServer, helper.String(kubeAPIServer))
				}
			}
		}
		if v, ok := extraArgsMap["kube_controller_manager"]; ok {
			kubeControllerManagerSet := v.([]interface{})
			for i := range kubeControllerManagerSet {
				if kubeControllerManager, ok := kubeControllerManagerSet[i].(string); ok && kubeControllerManager != "" {
					clusterExtraArgs.KubeControllerManager = append(clusterExtraArgs.KubeControllerManager, helper.String(kubeControllerManager))
				}
			}
		}
		if v, ok := extraArgsMap["kube_scheduler"]; ok {
			kubeSchedulerSet := v.([]interface{})
			for i := range kubeSchedulerSet {
				if kubeScheduler, ok := kubeSchedulerSet[i].(string); ok && kubeScheduler != "" {
					clusterExtraArgs.KubeScheduler = append(clusterExtraArgs.KubeScheduler, helper.String(kubeScheduler))
				}
			}
		}
		clusterAdvancedSettings.ExtraArgs = &clusterExtraArgs
	}
	if v, ok := d.GetOkExists("is_dual_stack"); ok {
		clusterAdvancedSettings.IsDualStack = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("network_type"); ok {
		clusterAdvancedSettings.NetworkType = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("is_non_static_ip_mode"); ok {
		clusterAdvancedSettings.IsNonStaticIpMode = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("data_plane_v2"); ok {
		clusterAdvancedSettings.DataPlaneV2 = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("deletion_protection"); ok {
		clusterAdvancedSettings.DeletionProtection = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOk("kube_proxy_mode"); ok {
		clusterAdvancedSettings.KubeProxyMode = helper.String(v.(string))
	}
	if v, ok := d.GetOk("runtime_version"); ok {
		clusterAdvancedSettings.RuntimeVersion = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("enable_customized_pod_cidr"); ok {
		clusterAdvancedSettings.EnableCustomizedPodCIDR = helper.Bool(v.(bool))
	}
	if v, ok := d.GetOkExists("base_pod_num"); ok {
		clusterAdvancedSettings.BasePodNumber = helper.IntInt64(v.(int))
	}
	request.ClusterAdvancedSettings = &clusterAdvancedSettings

	instanceAdvancedSettings := tkev20180525.InstanceAdvancedSettings{}
	if v, ok := d.GetOkExists("globe_desired_pod_num"); ok {
		instanceAdvancedSettings.DesiredPodNumber = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("mount_target"); ok {
		instanceAdvancedSettings.MountTarget = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("unschedulable"); ok {
		instanceAdvancedSettings.Unschedulable = helper.IntInt64(v.(int))
	}
	request.InstanceAdvancedSettings = &instanceAdvancedSettings

	if v, ok := d.GetOk("extension_addon"); ok {
		for _, item := range v.([]interface{}) {
			extensionAddonsMap := item.(map[string]interface{})
			extensionAddon := tkev20180525.ExtensionAddon{}
			if v, ok := extensionAddonsMap["name"]; ok {
				extensionAddon.AddonName = helper.String(v.(string))
			}
			if v, ok := extensionAddonsMap["param"]; ok {
				extensionAddon.AddonParam = helper.String(v.(string))
			}
			request.ExtensionAddons = append(request.ExtensionAddons, &extensionAddon)
		}
	}

	if v, ok := d.GetOk("disable_addons"); ok {
		for _, item := range v.([]interface{}) {
			if disableAddon, ok := item.(string); ok {
				request.DisableAddons = append(request.DisableAddons, &disableAddon)
			}
		}
	}

	if err := resourceTencentCloudKubernetesClusterCreatePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().CreateClusterWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes cluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId

	if err := resourceTencentCloudKubernetesClusterCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(clusterId)

	return resourceTencentCloudKubernetesClusterRead(d, meta)
}

func resourceTencentCloudKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Id()

	respData, err := service.DescribeKubernetesClusterById(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData.CdcId != nil {
		_ = d.Set("cdc_id", respData.CdcId)
	}

	if respData.ClusterName != nil {
		_ = d.Set("cluster_name", respData.ClusterName)
	}

	if respData.ClusterDescription != nil {
		_ = d.Set("cluster_desc", respData.ClusterDescription)
	}

	if respData.ClusterVersion != nil {
		_ = d.Set("cluster_version", respData.ClusterVersion)
	}

	if respData.ClusterType != nil {
		_ = d.Set("cluster_deploy_type", respData.ClusterType)
	}

	if respData.ClusterNetworkSettings != nil {
		if respData.ClusterNetworkSettings.ClusterCIDR != nil {
			_ = d.Set("cluster_cidr", respData.ClusterNetworkSettings.ClusterCIDR)
		}

		if respData.ClusterNetworkSettings.IgnoreClusterCIDRConflict != nil {
			_ = d.Set("ignore_cluster_cidr_conflict", respData.ClusterNetworkSettings.IgnoreClusterCIDRConflict)
		}

		if respData.ClusterNetworkSettings.IgnoreServiceCIDRConflict != nil {
			_ = d.Set("ignore_service_cidr_conflict", respData.ClusterNetworkSettings.IgnoreServiceCIDRConflict)
		}

		if respData.ClusterNetworkSettings.MaxNodePodNum != nil {
			_ = d.Set("cluster_max_pod_num", respData.ClusterNetworkSettings.MaxNodePodNum)
		}

		if respData.ClusterNetworkSettings.MaxClusterServiceNum != nil {
			_ = d.Set("cluster_max_service_num", respData.ClusterNetworkSettings.MaxClusterServiceNum)
		}

		if respData.ClusterNetworkSettings.Ipvs != nil {
			_ = d.Set("cluster_ipvs", respData.ClusterNetworkSettings.Ipvs)
		}

		if respData.ClusterNetworkSettings.VpcId != nil {
			_ = d.Set("vpc_id", respData.ClusterNetworkSettings.VpcId)
		}

		if respData.ClusterNetworkSettings.Subnets != nil {
			_ = d.Set("eni_subnet_ids", respData.ClusterNetworkSettings.Subnets)
		}

		if respData.ClusterNetworkSettings.DataPlaneV2 != nil {
			_ = d.Set("data_plane_v2", respData.ClusterNetworkSettings.DataPlaneV2)
		}
	}

	if respData.ClusterNodeNum != nil {
		_ = d.Set("cluster_node_num", respData.ClusterNodeNum)
	}

	if respData.ProjectId != nil {
		_ = d.Set("project_id", respData.ProjectId)
	}

	if respData.DeletionProtection != nil {
		_ = d.Set("deletion_protection", respData.DeletionProtection)
	}

	if respData.ClusterLevel != nil {
		_ = d.Set("cluster_level", respData.ClusterLevel)
	}

	if err := resourceTencentCloudKubernetesClusterReadPostHandleResponse0(ctx, respData); err != nil {
		return err
	}

	var respData1 *tkev20180525.DescribeClusterInstancesResponseParams
	reqErr1 := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterById1(ctx, clusterId)
		if e != nil {
			if err := resourceTencentCloudKubernetesClusterReadRequestOnError1(ctx, result, e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		}
		respData1 = result
		return nil
	})
	if reqErr1 != nil {
		log.Printf("[CRITAL]%s read kubernetes cluster failed, reason:%+v", logId, reqErr1)
		return reqErr1
	}

	if respData1 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	instanceSetList := make([]map[string]interface{}, 0, len(respData1.InstanceSet))
	if respData1.InstanceSet != nil {
		for _, instanceSet := range respData1.InstanceSet {
			instanceSetMap := map[string]interface{}{}

			if instanceSet.InstanceId != nil {
				instanceSetMap["instance_id"] = instanceSet.InstanceId
			}

			if instanceSet.InstanceRole != nil {
				instanceSetMap["instance_role"] = instanceSet.InstanceRole
			}

			if instanceSet.InstanceState != nil {
				instanceSetMap["instance_state"] = instanceSet.InstanceState
			}

			if instanceSet.FailedReason != nil {
				instanceSetMap["failed_reason"] = instanceSet.FailedReason
			}

			if instanceSet.LanIP != nil {
				instanceSetMap["lan_ip"] = instanceSet.LanIP
			}

			instanceSetList = append(instanceSetList, instanceSetMap)
		}

		_ = d.Set("worker_instances_list", instanceSetList)
	}

	var respData2 *tkev20180525.DescribeClusterSecurityResponseParams
	reqErr2 := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesClusterById2(ctx, clusterId)
		if e != nil {
			if err := resourceTencentCloudKubernetesClusterReadRequestOnError2(ctx, result, e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		}
		respData2 = result
		return nil
	})
	if reqErr2 != nil {
		log.Printf("[CRITAL]%s read kubernetes cluster failed, reason:%+v", logId, reqErr2)
		return reqErr2
	}

	if respData2 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_cluster` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData2.UserName != nil {
		_ = d.Set("user_name", respData2.UserName)
	}

	if respData2.Password != nil {
		_ = d.Set("password", respData2.Password)
	}

	if respData2.CertificationAuthority != nil {
		_ = d.Set("certification_authority", respData2.CertificationAuthority)
	}

	if respData2.ClusterExternalEndpoint != nil {
		_ = d.Set("cluster_external_endpoint", respData2.ClusterExternalEndpoint)
	}

	if respData2.Domain != nil {
		_ = d.Set("domain", respData2.Domain)
	}

	if respData2.PgwEndpoint != nil {
		_ = d.Set("pgw_endpoint", respData2.PgwEndpoint)
	}

	if err := resourceTencentCloudKubernetesClusterReadPostHandleResponse2(ctx, respData2); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudKubernetesClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	immutableArgs := []string{"cdc_id", "extension_addon", "disable_addons"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	clusterId := d.Id()

	if err := resourceTencentCloudKubernetesClusterUpdateOnStart(ctx); err != nil {
		return err
	}

	needChange := false
	mutableArgs := []string{"project_id", "cluster_name", "cluster_desc", "cluster_level", "auto_upgrade_cluster_level"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := tkev20180525.NewModifyClusterAttributeRequest()

		request.ClusterId = helper.String(clusterId)

		if v, ok := d.GetOkExists("project_id"); ok {
			request.ProjectId = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("cluster_name"); ok {
			request.ClusterName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cluster_desc"); ok {
			request.ClusterDesc = helper.String(v.(string))
		}

		if err := resourceTencentCloudKubernetesClusterUpdatePostFillRequest0(ctx, request); err != nil {
			return err
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterAttributeWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster failed, reason:%+v", logId, err)
			return err
		}
	}

	needChange1 := false
	mutableArgs1 := []string{"cluster_version"}
	for _, v := range mutableArgs1 {
		if d.HasChange(v) {
			needChange1 = true
			break
		}
	}

	if needChange1 {
		request1 := tkev20180525.NewUpdateClusterVersionRequest()

		response1 := tkev20180525.NewUpdateClusterVersionResponse()

		request1.ClusterId = helper.String(clusterId)

		if v, ok := d.GetOk("cluster_version"); ok {
			request1.DstVersion = helper.String(v.(string))
		}

		if err := resourceTencentCloudKubernetesClusterUpdatePostFillRequest1(ctx, request1); err != nil {
			return err
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().UpdateClusterVersionWithContext(ctx, request1)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request1.GetAction(), request1.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster failed, reason:%+v", logId, err)
			return err
		}
		if err := resourceTencentCloudKubernetesClusterUpdatePostHandleResponse1(ctx, response1); err != nil {
			return err
		}

	}

	// upgrade node version(instances)
	if v, ok := d.GetOkExists("upgrade_instances_follow_cluster"); ok {
		if v.(bool) {
			tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			err := upgradeClusterInstances(tkeService, ctx, clusterId)
			if err != nil {
				return err
			}
		}
	}

	needChange2 := false
	mutableArgs2 := []string{"node_pool_global_config"}
	for _, v := range mutableArgs2 {
		if d.HasChange(v) {
			needChange2 = true
			break
		}
	}

	if needChange2 {
		request2 := tkev20180525.NewModifyClusterAsGroupOptionAttributeRequest()

		request2.ClusterId = helper.String(clusterId)

		if clusterAsGroupOptionMap, ok := helper.InterfacesHeadMap(d, "node_pool_global_config"); ok {
			clusterAsGroupOption := tkev20180525.ClusterAsGroupOption{}
			if v, ok := clusterAsGroupOptionMap["is_scale_in_enabled"]; ok {
				clusterAsGroupOption.IsScaleDownEnabled = helper.Bool(v.(bool))
			}
			if v, ok := clusterAsGroupOptionMap["expander"]; ok {
				clusterAsGroupOption.Expander = helper.String(v.(string))
			}
			if v, ok := clusterAsGroupOptionMap["max_concurrent_scale_in"]; ok {
				clusterAsGroupOption.MaxEmptyBulkDelete = helper.IntInt64(v.(int))
			}
			if v, ok := clusterAsGroupOptionMap["scale_in_delay"]; ok {
				clusterAsGroupOption.ScaleDownDelay = helper.IntInt64(v.(int))
			}
			if v, ok := clusterAsGroupOptionMap["scale_in_unneeded_time"]; ok {
				clusterAsGroupOption.ScaleDownUnneededTime = helper.IntInt64(v.(int))
			}
			if v, ok := clusterAsGroupOptionMap["scale_in_utilization_threshold"]; ok {
				clusterAsGroupOption.ScaleDownUtilizationThreshold = helper.IntInt64(v.(int))
			}
			if v, ok := clusterAsGroupOptionMap["ignore_daemon_sets_utilization"]; ok {
				clusterAsGroupOption.IgnoreDaemonSetsUtilization = helper.Bool(v.(bool))
			}
			if v, ok := clusterAsGroupOptionMap["skip_nodes_with_local_storage"]; ok {
				clusterAsGroupOption.SkipNodesWithLocalStorage = helper.Bool(v.(bool))
			}
			if v, ok := clusterAsGroupOptionMap["skip_nodes_with_system_pods"]; ok {
				clusterAsGroupOption.SkipNodesWithSystemPods = helper.Bool(v.(bool))
			}
			request2.ClusterAsGroupOption = &clusterAsGroupOption
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterAsGroupOptionAttributeWithContext(ctx, request2)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request2.GetAction(), request2.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster failed, reason:%+v", logId, err)
			return err
		}
	}

	needChange3 := false
	mutableArgs3 := []string{"cluster_os"}
	for _, v := range mutableArgs3 {
		if d.HasChange(v) {
			needChange3 = true
			break
		}
	}

	if needChange3 {
		request3 := tkev20180525.NewModifyClusterImageRequest()

		request3.ClusterId = helper.String(clusterId)

		if v, ok := d.GetOk("cluster_os"); ok {
			request3.ImageId = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterImageWithContext(ctx, request3)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request3.GetAction(), request3.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes cluster failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("exist_instance") {
		tkeService := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		cvmService := svccvm.NewCvmService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

		oldInterface, newInterface := d.GetChange("exist_instance")
		oldInstances := oldInterface.(*schema.Set)
		newInstances := newInterface.(*schema.Set)

		remove := oldInstances.Difference(newInstances).List()
		add := newInstances.Difference(oldInstances).List()

		// scale out first
		if len(add) > 0 {
			tmpNew := make([]*tkev20180525.ExistedInstancesForNode, 0, len(add))
			instanceIds := make([]*string, 0)
			instanceInfo := make([]map[string]interface{}, 0)
			for index := range add {
				if add[index] != nil {
					instance := add[index].(map[string]interface{})
					existedInstance, _ := tkeGetCvmExistInstancesPara(instance)
					tmpNew = append(tmpNew, &existedInstance)

					// get all new cvm IDs
					if len(existedInstance.ExistedInstancesPara.InstanceIds) > 0 {
						dMap := make(map[string]interface{}, 0)
						instanceIds = append(instanceIds, existedInstance.ExistedInstancesPara.InstanceIds...)
						dMap["instance_ids"] = instanceIds
						dMap["node_role"] = existedInstance.NodeRole
						instanceInfo = append(instanceInfo, dMap)
					}
				}
			}

			if len(tmpNew) > 0 {
				request := tkev20180525.NewScaleOutClusterMasterRequest()
				request.ClusterId = &clusterId
				request.ExistedInstancesForNode = tmpNew
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ScaleOutClusterMasterWithContext(ctx, request)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s scale out cluster failed, reason:%+v", logId, err)
					return err
				}

				// wait for cvm status
				err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := cvmService.DescribeInstanceByFilter(ctx, instanceIds, nil)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					initFlag := true
					if result != nil {
						for _, item := range result {
							if item.InstanceState != nil && *item.InstanceState != "RUNNING" {
								initFlag = false
								break
							}
						}

						if initFlag {
							return nil
						}
					}

					return resource.RetryableError(fmt.Errorf("cvm instance status is not RUNNING, retry..."))
				})

				if err != nil {
					return err
				}

				// wait for tke node init
				for _, item := range instanceInfo {
					tmpInsIds := item["instance_ids"].([]*string)
					nodeRole := item["node_role"].(*string)
					err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
						result, e := tkeService.DescribeKubernetesClusterMasterAttachmentByIds(ctx, clusterId, tmpInsIds, nodeRole)
						if e != nil {
							return tccommon.RetryError(e, tccommon.InternalError)
						}

						initFlag := true
						if result != nil && result.InstanceSet != nil {
							for _, item := range result.InstanceSet {
								if item.InstanceState != nil && *item.InstanceState != "running" {
									initFlag = false
									break
								}
							}

							if initFlag {
								return nil
							}
						}

						return resource.RetryableError(fmt.Errorf("tke master node cvm instance status is not running, retry..."))
					})

					if err != nil {
						return err
					}
				}

				// wait for tke cluster status
				err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := tkeService.DescribeKubernetesClusterById(ctx, clusterId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					if result == nil {

					}

					if result.ClusterStatus != nil && *result.ClusterStatus == "Running" {
						return nil
					}

					return resource.RetryableError(fmt.Errorf("tke status is not RUNNING, retry..."))
				})

				if err != nil {
					return err
				}

			}
		}

		// scale in
		if len(remove) > 0 {
			tmpOld := make([]map[string]interface{}, 0)
			for index := range remove {
				if remove[index] != nil {
					instance := remove[index].(map[string]interface{})
					existedInstance, _ := tkeGetCvmExistInstancesPara(instance)

					insMap := make(map[string]interface{})
					if existedInstance.NodeRole != nil {
						insMap["node_role"] = *existedInstance.NodeRole
					}

					if len(existedInstance.ExistedInstancesPara.InstanceIds) > 0 {
						for _, item := range existedInstance.ExistedInstancesPara.InstanceIds {
							if item != nil {
								insMap["instance_id"] = *item
							}

							tmpOld = append(tmpOld, insMap)
						}
					}
				}
			}

			if len(tmpOld) > 0 {
				request := tkev20180525.NewScaleInClusterMasterRequest()
				request.ClusterId = &clusterId
				for _, item := range tmpOld {
					tmp := tkev20180525.ScaleInMaster{}
					if v, ok := item["node_role"].(string); ok && v != "" {
						tmp.NodeRole = &v
					}

					if v, ok := item["instance_id"].(string); ok && v != "" {
						tmp.InstanceId = &v
					}

					tmp.InstanceDeleteMode = helper.String("retain")
					request.ScaleInMasters = append(request.ScaleInMasters, &tmp)
				}

				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ScaleInClusterMasterWithContext(ctx, request)
					if e != nil {
						if sdkErr, ok := e.(*errors.TencentCloudSDKError); ok {
							if sdkErr.GetCode() == "ResourceNotFound" {
								return nil
							}

							if sdkErr.GetCode() == "InvalidParameter" && strings.Contains(sdkErr.GetMessage(), `is not exist`) {
								return nil
							}
						}

						return tccommon.RetryError(e, tccommon.InternalError)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s scale in cluster failed, reason:%+v", logId, err)
					return err
				}

				// wait for tke cluster status
				err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
					result, e := tkeService.DescribeKubernetesClusterById(ctx, clusterId)
					if e != nil {
						return tccommon.RetryError(e, tccommon.InternalError)
					}

					if result == nil {

					}

					if result.ClusterStatus != nil && *result.ClusterStatus == "Running" {
						return nil
					}

					return resource.RetryableError(fmt.Errorf("tke status is not RUNNING, retry..."))
				})

				if err != nil {
					return err
				}
			}
		}
	}

	if err := resourceTencentCloudKubernetesClusterUpdateOnExit(ctx); err != nil {
		return err
	}

	return resourceTencentCloudKubernetesClusterRead(d, meta)
}

func resourceTencentCloudKubernetesClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	clusterId := d.Id()

	var (
		request  = tkev20180525.NewDeleteClusterRequest()
		response = tkev20180525.NewDeleteClusterResponse()
	)

	request.ClusterId = helper.String(clusterId)

	instanceDeleteMode := "terminate"
	if v, ok := d.GetOk("instance_delete_mode"); ok {
		instanceDeleteMode = v.(string)
	}

	request.InstanceDeleteMode = &instanceDeleteMode

	if v, ok := d.GetOk("resource_delete_options"); ok {
		for _, item := range v.(*schema.Set).List() {
			resourceDeleteOptionsMap := item.(map[string]interface{})
			resourceDeleteOption := tkev20180525.ResourceDeleteOption{}
			if v, ok := resourceDeleteOptionsMap["resource_type"]; ok {
				resourceDeleteOption.ResourceType = helper.String(v.(string))
			}
			if v, ok := resourceDeleteOptionsMap["delete_mode"]; ok {
				resourceDeleteOption.DeleteMode = helper.String(v.(string))
			}
			if v, ok := resourceDeleteOptionsMap["skip_deletion_protection"]; ok {
				resourceDeleteOption.SkipDeletionProtection = helper.Bool(v.(bool))
			}
			request.ResourceDeleteOptions = append(request.ResourceDeleteOptions, &resourceDeleteOption)
		}
	}

	if err := resourceTencentCloudKubernetesClusterDeletePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DeleteClusterWithContext(ctx, request)
		if e != nil {
			if err := resourceTencentCloudKubernetesClusterDeleteRequestOnError0(ctx, e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete kubernetes cluster failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	if err := resourceTencentCloudKubernetesClusterDeletePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	return nil
}

func customizeDiffForContainerRuntimeDefault(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// example 1.22.5(maybe 1.22.5-tke.21)
	if clusterVersion, ok := d.GetOk("cluster_version"); ok {
		version := clusterVersion.(string)
		parts := strings.Split(version, ".")
		fmt.Println(parts)
		if len(parts) < 2 {
			log.Printf("[WARN] Invalid cluster version format: %s", version)
			return nil
		}

		mainVersionStr := strings.Split(parts[1], "-")[0]
		mainVersion, err := strconv.Atoi(mainVersionStr)
		fmt.Println(mainVersion)
		if err != nil {
			log.Printf("[WARN] Failed to parse cluster version: %v", err)
			return nil
		}

		runtimeValue := "docker"
		if mainVersion >= 24 {
			runtimeValue = "containerd"
		}

		if _, ok := d.GetOk("container_runtime"); !ok {
			if err := d.SetNew("container_runtime", runtimeValue); err != nil {
				return err
			}
		}
	} else {
		if _, ok := d.GetOk("container_runtime"); !ok {
			if err := d.SetNew("container_runtime", "containerd"); err != nil {
				return err
			}
		}
	}

	return nil
}
