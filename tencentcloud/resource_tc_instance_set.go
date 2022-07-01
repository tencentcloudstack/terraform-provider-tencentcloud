/*
Provides a CVM instance set resource.

~> **NOTE:** You can launch an CVM instance for a VPC network via specifying parameter `vpc_id`. One instance can only belong to one VPC.

~> **NOTE:** This resource is designed to cater for the scenario of creating CVM in large batches.

~> **NOTE:** After run command `terraform apply`, must wait all cvms is ready, then run command `terraform plan`, either it will casue state change.

Example Usage

```hcl
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "Tencent Linux release 3.2 (Final)"
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

// Create 10 CVM instances to host awesome_app
resource "tencentcloud_instance_set" "my_awesome_app" {
  timeouts {
			create = "5m"
			read   = "20s"
			delete = "1h"
  }

  instance_count             = 10
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
}
```

*/
package tencentcloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func resourceTencentCloudInstanceSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudInstanceSetCreate,
		Read:   resourceTencentCloudInstanceSetRead,
		Update: resourceTencentCloudInstanceSetUpdate,
		Delete: resourceTencentCloudInstanceSetDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(600 * time.Second),
			Read:   schema.DefaultTimeout(600 * time.Second),
			Delete: schema.DefaultTimeout(600 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The image to use for the instance. Changing `image_id` will cause the instance reset.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The available zone for the CVM instance.",
			},
			"instance_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The number of instances to be purchased. Value range:[1,100]; default value: 1.",
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
				Default:      CVM_CHARGE_TYPE_POSTPAID,
				ValidateFunc: validateAllowedStringValue(CVM_CHARGE_TYPE),
				Description:  "The charge type of instance. Only support `POSTPAID_BY_HOUR`.",
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
				Default:      CVM_DISK_TYPE_CLOUD_PREMIUM,
				ValidateFunc: validateAllowedStringValue(CVM_DISK_TYPE),
				Description:  "System disk type. For more information on limits of system disk types, see [Storage Overview](https://intl.cloud.tencent.com/document/product/213/4952). Valid values: `LOCAL_BASIC`: local disk, `LOCAL_SSD`: local SSD disk, `CLOUD_SSD`: SSD, `CLOUD_PREMIUM`: Premium Cloud Storage. NOTE: 1. `CLOUD_BASIC`, `LOCAL_BASIC` and `LOCAL_SSD` are deprecated; 2. If modified, the instance may force stop.",
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: validateIntegerInRange(50, 1000),
				Description:  "Size of the system disk. Valid value ranges: (50~1000). and unit is GB. Default is 50GB. If modified, the instance may force stop.",
			},
			"system_disk_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "System disk snapshot ID used to initialize the system disk. When system disk type is `LOCAL_BASIC` and `LOCAL_SSD`, disk id is not supported.",
			},
			// enhance services
			"disable_security_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable enhance service for security, it is enabled by default. When this options is set, security agent won't be installed. Modifying will cause the instance reset.",
			},
			"disable_monitor_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable enhance service for monitor, it is enabled by default. When this options is set, monitor agent won't be installed. Modifying will cause the instance reset.",
			},
			// login
			"key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The key pair to use for the instance, it looks like `skey-16jig7tx`. Modifying will cause the instance reset.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password for the instance. In order for the new password to take effect, the instance will be restarted after the password change. Modifying will cause the instance reset.",
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
			"instance_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "instance id list.",
			},
		},
	}
}

func resourceTencentCloudInstanceSetCreate(d *schema.ResourceData, meta interface{}) error {
	doneChan := make(chan struct{}, 1)
	rspChan := make(chan error, 1)

	timeout := d.Timeout(schema.TimeoutCreate)

	go func(d *schema.ResourceData, meta interface{}) {
		e := doResourceTencentCloudInstanceSetCreate(d, meta)
		doneChan <- struct{}{}
		rspChan <- e
	}(d, meta)

	select {
	case <-doneChan:
		return <-rspChan
	case <-time.After(timeout):
		return fmt.Errorf("Do cvm instance set create action timeout, current timeout :[%.3f]s", timeout.Seconds())
	}
}

func resourceTencentCloudInstanceSetRead(d *schema.ResourceData, meta interface{}) error {
	doneChan := make(chan struct{}, 1)
	rspChan := make(chan error, 1)

	timeout := d.Timeout(schema.TimeoutRead)

	go func(d *schema.ResourceData, meta interface{}) {
		e := doResourceTencentCloudInstanceSetRead(d, meta)
		doneChan <- struct{}{}
		rspChan <- e
	}(d, meta)

	select {
	case <-doneChan:
		return <-rspChan
	case <-time.After(timeout):
		return fmt.Errorf("Do cvm instance set read action timeout, current timeout :[%.3f]s", timeout.Seconds())
	}
}

func resourceTencentCloudInstanceSetUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	defer logElapsed("resource.tencentcloud_instance_set.update")()

	return fmt.Errorf("`resource_instance_set` do not support change now.")
}

func resourceTencentCloudInstanceSetDelete(d *schema.ResourceData, meta interface{}) error {
	doneChan := make(chan struct{}, 1)
	rspChan := make(chan error, 1)

	timeout := d.Timeout(schema.TimeoutDelete)

	go func(d *schema.ResourceData, meta interface{}) {
		e := doResourceTencentCloudInstanceSetDelete(d, meta)
		doneChan <- struct{}{}
		rspChan <- e
	}(d, meta)

	select {
	case <-doneChan:
		return <-rspChan
	case <-time.After(timeout):
		return fmt.Errorf("Do cvm instance set delete action timeout, current timeout :[%.3f]s", timeout.Seconds())
	}
}

func doResourceTencentCloudInstanceSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_instance_set.create")()
	logId := getLogId(contextNil)

	var instanceCount int

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
		instanceCount = v.(int)
		request.InstanceCount = helper.Int64(int64(instanceCount))
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

	instanceIds := make([]*string, 0)

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check("create")
		response, err := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().RunInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			e, ok := err.(*sdkErrors.TencentCloudSDKError)
			if ok && IsContains(CVM_RETRYABLE_ERROR, e.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("cvm create error: %s, retrying", e.Error()))
			}
			return resource.NonRetryableError(err)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		if len(response.Response.InstanceIdSet) < instanceCount {
			err = fmt.Errorf("number of instances is less than %s", strconv.Itoa(instanceCount))
			return resource.NonRetryableError(err)
		}
		instanceIds = response.Response.InstanceIdSet

		return nil
	})
	if err != nil {
		return err
	}

	_ = d.Set("instance_ids", instanceIds)
	d.SetId(helper.StrListToStr(instanceIds))

	return nil
}

func doResourceTencentCloudInstanceSetRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_instance_set.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Id()

	client := meta.(*TencentCloudClient).apiV3Conn
	cvmService := CvmService{
		client: client,
	}
	var instanceSet []*cvm.Instance
	var errRet error
	instanceSet, errRet = cvmService.DescribeInstanceSetByIds(ctx, instanceId)
	if errRet != nil {
		return errRet
	}

	if instanceSet == nil {
		d.SetId("")
		return nil
	}

	instance := instanceSet[0]

	_ = d.Set("instance_count", len(instanceSet))
	_ = d.Set("image_id", instance.ImageId)
	_ = d.Set("availability_zone", instance.Placement.Zone)
	_ = d.Set("instance_name", d.Get("instance_name"))
	_ = d.Set("instance_type", instance.InstanceType)
	_ = d.Set("project_id", instance.Placement.ProjectId)
	_ = d.Set("instance_charge_type", instance.InstanceChargeType)
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

	if _, ok := d.GetOkExists("allocate_public_ip"); !ok {
		_ = d.Set("allocate_public_ip", len(instance.PublicIpAddresses) > 0)
	}

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

	return nil
}

func doResourceTencentCloudInstanceSetDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_instance_set.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceSetIds := d.Id()

	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	// delete
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := cvmService.DeleteInstanceSetByIds(ctx, instanceSetIds)
		if errRet != nil {
			log.Printf("[CRITAL][first delete]%s api[%s] fail, reason[%s]\n",
				logId, "delete", errRet.Error())
			e, ok := errRet.(*sdkErrors.TencentCloudSDKError)
			if ok && IsContains(CVM_RETRYABLE_ERROR, e.Code) {
				time.Sleep(1 * time.Second) // 需要重试的话，等待1s进行重试
				return resource.RetryableError(fmt.Errorf("[first delete]cvm delete error: %s, retrying", e.Error()))
			}
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
