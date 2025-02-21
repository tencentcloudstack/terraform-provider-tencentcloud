package emr

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEmrJobStatusDetail() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEmrJobStatusDetailRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EMR Instance ID.",
			},

			"flow_param": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Flow-related Parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"f_key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Process Parameter Key: value range: TraceId: Query by TraceId FlowId: Query by FlowId.",
						},
						"f_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter Value.",
						},
					},
				},
			},

			"need_extra_detail": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to return additional task information.",
			},

			"stage_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Task Information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Step ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Step Name.",
						},
						"is_show": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to display the flow.",
						},
						"is_sub_flow": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it is a sub-flow.",
						},
						"sub_flow_flag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sub-Flow Flag.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Flow Execution Status: 0: Not Started, 1: In Progress, 2: Completed, 3: Partially Completed, -1: Failed.",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flow Execution Status Description.",
						},
						"progress": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Flow Execution Progress.",
						},
						"starttime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flow Execution Start Time.",
						},
						"endtime": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flow Execution End Time.",
						},
						"had_wood_detail": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to return additional task information.",
						},
						"wood_job_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Wood Subprocess ID.",
						},
						"language_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Multilingual Version Key.",
						},
						"failed_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flow Execution Failure Reason.",
						},
						"time_consuming": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flow Execution Time Consuming.",
						},
					},
				},
			},

			"flow_desc": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Flow Parameter Description.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"p_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter Key.",
						},
						"p_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter Value.",
						},
					},
				},
			},

			"flow_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Flow Name.",
			},

			"flow_total_progress": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Flow Total Execution Progress.",
			},

			"flow_total_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Flow Total Execution Status, 0: Initialized, 1: Running, 2: Completed, 3: Completed (with skipped steps), -1: Failed, -3: Blocke.",
			},

			"flow_extra_detail": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Flow Extra Execution Detail,Return when NeedExtraDetail is true.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Flow Extra Execution Detail Title.",
						},
						"detail": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Flow Extra Execution Detail.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"p_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter Key.",
									},
									"p_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Parameter Value.",
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

func dataSourceTencentCloudEmrJobStatusDetailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_emr_job_status_detail.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := EMRService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(instanceId)
	}

	if flowParamMap, ok := helper.InterfacesHeadMap(d, "flow_param"); ok {
		flowParam := emr.FlowParam{}
		if v, ok := flowParamMap["f_key"]; ok {
			flowParam.FKey = helper.String(v.(string))
		}
		if v, ok := flowParamMap["f_value"]; ok {
			flowParam.FValue = helper.String(v.(string))
		}
		paramMap["FlowParam"] = &flowParam
	}

	if v, ok := d.GetOkExists("need_extra_detail"); ok {
		paramMap["NeedExtraDetail"] = helper.Bool(v.(bool))
	}

	var respData *emr.DescribeClusterFlowStatusDetailResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEmrJobStatusDetailByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	stageDetailsList := make([]map[string]interface{}, 0, len(respData.StageDetails))
	if respData.StageDetails != nil {
		for _, stageDetails := range respData.StageDetails {
			stageDetailsMap := map[string]interface{}{}

			if stageDetails.Stage != nil {
				stageDetailsMap["stage"] = stageDetails.Stage
			}

			if stageDetails.Name != nil {
				stageDetailsMap["name"] = stageDetails.Name
			}

			if stageDetails.IsShow != nil {
				stageDetailsMap["is_show"] = stageDetails.IsShow
			}

			if stageDetails.IsSubFlow != nil {
				stageDetailsMap["is_sub_flow"] = stageDetails.IsSubFlow
			}

			if stageDetails.SubFlowFlag != nil {
				stageDetailsMap["sub_flow_flag"] = stageDetails.SubFlowFlag
			}

			if stageDetails.Status != nil {
				stageDetailsMap["status"] = stageDetails.Status
			}

			if stageDetails.Desc != nil {
				stageDetailsMap["desc"] = stageDetails.Desc
			}

			if stageDetails.Progress != nil {
				stageDetailsMap["progress"] = stageDetails.Progress
			}

			if stageDetails.Starttime != nil {
				stageDetailsMap["starttime"] = stageDetails.Starttime
			}

			if stageDetails.Endtime != nil {
				stageDetailsMap["endtime"] = stageDetails.Endtime
			}

			if stageDetails.HadWoodDetail != nil {
				stageDetailsMap["had_wood_detail"] = stageDetails.HadWoodDetail
			}

			if stageDetails.WoodJobId != nil {
				stageDetailsMap["wood_job_id"] = stageDetails.WoodJobId
			}

			if stageDetails.LanguageKey != nil {
				stageDetailsMap["language_key"] = stageDetails.LanguageKey
			}

			if stageDetails.FailedReason != nil {
				stageDetailsMap["failed_reason"] = stageDetails.FailedReason
			}

			if stageDetails.TimeConsuming != nil {
				stageDetailsMap["time_consuming"] = stageDetails.TimeConsuming
			}

			stageDetailsList = append(stageDetailsList, stageDetailsMap)
		}

		_ = d.Set("stage_details", stageDetailsList)
	}

	flowDescList := make([]map[string]interface{}, 0, len(respData.FlowDesc))
	if respData.FlowDesc != nil {
		for _, flowDesc := range respData.FlowDesc {
			flowDescMap := map[string]interface{}{}

			if flowDesc.PKey != nil {
				flowDescMap["p_key"] = flowDesc.PKey
			}

			if flowDesc.PValue != nil {
				flowDescMap["p_value"] = flowDesc.PValue
			}

			flowDescList = append(flowDescList, flowDescMap)
		}

		_ = d.Set("flow_desc", flowDescList)
	}

	if respData.FlowName != nil {
		_ = d.Set("flow_name", respData.FlowName)
	}

	if respData.FlowTotalProgress != nil {
		_ = d.Set("flow_total_progress", respData.FlowTotalProgress)
	}

	if respData.FlowTotalStatus != nil {
		_ = d.Set("flow_total_status", respData.FlowTotalStatus)
	}

	flowExtraDetailList := make([]map[string]interface{}, 0, len(respData.FlowExtraDetail))
	if respData.FlowExtraDetail != nil {
		for _, flowExtraDetail := range respData.FlowExtraDetail {
			flowExtraDetailMap := map[string]interface{}{}

			if flowExtraDetail.Title != nil {
				flowExtraDetailMap["title"] = flowExtraDetail.Title
			}

			detailList := make([]map[string]interface{}, 0, len(flowExtraDetail.Detail))
			if flowExtraDetail.Detail != nil {
				for _, detail := range flowExtraDetail.Detail {
					detailMap := map[string]interface{}{}

					if detail.PKey != nil {
						detailMap["p_key"] = detail.PKey
					}

					if detail.PValue != nil {
						detailMap["p_value"] = detail.PValue
					}

					detailList = append(detailList, detailMap)
				}

				flowExtraDetailMap["detail"] = detailList
			}
			flowExtraDetailList = append(flowExtraDetailList, flowExtraDetailMap)
		}

		_ = d.Set("flow_extra_detail", flowExtraDetailList)
	}

	d.SetId(instanceId)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
