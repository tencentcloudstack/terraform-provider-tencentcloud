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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
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
		paramMap["instance_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		paramMap["db_name"] = helper.String(v.(string))
	}

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}
	var result *dcdb.DescribeDatabaseObjectsResponseParams

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var e error
		result, e = service.DescribeDcdbDBObjectsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0)
	data := make(map[string]interface{})

	if result != nil {
		tables := result.Tables
		tabList := make([]map[string]interface{}, 0, len(tables))
		if tables != nil {
			for _, databaseTable := range tables {
				databaseTableMap := map[string]interface{}{}

				if databaseTable.Table != nil {
					databaseTableMap["table"] = databaseTable.Table
				}
				tabList = append(tabList, databaseTableMap)
			}
			_ = d.Set("tables", tabList)
			data["tables"] = tabList
		}

		views := result.Views
		viewList := make([]map[string]interface{}, 0, len(views))
		if views != nil {
			for _, databaseView := range views {
				databaseViewMap := map[string]interface{}{}

				if databaseView.View != nil {
					databaseViewMap["view"] = databaseView.View
				}
				viewList = append(viewList, databaseViewMap)
			}
			_ = d.Set("views", viewList)
			data["views"] = viewList
		}

		procs := result.Procs
		procList := make([]map[string]interface{}, 0, len(procs))
		if procs != nil {
			for _, databaseProcedure := range procs {
				databaseProcedureMap := map[string]interface{}{}

				if databaseProcedure.Proc != nil {
					databaseProcedureMap["proc"] = databaseProcedure.Proc
				}
				procList = append(procList, databaseProcedureMap)
			}
			_ = d.Set("procs", procList)
			data["procs"] = procList
		}

		funcs := result.Funcs
		funcList := make([]map[string]interface{}, 0, len(funcs))
		if funcs != nil {
			for _, databaseFunction := range funcs {
				databaseFunctionMap := map[string]interface{}{}

				if databaseFunction.Func != nil {
					databaseFunctionMap["func"] = databaseFunction.Func
				}
				funcList = append(funcList, databaseFunctionMap)
			}
			_ = d.Set("funcs", funcList)
			data["funcs"] = funcList
		}
	}

	ids = append(ids, strings.Join([]string{*result.InstanceId, *result.DbName}, FILED_SP))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), data); e != nil {
			return e
		}
	}
	return nil
}
