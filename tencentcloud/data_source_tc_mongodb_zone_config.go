package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudMongodbZoneConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbZoneConfigRead,

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"machine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudMongodbZoneConfigRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	zone := ""
	if v, ok := d.GetOk("zone"); ok {
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
		for _, item := range info.SpecItem {
			mapping := map[string]interface{}{
				"zone":            info.Zone,
				"cluster_type":    item.ClusterType,
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
