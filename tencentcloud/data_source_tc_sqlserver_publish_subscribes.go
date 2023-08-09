/*
Use this data source to query Publish Subscribe resources for the specific SQL Server instance.

Example Usage

```hcl
data "tencentcloud_sqlserver_publish_subscribes" "example" {
  instance_id = tencentcloud_sqlserver_publish_subscribe.example.publish_instance_id
}

data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-example"
  description = "desc."
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example_pub" {
  name                 = "tf-example-pub"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  db_version           = "2008R2"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example_sub" {
  name                 = "tf-example-sub"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = tencentcloud_subnet.subnet.id
  vpc_id               = tencentcloud_vpc.vpc.id
  db_version           = "2008R2"
  security_group_list  = [tencentcloud_security_group.security_group.id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}

resource "tencentcloud_sqlserver_db" "example_pub" {
  instance_id = tencentcloud_sqlserver_general_cloud_instance.example_pub.id
  name        = "tf_example_db_pub"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_db" "example_sub" {
  instance_id = tencentcloud_sqlserver_general_cloud_instance.example_sub.id
  name        = "tf_example_db_sub"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

resource "tencentcloud_sqlserver_publish_subscribe" "example" {
  publish_instance_id    = tencentcloud_sqlserver_general_cloud_instance.example_pub.id
  subscribe_instance_id  = tencentcloud_sqlserver_general_cloud_instance.example_sub.id
  publish_subscribe_name = "example"
  delete_subscribe_db    = false
  database_tuples {
    publish_database   = tencentcloud_sqlserver_db.example_pub.name
    subscribe_database = tencentcloud_sqlserver_db.example_sub.name
  }
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentSqlserverPublishSubscribes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentSqlserverPublishSubscribesRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the SQL Server instance.",
			},
			"pub_or_sub_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The subscribe/publish instance ID. It is related to whether the `instance_id` is a publish instance or a subscribe instance. when `instance_id` is a publish instance, this field is filtered according to the subscribe instance ID; when `instance_id` is a subscribe instance, this field is filtering according to the publish instance ID.",
			},
			"pub_or_sub_instance_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The intranet IP of the subscribe/publish instance. It is related to whether the `instance_id` is a publish instance or a subscribe instance. when `instance_id` is a publish instance, this field is filtered according to the intranet IP of the subscribe instance; when `instance_id` is a subscribe instance, this field is based on the publish instance intranet IP filter.",
			},
			"publish_subscribe_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The id of the Publish and Subscribe.",
			},
			"publish_subscribe_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Publish and Subscribe.",
			},
			"publish_database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of publish database.",
			},
			"subscribe_database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of subscribe database.",
			},
			"publish_subscribe_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Publish and subscribe list. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"publish_subscribe_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The id of the Publish and Subscribe.",
						},
						"publish_subscribe_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Publish and Subscribe.",
						},
						"publish_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the SQL Server instance which publish.",
						},
						"publish_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SQL Server instance which publish.",
						},
						"publish_instance_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP of the the SQL Server instance which publish.",
						},
						"subscribe_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the SQL Server instance which subscribe.",
						},
						"subscribe_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SQL Server instance which subscribe.",
						},
						"subscribe_instance_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP of the SQL Server instance which subscribe.",
						},
						"database_tuples": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Database Publish and Publish relationship list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"publish_database": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the publish SQL Server instance.",
									},
									"subscribe_database": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the subscribe SQL Server instance.",
									},
									"last_sync_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last sync time.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Publish and subscribe status between databases, valid values are `running`, `success`, `fail`, `unknow`.",
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

func dataSourceTencentSqlserverPublishSubscribesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_publish_subscribes.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
	paramMap := make(map[string]interface{})
	paramMap["instanceId"] = d.Get("instance_id").(string)
	if v, ok := d.GetOk("pub_or_sub_instance_id"); ok {
		paramMap["pubOrSubInstanceId"] = v.(string)
	}
	if v, ok := d.GetOk("pub_or_sub_instance_ip"); ok {
		paramMap["pubOrSubInstanceIp"] = v.(string)
	}
	if v, ok := d.GetOk("publish_subscribe_name"); ok {
		paramMap["publishSubscribeName"] = v.(string)
	}
	if v, ok := d.GetOk("publish_subscribe_id"); ok {
		id := v.(int)
		paramMap["publishSubscribeId"] = *helper.IntUint64(id)
	} else {
		paramMap["publishSubscribeId"] = *helper.IntUint64(0)
	}
	if v, ok := d.GetOk("publish_database"); ok {
		paramMap["publishDBName"] = v.(string)
		paramMap["subscribeDBName"] = v.(string)
	}

	publishSubscribes, err := sqlserverService.DescribeSqlserverPublishSubscribes(ctx, paramMap)
	if err != nil {
		return err
	}
	instanceList := make([]map[string]interface{}, 0, len(publishSubscribes))
	ids := make([]string, 0, len(publishSubscribes))

	for _, publishSubscribe := range publishSubscribes {
		var databaseTupleStatus []map[string]interface{}
		for _, inst := range publishSubscribe.DatabaseTupleSet {
			databaseTuple := map[string]interface{}{
				"publish_database":   inst.PublishDatabase,
				"subscribe_database": inst.SubscribeDatabase,
				"last_sync_time":     inst.LastSyncTime,
				"status":             inst.Status,
			}
			databaseTupleStatus = append(databaseTupleStatus, databaseTuple)
		}
		instance := map[string]interface{}{
			"publish_subscribe_id":    publishSubscribe.Id,
			"publish_subscribe_name":  publishSubscribe.Name,
			"publish_instance_id":     publishSubscribe.PublishInstanceId,
			"publish_instance_ip":     publishSubscribe.PublishInstanceIp,
			"publish_instance_name":   publishSubscribe.PublishInstanceName,
			"subscribe_instance_id":   publishSubscribe.SubscribeInstanceId,
			"subscribe_instance_ip":   publishSubscribe.SubscribeInstanceIp,
			"subscribe_instance_name": publishSubscribe.SubscribeInstanceName,
			"database_tuples":         databaseTupleStatus,
		}
		resourceId := *publishSubscribe.PublishInstanceId + FILED_SP + *publishSubscribe.SubscribeInstanceId
		instanceList = append(instanceList, instance)
		ids = append(ids, resourceId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if err = d.Set("publish_subscribe_list", instanceList); err != nil {
		log.Printf("[CRITAL]%s provider set sql server publish and subscribe list fail, reason:%s ", logId, err.Error())
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
