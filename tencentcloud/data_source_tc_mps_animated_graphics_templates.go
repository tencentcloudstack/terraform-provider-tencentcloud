/*
Use this data source to query detailed information of mps animated_graphics_templates

Example Usage

```hcl
data "tencentcloud_mps_animated_graphics_templates" "animated_graphics_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  animated_graphics_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		resolution_adaptive = &lt;nil&gt;
		format = &lt;nil&gt;
		fps = &lt;nil&gt;
		quality = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;

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

func dataSourceTencentCloudMpsAnimatedGraphicsTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMpsAnimatedGraphicsTemplatesRead,
		Schema: map[string]*schema.Schema{
			"definitions": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Animated graphics template uniquely identifies filter conditions, array length limit: 100.",
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

			"animated_graphics_template_set": {
				Type:        schema.TypeList,
				Description: "Animated graphics template details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"definition": {
							Type:        schema.TypeInt,
							Description: "The unique identifier of the animated graphics template.",
						},
						"type": {
							Type:        schema.TypeString,
							Description: "Template type, optional value:Preset: system preset template.Custom: user-defined template.",
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Animated graphics template name.",
						},
						"comment": {
							Type:        schema.TypeString,
							Description: "The description information of animated graphics template.",
						},
						"width": {
							Type:        schema.TypeInt,
							Description: "The maximum value of the animation width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
						},
						"height": {
							Type:        schema.TypeInt,
							Description: "The maximum value of the animation height (or short side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.",
						},
						"resolution_adaptive": {
							Type:        schema.TypeString,
							Description: "Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.",
						},
						"format": {
							Type:        schema.TypeString,
							Description: "Animation format.",
						},
						"fps": {
							Type:        schema.TypeInt,
							Description: "Frame rate.",
						},
						"quality": {
							Type:        schema.TypeFloat,
							Description: "Image quality.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Description: "Template creation time, in [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
						},
						"update_time": {
							Type:        schema.TypeString,
							Description: "Template last modified time, using [ISO date format](https://cloud.tencent.com/document/product/862/37710#52).",
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

func dataSourceTencentCloudMpsAnimatedGraphicsTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mps_animated_graphics_templates.read")()
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

	if v, ok := d.GetOk("animated_graphics_template_set"); ok {
		animatedGraphicsTemplateSetSet := v.([]interface{})
		tmpSet := make([]*mps.AnimatedGraphicsTemplate, 0, len(animatedGraphicsTemplateSetSet))

		for _, item := range animatedGraphicsTemplateSetSet {
			animatedGraphicsTemplate := mps.AnimatedGraphicsTemplate{}
			animatedGraphicsTemplateMap := item.(map[string]interface{})

			if v, ok := animatedGraphicsTemplateMap["definition"]; ok {
				animatedGraphicsTemplate.Definition = helper.IntUint64(v.(int))
			}
			if v, ok := animatedGraphicsTemplateMap["type"]; ok {
				animatedGraphicsTemplate.Type = helper.String(v.(string))
			}
			if v, ok := animatedGraphicsTemplateMap["name"]; ok {
				animatedGraphicsTemplate.Name = helper.String(v.(string))
			}
			if v, ok := animatedGraphicsTemplateMap["comment"]; ok {
				animatedGraphicsTemplate.Comment = helper.String(v.(string))
			}
			if v, ok := animatedGraphicsTemplateMap["width"]; ok {
				animatedGraphicsTemplate.Width = helper.IntUint64(v.(int))
			}
			if v, ok := animatedGraphicsTemplateMap["height"]; ok {
				animatedGraphicsTemplate.Height = helper.IntUint64(v.(int))
			}
			if v, ok := animatedGraphicsTemplateMap["resolution_adaptive"]; ok {
				animatedGraphicsTemplate.ResolutionAdaptive = helper.String(v.(string))
			}
			if v, ok := animatedGraphicsTemplateMap["format"]; ok {
				animatedGraphicsTemplate.Format = helper.String(v.(string))
			}
			if v, ok := animatedGraphicsTemplateMap["fps"]; ok {
				animatedGraphicsTemplate.Fps = helper.IntUint64(v.(int))
			}
			if v, ok := animatedGraphicsTemplateMap["quality"]; ok {
				animatedGraphicsTemplate.Quality = helper.Float64(v.(float64))
			}
			if v, ok := animatedGraphicsTemplateMap["create_time"]; ok {
				animatedGraphicsTemplate.CreateTime = helper.String(v.(string))
			}
			if v, ok := animatedGraphicsTemplateMap["update_time"]; ok {
				animatedGraphicsTemplate.UpdateTime = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &animatedGraphicsTemplate)
		}
		paramMap["animated_graphics_template_set"] = tmpSet
	}

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var animatedGraphicsTemplateSet []*mps.AnimatedGraphicsTemplate

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMpsAnimatedGraphicsTemplatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		animatedGraphicsTemplateSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(animatedGraphicsTemplateSet))
	tmpList := make([]map[string]interface{}, 0, len(animatedGraphicsTemplateSet))

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
