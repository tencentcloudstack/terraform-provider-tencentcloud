/*
Use this data source to query detailed information of Mongodb instances.

Example Usage

```hcl
data "tencentcloud_mongodb_instances" "mongodb" {
  instance_id  = "cmgo-l6lwdsel"
  cluster_type = "REPLSET"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudMongodbInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMongodbInstancesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the Mongodb instance to be queried.",
			},
			"instance_name_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name prefix of the Mongodb instance.",
			},
			"cluster_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(MONGODB_CLUSTER_TYPE),
				Description:  "Type of Mongodb cluster, and available values include replica set cluster(expressed with `REPLSET`), sharding cluster(expressed with `SHARD`).",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the Mongodb instance to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},

			// computed
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of instances. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the Mongodb instance.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Mongodb instance.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project which the instance belongs.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of Mongodb cluster.",
						},
						"available_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone of the Mongodb.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the VPC.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the subnet.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the Mongodb, and available values include pending initialization(expressed with 0),  processing(expressed with 1), running(expressed with 2) and expired(expressed with -2).",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP of the Mongodb instance.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP port of the Mongodb instance.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the Mongodb instance.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the Mongodb engine.",
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
						"volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size.",
						},
						"machine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of Mongodb instance.",
						},
						"shard_quantity": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of sharding.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of the Mongodb instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudMongodbInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mongodb_instances.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	instanceId := ""
	clusterType := -1
	name := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	if v, ok := d.GetOk("cluster_type"); ok {
		vv := v.(string)
		if vv == MONGODB_CLUSTER_TYPE_REPLSET {
			clusterType = 0
		} else if vv == MONGODB_CLUSTER_TYPE_SHARD {
			clusterType = 1
		}
	}
	if v, ok := d.GetOk("instance_name_prefix"); ok {
		name = v.(string)
	}

	tags := getTags(d, "tags")

	mongodbService := MongodbService{client: meta.(*TencentCloudClient).apiV3Conn}
	mongodbs, err := mongodbService.DescribeInstancesByFilter(ctx, instanceId, clusterType)
	if err != nil {
		return err
	}

	instanceList := make([]map[string]interface{}, 0, len(mongodbs))
	ids := make([]string, 0, len(mongodbs))

instancesLoop:
	for _, mongo := range mongodbs {
		if nilFields := CheckNil(mongo, map[string]string{
			"InstanceId":        "instance id",
			"InstanceName":      "instance name",
			"ProjectId":         "project id",
			"ClusterType":       "cluster type",
			"Zone":              "available zone",
			"VpcId":             "vpc id",
			"SubnetId":          "subnet id",
			"Status":            "status",
			"Vip":               "vip",
			"Vport":             "vport",
			"CreateTime":        "create time",
			"MongoVersion":      "engine version",
			"CpuNum":            "cpu number",
			"Memory":            "memory",
			"Volume":            "volume",
			"MachineType":       "machine type",
			"ReplicationSetNum": "shard quantity",
		}); len(nilFields) > 0 {
			return fmt.Errorf("mongodb %v are nil", nilFields)
		}

		if !strings.Contains(*mongo.InstanceName, name) {
			continue
		}

		respTags := make(map[string]string, len(mongo.Tags))
		for _, t := range mongo.Tags {
			if t.TagKey == nil {
				return errors.New("mongodb tag key is nil")
			}
			if t.TagValue == nil {
				return errors.New("mongodb tag value is nil")
			}

			respTags[*t.TagKey] = *t.TagValue
		}

		for k, v := range tags {
			if value, ok := respTags[k]; !ok || v != value {
				continue instancesLoop
			}
		}

		switch *mongo.MachineType {
		case "HIO10G":
			*mongo.MachineType = MONGODB_MACHINE_TYPE_TGIO

		case "HIO":
			*mongo.MachineType = MONGODB_MACHINE_TYPE_GIO
		}

		clusterType := MONGODB_CLUSTER_TYPE_REPLSET
		if *mongo.ClusterType == 1 {
			clusterType = MONGODB_CLUSTER_TYPE_SHARD
		}

		instance := map[string]interface{}{
			"instance_id":    mongo.InstanceId,
			"instance_name":  mongo.InstanceName,
			"project_id":     mongo.ProjectId,
			"cluster_type":   clusterType,
			"available_zone": mongo.Zone,
			"vpc_id":         mongo.VpcId,
			"subnet_id":      mongo.SubnetId,
			"status":         mongo.Status,
			"vip":            mongo.Vip,
			"vport":          mongo.Vport,
			"create_time":    mongo.CreateTime,
			"engine_version": mongo.MongoVersion,
			"cpu":            mongo.CpuNum,
			"memory":         *mongo.Memory / 1024,
			"volume":         *mongo.Volume / 1024,
			"machine_type":   mongo.MachineType,
			"shard_quantity": mongo.ReplicationSetNum,
			"tags":           respTags,
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
