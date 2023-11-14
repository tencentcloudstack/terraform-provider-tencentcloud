/*
Use this data source to query detailed information of dcdb database_table

Example Usage

```hcl
data "tencentcloud_dcdb_database_table" "database_table" {
  instance_id = "dcdbt-ow7t8lmc"
  db_name = &lt;nil&gt;
  table = &lt;nil&gt;
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

func dataSourceTencentCloudDcdbDatabaseTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDcdbDatabaseTableRead,
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

			"table": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Table name, obtained through the DescribeDatabaseObjects api.",
			},

			"cols": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Column information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"col": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of column.",
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

func dataSourceTencentCloudDcdbDatabaseTableRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dcdb_database_table.read")()
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

	if v, ok := d.GetOk("table"); ok {
		paramMap["Table"] = helper.String(v.(string))
	}

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDcdbDatabaseTableByFilter(ctx, paramMap)
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
