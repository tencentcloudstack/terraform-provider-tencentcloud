package dbbrain

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainSqlTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainSqlTemplatesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"schema": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "database name.",
			},

			"sql_text": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "SQL statements.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
			},

			"sql_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "sql type.",
			},

			"sql_template": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "SQL template content.",
			},

			"sql_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "SQL template ID.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainSqlTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_sql_templates.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var instanceId string

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("schema"); ok {
		paramMap["Schema"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sql_text"); ok {
		paramMap["SqlText"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var rows *dbbrain.DescribeSqlTemplateResponseParams

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainSqlTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		rows = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := []map[string]interface{}{}

	if rows != nil {

		if rows.SqlType != nil {
			_ = d.Set("sql_type", rows.SqlType)
		}

		if rows.SqlTemplate != nil {
			_ = d.Set("sql_template", rows.SqlTemplate)
		}

		if rows.SqlId != nil {
			_ = d.Set("sql_id", rows.SqlId)
		}
		tmpList = append(tmpList, map[string]interface{}{
			"sql_type":     rows.SqlType,
			"sql_template": rows.SqlTemplate,
			"sql_id":       rows.SqlId,
		})
	}

	d.SetId(helper.DataResourceIdHash(instanceId))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
