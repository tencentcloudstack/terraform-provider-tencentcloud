package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
)

func dataSourceTencentCloudMariadbDatabaseObjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDatabaseObjectsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"db_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "database name.",
			},

			"tables": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "table list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "table name.",
						},
					},
				},
			},

			"views": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "view list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"view": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "view name.",
						},
					},
				},
			},

			"procs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "proc list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "proc name.",
						},
					},
				},
			},

			"funcs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "func list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"func": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "func name.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMariadbDatabaseObjectsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_database_objects.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := ""
	dbName := ""

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("db_name"); ok {
		dbName = v.(string)
	}

	var databaseObjects *mariadb.DescribeDatabaseObjectsResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbDatabaseObjectsByFilter(ctx, instanceId, dbName)
		if e != nil {
			return retryError(e)
		}
		databaseObjects = result
		return nil
	})
	if err != nil {
		return err
	}

	_ = d.Set("instance_id", instanceId)
	_ = d.Set("db_name", dbName)

	if databaseObjects.Tables != nil {
		tmpList := make([]map[string]interface{}, 0, len(databaseObjects.Tables))
		for _, databaseTable := range databaseObjects.Tables {
			databaseTableMap := map[string]interface{}{}

			if databaseTable.Table != nil {
				databaseTableMap["table"] = databaseTable.Table
			}
			tmpList = append(tmpList, databaseTableMap)
		}

		_ = d.Set("tables", tmpList)
	}

	if databaseObjects.Views != nil {
		tmpList := make([]map[string]interface{}, 0, len(databaseObjects.Views))
		for _, databaseView := range databaseObjects.Views {
			databaseViewMap := map[string]interface{}{}

			if databaseView.View != nil {
				databaseViewMap["view"] = databaseView.View
			}

			tmpList = append(tmpList, databaseViewMap)
		}

		_ = d.Set("views", tmpList)
	}

	if databaseObjects.Procs != nil {
		tmpList := make([]map[string]interface{}, 0, len(databaseObjects.Procs))
		for _, databaseProcedure := range databaseObjects.Procs {
			databaseProcedureMap := map[string]interface{}{}

			if databaseProcedure.Proc != nil {
				databaseProcedureMap["proc"] = databaseProcedure.Proc
			}

			tmpList = append(tmpList, databaseProcedureMap)
		}

		_ = d.Set("procs", tmpList)
	}

	if databaseObjects.Funcs != nil {
		tmpList := make([]map[string]interface{}, 0, len(databaseObjects.Funcs))
		for _, databaseFunction := range databaseObjects.Funcs {
			databaseFunctionMap := map[string]interface{}{}

			if databaseFunction.Func != nil {
				databaseFunctionMap["func"] = databaseFunction.Func
			}

			tmpList = append(tmpList, databaseFunctionMap)
		}

		_ = d.Set("funcs", tmpList)
	}

	d.SetId(instanceId + FILED_SP + dbName)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
