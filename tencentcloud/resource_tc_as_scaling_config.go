/*
Provides a resource to create a configuration for an AS (Auto scaling) instance.

Example Usage

```hcl
resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "launch-configuration"
  image_id           = "img-9qabwvbn"
  instance_types     = ["SA1.SMALL1"]
  project_id         = 0
  system_disk_type   = "CLOUD_PREMIUM"
  system_disk_size   = "50"

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
  user_data                  = "dGVzdA=="

  instance_tags = {
    tag = "as"
  }
}
```

Import

AutoScaling Configuration can be imported using the id, e.g.

```
$ terraform import tencentcloud_as_scaling_config.scaling_config asc-n32ymck2
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudAsScalingConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudAsScalingConfigCreate,
		Read:   resourceTencentCloudAsScalingConfigRead,
		Update: resourceTencentCloudAsScalingConfigUpdate,
		Delete: resourceTencentCloudAsScalingConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"configuration_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 60),
				Description:  "Name of a launch configuration.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "An available image ID for a cvm instance.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Specifys to which project the configuration belongs.",
			},
			"instance_types": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    5,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specified types of CVM instances.",
			},
			"system_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
				ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
				Description:  "Type of a CVM disk, and available values include CLOUD_PREMIUM and CLOUD_SSD. Default is CLOUD_PREMIUM.",
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: validateIntegerInRange(50, 500),
				Description:  "Volume of system disk in GB. Default is 50.",
			},
			"data_disk": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    11,
				Description: "Configurations of data disk.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
							ValidateFunc: validateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
							Description:  "Types of disk, available values: CLOUD_PREMIUM and CLOUD_SSD.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Volume of disk in GB. Default is 0.",
						},
						"snapshot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Data disk snapshot ID.",
						},
					},
				},
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID_BY_HOUR,
				ValidateFunc: validateAllowedStringValue(INTERNET_CHARGE_ALLOW_TYPE),
				Description:  "Charge types for network traffic. Available values include `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
			},
			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Max bandwidth of Internet access in Mbps. Default is 0.",
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
				ValidateFunc:  validateAsConfigPassword,
				ConflictsWith: []string{"keep_image_login"},
				Description:   "Password to access.",
			},
			"key_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"keep_image_login"},
				Description:   "ID list of keys.",
			},
			"keep_image_login": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"password", "key_ids"},
				Description:   "Specify whether to keep original settings of a CVM image. And it can't be used with password or key_ids together.",
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
				Description: "To specify whether to enable cloud security service. Default is TRUE.",
			},
			"enhanced_monitor_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "To specify whether to enable cloud monitor service. Default is TRUE.",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ase64-encoded User Data text, the length limit is 16KB.",
			},
			"instance_tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "A list of tags used to associate different resources.",
			},

			// Computed values
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current statues of a launch configuration.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the launch configuration was created.",
			},
		},
	}
}

func resourceTencentCloudAsScalingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_config.create")()

	logId := getLogId(contextNil)
	request := as.NewCreateLaunchConfigurationRequest()

	v := d.Get("configuration_name")
	request.LaunchConfigurationName = helper.String(v.(string))

	v = d.Get("image_id")
	request.ImageId = helper.String(v.(string))

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntUint64(v.(int))
	}

	v = d.Get("instance_types")
	instanceTypes := v.([]interface{})
	request.InstanceTypes = make([]*string, 0, len(instanceTypes))
	for i := range instanceTypes {
		instanceType := instanceTypes[i].(string)
		request.InstanceTypes = append(request.InstanceTypes, &instanceType)
	}

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := d.GetOk("system_disk_type"); ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("system_disk_size"); ok {
		request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("data_disk"); ok {
		dataDisks := v.([]interface{})
		request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["disk_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			snapshotId := value["snapshot_id"].(string)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
				DiskSize: &diskSize,
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	request.InternetAccessible = &as.InternetAccessible{}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOkExists("public_ip_assigned"); ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	request.LoginSettings = &as.LoginSettings{}
	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = helper.String(v.(string))
	}
	if v, ok := d.GetOk("key_ids"); ok {
		keyIds := v.([]interface{})
		request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
		}
	}
	if v, ok := d.GetOk("keep_image_login"); ok {
		keepImageLogin := v.(bool)
		request.LoginSettings.KeepImageLogin = &keepImageLogin
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
		}
	}

	request.EnhancedService = &as.EnhancedService{}

	if v, ok := d.GetOkExists("enhanced_security_service"); ok {
		securityService := v.(bool)
		request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := d.GetOkExists("enhanced_monitor_service"); ok {
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &as.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}

	chargeType := INSTANCE_CHARGE_TYPE_POSTPAID
	request.InstanceChargeType = &chargeType

	if v, ok := d.GetOk("instance_types_check_policy"); ok {
		request.InstanceTypesCheckPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_tags"); ok {
		tags := v.(map[string]interface{})
		request.InstanceTags = make([]*as.InstanceTag, 0, len(tags))
		for k, t := range tags {
			key := k
			value := t.(string)
			tag := as.InstanceTag{
				Key:   &key,
				Value: &value,
			}
			request.InstanceTags = append(request.InstanceTags, &tag)
		}
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().CreateLaunchConfiguration(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	if response.Response.LaunchConfigurationId == nil {
		return fmt.Errorf("Launch configuration id is nil")
	}
	d.SetId(*response.Response.LaunchConfigurationId)

	return resourceTencentCloudAsScalingConfigRead(d, meta)
}

func resourceTencentCloudAsScalingConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	configurationId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		config, has, e := asService.DescribeLaunchConfigurationById(ctx, configurationId)
		if e != nil {
			return retryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}
		_ = d.Set("configuration_name", *config.LaunchConfigurationName)
		_ = d.Set("status", *config.LaunchConfigurationStatus)
		_ = d.Set("image_id", *config.ImageId)
		_ = d.Set("project_id", *config.ProjectId)
		_ = d.Set("instance_types", helper.StringsInterfaces(config.InstanceTypes))
		_ = d.Set("system_disk_type", *config.SystemDisk.DiskType)
		_ = d.Set("system_disk_size", *config.SystemDisk.DiskSize)
		_ = d.Set("data_disk", flattenDataDiskMappings(config.DataDisks))
		_ = d.Set("internet_charge_type", *config.InternetAccessible.InternetChargeType)
		_ = d.Set("internet_max_bandwidth_out", *config.InternetAccessible.InternetMaxBandwidthOut)
		_ = d.Set("public_ip_assigned", *config.InternetAccessible.PublicIpAssigned)
		_ = d.Set("login_settings.key_ids", helper.StringsInterfaces(config.LoginSettings.KeyIds))
		_ = d.Set("security_group_ids", helper.StringsInterfaces(config.SecurityGroupIds))
		_ = d.Set("enhanced_security_service", *config.EnhancedService.SecurityService.Enabled)
		_ = d.Set("enhanced_monitor_service", *config.EnhancedService.MonitorService.Enabled)
		_ = d.Set("user_data", helper.PString(config.UserData))
		_ = d.Set("instance_tags", flattenInstanceTagsMapping(config.InstanceTags))
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAsScalingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_config.update")()

	logId := getLogId(contextNil)
	request := as.NewUpgradeLaunchConfigurationRequest()

	configurationId := d.Id()
	request.LaunchConfigurationId = &configurationId

	v := d.Get("configuration_name")
	request.LaunchConfigurationName = helper.String(v.(string))

	v = d.Get("image_id")
	request.ImageId = helper.String(v.(string))

	if v, ok := d.GetOk("project_id"); ok {
		projectId := int64(v.(int))
		request.ProjectId = &projectId
	}

	v = d.Get("instance_types")
	instanceTypes := v.([]interface{})
	request.InstanceTypes = make([]*string, 0, len(instanceTypes))
	for i := range instanceTypes {
		instanceType := instanceTypes[i].(string)
		request.InstanceTypes = append(request.InstanceTypes, &instanceType)
	}

	request.SystemDisk = &as.SystemDisk{}
	if v, ok := d.GetOk("system_disk_type"); ok {
		request.SystemDisk.DiskType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("system_disk_size"); ok {
		request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("data_disk"); ok {
		dataDisks := v.([]interface{})
		request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
		for _, d := range dataDisks {
			value := d.(map[string]interface{})
			diskType := value["disk_type"].(string)
			diskSize := uint64(value["disk_size"].(int))
			snapshotId := value["snapshot_id"].(string)
			dataDisk := as.DataDisk{
				DiskType: &diskType,
				DiskSize: &diskSize,
			}
			if snapshotId != "" {
				dataDisk.SnapshotId = &snapshotId
			}
			request.DataDisks = append(request.DataDisks, &dataDisk)
		}
	}

	request.InternetAccessible = &as.InternetAccessible{}
	if v, ok := d.GetOk("internet_charge_type"); ok {
		request.InternetAccessible.InternetChargeType = helper.String(v.(string))
	}
	if v, ok := d.GetOk("internet_max_bandwidth_out"); ok {
		request.InternetAccessible.InternetMaxBandwidthOut = helper.IntUint64(v.(int))
	}
	if v, ok := d.GetOkExists("public_ip_assigned"); ok {
		publicIpAssigned := v.(bool)
		request.InternetAccessible.PublicIpAssigned = &publicIpAssigned
	}

	request.LoginSettings = &as.LoginSettings{}
	if v, ok := d.GetOk("password"); ok {
		request.LoginSettings.Password = helper.String(v.(string))
	}
	if v, ok := d.GetOk("key_ids"); ok {
		keyIds := v.([]interface{})
		request.LoginSettings.KeyIds = make([]*string, 0, len(keyIds))
		for i := range keyIds {
			keyId := keyIds[i].(string)
			request.LoginSettings.KeyIds = append(request.LoginSettings.KeyIds, &keyId)
		}
	}
	if v, ok := d.GetOk("keep_image_login"); ok {
		keepImageLogin := v.(bool)
		request.LoginSettings.KeepImageLogin = &keepImageLogin
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		securityGroups := v.([]interface{})
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for i := range securityGroups {
			securityGroup := securityGroups[i].(string)
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
		}
	}

	request.EnhancedService = &as.EnhancedService{}

	if v, ok := d.GetOkExists("enhanced_security_service"); ok {
		securityService := v.(bool)
		request.EnhancedService.SecurityService = &as.RunSecurityServiceEnabled{
			Enabled: &securityService,
		}
	}
	if v, ok := d.GetOkExists("enhanced_monitor_service"); ok {
		monitorService := v.(bool)
		request.EnhancedService.MonitorService = &as.RunMonitorServiceEnabled{
			Enabled: &monitorService,
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}

	chargeType := INSTANCE_CHARGE_TYPE_POSTPAID
	request.InstanceChargeType = &chargeType

	if v, ok := d.GetOk("instance_types_check_policy"); ok {
		request.InstanceTypesCheckPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_tags"); ok {
		tags := v.(map[string]interface{})
		request.InstanceTags = make([]*as.InstanceTag, 0, len(tags))
		for k, t := range tags {
			key := k
			value := t.(string)
			tag := as.InstanceTag{
				Key:   &key,
				Value: &value,
			}
			request.InstanceTags = append(request.InstanceTags, &tag)
		}
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseAsClient().UpgradeLaunchConfiguration(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}

	return nil
}

func resourceTencentCloudAsScalingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_as_scaling_config.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	configurationId := d.Id()
	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := asService.DeleteLaunchConfiguration(ctx, configurationId)
	if err != nil {
		return err
	}

	return nil
}
