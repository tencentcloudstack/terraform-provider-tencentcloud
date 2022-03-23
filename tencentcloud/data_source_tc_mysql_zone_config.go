/*
Use this data source to query the available database specifications for different regions. And a maximum of 20 requests can be initiated per second for this query.

Example Usage

```hcl
data "tencentcloud_mysql_zone_config" "mysql" {
  region             = "ap-guangzhou"
  result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func TencentMysqlSellType() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cdb_type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "",
		},
		"mem_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Memory size (in MB).",
		},
		"min_volume_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Minimum disk size (in GB).",
		},
		"max_volume_size": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Maximum disk size (in GB).",
		},
		"volume_step": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Disk increment (in GB).",
		},
		"qps": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Queries per second.",
		},
	}
}

func TencentMysqlZoneConfig() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The name of available zone which is equal to a specific datacenter.",
		},
		"is_default": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Indicates whether the current DC is the default DC for the region. Possible returned values: `0` - no; `1` - yes.",
		},
		"is_support_disaster_recovery": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Indicates whether recovery is supported: `0` - No; `1` - Yes.",
		},
		"is_support_vpc": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Indicates whether VPC is supported: `0` - No; `1` - Yes.",
		},
		"engine_versions": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "The version number of the database engine to use. Supported versions include `5.5`/`5.6`/`5.7`.",
		},
		"pay_type": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Computed:    true,
			Description: "",
		},
		"hour_instance_sale_max_num": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "",
		},
		"support_slave_sync_modes": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Computed:    true,
			Description: "Data replication mode. `0` - Async replication; `1` - Semisync replication; `2` - Strongsync replication.",
		},
		"disaster_recovery_zones": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "Information about available zones of recovery.",
		},
		"slave_deploy_modes": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Computed:    true,
			Description: "Availability zone deployment method. Available values: `0` - Single availability zone; `1` - Multiple availability zones.",
		},
		"first_slave_zones": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "Zone information about first slave instance.",
		},
		"second_slave_zones": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "Zone information about second slave instance.",
		},
		"remote_ro_zones": {
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Description: "Zone information about remote ro instance.",
		},
		"sells": {Type: schema.TypeList,
			Computed:    true,
			Description: "A list of supported instance types for sell:",
			Elem: &schema.Resource{
				Schema: TencentMysqlSellType(),
			},
		},
	}
}

func dataSourceTencentMysqlZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentMysqlZoneConfigRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region parameter, which is used to identify the region to which the data you want to work with belongs.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of zone config. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: TencentMysqlZoneConfig(),
				},
			},
		},
	}
}

func dataSourceTencentMysqlZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_zone_config.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region
	if regionInterface, ok := d.GetOk("region"); ok {
		region = regionInterface.(string)
	} else {
		log.Printf("[INFO]%s region is not set,so we use [%s] from env\n ", logId, region)
	}

	sellConfigures, err := mysqlService.DescribeDBZoneConfig(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeBackups]fail, return %s", err.Error())
	}
	var regionItem *cdb.RegionSellConf
	for _, regionItem = range sellConfigures {
		if *regionItem.Region == region {
			break
		}
	}
	if regionItem == nil {
		return nil
	}
	var zoneConfigs []interface{}
	for _, sellItem := range regionItem.ZonesConf {
		if *sellItem.Status != ZONE_SELL_STATUS_ONLINE && *sellItem.Status != ZONE_SELL_STATUS_NEW {
			continue
		}
		var zoneConfig = make(map[string]interface{})
		zoneConfig["name"] = *sellItem.Zone
		if sellItem.HourInstanceSaleMaxNum != nil {
			zoneConfig["hour_instance_sale_max_num"] = *sellItem.HourInstanceSaleMaxNum
		}

		if sellItem.IsDefaultZone != nil {
			if *sellItem.IsDefaultZone {
				zoneConfig["is_default"] = 1
			} else {
				zoneConfig["is_default"] = 0
			}
		}

		if sellItem.IsSupportDr != nil {
			if *sellItem.IsSupportDr {
				zoneConfig["is_support_disaster_recovery"] = 1
			} else {
				zoneConfig["is_support_disaster_recovery"] = 0
			}
		}

		if sellItem.IsSupportVpc != nil {
			if *sellItem.IsSupportVpc {
				zoneConfig["is_support_vpc"] = 1
			} else {
				zoneConfig["is_support_vpc"] = 0
			}
		}

		payTypes := make([]int, len(sellItem.PayType))
		for index, strPtr := range sellItem.PayType {
			if tempInt, err := strconv.ParseInt(*strPtr, 10, 64); err != nil {
				errRet := fmt.Errorf("api[DescribeDBZoneConfig]return PayType error,not int")
				log.Printf("[CRITAL]%s %s\n ", logId, errRet.Error())
				return errRet
			} else {
				payTypes[index] = int(tempInt)
			}
		}
		zoneConfig["pay_type"] = payTypes

		supportSlaveSyncModes := make([]string, len(sellItem.ProtectMode))
		for index, intPtr := range sellItem.ProtectMode {
			supportSlaveSyncModes[index] = *intPtr
		}
		zoneConfig["support_slave_sync_modes"] = payTypes

		disasterRecoveryZones := make([]string, len(sellItem.DrZone))
		for index, strPtr := range sellItem.DrZone {
			disasterRecoveryZones[index] = *strPtr
		}
		zoneConfig["disaster_recovery_zones"] = disasterRecoveryZones

		var (
			slaveDeployModes                                                 []int
			firstSlaveZones, secondSlaveZones, engineVersions, remoteRoZones []string
			sells                                                            []interface{}
		)
		if sellItem.ZoneConf != nil {
			for _, mode := range sellItem.ZoneConf.DeployMode {
				slaveDeployModes = append(slaveDeployModes, int(*mode))
			}
			for _, zoneName := range sellItem.ZoneConf.SlaveZone {
				firstSlaveZones = append(firstSlaveZones, *zoneName)
			}
			for _, zoneName := range sellItem.ZoneConf.BackupZone {
				secondSlaveZones = append(secondSlaveZones, *zoneName)
			}
			for _, zoneName := range sellItem.RemoteRoZone {
				remoteRoZones = append(remoteRoZones, *zoneName)
			}
		}
		zoneConfig["slave_deploy_modes"] = slaveDeployModes
		zoneConfig["first_slave_zones"] = firstSlaveZones
		zoneConfig["second_slave_zones"] = secondSlaveZones
		zoneConfig["remote_ro_zones"] = remoteRoZones

		for _, mysqlConfigs := range sellItem.SellType {
			for _, strPtr := range mysqlConfigs.EngineVersion {
				engineVersions = append(engineVersions, *strPtr)
			}
			for _, mysqlConfig := range mysqlConfigs.Configs {
				var showConfigMap = make(map[string]interface{})
				showConfigMap["cdb_type"] = *mysqlConfig.CdbType
				showConfigMap["mem_size"] = int(*mysqlConfig.Memory)
				showConfigMap["max_volume_size"] = int(*mysqlConfig.VolumeMax)
				showConfigMap["min_volume_size"] = int(*mysqlConfig.VolumeMin)
				showConfigMap["volume_step"] = int(*mysqlConfig.VolumeStep)
				showConfigMap["qps"] = int(*mysqlConfig.Qps)
				sells = append(sells, showConfigMap)
			}
		}
		zoneConfig["engine_versions"] = engineVersions
		zoneConfig["sells"] = sells

		zoneConfigs = append(zoneConfigs, zoneConfig)
	}

	if err := d.Set("list", zoneConfigs); err != nil {
		log.Printf("[CRITAL]%s provider set zoneConfigs fail, reason:%s\n ", logId, err.Error())
	}
	d.SetId("zoneconfig" + region)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), zoneConfigs); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}

	}
	return nil
}
