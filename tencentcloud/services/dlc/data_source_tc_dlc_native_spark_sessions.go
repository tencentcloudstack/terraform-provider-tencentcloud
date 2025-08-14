package dlc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcNativeSparkSessions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcNativeSparkSessionsRead,
		Schema: map[string]*schema.Schema{
			"data_engine_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Data engine id.",
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource group id.",
			},

			"spark_sessions_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Spark sessions list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spark_session_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Spark session id.",
						},
						"spark_session_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Spark session name.",
						},
						"resource_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource groupid.",
						},
						"engine_session_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine session id.",
						},
						"engine_session_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Engine session name.",
						},
						"idle_timeout_min": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Idle timeout min.",
						},
						"driver_spec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Driver specifications.",
						},
						"executor_spec": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Executor specifications.",
						},
						"executor_num_min": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Minimum number of executors.",
						},
						"executor_num_max": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum number of executors.",
						},
						"total_spec_min": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Minimum specifications.",
						},
						"total_spec_max": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Maximum specifications.",
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

func dataSourceTencentCloudDlcNativeSparkSessionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_native_spark_sessions.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("data_engine_id"); ok {
		paramMap["DataEngineId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_group_id"); ok {
		paramMap["ResourceGroupId"] = helper.String(v.(string))
	}

	var respData []*dlcv20210125.SparkSessionInfo
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcNativeSparkSessionsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	sparkSessionsListList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, sparkSessionsList := range respData {
			sparkSessionsListMap := map[string]interface{}{}
			if sparkSessionsList.SparkSessionId != nil {
				sparkSessionsListMap["spark_session_id"] = sparkSessionsList.SparkSessionId
			}

			if sparkSessionsList.SparkSessionName != nil {
				sparkSessionsListMap["spark_session_name"] = sparkSessionsList.SparkSessionName
			}

			if sparkSessionsList.ResourceGroupId != nil {
				sparkSessionsListMap["resource_group_id"] = sparkSessionsList.ResourceGroupId
			}

			if sparkSessionsList.EngineSessionId != nil {
				sparkSessionsListMap["engine_session_id"] = sparkSessionsList.EngineSessionId
			}

			if sparkSessionsList.EngineSessionName != nil {
				sparkSessionsListMap["engine_session_name"] = sparkSessionsList.EngineSessionName
			}

			if sparkSessionsList.IdleTimeoutMin != nil {
				sparkSessionsListMap["idle_timeout_min"] = sparkSessionsList.IdleTimeoutMin
			}

			if sparkSessionsList.DriverSpec != nil {
				sparkSessionsListMap["driver_spec"] = sparkSessionsList.DriverSpec
			}

			if sparkSessionsList.ExecutorSpec != nil {
				sparkSessionsListMap["executor_spec"] = sparkSessionsList.ExecutorSpec
			}

			if sparkSessionsList.ExecutorNumMin != nil {
				sparkSessionsListMap["executor_num_min"] = sparkSessionsList.ExecutorNumMin
			}

			if sparkSessionsList.ExecutorNumMax != nil {
				sparkSessionsListMap["executor_num_max"] = sparkSessionsList.ExecutorNumMax
			}

			if sparkSessionsList.TotalSpecMin != nil {
				sparkSessionsListMap["total_spec_min"] = sparkSessionsList.TotalSpecMin
			}

			if sparkSessionsList.TotalSpecMax != nil {
				sparkSessionsListMap["total_spec_max"] = sparkSessionsList.TotalSpecMax
			}

			sparkSessionsListList = append(sparkSessionsListList, sparkSessionsListMap)
		}

		_ = d.Set("spark_sessions_list", sparkSessionsListList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), sparkSessionsListList); e != nil {
			return e
		}
	}

	return nil
}
