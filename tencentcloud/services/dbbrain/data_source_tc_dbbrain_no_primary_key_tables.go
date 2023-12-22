package dbbrain

import (
	"context"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbbrainNoPrimaryKeyTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainNoPrimaryKeyTablesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Query date, such as 2021-05-27, the earliest date is 30 days ago.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values: `mysql` - ApsaraDB for MySQL, the default is `mysql`.",
			},

			"timestamp": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Collection timestamp (seconds).",
			},

			"no_primary_key_table_count_diff": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The difference with yesterday&amp;#39;s scan of the table without a primary key. A positive number means an increase, a negative number means a decrease, and 0 means no change.",
			},

			"no_primary_key_tables": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "A list of tables without primary keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "library name.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "tableName.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage engine for database tables.",
						},
						"table_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "rows.",
						},
						"total_length": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total space used (MB).",
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

func dataSourceTencentCloudDbbrainNoPrimaryKeyTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_no_primary_key_tables.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	var (
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("date"); ok {
		paramMap["Date"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var (
		tables []*dbbrain.Table
		resp   *dbbrain.DescribeNoPrimaryKeyTablesResponseParams
		e      error
	)
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		tables, resp, e = service.DescribeDbbrainNoPrimaryKeyTablesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(tables))
	tmpList := make([]map[string]interface{}, 0, len(tables))

	if resp != nil {
		if resp.NoPrimaryKeyTableCountDiff != nil {
			_ = d.Set("no_primary_key_table_count_diff", resp.NoPrimaryKeyTableCountDiff)
		}

		if resp.Timestamp != nil {
			_ = d.Set("timestamp", resp.Timestamp)
		}
	}

	if tables != nil {
		for _, table := range tables {
			tableMap := map[string]interface{}{}

			if table.TableSchema != nil {
				tableMap["table_schema"] = table.TableSchema
			}

			if table.TableName != nil {
				tableMap["table_name"] = table.TableName
			}

			if table.Engine != nil {
				tableMap["engine"] = table.Engine
			}

			if table.TableRows != nil {
				tableMap["table_rows"] = table.TableRows
			}

			if table.TotalLength != nil {
				tableMap["total_length"] = table.TotalLength
			}

			ids = append(ids, strings.Join([]string{instanceId, *table.TableSchema, *table.TableName}, tccommon.FILED_SP))
			tmpList = append(tmpList, tableMap)
		}

		_ = d.Set("no_primary_key_tables", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
