package as

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccvm "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cvm"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	as "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/as/v20180419"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudAsScalingConfig() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateStringLengthInRange(1, 60),
				Description:  "Name of a launch configuration.",
			},
			"image_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"image_id", "image_family"},
				Description:  "An available image ID for a cvm instance.",
			},
			"image_family": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"image_id", "image_family"},
				Description:  "Image Family Name. Either Image ID or Image Family Name must be provided, but not both.",
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
				MaxItems:    10,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specified types of CVM instances.",
			},
			"system_disk_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      SYSTEM_DISK_TYPE_CLOUD_PREMIUM,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
				Description:  "Type of a CVM disk. Valid values: `CLOUD_PREMIUM` and `CLOUD_SSD`. Default is `CLOUD_PREMIUM`. valid when disk_type_policy is ORIGINAL.",
			},
			"system_disk_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      50,
				ValidateFunc: tccommon.ValidateIntegerInRange(50, 500),
				Description:  "Volume of system disk in GB. Default is `50`.",
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
							ValidateFunc: tccommon.ValidateAllowedStringValue(SYSTEM_DISK_ALLOW_TYPE),
							Description:  "Types of disk. Valid values: `CLOUD_PREMIUM` and `CLOUD_SSD`. valid when disk_type_policy is ORIGINAL.",
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
							Description: "Data disk snapshot ID.",
						},
						"delete_with_instance": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Indicates whether the disk remove after instance terminated. Default is `false`.",
						},
					},
				},
			},
			// payment
			"instance_charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Charge type of instance. Valid values are `PREPAID`, `POSTPAID_BY_HOUR`, `SPOTPAID`, `CDCPAID`. The default is `POSTPAID_BY_HOUR`. NOTE: `SPOTPAID` instance must set `spot_instance_type` and `spot_max_price` at the same time.",
			},
			"instance_charge_type_prepaid_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue(svccvm.CVM_PREPAID_PERIOD),
				Description:  "The tenancy (in month) of the prepaid instance, NOTE: it only works when instance_charge_type is set to `PREPAID`. Valid values are `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `10`, `11`, `12`, `24`, `36`.",
			},
			"instance_charge_type_prepaid_renew_flag": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(svccvm.CVM_PREPAID_RENEW_FLAG),
				Description:  "Auto renewal flag. Valid values: `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically. Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. NOTE: it only works when instance_charge_type is set to `PREPAID`.",
			},
			"spot_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"one-time"}),
				Description:  "Type of spot instance, only support `one-time` now. Note: it only works when instance_charge_type is set to `SPOTPAID`.",
			},
			"spot_max_price": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateStringNumber,
				Description:  "Max price of a spot instance, is the format of decimal string, for example \"0.50\". Note: it only works when instance_charge_type is set to `SPOTPAID`.",
			},
			"internet_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      INTERNET_CHARGE_TYPE_TRAFFIC_POSTPAID_BY_HOUR,
				ValidateFunc: tccommon.ValidateAllowedStringValue(INTERNET_CHARGE_ALLOW_TYPE),
				Description:  "Charge types for network traffic. Valid values: `BANDWIDTH_PREPAID`, `TRAFFIC_POSTPAID_BY_HOUR` and `BANDWIDTH_PACKAGE`.",
			},
			"internet_max_bandwidth_out": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Max bandwidth of Internet access in Mbps. Default is `0`.",
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
				ValidateFunc:  tccommon.ValidateAsConfigPassword,
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
				Description: "To specify whether to enable cloud security service. Default is `TRUE`.",
			},
			"enhanced_monitor_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "To specify whether to enable cloud monitor service. Default is `TRUE`.",
			},
			"enhanced_automation_tools_service": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "To specify whether to enable cloud automation tools service.",
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
			"disk_type_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      SCALING_DISK_TYPE_POLICY_ORIGINAL,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SCALING_DISK_TYPE_ALLOW_POLICY),
				Description:  "Policy of cloud disk type. Valid values: `ORIGINAL` and `AUTOMATIC`. Default is `ORIGINAL`.",
			},
			"cam_role_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "CAM role name authorized to access.",
			},
			"host_name_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Related settings of the cloud server hostname (HostName).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The host name of the cloud server; dots (.) and dashes (-) cannot be used as the first and last characters of HostName, and cannot be used consecutively; Windows instances are not supported; other types (Linux, etc.) instances: the character length is [2, 40], it is allowed to support multiple dots, and there is a paragraph between the dots, and each paragraph is allowed to consist of letters (no uppercase and lowercase restrictions), numbers and dashes (-). Pure numbers are not allowed.",
						},
						"host_name_style": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The style of the host name of the cloud server, the value range includes `ORIGINAL` and `UNIQUE`, the default is `ORIGINAL`; `ORIGINAL`, the AS directly passes the HostName filled in the input parameter to the CVM, and the CVM may append a sequence to the HostName number, the HostName of the instance in the scaling group will conflict; `UNIQUE`, the HostName filled in as a parameter is equivalent to the host name prefix, AS and CVM will expand it, and the HostName of the instance in the scaling group can be guaranteed to be unique.",
						},
					},
				},
			},
			"instance_name_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Settings of CVM instance names.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "CVM instance name.",
						},
						"instance_name_style": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(INSTANCE_NAME_STYLE),
							Default:      INSTANCE_NAME_ORIGINAL,
							Description:  "Type of CVM instance name. Valid values: `ORIGINAL` and `UNIQUE`. Default is `ORIGINAL`.",
						},
					},
				},
			},

			"disaster_recover_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Placement group ID. Only one is allowed.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dedicated_cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Dedicated Cluster ID.",
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
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_config.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := as.NewCreateLaunchConfigurationRequest()

	v := d.Get("configuration_name")
	request.LaunchConfigurationName = helper.String(v.(string))

	if v, ok := d.GetOk("image_id"); ok {
		request.ImageId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("image_family"); ok {
		request.ImageFamily = helper.String(v.(string))
	}

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
			deleteWithInstance := value["delete_with_instance"].(bool)
			dataDisk := as.DataDisk{
				DiskType:           &diskType,
				DiskSize:           &diskSize,
				DeleteWithInstance: &deleteWithInstance,
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
	if v, ok := d.GetOkExists("enhanced_automation_tools_service"); ok {
		automationToolsService := v.(bool)
		request.EnhancedService.AutomationToolsService = &as.RunAutomationServiceEnabled{
			Enabled: &automationToolsService,
		}
	}

	if v, ok := d.GetOk("user_data"); ok {
		request.UserData = helper.String(v.(string))
	}

	chargeType, ok := d.Get("instance_charge_type").(string)
	if !ok || chargeType == "" {
		chargeType = INSTANCE_CHARGE_TYPE_POSTPAID
	}

	if chargeType == INSTANCE_CHARGE_TYPE_SPOTPAID {
		spotMaxPrice := d.Get("spot_max_price").(string)
		spotInstanceType := d.Get("spot_instance_type").(string)
		request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
			MarketType: helper.String("spot"),
			SpotOptions: &as.SpotMarketOptions{
				MaxPrice:         &spotMaxPrice,
				SpotInstanceType: &spotInstanceType,
			},
		}
	}

	if chargeType == INSTANCE_CHARGE_TYPE_PREPAID {
		period := d.Get("instance_charge_type_prepaid_period").(int)
		renewFlag := d.Get("instance_charge_type_prepaid_renew_flag").(string)
		request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
			Period:    helper.IntInt64(period),
			RenewFlag: &renewFlag,
		}
	}

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

	if v, ok := d.GetOk("disk_type_policy"); ok {
		request.DiskTypePolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cam_role_name"); ok {
		request.CamRoleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host_name_settings"); ok {
		settings := make([]*as.HostNameSettings, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			settingsInfo := as.HostNameSettings{}
			if hostName, ok := dMap["host_name"]; ok {
				settingsInfo.HostName = helper.String(hostName.(string))
			}
			if hostNameStyle, ok := dMap["host_name_style"]; ok {
				settingsInfo.HostNameStyle = helper.String(hostNameStyle.(string))
			}
			settings = append(settings, &settingsInfo)
		}
		request.HostNameSettings = settings[0]
	}

	if v, ok := d.GetOk("instance_name_settings"); ok {
		settings := make([]*as.InstanceNameSettings, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			settingsInfo := as.InstanceNameSettings{}
			if instanceName, ok := dMap["instance_name"]; ok {
				settingsInfo.InstanceName = helper.String(instanceName.(string))
			}
			if instanceNameStyle, ok := dMap["instance_name_style"]; ok {
				settingsInfo.InstanceNameStyle = helper.String(instanceNameStyle.(string))
			}
			settings = append(settings, &settingsInfo)
		}
		request.InstanceNameSettings = settings[0]
	}

	if v, ok := d.GetOk("disaster_recover_group_ids"); ok {
		disasterRecoverGroupIds := v.([]interface{})
		request.DisasterRecoverGroupIds = make([]*string, 0, len(disasterRecoverGroupIds))
		for i := range disasterRecoverGroupIds {
			subnetId := disasterRecoverGroupIds[i].(string)
			request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &subnetId)
		}
	}

	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		request.DedicatedClusterId = helper.String(v.(string))
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().CreateLaunchConfiguration(request)
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
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	configurationId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		config, has, e := asService.DescribeLaunchConfigurationById(ctx, configurationId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if has == 0 {
			d.SetId("")
			return nil
		}
		_ = d.Set("configuration_name", *config.LaunchConfigurationName)
		_ = d.Set("status", *config.LaunchConfigurationStatus)

		if config.ImageId != nil {
			_ = d.Set("image_id", *config.ImageId)
		}
		if config.ImageFamily != nil {
			_ = d.Set("image_family", *config.ImageFamily)
		}

		_ = d.Set("project_id", *config.ProjectId)
		_ = d.Set("instance_types", helper.StringsInterfaces(config.InstanceTypes))
		_ = d.Set("system_disk_size", *config.SystemDisk.DiskSize)
		_ = d.Set("data_disk", flattenDataDiskMappings(config.DataDisks))
		_ = d.Set("internet_charge_type", *config.InternetAccessible.InternetChargeType)
		_ = d.Set("internet_max_bandwidth_out", *config.InternetAccessible.InternetMaxBandwidthOut)
		_ = d.Set("public_ip_assigned", *config.InternetAccessible.PublicIpAssigned)
		_ = d.Set("key_ids", helper.StringsInterfaces(config.LoginSettings.KeyIds))
		_ = d.Set("security_group_ids", helper.StringsInterfaces(config.SecurityGroupIds))
		_ = d.Set("enhanced_security_service", *config.EnhancedService.SecurityService.Enabled)
		_ = d.Set("enhanced_monitor_service", *config.EnhancedService.MonitorService.Enabled)
		if config.EnhancedService.AutomationToolsService.Enabled != nil {
			_ = d.Set("enhanced_automation_tools_service", *config.EnhancedService.AutomationToolsService.Enabled)
		}
		_ = d.Set("user_data", helper.PString(config.UserData))
		_ = d.Set("instance_tags", flattenInstanceTagsMapping(config.InstanceTags))
		_ = d.Set("disk_type_policy", *config.DiskTypePolicy)

		_ = d.Set("cam_role_name", *config.CamRoleName)

		if config.HostNameSettings != nil {
			isEmptySettings := true
			settings := map[string]interface{}{}
			if config.HostNameSettings.HostName != nil {
				isEmptySettings = false
				settings["host_name"] = config.HostNameSettings.HostName
			}
			if config.HostNameSettings.HostNameStyle != nil {
				isEmptySettings = false
				settings["host_name_style"] = config.HostNameSettings.HostNameStyle
			}
			if !isEmptySettings {
				_ = d.Set("host_name_settings", []interface{}{settings})
			}
		}

		if config.InstanceNameSettings != nil {
			settings := make([]map[string]interface{}, 0)
			setting := map[string]interface{}{
				"instance_name":       config.InstanceNameSettings.InstanceName,
				"instance_name_style": config.InstanceNameSettings.InstanceNameStyle,
			}
			name, nameOk := setting["instance_name"].(string)
			style, styleOk := setting["instance_name_style"].(string)
			if nameOk && name != "" || styleOk && style != "" {
				settings = append(settings, setting)
				_ = d.Set("instance_name_settings", settings)
			}
		}

		if config.SystemDisk.DiskType != nil {
			_ = d.Set("system_disk_type", *config.SystemDisk.DiskType)
		}

		if _, ok := d.GetOk("instance_charge_type"); ok || *config.InstanceChargeType != INSTANCE_CHARGE_TYPE_POSTPAID {
			_ = d.Set("instance_charge_type", *config.InstanceChargeType)
		}

		if config.InstanceMarketOptions != nil && config.InstanceMarketOptions.SpotOptions != nil {
			_ = d.Set("spot_instance_type", config.InstanceMarketOptions.SpotOptions.SpotInstanceType)
			_ = d.Set("spot_max_price", config.InstanceMarketOptions.SpotOptions.MaxPrice)
		}

		if config.InstanceChargePrepaid != nil {
			_ = d.Set("instance_charge_type_prepaid_renew_flag", config.InstanceChargePrepaid.RenewFlag)
		}

		if len(config.DisasterRecoverGroupIds) > 0 {
			_ = d.Set("disaster_recover_group_ids", helper.StringsInterfaces(config.DisasterRecoverGroupIds))
		} else {
			_ = d.Set("disaster_recover_group_ids", []string{})
		}

		if config.DedicatedClusterId != nil {
			_ = d.Set("dedicated_cluster_id", config.DedicatedClusterId)
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudAsScalingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_config.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := as.NewModifyLaunchConfigurationAttributesRequest()

	configurationId := d.Id()
	request.LaunchConfigurationId = &configurationId

	if d.HasChange("configuration_name") {
		if v, ok := d.GetOk("configuration_name"); ok {
			request.LaunchConfigurationName = helper.String(v.(string))
		}
	}

	if d.HasChange("image_id") {
		if v, ok := d.GetOk("image_id"); ok {
			request.ImageId = helper.String(v.(string))
		}
	}

	if d.HasChange("project_id") {
		return fmt.Errorf("`project_id` do not support change now.")
	}

	if d.HasChange("instance_types") {
		if v, ok := d.GetOk("instance_types"); ok {
			instanceTypes := v.([]interface{})
			request.InstanceTypes = make([]*string, 0, len(instanceTypes))
			for i := range instanceTypes {
				instanceType := instanceTypes[i].(string)
				request.InstanceTypes = append(request.InstanceTypes, &instanceType)
			}
		}
	}

	if d.HasChange("system_disk_type") || d.HasChange("system_disk_size") {
		request.SystemDisk = &as.SystemDisk{}
		if v, ok := d.GetOk("system_disk_type"); ok {
			request.SystemDisk.DiskType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("system_disk_size"); ok {
			request.SystemDisk.DiskSize = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("data_disk") {
		if v, ok := d.GetOk("data_disk"); ok {
			dataDisks := v.([]interface{})
			request.DataDisks = make([]*as.DataDisk, 0, len(dataDisks))
			for _, d := range dataDisks {
				value := d.(map[string]interface{})
				diskType := value["disk_type"].(string)
				diskSize := uint64(value["disk_size"].(int))
				snapshotId := value["snapshot_id"].(string)
				deleteWithInstance := value["delete_with_instance"].(bool)
				dataDisk := as.DataDisk{
					DiskType:           &diskType,
					DiskSize:           &diskSize,
					DeleteWithInstance: &deleteWithInstance,
				}
				if snapshotId != "" {
					dataDisk.SnapshotId = &snapshotId
				}
				request.DataDisks = append(request.DataDisks, &dataDisk)
			}
		}
	}

	if d.HasChange("internet_charge_type") || d.HasChange("internet_max_bandwidth_out") || d.HasChange("public_ip_assigned") {
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
	}

	if d.HasChange("security_group_ids") {
		if v, ok := d.GetOk("security_group_ids"); ok {
			securityGroups := v.([]interface{})
			request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
			for i := range securityGroups {
				securityGroup := securityGroups[i].(string)
				request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroup)
			}
		}
	}

	if d.HasChange("enhanced_security_service") || d.HasChange("enhanced_monitor_service") || d.HasChange("enhanced_automation_tools_service") {
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
		if v, ok := d.GetOkExists("enhanced_automation_tools_service"); ok {
			automationToolsService := v.(bool)
			request.EnhancedService.AutomationToolsService = &as.RunAutomationServiceEnabled{
				Enabled: &automationToolsService,
			}
		}
	}

	if d.HasChange("user_data") {
		if v, ok := d.GetOk("user_data"); ok {
			request.UserData = helper.String(v.(string))
		}
	}

	if d.HasChange("instance_charge_type") {
		chargeType, ok := d.Get("instance_charge_type").(string)
		if !ok || chargeType == "" {
			chargeType = INSTANCE_CHARGE_TYPE_POSTPAID
		}

		if chargeType == INSTANCE_CHARGE_TYPE_SPOTPAID {
			spotMaxPrice := d.Get("spot_max_price").(string)
			spotInstanceType := d.Get("spot_instance_type").(string)
			request.InstanceMarketOptions = &as.InstanceMarketOptionsRequest{
				MarketType: helper.String("spot"),
				SpotOptions: &as.SpotMarketOptions{
					MaxPrice:         &spotMaxPrice,
					SpotInstanceType: &spotInstanceType,
				},
			}
		}

		if chargeType == INSTANCE_CHARGE_TYPE_PREPAID {
			period := d.Get("instance_charge_type_prepaid_period").(int)
			renewFlag := d.Get("instance_charge_type_prepaid_renew_flag").(string)
			request.InstanceChargePrepaid = &as.InstanceChargePrepaid{
				Period:    helper.IntInt64(period),
				RenewFlag: &renewFlag,
			}
		}

		request.InstanceChargeType = &chargeType
	}

	if d.HasChange("instance_types_check_policy") {
		if v, ok := d.GetOk("instance_types_check_policy"); ok {
			request.InstanceTypesCheckPolicy = helper.String(v.(string))
		}
	}

	if d.HasChange("instance_tags") {
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
	}

	if d.HasChange("disk_type_policy") {
		if v, ok := d.GetOk("disk_type_policy"); ok {
			request.DiskTypePolicy = helper.String(v.(string))
		}
	}

	if d.HasChange("cam_role_name") {
		if v, ok := d.GetOk("cam_role_name"); ok {
			request.CamRoleName = helper.String(v.(string))
		}
	}

	if d.HasChange("host_name_settings") {
		if v, ok := d.GetOk("host_name_settings"); ok {
			settings := make([]*as.HostNameSettings, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				settingsInfo := as.HostNameSettings{}
				if hostName, ok := dMap["host_name"]; ok {
					settingsInfo.HostName = helper.String(hostName.(string))
				}
				if hostNameStyle, ok := dMap["host_name_style"]; ok {
					settingsInfo.HostNameStyle = helper.String(hostNameStyle.(string))
				}
				settings = append(settings, &settingsInfo)
			}
			request.HostNameSettings = settings[0]
		}
	}

	if d.HasChange("instance_name_settings") {
		if v, ok := d.GetOk("instance_name_settings"); ok {
			settings := make([]*as.InstanceNameSettings, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				settingsInfo := as.InstanceNameSettings{}
				if instanceName, ok := dMap["instance_name"]; ok {
					settingsInfo.InstanceName = helper.String(instanceName.(string))
				}
				if instanceNameStyle, ok := dMap["instance_name_style"]; ok {
					settingsInfo.InstanceNameStyle = helper.String(instanceNameStyle.(string))
				}
				settings = append(settings, &settingsInfo)
			}
			request.InstanceNameSettings = settings[0]
		}
	}

	if d.HasChange("image_family") {
		if v, ok := d.GetOk("image_family"); ok {
			request.ImageFamily = helper.String(v.(string))
		}
	}

	if d.HasChange("password") || d.HasChange("key_ids") || d.HasChange("keep_image_login") {
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
	}

	if d.HasChange("disaster_recover_group_ids") {
		if v, ok := d.GetOk("disaster_recover_group_ids"); ok {
			disasterRecoverGroupIds := v.([]interface{})
			request.DisasterRecoverGroupIds = make([]*string, 0, len(disasterRecoverGroupIds))
			for i := range disasterRecoverGroupIds {
				subnetId := disasterRecoverGroupIds[i].(string)
				request.DisasterRecoverGroupIds = append(request.DisasterRecoverGroupIds, &subnetId)
			}
		}
	}

	if d.HasChange("dedicated_cluster_id") {
		if v, ok := d.GetOk("dedicated_cluster_id"); ok {
			request.DedicatedClusterId = helper.String(v.(string))
		}
	}

	response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseAsClient().ModifyLaunchConfigurationAttributes(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}

	return resourceTencentCloudAsScalingConfigRead(d, meta)
}

func resourceTencentCloudAsScalingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_as_scaling_config.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	configurationId := d.Id()
	asService := AsService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := asService.DeleteLaunchConfiguration(ctx, configurationId)
	if err != nil {
		return err
	}

	return nil
}
