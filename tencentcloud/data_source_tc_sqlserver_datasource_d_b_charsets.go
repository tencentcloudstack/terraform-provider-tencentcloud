/*
Use this data source to query detailed information of sqlserver datasource_d_b_charsets

Example Usage

```hcl
data "tencentcloud_sqlserver_datasource_d_b_charsets" "datasource_d_b_charsets" {
  instance_id = "mssql-j8kv137v"
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

func dataSourceTencentCloudSqlserverDatasourceDBCharsets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverDatasourceDBCharsetsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of mssql-j8kv137v.",
			},

			"database_charsets": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func dataSourceTencentCloudSqlserverDatasourceDBCharsetsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_datasource_d_b_charsets.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	var databaseCharsets []*string

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

	ids := make([]string, 0, len(databaseCharsets))
	if databaseCharsets != nil {
		_ = d.Set("database_charsets", databaseCharsets)
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
