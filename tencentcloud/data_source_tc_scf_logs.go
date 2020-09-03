/*
Use this data source to query SCF function logs.

Example Usage

```hcl
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cos_bucket_name   = "scf-code-1234567890"
  cos_object_name   = "code.zip"
  cos_bucket_region = "ap-guangzhou"
}

data "tencentcloud_scf_logs" "foo" {
  function_name = tencentcloud_scf_function.foo.name
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudScfLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfLogsRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the SCF function to be queried.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Log offset, default is `0`, offset+limit cannot be greater than 10000.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10000,
				Description: "Number of logs, the default is `10000`, offset+limit cannot be greater than 10000.",
			},
			"order": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(SCF_LOGS_ORDERS),
				Default:      SCF_LOGS_ORDER_DESC,
				Description:  "Order to sort the log, optional values `desc` and `asc`, default `desc`.",
			},
			"order_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(SCF_LOGS_ORDER_BY),
				Default:      SCF_LOGS_ORDER_BY_START_TIME,
				Description:  "Sort the logs according to the following fields: `function_name`, `duration`, `mem_usage`, `start_time`, default `start_time`.",
			},
			"ret_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(SCF_LOGS_RET_CODES),
				Description:  "Use to filter log, optional value: `not0` only returns the error log. `is0` only returns the correct log. `TimeLimitExceeded` returns the log of the function call timeout. `ResourceLimitExceeded` returns the function call generation resource overrun log. `UserCodeException` returns logs of the user code error that occurred in the function call. Not passing the parameter means returning all logs.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Namespace of the SCF function to be queried.",
			},
			"invoke_request_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Corresponding requestId when executing function.",
			},
			"start_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateTime(SCF_LOGS_DESCRIBE_TIME_FORMAT),
				Description:  "The start time of the query, the format is `2017-05-16 20:00:00`, which can only be within one day from `end_time`.",
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateTime(SCF_LOGS_DESCRIBE_TIME_FORMAT),
				Description:  "The end time of the query, the format is `2017-05-16 20:00:00`, which can only be within one day from `start_time`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of logs. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SCF function.",
						},
						"ret_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Return value after function execution is completed.",
						},
						"request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execute the requestId corresponding to the function.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Point in time at which the function begins execution.",
						},
						"ret_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Execution result of function, `0` means the execution is successful, other values indicate failure.",
						},
						"invoke_finished": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the function call ends, `1` means the execution ends, other values indicate the call exception.",
						},
						"duration": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Function execution time-consuming, unit is ms.",
						},
						"bill_duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Function billing time, according to duration up to the last 100ms, unit is ms.",
						},
						"mem_usage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The actual memory size consumed in the execution of the function, unit is Byte.",
						},
						"log": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log output during function execution.",
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log level.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log source.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudScfLogsRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_scf_logs.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := ScfService{client: m.(*TencentCloudClient).apiV3Conn}

	functionName := d.Get("function_name").(string)
	namespace := d.Get("namespace").(string)

	offset := d.Get("offset").(int)
	limit := d.Get("limit").(int)
	if offset+limit > 10000 {
		return errors.New("offset + limit can't greater than 10000")
	}

	order := d.Get("order").(string)
	orderBy := d.Get("order_by").(string)

	var (
		retCode         *string
		invokeRequestId *string
		startTime       *string
		endTime         *string
	)

	if raw, ok := d.GetOk("ret_code"); ok {
		retCode = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("invoke_request_id"); ok {
		invokeRequestId = helper.String(raw.(string))
	}

	if raw, ok := d.GetOk("start_time"); ok {
		startTime = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("end_time"); ok {
		endTime = helper.String(raw.(string))
	}
	if err := helper.CheckIfSetTogether(d, "start_time", "end_time"); err != nil {
		return err
	}

	if startTime != nil && endTime != nil {
		startTime, _ := time.Parse(SCF_LOGS_DESCRIBE_TIME_FORMAT, *startTime)
		endTime, _ := time.Parse(SCF_LOGS_DESCRIBE_TIME_FORMAT, *endTime)

		if endTime.Sub(startTime) > 24*time.Hour {
			return errors.New("end_time - start_time can't greater then 1 day")
		}
	}

	respLogs, err := service.DescribeLogs(ctx,
		functionName, namespace, order, orderBy,
		offset, limit,
		retCode, invokeRequestId, startTime, endTime,
	)
	if err != nil {
		log.Printf("[CRITAL]%s read function logs failed: %+v", logId, err)
		return err
	}

	logs := make([]map[string]interface{}, 0, len(respLogs))
	ids := make([]string, 0, len(respLogs))
	for _, l := range respLogs {
		ids = append(ids, *l.RequestId)

		logs = append(logs, map[string]interface{}{
			"function_name":   l.FunctionName,
			"ret_msg":         l.RetMsg,
			"request_id":      l.RequestId,
			"start_time":      l.StartTime,
			"ret_code":        l.RetCode,
			"invoke_finished": l.InvokeFinished,
			"duration":        l.Duration,
			"bill_duration":   l.BillDuration,
			"mem_usage":       l.MemUsage,
			"log":             l.Log,
			"level":           l.Level,
			"source":          l.Source,
		})
	}

	_ = d.Set("logs", logs)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), logs); err != nil {
			err = errors.WithStack(err)
			log.Printf("[CRITAL]%s output file[%s] fail, reason: %+v", logId, output.(string), err)
			return err
		}
	}

	return nil
}
