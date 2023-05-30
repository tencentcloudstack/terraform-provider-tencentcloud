/*
Use this data source to query detailed information of postgresql recovery_time

Example Usage

```hcl
data "tencentcloud_postgresql_recovery_time" "recovery_time" {
  db_instance_id = local.pgsql_id
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudPostgresqlRecoveryTime() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlRecoveryTimeRead,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"recovery_begin_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The earliest restoration time (UTC+8).",
			},

			"recovery_end_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The latest restoration time (UTC+8).",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudPostgresqlRecoveryTimeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_postgresql_recovery_time.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		result       *postgresql.DescribeAvailableRecoveryTimeResponseParams
		e            error
		dbInstanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("db_instance_id"); ok {
		paramMap["DBInstanceId"] = helper.String(v.(string))
		dbInstanceId = v.(string)
	}

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e = service.DescribePostgresqlRecoveryTimeByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result.RecoveryBeginTime != nil {
		_ = d.Set("recovery_begin_time", result.RecoveryBeginTime)
	}

	if result.RecoveryEndTime != nil {
		_ = d.Set("recovery_end_time", result.RecoveryEndTime)
	}

	d.SetId(helper.DataResourceIdHash(dbInstanceId))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
