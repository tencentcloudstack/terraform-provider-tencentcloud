package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"database": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},
			"table": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Data Table Name.",
			},
			"table_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Data table type: view: only return view, base_ Table: only returns the basic table, all: returns the view and table.",
			},
			"tables": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Note to the data table list: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database name note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"tables": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Table Name List Note: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudCynosdbClusterRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_cluster.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		tables    []*cynosdb.DatabaseTables
		clusterId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("database"); ok {
		paramMap["Database"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table"); ok {
		paramMap["Table"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table_type"); ok {
		paramMap["TableType"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbClusterByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		tables = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(tables))

	if tables != nil {
		for _, databaseTables := range tables {
			databaseTablesMap := map[string]interface{}{}

			if databaseTables.Database != nil {
				databaseTablesMap["database"] = databaseTables.Database
			}

			if databaseTables.Tables != nil {
				databaseTablesMap["tables"] = databaseTables.Tables
			}

			tmpList = append(tmpList, databaseTablesMap)
		}

		_ = d.Set("tables", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
