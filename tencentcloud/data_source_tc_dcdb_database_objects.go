/*
Use this data source to query detailed information of dcdb database_objects

Example Usage

```hcl
data "tencentcloud_dcdb_database_objects" "database_objects" {
  instance_id = "dcdbt-ow7t8lmc"
  db_name = &lt;nil&gt;
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

func dataSourceTencentCloudDcdbDatabaseObjects() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbDatabaseObjectsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"db_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database name, obtained through the DescribeDatabases api.",
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
							Description: "The name of table.",
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
							Description: "The name of view.",
						},
					},
				},
			},

			"procs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Procedure list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of procedure.",
						},
					},
				},
			},

			"funcs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Function list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"func": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of function.",
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

func dataSourceTencentCloudDcdbDatabaseObjectsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_database_objects.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		paramMap["DbName"] = helper.String(v.(string))
	}

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcdbDatabaseObjectsByFilter(ctx, paramMap)
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
