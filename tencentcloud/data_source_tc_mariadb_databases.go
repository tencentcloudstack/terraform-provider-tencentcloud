/*
Use this data source to query detailed information of mariadb databases

Example Usage

```hcl
data "tencentcloud_mariadb_databases" "databases" {
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMariadbDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDatabasesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"databases": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The database list of this instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name.",
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

func dataSourceTencentCloudMariadbDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mariadb_databases.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var databases []*mariadb.Database

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbDatabasesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		databases = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(databases))
	tmpList := make([]map[string]interface{}, 0, len(databases))

	if instanceId != nil {
		_ = d.Set("instance_id", instanceId)
	}

	if databases != nil {
		for _, database := range databases {
			databaseMap := map[string]interface{}{}

			if database.DbName != nil {
				databaseMap["db_name"] = database.DbName
			}

			ids = append(ids, *database.InstanceId)
			tmpList = append(tmpList, databaseMap)
		}

		_ = d.Set("databases", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
