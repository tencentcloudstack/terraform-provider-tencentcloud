package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDlcDescribeDataEnginePythonSparkImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcDescribeDataEnginePythonSparkImagesRead,
		Schema: map[string]*schema.Schema{
			"child_image_version_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Engine Image version id.",
			},

			"python_spark_images": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Pyspark image list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spark_image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Spark image unique id.",
						},
						"child_image_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Engine Image version id.",
						},
						"spark_image_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Spark image name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Spark image description information.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Spark image create time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Spark image update time.",
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

func dataSourceTencentCloudDlcDescribeDataEnginePythonSparkImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dlc_describe_data_engine_python_spark_images.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("child_image_version_id"); ok {
		paramMap["ChildImageVersionId"] = helper.String(v.(string))
	}

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	var pythonSparkImages []*dlc.PythonSparkImage

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcDataEnginePythonSparkImagesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		pythonSparkImages = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(pythonSparkImages))
	tmpList := make([]map[string]interface{}, 0, len(pythonSparkImages))

	if pythonSparkImages != nil {
		for _, pythonSparkImage := range pythonSparkImages {
			pythonSparkImageMap := map[string]interface{}{}

			if pythonSparkImage.SparkImageId != nil {
				pythonSparkImageMap["spark_image_id"] = pythonSparkImage.SparkImageId
			}

			if pythonSparkImage.ChildImageVersionId != nil {
				pythonSparkImageMap["child_image_version_id"] = pythonSparkImage.ChildImageVersionId
			}

			if pythonSparkImage.SparkImageVersion != nil {
				pythonSparkImageMap["spark_image_version"] = pythonSparkImage.SparkImageVersion
			}

			if pythonSparkImage.Description != nil {
				pythonSparkImageMap["description"] = pythonSparkImage.Description
			}

			if pythonSparkImage.CreateTime != nil {
				pythonSparkImageMap["create_time"] = pythonSparkImage.CreateTime
			}

			if pythonSparkImage.UpdateTime != nil {
				pythonSparkImageMap["update_time"] = pythonSparkImage.UpdateTime
			}

			ids = append(ids, *pythonSparkImage.ChildImageVersionId)
			tmpList = append(tmpList, pythonSparkImageMap)
		}

		_ = d.Set("python_spark_images", tmpList)
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
