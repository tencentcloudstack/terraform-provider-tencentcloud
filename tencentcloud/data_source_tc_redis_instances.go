/*
Use this data source to query the detail information of redis instance.

Example Usage

```hcl
data "tencentcloud_redis_instances" "redislab" {
  zone               = "ap-hongkong-1"
  search_key         = "myredis"
  project_id         = 0
  limit              = 20
  result_output_file = "/tmp/redis_instances"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentRedisInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentRedisInstancesRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of an available zone.",
			},
			"search_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key words used to match the results, and the key words can be: instance ID, instance name and IP address.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the project to which redis instance belongs.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number limitation of results for a query.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of redis instance.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of redis instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redis_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of a redis instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of a redis instance.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available zone to which a redis instance belongs.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of the project to which a redis instance belongs.",
						},
						"type_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance type. Refer to `data.tencentcloud_redis_zone_config.list.type_id` get available values.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Deprecated:  "It has been deprecated from version 1.33.1. Please use 'type_id' instead.",
							Description: "Instance type. Available values: `master_slave_redis`, `master_slave_ckv`, `cluster_ckv`, `cluster_redis` and `standalone_redis`.",
						},
						"redis_shard_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of instance shard.",
						},
						"redis_replicas_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of instance copies.",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size in MB.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current status of an instance, maybe: `init`, `processing`, `online`, `isolate` and `todelete`.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the vpc with which the instance is associated.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the vpc subnet.",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP address of an instance.",
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port used to access a redis instance.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the instance is created.",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Tags of an instance.",
						},
						// payment
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charge type of instance. Valid values are `POSTPAID` and `PREPAID`.",
						},
						"node_info": {
							Type:		 schema.TypeList,
							Computed: 	 true,
							Description: "List of instance node information. Currently, information about the node type (master or replica) and node availability zone can be passed in.",
							Elem:        &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master": {
										Type: schema.TypeBool,
										Computed: true,
										Description: "Indicates whether the node is master.",
									},
									"id": {
										Type:	schema.TypeInt,
										Computed: true,
										Description: "ID of the master or replica node.",
									},
									"zone_id": {
										Type: schema.TypeInt,
										Computed: true,
										Description: "ID of the availability zone of the master or replica node.",
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

func dataSourceTencentRedisInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_redis_instances.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	client := meta.(*TencentCloudClient).apiV3Conn
	service := RedisService{client: client}
	region := client.Region

	var (
		zone      string
		searchKey string
		projectId int64 = -1
		limit     int64 = -1
	)

	if temp, ok := d.GetOk("zone"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			zone = tempStr
		}
	}
	if temp, ok := d.GetOk("search_key"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			searchKey = tempStr
		}
	}

	if temp, ok := d.GetOkExists("project_id"); ok {
		tempInt := temp.(int)
		if tempInt >= 0 {
			projectId = int64(tempInt)
		}
	}

	if temp, ok := d.GetOk("limit"); ok {
		tempInt := temp.(int)
		if tempInt >= 0 {
			limit = int64(tempInt)
		}
	}

	tags := helper.GetTags(d, "tags")

	instances, err := service.DescribeInstances(ctx, zone, searchKey, projectId, limit)
	if err != nil {
		return err
	}

	var instanceList = make([]interface{}, 0, len(instances))

instanceLoop:
	for _, instance := range instances {
		if len(tags) > 0 {
			// filter by tags, must match all tags
			for k, v := range tags {
				if instance.Tags[k] != v {
					continue instanceLoop
				}
			}
		}

		var instanceDes = make(map[string]interface{})
		instanceDes["redis_id"] = instance.RedisId
		instanceDes["name"] = instance.Name
		instanceDes["zone"] = instance.Zone
		instanceDes["project_id"] = instance.ProjectId
		instanceDes["type"] = instance.Type
		instanceDes["mem_size"] = instance.MemSize
		instanceDes["status"] = instance.Status
		instanceDes["vpc_id"] = instance.VpcId
		instanceDes["subnet_id"] = instance.SubnetId
		instanceDes["ip"] = instance.Ip
		instanceDes["port"] = instance.Port
		instanceDes["create_time"] = instance.CreateTime
		instanceDes["tags"] = instance.Tags
		instanceDes["redis_shard_num"] = instance.RedisShardNum
		instanceDes["redis_replicas_num"] = instance.RedisReplicasNum
		instanceDes["type_id"] = instance.TypeId
		instanceDes["charge_type"] = instance.BillingMode
		instanceDes["node_info"] = instance.NodeInfo
		instanceList = append(instanceList, instanceDes)
	}

	if err := d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set redis instances fail, reason:%s\n", logId, err.Error())
	}
	d.SetId("redis_instances_list" + region)

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), instanceList); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
		}
	}
	return nil
}
