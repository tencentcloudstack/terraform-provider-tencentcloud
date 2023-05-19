/*
Use this data source to query detailed information of dbbrain top_space_tables

Example Usage

Sort by PhysicalFileSize
```hcl
data "tencentcloud_dbbrain_top_space_tables" "top_space_tables" {
  instance_id = "%s"
  sort_by = "PhysicalFileSize"
  product = "mysql"
}
```

Sort by TotalLength
```hcl
data "tencentcloud_dbbrain_top_space_tables" "top_space_tables" {
  instance_id = "%s"
  sort_by = "PhysicalFileSize"
  product = "mysql"
}
```
*/
package tencentcloud

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainTopSpaceTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainTopSpaceTablesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance id.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     20,
				Description: "The number of Top tables returned, the maximum value is 100, and the default is 20.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The sorting field used to filter the Top table. The optional fields include DataLength, IndexLength, TotalLength, DataFree, FragRatio, TableRows, and PhysicalFileSize (only supported by ApsaraDB for MySQL instances). The default for ApsaraDB for MySQL instances is PhysicalFileSize, and the default for other product instances is TotalLength.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
			},

			"top_space_tables": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The list of Top tablespace statistics returned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "table name.",
						},
						"table_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "database name.",
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage engine for database tables.",
						},
						"data_length": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "data space (MB).",
						},
						"index_length": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Index space (MB).",
						},
						"data_free": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Fragmentation space (MB).",
						},
						"total_length": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Total space used (MB).",
						},
						"frag_ratio": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Fragmentation rate (%).",
						},
						"table_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of lines.",
						},
						"physical_file_size": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The independent physical file size (MB) corresponding to the table.",
						},
					},
				},
			},

			"timestamp": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The timestamp (in seconds) of collecting tablespace data.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainTopSpaceTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_top_space_tables.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var instanceId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var topSpaceTables []*dbbrain.TableSpaceData
	var timestamp *int64

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, ts, e := service.DescribeDbbrainTopSpaceTablesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		topSpaceTables = result
		timestamp = ts
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(topSpaceTables))
	tmpList := make([]map[string]interface{}, 0, len(topSpaceTables))

	if topSpaceTables != nil {
		for _, tableSpaceData := range topSpaceTables {
			tableSpaceDataMap := map[string]interface{}{}

			if tableSpaceData.TableName != nil {
				tableSpaceDataMap["table_name"] = tableSpaceData.TableName
			}

			if tableSpaceData.TableSchema != nil {
				tableSpaceDataMap["table_schema"] = tableSpaceData.TableSchema
			}

			if tableSpaceData.Engine != nil {
				tableSpaceDataMap["engine"] = tableSpaceData.Engine
			}

			if tableSpaceData.DataLength != nil {
				tableSpaceDataMap["data_length"] = tableSpaceData.DataLength
			}

			if tableSpaceData.IndexLength != nil {
				tableSpaceDataMap["index_length"] = tableSpaceData.IndexLength
			}

			if tableSpaceData.DataFree != nil {
				tableSpaceDataMap["data_free"] = tableSpaceData.DataFree
			}

			if tableSpaceData.TotalLength != nil {
				tableSpaceDataMap["total_length"] = tableSpaceData.TotalLength
			}

			if tableSpaceData.FragRatio != nil {
				tableSpaceDataMap["frag_ratio"] = tableSpaceData.FragRatio
			}

			if tableSpaceData.TableRows != nil {
				tableSpaceDataMap["table_rows"] = tableSpaceData.TableRows
			}

			if tableSpaceData.PhysicalFileSize != nil {
				tableSpaceDataMap["physical_file_size"] = tableSpaceData.PhysicalFileSize
			}

			ids = append(ids, strings.Join([]string{instanceId, *tableSpaceData.TableSchema, *tableSpaceData.TableName}, FILED_SP))
			tmpList = append(tmpList, tableSpaceDataMap)
		}

		_ = d.Set("top_space_tables", tmpList)
	}

	if timestamp != nil {
		_ = d.Set("timestamp", timestamp)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
