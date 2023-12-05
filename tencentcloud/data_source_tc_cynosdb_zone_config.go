package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTencentCynosdbZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCynosdbZoneConfigRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"list": {
				Type:        schema.TypeList,
				Description: "A list of zone. Each element contains the following attributes:",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance CPU, unit: core.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance memory, unit: GB.",
						},
						"max_storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum available storage for the instance, unit GB.",
						},
						"min_storage_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum available storage of the instance, unit: GB.",
						},
						"machine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Machine type.",
						},
						"max_io_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max io bandwidth.",
						},
						"zone_stock_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Regional inventory information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"zone": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Availability zone.",
									},
									"has_stock": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Has stock.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCynosdbZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_zone_config.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region

	instanceSpecSet, err := service.DescribeRedisZoneConfig(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeRedisZoneConfig]fail, return %s", err.Error())
	}

	result := make([]map[string]interface{}, 0)

	for _, instanceSpec := range instanceSpecSet {
		resultItem := make(map[string]interface{})
		resultItem["cpu"] = *instanceSpec.Cpu
		resultItem["memory"] = *instanceSpec.Memory
		resultItem["max_storage_size"] = *instanceSpec.MaxStorageSize
		resultItem["min_storage_size"] = *instanceSpec.MinStorageSize
		resultItem["machine_type"] = *instanceSpec.MachineType
		resultItem["max_io_bandwidth"] = *instanceSpec.MaxIoBandWidth
		zoneStockInfos := make([]map[string]interface{}, 0)
		for _, zoneStockInfoItem := range instanceSpec.ZoneStockInfos {
			zoneStockInfo := make(map[string]interface{})
			zoneStockInfo["zone"] = *zoneStockInfoItem.Zone
			zoneStockInfo["has_stock"] = *zoneStockInfoItem.HasStock

			zoneStockInfos = append(zoneStockInfos, zoneStockInfo)
		}
		resultItem["zone_stock_infos"] = zoneStockInfos
		result = append(result, resultItem)
	}

	id := "cynosdb_zoneconfig_" + region
	d.SetId(id)
	_ = d.Set("list", result)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {

		if err := writeToFile(output.(string), result); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}

	}
	return nil
}
