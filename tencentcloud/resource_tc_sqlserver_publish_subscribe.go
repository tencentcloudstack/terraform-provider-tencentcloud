/*
Provides a SQL Server PublishSubscribe resource belongs to SQL Server instance.

Example Usage

```hcl
resource "tencentcloud_sqlserver_publish_subscribe" "example" {
	publish_instance_id             = tencentcloud_sqlserver_instance.publish_instance.id
	subscribe_instance_id           = tencentcloud_sqlserver_instance.subscribe_instance.id
	publish_subscribe_name          = "example"
	delete_subscribe_db             = false
	database_tuples {
		publish_database            = tencentcloud_sqlserver_db.test_publish_subscribe.name
		subscribe_database          = tencentcloud_sqlserver_db.test_publish_subscribe.name
	}
}
```

Import

SQL Server PublishSubscribe can be imported using the publish_sqlserver_id#subscribe_sqlserver_id, e.g.

```
$ terraform import tencentcloud_sqlserver_publish_subscribe.foo publish_sqlserver_id#subscribe_sqlserver_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
)

func resourceTencentCloudSqlserverPublishSubscribe() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverPublishSubscribeCreate,
		Read:   resourceTencentCloudSqlserverPublishSubscribeRead,
		Update: resourceTencentCloudSqlserverPublishSubscribeUpdate,
		Delete: resourceTencentCloudSqlserverPublishSubscribeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"publish_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Publish the instance ID in the SQLServer instance.",
			},
			"subscribe_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Subscribe the instance ID in the SQLServer instance.",
			},
			"publish_subscribe_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default_name",
				Description: "The name of the Publish and Subscribe in the SQLServer instance. default is `default_name`.",
			},
			"delete_subscribe_db": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to delete the subscriber database when deleting the Publish and Subscribe in the SQLServer instance. `true` for deletes the subscribe database, `false` for does not delete the subscribe database. default is `false`.",
			},
			"database_tuples": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    80,
				Description: "Database Publish and Publish relationship list. The elements inside can be deleted and added individually, but modify is not allowed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"publish_database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Publish the database.",
						},
						"subscribe_database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subscribe to the database.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudSqlserverPublishSubscribeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_publish_subscribe.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		sqlserverService     = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		publishInstanceId    = d.Get("publish_instance_id").(string)
		subscribeInstanceId  = d.Get("subscribe_instance_id").(string)
		publishSubscribeName = d.Get("publish_subscribe_name").(string)
		databaseTuples       = d.Get("database_tuples").(*schema.Set).List()
	)

	//check publishInstanceId exist and status
	_, has, err := sqlserverService.DescribeSqlserverInstanceById(ctx, publishInstanceId)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s DescribeSqlserverInstanceById fail, reason:%s\n", logId, err)
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s SQL Server Publish instance %s dose not exist for Publish Subscribe Create", logId, publishInstanceId)
	}

	//check subscribeInstanceId exist and status
	_, has, err = sqlserverService.DescribeSqlserverInstanceById(ctx, subscribeInstanceId)
	if err != nil {
		return fmt.Errorf("[CRITAL]%s DescribeSqlserverInstanceById fail, reason:%s\n", logId, err)
	}
	if !has {
		return fmt.Errorf("[CRITAL]%s SQL Server Subscribe %s dose not exist for Publish Subscribe Create", logId, subscribeInstanceId)
	}

	if err := sqlserverService.CreateSqlserverPublishSubscribe(ctx, publishInstanceId, subscribeInstanceId, publishSubscribeName, databaseTuples); err != nil {
		return err
	}
	_, hasExist, err := sqlserverService.DescribeSqlserverPublishSubscribeById(ctx, publishInstanceId, subscribeInstanceId)
	if err != nil {
		return err
	}
	if !hasExist {
		return fmt.Errorf("this Sqlserver Publish Subscribe %s Create Failed", subscribeInstanceId)
	}
	d.SetId(publishInstanceId + FILED_SP + subscribeInstanceId)
	return resourceTencentCloudSqlserverPublishSubscribeRead(d, meta)
}

func resourceTencentCloudSqlserverPublishSubscribeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_publish_subscribe.read")()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		sqlserverService = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		id               = d.Id()
		idItem           = strings.Split(id, FILED_SP)
	)

	if len(idItem) < 2 {
		return fmt.Errorf("broken ID %s of SQL Server Publish Subscribe", id)
	}
	publishInstanceId := idItem[0]
	subscribeInstanceId := idItem[1]
	publishSubscribe, has, err := sqlserverService.DescribeSqlserverPublishSubscribeById(ctx, publishInstanceId, subscribeInstanceId)
	if err != nil {
		return err
	}
	if !has {
		d.SetId("")
		return nil
	}
	_ = d.Set("publish_subscribe_id", publishSubscribe.Id)
	_ = d.Set("publish_subscribe_name", publishSubscribe.Name)
	_ = d.Set("publish_instance_id", publishSubscribe.PublishInstanceId)
	_ = d.Set("publish_instance_name", publishSubscribe.PublishInstanceName)
	_ = d.Set("publish_instance_ip", publishSubscribe.PublishInstanceIp)
	_ = d.Set("subscribe_instance_id", publishSubscribe.SubscribeInstanceId)
	_ = d.Set("subscribe_instance_name", publishSubscribe.SubscribeInstanceName)
	_ = d.Set("subscribe_instance_ip", publishSubscribe.SubscribeInstanceIp)
	var databaseTupleSet []map[string]interface{}
	for _, inst_ := range publishSubscribe.DatabaseTupleSet {
		databaseTuple := map[string]interface{}{
			"publish_database":   inst_.PublishDatabase,
			"subscribe_database": inst_.SubscribeDatabase,
		}
		databaseTupleSet = append(databaseTupleSet, databaseTuple)
	}
	_ = d.Set("database_tuples", databaseTupleSet)
	return nil
}

func resourceTencentCloudSqlserverPublishSubscribeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_publish_subscribe.update")()

	var (
		logId               = getLogId(contextNil)
		ctx                 = context.WithValue(context.TODO(), logIdKey, logId)
		sqlserverService    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		publishInstanceId   = d.Get("publish_instance_id").(string)
		subscribeInstanceId = d.Get("subscribe_instance_id").(string)
		deleteSubscribeDb   = d.Get("delete_subscribe_db").(bool)
	)

	publishSubscribe, _, err := sqlserverService.DescribeSqlserverPublishSubscribeById(ctx, publishInstanceId, subscribeInstanceId)
	if err != nil {
		return err
	}
	if d.HasChange("publish_subscribe_name") {
		publishSubscribeName := d.Get("publish_subscribe_name").(string)
		if err := sqlserverService.ModifyPublishSubscribeName(ctx, *publishSubscribe.Id, publishSubscribeName); err != nil {
			return err
		}
	}

	if d.HasChange("database_tuples") {
		var deleteDatabaseTupleSet []*sqlserver.DatabaseTuple
		var subscribeDatabases []*string
		oldSet, newSet := d.GetChange("database_tuples")
		//get new DatabaseTupleSet
		var newDatabaseTupleSet []*sqlserver.DatabaseTuple
		for _, inst_ := range newSet.(*schema.Set).List() {
			inst := inst_.(map[string]interface{})
			newDatabaseTupleSet = append(newDatabaseTupleSet, sqlServerNewDatabaseTuple(inst["publish_database"], inst["subscribe_database"]))
		}
		//get old DatabaseTupleSet
		var oldDatabaseTupleSet []*sqlserver.DatabaseTuple
		for _, inst_ := range oldSet.(*schema.Set).List() {
			inst := inst_.(map[string]interface{})
			oldDatabaseTupleSet = append(oldDatabaseTupleSet, sqlServerNewDatabaseTuple(inst["publish_database"], inst["subscribe_database"]))
		}

		for _, oldInstance := range oldDatabaseTupleSet {
			var exist = false
			for _, newInstance := range newDatabaseTupleSet {
				if *newInstance.SubscribeDatabase == *oldInstance.SubscribeDatabase && *newInstance.PublishDatabase == *oldInstance.PublishDatabase {
					exist = true
					break
				}
			}
			if !exist {
				databaseTuple := sqlserver.DatabaseTuple{
					PublishDatabase:   oldInstance.PublishDatabase,
					SubscribeDatabase: oldInstance.SubscribeDatabase,
				}
				deleteDatabaseTupleSet = append(deleteDatabaseTupleSet, &databaseTuple)
				subDatabase := *oldInstance.SubscribeDatabase
				subscribeDatabases = append(subscribeDatabases, &subDatabase)
			}
		}
		if deleteDatabaseTupleSet == nil {
			return fmt.Errorf("[CRITAL]%s resourceTencentCloudSqlserverPublishSubscribeUpdate fail, reason: DatabaseTupleSet does not allow modify", logId)
		}
		if err := sqlserverService.DeletePublishSubscribe(ctx, publishSubscribe, deleteDatabaseTupleSet); err != nil {
			return err
		}
		if deleteSubscribeDb {
			//delete subscribe databases
			if err = sqlserverService.DeleteSqlserverDB(ctx, subscribeInstanceId, subscribeDatabases); err != nil {
				return err
			}
		}
	}

	return resourceTencentCloudSqlserverPublishSubscribeRead(d, meta)
}

func resourceTencentCloudSqlserverPublishSubscribeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_publish_subscribe.delete")()

	var (
		logId               = getLogId(contextNil)
		ctx                 = context.WithValue(context.TODO(), logIdKey, logId)
		sqlserverService    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		publishInstanceId   = d.Get("publish_instance_id").(string)
		subscribeInstanceId = d.Get("subscribe_instance_id").(string)
		deleteSubscribeDb   = d.Get("delete_subscribe_db").(bool)
	)
	publishSubscribe, has, err := sqlserverService.DescribeSqlserverPublishSubscribeById(ctx, publishInstanceId, subscribeInstanceId)
	if err != nil {
		return err
	}
	if !has {
		return nil
	}

	oldDatabaseTuples, _ := d.GetChange("database_tuples")
	var oldDatabaseTupleSet []*sqlserver.DatabaseTuple
	var subscribeDatabases []*string
	for _, inst_ := range oldDatabaseTuples.(*schema.Set).List() {
		inst := inst_.(map[string]interface{})
		subDatabase := inst["subscribe_database"].(string)
		oldDatabaseTupleSet = append(oldDatabaseTupleSet, sqlServerNewDatabaseTuple(inst["publish_database"], inst["subscribe_database"]))
		subscribeDatabases = append(subscribeDatabases, &subDatabase)
	}

	if err := sqlserverService.DeletePublishSubscribe(ctx, publishSubscribe, oldDatabaseTupleSet); err != nil {
		return err
	}
	if deleteSubscribeDb {
		//delete subscribe databases
		if err = sqlserverService.DeleteSqlserverDB(ctx, subscribeInstanceId, subscribeDatabases); err != nil {
			return err
		}
	}
	return nil
}
