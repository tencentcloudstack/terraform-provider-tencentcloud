package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoInferenceServiceDeploymentLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoInferenceServiceDeploymentLogsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Inference service ID.",
			},

			"record_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Deployment record ID.",
			},

			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Start time of the logs to be retrieved.",
			},

			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "End time of the logs to be retrieved. The default time range (EndTime - StartTime) is the last 7 days.",
			},

			"sort_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort field. Valid value: `timestamp` (log generation time). Default value: `timestamp`.",
			},

			"sort_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sort order. Valid values: `asc` (ascending), `desc` (descending). Default value: `desc`.",
			},

			"deployment_log_info_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Deployment log list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log message content.",
						},
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log generation time.",
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

func dataSourceTencentCloudTeoInferenceServiceDeploymentLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_inference_service_deployment_logs.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_id"); ok {
		paramMap["ServiceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("record_id"); ok {
		paramMap["RecordId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_by"); ok {
		paramMap["SortBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_order"); ok {
		paramMap["SortOrder"] = helper.String(v.(string))
	}

	var respData []*teov20220901.InferenceServiceDeploymentLogInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoInferenceServiceDeploymentLogs(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if len(result) == 0 {
			log.Printf("[DATASOURCE] read empty, skip SetId, teo_inference_service_deployment_logs paramMap=%v", paramMap)
			return resource.NonRetryableError(fmt.Errorf("teo_inference_service_deployment_logs DescribeInferenceServiceDeploymentLogs response is empty"))
		}
		respData = result
		return nil
	})
	if reqErr != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		return reqErr
	}

	ids := make([]string, 0, len(respData))
	deploymentLogInfoSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, deploymentLogInfo := range respData {
			deploymentLogInfoMap := map[string]interface{}{}

			if deploymentLogInfo.LogMessage != nil {
				deploymentLogInfoMap["log_message"] = deploymentLogInfo.LogMessage
			}

			if deploymentLogInfo.Timestamp != nil {
				deploymentLogInfoMap["timestamp"] = deploymentLogInfo.Timestamp
				ids = append(ids, *deploymentLogInfo.Timestamp)
			}

			deploymentLogInfoSetList = append(deploymentLogInfoSetList, deploymentLogInfoMap)
		}

		_ = d.Set("deployment_log_info_set", deploymentLogInfoSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), deploymentLogInfoSetList); e != nil {
			return e
		}
	}

	return nil
}
