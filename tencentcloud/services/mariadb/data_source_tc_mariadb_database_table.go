package mariadb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMariadbDatabaseTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMariadbDatabaseTableRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"db_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "database name.",
			},

			"table": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "table name.",
			},

			"cols": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "column list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"col": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "column name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "column type.",
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

func dataSourceTencentCloudMariadbDatabaseTableRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mariadb_database_table.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		instanceId string
		dbName     string
		table      string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_name"); ok {
		dbName = v.(string)
		paramMap["DbName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table"); ok {
		table = v.(string)
		paramMap["Table"] = helper.String(v.(string))
	}

	service := MariadbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var cols []*mariadb.TableColumn
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMariadbDatabaseTableByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		cols = result
		return nil
	})
	if err != nil {
		return err
	}

	if cols != nil {
		tmpList := make([]map[string]interface{}, 0, len(cols))
		for _, tableColumn := range cols {
			tableColumnMap := map[string]interface{}{}

			if tableColumn.Col != nil {
				tableColumnMap["col"] = tableColumn.Col
			}

			if tableColumn.Type != nil {
				tableColumnMap["type"] = tableColumn.Type
			}

			tmpList = append(tmpList, tableColumnMap)
		}

		_ = d.Set("cols", tmpList)
	}

	d.SetId(instanceId + tccommon.FILED_SP + dbName + tccommon.FILED_SP + table)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
