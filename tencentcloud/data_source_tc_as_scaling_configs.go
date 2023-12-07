package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudAsScalingConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsScalingConfigRead,

		Schema: map[string]*schema.Schema{
			"configuration_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Launch configuration ID.",
			},
			"configuration_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Launch configuration name.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"configuration_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of configuration. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Launch configuration ID.",
						},
						"configuration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Launch configuration name.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of available image, for example `img-8toqc6s3`.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project to which the configuration belongs. Default value is 0.",
						},
						"instance_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Instance type list of the scaling configuration.",
						},
						"system_disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "System disk category of the scaling configuration.",
						},
						"system_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "System disk size of the scaling configuration in GB.",
						},
						"data_disk": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Configurations of data disk.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of disk.",
									},
									"disk_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Volume of disk in GB. Default is `0`.",
									},
									"snapshot_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Data disk snapshot ID.",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicates whether the disk remove after instance terminated.",
									},
								},
							},
						},
						"internet_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge types for network traffic.",
						},
						"internet_max_bandwidth_out": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max bandwidth of Internet access in Mbps.",
						},
						"public_ip_assigned": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specify whether to assign an Internet IP address.",
						},
						"key_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "ID list of login keys.",
						},
						"security_group_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Security groups to which the instance belongs.",
						},
						"enhanced_security_service": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to activate cloud security service.",
						},
						"enhanced_monitor_service": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to activate cloud monitor service.",
						},
						"user_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Base64-encoded User Data text.",
						},
						"instance_tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "A tag list associates with an instance.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current status of a launch configuration.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the launch configuration was created.",
						},
						"disk_type_policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy of cloud disk type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAsScalingConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_as_scaling_configs.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	asService := AsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	configurationId := ""
	configurationName := ""
	if v, ok := d.GetOk("configuration_id"); ok {
		configurationId = v.(string)
	}
	if v, ok := d.GetOk("configuration_name"); ok {
		configurationName = v.(string)
	}

	configs, err := asService.DescribeLaunchConfigurationByFilter(ctx, configurationId, configurationName)
	if err != nil {
		return err
	}

	configurationList := make([]map[string]interface{}, 0, len(configs))
	for _, config := range configs {
		mapping := map[string]interface{}{
			"configuration_id":           *config.LaunchConfigurationId,
			"configuration_name":         *config.LaunchConfigurationName,
			"image_id":                   *config.ImageId,
			"project_id":                 *config.ProjectId,
			"instance_types":             helper.StringsInterfaces(config.InstanceTypes),
			"system_disk_size":           *config.SystemDisk.DiskSize,
			"data_disk":                  flattenDataDiskMappings(config.DataDisks),
			"internet_charge_type":       *config.InternetAccessible.InternetChargeType,
			"internet_max_bandwidth_out": *config.InternetAccessible.InternetMaxBandwidthOut,
			"public_ip_assigned":         *config.InternetAccessible.PublicIpAssigned,
			"key_ids":                    helper.StringsInterfaces(config.LoginSettings.KeyIds),
			"security_group_ids":         helper.StringsInterfaces(config.SecurityGroupIds),
			"enhanced_security_service":  *config.EnhancedService.SecurityService.Enabled,
			"enhanced_monitor_service":   *config.EnhancedService.MonitorService.Enabled,
			"user_data":                  helper.PString(config.UserData),
			"instance_tags":              flattenInstanceTagsMapping(config.InstanceTags),
			"status":                     *config.LaunchConfigurationStatus,
			"create_time":                *config.CreatedTime,
			"disk_type_policy":           *config.DiskTypePolicy,
		}
		if config.SystemDisk.DiskType != nil {
			mapping["system_disk_type"] = *config.SystemDisk.DiskType
		}
		configurationList = append(configurationList, mapping)
	}

	d.SetId("ConfigurationList" + configurationId + configurationName)
	err = d.Set("configuration_list", configurationList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set configuration list fail, reason:%s\n ", logId, err.Error())
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err = writeToFile(output.(string), configurationList); err != nil {
			return err
		}
	}

	return nil
}
