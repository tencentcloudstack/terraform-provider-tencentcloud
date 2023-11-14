/*
Use this data source to query detailed information of dbbrain slow_log_user_sql_advice

Example Usage

```hcl
data "tencentcloud_dbbrain_slow_log_user_sql_advice" "slow_log_user_sql_advice" {
  instance_id = ""
  sql_text = ""
  schema = ""
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

func dataSourceTencentCloudDbbrainSlowLogUserSqlAdvice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSlowLogUserSqlAdviceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"sql_text": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SQL statements.",
			},

			"schema": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Library name.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values： mysql - cloud database MySQL; cynosdb - cloud database TDSQL-C for MySQL; dbbrain-mysql - self-built MySQL, the default is mysql.",
			},

			"advices": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "SQL optimization suggestion, which can be parsed into a JSON array, and the output is empty when no optimization is required.",
			},

			"comments": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "SQL optimization suggestion remarks, which can be parsed into a String array, and the output is empty when optimization is not required.",
			},

			"tables": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The DDL information of related tables can be parsed into a JSON array.",
			},

			"sql_plan": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The SQL execution plan can be parsed into JSON, and the output is empty when no optimization is required.",
			},

			"cost": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The cost saving details after SQL optimization can be parsed as JSON, and the output is empty when no optimization is required.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainSlowLogUserSqlAdviceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_slow_log_user_sql_advice.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sql_text"); ok {
		paramMap["SqlText"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("schema"); ok {
		paramMap["Schema"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainSlowLogUserSqlAdviceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		advices = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(advices))
	if advices != nil {
		_ = d.Set("advices", advices)
	}

	if comments != nil {
		_ = d.Set("comments", comments)
	}

	if tables != nil {
		_ = d.Set("tables", tables)
	}

	if sqlPlan != nil {
		_ = d.Set("sql_plan", sqlPlan)
	}

	if cost != nil {
		_ = d.Set("cost", cost)
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
