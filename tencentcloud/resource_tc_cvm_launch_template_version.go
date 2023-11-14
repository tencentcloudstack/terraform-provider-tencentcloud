/*
Provides a resource to create a cvm launch_template_version

Example Usage

```hcl
resource "tencentcloud_cvm_launch_template_version" "launch_template_version" {
  placement {
		zone = "ap-guangzhou-6"
		project_id = 0
		host_ids =
		host_ips =
		host_id = "host-dab6ejhx"

  }
  launch_template_id = "lt-lobxe2yo"
  launch_template_version = 1
  launch_template_version_description = ""
  instance_type = "S5.MEDIUM4"
  image_id = "img-eb30mz89"
  system_disk {
		disk_type = "CLOUD_PREMIUM"
		disk_id = ""
		disk_size = 50
		cdc_id = "cdc-b9pbd3px"

  }
  data_disks {
		disk_size = 50
		disk_type = "CLOUD_PREMIUM"
		disk_id = ""
		delete_with_instance = false
		snapshot_id = "snap-r9unnd89"
		encrypt = false
		kms_key_id = "kms-abcd1234"
		throughput_performance = 2
		cdc_id = "cdc-b9pbd3px"

  }
  virtual_private_cloud {
		vpc_id = "vpc-x2e4dam7"
		subnet_id = "subnet-g73bdf1r"
		as_vpc_gateway = false
		private_ip_addresses =
		ipv6_address_count = 1

  }
  internet_accessible {
		internet_charge_type = "TRAFFIC_POSTPAID_BY_HOUR"
		internet_max_bandwidth_out = 0
		public_ip_assigned = false
		bandwidth_package_id = ""

  }
  instance_count = 1
  instance_name = ""
  login_settings {
		password = "1@34qwer"
		key_ids =
		keep_image_login = "FALSE"

  }
  security_group_ids =
  enhanced_service {
		security_service {
			enabled = true
		}
		monitor_service {
			enabled = true
		}
		automation_service {
			enabled = false
		}

  }
  client_token = "cg_1678809600005774197"
  host_name = ""
  action_timer {
		timer_action = "TerminateInstances"
		action_time = "2018-05-29T11:26:40Z"
		externals {
			release_address = false
			unsupport_networks =
			storage_block_attr {
				type = "LOCAL_PRO"
				min_size = 10
				max_size = 100
			}
		}

  }
  disaster_recover_group_ids =
  tag_specification {
		resource_type = "instance"
		tags {
			key = "tagKey"
			value = "tagValue"
		}

  }
  instance_market_options {
		spot_options {
			max_price = "1.99"
			spot_instance_type = "one-time"
		}
		market_type = "spot"

  }
  user_data = "IyEvdXNyL2Jpbi9lbnYgcHl0a"
  dry_run = false
  cam_role_name = ""
  hpc_cluster_id = "hpc-bwu6b3e2"
  instance_charge_type = "POSTPAID_BY_HOUR"
  instance_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_AUTO_RENEW"

  }
  disable_api_termination = false
}
```

Import

cvm launch_template_version can be imported using the id, e.g.

```
terraform import tencentcloud_cvm_launch_template_version.launch_template_version launch_template_version_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCvmLaunchTemplateVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCvmLaunchTemplateVersionCreate,
		Read:   resourceTencentCloudCvmLaunchTemplateVersionRead,
		Update: resourceTencentCloudCvmLaunchTemplateVersionUpdate,
		Delete: resourceTencentCloudCvmLaunchTemplateVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"placement": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Location of the instance. You can use this parameter to specify the attributes of the instance, such as its availability zone, project, and CDH (for dedicated CVMs).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the availability zone where the instance resides. You can call the DescribeZones API and obtain the ID in the returned Zone field.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "ID of the project to which the instance belongs. This parameter can be obtained from the projectId returned by DescribeProject. If this is left empty, the default project is used.",
						},
						"host_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "ID list of CDHs from which the instance can be created. If you have purchased CDHs and specify this parameter, the instances you purchase will be randomly deployed on the CDHs.",
						},
						"host_ips": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "IPs of the hosts to create CVMs.",
						},
						"host_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the CDH to which the instance belongs, only used as an output parameter.",
						},
					},
				},
			},

			"launch_template_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance launch template ID. This parameter is used as a basis for creating new template versions.",
			},

			"launch_template_version": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "This parameter, when specified, is used to create instance launch templates. If this parameter is not specified, the default version will be used.",
			},

			"launch_template_version_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description of instance launch template versions. This parameter can contain 2-256 characters.",
			},

			"instance_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The type of the instance. If this parameter is not specified, the system will dynamically specify the default model according to the resource sales in the current region.",
			},

			"image_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Image ID.",
			},

			"system_disk": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "System disk configuration information of the instance. If this parameter is not specified, it is assigned according to the system default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of system disk.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "System disk ID. System disks whose type is LOCAL_BASIC or LOCAL_SSD do not have an ID and do not support this parameter. It is only used as a response parameter for APIs such as DescribeInstances, and cannot be used as a request parameter for APIs such as RunInstances.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "System disk size; unit: GB; default value: 50 GB.",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the dedicated cluster to which the instance belongs.",
						},
					},
				},
			},

			"data_disks": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "The configuration information of instance data disks. If this parameter is not specified, no data disk will be purchased by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Data disk size (in GB). The minimum adjustment increment is 10 GB. The value range varies by data disk type.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The type of data disk.",
						},
						"disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "System disk ID. System disks whose type is LOCAL_BASIC or LOCAL_SSD do not have an ID and do not support this parameter. It is only used as a response parameter for APIs such as DescribeInstances, and cannot be used as a request parameter for APIs such as RunInstances.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to terminate the data disk when its CVM is terminated.",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data disk snapshot ID. The size of the selected data disk snapshot must be smaller than that of the data disk. Note: This field may return null, indicating that no valid value is found.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Specifies whether the data disk is encrypted.",
						},
						"kms_key_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the custom CMK in the format of UUID or “kms-abcd1234”. .",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Cloud disk performance, in MB/s .",
						},
						"cdc_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the dedicated cluster to which the instance belongs.",
						},
					},
				},
			},

			"virtual_private_cloud": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes information on VPC, including subnets, IP addresses, etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC ID in the format of vpc-xxx, if you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "VPC subnet ID in the format subnet-xxx, if you specify DEFAULT for both VpcId and SubnetId when creating an instance, the default VPC will be used.",
						},
						"as_vpc_gateway": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to use a CVM instance as a public gateway. The public gateway is only available when the instance has a public IP and resides in a VPC.",
						},
						"private_ip_addresses": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Array of VPC subnet IPs. You can use this parameter when creating instances or modifying VPC attributes of instances. Currently you can specify multiple IPs in one subnet only when creating multiple instances at the same time.",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of IPv6 addresses randomly generated for the ENI.",
						},
					},
				},
			},

			"internet_accessible": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes the accessibility of an instance in the public network, including its network billing method, maximum bandwidth, etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"internet_charge_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Network connection billing plan.",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The maximum outbound bandwidth of the public network, in Mbps. The default value is 0 Mbps.",
						},
						"public_ip_assigned": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to assign a public IP.",
						},
						"bandwidth_package_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bandwidth package ID.",
						},
					},
				},
			},

			"instance_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The number of instances to be purchased.",
			},

			"instance_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance name to be displayed.",
			},

			"login_settings": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes login settings of an instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Login password of the instance.",
						},
						"key_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "List of key IDs. After an instance is associated with a key, you can access the instance with the private key in the key pair.",
						},
						"keep_image_login": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to keep the original settings of an image.",
						},
					},
				},
			},

			"security_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security groups to which the instance belongs. If this parameter is not specified, the instance will be associated with default security groups.",
			},

			"enhanced_service": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Enhanced service. You can use this parameter to specify whether to enable services such as Anti-DDoS and Cloud Monitor. If this parameter is not specified, Cloud Monitor and Anti-DDoS are enabled for public images by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Enables cloud security service. If this parameter is not specified, the cloud security service will be enabled by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to enable Cloud Security.",
									},
								},
							},
						},
						"monitor_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Enables cloud monitor service. If this parameter is not specified, the cloud monitor service will be enabled by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether to enable Cloud Monitor.",
									},
								},
							},
						},
						"automation_service": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Whether to enable the TAT service. If this parameter is not specified, the TAT service is enabled for public images and disabled for other images by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
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
				Type:        schema.TypeString,
				Description: "A unique string supplied by the client to ensure that the request is idempotent. Its maximum length is 64 ASCII characters. If this parameter is not specified, the idem-potency of the request cannot be guaranteed.",
			},

			"host_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Hostname of a CVM.",
			},

			"action_timer": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Scheduled tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timer_action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Timer name. Currently TerminateInstances is the only supported value.",
						},
						"action_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Execution time, displayed according to ISO8601 standard, and UTC time is used. The format is YYYY-MM-DDThh:mm:ssZ. For example, 2018-05-29T11:26:40Z, the execution must be at least 5 minutes later than the current time.",
						},
						"externals": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Additional data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"release_address": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Release address.",
									},
									"unsupport_networks": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Not supported network.",
									},
									"storage_block_attr": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Information on local HDD storage.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Local HDD storage type. Value: LOCAL_PRO.",
												},
												"min_size": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Minimum capacity of local HDD storage.",
												},
												"max_size": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Maximum capacity of local HDD storage .",
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
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Placement group ID. You can only specify one.",
			},

			"tag_specification": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Description of tags associated with resource instances during instance creation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
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
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
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
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Options related to bidding requests.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spot_options": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Options related to bidding.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_price": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Bidding price.",
									},
									"spot_instance_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Bidding request type. Currently only one-time is supported.",
									},
								},
							},
						},
						"market_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Market option type. Currently spot is the only supported value.",
						},
					},
				},
			},

			"user_data": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User data provided to the instance. This parameter needs to be encoded in base64 format with the maximum size of 16 KB.",
			},

			"dry_run": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether the request is a dry run only. .",
			},

			"cam_role_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The role name of CAM.",
			},

			"hpc_cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "HPC cluster ID. The HPC cluster must and can only be specified for a high-performance computing instance.",
			},

			"instance_charge_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The charge type of instance.",
			},

			"instance_charge_prepaid": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Describes the billing method of an instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Subscription period; unit: month; valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60.",
						},
						"renew_flag": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Auto renewal flag. Valid values: &amp;lt;br&amp;gt;&amp;lt;li&amp;gt;NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically &amp;lt;br&amp;gt;&amp;lt;li&amp;gt;NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically &amp;lt;br&amp;gt;&amp;lt;li&amp;gt;DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically &amp;lt;br&amp;gt;&amp;lt;br&amp;gt;Default value: NOTIFY_AND_MANUAL_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient.",
						},
					},
				},
			},

			"disable_api_termination": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether the termination protection is enabled. Values: &amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;`TRUE`: Enable instance protection, which means that this instance can not be deleted by an API action.&amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;li&amp;amp;gt;`FALSE`: Do not enable the instance protection.&amp;amp;lt;br&amp;amp;gt;&amp;amp;lt;br&amp;amp;gt;Default value: `FALSE`.",
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
		if v, ok := dMap["host_id"]; ok {
			placement.HostId = helper.String(v.(string))
		}
		request.Placement = &placement
	}

	if v, ok := d.GetOk("launch_template_id"); ok {
		launchTemplateId = v.(string)
		request.LaunchTemplateId = helper.String(v.(string))
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

	launchTemplateId = *response.Response.LaunchTemplateId
	d.SetId(launchTemplateId)

	return resourceTencentCloudCvmLaunchTemplateVersionRead(d, meta)
}

func resourceTencentCloudCvmLaunchTemplateVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_version.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}

	launchTemplateVersionId := d.Id()

	launchTemplateVersion, err := service.DescribeCvmLaunchTemplateVersionById(ctx, launchTemplateId)
	if err != nil {
		return err
	}

	if launchTemplateVersion == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CvmLaunchTemplateVersion` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if launchTemplateVersion.Placement != nil {
		placementMap := map[string]interface{}{}

		if launchTemplateVersion.Placement.Zone != nil {
			placementMap["zone"] = launchTemplateVersion.Placement.Zone
		}

		if launchTemplateVersion.Placement.ProjectId != nil {
			placementMap["project_id"] = launchTemplateVersion.Placement.ProjectId
		}

		if launchTemplateVersion.Placement.HostIds != nil {
			placementMap["host_ids"] = launchTemplateVersion.Placement.HostIds
		}

		if launchTemplateVersion.Placement.HostIps != nil {
			placementMap["host_ips"] = launchTemplateVersion.Placement.HostIps
		}

		if launchTemplateVersion.Placement.HostId != nil {
			placementMap["host_id"] = launchTemplateVersion.Placement.HostId
		}

		_ = d.Set("placement", []interface{}{placementMap})
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

	if launchTemplateVersion.InstanceType != nil {
		_ = d.Set("instance_type", launchTemplateVersion.InstanceType)
	}

	if launchTemplateVersion.ImageId != nil {
		_ = d.Set("image_id", launchTemplateVersion.ImageId)
	}

	if launchTemplateVersion.SystemDisk != nil {
		systemDiskMap := map[string]interface{}{}

		if launchTemplateVersion.SystemDisk.DiskType != nil {
			systemDiskMap["disk_type"] = launchTemplateVersion.SystemDisk.DiskType
		}

		if launchTemplateVersion.SystemDisk.DiskId != nil {
			systemDiskMap["disk_id"] = launchTemplateVersion.SystemDisk.DiskId
		}

		if launchTemplateVersion.SystemDisk.DiskSize != nil {
			systemDiskMap["disk_size"] = launchTemplateVersion.SystemDisk.DiskSize
		}

		if launchTemplateVersion.SystemDisk.CdcId != nil {
			systemDiskMap["cdc_id"] = launchTemplateVersion.SystemDisk.CdcId
		}

		_ = d.Set("system_disk", []interface{}{systemDiskMap})
	}

	if launchTemplateVersion.DataDisks != nil {
		dataDisksList := []interface{}{}
		for _, dataDisks := range launchTemplateVersion.DataDisks {
			dataDisksMap := map[string]interface{}{}

			if launchTemplateVersion.DataDisks.DiskSize != nil {
				dataDisksMap["disk_size"] = launchTemplateVersion.DataDisks.DiskSize
			}

			if launchTemplateVersion.DataDisks.DiskType != nil {
				dataDisksMap["disk_type"] = launchTemplateVersion.DataDisks.DiskType
			}

			if launchTemplateVersion.DataDisks.DiskId != nil {
				dataDisksMap["disk_id"] = launchTemplateVersion.DataDisks.DiskId
			}

			if launchTemplateVersion.DataDisks.DeleteWithInstance != nil {
				dataDisksMap["delete_with_instance"] = launchTemplateVersion.DataDisks.DeleteWithInstance
			}

			if launchTemplateVersion.DataDisks.SnapshotId != nil {
				dataDisksMap["snapshot_id"] = launchTemplateVersion.DataDisks.SnapshotId
			}

			if launchTemplateVersion.DataDisks.Encrypt != nil {
				dataDisksMap["encrypt"] = launchTemplateVersion.DataDisks.Encrypt
			}

			if launchTemplateVersion.DataDisks.KmsKeyId != nil {
				dataDisksMap["kms_key_id"] = launchTemplateVersion.DataDisks.KmsKeyId
			}

			if launchTemplateVersion.DataDisks.ThroughputPerformance != nil {
				dataDisksMap["throughput_performance"] = launchTemplateVersion.DataDisks.ThroughputPerformance
			}

			if launchTemplateVersion.DataDisks.CdcId != nil {
				dataDisksMap["cdc_id"] = launchTemplateVersion.DataDisks.CdcId
			}

			dataDisksList = append(dataDisksList, dataDisksMap)
		}

		_ = d.Set("data_disks", dataDisksList)

	}

	if launchTemplateVersion.VirtualPrivateCloud != nil {
		virtualPrivateCloudMap := map[string]interface{}{}

		if launchTemplateVersion.VirtualPrivateCloud.VpcId != nil {
			virtualPrivateCloudMap["vpc_id"] = launchTemplateVersion.VirtualPrivateCloud.VpcId
		}

		if launchTemplateVersion.VirtualPrivateCloud.SubnetId != nil {
			virtualPrivateCloudMap["subnet_id"] = launchTemplateVersion.VirtualPrivateCloud.SubnetId
		}

		if launchTemplateVersion.VirtualPrivateCloud.AsVpcGateway != nil {
			virtualPrivateCloudMap["as_vpc_gateway"] = launchTemplateVersion.VirtualPrivateCloud.AsVpcGateway
		}

		if launchTemplateVersion.VirtualPrivateCloud.PrivateIpAddresses != nil {
			virtualPrivateCloudMap["private_ip_addresses"] = launchTemplateVersion.VirtualPrivateCloud.PrivateIpAddresses
		}

		if launchTemplateVersion.VirtualPrivateCloud.Ipv6AddressCount != nil {
			virtualPrivateCloudMap["ipv6_address_count"] = launchTemplateVersion.VirtualPrivateCloud.Ipv6AddressCount
		}

		_ = d.Set("virtual_private_cloud", []interface{}{virtualPrivateCloudMap})
	}

	if launchTemplateVersion.InternetAccessible != nil {
		internetAccessibleMap := map[string]interface{}{}

		if launchTemplateVersion.InternetAccessible.InternetChargeType != nil {
			internetAccessibleMap["internet_charge_type"] = launchTemplateVersion.InternetAccessible.InternetChargeType
		}

		if launchTemplateVersion.InternetAccessible.InternetMaxBandwidthOut != nil {
			internetAccessibleMap["internet_max_bandwidth_out"] = launchTemplateVersion.InternetAccessible.InternetMaxBandwidthOut
		}

		if launchTemplateVersion.InternetAccessible.PublicIpAssigned != nil {
			internetAccessibleMap["public_ip_assigned"] = launchTemplateVersion.InternetAccessible.PublicIpAssigned
		}

		if launchTemplateVersion.InternetAccessible.BandwidthPackageId != nil {
			internetAccessibleMap["bandwidth_package_id"] = launchTemplateVersion.InternetAccessible.BandwidthPackageId
		}

		_ = d.Set("internet_accessible", []interface{}{internetAccessibleMap})
	}

	if launchTemplateVersion.InstanceCount != nil {
		_ = d.Set("instance_count", launchTemplateVersion.InstanceCount)
	}

	if launchTemplateVersion.InstanceName != nil {
		_ = d.Set("instance_name", launchTemplateVersion.InstanceName)
	}

	if launchTemplateVersion.LoginSettings != nil {
		loginSettingsMap := map[string]interface{}{}

		if launchTemplateVersion.LoginSettings.Password != nil {
			loginSettingsMap["password"] = launchTemplateVersion.LoginSettings.Password
		}

		if launchTemplateVersion.LoginSettings.KeyIds != nil {
			loginSettingsMap["key_ids"] = launchTemplateVersion.LoginSettings.KeyIds
		}

		if launchTemplateVersion.LoginSettings.KeepImageLogin != nil {
			loginSettingsMap["keep_image_login"] = launchTemplateVersion.LoginSettings.KeepImageLogin
		}

		_ = d.Set("login_settings", []interface{}{loginSettingsMap})
	}

	if launchTemplateVersion.SecurityGroupIds != nil {
		_ = d.Set("security_group_ids", launchTemplateVersion.SecurityGroupIds)
	}

	if launchTemplateVersion.EnhancedService != nil {
		enhancedServiceMap := map[string]interface{}{}

		if launchTemplateVersion.EnhancedService.SecurityService != nil {
			securityServiceMap := map[string]interface{}{}

			if launchTemplateVersion.EnhancedService.SecurityService.Enabled != nil {
				securityServiceMap["enabled"] = launchTemplateVersion.EnhancedService.SecurityService.Enabled
			}

			enhancedServiceMap["security_service"] = []interface{}{securityServiceMap}
		}

		if launchTemplateVersion.EnhancedService.MonitorService != nil {
			monitorServiceMap := map[string]interface{}{}

			if launchTemplateVersion.EnhancedService.MonitorService.Enabled != nil {
				monitorServiceMap["enabled"] = launchTemplateVersion.EnhancedService.MonitorService.Enabled
			}

			enhancedServiceMap["monitor_service"] = []interface{}{monitorServiceMap}
		}

		if launchTemplateVersion.EnhancedService.AutomationService != nil {
			automationServiceMap := map[string]interface{}{}

			if launchTemplateVersion.EnhancedService.AutomationService.Enabled != nil {
				automationServiceMap["enabled"] = launchTemplateVersion.EnhancedService.AutomationService.Enabled
			}

			enhancedServiceMap["automation_service"] = []interface{}{automationServiceMap}
		}

		_ = d.Set("enhanced_service", []interface{}{enhancedServiceMap})
	}

	if launchTemplateVersion.ClientToken != nil {
		_ = d.Set("client_token", launchTemplateVersion.ClientToken)
	}

	if launchTemplateVersion.HostName != nil {
		_ = d.Set("host_name", launchTemplateVersion.HostName)
	}

	if launchTemplateVersion.ActionTimer != nil {
		actionTimerMap := map[string]interface{}{}

		if launchTemplateVersion.ActionTimer.TimerAction != nil {
			actionTimerMap["timer_action"] = launchTemplateVersion.ActionTimer.TimerAction
		}

		if launchTemplateVersion.ActionTimer.ActionTime != nil {
			actionTimerMap["action_time"] = launchTemplateVersion.ActionTimer.ActionTime
		}

		if launchTemplateVersion.ActionTimer.Externals != nil {
			externalsMap := map[string]interface{}{}

			if launchTemplateVersion.ActionTimer.Externals.ReleaseAddress != nil {
				externalsMap["release_address"] = launchTemplateVersion.ActionTimer.Externals.ReleaseAddress
			}

			if launchTemplateVersion.ActionTimer.Externals.UnsupportNetworks != nil {
				externalsMap["unsupport_networks"] = launchTemplateVersion.ActionTimer.Externals.UnsupportNetworks
			}

			if launchTemplateVersion.ActionTimer.Externals.StorageBlockAttr != nil {
				storageBlockAttrMap := map[string]interface{}{}

				if launchTemplateVersion.ActionTimer.Externals.StorageBlockAttr.Type != nil {
					storageBlockAttrMap["type"] = launchTemplateVersion.ActionTimer.Externals.StorageBlockAttr.Type
				}

				if launchTemplateVersion.ActionTimer.Externals.StorageBlockAttr.MinSize != nil {
					storageBlockAttrMap["min_size"] = launchTemplateVersion.ActionTimer.Externals.StorageBlockAttr.MinSize
				}

				if launchTemplateVersion.ActionTimer.Externals.StorageBlockAttr.MaxSize != nil {
					storageBlockAttrMap["max_size"] = launchTemplateVersion.ActionTimer.Externals.StorageBlockAttr.MaxSize
				}

				externalsMap["storage_block_attr"] = []interface{}{storageBlockAttrMap}
			}

			actionTimerMap["externals"] = []interface{}{externalsMap}
		}

		_ = d.Set("action_timer", []interface{}{actionTimerMap})
	}

	if launchTemplateVersion.DisasterRecoverGroupIds != nil {
		_ = d.Set("disaster_recover_group_ids", launchTemplateVersion.DisasterRecoverGroupIds)
	}

	if launchTemplateVersion.TagSpecification != nil {
		tagSpecificationList := []interface{}{}
		for _, tagSpecification := range launchTemplateVersion.TagSpecification {
			tagSpecificationMap := map[string]interface{}{}

			if launchTemplateVersion.TagSpecification.ResourceType != nil {
				tagSpecificationMap["resource_type"] = launchTemplateVersion.TagSpecification.ResourceType
			}

			if launchTemplateVersion.TagSpecification.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range launchTemplateVersion.TagSpecification.Tags {
					tagsMap := map[string]interface{}{}

					if tags.Key != nil {
						tagsMap["key"] = tags.Key
					}

					if tags.Value != nil {
						tagsMap["value"] = tags.Value
					}

					tagsList = append(tagsList, tagsMap)
				}

				tagSpecificationMap["tags"] = []interface{}{tagsList}
			}

			tagSpecificationList = append(tagSpecificationList, tagSpecificationMap)
		}

		_ = d.Set("tag_specification", tagSpecificationList)

	}

	if launchTemplateVersion.InstanceMarketOptions != nil {
		instanceMarketOptionsMap := map[string]interface{}{}

		if launchTemplateVersion.InstanceMarketOptions.SpotOptions != nil {
			spotOptionsMap := map[string]interface{}{}

			if launchTemplateVersion.InstanceMarketOptions.SpotOptions.MaxPrice != nil {
				spotOptionsMap["max_price"] = launchTemplateVersion.InstanceMarketOptions.SpotOptions.MaxPrice
			}

			if launchTemplateVersion.InstanceMarketOptions.SpotOptions.SpotInstanceType != nil {
				spotOptionsMap["spot_instance_type"] = launchTemplateVersion.InstanceMarketOptions.SpotOptions.SpotInstanceType
			}

			instanceMarketOptionsMap["spot_options"] = []interface{}{spotOptionsMap}
		}

		if launchTemplateVersion.InstanceMarketOptions.MarketType != nil {
			instanceMarketOptionsMap["market_type"] = launchTemplateVersion.InstanceMarketOptions.MarketType
		}

		_ = d.Set("instance_market_options", []interface{}{instanceMarketOptionsMap})
	}

	if launchTemplateVersion.UserData != nil {
		_ = d.Set("user_data", launchTemplateVersion.UserData)
	}

	if launchTemplateVersion.DryRun != nil {
		_ = d.Set("dry_run", launchTemplateVersion.DryRun)
	}

	if launchTemplateVersion.CamRoleName != nil {
		_ = d.Set("cam_role_name", launchTemplateVersion.CamRoleName)
	}

	if launchTemplateVersion.HpcClusterId != nil {
		_ = d.Set("hpc_cluster_id", launchTemplateVersion.HpcClusterId)
	}

	if launchTemplateVersion.InstanceChargeType != nil {
		_ = d.Set("instance_charge_type", launchTemplateVersion.InstanceChargeType)
	}

	if launchTemplateVersion.InstanceChargePrepaid != nil {
		instanceChargePrepaidMap := map[string]interface{}{}

		if launchTemplateVersion.InstanceChargePrepaid.Period != nil {
			instanceChargePrepaidMap["period"] = launchTemplateVersion.InstanceChargePrepaid.Period
		}

		if launchTemplateVersion.InstanceChargePrepaid.RenewFlag != nil {
			instanceChargePrepaidMap["renew_flag"] = launchTemplateVersion.InstanceChargePrepaid.RenewFlag
		}

		_ = d.Set("instance_charge_prepaid", []interface{}{instanceChargePrepaidMap})
	}

	if launchTemplateVersion.DisableApiTermination != nil {
		_ = d.Set("disable_api_termination", launchTemplateVersion.DisableApiTermination)
	}

	return nil
}

func resourceTencentCloudCvmLaunchTemplateVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_version.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := cvm.NewModifyLaunchTemplateDefaultVersionRequest()

	launchTemplateVersionId := d.Id()

	request.LaunchTemplateId = &launchTemplateId

	immutableArgs := []string{"placement", "launch_template_id", "launch_template_version", "launch_template_version_description", "instance_type", "image_id", "system_disk", "data_disks", "virtual_private_cloud", "internet_accessible", "instance_count", "instance_name", "login_settings", "security_group_ids", "enhanced_service", "client_token", "host_name", "action_timer", "disaster_recover_group_ids", "tag_specification", "instance_market_options", "user_data", "dry_run", "cam_role_name", "hpc_cluster_id", "instance_charge_type", "instance_charge_prepaid", "disable_api_termination"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("launch_template_id") {
		if v, ok := d.GetOk("launch_template_id"); ok {
			request.LaunchTemplateId = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ModifyLaunchTemplateDefaultVersion(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cvm launchTemplateVersion failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCvmLaunchTemplateVersionRead(d, meta)
}

func resourceTencentCloudCvmLaunchTemplateVersionDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cvm_launch_template_version.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CvmService{client: meta.(*TencentCloudClient).apiV3Conn}
	launchTemplateVersionId := d.Id()

	if err := service.DeleteCvmLaunchTemplateVersionById(ctx, launchTemplateId); err != nil {
		return err
	}

	return nil
}
