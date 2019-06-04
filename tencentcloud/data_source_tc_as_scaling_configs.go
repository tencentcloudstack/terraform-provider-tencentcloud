package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudAsScalingConfigs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAsScalingConfigRead,

		Schema: map[string]*schema.Schema{
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configuration_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configuration_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"configuration_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"configuration_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"system_disk_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"system_disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_disk": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disk_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"disk_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"snapshot_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_max_bandwidth_out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"public_ip_assigned": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"key_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enhanced_security_service": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enhanced_monitor_service": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"user_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAsScalingConfigRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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
			"instance_types":             flattenStringList(config.InstanceTypes),
			"system_disk_type":           *config.SystemDisk.DiskType,
			"system_disk_size":           *config.SystemDisk.DiskSize,
			"data_disk":                  flattenDataDiskMappings(config.DataDisks),
			"internet_charge_type":       *config.InternetAccessible.InternetChargeType,
			"internet_max_bandwidth_out": *config.InternetAccessible.InternetMaxBandwidthOut,
			"public_ip_assigned":         *config.InternetAccessible.PublicIpAssigned,
			"key_ids":                    flattenStringList(config.LoginSettings.KeyIds),
			"security_group_ids":         flattenStringList(config.SecurityGroupIds),
			"enhanced_security_service":  *config.EnhancedService.SecurityService.Enabled,
			"enhanced_monitor_service":   *config.EnhancedService.MonitorService.Enabled,
			"user_data":                  pointerToString(config.UserData),
			"instance_tags":              flattenInstanceTagsMapping(config.InstanceTags),
			"status":                     *config.LaunchConfigurationStatus,
			"create_time":                *config.CreatedTime,
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
		writeToFile(output.(string), configurationList)
	}

	return nil
}
