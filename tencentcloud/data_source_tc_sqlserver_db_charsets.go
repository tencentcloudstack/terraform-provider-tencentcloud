/*
Use this data source to query detailed information of sqlserver datasource_d_b_charsets

Example Usage

```hcl
data "tencentcloud_sqlserver_db_charsets" "db_charsets" {
  instance_id = "mssql-qelbzgwf"
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

func dataSourceTencentCloudSqlserverDBCharsets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverDBCharsetsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of mssql-j8kv137v.",
			},
			"database_charsets": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Database character set list.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSqlserverDBCharsetsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_db_charsets.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		service          = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId       string
		databaseCharsets []*string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverDatasourceDBCharsetsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		databaseCharsets = result
		return nil
	})

	if err != nil {
		return err
	}

	if databaseCharsets != nil {
		_ = d.Set("database_charsets", databaseCharsets)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
