package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ses "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ses/v20201002"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSesSendTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSesSendTasksRead,
		Schema: map[string]*schema.Schema{
			"status": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Task status. `1`: to start; `5`: sending; `6`: sending suspended today; `7`: sending error; `10`: sent. To query tasks in all states, do not pass in this parameter.",
			},

			"receiver_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Recipient group ID.",
			},

			"task_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Task type. `1`: immediate; `2`: scheduled; `3`: recurring. To query tasks of all types, do not pass in this parameter.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Data record.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task ID.",
						},
						"from_email_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sender address.",
						},
						"receiver_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Recipient group ID.",
						},
						"task_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task status. `1`: to start; `5`: sending; `6`: sending suspended today; `7`: sending error; `10`: sent.",
						},
						"task_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task type. `1`: immediate; `2`: scheduled; `3`: recurring.",
						},
						"request_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of emails requested to be sent.",
						},
						"send_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of emails sent.",
						},
						"cache_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of emails cached.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task update time.",
						},
						"subject": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Email subject.",
						},
						"template": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Template and template dataNote: This field may return `null`, indicating that no valid value can be found.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"template_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Template ID. If you do not have any template, please create one.",
									},
									"template_data": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Variable parameters in the template. Please use `json.dump` to format the JSON object into a string type. The object is a set of key-value pairs. Each key denotes a variable, which is represented by {{key}}. The key will be replaced with the corresponding value (represented by {{value}}) when sending the email.Note: The parameter value cannot be data of a complex type such as HTML.Example: {name:xxx,age:xx}.",
									},
								},
							},
						},
						"cycle_param": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Parameters of a recurring taskNote: This field may return `null`, indicating that no valid value can be found.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"begin_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Start time of the task.",
									},
									"interval_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Task recurrence in hours.",
									},
									"term_cycle": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies whether to end the cycle. This parameter is used to update the task. Valid values: 0: No; 1: Yes.",
									},
								},
							},
						},
						"timed_param": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Parameters of a scheduled taskNote: This field may return `null`, indicating that no valid value can be found.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"begin_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Start time of a scheduled sending task.",
									},
								},
							},
						},
						"err_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task exception informationNote: This field may return `null`, indicating that no valid value can be found.",
						},
						"receivers_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Recipient group name.",
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

func dataSourceTencentCloudSesSendTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ses_send_tasks.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("status"); v != nil {
		paramMap["Status"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("receiver_id"); v != nil {
		paramMap["ReceiverId"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("task_type"); v != nil {
		paramMap["TaskType"] = helper.IntUint64(v.(int))
	}

	service := SesService{client: meta.(*TencentCloudClient).apiV3Conn}

	var data []*ses.SendTaskData

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSesSendTasksByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(data))
	tmpList := make([]map[string]interface{}, 0, len(data))

	if data != nil {
		for _, sendTaskData := range data {
			sendTaskDataMap := map[string]interface{}{}

			if sendTaskData.TaskId != nil {
				sendTaskDataMap["task_id"] = sendTaskData.TaskId
			}

			if sendTaskData.FromEmailAddress != nil {
				sendTaskDataMap["from_email_address"] = sendTaskData.FromEmailAddress
			}

			if sendTaskData.ReceiverId != nil {
				sendTaskDataMap["receiver_id"] = sendTaskData.ReceiverId
			}

			if sendTaskData.TaskStatus != nil {
				sendTaskDataMap["task_status"] = sendTaskData.TaskStatus
			}

			if sendTaskData.TaskType != nil {
				sendTaskDataMap["task_type"] = sendTaskData.TaskType
			}

			if sendTaskData.RequestCount != nil {
				sendTaskDataMap["request_count"] = sendTaskData.RequestCount
			}

			if sendTaskData.SendCount != nil {
				sendTaskDataMap["send_count"] = sendTaskData.SendCount
			}

			if sendTaskData.CacheCount != nil {
				sendTaskDataMap["cache_count"] = sendTaskData.CacheCount
			}

			if sendTaskData.CreateTime != nil {
				sendTaskDataMap["create_time"] = sendTaskData.CreateTime
			}

			if sendTaskData.UpdateTime != nil {
				sendTaskDataMap["update_time"] = sendTaskData.UpdateTime
			}

			if sendTaskData.Subject != nil {
				sendTaskDataMap["subject"] = sendTaskData.Subject
			}

			if sendTaskData.Template != nil {
				templateMap := map[string]interface{}{}

				if sendTaskData.Template.TemplateID != nil {
					templateMap["template_id"] = sendTaskData.Template.TemplateID
				}

				if sendTaskData.Template.TemplateData != nil {
					templateMap["template_data"] = sendTaskData.Template.TemplateData
				}

				sendTaskDataMap["template"] = []interface{}{templateMap}
			}

			if sendTaskData.CycleParam != nil {
				cycleParamMap := map[string]interface{}{}

				if sendTaskData.CycleParam.BeginTime != nil {
					cycleParamMap["begin_time"] = sendTaskData.CycleParam.BeginTime
				}

				if sendTaskData.CycleParam.IntervalTime != nil {
					cycleParamMap["interval_time"] = sendTaskData.CycleParam.IntervalTime
				}

				if sendTaskData.CycleParam.TermCycle != nil {
					cycleParamMap["term_cycle"] = sendTaskData.CycleParam.TermCycle
				}

				sendTaskDataMap["cycle_param"] = []interface{}{cycleParamMap}
			}

			if sendTaskData.TimedParam != nil {
				timedParamMap := map[string]interface{}{}

				if sendTaskData.TimedParam.BeginTime != nil {
					timedParamMap["begin_time"] = sendTaskData.TimedParam.BeginTime
				}

				sendTaskDataMap["timed_param"] = []interface{}{timedParamMap}
			}

			if sendTaskData.ErrMsg != nil {
				sendTaskDataMap["err_msg"] = sendTaskData.ErrMsg
			}

			if sendTaskData.ReceiversName != nil {
				sendTaskDataMap["receivers_name"] = sendTaskData.ReceiversName
			}

			ids = append(ids, strconv.Itoa(int(*sendTaskData.TaskId)))
			tmpList = append(tmpList, sendTaskDataMap)
		}

		_ = d.Set("data", tmpList)
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
