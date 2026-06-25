package dbdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomImages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomImagesRead,
		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"image_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DB Custom available OS image list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image ID.",
						},
						"os_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OS name.",
						},
						"image_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Image type. Values: PUBLIC_IMAGE (TencentCloud official image), PRIVATE_IMAGE (customer dedicated image).",
						},
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OS architecture. Values: x86_64, arm64.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDbdcDbCustomImagesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_images.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})

	var respData []*dbdcv20201029.DBCustomImage
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, _, e := service.DescribeDBCustomImagesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	imageSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, image := range respData {
			imageMap := map[string]interface{}{}
			if image.ImageId != nil {
				imageMap["image_id"] = image.ImageId
			}

			if image.OsName != nil {
				imageMap["os_name"] = image.OsName
			}

			if image.ImageType != nil {
				imageMap["image_type"] = image.ImageType
			}

			if image.Architecture != nil {
				imageMap["architecture"] = image.Architecture
			}

			imageSetList = append(imageSetList, imageMap)
		}

		_ = d.Set("image_set", imageSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
