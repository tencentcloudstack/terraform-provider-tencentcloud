package tencentcloud

import (
	"context"

	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentRedisInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentRedisInstancesRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"search_key": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
			},
			"result_output_file": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			// Computed values
			"instance_list": {Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"redis_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mem_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
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
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentRedisInstancesRead(d *schema.ResourceData, meta interface{}) error {

	defer LogElapsed("data_source.tencentcloud_redis_instances.read")()

	logId := GetLogId(nil)
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
