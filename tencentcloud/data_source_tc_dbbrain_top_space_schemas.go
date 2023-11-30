package tencentcloud

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainTopSpaceSchemas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainTopSpaceSchemasRead,
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
				Description: "The number of Top libraries to return, the maximum value is 100, and the default is 20.",
			},

			"sort_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The sorting field used to filter the Top library. The optional fields include DataLength, IndexLength, TotalLength, DataFree, FragRatio, TableRows, and PhysicalFileSize (only supported by ApsaraDB for MySQL instances). The default for ApsaraDB for MySQL instances is PhysicalFileSize, and the default for other product instances is TotalLength.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
			},

			"top_space_schemas": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The returned list of top library space statistics.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "library name.",
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
							Description: "The sum (MB) of the independent physical file sizes corresponding to all tables in the library. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"timestamp": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Timestamp (in seconds) when library space data is collected.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDbbrainTopSpaceSchemasRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_top_space_schemas.read")()
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

	var topSpaceSchemas []*dbbrain.SchemaSpaceData
	var timestamp *int64

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, ts, e := service.DescribeDbbrainTopSpaceSchemasByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		topSpaceSchemas = result
		timestamp = ts
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(topSpaceSchemas))
	tmpList := make([]map[string]interface{}, 0, len(topSpaceSchemas))

	if topSpaceSchemas != nil {
		for _, schemaSpaceData := range topSpaceSchemas {
			schemaSpaceDataMap := map[string]interface{}{}

			if schemaSpaceData.TableSchema != nil {
				schemaSpaceDataMap["table_schema"] = schemaSpaceData.TableSchema
			}

			if schemaSpaceData.DataLength != nil {
				schemaSpaceDataMap["data_length"] = schemaSpaceData.DataLength
			}

			if schemaSpaceData.IndexLength != nil {
				schemaSpaceDataMap["index_length"] = schemaSpaceData.IndexLength
			}

			if schemaSpaceData.DataFree != nil {
				schemaSpaceDataMap["data_free"] = schemaSpaceData.DataFree
			}

			if schemaSpaceData.TotalLength != nil {
				schemaSpaceDataMap["total_length"] = schemaSpaceData.TotalLength
			}

			if schemaSpaceData.FragRatio != nil {
				schemaSpaceDataMap["frag_ratio"] = schemaSpaceData.FragRatio
			}

			if schemaSpaceData.TableRows != nil {
				schemaSpaceDataMap["table_rows"] = schemaSpaceData.TableRows
			}

			if schemaSpaceData.PhysicalFileSize != nil {
				schemaSpaceDataMap["physical_file_size"] = schemaSpaceData.PhysicalFileSize
			}

			ids = append(ids, strings.Join([]string{instanceId, *schemaSpaceData.TableSchema}, FILED_SP))
			tmpList = append(tmpList, schemaSpaceDataMap)
		}

		_ = d.Set("top_space_schemas", tmpList)
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
