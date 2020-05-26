/*
Use this data source to query the available mongodb specifications for different zone.

Example Usage

```hcl
data "tencentcloud_mongodb_zone_config" "mongodb" {
  available_zone = "ap-guangzhou-2"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudMongodbZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbZoneConfigRead,

		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The available zone of the Mongodb.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of zone config. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"available_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone of the Mongodb.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of Mongodb cluster.",
						},
						"machine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of Mongodb instance.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of cpu's core.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size.",
						},
						"default_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Default disk size.",
						},
						"min_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum sie of the disk.",
						},
						"max_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum size of the disk.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the Mongodb version.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudMongodbZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mongodb_zone_config.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	zone := ""
	if v, ok := d.GetOk("available_zone"); ok {
		zone = v.(string)
	}
	mongodbService := MongodbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	infos, err := mongodbService.DescribeSpecInfo(ctx, zone)
	if err != nil {
		return err
	}

	configList := make([]map[string]interface{}, 0)
	for _, info := range infos {
		for _, item := range info.SpecItems {
			clusterType := MONGODB_CLUSTER_TYPE_REPLSET
			if *item.ClusterType == 1 {
				clusterType = MONGODB_CLUSTER_TYPE_SHARD
			}
			mapping := map[string]interface{}{
				"available_zone":  info.Zone,
				"cluster_type":    clusterType,
				"machine_type":    item.MachineType,
				"cpu":             item.Cpu,
				"memory":          item.Memory,
				"default_storage": item.DefaultStorage,
				"min_storage":     item.MinStorage,
				"max_storage":     item.MaxStorage,
				"engine_version":  item.Version,
			}
			configList = append(configList, mapping)
		}
	}

	id := zone
	if id == "" {
		id = meta.(*TencentCloudClient).apiV3Conn.Region
	}
	d.SetId(id)
	if err = d.Set("list", configList); err != nil {
		log.Printf("[CRITAL]%s provider set mongodb zone config list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), configList); err != nil {
			return err
		}
	}

	return nil
}
