// Code generated by iacg; DO NOT EDIT.
package tke

import (
	"context"
	"fmt"
	"log"
	"strings"

	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tkev20180525 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcas "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/as"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"
)

func ResourceTencentCloudKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesNodePoolCreate,
		Read:   resourceTencentCloudKubernetesNodePoolRead,
		Update: resourceTencentCloudKubernetesNodePoolUpdate,
		Delete: resourceTencentCloudKubernetesNodePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: nodePoolCustomResourceImporter,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
				Description:  "Maximum number of node.",
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 2000),
			},

			"min_size": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Minimum number of node.",
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 2000),
			},

			"desired_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Desired capacity of the node. If `enable_auto_scale` is set `true`, this will be a computed parameter.",
				ValidateFunc: tccommon.ValidateIntegerInRange(0, 2000),
			},

			"enable_auto_scale": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate whether to enable auto scaling or not.",
			},

			"retry_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "IMMEDIATE_RETRY",
				Description:  "Available values for retry policies include `IMMEDIATE_RETRY` and `INCREMENTAL_INTERVALS`.",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{svcas.SCALING_GROUP_RETRY_POLICY_IMMEDIATE_RETRY, svcas.SCALING_GROUP_RETRY_POLICY_INCREMENTAL_INTERVALS, svcas.SCALING_GROUP_RETRY_POLICY_NO_RETRY}),
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
				Description: "ID list of subnet, and for VPC it is required.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"scaling_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Auto scaling mode. Valid values are `CLASSIC_SCALING`(scaling by create/destroy instances), `WAKE_UP_STOPPED_SCALING`(Boot priority for expansion. When expanding the capacity, the shutdown operation is given priority to the shutdown of the instance. If the number of instances is still lower than the expected number of instances after the startup, the instance will be created, and the method of destroying the instance will still be used for shrinking).",
			},

			"multi_zone_subnet_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Multi-availability zone/subnet policy. Valid values: PRIORITY and EQUALITY. Default value: PRIORITY.",
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{svcas.MultiZoneSubnetPolicyPriority, svcas.MultiZoneSubnetPolicyEquality}),
			},

			"node_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Node config.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_target": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Mount target. Default is not mounting.",
						},
						"docker_graph_path": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Default:     "/var/lib/docker",
							Description: "Docker graph path. Default is `/var/lib/docker`.",
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
										Description:  "Types of disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD` and `CLOUD_BSSD`.",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
										ForceNew:    true,
										Default:     0,
										Description: "Volume of disk in GB. Default is `0`.",
									},
									"file_system": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Default:     "",
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
										Default:     "",
										Description: "Mount target.",
									},
									"disk_partition": {
										Type:        schema.TypeString,
										Optional:    true,
										ForceNew:    true,
										Description: "The name of the device or partition to mount. NOTE: this argument doesn't support setting in node pool, or will leads to mount error.",
									},
								},
							},
						},
						"extra_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							Description: "Custom parameter information related to the node. This is a white-list parameter.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Base64-encoded User Data text, the length limit is 16KB.",
						},
						"pre_start_user_script": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Base64-encoded user script, executed before initializing the node, currently only effective for adding existing nodes.",
						},
						"is_schedule": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Default:     true,
							Description: "Indicate to schedule the adding node or not. Default is true.",
						},
						"desired_pod_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Indicate to set desired pod number in node. valid when the cluster is podCIDR.",
						},
						"gpu_args": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "GPU driver parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mig_enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
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
					},
				},
			},

			"auto_scaling_config": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Auto scaling config parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"system_disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "CLOUD_PREMIUM",
							Description:  "Type of a CVM disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD` and `CLOUD_BSSD`. Default is `CLOUD_PREMIUM`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
						},
						"system_disk_size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      50,
							Description:  "Volume of system disk in GB. Default is `50`.",
							ValidateFunc: tccommon.ValidateIntegerInRange(20, 1024),
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
										Default:      "CLOUD_PREMIUM",
										Description:  "Types of disk. Valid value: `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_PREMIUM`, `CLOUD_SSD`, `CLOUD_HSSD`, `CLOUD_TSSD` and `CLOUD_BSSD`.",
										ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.SYSTEM_DISK_ALLOW_TYPE),
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Optional:    true,
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
										Description: "Indicates whether the disk remove after instance terminated. Default is `false`.",
									},
									"encrypt": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Specify whether to encrypt data disk, default: false. NOTE: Make sure the instance type is offering and the cam role `QcloudKMSAccessForCVMRole` was provided.",
									},
									"throughput_performance": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD` and `data_size` > 460GB.",
									},
								},
							},
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID`. The default is `POSTPAID_BY_HOUR`. NOTE: `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time.",
						},
						"instance_charge_type_prepaid_period": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "The tenancy (in month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
							ValidateFunc: tccommon.ValidateAllowedIntValue(svccvm.CVM_PREPAID_PERIOD),
						},
						"instance_charge_type_prepaid_renew_flag": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svccvm.CVM_PREPAID_RENEW_FLAG),
						},
						"spot_instance_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Type of spot instance, only support `one-time` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"one-time"}),
						},
						"spot_max_price": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Max price of a spot instance, is the format of decimal string, for example \"0.50\". Note: it only works when instance_charge_type is set to `SPOTPAID`.",
							ValidateFunc: tccommon.ValidateStringNumber,
						},
						"internet_charge_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "TRAFFIC_POSTPAID_BY_HOUR",
							Description:  "Charge types for network traffic. Valid value: `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
							ValidateFunc: tccommon.ValidateAllowedStringValue(svcas.INTERNET_CHARGE_ALLOW_TYPE),
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
							ForceNew:      true,
							Sensitive:     true,
							ConflictsWith: []string{"auto_scaling_config.0.key_ids"},
							Description:   "Password to access.",
							ValidateFunc:  tccommon.ValidateAsConfigPassword,
						},
						"key_ids": {
							Type:          schema.TypeList,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"auto_scaling_config.0.password"},
							Description:   "ID list of keys.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"security_group_ids": {
							Type:          schema.TypeSet,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"auto_scaling_config.0.orderly_security_group_ids"},
							Description:   "Security groups to which a CVM instance belongs.",
							Deprecated:    "The order of elements in this field cannot be guaranteed. Use `orderly_security_group_ids` instead.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"orderly_security_group_ids": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"auto_scaling_config.0.security_group_ids"},
							Description:   "Ordered security groups to which a CVM instance belongs.",
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
							ForceNew:    true,
							Default:     true,
							Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
						},
						"cam_role_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Name of cam role.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Instance name, no more than 60 characters. For usage, refer to `InstanceNameSettings` in https://www.tencentcloud.com/document/product/377/31001.",
						},
						"instance_name_style": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Type of CVM instance name. Valid values: `ORIGINAL` and `UNIQUE`. Default value: `ORIGINAL`. For usage, refer to `InstanceNameSettings` in https://www.tencentcloud.com/document/product/377/31001.",
						},
						"host_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The hostname of the cloud server, dot (.) and dash (-) cannot be used as the first and last characters of HostName and cannot be used consecutively. Windows instances are not supported. Examples of other types (Linux, etc.): The character length is [2, 40], multiple periods are allowed, and there is a paragraph between the dots, and each paragraph is allowed to consist of letters (unlimited case), numbers and dashes (-). Pure numbers are not allowed. For usage, refer to `HostNameSettings` in https://www.tencentcloud.com/document/product/377/31001.",
						},
						"host_name_style": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The style of the host name of the cloud server, the value range includes ORIGINAL and UNIQUE, and the default is ORIGINAL. For usage, refer to `HostNameSettings` in https://www.tencentcloud.com/document/product/377/31001.",
						},
					},
				},
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
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Taints of kubernetes node pool created nodes.",
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
			},

			"delete_keep_instance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicate to keep the CVM instance when delete the node pool. Default is `true`.",
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether the node pool deletion protection is enabled.",
			},

			"node_os": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tlinux2.4x86_64",
				Description: "Operating system of the cluster. Please refer to [TencentCloud Documentation](https://www.tencentcloud.com/document/product/457/46750?lang=en&pg=#list-of-public-images-supported-by-tke) for available values. Default is 'tlinux2.4x86_64'. This parameter will only affect new nodes, not including the existing nodes.",
			},

			"node_os_type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "GENERAL",
				Description:      "The image version of the node. Valida values are `DOCKER_CUSTOMIZE` and `GENERAL`. Default is `GENERAL`. This parameter will only affect new nodes, not including the existing nodes.",
				DiffSuppressFunc: nodeOsTypeDiffSuppressFunc,
			},

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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Policy of scaling group termination. Available values: `[\"OLDEST_INSTANCE\"]`, `[\"NEWEST_INSTANCE\"]`.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Node pool tag specifications, will passthroughs to the scaling instances.",
			},

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
	}
}

func resourceTencentCloudKubernetesNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_node_pool.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		clusterId  string
		nodePoolId string
	)
	var (
		request  = tkev20180525.NewCreateClusterNodePoolRequest()
		response = tkev20180525.NewCreateClusterNodePoolResponse()
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	request.ClusterId = helper.String(clusterId)

	if v, ok := d.GetOkExists("enable_auto_scale"); ok {
		request.EnableAutoscale = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("taints"); ok {
		for _, item := range v.([]interface{}) {
			taintsMap := item.(map[string]interface{})
			taint := tkev20180525.Taint{}
			if v, ok := taintsMap["key"]; ok {
				taint.Key = helper.String(v.(string))
			}
			if v, ok := taintsMap["value"]; ok {
				taint.Value = helper.String(v.(string))
			}
			if v, ok := taintsMap["effect"]; ok {
				taint.Effect = helper.String(v.(string))
			}
			request.Taints = append(request.Taints, &taint)
		}
	}

	if v, ok := d.GetOkExists("deletion_protection"); ok {
		request.DeletionProtection = helper.Bool(v.(bool))
	}

	if err := resourceTencentCloudKubernetesNodePoolCreatePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().CreateClusterNodePoolWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create kubernetes node pool failed, reason:%+v", logId, err)
		return err
	}

	nodePoolId = *response.Response.NodePoolId

	if err := resourceTencentCloudKubernetesNodePoolCreatePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{clusterId, nodePoolId}, tccommon.FILED_SP))

	return resourceTencentCloudKubernetesNodePoolRead(d, meta)
}

func resourceTencentCloudKubernetesNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_node_pool.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	nodePoolId := idSplit[1]

	_ = d.Set("cluster_id", clusterId)

	respData, err := service.DescribeKubernetesNodePoolById(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_node_pool` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	var respData1 *tkev20180525.NodePool
	reqErr1 := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKubernetesNodePoolById1(ctx, clusterId, nodePoolId)
		if e != nil {
			if err := resourceTencentCloudKubernetesNodePoolReadRequestOnError1(ctx, result, e); err != nil {
				return err
			}
			return tccommon.RetryError(e)
		}
		if err := resourceTencentCloudKubernetesNodePoolReadRequestOnSuccess1(ctx, result); err != nil {
			return err
		}
		respData1 = result
		return nil
	})
	if reqErr1 != nil {
		log.Printf("[CRITAL]%s read kubernetes node pool failed, reason:%+v", logId, reqErr1)
		return reqErr1
	}

	if respData1 == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `kubernetes_node_pool` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if respData1.Name != nil {
		_ = d.Set("name", respData1.Name)
	}

	if respData1.LifeState != nil {
		_ = d.Set("status", respData1.LifeState)
	}

	if respData1.LaunchConfigurationId != nil {
		_ = d.Set("launch_config_id", respData1.LaunchConfigurationId)
	}

	if respData1.AutoscalingGroupId != nil {
		_ = d.Set("auto_scaling_group_id", respData1.AutoscalingGroupId)
	}

	taintsList := make([]map[string]interface{}, 0, len(respData1.Taints))
	if respData1.Taints != nil {
		for _, taints := range respData1.Taints {
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

		_ = d.Set("taints", taintsList)
	}

	if respData1.NodeCountSummary != nil {
		if respData1.NodeCountSummary.ManuallyAdded != nil {
			if respData1.NodeCountSummary.ManuallyAdded.Total != nil {
				_ = d.Set("manually_added_total", respData1.NodeCountSummary.ManuallyAdded.Total)
			}

		}

		if respData1.NodeCountSummary.AutoscalingAdded != nil {
			if respData1.NodeCountSummary.AutoscalingAdded.Total != nil {
				_ = d.Set("autoscaling_added_total", respData1.NodeCountSummary.AutoscalingAdded.Total)
			}

		}

	}

	if respData1.MaxNodesNum != nil {
		_ = d.Set("max_size", respData1.MaxNodesNum)
	}

	if respData1.MinNodesNum != nil {
		_ = d.Set("min_size", respData1.MinNodesNum)
	}

	if respData1.DesiredNodesNum != nil {
		_ = d.Set("desired_capacity", respData1.DesiredNodesNum)
	}

	if respData1.DeletionProtection != nil {
		_ = d.Set("deletion_protection", respData1.DeletionProtection)
	}

	if err := resourceTencentCloudKubernetesNodePoolReadPostHandleResponse1(ctx, respData1); err != nil {
		return err
	}

	return nil
}

func resourceTencentCloudKubernetesNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_node_pool.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	nodePoolId := idSplit[1]

	if err := resourceTencentCloudKubernetesNodePoolUpdateOnStart(ctx); err != nil {
		return err
	}

	needChange := false
	mutableArgs := []string{"name", "max_size", "min_size", "enable_auto_scale", "deletion_protection"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := tkev20180525.NewModifyClusterNodePoolRequest()

		request.ClusterId = helper.String(clusterId)

		request.NodePoolId = helper.String(nodePoolId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("max_size"); ok {
			request.MaxNodesNum = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("min_size"); ok {
			request.MinNodesNum = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOkExists("enable_auto_scale"); ok {
			request.EnableAutoscale = helper.Bool(v.(bool))
		}

		if v, ok := d.GetOkExists("deletion_protection"); ok {
			request.DeletionProtection = helper.Bool(v.(bool))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterNodePoolWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update kubernetes node pool failed, reason:%+v", logId, err)
			return err
		}
	}

	if d.HasChange("taints") {
		_, n := d.GetChange("taints")

		// clean taints
		if len(n.([]interface{})) == 0 {
			body := map[string]interface{}{
				"ClusterId":  clusterId,
				"NodePoolId": nodePoolId,
				"Taints":     []interface{}{},
			}

			client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOmitNilClient("tke")
			request := tchttp.NewCommonRequest("tke", "2018-05-25", "ModifyClusterNodePool")
			err := request.SetActionParameters(body)
			if err != nil {
				return err
			}

			response := tchttp.NewCommonResponse()
			err = client.Send(request, response)
			if err != nil {
				fmt.Printf("update kubernetes node pool taints failed: %v \n", err)
				return err
			}
		} else {
			request := tkev20180525.NewModifyClusterNodePoolRequest()
			request.ClusterId = helper.String(clusterId)
			request.NodePoolId = helper.String(nodePoolId)

			if v, ok := d.GetOk("taints"); ok {
				for _, item := range v.([]interface{}) {
					taintsMap := item.(map[string]interface{})
					taint := tkev20180525.Taint{}
					if v, ok := taintsMap["key"]; ok {
						taint.Key = helper.String(v.(string))
					}

					if v, ok := taintsMap["value"]; ok {
						taint.Value = helper.String(v.(string))
					}

					if v, ok := taintsMap["effect"]; ok {
						taint.Effect = helper.String(v.(string))
					}

					request.Taints = append(request.Taints, &taint)
				}
			}

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterNodePoolWithContext(ctx, request)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update kubernetes node pool taints failed, reason:%+v", logId, err)
				return err
			}
		}

		return nil
	}

	if err := resourceTencentCloudKubernetesNodePoolUpdateOnExit(ctx); err != nil {
		return err
	}

	return resourceTencentCloudKubernetesNodePoolRead(d, meta)
}

func resourceTencentCloudKubernetesNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_node_pool.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	clusterId := idSplit[0]
	nodePoolId := idSplit[1]

	var (
		request  = tkev20180525.NewDeleteClusterNodePoolRequest()
		response = tkev20180525.NewDeleteClusterNodePoolResponse()
	)

	request.ClusterId = helper.String(clusterId)

	request.NodePoolIds = []*string{helper.String(nodePoolId)}

	if v, ok := d.GetOkExists("delete_keep_instance"); ok {
		request.KeepInstance = helper.Bool(v.(bool))
	}

	if err := resourceTencentCloudKubernetesNodePoolDeletePostFillRequest0(ctx, request); err != nil {
		return err
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DeleteClusterNodePoolWithContext(ctx, request)
		if e != nil {
			if err := resourceTencentCloudKubernetesNodePoolDeleteRequestOnError0(ctx, e); err != nil {
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
		log.Printf("[CRITAL]%s delete kubernetes node pool failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	if err := resourceTencentCloudKubernetesNodePoolDeletePostHandleResponse0(ctx, response); err != nil {
		return err
	}

	return nil
}
