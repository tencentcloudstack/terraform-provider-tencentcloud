package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseDiskConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseDiskConfigRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter list.zoneFilter by availability zone.Type: StringRequired: no.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter value of field.",
						},
					},
				},
			},

			"disk_config_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of cloud disk configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"disk_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud disk type.",
						},
						"disk_sales_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud disk sale status.",
						},
						"max_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum cloud disk size.",
						},
						"min_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum cloud disk size.",
						},
						"disk_step_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cloud disk increment.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudLighthouseDiskConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_disk_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*lighthouse.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := lighthouse.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var diskConfigSet []*lighthouse.DiskConfig

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseDiskConfigByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		diskConfigSet = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(diskConfigSet))

	if diskConfigSet != nil {
		for _, diskConfig := range diskConfigSet {
			diskConfigMap := map[string]interface{}{}

			if diskConfig.Zone != nil {
				diskConfigMap["zone"] = diskConfig.Zone
			}

			if diskConfig.DiskType != nil {
				diskConfigMap["disk_type"] = diskConfig.DiskType
			}

			if diskConfig.DiskSalesState != nil {
				diskConfigMap["disk_sales_state"] = diskConfig.DiskSalesState
			}

			if diskConfig.MaxDiskSize != nil {
				diskConfigMap["max_disk_size"] = diskConfig.MaxDiskSize
			}

			if diskConfig.MinDiskSize != nil {
				diskConfigMap["min_disk_size"] = diskConfig.MinDiskSize
			}

			if diskConfig.DiskStepSize != nil {
				diskConfigMap["disk_step_size"] = diskConfig.DiskStepSize
			}

			tmpList = append(tmpList, diskConfigMap)
		}

		_ = d.Set("disk_config_set", tmpList)
	}
	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
