/*
Use this data source to query detailed information of eb search

Example Usage

```hcl
data "tencentcloud_eb_search" "search" {
  start_time =
  end_time =
  event_bus_id = ""
  group_field = ""
  filter {
		key = ""
		operator = ""
		value = ""
		type = ""
		filters {
			key = ""
			operator = ""
			value = ""
		}

  }
  order_fields =
  order_by = ""
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEbSearch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbSearchRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Start time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "End time.",
			},

			"event_bus_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Event bus Id.",
			},

			"group_field": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Aggregate field.",
			},

			"filter": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter field name.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Operator, congruent eq, not equal neq, similar like, exclude similar not like, less than lt, less than and equal to lte, greater than gt, greater than and equal to gte, in range range, not in range norange.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter value, range operation needs to enter two values at the same time, separated by commas.",
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The logical relationship of the level filters, the value AND or OR.",
						},
						"filters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "LogFilters array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Filter field name.",
									},
									"operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Operator, congruent eq, not equal neq, similar like, exclude similar not like, less than lt, less than and equal to lte, greater than gt, greater than and equal to gte, within range range, not within range norange.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Filter values, range operations need to enter two values at the same time, separated by commas.",
									},
								},
							},
						},
					},
				},
			},

			"order_fields": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Sort array.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort by, asc from old to new, desc from new to old.",
			},

			"results": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Log search results, note: this field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reporting time of a single log, note: this field may return null, indicating that no valid value can be obtained.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log content details, note: this field may return null, indicating that no valid value can be obtained.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event source, note: this field may return null, indicating that no valid value can be obtained.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event type, note: this field may return null, indicating that no valid value can be obtained.",
						},
						"rule_ids": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event matching rules, note: this field may return null, indicating that no valid value can be obtained.",
						},
						"subject": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID, note: this field may return null, indicating that no valid value can be obtained.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region, Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event status, note: this field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudEbSearchRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_eb_search.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("start_time"); v != nil {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("event_bus_id"); ok {
		paramMap["EventBusId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_field"); ok {
		paramMap["GroupField"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filter"); ok {
		filterSet := v.([]interface{})
		tmpSet := make([]*eb.LogFilter, 0, len(filterSet))

		for _, item := range filterSet {
			logFilter := eb.LogFilter{}
			logFilterMap := item.(map[string]interface{})

			if v, ok := logFilterMap["key"]; ok {
				logFilter.Key = helper.String(v.(string))
			}
			if v, ok := logFilterMap["operator"]; ok {
				logFilter.Operator = helper.String(v.(string))
			}
			if v, ok := logFilterMap["value"]; ok {
				logFilter.Value = helper.String(v.(string))
			}
			if v, ok := logFilterMap["type"]; ok {
				logFilter.Type = helper.String(v.(string))
			}
			if v, ok := logFilterMap["filters"]; ok {
				for _, item := range v.([]interface{}) {
					filtersMap := item.(map[string]interface{})
					logFilters := eb.LogFilters{}
					if v, ok := filtersMap["key"]; ok {
						logFilters.Key = helper.String(v.(string))
					}
					if v, ok := filtersMap["operator"]; ok {
						logFilters.Operator = helper.String(v.(string))
					}
					if v, ok := filtersMap["value"]; ok {
						logFilters.Value = helper.String(v.(string))
					}
					logFilter.Filters = append(logFilter.Filters, &logFilters)
				}
			}
			tmpSet = append(tmpSet, &logFilter)
		}
		paramMap["filter"] = tmpSet
	}

	if v, ok := d.GetOk("order_fields"); ok {
		orderFieldsSet := v.(*schema.Set).List()
		paramMap["OrderFields"] = helper.InterfacesStringsPoint(orderFieldsSet)
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	var results []*string

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbSearchByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		results = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(results))
	if results != nil {
		for _, searchLogResult := range results {
			searchLogResultMap := map[string]interface{}{}

			if searchLogResult.Timestamp != nil {
				searchLogResultMap["timestamp"] = searchLogResult.Timestamp
			}

			if searchLogResult.Message != nil {
				searchLogResultMap["message"] = searchLogResult.Message
			}

			if searchLogResult.Source != nil {
				searchLogResultMap["source"] = searchLogResult.Source
			}

			if searchLogResult.Type != nil {
				searchLogResultMap["type"] = searchLogResult.Type
			}

			if searchLogResult.RuleIds != nil {
				searchLogResultMap["rule_ids"] = searchLogResult.RuleIds
			}

			if searchLogResult.Subject != nil {
				searchLogResultMap["subject"] = searchLogResult.Subject
			}

			if searchLogResult.Region != nil {
				searchLogResultMap["region"] = searchLogResult.Region
			}

			if searchLogResult.Status != nil {
				searchLogResultMap["status"] = searchLogResult.Status
			}

			ids = append(ids, *searchLogResult.EventBusId)
			tmpList = append(tmpList, searchLogResultMap)
		}

		_ = d.Set("results", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
