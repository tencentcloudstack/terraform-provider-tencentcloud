/*
Use this data source to query detailed information of eb eb_search

Example Usage

```hcl
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_put_events" "put_events" {
  event_list {
    source = "ckafka.cloud.tencent"
    data = jsonencode(
      {
        "topic" : "test-topic",
        "Partition" : 1,
        "offset" : 37,
        "msgKey" : "test",
        "msgBody" : "Hello from Ckafka again!"
      }
    )
    type    = "connector:ckafka"
    subject = "qcs::ckafka:ap-guangzhou:uin/1250000000:ckafkaId/uin/1250000000/ckafka-123456"
    time    = 1691572461939

  }
  event_bus_id = tencentcloud_eb_event_bus.foo.id
}

data "tencentcloud_eb_search" "eb_search" {
  start_time   = 1691637288422
  end_time     = 1691648088422
  event_bus_id = "eb-jzytzr4e"
  group_field = "RuleIds"
  filter {
  	type = "OR"
  	filters {
  		key = "status"
  		operator = "eq"
  		value = "1"
  	}
  }

  filter {
  	type = "OR"
  	filters {
  		key = "type"
  		operator = "eq"
  		value = "connector:ckafka"
  	}
  }
  # order_fields = [""]
  order_by = "desc"
}
```
*/
package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	eb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/eb/v20210416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEbSearch() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbSearchRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "start time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "end time.",
			},

			"event_bus_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "event bus Id.",
			},

			"group_field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "aggregate field, When querying the log index dimension value, you must enter.",
			},

			"filter": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "filter criteria.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "filter field name.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "operator, congruent eq, not equal neq, similar like, exclude similar not like, less than lt, less than and equal to lte, greater than gt, greater than and equal to gte, in range range, not in range norange.",
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
										Description: "filter field name.",
									},
									"operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "operator, congruent eq, not equal neq, similar like, exclude similar not like, less than lt, less than and equal to lte, greater than gt, greater than and equal to gte, within range range, not within range norange.",
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
				Description: "sort array, take effect when the log is retrieved.",
			},

			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort by, asc from old to new, desc from new to old, take effect when the log is retrieved.",
			},

			"dimension_values": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Index retrieves dimension values.",
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

	var (
		startTime  string
		endTime    string
		eventBusId string
		groupField string
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("start_time"); v != nil {
		startTime = strconv.Itoa(v.(int))
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("end_time"); v != nil {
		endTime = strconv.Itoa(v.(int))
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("event_bus_id"); ok {
		eventBusId = v.(string)
		paramMap["EventBusId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_field"); ok {
		groupField = v.(string)
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
		paramMap["Filter"] = tmpSet
	}

	if v, ok := d.GetOk("order_fields"); ok {
		orderFieldsSet := v.(*schema.Set).List()
		paramMap["OrderFields"] = helper.InterfacesStringsPoint(orderFieldsSet)
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	service := EbService{client: meta.(*TencentCloudClient).apiV3Conn}

	if groupField != "" {
		var searchResults []*string
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			response, e := service.DescribeEbSearchByFilter(ctx, paramMap)
			if e != nil {
				return retryError(e)
			}
			searchResults = response
			return nil
		})
		if err != nil {
			return err
		}

		if searchResults != nil {
			_ = d.Set("dimension_values", searchResults)
		}
	}

	var results []*eb.SearchLogResult
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeEbSearchLogByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		results = response
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(results))
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

			tmpList = append(tmpList, searchLogResultMap)
		}

		_ = d.Set("results", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash([]string{startTime, endTime, eventBusId}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}
	return nil
}
