/*
Use this data source to query detailed information of mariadb database_objects

Example Usage

```hcl
data "tencentcloud_mariadb_database_objects" "database_objects" {
            }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMariadbDatabaseObjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDatabaseObjectsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"db_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},

			"tables": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Table list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Table name.",
						},
					},
				},
			},

			"views": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "View list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"view": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "View name.",
						},
					},
				},
			},

			"procs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Proc list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proc name.",
						},
					},
				},
			},

			"funcs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Func list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"func": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Func name.",
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

	paramMap := make(map[string]interface{})
	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbDatabaseObjectsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceId = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceId))
	if instanceId != nil {
		_ = d.Set("instance_id", instanceId)
	}

	if dbName != nil {
		_ = d.Set("db_name", dbName)
	}

	if tables != nil {
		for _, databaseTable := range tables {
			databaseTableMap := map[string]interface{}{}

			if databaseTable.Table != nil {
				databaseTableMap["table"] = databaseTable.Table
			}

			ids = append(ids, *databaseTable.InstanceId)
			tmpList = append(tmpList, databaseTableMap)
		}

		_ = d.Set("tables", tmpList)
	}

	if views != nil {
		for _, databaseView := range views {
			databaseViewMap := map[string]interface{}{}

			if databaseView.View != nil {
				databaseViewMap["view"] = databaseView.View
			}

			ids = append(ids, *databaseView.InstanceId)
			tmpList = append(tmpList, databaseViewMap)
		}

		_ = d.Set("views", tmpList)
	}

	if procs != nil {
		for _, databaseProcedure := range procs {
			databaseProcedureMap := map[string]interface{}{}

			if databaseProcedure.Proc != nil {
				databaseProcedureMap["proc"] = databaseProcedure.Proc
			}

			ids = append(ids, *databaseProcedure.InstanceId)
			tmpList = append(tmpList, databaseProcedureMap)
		}

		_ = d.Set("procs", tmpList)
	}

	if funcs != nil {
		for _, databaseFunction := range funcs {
			databaseFunctionMap := map[string]interface{}{}

			if databaseFunction.Func != nil {
				databaseFunctionMap["func"] = databaseFunction.Func
			}

			ids = append(ids, *databaseFunction.InstanceId)
			tmpList = append(tmpList, databaseFunctionMap)
		}

		_ = d.Set("funcs", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
