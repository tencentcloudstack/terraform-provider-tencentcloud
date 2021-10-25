/*
Provides a CVM instance resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** At present, 'PREPAID' instance cannot be deleted and must wait it to be outdated and released automatically.

Example Usage

```hcl
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "my_favorite_zones" {
}

// Create VPC resource
resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = tencentcloud_vpc.app.id
  availability_zone = data.tencentcloud_availability_zones.my_favorite_zones.zones.0.name
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

// Create 2 CVM instances to host awesome_app
resource "tencentcloud_instance" "my_awesome_app" {
  instance_name              = "awesome_app"
  availability_zone          = data.tencentcloud_availability_zones.my_favorite_zones.zones.0.name
  image_id                   = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.app.id
  subnet_id                  = tencentcloud_subnet.app.id
  internet_max_bandwidth_out = 20
  count                      = 2
  cam_role_name              = "CVM_QcsRole"

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
	encrypt = false
  }

  tags = {
    tagKey = "tagValue"
  }
}
```

Create CVM instance based on CDH
```hcl
variable "availability_zone" {
  default = "ap-shanghai-4"
}

resource "tencentcloud_cdh_instance" "foo" {
  availability_zone = var.availability_zone
  host_type = "HM50"
  charge_type = "PREPAID"
  instance_charge_type_prepaid_period = 1
  hostname = "test"
  prepaid_renew_flag = "DISABLE_NOTIFY_AND_MANUAL_RENEW"
}

data "tencentcloud_cdh_instances" "list" {
  availability_zone = var.availability_zone
  host_id = tencentcloud_cdh_instance.foo.id
  hostname = "test"
  host_state = "RUNNING"
}

resource "tencentcloud_key_pair" "random_key" {
  key_name   = "tf_example_key6"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}

resource "tencentcloud_placement_group" "foo" {
  name = "test"
  type = "HOST"
}

resource "tencentcloud_instance" "foo" {
  availability_zone = var.availability_zone
  instance_name     = "terraform-testing"
  image_id          = "img-ix05e4px"
  key_name          = tencentcloud_key_pair.random_key.id
  placement_group_id = tencentcloud_placement_group.foo.id
  security_groups               = ["sg-9c3f33xk"]
  system_disk_type  = "CLOUD_PREMIUM"

  instance_charge_type = "CDHPAID"
  cdh_instance_type     = "CDH_10C10G"
  cdh_host_id = tencentcloud_cdh_instance.foo.id

  vpc_id                     = "vpc-31zmeluu"
  subnet_id                  = "subnet-aujc02np"
  allocate_public_ip    = true
  internet_max_bandwidth_out = 2
  count                      = 3

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt = false
  }
}
```

Import

CVM instance can be imported using the id, e.g.

```
terraform import tencentcloud_instance.foo ins-2qol3a80
```
*/
package tencentcloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudInstanceCreate,
		Read:   resourceTencentCloudInstanceRead,
		Update: resourceTencentCloudInstanceUpdate,
		Delete: resourceTencentCloudInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image to use for the instance. Changing `image_id` will cause the instance reset.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone for the CVM instance.",
			},
			"instance_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				Deprecated:   "It has been deprecated from version 1.59.18. Use built-in `count` instead.",
				ValidateFunc: validateIntegerInRange(1, 100),
				Description:  "The number of instances to be purchased. Value range:[1,100]; default value: 1.",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Terraform-CVM-Instance",
				ValidateFunc: validateStringLengthInRange(2, 128),
				Description:  "The name of the instance. The max length of instance_name is 60, and default value is `Terraform-CVM-Instance`.",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateInstanceType,
				Description:  "The type of the instance.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The hostname of the instance. Windows instance: The name should be a combination of 2 to 15 characters comprised of letters (case insensitive), numbers, and hyphens (-). Period (.) is not supported, and the name cannot be a string of pure numbers. Other types (such as Linux) of instances: The name should be a combination of 2 to 60 characters, supporting multiple periods (.). The piece between two periods is composed of letters (case insensitive), numbers, and hyphens (-). Modifying will cause the instance reset.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The project the instance belongs to, default to 0.",
			},
			"running_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Set instance to running or stop. Default value is true, the instance will shutdown when this flag is false.",
			},
			"stopped_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Billing method of a pay-as-you-go instance after shutdown. Available values: `KEEP_CHARGING`,`STOP_CHARGING`. Default `KEEP_CHARGING`.",
				ValidateFunc: validateAllowedStringValue([]string{
					CVM_STOP_MODE_KEEP_CHARGING,
					CVM_STOP_MODE_STOP_CHARGING,
				}),
			},
			"placement_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of a placement group.",
			},
			// payment
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      CVM_CHARGE_TYPE_POSTPAID,
				ValidateFunc: validateAllowedStringValue(CVM_CHARGE_TYPE),
				Description:  "The charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID` and `CDHPAID`. The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR` and `CDHPAID`. `PREPAID` instance may not allow to delete before expired. `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time. `CDHPAID` instance must set `cdh_instance_type` and `cdh_host_id`.",
			},
			"instance_charge_type_prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue(CVM_PREPAID_PERIOD),
				Description:  "The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`. Modifying will cause the instance reset.",
			},
			"instance_charge_type_prepaid_renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue(CVM_PREPAID_RENEW_FLAG),
				Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`. Modifying will cause the instance reset.",
			},
			"spot_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CVM_SPOT_INSTANCE_TYPE),
				Description:  "Type of spot instance, only support `ONE-TIME` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.",
			},
			"spot_max_price": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateStringNumber,
				Description:  "Max price of a spot instance, is the format of decimal string, for example \"0.50\". Note: it only works when instance_charge_type is set to `SPOTPAID`.",
			},
			"cdh_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringPrefix("CDH_"),
				Description:  "Type of instance created on cdh, the value of this parameter is in the format of CDH_XCXG based on the number of CPU cores and memory capacity. Note: it only works when instance_charge_type is set to `CDHPAID`.",
			},
			"cdh_host_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Id of cdh instance. Note: it only works when instance_charge_type is set to `CDHPAID`.",
			},
			// network
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CVM_INTERNET_CHARGE_TYPE),
				Description:  "Internet charge type of the instance, Valid values are `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`. This value does not need to be set when `allocate_public_ip` is false.",
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "bandwidth package id. if user is standard user, then the bandwidth_package_id is needed, or default has bandwidth_package_id.",
			},
			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bits per second). This value does not need to be set when `allocate_public_ip` is false.",
			},
			"allocate_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Associate a public IP address with an instance in a VPC or Classic. Boolean value, Default is false.",
			},
			// vpc
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of a VPC network. If you want to create instances in a VPC network, this parameter must be set.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of a VPC subnet. If you want to create instances in a VPC network, this parameter must be set.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The private IP to be assigned to this instance, must be in the provided subnet and available.",
			},
			// security group
			"security_groups": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "A list of security group IDs to associate with.",
			},
			// storage
			"system_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CVM_DISK_TYPE_CLOUD_BASIC,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CVM_DISK_TYPE),
				Description:  "System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_BASIC`: HDD cloud disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.",
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(50, 1000),
				Description:  "Size of the system disk. Valid value ranges: (50~1000). and unit is GB. Default is 50GB.",
			},
			"system_disk_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "System disk snapshot ID used to initialize the system disk. When system disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.",
			},
			"data_disks": {
				Type:        schema.TypeList,
				MaxItems:    10,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Settings for data disks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_disk_type": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Data disk type. For more information about limits on different data disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_BASIC`: HDD cloud disk, `CLOUD_PREMIUM`: Premium Cloud Storage, `CLOUD_SSD`: SSD, `CLOUD_HSSD`: Enhanced SSD. NOTE: `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.",
						},
						"data_disk_size": {
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
							Description: "Size of the data disk, and unit is GB. If disk type is `CLOUD_SSD`, the size range is [100, 16000], and the others are [10-16000].",
						},
						"data_disk_snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Snapshot ID of the data disk. The selected data disk snapshot size must be smaller than the data disk size.",
						},
						"data_disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Data disk ID used to initialize the data disk. When data disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: "Decides whether the disk is deleted with instance(only applied to `CLOUD_BASIC`, `CLOUD_SSD` and `CLOUD_PREMIUM` disk with `POSTPAID_BY_HOUR` instance), default is true.",
						},
						"encrypt": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Decides whether the disk is encrypted. Default is `false`.",
						},
						"throughput_performance": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							ForceNew:    true,
							Description: "Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD`.",
						},
					},
				},
			},
			// enhance services
			"disable_security_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable enhance service for security, it is enabled by default. When this options is set, security agent won't be installed.",
			},
			"disable_monitor_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable enhance service for monitor, it is enabled by default. When this options is set, monitor agent won't be installed.",
			},
			// login
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The key pair to use for the instance, it looks like `skey-16jig7tx`.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password for the instance. In order for the new password to take effect, the instance will be restarted after the password change.",
			},
			"keep_image_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "false" && old == "" || old == "false" && new == "" {
						return true
					} else {
						return old == new
					}
				},
				ConflictsWith: []string{"key_name", "password"},
				Description:   "Whether to keep image login or not, default is `false`. When the image type is private or shared or imported, this parameter can be set `true`. Modifying will cause the instance reset.",
			},
			"user_data": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user_data_raw"},
				Description:   "The user data to be injected into this instance. Must be base64 encoded and up to 16 KB.",
			},
			"user_data_raw": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user_data"},
				Description:   "The user data to be injected into this instance, in plain text. Conflicts with `user_data`. Up to 16 KB after base64 encoded.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A mapping of tags to assign to the resource. For tag limits, please refer to [Use Limits](https://intl.cloud.tencent.com/document/product/651/13354).",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate whether to force delete the instance. Default is `false`. If set true, the instance will be permanently deleted instead of being moved into the recycle bin. Note: only works for `PREPAID` instance.",
			},
			// role
			"cam_role_name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "CAM role name authorized to access.",
			},
			// Computed values.
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the instance.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP of the instance.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Create time of the instance.",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expired time of the instance.",
			},
		},
	}
}

func resourceTencentCloudInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_instance.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := cvm.NewRunInstancesRequest()
	request.ImageId = helper.String(d.Get("image_id").(string))
	request.Placement = &cvm.Placement{
		Zone: helper.String(d.Get("availability_zone").(string)),
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.Placement.ProjectId = &projectId
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("instance_count"); ok {
		request.InstanceCount = helper.Int64(int64(v.(int)))
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("hostname"); ok {
		request.HostName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("cam_role_name"); ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		instanceChargeType := v.(string)
		request.InstanceChargeType = &instanceChargeType
		if instanceChargeType == CVM_CHARGE_TYPE_PREPAID {
			request.InstanceChargePrepaid = &cvm.InstanceChargePrepaid{}
			if period, ok := d.GetOk("instance_charge_type_prepaid_period"); ok {
				periodInt64 := int64(period.(int))
				request.InstanceChargePrepaid.Period = &periodInt64
			} else {
				return fmt.Errorf("instance charge type prepaid period can not be empty when charge type is %s",
					instanceChargeType)
			}
			if renewFlag, ok := d.GetOk("instance_charge_type_prepaid_renew_flag"); ok {
				request.InstanceChargePrepaid.RenewFlag = helper.String(renewFlag.(string))
			}
		}
		if instanceChargeType == CVM_CHARGE_TYPE_SPOTPAID {
			request.InstanceMarketOptions = &cvm.InstanceMarketOptionsRequest{}
			request.InstanceMarketOptions.MarketType = helper.String(CVM_MARKET_TYPE_SPOT)
			request.InstanceMarketOptions.SpotOptions = &cvm.SpotMarketOptions{}
			if v, ok := d.GetOk("spot_instance_type"); ok {
				request.InstanceMarketOptions.SpotOptions.SpotInstanceType = helper.String(strings.ToLower(v.(string)))
			} else {
				return fmt.Errorf("spot_instance_type can not be empty when instance_charge_type is %s", instanceChargeType)
			}
			if v, ok := d.GetOk("spot_max_price"); ok {
				request.InstanceMarketOptions.SpotOptions.MaxPrice = helper.String(v.(string))
			} else {
				return fmt.Errorf("spot_max_price can not be empty when instance_charge_type is %s", instanceChargeType)
			}
		}
		if instanceChargeType == CVM_CHARGE_TYPE_CDHPAID {
			if v, ok := d.GetOk("cdh_instance_type"); ok {
				request.InstanceType = helper.String(v.(string))
			} else {
				return fmt.Errorf("cdh_instance_type can not be empty when instance_charge_type is %s", instanceChargeType)
			}
			if v, ok := d.GetOk("cdh_host_id"); ok {
				request.Placement.HostIds = append(request.Placement.HostIds, helper.String(v.(string)))
			} else {
				return fmt.Errorf("cdh_host_id can not be empty when instance_charge_type is %s", instanceChargeType)
			}
		}
	}
	if v, ok := d.GetOk("placement_group_id"); ok {
		request.DisasterRecoverGroupIds = []*string{helper.String(v.(string))}
	}

	// network
	request.InternetAccessible = &cvm.InternetAccessible{}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		maxBandwidthOut := int64(v.(int))
		request.InternetAccessible.InternetMaxBandwidthOut = &maxBandwidthOut
	}
	if v, ok := d.GetOk("bandwidth_package_id"); ok {
		request.InternetAccessible.BandwidthPackageId = helper.String(v.(string))
	}
	if v, ok := d.GetOkExists("allocate_public_ip"); ok {
		allocatePublicIp := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &allocatePublicIp
	}

	// vpc
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{}
		request.VirtualPrivateCloud.VpcId = helper.String(v.(string))

		if v, ok = d.GetOk("subnet_id"); ok {
			request.VirtualPrivateCloud.SubnetId = helper.String(v.(string))
		}

		if v, ok = d.GetOk("private_ip"); ok {
			request.VirtualPrivateCloud.PrivateIpAddresses = []*string{helper.String(v.(string))}
		}
	}

	if v, ok := d.GetOk("security_groups"); ok {
		securityGroups := v.(*schema.Set).List()
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for _, securityGroup := range securityGroups {
			request.SecurityGroupIds = append(request.SecurityGroupIds, helper.String(securityGroup.(string)))
		}
	}

	// storage
	request.SystemDisk = &cvm.SystemDisk{}
	if v, ok := d.GetOk("system_disk_type"); ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("system_disk_size"); ok {
		diskSize := int64(v.(int))
		request.SystemDisk.DiskSize = &diskSize
	}
	if v, ok := d.GetOk("system_disk_id"); ok {
		request.SystemDisk.DiskId = helper.String(v.(string))
	}
	if v, ok := d.GetOk("data_disks"); ok {
		dataDisks := v.([]interface{})
		request.DataDisks = make([]*cvm.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["data_disk_type"].(string)
			diskSize := int64(value["data_disk_size"].(int))
			throughputPerformance := int64(value["throughput_performance"].(int))
			dataDisk := cvm.DataDisk{
				DiskType:              &diskType,
				DiskSize:              &diskSize,
				ThroughputPerformance: &throughputPerformance,
			}
			if v, ok := value["data_disk_snapshot_id"]; ok && v != nil {
				snapshotId := v.(string)
				if snapshotId != "" {
					dataDisk.SnapshotId = helper.String(snapshotId)
				}
			}
			if value["data_disk_id"] != "" {
				dataDisk.DiskId = helper.String(value["data_disk_id"].(string))
			}
			if deleteWithInstance, ok := value["delete_with_instance"]; ok {
				deleteWithInstanceBool := deleteWithInstance.(bool)
				dataDisk.DeleteWithInstance = &deleteWithInstanceBool
			}

			if encrypt, ok := value["encrypt"]; ok {
				encryptBool := encrypt.(bool)
				dataDisk.Encrypt = &encryptBool
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	// enhanced service
	request.EnhancedService = &cvm.EnhancedService{}
	if v, ok := d.GetOkExists("disable_security_service"); ok {
		securityService := !(v.(bool))
		request.EnhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := d.GetOkExists("disable_monitor_service"); ok {
		monitorService := !(v.(bool))
		request.EnhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}

	// login
	request.LoginSettings = &cvm.LoginSettings{}
	if v, ok := d.GetOk("key_name"); ok {
		request.LoginSettings.KeyIds = []*string{helper.String(v.(string))}
	}
	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = helper.String(v.(string))
	}
	v := d.Get("keep_image_login").(bool)
	if v {
		request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN)
	} else {
		request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN_NOT)
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}
	if v, ok := d.GetOk("user_data_raw"); ok {
		userData := base64.StdEncoding.EncodeToString([]byte(v.(string)))
		request.UserData = &userData
	}

	instanceId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check("create")
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().RunInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			e, ok := err.(*sdkErrors.TencentCloudSDKError)
			if ok && e.Code == CVM_CLOUD_DISK_SOLD_OUT_ERROR {
				return resource.NonRetryableError(e)
			}
			return retryError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if len(response.Response.InstanceIdSet) < 1 {
			err = fmt.Errorf("instance id is nil")
			return resource.NonRetryableError(err)
		}
		instanceId = *response.Response.InstanceIdSet[0]

		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(instanceId)

	// wait for status
	//get system disk ID and data disk ID
	var systemDiskId string
	var dataDiskIds []string
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance != nil && (*instance.InstanceState == CVM_STATUS_RUNNING ||
			*instance.InstanceState == CVM_STATUS_LAUNCH_FAILED) {
			//get system disk ID
			if instance.SystemDisk != nil && instance.SystemDisk.DiskId != nil {
				systemDiskId = *instance.SystemDisk.DiskId
			}
			if instance.DataDisks != nil {
				for _, dataDisk := range instance.DataDisks {
					if dataDisk != nil && dataDisk.DiskId != nil {
						dataDiskIds = append(dataDiskIds, *dataDisk.DiskId)
					}
				}
			}
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})

	if err != nil {
		return err
	}

	// Wait for the tags attached to the vm since tags attachment it's async while vm creation.
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("cvm", "instance", tcClient.Region, instanceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			// If tags attachment failed, the user will be notified, then plan/apply/update with terraform.
			return err
		}

		//except instance ,system disk and data disk will be tagged
		//keep logical consistence with the console
		//tag system disk
		if systemDiskId != "" {
			resourceName = BuildTagResourceName("cvm", "volume", tcClient.Region, systemDiskId)
			if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
				// If tags attachment failed, the user will be notified, then plan/apply/update with terraform.
				return err
			}
		}
		//tag disk ids
		for _, diskId := range dataDiskIds {
			if diskId != "" {
				resourceName = BuildTagResourceName("cvm", "volume", tcClient.Region, diskId)
				if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
					// If tags attachment failed, the user will be notified, then plan/apply/update with terraform.
					return err
				}
			}
		}
	}

	if !(d.Get("running_flag").(bool)) {
		stoppedMode := d.Get("stopped_mode").(string)
		err = cvmService.StopInstance(ctx, instanceId, stoppedMode)
		if err != nil {
			return err
		}

		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			if instance != nil && *instance.InstanceState == CVM_STATUS_STOPPED {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
		})
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudInstanceRead(d, meta)
}

func resourceTencentCloudInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	forceDelete := false
	if v, ok := d.GetOkExists("force_delete"); ok {
		forceDelete = v.(bool)
		_ = d.Set("force_delete", forceDelete)
	}

	client := meta.(*TencentCloudClient).apiV3Conn
	cvmService := CvmService{
		client: client,
	}
	var instance *cvm.Instance
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, errRet = cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if instance == nil || *instance.InstanceState == CVM_STATUS_LAUNCH_FAILED {
		d.SetId("")
		return nil
	}

	var cvmImages []string
	var response *cvm.DescribeImagesResponse
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		request := cvm.NewDescribeImagesRequest()
		response, errRet = client.UseCvmClient().DescribeImages(request)
		if *response.Response.TotalCount > 0 {
			for i := range response.Response.ImageSet {
				image := response.Response.ImageSet[i]
				cvmImages = append(cvmImages, *image.ImageId)
			}
		}
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	if d.Get("image_id").(string) == "" || !IsContains(cvmImages, *instance.ImageId) {
		_ = d.Set("image_id", instance.ImageId)
	}

	_ = d.Set("availability_zone", instance.Placement.Zone)
	_ = d.Set("instance_name", instance.InstanceName)
	_ = d.Set("instance_type", instance.InstanceType)
	_ = d.Set("project_id", instance.Placement.ProjectId)
	_ = d.Set("instance_charge_type", instance.InstanceChargeType)
	_ = d.Set("instance_charge_type_prepaid_renew_flag", instance.RenewFlag)
	_ = d.Set("internet_charge_type", instance.InternetAccessible.InternetChargeType)
	_ = d.Set("internet_max_bandwidth_out", instance.InternetAccessible.InternetMaxBandwidthOut)
	_ = d.Set("vpc_id", instance.VirtualPrivateCloud.VpcId)
	_ = d.Set("subnet_id", instance.VirtualPrivateCloud.SubnetId)
	_ = d.Set("security_groups", helper.StringsInterfaces(instance.SecurityGroupIds))
	_ = d.Set("system_disk_type", instance.SystemDisk.DiskType)
	_ = d.Set("system_disk_size", instance.SystemDisk.DiskSize)
	_ = d.Set("system_disk_id", instance.SystemDisk.DiskId)
	_ = d.Set("instance_status", instance.InstanceState)
	_ = d.Set("create_time", instance.CreatedTime)
	_ = d.Set("expired_time", instance.ExpiredTime)
	_ = d.Set("cam_role_name", instance.CamRoleName)

	if *instance.InstanceChargeType == CVM_CHARGE_TYPE_CDHPAID {
		_ = d.Set("cdh_instance_type", instance.InstanceType)
	}

	if _, ok := d.GetOkExists("allocate_public_ip"); !ok {
		_ = d.Set("allocate_public_ip", len(instance.PublicIpAddresses) > 0)
	}

	// as attachment add tencentcloud:autoscaling:auto-scaling-group-id tag automatically
	// we should remove this tag, otherwise it will cause terraform state change
	tags := flattenCvmTagsMapping(instance.Tags)
	delete(tags, "tencentcloud:autoscaling:auto-scaling-group-id")
	_ = d.Set("tags", tags)

	//set data_disks
	dataDiskList := make([]map[string]interface{}, 0, len(instance.DataDisks))
	for _, disk := range instance.DataDisks {
		dataDisk := make(map[string]interface{}, 5)
		dataDisk["data_disk_snapshot_id"] = disk.SnapshotId
		dataDisk["data_disk_type"] = disk.DiskType
		dataDisk["data_disk_size"] = disk.DiskSize
		dataDisk["data_disk_id"] = disk.DiskId
		dataDisk["delete_with_instance"] = disk.DeleteWithInstance
		dataDisk["encrypt"] = disk.Encrypt
		dataDisk["throughput_performance"] = disk.ThroughputPerformance
		dataDiskList = append(dataDiskList, dataDisk)
	}
	_ = d.Set("data_disks", dataDiskList)

	if len(instance.PrivateIpAddresses) > 0 {
		_ = d.Set("private_ip", instance.PrivateIpAddresses[0])
	}
	if len(instance.PublicIpAddresses) > 0 {
		_ = d.Set("public_ip", instance.PublicIpAddresses[0])
	}
	if len(instance.LoginSettings.KeyIds) > 0 {
		_ = d.Set("key_name", instance.LoginSettings.KeyIds[0])
	} else {
		_ = d.Set("key_name", "")
	}
	if instance.LoginSettings.KeepImageLogin != nil {
		_ = d.Set("keep_image_login", *instance.LoginSettings.KeepImageLogin == CVM_IMAGE_LOGIN)
	}

	if *instance.InstanceState == CVM_STATUS_STOPPED {
		_ = d.Set("running_flag", false)
	} else {
		_ = d.Set("running_flag", true)
	}

	return nil
}

func resourceTencentCloudInstanceUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer logElapsed("resource.tencentcloud_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	instanceId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	d.Partial(true)

	unsupportedUpdateFields := []string{
		"instance_charge_type_prepaid_period",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_cvm_instance update on %s is not support yet", field)
		}
	}

	if d.HasChange("instance_charge_type_prepaid_renew_flag") {
		//check
		chargeType := d.Get("instance_charge_type").(string)
		if chargeType != CVM_CHARGE_TYPE_PREPAID {
			return fmt.Errorf("tencentcloud_cvm_instance update on instance_charge_type_prepaid_period or instance_charge_type_prepaid_renew_flag is only supported with charge type PREPAID")
		}

		//renew api
		err := cvmService.ModifyRenewParam(ctx, instanceId, d.Get("instance_charge_type_prepaid_renew_flag").(string))
		if err != nil {
			return err
		}

		//check success
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			if instance != nil && instance.LatestOperationState != nil {
				if *instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_SUCCESS {
					return nil
				} else if *instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_FAILED {
					return resource.NonRetryableError(fmt.Errorf("update instance %s prepaid charge type failed", instanceId))
				} else {
					return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
				}
			} else {
				return resource.RetryableError(fmt.Errorf("cvm instance %s returns nil status", instanceId))
			}

		})
		if err != nil {
			return err
		}

		time.Sleep(readRetryTimeout)
		d.SetPartial("instance_charge_type_prepaid_renew_flag")
	}

	if d.HasChange("instance_name") {
		err := cvmService.ModifyInstanceName(ctx, instanceId, d.Get("instance_name").(string))
		if err != nil {
			return err
		}
		d.SetPartial("instance_name")
	}

	if d.HasChange("security_groups") {
		securityGroups := d.Get("security_groups").(*schema.Set).List()
		securityGroupIds := make([]*string, 0, len(securityGroups))
		for _, securityGroup := range securityGroups {
			securityGroupIds = append(securityGroupIds, helper.String(securityGroup.(string)))
		}
		err := cvmService.ModifySecurityGroups(ctx, instanceId, securityGroupIds)
		if err != nil {
			return err
		}
		d.SetPartial("security_groups")
	}

	if d.HasChange("project_id") {
		projectId := d.Get("project_id").(int)
		err := cvmService.ModifyProjectId(ctx, instanceId, int64(projectId))
		if err != nil {
			return err
		}
		d.SetPartial("project_id")
	}

	// Reset Instance
	// Keep Login Info
	if d.HasChange("image_id") ||
		d.HasChange("host_name") ||
		d.HasChange("disable_security_service") ||
		d.HasChange("disable_monitor_service") ||
		d.HasChange("keep_image_login") {

		var updateAttr []string

		request := cvm.NewResetInstanceRequest()
		request.InstanceId = helper.String(d.Id())

		if v, ok := d.GetOk("image_id"); ok {
			updateAttr = append(updateAttr, "image_id")
			request.ImageId = helper.String(v.(string))
		}
		if v, ok := d.GetOk("hostname"); ok {
			updateAttr = append(updateAttr, "hostname")
			request.HostName = helper.String(v.(string))
		}

		// enhanced service
		request.EnhancedService = &cvm.EnhancedService{}
		if d.HasChange("disable_security_service") {
			updateAttr = append(updateAttr, "disable_security_service")
			v := d.Get("disable_security_service")
			securityService := v.(bool)
			request.EnhancedService.SecurityService = &cvm.RunSecurityServiceEnabled{
				Enabled: &securityService,
			}
		}

		if d.HasChange("disable_monitor_service") {
			updateAttr = append(updateAttr, "disable_monitor_service")
			v := d.Get("disable_monitor_service")
			monitorService := !(v.(bool))
			request.EnhancedService.MonitorService = &cvm.RunMonitorServiceEnabled{
				Enabled: &monitorService,
			}
		}

		// Modify or keep login info when instance reset
		request.LoginSettings = &cvm.LoginSettings{}

		if v, ok := d.GetOk("password"); ok {
			updateAttr = append(updateAttr, "password")
			request.LoginSettings.Password = helper.String(v.(string))
		}

		if v, ok := d.GetOk("key_name"); ok {
			updateAttr = append(updateAttr, "key_name")
			request.LoginSettings.KeyIds = []*string{helper.String(v.(string))}
		}

		if d.HasChange("keep_image_login") {
			updateAttr = append(updateAttr, "keep_image_login")
		}

		if v := d.Get("keep_image_login").(bool); v {
			request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN)
		} else {
			request.LoginSettings.KeepImageLogin = helper.String(CVM_IMAGE_LOGIN_NOT)
		}

		if err := cvmService.ResetInstance(ctx, request); err != nil {
			return err
		}

		for _, attr := range updateAttr {
			d.SetPartial(attr)
		}

		// Modify Login Info Directly
	} else {
		if d.HasChange("password") {
			err := cvmService.ModifyPassword(ctx, instanceId, d.Get("password").(string))
			if err != nil {
				return err
			}
			d.SetPartial("password")
			time.Sleep(10 * time.Second)
			err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
				instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
				if errRet != nil {
					return retryError(errRet, InternalError)
				}
				if instance != nil && *instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_OPERATING {
					return resource.RetryableError(fmt.Errorf("cvm instance latest operetion status is %s, retry...", *instance.LatestOperationState))
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		if d.HasChange("key_name") {
			old, new := d.GetChange("key_name")
			oldKeyId := old.(string)
			keyId := new.(string)
			if oldKeyId != "" {
				err := cvmService.UnbindKeyPair(ctx, oldKeyId, []*string{&instanceId})
				if err != nil {
					return err
				}
				time.Sleep(10 * time.Second)
				err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
					instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
					if errRet != nil {
						return retryError(errRet, InternalError)
					}
					if instance != nil && *instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_OPERATING {
						return resource.RetryableError(fmt.Errorf("cvm instance latest operetion status is %s, retry...", *instance.LatestOperationState))
					}
					return nil
				})
				if err != nil {
					return err
				}
			}

			if keyId != "" {
				err = cvmService.BindKeyPair(ctx, keyId, instanceId)
				if err != nil {
					return err
				}
				time.Sleep(10 * time.Second)
				err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
					instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
					if errRet != nil {
						return retryError(errRet, InternalError)
					}
					if instance != nil && *instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_OPERATING {
						return resource.RetryableError(fmt.Errorf("cvm instance latest operetion status is %s, retry...", *instance.LatestOperationState))
					}
					return nil
				})
				if err != nil {
					return err
				}
			}
			d.SetPartial("key_name")
		}
	}

	var flag bool
	if d.HasChange("running_flag") {
		flag = d.Get("running_flag").(bool)
		if flag {
			err = cvmService.StartInstance(ctx, instanceId)
			if err != nil {
				return err
			}
			err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
				instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
				if errRet != nil {
					return retryError(errRet, InternalError)
				}
				if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
					return nil
				}
				return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
			})
			if err != nil {
				return err
			}
		} else {
			stoppedMode := d.Get("stopped_mode").(string)
			err = cvmService.StopInstance(ctx, instanceId, stoppedMode)
			if err != nil {
				return err
			}
			err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
				instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
				if errRet != nil {
					return retryError(errRet, InternalError)
				}
				if instance != nil && *instance.InstanceState == CVM_STATUS_STOPPED {
					return nil
				}
				return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
			})
			if err != nil {
				return err
			}
		}
		d.SetPartial("running_flag")
	}

	if d.HasChange("instance_type") {
		err := cvmService.ModifyInstanceType(ctx, instanceId, d.Get("instance_type").(string))
		if err != nil {
			return err
		}
		d.SetPartial("instance_type")

		// wait for status
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			// Modifying instance type need restart the instance
			// so status of CVM must be running when running flag is true
			if instance != nil && (*instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_OPERATING ||
				(flag && *instance.InstanceState != CVM_STATUS_RUNNING)) {
				return resource.RetryableError(fmt.Errorf("cvm instance latest operetion status is %s, retry...", *instance.LatestOperationState))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("cdh_instance_type") {
		err := cvmService.ModifyInstanceType(ctx, instanceId, d.Get("cdh_instance_type").(string))
		if err != nil {
			return err
		}
		d.SetPartial("cdh_instance_type")

		// wait for status
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			// Modifying instance type need restart the instance
			// so status of CVM must be running when running flag is true
			if instance != nil && (*instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_OPERATING ||
				(flag && *instance.InstanceState != CVM_STATUS_RUNNING)) {
				return resource.RetryableError(fmt.Errorf("cvm instance latest operetion status is %s, retry...", *instance.LatestOperationState))
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("vpc_id") || d.HasChange("subnet_id") || d.HasChange("private_ip") {
		vpcId := d.Get("vpc_id").(string)
		subnetId := d.Get("subnet_id").(string)
		privateIp := d.Get("private_ip").(string)
		err := cvmService.ModifyVpc(ctx, instanceId, vpcId, subnetId, privateIp)
		if err != nil {
			return err
		}
		if d.HasChange("vpc_id") {
			d.SetPartial("vpc_id")
		}
		if d.HasChange("subnet_id") {
			d.SetPartial("subnet_id")
		}
		if d.HasChange("private_ip") {
			d.SetPartial("private_ip")
		}
	}

	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("cvm", "instance", region, instanceId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		//except instance ,system disk and data disk will be tagged
		//keep logical consistence with the console
		//tag system disk
		if systemDiskId, ok := d.GetOk("system_disk_id"); ok {
			if systemDiskId.(string) != "" {
				resourceName = BuildTagResourceName("cvm", "volume", region, systemDiskId.(string))
				if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
					return err
				}
			}
		}
		//tag disk ids
		if dataDisks, ok := d.GetOk("date_disk"); ok {
			dataDiskList := dataDisks.([]map[string]interface{})
			for _, disk := range dataDiskList {
				dataDiskId := disk["data_disk_id"].(string)
				resourceName = BuildTagResourceName("cvm", "volume", region, dataDiskId)
				if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
					return err
				}
			}
		}
		d.SetPartial("tags")
	}

	if d.HasChange("internet_max_bandwidth_out") {
		chargeType := d.Get("internet_charge_type").(string)
		if chargeType != "TRAFFIC_POSTPAID_BY_HOUR" && chargeType != "BANDWIDTH_POSTPAID_BY_HOUR" && chargeType != "BANDWIDTH_PACKAGE" {
			return fmt.Errorf("charge type should be one of `TRAFFIC_POSTPAID_BY_HOUR BANDWIDTH_POSTPAID_BY_HOUR BANDWIDTH_PACKAGE` when adjusting internet_max_bandwidth_out")
		}

		err := cvmService.ModifyInternetMaxBandwidthOut(ctx, instanceId, chargeType, int64(d.Get("internet_max_bandwidth_out").(int)))
		if err != nil {
			return err
		}
		d.SetPartial("internet_max_bandwidth_out")
		time.Sleep(1 * time.Second)
		err = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, InternalError)
			}
			if instance != nil && *instance.LatestOperationState == CVM_LATEST_OPERATION_STATE_OPERATING {
				return resource.RetryableError(fmt.Errorf("cvm instance latest operetion status is %s, retry...", *instance.LatestOperationState))
			}
			return nil
		})
		if err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudInstanceRead(d, meta)
}

func resourceTencentCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()
	//check is force delete or not
	forceDelete := d.Get("force_delete").(bool)

	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := cvmService.DeleteInstance(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	//check recycling
	notExist := false

	//check exist
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance == nil {
			notExist = true
			return nil
		}
		if *instance.InstanceState == CVM_STATUS_SHUTDOWN && *instance.LatestOperationState != CVM_LATEST_OPERATION_STATE_OPERATING {
			//in recycling
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})
	if err != nil {
		return err
	}

	if notExist || !forceDelete {
		return nil
	}

	// exist in recycle, delete again
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := cvmService.DeleteInstance(ctx, instanceId)
		//when state is terminating, do not delete but check exist
		if errRet != nil {
			//check InvalidInstanceState.Terminating
			ee, ok := errRet.(*errors.TencentCloudSDKError)
			if !ok {
				return retryError(errRet)
			}
			if ee.Code == "InvalidInstanceState.Terminating" {
				return nil
			}
			return retryError(errRet, "OperationDenied.InstanceOperationInProgress")
		}
		return nil
	})
	if err != nil {
		return err
	}

	//describe and check not exist
	err = resource.Retry(5*readRetryTimeout, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		if instance == nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})
	if err != nil {
		return err
	}

	return nil
}
