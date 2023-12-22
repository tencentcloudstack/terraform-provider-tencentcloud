package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainSlowLogUserSqlAdvice() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSlowLogUserSqlAdviceRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"sql_text": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SQL statements.",
			},

			"schema": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "library name.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values: `mysql` - cloud database MySQL; `cynosdb` - cloud database TDSQL-C for MySQL; `dbbrain-mysql` - self-built MySQL, the default is `mysql`.",
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

			"request_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Unique request ID, returned for every request. The RequestId of the request needs to be provided when locating the problem.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_slow_log_user_sql_advice.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var id string
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		id = v.(string)
		paramMap["instance_id"] = helper.String(id)
	}

	if v, ok := d.GetOk("sql_text"); ok {
		paramMap["sql_text"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("schema"); ok {
		paramMap["schema"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	var result *dbbrain.DescribeUserSqlAdviceResponseParams
	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		var e error
		result, e = service.DescribeDbbrainSlowLogUserSqlAdviceByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if result != nil {
		if result.Advices != nil {
			_ = d.Set("advices", result.Advices)
		}

		if result.Comments != nil {
			_ = d.Set("comments", result.Comments)
		}

		if result.SqlText != nil {
			_ = d.Set("sql_text", result.SqlText)
		}

		if result.Schema != nil {
			_ = d.Set("schema", result.Schema)
		}

		if result.Tables != nil {
			_ = d.Set("tables", result.Tables)
		}

		if result.SqlPlan != nil {
			_ = d.Set("sql_plan", result.SqlPlan)
		}

		if result.Cost != nil {
			_ = d.Set("cost", result.Cost)
		}

	}

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
