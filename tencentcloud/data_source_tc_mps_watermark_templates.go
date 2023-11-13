/*
Use this data source to query detailed information of mps watermark_templates

Example Usage

```hcl
data "tencentcloud_mps_watermark_templates" "watermark_templates" {
  definitions = &lt;nil&gt;
  type = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  total_count = &lt;nil&gt;
  watermark_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		x_pos = &lt;nil&gt;
		y_pos = &lt;nil&gt;
		image_template {
			image_url = &lt;nil&gt;
			width = &lt;nil&gt;
			height = &lt;nil&gt;
			repeat_type = &lt;nil&gt;
		}
		text_template {
			font_type = &lt;nil&gt;
			font_size = &lt;nil&gt;
			font_color = &lt;nil&gt;
			font_alpha = &lt;nil&gt;
		}
		svg_template {
			width = &lt;nil&gt;
			height = &lt;nil&gt;
		}
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		coordinate_origin = &lt;nil&gt;

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

func dataSourceTencentCloudMpsWatermarkTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsWatermarkTemplatesRead,
		Schema: map[string]*schema.Schema{
			"definitions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Watermark template uniquely identifies filter conditions, array length limit: 100.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Watermark type filter condition, optional value:image/text.",
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

			"total_count": {
				Type:        schema.TypeInt,
				Description: "Total number of records matching filter condition.",
			},

			"watermark_template_set": {
				Type:        schema.TypeList,
				Description: "Watermark template details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Description: "The unique identifier of the watermark template.",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Watermark type, optional value:image/text.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Watermark template name.",
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "Template description information.",
						},
						"x_pos": {
							Type:        schema.TypeString,
							Description: "The horizontal position from the origin of the watermark image to the origin of the video image.When the string ends with %, it means that the watermark Left is the specified percentage of the video width, such as 10% means that the Left is 10% of the video width.When the string ends with px, it means that the watermark Left is the pixel position of the video width, such as 100px means that the Left is 100 pixels.",
						},
						"y_pos": {
							Type:        schema.TypeString,
							Description: "The vertical position between the origin of the watermark image and the origin of the video image.When the string ends with %, it means that the watermark Top is the specified percentage of the video height, such as 10% means that Top is 10% of the video height.When the string ends with px, it means that the watermark Top specifies the pixel position of the video height, such as 100px means that the Top is 100 pixels.",
						},
						"image_template": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Image watermark template, only when Type is image, this field is valid.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"image_url": {
										Type:        schema.TypeString,
										Description: "Watermark image address.",
									},
									"width": {
										Type:        schema.TypeString,
										Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixel, such as 100px means that the Width is 100 pixels.",
									},
									"height": {
										Type:        schema.TypeString,
										Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels.0px: Indicates that Height is scaled according to the aspect ratio of the original watermark image.",
									},
									"repeat_type": {
										Type:        schema.TypeString,
										Description: "Watermark repeat type. Usage scenario: The watermark is a dynamic image. Ranges:once: After the dynamic watermark is played, it will no longer appear.repeat_last_frame: After the watermark is played, stay on the last frame.repeat: the watermark loops until the end of the video (default).",
									},
								},
							},
						},
						"text_template": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Text watermark template, only when Type is text, this field is valid.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"font_type": {
										Type:        schema.TypeString,
										Description: "Font type, currently supports two:simkai.ttf: can support Chinese and English.arial.ttf: English only.",
									},
									"font_size": {
										Type:        schema.TypeString,
										Description: "Font size, format: Npx, N is a number.",
									},
									"font_color": {
										Type:        schema.TypeString,
										Description: "Font color, format: 0xRRGGBB, default value: 0xFFFFFF (white).",
									},
									"font_alpha": {
										Type:        schema.TypeFloat,
										Description: "Text transparency, value range: (0, 1].0: fully transparent.1: fully opaque.Default value: 1.",
									},
								},
							},
						},
						"svg_template": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "SVG watermark template, only when Type is svg, this field is valid.Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"width": {
										Type:        schema.TypeString,
										Description: "The width of the watermark, supports px, %, W%, H%, S%, L% six formats.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels; when filling 0px and the Height is not 0px, it means that the width of the watermark is proportionally scaled according to the original SVG image; when both Width and Height are filled When 0px, it means that the width of the watermark takes the width of the original SVG image.When the string ends with W%, it means that the watermark Width is a percentage of the video width, such as 10W% means that the Width is 10% of the video width.When the string ends with H%, it means that the watermark Width is a percentage of the video height, such as 10H% means that the Width is 10% of the video height.When the string ends with S%, it means that the watermark Width is the percentage size of the short side of the video, such as 10S% means that the Width is 10% of the short side of the video.When the string ends with L%, it means that the watermark Width is the percentage size of the long side of the video, such as 10L% means that the Width is 10% of the long side of the video.When the string ends with %, it has the same meaning as W%.Default value: 10W%.",
									},
									"height": {
										Type:        schema.TypeString,
										Description: "The height of the watermark, supports px, W%, H%, S%, L% six formats:When the string ends with px, it means that the watermark Height unit is pixels, such as 100px means that the Height is 100 pixels; when filling 0px and Width is not 0px, it means that the height of the watermark is proportionally scaled according to the original SVG image; when both Width and Height are filled When 0px, it means that the height of the watermark takes the height of the original SVG image.When the string ends with W%, it means that the watermark Height is a percentage of the video width, such as 10W% means that the Height is 10% of the video width.When the string ends with H%, it means that the watermark Height is the percentage size of the video height, such as 10H% means that the Height is 10% of the video height.When the string ends with S%, it means that the watermark Height is the percentage size of the short side of the video, such as 10S% means that the Height is 10% of the short side of the video.When the string ends with L%, it means that the watermark Height is the percentage size of the long side of the video, such as 10L% means that the Height is 10% of the long side of the video.When the string ends with %, the meaning is the same as H%.Default value: 0px.",
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Template creation time, in [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Template last modified time, using [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"coordinate_origin": {
							Type:        schema.TypeString,
							Description: "Origin position, optional value:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.TopRight: Indicates that the origin of the coordinates is at the upper right corner of the video image, and the origin of the watermark is at the upper right corner of the picture or text.BottomLeft: Indicates that the origin of the coordinates is at the lower left corner of the video image, and the origin of the watermark is the lower left corner of the picture or text.BottomRight: Indicates that the origin of the coordinates is at the lower right corner of the video image, and the origin of the watermark is at the lower right corner of the picture or text.",
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

func dataSourceTencentCloudMpsWatermarkTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_watermark_templates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("definitions"); ok {
		definitionsSet := v.(*schema.Set).List()
		for i := range definitionsSet {
			definitions := definitionsSet[i].(int)
			paramMap["Definitions"] = append(paramMap["Definitions"], helper.IntInt64(definitions))
		}
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("offset"); v != nil {
		paramMap["Offset"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("limit"); v != nil {
		paramMap["Limit"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("total_count"); v != nil {
		paramMap["TotalCount"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("watermark_template_set"); ok {
		watermarkTemplateSetSet := v.([]interface{})
		tmpSet := make([]*mps.WatermarkTemplate, 0, len(watermarkTemplateSetSet))

		for _, item := range watermarkTemplateSetSet {
			watermarkTemplate := mps.WatermarkTemplate{}
			watermarkTemplateMap := item.(map[string]interface{})

			if v, ok := watermarkTemplateMap["definition"]; ok {
				watermarkTemplate.Definition = helper.IntInt64(v.(int))
			}
			if v, ok := watermarkTemplateMap["type"]; ok {
				watermarkTemplate.Type = helper.String(v.(string))
			}
			if v, ok := watermarkTemplateMap["name"]; ok {
				watermarkTemplate.Name = helper.String(v.(string))
			}
			if v, ok := watermarkTemplateMap["comment"]; ok {
				watermarkTemplate.Comment = helper.String(v.(string))
			}
			if v, ok := watermarkTemplateMap["x_pos"]; ok {
				watermarkTemplate.XPos = helper.String(v.(string))
			}
			if v, ok := watermarkTemplateMap["y_pos"]; ok {
				watermarkTemplate.YPos = helper.String(v.(string))
			}
			if imageTemplateMap, ok := helper.InterfaceToMap(watermarkTemplateMap, "image_template"); ok {
				imageWatermarkTemplate := mps.ImageWatermarkTemplate{}
				if v, ok := imageTemplateMap["image_url"]; ok {
					imageWatermarkTemplate.ImageUrl = helper.String(v.(string))
				}
				if v, ok := imageTemplateMap["width"]; ok {
					imageWatermarkTemplate.Width = helper.String(v.(string))
				}
				if v, ok := imageTemplateMap["height"]; ok {
					imageWatermarkTemplate.Height = helper.String(v.(string))
				}
				if v, ok := imageTemplateMap["repeat_type"]; ok {
					imageWatermarkTemplate.RepeatType = helper.String(v.(string))
				}
				watermarkTemplate.ImageTemplate = &imageWatermarkTemplate
			}
			if textTemplateMap, ok := helper.InterfaceToMap(watermarkTemplateMap, "text_template"); ok {
				textWatermarkTemplateInput := mps.TextWatermarkTemplateInput{}
				if v, ok := textTemplateMap["font_type"]; ok {
					textWatermarkTemplateInput.FontType = helper.String(v.(string))
				}
				if v, ok := textTemplateMap["font_size"]; ok {
					textWatermarkTemplateInput.FontSize = helper.String(v.(string))
				}
				if v, ok := textTemplateMap["font_color"]; ok {
					textWatermarkTemplateInput.FontColor = helper.String(v.(string))
				}
				if v, ok := textTemplateMap["font_alpha"]; ok {
					textWatermarkTemplateInput.FontAlpha = helper.Float64(v.(float64))
				}
				watermarkTemplate.TextTemplate = &textWatermarkTemplateInput
			}
			if svgTemplateMap, ok := helper.InterfaceToMap(watermarkTemplateMap, "svg_template"); ok {
				svgWatermarkInput := mps.SvgWatermarkInput{}
				if v, ok := svgTemplateMap["width"]; ok {
					svgWatermarkInput.Width = helper.String(v.(string))
				}
				if v, ok := svgTemplateMap["height"]; ok {
					svgWatermarkInput.Height = helper.String(v.(string))
				}
				watermarkTemplate.SvgTemplate = &svgWatermarkInput
			}
			if v, ok := watermarkTemplateMap["create_time"]; ok {
				watermarkTemplate.CreateTime = helper.String(v.(string))
			}
			if v, ok := watermarkTemplateMap["update_time"]; ok {
				watermarkTemplate.UpdateTime = helper.String(v.(string))
			}
			if v, ok := watermarkTemplateMap["coordinate_origin"]; ok {
				watermarkTemplate.CoordinateOrigin = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &watermarkTemplate)
		}
		paramMap["watermark_template_set"] = tmpSet
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var watermarkTemplateSet []*mps.WatermarkTemplate

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsWatermarkTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		watermarkTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(watermarkTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(watermarkTemplateSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
