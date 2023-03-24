/*
Provides a resource to create a cvm launch_template_version

Example Usage

```hcl
resource "tencentcloud_cvm_launch_template_version" "foo" {
  placement {
		zone = "ap-guangzhou-6"
		project_id = 0

  }
  launch_template_id = "lt-r9ajalbi"
  launch_template_version_description = "version description"
  disable_api_termination = false
  instance_type = "S5.MEDIUM4"
  image_id = "img-9qrfy1xt"
}
```

Import

cvm launch_template_version can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_launch_template_version.launch_template_version ${launch_template_id}#${launch_template_version}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCvmLaunchTemplateVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmLaunchTemplateVersionCreate,
		Read:   resourceTencentCloudCvmLaunchTemplateVersionRead,
		Delete: resourceTencentCloudCvmLaunchTemplateVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"placement": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Location of the instance. You can use this parameter to specify the attributes of the instance, such as its availability zone, project, and CDH (for dedicated CVMs).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "ID of the availability zone where the instance resides. You can call the DescribeZones API and obtain the ID in the returned Zone field.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID of the project to which the instance belongs. This parameter can be obtained from the projectId returned by DescribeProject. If this is left empty, the default project is used.",
						},
						"host_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID list of CDHs from which the instance can be created. If you have purchased CDHs and specify this parameter, the instances you purchase will be randomly deployed on the CDHs.",
						},
						"host_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "IPs of the hosts to create CVMs.",
						},
					},
				},
			},

			"launch_template_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance launch template ID. This parameter is used as a basis for creating new template versions.",
			},

			"launch_template_version": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "This parameter, when specified, is used to create instance launch templates. If this parameter is not specified, the default version will be used.",
			},

			"launch_template_version_description": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Description of instance launch template versions. This parameter can contain 2-256 characters.",
			},

			"instance_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The type of the instance. If this parameter is not specified, the system will dynamically specify the default model according to the resource sales in the current region.",
			},

			"image_id": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Image ID.",
			},

			"system_disk": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "System disk configuration information of the instance. If this parameter is not specified, it is assigned according to the system default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The type of system disk. Default value: the type of hard disk currently in stock.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "System disk ID. System disks whose type is LOCAL_BASIC or LOCAL_SSD do not have an ID and do not support this parameter. It is only used as a response parameter for APIs such as DescribeInstances, and cannot be used as a request parameter for APIs such as RunInstances.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "System disk size; unit: GB; default value: 50 GB.",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID of the dedicated cluster to which the instance belongs.",
						},
					},
				},
			},

			"data_disks": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "The configuration information of instance data disks. If this parameter is not specified, no data disk will be purchased by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "Data disk size (in GB). The minimum adjustment increment is 10 GB. The value range varies by data disk type.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The type of data disk.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "System disk ID. System disks whose type is LOCAL_BASIC or LOCAL_SSD do not have an ID and do not support this parameter. It is only used as a response parameter for APIs such as DescribeInstances, and cannot be used as a request parameter for APIs such as RunInstances.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Whether to terminate the data disk when its CVM is terminated. Default value: `true`.",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Data disk snapshot ID. The size of the selected data disk snapshot must be smaller than that of the data disk. Note: This field may return null, indicating that no valid value is found.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Specifies whether the data disk is encrypted.",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID of the custom CMK in the format of UUID or `kms-abcd1234`.",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Cloud disk performance in MB/s.",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "ID of the dedicated cluster to which the instance belongs.",
						},
					},
				},
			},

			"virtual_private_cloud": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes information on VPC, including subnets, IP addresses, etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC ID in the format of vpc-xxx, if you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "VPC subnet ID in the format subnet-xxx, if you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.",
						},
						"as_vpc_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC.",
						},
						"private_ip_addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Number of IPv6 addresses randomly generated for the ENI.",
						},
					},
				},
			},

			"internet_accessible": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes the accessibility of an instance in the public network, including its network billing method, maximum bandwidth, etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"internet_charge_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Network connection billing plan.",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "The maximum outbound bandwidth of the public network, in Mbps. The default value is 0 Mbps.",
						},
						"public_ip_assigned": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Whether to assign a public IP.",
						},
						"bandwidth_package_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Bandwidth package ID.",
						},
					},
				},
			},

			"instance_count": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The number of instances to be purchased.",
			},

			"instance_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance name to be displayed.",
			},

			"login_settings": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes login settings of an instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Login password of the instance.",
						},
						"key_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "List of key IDs. After an instance is associated with a key, you can access the instance with the private key in the key pair.",
						},
						"keep_image_login": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Whether to keep the original settings of an image.",
						},
					},
				},
			},

			"security_group_ids": {
				Optional: true,
				Computed: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security groups to which the instance belongs. If this parameter is not specified, the instance will be associated with default security groups.",
			},

			"enhanced_service": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Enhanced service. You can use this parameter to specify whether to enable services such as Anti-DDoS and Cloud Monitor. If this parameter is not specified, Cloud Monitor and Anti-DDoS are enabled for public images by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Enables cloud security service. If this parameter is not specified, the cloud security service will be enabled by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Whether to enable Cloud Security.",
									},
								},
							},
						},
						"monitor_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Enables cloud monitor service. If this parameter is not specified, the cloud monitor service will be enabled by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Whether to enable Cloud Monitor.",
									},
								},
							},
						},
						"automation_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Whether to enable the TAT service. If this parameter is not specified, the TAT service is enabled for public images and disabled for other images by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Whether to enable the TAT service.",
									},
								},
							},
						},
					},
				},
			},

			"client_token": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idem-potency of the request cannot be guaranteed.",
			},

			"host_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Hostname of a CVM.",
			},

			"action_timer": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Scheduled tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timer_action": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Timer name. Currently TerminateInstances is the only supported value.",
						},
						"action_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Execution time, displayed according to ISO8601 standard, and UTC time is used. The format is YYYY-MM-DDThh:mm:ssZ. For example, 2018-05-29T11:26:40Z, the execution must be at least 5 minutes later than the current time.",
						},
						"externals": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Additional data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"release_address": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Release address.",
									},
									"unsupport_networks": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Not supported network.",
									},
									"storage_block_attr": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Information on local HDD storage.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "Local HDD storage type. Value: LOCAL_PRO.",
												},
												"min_size": {
													Type:        schema.TypeInt,
													Required:    true,
													ForceNew:    true,
													Description: "Minimum capacity of local HDD storage.",
												},
												"max_size": {
													Type:        schema.TypeInt,
													Required:    true,
													ForceNew:    true,
													Description: "Maximum capacity of local HDD storage.",
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

			"disaster_recover_group_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Placement group ID. You can only specify one.",
			},

			"tag_specification": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Description of tags associated with resource instances during instance creation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The type of resource that the tag is bound to.",
						},
						"tags": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Tag value.",
									},
								},
							},
						},
					},
				},
			},

			"instance_market_options": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Options related to bidding requests.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spot_options": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							ForceNew:    true,
							Description: "Options related to bidding.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_price": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Bidding price.",
									},
									"spot_instance_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										ForceNew:    true,
										Description: "Bidding request type. Currently only one-time is supported.",
									},
								},
							},
						},
						"market_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Market option type. Currently spot is the only supported value.",
						},
					},
				},
			},

			"user_data": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16 KB.",
			},

			"dry_run": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether the request is a dry run only.",
			},

			"cam_role_name": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The role name of CAM.",
			},

			"hpc_cluster_id": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "HPC cluster ID. The HPC cluster must and can only be specified for a high-performance computing instance.",
			},

			"instance_charge_type": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The charge type of instance.",
			},

			"instance_charge_prepaid": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes the billing method of an instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "Subscription period; unit: month; valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: "Auto renewal flag. Valid values: NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically &lt;br&gt;&lt;br&gt;Default value: NOTIFY_AND_MANUAL_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient.",
						},
					},
				},
			},

			"disable_api_termination": {
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether the termination protection is enabled. `TRUE`: Enable instance protection, which means that this instance can not be deleted by an API action.`FALSE`: Do not enable the instance protection. Default value: `FALSE`.",
			},
		},
	}
}

func resourceTencentCloudCvmLaunchTemplateVersionCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_version.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = cvm.NewCreateLaunchTemplateVersionRequest()
		response         = cvm.NewCreateLaunchTemplateVersionResponse()
		launchTemplateId string
	)
	if dMap, ok := helper.InterfacesHeadMap(d, "placement"); ok {
		placement := cvm.Placement{}
		if v, ok := dMap["zone"]; ok {
			placement.Zone = helper.String(v.(string))
		}
		if v, ok := dMap["project_id"]; ok {
			placement.ProjectId = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["host_ids"]; ok {
			hostIdsSet := v.(*schema.Set).List()
			for i := range hostIdsSet {
				hostIds := hostIdsSet[i].(string)
				placement.HostIds = append(placement.HostIds, &hostIds)
			}
		}
		if v, ok := dMap["host_ips"]; ok {
			hostIpsSet := v.(*schema.Set).List()
			for i := range hostIpsSet {
				hostIps := hostIpsSet[i].(string)
				placement.HostIps = append(placement.HostIps, &hostIps)
			}
		}
		request.Placement = &placement
	}

	if v, ok := d.GetOk("launch_template_id"); ok {
		launchTemplateId = v.(string)
		request.LaunchTemplateId = helper.String(launchTemplateId)
	}

	if v, ok := d.GetOkExists("launch_template_version"); ok {
		request.LaunchTemplateVersion = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("launch_template_version_description"); ok {
		request.LaunchTemplateVersionDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "system_disk"); ok {
		systemDisk := cvm.SystemDisk{}
		if v, ok := dMap["disk_type"]; ok {
			systemDisk.DiskType = helper.String(v.(string))
		}
		if v, ok := dMap["disk_id"]; ok {
			systemDisk.DiskId = helper.String(v.(string))
		}
		if v, ok := dMap["disk_size"]; ok {
			systemDisk.DiskSize = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["cdc_id"]; ok {
			systemDisk.CdcId = helper.String(v.(string))
		}
		request.SystemDisk = &systemDisk
	}

	if v, ok := d.GetOk("data_disks"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			dataDisk := cvm.DataDisk{}
			if v, ok := dMap["disk_size"]; ok {
				dataDisk.DiskSize = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["disk_type"]; ok {
				dataDisk.DiskType = helper.String(v.(string))
			}
			if v, ok := dMap["disk_id"]; ok {
				dataDisk.DiskId = helper.String(v.(string))
			}
			if v, ok := dMap["delete_with_instance"]; ok {
				dataDisk.DeleteWithInstance = helper.Bool(v.(bool))
			}
			if v, ok := dMap["snapshot_id"]; ok {
				dataDisk.SnapshotId = helper.String(v.(string))
			}
			if v, ok := dMap["encrypt"]; ok {
				dataDisk.Encrypt = helper.Bool(v.(bool))
			}
			if v, ok := dMap["kms_key_id"]; ok {
				dataDisk.KmsKeyId = helper.String(v.(string))
			}
			if v, ok := dMap["throughput_performance"]; ok {
				dataDisk.ThroughputPerformance = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["cdc_id"]; ok {
				dataDisk.CdcId = helper.String(v.(string))
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "virtual_private_cloud"); ok {
		virtualPrivateCloud := cvm.VirtualPrivateCloud{}
		if v, ok := dMap["vpc_id"]; ok {
			virtualPrivateCloud.VpcId = helper.String(v.(string))
		}
		if v, ok := dMap["subnet_id"]; ok {
			virtualPrivateCloud.SubnetId = helper.String(v.(string))
		}
		if v, ok := dMap["as_vpc_gateway"]; ok {
			virtualPrivateCloud.AsVpcGateway = helper.Bool(v.(bool))
		}
		if v, ok := dMap["private_ip_addresses"]; ok {
			privateIpAddressesSet := v.(*schema.Set).List()
			for i := range privateIpAddressesSet {
				privateIpAddresses := privateIpAddressesSet[i].(string)
				virtualPrivateCloud.PrivateIpAddresses = append(virtualPrivateCloud.PrivateIpAddresses, &privateIpAddresses)
			}
		}
		if v, ok := dMap["ipv6_address_count"]; ok {
			virtualPrivateCloud.Ipv6AddressCount = helper.IntUint64(v.(int))
		}
		request.VirtualPrivateCloud = &virtualPrivateCloud
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "internet_accessible"); ok {
		internetAccessible := cvm.InternetAccessible{}
		if v, ok := dMap["internet_charge_type"]; ok {
			internetAccessible.InternetChargeType = helper.String(v.(string))
		}
		if v, ok := dMap["internet_max_bandwidth_out"]; ok {
			internetAccessible.InternetMaxBandwidthOut = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["public_ip_assigned"]; ok {
			internetAccessible.PublicIpAssigned = helper.Bool(v.(bool))
		}
		if v, ok := dMap["bandwidth_package_id"]; ok {
			internetAccessible.BandwidthPackageId = helper.String(v.(string))
		}
		request.InternetAccessible = &internetAccessible
	}

	if v, ok := d.GetOkExists("instance_count"); ok {
		request.InstanceCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "login_settings"); ok {
		loginSettings := cvm.LoginSettings{}
		if v, ok := dMap["password"]; ok {
			loginSettings.Password = helper.String(v.(string))
		}
		if v, ok := dMap["key_ids"]; ok {
			keyIdsSet := v.(*schema.Set).List()
			for i := range keyIdsSet {
				keyIds := keyIdsSet[i].(string)
				loginSettings.KeyIds = append(loginSettings.KeyIds, &keyIds)
			}
		}
		if v, ok := dMap["keep_image_login"]; ok {
			loginSettings.KeepImageLogin = helper.String(v.(string))
		}
		request.LoginSettings = &loginSettings
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroupIdsSet := v.(*schema.Set).List()
		for i := range securityGroupIdsSet {
			securityGroupIds := securityGroupIdsSet[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroupIds)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "enhanced_service"); ok {
		enhancedService := cvm.EnhancedService{}
		if securityServiceMap, ok := helper.InterfaceToMap(dMap, "security_service"); ok {
			runSecurityServiceEnabled := cvm.RunSecurityServiceEnabled{}
			if v, ok := securityServiceMap["enabled"]; ok {
				runSecurityServiceEnabled.Enabled = helper.Bool(v.(bool))
			}
			enhancedService.SecurityService = &runSecurityServiceEnabled
		}
		if monitorServiceMap, ok := helper.InterfaceToMap(dMap, "monitor_service"); ok {
			runMonitorServiceEnabled := cvm.RunMonitorServiceEnabled{}
			if v, ok := monitorServiceMap["enabled"]; ok {
				runMonitorServiceEnabled.Enabled = helper.Bool(v.(bool))
			}
			enhancedService.MonitorService = &runMonitorServiceEnabled
		}
		if automationServiceMap, ok := helper.InterfaceToMap(dMap, "automation_service"); ok {
			runAutomationServiceEnabled := cvm.RunAutomationServiceEnabled{}
			if v, ok := automationServiceMap["enabled"]; ok {
				runAutomationServiceEnabled.Enabled = helper.Bool(v.(bool))
			}
			enhancedService.AutomationService = &runAutomationServiceEnabled
		}
		request.EnhancedService = &enhancedService
	}

	if v, ok := d.GetOk("client_token"); ok {
		request.ClientToken = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host_name"); ok {
		request.HostName = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "action_timer"); ok {
		actionTimer := cvm.ActionTimer{}
		if v, ok := dMap["timer_action"]; ok {
			actionTimer.TimerAction = helper.String(v.(string))
		}
		if v, ok := dMap["action_time"]; ok {
			actionTimer.ActionTime = helper.String(v.(string))
		}
		if externalsMap, ok := helper.InterfaceToMap(dMap, "externals"); ok {
			externals := cvm.Externals{}
			if v, ok := externalsMap["release_address"]; ok {
				externals.ReleaseAddress = helper.Bool(v.(bool))
			}
			if v, ok := externalsMap["unsupport_networks"]; ok {
				unsupportNetworksSet := v.(*schema.Set).List()
				for i := range unsupportNetworksSet {
					unsupportNetworks := unsupportNetworksSet[i].(string)
					externals.UnsupportNetworks = append(externals.UnsupportNetworks, &unsupportNetworks)
				}
			}
			if storageBlockAttrMap, ok := helper.InterfaceToMap(externalsMap, "storage_block_attr"); ok {
				storageBlock := cvm.StorageBlock{}
				if v, ok := storageBlockAttrMap["type"]; ok {
					storageBlock.Type = helper.String(v.(string))
				}
				if v, ok := storageBlockAttrMap["min_size"]; ok {
					storageBlock.MinSize = helper.IntInt64(v.(int))
				}
				if v, ok := storageBlockAttrMap["max_size"]; ok {
					storageBlock.MaxSize = helper.IntInt64(v.(int))
				}
				externals.StorageBlockAttr = &storageBlock
			}
			actionTimer.Externals = &externals
		}
		request.ActionTimer = &actionTimer
	}

	if v, ok := d.GetOk("disaster_recover_group_ids"); ok {
		disasterRecoverGroupIdsSet := v.(*schema.Set).List()
		for i := range disasterRecoverGroupIdsSet {
			disasterRecoverGroupIds := disasterRecoverGroupIdsSet[i].(string)
			request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &disasterRecoverGroupIds)
		}
	}

	if v, ok := d.GetOk("tag_specification"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			tagSpecification := cvm.TagSpecification{}
			if v, ok := dMap["resource_type"]; ok {
				tagSpecification.ResourceType = helper.String(v.(string))
			}
			if v, ok := dMap["tags"]; ok {
				for _, item := range v.([]interface{}) {
					tagsMap := item.(map[string]interface{})
					tag := cvm.Tag{}
					if v, ok := tagsMap["key"]; ok {
						tag.Key = helper.String(v.(string))
					}
					if v, ok := tagsMap["value"]; ok {
						tag.Value = helper.String(v.(string))
					}
					tagSpecification.Tags = append(tagSpecification.Tags, &tag)
				}
			}
			request.TagSpecification = append(request.TagSpecification, &tagSpecification)
		}
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_market_options"); ok {
		instanceMarketOptionsRequest := cvm.InstanceMarketOptionsRequest{}
		if spotOptionsMap, ok := helper.InterfaceToMap(dMap, "spot_options"); ok {
			spotMarketOptions := cvm.SpotMarketOptions{}
			if v, ok := spotOptionsMap["max_price"]; ok {
				spotMarketOptions.MaxPrice = helper.String(v.(string))
			}
			if v, ok := spotOptionsMap["spot_instance_type"]; ok {
				spotMarketOptions.SpotInstanceType = helper.String(v.(string))
			}
			instanceMarketOptionsRequest.SpotOptions = &spotMarketOptions
		}
		if v, ok := dMap["market_type"]; ok {
			instanceMarketOptionsRequest.MarketType = helper.String(v.(string))
		}
		request.InstanceMarketOptions = &instanceMarketOptionsRequest
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("dry_run"); ok {
		request.DryRun = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("cam_role_name"); ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("hpc_cluster_id"); ok {
		request.HpcClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "instance_charge_prepaid"); ok {
		instanceChargePrepaid := cvm.InstanceChargePrepaid{}
		if v, ok := dMap["period"]; ok {
			instanceChargePrepaid.Period = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["renew_flag"]; ok {
			instanceChargePrepaid.RenewFlag = helper.String(v.(string))
		}
		request.InstanceChargePrepaid = &instanceChargePrepaid
	}

	if v, ok := d.GetOkExists("disable_api_termination"); ok {
		request.DisableApiTermination = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().CreateLaunchTemplateVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create cvm launchTemplateVersion failed, reason:%+v", logId, err)
		return err
	}

	launchTemplateVersionNumber := *response.Response.LaunchTemplateVersionNumber
	launchTemplateVersionNumberString := strconv.FormatInt(launchTemplateVersionNumber, 10)
	d.SetId(launchTemplateId + FILED_SP + launchTemplateVersionNumberString)

	return resourceTencentCloudCvmLaunchTemplateVersionRead(d, meta)
}

func resourceTencentCloudCvmLaunchTemplateVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	launchTemplateId := idSplit[0]
	launchTemplateVersionNumber := idSplit[1]

	launchTemplateVersion, err := service.DescribeCvmLaunchTemplateVersionById(ctx, launchTemplateId, launchTemplateVersionNumber)
	if err != nil {
		return err
	}

	if launchTemplateVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmLaunchTemplateVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if launchTemplateVersion.LaunchTemplateId != nil {
		_ = d.Set("launch_template_id", launchTemplateVersion.LaunchTemplateId)
	}

	if launchTemplateVersion.LaunchTemplateVersion != nil {
		_ = d.Set("launch_template_version", launchTemplateVersion.LaunchTemplateVersion)
	}

	if launchTemplateVersion.LaunchTemplateVersionDescription != nil {
		_ = d.Set("launch_template_version_description", launchTemplateVersion.LaunchTemplateVersionDescription)
	}

	if launchTemplateVersion.LaunchTemplateVersionData != nil {
		if launchTemplateVersion.LaunchTemplateVersionData.InstanceType != nil {
			_ = d.Set("instance_type", launchTemplateVersion.LaunchTemplateVersionData.InstanceType)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.ImageId != nil {
			_ = d.Set("image_id", launchTemplateVersion.LaunchTemplateVersionData.ImageId)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.Placement != nil {
			placementMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.Placement.Zone != nil {
				placementMap["zone"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.Zone
			}

			if launchTemplateVersion.LaunchTemplateVersionData.Placement.ProjectId != nil {
				placementMap["project_id"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.ProjectId
			}

			if launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIds != nil {
				placementMap["host_ids"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIds
			}

			if launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIps != nil {
				placementMap["host_ips"] = launchTemplateVersion.LaunchTemplateVersionData.Placement.HostIps
			}

			_ = d.Set("placement", []interface{}{placementMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk != nil {
			systemDiskMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskType != nil {
				systemDiskMap["disk_type"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskType
			}

			if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskId != nil {
				systemDiskMap["disk_id"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskId
			}

			if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskSize != nil {
				systemDiskMap["disk_size"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.DiskSize
			}

			if launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.CdcId != nil {
				systemDiskMap["cdc_id"] = launchTemplateVersion.LaunchTemplateVersionData.SystemDisk.CdcId
			}

			_ = d.Set("system_disk", []interface{}{systemDiskMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.DataDisks != nil {
			dataDisksList := []interface{}{}
			for _, dataDisk := range launchTemplateVersion.LaunchTemplateVersionData.DataDisks {
				dataDisksMap := map[string]interface{}{}

				if dataDisk.DiskSize != nil {
					dataDisksMap["disk_size"] = dataDisk.DiskSize
				}

				if dataDisk.DiskType != nil {
					dataDisksMap["disk_type"] = dataDisk.DiskType
				}

				if dataDisk.DiskId != nil {
					dataDisksMap["disk_id"] = dataDisk.DiskId
				}

				if dataDisk.DeleteWithInstance != nil {
					dataDisksMap["delete_with_instance"] = dataDisk.DeleteWithInstance
				}

				if dataDisk.SnapshotId != nil {
					dataDisksMap["snapshot_id"] = dataDisk.SnapshotId
				}

				if dataDisk.Encrypt != nil {
					dataDisksMap["encrypt"] = dataDisk.Encrypt
				}

				if dataDisk.KmsKeyId != nil {
					dataDisksMap["kms_key_id"] = dataDisk.KmsKeyId
				}

				if dataDisk.ThroughputPerformance != nil {
					dataDisksMap["throughput_performance"] = dataDisk.ThroughputPerformance
				}

				if dataDisk.CdcId != nil {
					dataDisksMap["cdc_id"] = dataDisk.CdcId
				}

				dataDisksList = append(dataDisksList, dataDisksMap)
			}

			_ = d.Set("data_disks", dataDisksList)

		}

		if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud != nil {
			virtualPrivateCloudMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.VpcId != nil {
				virtualPrivateCloudMap["vpc_id"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.VpcId
			}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.SubnetId != nil {
				virtualPrivateCloudMap["subnet_id"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.SubnetId
			}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.AsVpcGateway != nil {
				virtualPrivateCloudMap["as_vpc_gateway"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.AsVpcGateway
			}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.PrivateIpAddresses != nil {
				virtualPrivateCloudMap["private_ip_addresses"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.PrivateIpAddresses
			}

			if launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.Ipv6AddressCount != nil {
				virtualPrivateCloudMap["ipv6_address_count"] = launchTemplateVersion.LaunchTemplateVersionData.VirtualPrivateCloud.Ipv6AddressCount
			}

			_ = d.Set("virtual_private_cloud", []interface{}{virtualPrivateCloudMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible != nil {
			internetAccessibleMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetChargeType != nil {
				internetAccessibleMap["internet_charge_type"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetChargeType
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetMaxBandwidthOut != nil {
				internetAccessibleMap["internet_max_bandwidth_out"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.InternetMaxBandwidthOut
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.PublicIpAssigned != nil {
				internetAccessibleMap["public_ip_assigned"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.PublicIpAssigned
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.BandwidthPackageId != nil {
				internetAccessibleMap["bandwidth_package_id"] = launchTemplateVersion.LaunchTemplateVersionData.InternetAccessible.BandwidthPackageId
			}

			_ = d.Set("internet_accessible", []interface{}{internetAccessibleMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceCount != nil {
			_ = d.Set("instance_count", launchTemplateVersion.LaunchTemplateVersionData.InstanceCount)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceName != nil {
			_ = d.Set("instance_name", launchTemplateVersion.LaunchTemplateVersionData.InstanceName)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings != nil {
			loginSettingsMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.Password != nil {
				loginSettingsMap["password"] = launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.Password
			}

			if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeyIds != nil {
				loginSettingsMap["key_ids"] = launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeyIds
			}

			if launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeepImageLogin != nil {
				loginSettingsMap["keep_image_login"] = launchTemplateVersion.LaunchTemplateVersionData.LoginSettings.KeepImageLogin
			}

			_ = d.Set("login_settings", []interface{}{loginSettingsMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.SecurityGroupIds != nil {
			_ = d.Set("security_group_ids", launchTemplateVersion.LaunchTemplateVersionData.SecurityGroupIds)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService != nil {
			enhancedServiceMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.SecurityService != nil {
				securityServiceMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.SecurityService.Enabled != nil {
					securityServiceMap["enabled"] = launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.SecurityService.Enabled
				}

				enhancedServiceMap["security_service"] = []interface{}{securityServiceMap}
			}

			if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.MonitorService != nil {
				monitorServiceMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.MonitorService.Enabled != nil {
					monitorServiceMap["enabled"] = launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.MonitorService.Enabled
				}

				enhancedServiceMap["monitor_service"] = []interface{}{monitorServiceMap}
			}

			if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.AutomationService != nil {
				automationServiceMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.AutomationService.Enabled != nil {
					automationServiceMap["enabled"] = launchTemplateVersion.LaunchTemplateVersionData.EnhancedService.AutomationService.Enabled
				}

				enhancedServiceMap["automation_service"] = []interface{}{automationServiceMap}
			}

			_ = d.Set("enhanced_service", []interface{}{enhancedServiceMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.ClientToken != nil {
			_ = d.Set("client_token", launchTemplateVersion.LaunchTemplateVersionData.ClientToken)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.HostName != nil {
			_ = d.Set("host_name", launchTemplateVersion.LaunchTemplateVersionData.HostName)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer != nil {
			actionTimerMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.TimerAction != nil {
				actionTimerMap["timer_action"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.TimerAction
			}

			if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.ActionTime != nil {
				actionTimerMap["action_time"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.ActionTime
			}

			if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals != nil {
				externalsMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.ReleaseAddress != nil {
					externalsMap["release_address"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.ReleaseAddress
				}

				if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.UnsupportNetworks != nil {
					externalsMap["unsupport_networks"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.UnsupportNetworks
				}

				if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr != nil {
					storageBlockAttrMap := map[string]interface{}{}

					if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.Type != nil {
						storageBlockAttrMap["type"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.Type
					}

					if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MinSize != nil {
						storageBlockAttrMap["min_size"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MinSize
					}

					if launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MaxSize != nil {
						storageBlockAttrMap["max_size"] = launchTemplateVersion.LaunchTemplateVersionData.ActionTimer.Externals.StorageBlockAttr.MaxSize
					}

					externalsMap["storage_block_attr"] = []interface{}{storageBlockAttrMap}
				}

				actionTimerMap["externals"] = []interface{}{externalsMap}
			}

			_ = d.Set("action_timer", []interface{}{actionTimerMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.DisasterRecoverGroupIds != nil {
			_ = d.Set("disaster_recover_group_ids", launchTemplateVersion.LaunchTemplateVersionData.DisasterRecoverGroupIds)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.TagSpecification != nil {
			tagSpecificationList := []interface{}{}
			for _, tagSpecification := range launchTemplateVersion.LaunchTemplateVersionData.TagSpecification {
				tagSpecificationMap := map[string]interface{}{}

				if tagSpecification.ResourceType != nil {
					tagSpecificationMap["resource_type"] = tagSpecification.ResourceType
				}

				if tagSpecification.Tags != nil {
					tagsList := []interface{}{}
					for _, tag := range tagSpecification.Tags {
						tagsMap := map[string]interface{}{}

						if tag.Key != nil {
							tagsMap["key"] = tag.Key
						}

						if tag.Value != nil {
							tagsMap["value"] = tag.Value
						}

						tagsList = append(tagsList, tagsMap)
					}

					tagSpecificationMap["tags"] = []interface{}{tagsList}
				}

				tagSpecificationList = append(tagSpecificationList, tagSpecificationMap)
			}

			_ = d.Set("tag_specification", tagSpecificationList)

		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions != nil {
			instanceMarketOptionsMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions != nil {
				spotOptionsMap := map[string]interface{}{}

				if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.MaxPrice != nil {
					spotOptionsMap["max_price"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.MaxPrice
				}

				if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.SpotInstanceType != nil {
					spotOptionsMap["spot_instance_type"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.SpotOptions.SpotInstanceType
				}

				instanceMarketOptionsMap["spot_options"] = []interface{}{spotOptionsMap}
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.MarketType != nil {
				instanceMarketOptionsMap["market_type"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceMarketOptions.MarketType
			}

			_ = d.Set("instance_market_options", []interface{}{instanceMarketOptionsMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.UserData != nil {
			_ = d.Set("user_data", launchTemplateVersion.LaunchTemplateVersionData.UserData)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.CamRoleName != nil {
			_ = d.Set("cam_role_name", launchTemplateVersion.LaunchTemplateVersionData.CamRoleName)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.HpcClusterId != nil {
			_ = d.Set("hpc_cluster_id", launchTemplateVersion.LaunchTemplateVersionData.HpcClusterId)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargeType != nil {
			_ = d.Set("instance_charge_type", launchTemplateVersion.LaunchTemplateVersionData.InstanceChargeType)
		}

		if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid != nil {
			instanceChargePrepaidMap := map[string]interface{}{}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.Period != nil {
				instanceChargePrepaidMap["period"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.Period
			}

			if launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.RenewFlag != nil {
				instanceChargePrepaidMap["renew_flag"] = launchTemplateVersion.LaunchTemplateVersionData.InstanceChargePrepaid.RenewFlag
			}

			_ = d.Set("instance_charge_prepaid", []interface{}{instanceChargePrepaidMap})
		}

		if launchTemplateVersion.LaunchTemplateVersionData.DisableApiTermination != nil {
			_ = d.Set("disable_api_termination", launchTemplateVersion.LaunchTemplateVersionData.DisableApiTermination)
		}
	}

	return nil
}

func resourceTencentCloudCvmLaunchTemplateVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_version.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	launchTemplateId := idSplit[0]
	launchTemplateVersionNumber := idSplit[1]

	if err := service.DeleteCvmLaunchTemplateVersionById(ctx, launchTemplateId, launchTemplateVersionNumber); err != nil {
		return err
	}

	return nil
}
