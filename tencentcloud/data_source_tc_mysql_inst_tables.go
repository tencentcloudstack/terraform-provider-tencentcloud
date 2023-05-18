/*
Use this data source to query detailed information of mysql inst_tables

Example Usage

```hcl
data "tencentcloud_mysql_inst_tables" "inst_tables" {
  instance_id = "cdb-fitq5t9h"
  database = "tf_ci_test"
  # table_regexp = ""
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

func dataSourceTencentCloudMysqlInstTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstTablesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},

			"database": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The name of the database.",
			},

			"table_regexp": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Match the regular expression of the database table name, the rules are the same as MySQL official website.",
			},

			"items": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The returned database table information.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMysqlInstTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_inst_tables.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database"); ok {
		paramMap["Database"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table_regexp"); ok {
		paramMap["TableRegexp"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var tables []*string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlInstTablesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		tables = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(tables))
	if tables != nil {
		_ = d.Set("items", tables)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
