package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudMongodbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vport": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mongo_version": {
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
						"volume": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"machine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shard_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudMongodbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	instanceId := ""
	instanceType := -1
	namePrefix := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("instance_type"); ok {
		instanceType = v.(int)
	}
	if v, ok := d.GetOk("instance_name_prefix"); ok {
		namePrefix = v.(string)
	}
	mongodbService := MongodbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	mongodbs, err := mongodbService.DescribeInstancesByFilter(ctx, instanceId, instanceType)
	if err != nil {
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(mongodbs))
	ids := make([]string, 0, len(mongodbs))
	for _, mongo := range mongodbs {
		if namePrefix != "" && strings.HasPrefix(*mongo.InstanceName, namePrefix) {
			continue
		}

		switch *mongo.MachineType {
		case "HIO10G":
			*mongo.MachineType = "TGIO"

		case "HIO":
			*mongo.MachineType = "GIO"
		}

		instance := map[string]interface{}{
			"instance_id":    mongo.InstanceId,
			"instance_name":  mongo.InstanceName,
			"project_id":     mongo.ProjectId,
			"cluster_type":   mongo.ClusterType,
			"zone":           mongo.Zone,
			"vpc_id":         mongo.VpcId,
			"subnet_id":      mongo.SubnetId,
			"status":         mongo.Status,
			"vip":            mongo.Vip,
			"vport":          mongo.Vport,
			"create_time":    mongo.CreateTime,
			"mongo_version":  mongo.MongoVersion,
			"cpu":            mongo.CpuNum,
			"memory":         *mongo.Memory / 1024,
			"volume":         *mongo.Volume / 1024,
			"machine_type":   mongo.MachineType,
			"shard_quantity": mongo.ReplicationSetNum,
		}
		instanceList = append(instanceList, instance)
		ids = append(ids, *mongo.InstanceId)
	}

	d.SetId(dataResourceIdsHash(ids))
	if err = d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set mongodb instance list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), instanceList); err != nil {
			return err
		}
	}

	return nil
}
