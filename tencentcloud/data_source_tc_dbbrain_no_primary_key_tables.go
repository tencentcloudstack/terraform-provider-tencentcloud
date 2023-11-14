/*
Use this data source to query detailed information of dbbrain no_primary_key_tables

Example Usage

```hcl
data "tencentcloud_dbbrain_no_primary_key_tables" "no_primary_key_tables" {
  instance_id = ""
  date = ""
  product = ""
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

func dataSourceTencentCloudDbbrainNoPrimaryKeyTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainNoPrimaryKeyTablesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query date, such as 2021-05-27, the earliest date is 30 days ago.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values： mysql - ApsaraDB for MySQL, the default is mysql.",
			},

			"no_primary_key_table_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of tables without a primary key.",
			},

			"no_primary_key_table_count_diff": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The difference with yesterday&amp;amp;#39;s scan of the table without a primary key. A positive number means an increase, a negative number means a decrease, and 0 means no change.",
			},

			"no_primary_key_table_record_count": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The total number of recorded non-primary key tables (no more than the total number of non-primary key tables), which can be used for pagination query.",
			},

			"no_primary_key_tables": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "A list of tables without primary keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Library name.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "TableName.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage engine for database tables.",
						},
						"table_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rows.",
						},
						"total_length": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total space used (MB).",
						},
					},
				},
			},

			"timestamp": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Collection timestamp (seconds).",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainNoPrimaryKeyTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_no_primary_key_tables.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("date"); ok {
		paramMap["Date"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainNoPrimaryKeyTablesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		noPrimaryKeyTableCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(noPrimaryKeyTableCount))
	if noPrimaryKeyTableCount != nil {
		_ = d.Set("no_primary_key_table_count", noPrimaryKeyTableCount)
	}

	if noPrimaryKeyTableCountDiff != nil {
		_ = d.Set("no_primary_key_table_count_diff", noPrimaryKeyTableCountDiff)
	}

	if noPrimaryKeyTableRecordCount != nil {
		_ = d.Set("no_primary_key_table_record_count", noPrimaryKeyTableRecordCount)
	}

	if noPrimaryKeyTables != nil {
		for _, table := range noPrimaryKeyTables {
			tableMap := map[string]interface{}{}

			if table.TableSchema != nil {
				tableMap["table_schema"] = table.TableSchema
			}

			if table.TableName != nil {
				tableMap["table_name"] = table.TableName
			}

			if table.Engine != nil {
				tableMap["engine"] = table.Engine
			}

			if table.TableRows != nil {
				tableMap["table_rows"] = table.TableRows
			}

			if table.TotalLength != nil {
				tableMap["total_length"] = table.TotalLength
			}

			ids = append(ids, *table.InstanceId)
			tmpList = append(tmpList, tableMap)
		}

		_ = d.Set("no_primary_key_tables", tmpList)
	}

	if timestamp != nil {
		_ = d.Set("timestamp", timestamp)
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
