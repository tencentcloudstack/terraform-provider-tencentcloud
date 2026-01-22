package dlc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcDataEngineSessionParameters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDataEngineSessionParametersRead,
		Schema: map[string]*schema.Schema{
			"data_engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DataEngine Id.",
			},

			"data_engine_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Engine name. When the engine name is specified, the name is used first to obtain the configuration.",
			},

			"data_engine_parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Engine Session Configuration List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration ID.",
						},
						"child_image_version_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Minor version image ID.",
						},
						"engine_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Cluster type: SparkSQL/PrestoSQL/SparkBatch.",
						},
						"key_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Parameter key.",
						},
						"key_description": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Description of the key.",
						},
						"value_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of the value.",
						},
						"value_length_limit": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Length limit of the value.",
						},
						"value_regexp_limit": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Regular expression constraint for the value.",
						},
						"value_default": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Default value.",
						},
						"is_public": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Whether it is a public version: 1 for public; 2 for private.",
						},
						"parameter_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Configuration type: 1 for session config (default); 2 for common config; 3 for cluster config.",
						},
						"submit_method": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Submission method: User or BackGround.",
						},
						"operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Operator.",
						},
						"insert_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Insert time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Update time.",
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

func dataSourceTencentCloudDlcDataEngineSessionParametersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_data_engine_session_parameters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(nil)
		ctx          = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service      = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		dataEngineId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("data_engine_id"); ok {
		paramMap["DataEngineId"] = helper.String(v.(string))
		dataEngineId = v.(string)
	}

	if v, ok := d.GetOk("data_engine_name"); ok {
		paramMap["DataEngineName"] = helper.String(v.(string))
	}

	var respData []*dlcv20210125.DataEngineImageSessionParameter
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDataEngineSessionParametersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	dataEngineParametersList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, dataEngineParameters := range respData {
			dataEngineParametersMap := map[string]interface{}{}
			if dataEngineParameters.ParameterId != nil {
				dataEngineParametersMap["parameter_id"] = dataEngineParameters.ParameterId
			}

			if dataEngineParameters.ChildImageVersionId != nil {
				dataEngineParametersMap["child_image_version_id"] = dataEngineParameters.ChildImageVersionId
			}

			if dataEngineParameters.EngineType != nil {
				dataEngineParametersMap["engine_type"] = dataEngineParameters.EngineType
			}

			if dataEngineParameters.KeyName != nil {
				dataEngineParametersMap["key_name"] = dataEngineParameters.KeyName
			}

			if dataEngineParameters.KeyDescription != nil {
				dataEngineParametersMap["key_description"] = dataEngineParameters.KeyDescription
			}

			if dataEngineParameters.ValueType != nil {
				dataEngineParametersMap["value_type"] = dataEngineParameters.ValueType
			}

			if dataEngineParameters.ValueLengthLimit != nil {
				dataEngineParametersMap["value_length_limit"] = dataEngineParameters.ValueLengthLimit
			}

			if dataEngineParameters.ValueRegexpLimit != nil {
				dataEngineParametersMap["value_regexp_limit"] = dataEngineParameters.ValueRegexpLimit
			}

			if dataEngineParameters.ValueDefault != nil {
				dataEngineParametersMap["value_default"] = dataEngineParameters.ValueDefault
			}

			if dataEngineParameters.IsPublic != nil {
				dataEngineParametersMap["is_public"] = dataEngineParameters.IsPublic
			}

			if dataEngineParameters.ParameterType != nil {
				dataEngineParametersMap["parameter_type"] = dataEngineParameters.ParameterType
			}

			if dataEngineParameters.SubmitMethod != nil {
				dataEngineParametersMap["submit_method"] = dataEngineParameters.SubmitMethod
			}

			if dataEngineParameters.Operator != nil {
				dataEngineParametersMap["operator"] = dataEngineParameters.Operator
			}

			if dataEngineParameters.InsertTime != nil {
				dataEngineParametersMap["insert_time"] = dataEngineParameters.InsertTime
			}

			if dataEngineParameters.UpdateTime != nil {
				dataEngineParametersMap["update_time"] = dataEngineParameters.UpdateTime
			}

			dataEngineParametersList = append(dataEngineParametersList, dataEngineParametersMap)
		}

		_ = d.Set("data_engine_parameters", dataEngineParametersList)
	}

	d.SetId(dataEngineId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataEngineParametersList); e != nil {
			return e
		}
	}

	return nil
}
