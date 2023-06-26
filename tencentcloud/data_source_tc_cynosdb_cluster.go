/*
Use this data source to query detailed information of cynosdb cluster

Example Usage

```hcl
data "tencentcloud_cynosdb_cluster" "cluster" {
  cluster_id = "cynosdbmysql-bws8h88b"
  database   = "users"
  table      = "tb_user_name"
  table_type = "all"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"database": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},
			"table": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Data Table Name.",
			},
			"table_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Data table type: view: only return view, base_ Table: only returns the basic table, all: returns the view and table.",
			},
			"tables": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Note to the data table list: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"tables": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Table Name List Note: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudCynosdbClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_cluster.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		tables    []*cynosdb.DatabaseTables
		clusterId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("database"); ok {
		paramMap["Database"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table"); ok {
		paramMap["Table"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table_type"); ok {
		paramMap["TableType"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbClusterByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		tables = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(tables))

	if tables != nil {
		for _, databaseTables := range tables {
			databaseTablesMap := map[string]interface{}{}

			if databaseTables.Database != nil {
				databaseTablesMap["database"] = databaseTables.Database
			}

			if databaseTables.Tables != nil {
				databaseTablesMap["tables"] = databaseTables.Tables
			}

			tmpList = append(tmpList, databaseTablesMap)
		}

		_ = d.Set("tables", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
