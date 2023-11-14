/*
Use this data source to query detailed information of mariadb database_table

Example Usage

```hcl
data "tencentcloud_mariadb_database_table" "database_table" {
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

func dataSourceTencentCloudMariadbDatabaseTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDatabaseTableRead,
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

			"table": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Table name.",
			},

			"cols": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Column list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"col": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Column name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Column type.",
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

func dataSourceTencentCloudMariadbDatabaseTableRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_database_table.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbDatabaseTableByFilter(ctx, paramMap)
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

	if table != nil {
		_ = d.Set("table", table)
	}

	if cols != nil {
		for _, tableColumn := range cols {
			tableColumnMap := map[string]interface{}{}

			if tableColumn.Col != nil {
				tableColumnMap["col"] = tableColumn.Col
			}

			if tableColumn.Type != nil {
				tableColumnMap["type"] = tableColumn.Type
			}

			ids = append(ids, *tableColumn.InstanceId)
			tmpList = append(tmpList, tableColumnMap)
		}

		_ = d.Set("cols", tmpList)
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
