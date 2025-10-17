package wedata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	wedatav20250806 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20250806"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWedataResourceGroupMetrics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWedataResourceGroupMetricsRead,
		Schema: map[string]*schema.Schema{
			"resource_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Execution resource group ID.",
			},

			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Usage trend start time (milliseconds), default to the last hour.",
			},

			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Usage trend end time (milliseconds), default to current time.",
			},

			"metric_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Metric dimension.\n\n- all --- All\n- task --- Task metrics\n- system --- System metrics.",
			},

			"granularity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Metric collection granularity, unit in minutes, default 1 minute.",
			},

			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Execution group metric information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource group specification related: CPU count.",
						},
						"disk_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource group specification related: disk specification.",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource group specification related: memory size, unit: G.",
						},
						"life_cycle": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource group lifecycle, unit: days.",
						},
						"maximum_concurrency": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource group specification related: maximum concurrency.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Resource group status.\n\n- 0 --- Initializing\n- 1 --- Running\n- 2 --- Running abnormally\n- 3 --- Releasing\n- 4 --- Released\n- 5 --- Creating\n- 6 --- Creation failed\n- 7 --- Updating\n- 8 --- Update failed\n- 9 --- Expired\n- 10 --- Release failed\n- 11 --- In use\n- 12 --- Not in use.",
						},
						"metric_snapshots": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Metric details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Metric name.\n\n- ConcurrencyUsage --- Concurrency usage rate\n- CpuCoreUsage --- CPU usage rate\n- CpuLoad --- CPU load\n- DevelopQueueTask --- Number of development tasks in queue\n- DevelopRunningTask --- Number of running development tasks\n- DevelopSchedulingTask --- Number of scheduling development tasks\n- DiskUsage --- Disk usage\n- DiskUsed --- Disk used amount\n- MaximumConcurrency --- Maximum concurrency\n- MemoryLoad --- Memory load\n- MemoryUsage --- Memory usage.",
									},
									"snapshot_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Current value.",
									},
									"trend_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Metric trend.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timestamp": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Timestamp.",
												},
												"value": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Metric value.",
												},
											},
										},
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

func dataSourceTencentCloudWedataResourceGroupMetricsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_wedata_resource_group_metrics.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(nil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service         = WedataService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		resourceGroupId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("resource_group_id"); ok {
		paramMap["ResourceGroupId"] = helper.String(v.(string))
		resourceGroupId = v.(string)
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("metric_type"); ok {
		paramMap["MetricType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("granularity"); ok {
		paramMap["Granularity"] = helper.IntUint64(v.(int))
	}

	var respData *wedatav20250806.ResourceGroupMetrics
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWedataResourceGroupMetricsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataMap := map[string]interface{}{}
	if respData != nil {
		if respData.CpuNum != nil {
			dataMap["cpu_num"] = respData.CpuNum
		}

		if respData.DiskVolume != nil {
			dataMap["disk_volume"] = respData.DiskVolume
		}

		if respData.MemSize != nil {
			dataMap["mem_size"] = respData.MemSize
		}

		if respData.LifeCycle != nil {
			dataMap["life_cycle"] = respData.LifeCycle
		}

		if respData.MaximumConcurrency != nil {
			dataMap["maximum_concurrency"] = respData.MaximumConcurrency
		}

		if respData.Status != nil {
			dataMap["status"] = respData.Status
		}

		metricSnapshotsList := make([]map[string]interface{}, 0, len(respData.MetricSnapshots))
		if respData.MetricSnapshots != nil {
			for _, metricSnapshots := range respData.MetricSnapshots {
				metricSnapshotsMap := map[string]interface{}{}
				if metricSnapshots.MetricName != nil {
					metricSnapshotsMap["metric_name"] = metricSnapshots.MetricName
				}

				if metricSnapshots.SnapshotValue != nil {
					metricSnapshotsMap["snapshot_value"] = metricSnapshots.SnapshotValue
				}

				trendListList := make([]map[string]interface{}, 0, len(metricSnapshots.TrendList))
				if metricSnapshots.TrendList != nil {
					for _, trendList := range metricSnapshots.TrendList {
						trendListMap := map[string]interface{}{}

						if trendList.Timestamp != nil {
							trendListMap["timestamp"] = trendList.Timestamp
						}

						if trendList.Value != nil {
							trendListMap["value"] = trendList.Value
						}

						trendListList = append(trendListList, trendListMap)
					}

					metricSnapshotsMap["trend_list"] = trendListList
				}

				metricSnapshotsList = append(metricSnapshotsList, metricSnapshotsMap)
			}

			dataMap["metric_snapshots"] = metricSnapshotsList
		}

		_ = d.Set("data", []interface{}{dataMap})
	}

	d.SetId(resourceGroupId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataMap); e != nil {
			return e
		}
	}

	return nil
}
