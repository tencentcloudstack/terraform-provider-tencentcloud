/*
Use this data source to query detailed information of mps image_sprite_templates

Example Usage

```hcl
data "tencentcloud_mps_image_sprite_templates" "image_sprite_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  image_sprite_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		resolution_adaptive = &lt;nil&gt;
		sample_type = &lt;nil&gt;
		sample_interval = &lt;nil&gt;
		row_count = &lt;nil&gt;
		column_count = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		fill_type = &lt;nil&gt;
		comment = &lt;nil&gt;
		format = &lt;nil&gt;

  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMpsImageSpriteTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsImageSpriteTemplatesRead,
		Schema: map[string]*schema.Schema{
			"definitions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The image sprite template uniquely identifies the filter condition, and the array length limit: 100.",
			},

			"offset": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Page offset, default: 0.",
			},

			"limit": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Return the number of records, default value: 10, maximum value: 100.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template type filter condition, optional value:Preset: system preset template.Custom: user-defined template.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of records matching filter condition.",
			},

			"image_sprite_template_set": {
				Type:        schema.TypeList,
				Description: "Image sprite template details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Description: "The unique identifier of the image sprite template.",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Template type, optional value:Preset: system preset template.Custom: user-defined template.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Image sprite template name.",
						},
						"width": {
							Type:        schema.TypeInt,
							Description: "The maximum value of the width (or long side) of the small image in the sprite image, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
						},
						"height": {
							Type:        schema.TypeInt,
							Description: "The maximum value of the height (or short side) of the small image in the sprite image, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
						},
						"resolution_adaptive": {
							Type:        schema.TypeString,
							Description: "Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.",
						},
						"sample_type": {
							Type:        schema.TypeString,
							Description: "Sampling type.",
						},
						"sample_interval": {
							Type:        schema.TypeInt,
							Description: "Sampling interval.",
						},
						"row_count": {
							Type:        schema.TypeInt,
							Description: "The number of rows in the small image in the sprite.",
						},
						"column_count": {
							Type:        schema.TypeInt,
							Description: "The number of columns in the small image in the sprite.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Template creation time, in [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Template last modified time, using [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"fill_type": {
							Type:        schema.TypeString,
							Description: "Filling type, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling type:stretch: Stretching, stretching each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the video aspect ratio unchanged, and fill the rest of the edge with black.Default value: black.",
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "The description information of image sprite template.",
						},
						"format": {
							Type:        schema.TypeString,
							Description: "Image format.",
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

func dataSourceTencentCloudMpsImageSpriteTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_image_sprite_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("definitions"); ok {
		definitionsSet := v.(*schema.Set).List()
		for i := range definitionsSet {
			definitions := definitionsSet[i].(int)
			paramMap["Definitions"] = append(paramMap["Definitions"], helper.IntUint64(definitions))
		}
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("image_sprite_template_set"); ok {
		imageSpriteTemplateSetSet := v.([]interface{})
		tmpSet := make([]*mps.ImageSpriteTemplate, 0, len(imageSpriteTemplateSetSet))

		for _, item := range imageSpriteTemplateSetSet {
			imageSpriteTemplate := mps.ImageSpriteTemplate{}
			imageSpriteTemplateMap := item.(map[string]interface{})

			if v, ok := imageSpriteTemplateMap["definition"]; ok {
				imageSpriteTemplate.Definition = helper.IntUint64(v.(int))
			}
			if v, ok := imageSpriteTemplateMap["type"]; ok {
				imageSpriteTemplate.Type = helper.String(v.(string))
			}
			if v, ok := imageSpriteTemplateMap["name"]; ok {
				imageSpriteTemplate.Name = helper.String(v.(string))
			}
			if v, ok := imageSpriteTemplateMap["width"]; ok {
				imageSpriteTemplate.Width = helper.IntUint64(v.(int))
			}
			if v, ok := imageSpriteTemplateMap["height"]; ok {
				imageSpriteTemplate.Height = helper.IntUint64(v.(int))
			}
			if v, ok := imageSpriteTemplateMap["resolution_adaptive"]; ok {
				imageSpriteTemplate.ResolutionAdaptive = helper.String(v.(string))
			}
			if v, ok := imageSpriteTemplateMap["sample_type"]; ok {
				imageSpriteTemplate.SampleType = helper.String(v.(string))
			}
			if v, ok := imageSpriteTemplateMap["sample_interval"]; ok {
				imageSpriteTemplate.SampleInterval = helper.IntUint64(v.(int))
			}
			if v, ok := imageSpriteTemplateMap["row_count"]; ok {
				imageSpriteTemplate.RowCount = helper.IntUint64(v.(int))
			}
			if v, ok := imageSpriteTemplateMap["column_count"]; ok {
				imageSpriteTemplate.ColumnCount = helper.IntUint64(v.(int))
			}
			if v, ok := imageSpriteTemplateMap["create_time"]; ok {
				imageSpriteTemplate.CreateTime = helper.String(v.(string))
			}
			if v, ok := imageSpriteTemplateMap["update_time"]; ok {
				imageSpriteTemplate.UpdateTime = helper.String(v.(string))
			}
			if v, ok := imageSpriteTemplateMap["fill_type"]; ok {
				imageSpriteTemplate.FillType = helper.String(v.(string))
			}
			if v, ok := imageSpriteTemplateMap["comment"]; ok {
				imageSpriteTemplate.Comment = helper.String(v.(string))
			}
			if v, ok := imageSpriteTemplateMap["format"]; ok {
				imageSpriteTemplate.Format = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &imageSpriteTemplate)
		}
		paramMap["image_sprite_template_set"] = tmpSet
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var imageSpriteTemplateSet []*mps.ImageSpriteTemplate

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsImageSpriteTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		imageSpriteTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(imageSpriteTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(imageSpriteTemplateSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
