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

func DataSourceTencentCloudDbbrainTopSpaceSchemaTimeSeries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainTopSpaceSchemaTimeSeriesRead,
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

			"start_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The start date, such as 2021-01-01, the earliest is the 29th day before the current day, and the default is the 6th day before the deadline.",
			},

			"end_date": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The deadline, such as 2021-01-01, the earliest is the 29th day before the current day, and the default is the current day.",
			},

			"product": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL, the default is mysql.",
			},

			"top_space_schema_time_series": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The time series data list of the returned top library space statistics.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table_schema": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "databases name.",
						},
						"series_data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Spatial index data in unit time interval.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"series": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Monitor metrics.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Indicator name.",
												},
												"unit": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Indicator unit.",
												},
												"values": {
													Type:     schema.TypeSet,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeFloat,
													},
													Description: "Index value. Note: This field may return null, indicating that no valid value can be obtained.",
												},
											},
										},
									},
									"timestamp": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "The timestamp corresponding to the monitoring indicator.",
									},
								},
							},
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

func dataSourceTencentCloudDbbrainTopSpaceSchemaTimeSeriesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbbrain_top_space_schema_time_series.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
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

	if v, ok := d.GetOk("start_date"); ok {
		paramMap["StartDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		paramMap["EndDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["Product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var topSpaceSchemaTimeSeries []*dbbrain.SchemaSpaceTimeSeries

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainTopSpaceSchemaTimeSeriesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		topSpaceSchemaTimeSeries = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(topSpaceSchemaTimeSeries))
	tmpList := make([]map[string]interface{}, 0, len(topSpaceSchemaTimeSeries))

	if topSpaceSchemaTimeSeries != nil {
		for _, schemaSpaceTimeSeries := range topSpaceSchemaTimeSeries {
			schemaSpaceTimeSeriesMap := map[string]interface{}{}

			if schemaSpaceTimeSeries.TableSchema != nil {
				schemaSpaceTimeSeriesMap["table_schema"] = schemaSpaceTimeSeries.TableSchema
			}

			if schemaSpaceTimeSeries.SeriesData != nil {
				seriesDataMap := map[string]interface{}{}

				if schemaSpaceTimeSeries.SeriesData.Series != nil {
					seriesList := []interface{}{}
					for _, series := range schemaSpaceTimeSeries.SeriesData.Series {
						seriesMap := map[string]interface{}{}

						if series.Metric != nil {
							seriesMap["metric"] = series.Metric
						}

						if series.Unit != nil {
							seriesMap["unit"] = series.Unit
						}

						if series.Values != nil {
							seriesMap["values"] = series.Values
						}

						seriesList = append(seriesList, seriesMap)
					}
					seriesDataMap["series"] = seriesList
				}

				if schemaSpaceTimeSeries.SeriesData.Timestamp != nil {
					seriesDataMap["timestamp"] = schemaSpaceTimeSeries.SeriesData.Timestamp
				}

				schemaSpaceTimeSeriesMap["series_data"] = []interface{}{seriesDataMap}
			}

			ids = append(ids, strings.Join([]string{instanceId, *schemaSpaceTimeSeries.TableSchema}, tccommon.FILED_SP))
			tmpList = append(tmpList, schemaSpaceTimeSeriesMap)
		}

		_ = d.Set("top_space_schema_time_series", tmpList)
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
