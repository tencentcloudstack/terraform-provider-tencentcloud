package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func dataSourceTencentRedisZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentRedisZoneConfigRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(connectivity.MysqlSupportedRegions),
			},
			"result_output_file": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			// Computed values
			"list": {Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mem_sizes": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentRedisZoneConfigRead(d *schema.ResourceData, meta interface{}) error {

	defer LogElapsed("data_source.tencentcloud_redis_zone_config.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region

	if regionInterface, ok := d.GetOk("region"); ok {
		region = regionInterface.(string)
	} else {
		log.Printf("[INFO]%s region is not set,so we use [%s] from env\n ", logId, region)
	}

	sellConfigures, err := service.DescribeRedisZoneConfig(ctx)
	if err != nil {
		return fmt.Errorf("api[DescribeRedisZoneConfig]fail, return %s", err.Error())
	}

	var regionItem *redis.RegionConf
	for _, regionItem = range sellConfigures {
		if *regionItem.RegionId == region {
			break
		}
	}
	if regionItem == nil {
		return nil
	}
	var allZonesConfigs []interface{}

	for _, zones := range regionItem.ZoneSet {
		zoneName := *zones.ZoneId

		for _, products := range zones.ProductSet {
			//1- package year and month, 0- billing according to quantity.
			if *products.PayMode != "0" {
				continue
			}
			//this products sale out.
			if *products.Saleout {
				continue
			}
			//not support this type now .
			if REDIS_NAMES[*products.Type] == "" {
				continue
			}

			zoneConfigures := map[string]interface{}{}
			zoneConfigures["zone"] = zoneName
			zoneConfigures["version"] = *products.Version

			memSizes := make([]int64, 0, len(products.TotalSize))

			for _, size := range products.TotalSize {
				temp, err := strconv.ParseInt(*size, 10, 64)
				if err != nil {
					continue
				}
				memSizes = append(memSizes, temp*1024)
			}

			zoneConfigures["mem_sizes"] = memSizes
			zoneConfigures["type"] = REDIS_NAMES[*products.Type]

			allZonesConfigs = append(allZonesConfigs, zoneConfigures)
		}
	}

	if err := d.Set("list", allZonesConfigs); err != nil {
		log.Printf("[CRITAL]%s provider set  redis zoneConfigs fail, reason:%s\n ", logId, err.Error())
	}
	d.SetId("redis_zoneconfig" + region)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {

		if err := writeToFile(output.(string), allZonesConfigs); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}

	}
	return nil
}
