/*
Provides a resource to create a cynosdb rollback_cluster

Example Usage

```hcl
resource "tencentcloud_cynosdb_rollback_cluster" "rollback_cluster" {
  cluster_id = "cynosdbmysql-xxxxxxxx"
  rollback_strategy = "timeRollback"
  rollback_id = 1
  expect_time = "	2022-01-20 00:00:00"
  expect_time_thresh = 1
  rollback_databases {
		old_database = "old_db_1"
		new_database = "new_db_1"

  }
  rollback_tables {
		database = "old_db_1"
		tables {
			old_table = "old_tbl_1"
			new_table = "new_tbl_1"
		}

  }
  rollback_mode = "full"
}
```

Import

cynosdb rollback_cluster can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_rollback_cluster.rollback_cluster rollback_cluster_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudCynosdbRollbackCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbRollbackClusterCreate,
		Read:   resourceTencentCloudCynosdbRollbackClusterRead,
		Delete: resourceTencentCloudCynosdbRollbackClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of cluster.",
			},

			"rollback_strategy": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Backfile policy timeRollback - Backfile by point in time snapRollback - Backfile by backup file.",
			},

			"rollback_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Rollback ID.",
			},

			"expect_time": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Expected rollback Time.",
			},

			"expect_time_thresh": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Expected Threshold (Obsolete).",
			},

			"rollback_databases": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Database list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"old_database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Old database name.",
						},
						"new_database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "New database name.",
						},
					},
				},
			},

			"rollback_tables": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				Description: "Table list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "New database name.",
						},
						"tables": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Tables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"old_table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Old table name.",
									},
									"new_table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "New table name.",
									},
								},
							},
						},
					},
				},
			},

			"rollback_mode": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Rollback mode by time point, full: normal; Db: fast; Table: Extreme speed (default is normal).",
			},
		},
	}
}

func resourceTencentCloudCynosdbRollbackClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_rollback_cluster.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cynosdb.NewRollBackClusterRequest()
		response  = cynosdb.NewRollBackClusterResponse()
		clusterId string
	)
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		request.ClusterId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("rollback_strategy"); ok {
		request.RollbackStrategy = helper.String(v.(string))
	}

	if v, _ := d.GetOk("rollback_id"); v != nil {
		request.RollbackId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("expect_time"); ok {
		request.ExpectTime = helper.String(v.(string))
	}

	if v, _ := d.GetOk("expect_time_thresh"); v != nil {
		request.ExpectTimeThresh = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("rollback_databases"); ok {
		for _, item := range v.([]interface{}) {
			rollbackDatabase := cynosdb.RollbackDatabase{}
			if v, ok := dMap["old_database"]; ok {
				rollbackDatabase.OldDatabase = helper.String(v.(string))
			}
			if v, ok := dMap["new_database"]; ok {
				rollbackDatabase.NewDatabase = helper.String(v.(string))
			}
			request.RollbackDatabases = append(request.RollbackDatabases, &rollbackDatabase)
		}
	}

	if v, ok := d.GetOk("rollback_tables"); ok {
		for _, item := range v.([]interface{}) {
			rollbackTable := cynosdb.RollbackTable{}
			if v, ok := dMap["database"]; ok {
				rollbackTable.Database = helper.String(v.(string))
			}
			if v, ok := dMap["tables"]; ok {
				for _, item := range v.([]interface{}) {
					tablesMap := item.(map[string]interface{})
					rollbackTableInfo := cynosdb.RollbackTableInfo{}
					if v, ok := tablesMap["old_table"]; ok {
						rollbackTableInfo.OldTable = helper.String(v.(string))
					}
					if v, ok := tablesMap["new_table"]; ok {
						rollbackTableInfo.NewTable = helper.String(v.(string))
					}
					rollbackTable.Tables = append(rollbackTable.Tables, &rollbackTableInfo)
				}
			}
			request.RollbackTables = append(request.RollbackTables, &rollbackTable)
		}
	}

	if v, ok := d.GetOk("rollback_mode"); ok {
		request.RollbackMode = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCynosdbClient().RollBackCluster(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb rollbackCluster failed, reason:%+v", logId, err)
		return err
	}

	clusterId = *response.Response.ClusterId
	d.SetId(clusterId)

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"success"}, 30*readRetryTimeout, time.Second, service.CynosdbRollbackClusterStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudCynosdbRollbackClusterRead(d, meta)
}

func resourceTencentCloudCynosdbRollbackClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_rollback_cluster.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbRollbackClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cynosdb_rollback_cluster.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
