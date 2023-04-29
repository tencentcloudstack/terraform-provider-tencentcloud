/*
Use this data source to query detailed information of mariadb databases

Example Usage

```hcl
data "tencentcloud_mariadb_databases" "databases" {
  instance_id = "tdsql-e9tklsgz"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
)

func dataSourceTencentCloudMariadbDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDatabasesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
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

	instanceId := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MariadbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var databases []*mariadb.Database

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbDatabasesByFilter(ctx, instanceId)
		if e != nil {
			return retryError(e)
		}
		databases = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(databases))
	if databases != nil {
		for _, database := range databases {
			databaseMap := map[string]interface{}{}

			if database.DbName != nil {
				databaseMap["db_name"] = database.DbName
			}

			tmpList = append(tmpList, databaseMap)
		}

		_ = d.Set("databases", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
