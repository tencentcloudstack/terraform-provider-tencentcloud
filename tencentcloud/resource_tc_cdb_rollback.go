/*
Provides a resource to create a cdb rollback

Example Usage

```hcl
resource "tencentcloud_cdb_rollback" "rollback" {
  instances {
		instance_id = "cdb_xxx"
		strategy = ""
		rollback_time = ""
		databases {
			database_name = ""
			new_database_name = ""
		}
		tables {
			database = ""
			table {
				table_name = ""
				new_table_name = ""
			}
		}

  }
}
```

Import

cdb rollback can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_rollback.rollback rollback_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCdbRollback() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbRollbackCreate,
		Read:   resourceTencentCloudCdbRollbackRead,
		Delete: resourceTencentCloudCdbRollbackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instances": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Instance details for rollback.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cloud database instance ID.",
						},
						"strategy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rollback strategy. Available values are: table, db, full; the default value is full. table - Extremely fast rollback mode, only import the backup and binlog of the selected table level, if there is a cross-table operation, and the associated table is not selected at the same time, the rollback will fail. In this mode, the parameter Databases must be empty; db - Quick mode, only import the backup and binlog of the selected library level, if there is a cross-database operation, and the associated library is not selected at the same time, the rollback will fail; full - normal rollback mode, the backup and binlog of the entire instance will be imported , at a slower rate.",
						},
						"rollback_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Database rollback time, the time format is: yyyy-mm-dd hh:mm:ss.",
						},
						"databases": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The database information to be archived, indicating that the entire database is archived.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The original database name before rollback.",
									},
									"new_database_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The new database name after rollback.",
									},
								},
							},
						},
						"tables": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The database table information to be rolled back, indicating that the file is rolled back by table.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"database": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Database name.",
									},
									"table": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Database table details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"table_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The original database table name before rollback.",
												},
												"new_table_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "New database table name after rollback.",
												},
											},
										},
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

func resourceTencentCloudCdbRollbackCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_rollback.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewStartBatchRollbackRequest()
		response   = cdb.NewStartBatchRollbackResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instances"); ok {
		for _, item := range v.([]interface{}) {
			rollbackInstancesInfo := cdb.RollbackInstancesInfo{}
			if v, ok := dMap["instance_id"]; ok {
				rollbackInstancesInfo.InstanceId = helper.String(v.(string))
			}
			if v, ok := dMap["strategy"]; ok {
				rollbackInstancesInfo.Strategy = helper.String(v.(string))
			}
			if v, ok := dMap["rollback_time"]; ok {
				rollbackInstancesInfo.RollbackTime = helper.String(v.(string))
			}
			if v, ok := dMap["databases"]; ok {
				for _, item := range v.([]interface{}) {
					databasesMap := item.(map[string]interface{})
					rollbackDBName := cdb.RollbackDBName{}
					if v, ok := databasesMap["database_name"]; ok {
						rollbackDBName.DatabaseName = helper.String(v.(string))
					}
					if v, ok := databasesMap["new_database_name"]; ok {
						rollbackDBName.NewDatabaseName = helper.String(v.(string))
					}
					rollbackInstancesInfo.Databases = append(rollbackInstancesInfo.Databases, &rollbackDBName)
				}
			}
			if v, ok := dMap["tables"]; ok {
				for _, item := range v.([]interface{}) {
					tablesMap := item.(map[string]interface{})
					rollbackTables := cdb.RollbackTables{}
					if v, ok := tablesMap["database"]; ok {
						rollbackTables.Database = helper.String(v.(string))
					}
					if v, ok := tablesMap["table"]; ok {
						for _, item := range v.([]interface{}) {
							tableMap := item.(map[string]interface{})
							rollbackTableName := cdb.RollbackTableName{}
							if v, ok := tableMap["table_name"]; ok {
								rollbackTableName.TableName = helper.String(v.(string))
							}
							if v, ok := tableMap["new_table_name"]; ok {
								rollbackTableName.NewTableName = helper.String(v.(string))
							}
							rollbackTables.Table = append(rollbackTables.Table, &rollbackTableName)
						}
					}
					rollbackInstancesInfo.Tables = append(rollbackInstancesInfo.Tables, &rollbackTables)
				}
			}
			request.Instances = append(request.Instances, &rollbackInstancesInfo)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().StartBatchRollback(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cdb rollback failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCEED"}, 1*readRetryTimeout, time.Second, service.CdbRollbackStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCdbRollbackRead(d, meta)
}

func resourceTencentCloudCdbRollbackRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_rollback.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCdbRollbackDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_rollback.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
