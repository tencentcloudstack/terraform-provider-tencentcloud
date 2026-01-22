package dlc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcSessionImageVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcSessionImageVersionRead,
		Schema: map[string]*schema.Schema{
			"data_engine_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data engine ID.",
			},

			"framework_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Framework type: machine learning, Python, Spark ML.",
			},

			"engine_session_images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Engine session image information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spark_image_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Spark image ID.",
						},
						"spark_image_version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Spark image version.",
						},
						"spark_image_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Small version image type. 1: TensorFlow, 2: Pytorch, 3: SK-learn.",
						},
						"spark_image_tag": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Spark image tag.",
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

func dataSourceTencentCloudDlcSessionImageVersionRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_session_image_version.read")()
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

	if v, ok := d.GetOk("framework_type"); ok {
		paramMap["FrameworkType"] = helper.String(v.(string))
	}

	var respData []*dlcv20210125.EngineSessionImage
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcSessionImageVersionByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	engineSessionImagesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, engineSessionImages := range respData {
			engineSessionImagesMap := map[string]interface{}{}
			if engineSessionImages.SparkImageId != nil {
				engineSessionImagesMap["spark_image_id"] = engineSessionImages.SparkImageId
			}

			if engineSessionImages.SparkImageVersion != nil {
				engineSessionImagesMap["spark_image_version"] = engineSessionImages.SparkImageVersion
			}

			if engineSessionImages.SparkImageType != nil {
				engineSessionImagesMap["spark_image_type"] = engineSessionImages.SparkImageType
			}

			if engineSessionImages.SparkImageTag != nil {
				engineSessionImagesMap["spark_image_tag"] = engineSessionImages.SparkImageTag
			}

			engineSessionImagesList = append(engineSessionImagesList, engineSessionImagesMap)
		}

		_ = d.Set("engine_session_images", engineSessionImagesList)
	}

	d.SetId(dataEngineId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), engineSessionImagesList); e != nil {
			return e
		}
	}

	return nil
}
