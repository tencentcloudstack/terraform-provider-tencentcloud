/*
Use this data source to query Publish Subscribe resources for the specific SQL Server instance.

Example Usage

```hcl
resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = tencentcloud_sqlserver_instance.publish_instance.id
	subscribe_instance_id           = tencentcloud_sqlserver_instance.subscribe_instance.id
	publish_subscribe_name          = "example"
	database_tuples {
		publish_database            = tencentcloud_sqlserver_db.test_publish_subscribe.name
		subscribe_database          = tencentcloud_sqlserver_db.test_publish_subscribe.name
	}
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
				Description: "The subscribe/publish instance ID is related to whether the `instance_id` is a publish instance or a subscribe instance. when `instance_id` is a publish instance, this field is filtered according to the subscribe instance ID; when `instance_id` is a subscribe instance, this field is filtering according to the publish instance ID.",
			},
			"pub_or_sub_instance_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The intranet IP of the subscribe/publish instance is related to whether the `instance_id` is a publish instance or a subscribe instance. when `instance_id` is a publish instance, this field is filtered according to the intranet IP of the subscribe instance; when `instance_id` is a subscribe instance, this field is based on the publish instance intranet IP filter.",
			},
			"publish_subscribe_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The id of the Publish and Subscribe in the SQLServer instance.",
			},
			"publish_subscribe_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Publish and Subscribe in the SQLServer instance.",
			},
			"publish_database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Publish the database.",
			},
			"subscribe_database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Subscribe to the database.",
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
							Description: "The id of the Publish and Subscribe in the SQLServer instance.",
						},
						"publish_subscribe_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the Publish and Subscribe in the SQLServer instance.",
						},
						"publish_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Publish the instance ID in the SQLServer instance.",
						},
						"publish_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Publish the instance name in the SQLServer instance.",
						},
						"publish_instance_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Publish the instance IP in the SQLServer instance.",
						},
						"subscribe_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subscribe the instance ID in the SQLServer instance.",
						},
						"subscribe_instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subscribe the instance name in the SQLServer instance.",
						},
						"subscribe_instance_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subscribe the instance IP in the SQLServer instance.",
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
										Description: "Publish the database.",
									},
									"subscribe_database": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Subscribe to the database.",
									},
									"last_sync_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last sync time.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Publish and subscribe status between databases `running`, `success`, `fail`, `unknow`.",
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
	var (
		instanceId           = d.Get("instance_id").(string)
		pubOrSubInstanceId   string
		pubOrSubInstanceIp   string
		publishSubscribeName string
		publishSubscribeId   uint64 = 0
		publishDBName        string
		subscribeDBName      string
	)

	if v, ok := d.GetOk("pub_or_sub_instance_id"); ok {
		pubOrSubInstanceId = v.(string)
	}
	if v, ok := d.GetOk("pub_or_sub_instance_ip"); ok {
		pubOrSubInstanceIp = v.(string)
	}
	if v, ok := d.GetOk("publish_subscribe_name"); ok {
		publishSubscribeName = v.(string)
	}
	if v, ok := d.GetOk("publish_subscribe_id"); ok {
		id := v.(int)
		if id != 0 {
			publishSubscribeId = *helper.IntUint64(id)
		}
	}
	if v, ok := d.GetOk("publish_database"); ok {
		publishDBName = v.(string)
	}
	if v, ok := d.GetOk("subscribe_database"); ok {
		subscribeDBName = v.(string)
	}

	publishSubscribes, err := sqlserverService.DescribeSqlserverPublishSubscribes(ctx, instanceId, pubOrSubInstanceId, pubOrSubInstanceIp, publishSubscribeName, publishDBName, subscribeDBName, publishSubscribeId)
	if err != nil {
		return err
	}
	instanceList := make([]map[string]interface{}, 0, len(publishSubscribes))
	ids := make([]string, 0, len(publishSubscribes))

	for _, publishSubscribe := range publishSubscribes {
		var databaseTupleStatus []map[string]interface{}
		for _, inst_ := range publishSubscribe.DatabaseTupleSet {
			databaseTuple := map[string]interface{}{
				"publish_database":   inst_.PublishDatabase,
				"subscribe_database": inst_.SubscribeDatabase,
				"last_sync_time":     inst_.LastSyncTime,
				"status":             inst_.Status,
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
		log.Printf("[CRITAL]%s provider set sql server publish subscribe list fail, reason:%s\n ", logId, err.Error())
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
