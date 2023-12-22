package scf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudScfAsyncEventManagement() *schema.Resource {
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
	defer tccommon.LogElapsed("data_source.tencentcloud_scf_async_event_management.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("orderby"); ok {
		paramMap["Orderby"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("invoke_request_id"); ok {
		paramMap["InvokeRequestId"] = helper.String(v.(string))
	}

	service := ScfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var eventList []*scf.AsyncEvent

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeScfAsyncEventManagementByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
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

			ids = append(ids, *asyncEvent.InvokeRequestId)
			tmpList = append(tmpList, asyncEventMap)
		}

		_ = d.Set("event_list", tmpList)
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
