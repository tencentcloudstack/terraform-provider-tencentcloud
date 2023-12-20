package cynosdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbRollBackCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbRollBackClusterCreate,
		Read:   resourceTencentCloudCynosdbRollBackClusterRead,
		Delete: resourceTencentCloudCynosdbRollBackClusterDelete,
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

func resourceTencentCloudCynosdbRollBackClusterCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_roll_back_cluster.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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
			dMap := item.(map[string]interface{})
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
			dMap := item.(map[string]interface{})
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().RollBackCluster(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cynosdb rollBackCluster failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(clusterId)

	flowId := *response.Response.FlowId
	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err = resource.Retry(6*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ok, err := service.DescribeFlow(ctx, flowId)
		if err != nil {
			if _, ok := err.(*sdkErrors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if ok {
			return nil
		} else {
			return resource.RetryableError(fmt.Errorf("operate cynosdb rollBackCluster is processing"))
		}
	})

	if err != nil {
		log.Printf("[CRITAL]%s Open cynosdb wan fail, reason:%s\n", logId, err.Error())
		return err
	}

	return resourceTencentCloudCynosdbRollBackClusterRead(d, meta)
}

func resourceTencentCloudCynosdbRollBackClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_roll_back_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCynosdbRollBackClusterDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_roll_back_cluster.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
