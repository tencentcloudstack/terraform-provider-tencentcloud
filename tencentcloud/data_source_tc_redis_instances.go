/*
Use this data source to query the detail information of redis instance.

Example Usage

```hcl
data "tencentcloud_redis_instances" "redislab" {
    zone                = "ap-hongkong-1"
    search_key          = "myredis"
    project_id          = 0
    limit               = 20
    result_output_file  = "/tmp/redis_instances"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentRedisInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentRedisInstancesRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of an available zone.",
			},
			"search_key": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Key words used to match the results, and the key words can be: instance ID, instance name and IP address.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Description: "ID of the project to which  redis instance belongs.",
			},
			"limit": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Optional:    true,
				Description: "The number limitation of results for a query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values
			"instance_list": {Type: schema.TypeList,
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
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type. Available values: master_slave_redis, master_slave_ckv, cluster_ckv, cluster_redis and standalone_redis.",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size in MB",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current status of an instance，maybe: init, processing, online, isolate and todelete.",
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
					},
				},
			},
		},
	}
}

func dataSourceTencentRedisInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_redis_instances.read")()

	logId := getLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := RedisService{client: meta.(*TencentCloudClient).apiV3Conn}
	region := meta.(*TencentCloudClient).apiV3Conn.Region

	var (
		zone      string = ""
		searchKey string = ""
		projectId int64  = -1
		limit     int64  = -1
	)

	if temp, ok := d.GetOk("zone"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			zone = tempStr
			if !strings.Contains(zone, region) {
				return fmt.Errorf("zone[%s] not in region[%s]", zone, region)
			}
		}
	}
	if temp, ok := d.GetOk("search_key"); ok {
		tempStr := temp.(string)
		if tempStr != "" {
			searchKey = tempStr
		}
	}

	if temp, ok := d.GetOk("project_id"); ok {
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

	instances, err := service.DescribeInstances(ctx, zone, searchKey, projectId, limit)
	if err != nil {
		return err
	}

	var instanceList = make([]interface{}, 0, len(instances))

	for _, instance := range instances {

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

		instanceList = append(instanceList, instanceDes)
	}

	if err := d.Set("instance_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set  redis instances fail, reason:%s\n ", logId, err.Error())
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
