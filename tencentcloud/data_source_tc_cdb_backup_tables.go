/*
Use this data source to query detailed information of cdb backup_tables

Example Usage

```hcl
data "tencentcloud_cdb_backup_tables" "backup_tables" {
  instance_id = "cdb-c1nl9rpv"
  start_time = "2022-07-12 10:29:20"
  database_name = &lt;nil&gt;
  search_table = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  total_count = &lt;nil&gt;
  items {
		table_name = &lt;nil&gt;

  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCdbBackupTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbBackupTablesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of instance.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time.",
			},

			"database_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Specified database name.",
			},

			"search_table": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The prefix of the data table name to be queried.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page size, the minimum value is 1, and the maximum value is 2000.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "The number of returned data.",
			},

			"items": {
				Type:        schema.TypeList,
				Description: "An array of tables that meet the query condition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_name": {
							Type:        schema.TypeString,
							Description: "The name of table.",
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

func dataSourceTencentCloudCdbBackupTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_backup_tables.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database_name"); ok {
		paramMap["DatabaseName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_table"); ok {
		paramMap["SearchTable"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("items"); ok {
		itemsSet := v.([]interface{})
		tmpSet := make([]*cdb.TableName, 0, len(itemsSet))

		for _, item := range itemsSet {
			tableName := cdb.TableName{}
			tableNameMap := item.(map[string]interface{})

			if v, ok := tableNameMap["table_name"]; ok {
				tableName.TableName = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &tableName)
		}
		paramMap["items"] = tmpSet
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var items []*cdb.TableName

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbBackupTablesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(items))
	tmpList := make([]map[string]interface{}, 0, len(items))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
