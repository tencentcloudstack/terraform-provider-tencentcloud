/*
Use this data source to query detailed information of scf async_event_management

Example Usage

```hcl
data "tencentcloud_scf_async_event_management" "async_event_management" {
  function_name = "test_function"
  namespace = "test_namespace"
  qualifier = "$LATEST"
  invoke_type = &lt;nil&gt;
  status = &lt;nil&gt;
  start_time_interval {
		start = &lt;nil&gt;
		end = &lt;nil&gt;

  }
  end_time_interval {
		start = "2020-02-02 04:03:03"
		end = "2020-02-02 05:03:03"

  }
  order = "ASC"
  orderby = "StartTime"
  offset = 0
  limit = 20
  invoke_request_id = "xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudScfAsyncEventManagement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfAsyncEventManagementRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace.",
			},

			"qualifier": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter (function version).",
			},

			"invoke_type": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter (invocation type list), Values: CMQ, CKAFKA_TRIGGER, APIGW, COS, TRIGGER_TIMER, MPS_TRIGGER, CLS_TRIGGER, OTHERS.",
			},

			"status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter (event status list), Values: RUNNING, FINISHED, ABORTED, FAILED.",
			},

			"start_time_interval": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Filter (left-closed-right-open range of execution start time).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Start time (inclusive) the format of %Y-%m-%d %H:%M:%S.",
						},
						"end": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "End time (exclusive) the format of %Y-%m-%d %H:%M:%S.",
						},
					},
				},
			},

			"end_time_interval": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Filter (left-closed-right-open range of execution end time).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Start time (inclusive) in the format of %Y-%m-%d %H:%M:%S.",
						},
						"end": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "End time (exclusive) the format of %Y-%m-%d %H:%M:%S.",
						},
					},
				},
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Valid values: ASC, DESC. Default value: DESC.",
			},

			"orderby": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Valid values: StartTime, EndTime. Default value: StartTime.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Data offset. Default value: 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Number of results to be returned. Default value: 20. Maximum value: 100.",
			},

			"invoke_request_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter (event invocation request ID).",
			},

			"event_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Async event list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"invoke_request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invocation request ID.",
						},
						"invoke_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invocation type.",
						},
						"qualifier": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Function version.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event status. Values: `RUNNING`; `FINISHED` (invoked successfully); `ABORTED` (invocation ended); `FAILED` (invocation failed).",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invocation start time in the format of %Y-%m-%d %H:%M:%S.%f.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Invocation end time in the format of %Y-%m-%d %H:%M:%S.%f.",
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

func dataSourceTencentCloudScfAsyncEventManagementRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_async_event_management.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("qualifier"); ok {
		paramMap["Qualifier"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("invoke_type"); ok {
		invokeTypeSet := v.(*schema.Set).List()
		paramMap["InvokeType"] = helper.InterfacesStringsPoint(invokeTypeSet)
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		paramMap["Status"] = helper.InterfacesStringsPoint(statusSet)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "start_time_interval"); ok {
		timeInterval := scf.TimeInterval{}
		if v, ok := dMap["start"]; ok {
			timeInterval.Start = helper.String(v.(string))
		}
		if v, ok := dMap["end"]; ok {
			timeInterval.End = helper.String(v.(string))
		}
		paramMap["start_time_interval"] = &timeInterval
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "end_time_interval"); ok {
		timeInterval := scf.TimeInterval{}
		if v, ok := dMap["start"]; ok {
			timeInterval.Start = helper.String(v.(string))
		}
		if v, ok := dMap["end"]; ok {
			timeInterval.End = helper.String(v.(string))
		}
		paramMap["end_time_interval"] = &timeInterval
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("orderby"); ok {
		paramMap["Orderby"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("invoke_request_id"); ok {
		paramMap["InvokeRequestId"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(*TencentCloudClient).apiV3Conn}

	var eventList []*scf.AsyncEvent

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfAsyncEventManagementByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		eventList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(eventList))
	tmpList := make([]map[string]interface{}, 0, len(eventList))

	if eventList != nil {
		for _, asyncEvent := range eventList {
			asyncEventMap := map[string]interface{}{}

			if asyncEvent.InvokeRequestId != nil {
				asyncEventMap["invoke_request_id"] = asyncEvent.InvokeRequestId
			}

			if asyncEvent.InvokeType != nil {
				asyncEventMap["invoke_type"] = asyncEvent.InvokeType
			}

			if asyncEvent.Qualifier != nil {
				asyncEventMap["qualifier"] = asyncEvent.Qualifier
			}

			if asyncEvent.Status != nil {
				asyncEventMap["status"] = asyncEvent.Status
			}

			if asyncEvent.StartTime != nil {
				asyncEventMap["start_time"] = asyncEvent.StartTime
			}

			if asyncEvent.EndTime != nil {
				asyncEventMap["end_time"] = asyncEvent.EndTime
			}

			ids = append(ids, *asyncEvent.Namespace)
			tmpList = append(tmpList, asyncEventMap)
		}

		_ = d.Set("event_list", tmpList)
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
