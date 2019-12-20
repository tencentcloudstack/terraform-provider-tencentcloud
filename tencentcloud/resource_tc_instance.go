/*
Provides a CVM instance resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** At present, 'PREPAID' instance cannot be deleted and must wait it to be outdated and released automatically.

Example Usage

```hcl
data "tencentcloud_images" "my_favorate_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "my_favorate_zones" {
}

// Create VPC resource
resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = "${tencentcloud_vpc.app.id}"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

// Create 2 CVM instances to host awesome_app
resource "tencentcloud_instance" "my_awesome_app" {
  instance_name              = "awesome_app"
  availability_zone          = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id                   = "${data.tencentcloud_images.my_favorate_image.images.0.image_id}"
  instance_type              = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = "${tencentcloud_vpc.app.id}"
  subnet_id                  = "${tencentcloud_subnet.app.id}"
  internet_max_bandwidth_out = 20
  count                      = 2

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
  }

  tags = {
    tagKey = "tagValue"
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
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
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
				ForceNew:    true,
				Description: "The Image to use for the instance. Change 'image_id' will case instance destroy and re-created.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone that the CVM instance locates at.",
			},
			"instance_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Terrafrom-CVM-Instance",
				ValidateFunc: validateStringLengthInRange(2, 128),
				Description:  "The name of the CVM. The max length of instance_name is 60, and default value is `Terrafrom-CVM-Instance`.",
			},
			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateInstanceType,
				Description:  "The type of instance to start.",
			},
			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The hostname of CVM. Windows instance: The name should be a combination of 2 to 15 characters comprised of letters (case insensitive), numbers, and hyphens (-). Period (.) is not supported, and the name cannot be a string of pure numbers. Other types (such as Linux) of instances: The name should be a combination of 2 to 60 characters, supporting multiple periods (.). The piece between two periods is composed of letters (case insensitive), numbers, and hyphens (-).",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The project CVM belongs to, default to 0.",
			},
			"running_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Set instance to running or stop. Default value is true, the instance will shutdown when flag is false.",
			},
			"placement_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The id of a placement group.",
			},
			// payment
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      CVM_CHARGE_TYPE_POSTPAID,
				ValidateFunc: validateAllowedStringValue(CVM_CHARGE_TYPE),
				Description:  "The charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR` and `SPOTPAID`, The default is `POSTPAID_BY_HOUR`. Note: TencentCloud International only supports `POSTPAID_BY_HOUR`. `PREPAID` instance may not allow to delete before expired. `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time.",
			},
			"instance_charge_type_prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue(CVM_PREPAID_PERIOD),
				Description:  "The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.",
			},
			"instance_charge_type_prepaid_renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CVM_PREPAID_RENEW_FLAG),
				Description:  "When enabled, the CVM instance will be renew automatically when it reach the end of the prepaid tenancy. Valid values are `NOTIFY_AND_AUTO_RENEW`, `NOTIFY_AND_MANUAL_RENEW` and `DISABLE_NOTIFY_AND_MANUAL_RENEW`. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
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
				Description:  "Max price of spot instance, is the format of decimal string, for example \"0.50\". Note: it only works when instance_charge_type is set to `SPOTPAID`.",
			},
			// network
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CVM_INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CVM_INTERNET_CHARGE_TYPE),
				Description:  "Internet charge type of the instance, Valid values are `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `BANDWIDTH_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`. The default is `TRAFFIC_POSTPAID_BY_HOUR`.",
			},
			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				ForceNew:    true,
				Description: "Maximum outgoing bandwidth to the public network, measured in Mbps (Mega bit per second). If this value is not specified, then automatically sets it to 0 Mbps.",
			},
			"allocate_public_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "Associate a public ip address with an instance in a VPC or Classic. Boolean value, Default is false.",
			},
			// vpc
			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of a VPC network. If you want to create instances in VPC network, this parameter must be set.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of a VPC subnetwork. If you want to create instances in VPC network, this parameter must be set.",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The private ip to be assigned to this instance, must be in the provided subnet and available.",
			},
			// security group
			"security_groups": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: "A list of security group ids to associate with.",
			},
			// storage
			"system_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      CVM_DISK_TYPE_CLOUD_BASIC,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(CVM_DISK_TYPE),
				Description:  "Type of the system disk. Valid values are `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_SSD` and `CLOUD_PREMIUM`, default value is `CLOUD_BASIC`. NOTE: `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.",
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: validateIntegerInRange(50, 1000),
				Description:  "Size of the system disk. Value range: [50, 1000], and unit is GB. Default is 50GB.",
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
				Description: "Settings for data disk.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_disk_type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue(CVM_DISK_TYPE),
							Description:  "Type of the data disk. Valid values are `LOCAL_BASIC`, `LOCAL_SSD`, `CLOUD_BASIC`, `CLOUD_SSD` and `CLOUD_PREMIUM`. NOTE: `LOCAL_BASIC` and `LOCAL_SSD` are deprecated.",
						},
						"data_disk_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validateIntegerInRange(10, 16000),
							Description:  "Size of the data disk, and unit is GB. If disk type is `CLOUD_SSD`, the size range is [100, 16000], and the others are [10-16000].",
						},
						"data_disk_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Data disk snapshot ID used to initialize the data disk. When data disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Decides whether the disk is deleted with instance(only applied to cloud disk), default to true.",
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
				Description: "The key pair to use for the instance, it looks like skey-16jig7tx.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password to an instance. In order to take effect new password, the instance will be restarted after modifying the password.",
			},
			"user_data": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user_data_raw"},
				Description:   "The user data to be specified into this instance. Must be encrypted in base64 format and limited in 16 KB.",
			},
			"user_data_raw": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"user_data"},
				Description:   "The user data to be specified into this instance, plain text. Conflicts with `user_data`. Limited in 16 KB after encrypted in base64 format.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A mapping of tags to assign to the resource. For tag limits, please refer to [Use Limits](https://intl.cloud.tencent.com/document/product/651/13354).",
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
				Description: "Public ip of the instance.",
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
	ctx := context.WithValue(context.TODO(), "logId", logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	request := cvm.NewRunInstancesRequest()
	request.ImageId = stringToPointer(d.Get("image_id").(string))
	request.Placement = &cvm.Placement{
		Zone: stringToPointer(d.Get("availability_zone").(string)),
	}
	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.Placement.ProjectId = &projectId
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("hostname"); ok {
		request.HostName = stringToPointer(v.(string))
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
				request.InstanceChargePrepaid.RenewFlag = stringToPointer(renewFlag.(string))
			}
		}
		if instanceChargeType == CVM_CHARGE_TYPE_SPOTPAID {
			request.InstanceMarketOptions = &cvm.InstanceMarketOptionsRequest{}
			request.InstanceMarketOptions.MarketType = stringToPointer(CVM_MARKET_TYPE_SPOT)
			request.InstanceMarketOptions.SpotOptions = &cvm.SpotMarketOptions{}
			if v, ok := d.GetOk("spot_instance_type"); ok {
				request.InstanceMarketOptions.SpotOptions.SpotInstanceType = stringToPointer(strings.ToLower(v.(string)))
			} else {
				return fmt.Errorf("spot_instance_type can not be empty when instance_charge_type is %s", instanceChargeType)
			}
			if v, ok := d.GetOk("spot_max_price"); ok {
				request.InstanceMarketOptions.SpotOptions.MaxPrice = stringToPointer(v.(string))
			} else {
				return fmt.Errorf("spot_max_price can not be empty when instance_charge_type is %s", instanceChargeType)
			}
		}
	}
	if v, ok := d.GetOk("placement_group_id"); ok {
		request.DisasterRecoverGroupIds = []*string{stringToPointer(v.(string))}
	}

	// network
	request.InternetAccessible = &cvm.InternetAccessible{}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetAccessible.InternetChargeType = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		maxBandwidthOut := int64(v.(int))
		request.InternetAccessible.InternetMaxBandwidthOut = &maxBandwidthOut
	}
	if v, ok := d.GetOkExists("allocate_public_ip"); ok {
		allocatePublicIp := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &allocatePublicIp
	}

	// vpc
	request.VirtualPrivateCloud = &cvm.VirtualPrivateCloud{}
	if v, ok := d.GetOk("vpc_id"); ok {
		request.VirtualPrivateCloud.VpcId = stringToPointer(v.(string))
	} else {
		request.VirtualPrivateCloud.VpcId = stringToPointer("DEFAULT")
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		request.VirtualPrivateCloud.SubnetId = stringToPointer(v.(string))
	} else {
		request.VirtualPrivateCloud.SubnetId = stringToPointer("DEFAULT")
	}
	if v, ok := d.GetOk("private_ip"); ok {
		request.VirtualPrivateCloud.PrivateIpAddresses = []*string{stringToPointer(v.(string))}
	}

	if v, ok := d.GetOk("security_groups"); ok {
		securityGroups := v.(*schema.Set).List()
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for _, securityGroup := range securityGroups {
			request.SecurityGroupIds = append(request.SecurityGroupIds, stringToPointer(securityGroup.(string)))
		}
	}

	// storage
	request.SystemDisk = &cvm.SystemDisk{}
	if v, ok := d.GetOk("system_disk_type"); ok {
		request.SystemDisk.DiskType = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("system_disk_size"); ok {
		diskSize := int64(v.(int))
		request.SystemDisk.DiskSize = &diskSize
	}
	if v, ok := d.GetOk("system_disk_id"); ok {
		request.SystemDisk.DiskId = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("data_disks"); ok {
		dataDisks := v.([]interface{})
		request.DataDisks = make([]*cvm.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["data_disk_type"].(string)
			diskSize := int64(value["data_disk_size"].(int))
			dataDisk := cvm.DataDisk{
				DiskType: &diskType,
				DiskSize: &diskSize,
			}
			if value["data_disk_id"] != "" {
				dataDisk.DiskId = stringToPointer(value["data_disk_id"].(string))
			}
			if deleteWithInstance, ok := value["delete_with_instance"]; ok {
				deleteWithInstanceBool := deleteWithInstance.(bool)
				dataDisk.DeleteWithInstance = &deleteWithInstanceBool
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
		request.LoginSettings.KeyIds = []*string{stringToPointer(v.(string))}
	}
	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = stringToPointer(v.(string))
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = stringToPointer(v.(string))
	}
	if v, ok := d.GetOk("user_data_raw"); ok {
		userData := base64.StdEncoding.EncodeToString([]byte(v.(string)))
		request.UserData = &userData
	}

	// tags
	if v, ok := d.GetOk("tags"); ok {
		tags := make([]*cvm.Tag, 0)
		for key, value := range v.(map[string]interface{}) {
			tag := &cvm.Tag{
				Key:   stringToPointer(key),
				Value: stringToPointer(value.(string)),
			}
			tags = append(tags, tag)
		}
		tagSpecification := &cvm.TagSpecification{
			ResourceType: stringToPointer("instance"),
			Tags:         tags,
		}
		request.TagSpecification = []*cvm.TagSpecification{tagSpecification}
	}

	instanceId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check("create")
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().RunInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
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
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, "InternalError")
		}
		if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
	})

	if err != nil {
		return err
	}

	if !(d.Get("running_flag").(bool)) {
		err = cvmService.StopInstance(ctx, instanceId)
		if err != nil {
			return err
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, "InternalError")
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

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	instanceId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var instance *cvm.Instance
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, errRet = cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, "InternalError")
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

	_ = d.Set("image_id", instance.ImageId)
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
	_ = d.Set("security_groups", flattenStringList(instance.SecurityGroupIds))
	_ = d.Set("system_disk_type", instance.SystemDisk.DiskType)
	_ = d.Set("system_disk_size", instance.SystemDisk.DiskSize)
	_ = d.Set("system_disk_id", instance.SystemDisk.DiskId)
	_ = d.Set("tags", flattenCvmTagsMapping(instance.Tags))
	_ = d.Set("instance_status", instance.InstanceState)
	_ = d.Set("create_time", instance.CreatedTime)
	_ = d.Set("expired_time", instance.ExpiredTime)

	dataDiskList := make([]map[string]interface{}, 0, len(instance.DataDisks))
	for _, disk := range instance.DataDisks {
		dataDisk := make(map[string]interface{}, 4)
		dataDisk["data_disk_type"] = disk.DiskType
		dataDisk["data_disk_size"] = disk.DiskSize
		dataDisk["data_disk_id"] = disk.DiskId
		dataDisk["delete_with_instance"] = disk.DeleteWithInstance
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
		_ = d.Set("key_name", instance.LoginSettings.KeyIds)
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
	ctx := context.WithValue(context.TODO(), "logId", logId)
	instanceId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	d.Partial(true)

	unsupportedUpdateFields := []string{
		"instance_charge_type_prepaid_period",
		"instance_charge_type_prepaid_renew_flag",
		"internet_charge_type",
		"internet_max_bandwidth_out",
		"allocate_public_ip",
		"system_disk_size",
		"data_disks",
	}
	for _, field := range unsupportedUpdateFields {
		if d.HasChange(field) {
			return fmt.Errorf("tencentcloud_cvm_instance update on %s is not support yet", field)
		}
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
			securityGroupIds = append(securityGroupIds, stringToPointer(securityGroup.(string)))
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

	if d.HasChange("instance_type") {
		err := cvmService.ModifyInstanceType(ctx, instanceId, d.Get("instance_type").(string))
		if err != nil {
			return err
		}
		d.SetPartial("instance_type")

		// wait for status
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
			if errRet != nil {
				return retryError(errRet, "InternalError")
			}
			if instance != nil && *instance.InstanceState == CVM_STATUS_RUNNING {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("cvm instance status is %s, retry...", *instance.InstanceState))
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("password") {
		err := cvmService.ModifyPassword(ctx, instanceId, d.Get("password").(string))
		if err != nil {
			return err
		}
		d.SetPartial("password")
		// just wait, 95% change password will be reset
		// waiting for status changed is not work
		// a little stupid, but it's ok for now
		time.Sleep(1 * time.Minute)
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

	if d.HasChange("running_flag") {
		flag := d.Get("running_flag").(bool)
		if flag {
			err = cvmService.StartInstance(ctx, instanceId)
			if err != nil {
				return err
			}
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
				if errRet != nil {
					return retryError(errRet, "InternalError")
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
			err = cvmService.StopInstance(ctx, instanceId)
			if err != nil {
				return err
			}
			err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
				if errRet != nil {
					return retryError(errRet, "InternalError")
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

	if d.HasChange("tags") {
		old, new := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(old.(map[string]interface{}), new.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := fmt.Sprintf("qcs::cvm:%s:uin/:instance/%s", region, instanceId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	d.Partial(false)

	return resourceTencentCloudInstanceRead(d, meta)
}

func resourceTencentCloudInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	instanceId := d.Id()
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

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		instance, errRet := cvmService.DescribeInstanceById(ctx, instanceId)
		if errRet != nil {
			return retryError(errRet, "InternalError")
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
