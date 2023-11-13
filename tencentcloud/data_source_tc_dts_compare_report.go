/*
Use this data source to query detailed information of dts compare_report

Example Usage

```hcl
data "tencentcloud_dts_compare_report" "compare_report" {
  job_id = "dts-amm1jw5q"
  compare_task_id = "dts-amm1jw5q-cmp-bmuum7jk"
  difference_limit = 10
  difference_offset = 0
  difference_d_b = "db1"
  difference_table = "t1"
  skipped_limit = 10
  skipped_offset = 0
  skipped_d_b = "db1"
  skipped_table = "t1"
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dts/v20211206"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDtsCompareReport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDtsCompareReportRead,
		Schema: map[string]*schema.Schema{
			"job_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Job id.",
			},

			"compare_task_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Compare task id.",
			},

			"difference_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Limit for inconsistent.",
			},

			"difference_offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Offset for inconsistent.",
			},

			"difference_d_b": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Inconsistent database name.",
			},

			"difference_table": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Inconsistent table name.",
			},

			"skipped_limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Not execute limit.",
			},

			"skipped_offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Not execute offset.",
			},

			"skipped_d_b": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Not execute database name.",
			},

			"skipped_table": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Not execute table name.",
			},

			"abstract": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Abstract result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"options": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Compare config param.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Compare type,dataCheck(full compare)、sampleDataCheck(sample compare)、rowsCount(row count compare).",
									},
									"sample_rate": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Sample rate.",
									},
									"thread_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Thread count, default 1.",
									},
								},
							},
						},
						"objects": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Compare objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_mode": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Object mode, all-all instance, partial-partial object.",
									},
									"object_items": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Object items.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Database name.",
												},
												"db_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Database mode, all/partial.",
												},
												"schema_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Schema name.",
												},
												"table_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Table mode, all/partial.",
												},
												"tables": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Compare tables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"table_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Table name.",
															},
														},
													},
												},
												"view_mode": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "View mode, all/partial.",
												},
												"views": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Compare views.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"view_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "View name.",
															},
														},
													},
												},
											},
										},
									},
									"advanced_objects": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Advanced object, 如account/index/shardkey/schema.",
									},
								},
							},
						},
						"conclusion": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Compare conclusion,same/different.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task status, success/failed.",
						},
						"total_tables": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total tables.",
						},
						"checked_tables": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Checked tables.",
						},
						"different_tables": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Different tables.",
						},
						"skipped_tables": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Skipped tables.",
						},
						"nearly_table_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Nearly table count.",
						},
						"different_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Different rows.",
						},
						"src_sample_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Source rows on sample compare.",
						},
						"dst_sample_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Destination rows on sample compare.",
						},
						"started_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time.",
						},
						"finished_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Finish time.",
						},
					},
				},
			},

			"detail": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Compare detail.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"difference": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Different detail.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total count.",
									},
									"items": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Different detail list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Database.",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Table.",
												},
												"chunk": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Chunk.",
												},
												"src_item": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Src item value.",
												},
												"dst_item": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dst item value.",
												},
												"index_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Index name.",
												},
												"lower_boundary": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Lower boundary.",
												},
												"upper_boundary": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Upper boundary.",
												},
												"cost_time": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Cost time(ms).",
												},
												"finished_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Finish time.",
												},
											},
										},
									},
								},
							},
						},
						"skipped": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Skip table detail.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total count.",
									},
									"items": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Skip list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"db": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Database.",
												},
												"table": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Table.",
												},
												"reason": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Skip reason.",
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

func dataSourceTencentCloudDtsCompareReportRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dts_compare_report.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("job_id"); ok {
		paramMap["JobId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compare_task_id"); ok {
		paramMap["CompareTaskId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("difference_limit"); v != nil {
		paramMap["DifferenceLimit"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("difference_offset"); v != nil {
		paramMap["DifferenceOffset"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("difference_d_b"); ok {
		paramMap["DifferenceDB"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("difference_table"); ok {
		paramMap["DifferenceTable"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("skipped_limit"); v != nil {
		paramMap["SkippedLimit"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("skipped_offset"); v != nil {
		paramMap["SkippedOffset"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("skipped_d_b"); ok {
		paramMap["SkippedDB"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("skipped_table"); ok {
		paramMap["SkippedTable"] = helper.String(v.(string))
	}

	service := DtsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var abstract []*dts.CompareAbstractInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDtsCompareReportByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		abstract = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(abstract))
	if abstract != nil {
		compareAbstractInfoMap := map[string]interface{}{}

		if abstract.Options != nil {
			optionsMap := map[string]interface{}{}

			if abstract.Options.Method != nil {
				optionsMap["method"] = abstract.Options.Method
			}

			if abstract.Options.SampleRate != nil {
				optionsMap["sample_rate"] = abstract.Options.SampleRate
			}

			if abstract.Options.ThreadCount != nil {
				optionsMap["thread_count"] = abstract.Options.ThreadCount
			}

			compareAbstractInfoMap["options"] = []interface{}{optionsMap}
		}

		if abstract.Objects != nil {
			objectsMap := map[string]interface{}{}

			if abstract.Objects.ObjectMode != nil {
				objectsMap["object_mode"] = abstract.Objects.ObjectMode
			}

			if abstract.Objects.ObjectItems != nil {
				objectItemsList := []interface{}{}
				for _, objectItems := range abstract.Objects.ObjectItems {
					objectItemsMap := map[string]interface{}{}

					if objectItems.DbName != nil {
						objectItemsMap["db_name"] = objectItems.DbName
					}

					if objectItems.DbMode != nil {
						objectItemsMap["db_mode"] = objectItems.DbMode
					}

					if objectItems.SchemaName != nil {
						objectItemsMap["schema_name"] = objectItems.SchemaName
					}

					if objectItems.TableMode != nil {
						objectItemsMap["table_mode"] = objectItems.TableMode
					}

					if objectItems.Tables != nil {
						tablesList := []interface{}{}
						for _, tables := range objectItems.Tables {
							tablesMap := map[string]interface{}{}

							if tables.TableName != nil {
								tablesMap["table_name"] = tables.TableName
							}

							tablesList = append(tablesList, tablesMap)
						}

						objectItemsMap["tables"] = []interface{}{tablesList}
					}

					if objectItems.ViewMode != nil {
						objectItemsMap["view_mode"] = objectItems.ViewMode
					}

					if objectItems.Views != nil {
						viewsList := []interface{}{}
						for _, views := range objectItems.Views {
							viewsMap := map[string]interface{}{}

							if views.ViewName != nil {
								viewsMap["view_name"] = views.ViewName
							}

							viewsList = append(viewsList, viewsMap)
						}

						objectItemsMap["views"] = []interface{}{viewsList}
					}

					objectItemsList = append(objectItemsList, objectItemsMap)
				}

				objectsMap["object_items"] = []interface{}{objectItemsList}
			}

			if abstract.Objects.AdvancedObjects != nil {
				objectsMap["advanced_objects"] = abstract.Objects.AdvancedObjects
			}

			compareAbstractInfoMap["objects"] = []interface{}{objectsMap}
		}

		if abstract.Conclusion != nil {
			compareAbstractInfoMap["conclusion"] = abstract.Conclusion
		}

		if abstract.Status != nil {
			compareAbstractInfoMap["status"] = abstract.Status
		}

		if abstract.TotalTables != nil {
			compareAbstractInfoMap["total_tables"] = abstract.TotalTables
		}

		if abstract.CheckedTables != nil {
			compareAbstractInfoMap["checked_tables"] = abstract.CheckedTables
		}

		if abstract.DifferentTables != nil {
			compareAbstractInfoMap["different_tables"] = abstract.DifferentTables
		}

		if abstract.SkippedTables != nil {
			compareAbstractInfoMap["skipped_tables"] = abstract.SkippedTables
		}

		if abstract.NearlyTableCount != nil {
			compareAbstractInfoMap["nearly_table_count"] = abstract.NearlyTableCount
		}

		if abstract.DifferentRows != nil {
			compareAbstractInfoMap["different_rows"] = abstract.DifferentRows
		}

		if abstract.SrcSampleRows != nil {
			compareAbstractInfoMap["src_sample_rows"] = abstract.SrcSampleRows
		}

		if abstract.DstSampleRows != nil {
			compareAbstractInfoMap["dst_sample_rows"] = abstract.DstSampleRows
		}

		if abstract.StartedAt != nil {
			compareAbstractInfoMap["started_at"] = abstract.StartedAt
		}

		if abstract.FinishedAt != nil {
			compareAbstractInfoMap["finished_at"] = abstract.FinishedAt
		}

		ids = append(ids, *abstract.CompareTaskId)
		_ = d.Set("abstract", compareAbstractInfoMap)
	}

	if detail != nil {
		compareDetailInfoMap := map[string]interface{}{}

		if detail.Difference != nil {
			differenceMap := map[string]interface{}{}

			if detail.Difference.TotalCount != nil {
				differenceMap["total_count"] = detail.Difference.TotalCount
			}

			if detail.Difference.Items != nil {
				itemsList := []interface{}{}
				for _, items := range detail.Difference.Items {
					itemsMap := map[string]interface{}{}

					if items.Db != nil {
						itemsMap["db"] = items.Db
					}

					if items.Table != nil {
						itemsMap["table"] = items.Table
					}

					if items.Chunk != nil {
						itemsMap["chunk"] = items.Chunk
					}

					if items.SrcItem != nil {
						itemsMap["src_item"] = items.SrcItem
					}

					if items.DstItem != nil {
						itemsMap["dst_item"] = items.DstItem
					}

					if items.IndexName != nil {
						itemsMap["index_name"] = items.IndexName
					}

					if items.LowerBoundary != nil {
						itemsMap["lower_boundary"] = items.LowerBoundary
					}

					if items.UpperBoundary != nil {
						itemsMap["upper_boundary"] = items.UpperBoundary
					}

					if items.CostTime != nil {
						itemsMap["cost_time"] = items.CostTime
					}

					if items.FinishedAt != nil {
						itemsMap["finished_at"] = items.FinishedAt
					}

					itemsList = append(itemsList, itemsMap)
				}

				differenceMap["items"] = []interface{}{itemsList}
			}

			compareDetailInfoMap["difference"] = []interface{}{differenceMap}
		}

		if detail.Skipped != nil {
			skippedMap := map[string]interface{}{}

			if detail.Skipped.TotalCount != nil {
				skippedMap["total_count"] = detail.Skipped.TotalCount
			}

			if detail.Skipped.Items != nil {
				itemsList := []interface{}{}
				for _, items := range detail.Skipped.Items {
					itemsMap := map[string]interface{}{}

					if items.Db != nil {
						itemsMap["db"] = items.Db
					}

					if items.Table != nil {
						itemsMap["table"] = items.Table
					}

					if items.Reason != nil {
						itemsMap["reason"] = items.Reason
					}

					itemsList = append(itemsList, itemsMap)
				}

				skippedMap["items"] = []interface{}{itemsList}
			}

			compareDetailInfoMap["skipped"] = []interface{}{skippedMap}
		}

		ids = append(ids, *detail.CompareTaskId)
		_ = d.Set("detail", compareDetailInfoMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), compareAbstractInfoMap); e != nil {
			return e
		}
	}
	return nil
}
