package scf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfRequestStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudScfRequestStatusRead,
		Schema: map[string]*schema.Schema{
			"function_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Function name.",
			},

			"function_request_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the request to be queried.",
			},

			"namespace": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Function namespace.",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Start time of the query, for example `2017-05-16 20:00:00`. If it's left empty, it defaults to 15 minutes before the current time.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "End time of the query. such as `2017-05-16 20:59:59`. If `StartTime` is not specified, `EndTime` defaults to the current time. If `StartTime` is specified, `EndTime` is required, and it need to be later than the `StartTime`.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Details of the function running statusNote: this field may return `null`, indicating that no valid values can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"function_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Function name.",
						},
						"ret_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Return value after the function is executed.",
						},
						"request_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request ID.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Request start time.",
						},
						"ret_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Result of the request. `0`: succeeded, `1`: running, `-1`: exception.",
						},
						"duration": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Time consumed for the request in ms.",
						},
						"mem_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Time consumed by the request in MB.",
						},
						"retry_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retry Attempts.",
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

func dataSourceTencentCloudScfRequestStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_request_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("function_name"); ok {
		paramMap["FunctionName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("function_request_id"); ok {
		paramMap["FunctionRequestId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var data []*scf.RequestStatus

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfRequestStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
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
		for _, requestStatus := range data {
			requestStatusMap := map[string]interface{}{}

			if requestStatus.FunctionName != nil {
				requestStatusMap["function_name"] = requestStatus.FunctionName
			}

			if requestStatus.RetMsg != nil {
				requestStatusMap["ret_msg"] = requestStatus.RetMsg
			}

			if requestStatus.RequestId != nil {
				requestStatusMap["request_id"] = requestStatus.RequestId
			}

			if requestStatus.StartTime != nil {
				requestStatusMap["start_time"] = requestStatus.StartTime
			}

			if requestStatus.RetCode != nil {
				requestStatusMap["ret_code"] = requestStatus.RetCode
			}

			if requestStatus.Duration != nil {
				requestStatusMap["duration"] = requestStatus.Duration
			}

			if requestStatus.MemUsage != nil {
				requestStatusMap["mem_usage"] = requestStatus.MemUsage
			}

			if requestStatus.RetryNum != nil {
				requestStatusMap["retry_num"] = requestStatus.RetryNum
			}

			ids = append(ids, *requestStatus.FunctionName)
			tmpList = append(tmpList, requestStatusMap)
		}

		_ = d.Set("data", tmpList)
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
