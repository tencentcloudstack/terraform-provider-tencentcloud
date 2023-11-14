/*
Provides a resource to create a ci media_pic_process_template

Example Usage

```hcl
resource "tencentcloud_ci_media_pic_process_template" "media_pic_process_template" {
  name = &lt;nil&gt;
  watermark {
		type = &lt;nil&gt;
		pos = &lt;nil&gt;
		loc_mode = &lt;nil&gt;
		dx = &lt;nil&gt;
		dy = &lt;nil&gt;
		start_time = &lt;nil&gt;
		end_time = &lt;nil&gt;
		slide_config {
			slide_mode = &lt;nil&gt;
			x_slide_speed = &lt;nil&gt;
			x_slide_speed = &lt;nil&gt;
		}
		image {
			url = &lt;nil&gt;
			mode = &lt;nil&gt;
			width = &lt;nil&gt;
			height = &lt;nil&gt;
			transparency = &lt;nil&gt;
			background = &lt;nil&gt;
		}
		text {
			font_size = &lt;nil&gt;
			font_type = &lt;nil&gt;
			font_color = &lt;nil&gt;
			transparency = &lt;nil&gt;
			text = &lt;nil&gt;
		}

  }
}
```

Import

ci media_pic_process_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_pic_process_template.media_pic_process_template media_pic_process_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	ci "github.com/tencentyun/cos-go-sdk-v5"
	"log"
)

func resourceTencentCloudCiMediaPicProcessTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaPicProcessTemplateCreate,
		Read:   resourceTencentCloudCiMediaPicProcessTemplateRead,
		Update: resourceTencentCloudCiMediaPicProcessTemplateUpdate,
		Delete: resourceTencentCloudCiMediaPicProcessTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"watermark": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Container format.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Watermark type, Text: text watermark, Image: image watermark.",
						},
						"pos": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Reference position, TopRight, TopLeft, BottomRight, BottomLeft, Left, Right, Top, Bottom, Center.",
						},
						"loc_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Offset method, Relativity: proportional, Absolute: fixed position.",
						},
						"dx": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Horizontal offset, 1: In the picture watermark, if Background is true, when locMode is Relativity, it is %, value range: [-300 0]; when locMode is Absolute, it is px, value range: [-4096 0] ],2: In the picture watermark, if Background is false, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096],3: In text watermark, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096], 4: When Pos is Top, Bottom and Center , the parameter is invalid.",
						},
						"dy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vertical offset, 1: In the picture watermark, if Background is true, when locMode is Relativity, it is %, value range: [-300 0]; when locMode is Absolute, it is px, value range: [-4096 0] ],2: In the picture watermark, if Background is false, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096],3: In text watermark, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096], 4: When Pos is Left, Right and Center , the parameter is invalid.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Watermark start time, 1: [0 video duration], 2: unit is second, 3: support float format, execution accuracy is accurate to milliseconds.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Watermark end time, 1: [0 video duration], 2: unit is second, 3: support float format, execution accuracy is accurate to milliseconds.",
						},
						"slide_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Watermark sliding configuration, after configuring this parameter, the watermark displacement setting will not take effect, and the ultra-fast HD/H265 transcoding does not support this parameter for the time being.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"slide_mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Sliding mode, Default: Disabled by default, ScrollFromLeft: Scroll from left to right, if the ScrollFromLeft mode is set, the Watermark.Pos parameter will not take effecte.",
									},
									"x_slide_speed": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Horizontal sliding speed, value range: an integer in [0,10], the default is 0.",
									},
									"x_slide_speed": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Vertical sliding speed, value range: an integer in [0,10], the default is 0.",
									},
								},
							},
						},
						"image": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Image watermark node.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Address of watermark map (pass in after Urlencode is required).",
									},
									"mode": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Size mode, Original: original size, Proportion: proportional, Fixed: fixed size.",
									},
									"width": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Width, 1: When the Mode is Original, it does not support setting the width of the watermark image, 2: When the Mode is Proportion, the unit is %, the value range of the background image: [100 300]; the value range of the foreground image: [1 100], relative to Video width, up to 4096px, 3: When Mode is Fixed, the unit is px, value range: [8, 4096], 4: If only Width is set, Height is calculated according to the proportion of the watermark image.",
									},
									"height": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "High, 1: When the Mode is Original, it does not support setting the width of the watermark image, 2: When the Mode is Proportion, the unit is %, the value range of the background image: [100 300]; the value range of the foreground image: [1 100], relative to Video width, up to 4096px, 3: When Mode is Fixed, the unit is px, value range: [8, 4096], 4: If only Width is set, Height is calculated according to the proportion of the watermark image.",
									},
									"transparency": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Transparency, value range: [1 100], unit %.",
									},
									"background": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Whether the background image.",
									},
								},
							},
						},
						"text": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Text Watermark Node.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"font_size": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Font size, value range: [5 100], unit px.",
									},
									"font_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Font type.",
									},
									"font_color": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Font color, format: 0xRRGGBB.",
									},
									"transparency": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Transparency, value range: [1 100], unit %.",
									},
									"text": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Watermark content, the length does not exceed 64 characters, only supports Chinese, English, numbers, _, - and *.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaPicProcessTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_pic_process_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ci.NewCreateMediaWatermarkTemplateRequest()
		response   = ci.NewCreateMediaWatermarkTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "watermark"); ok {
		watermark := ci.Watermark{}
		if v, ok := dMap["type"]; ok {
			watermark.Type = helper.String(v.(string))
		}
		if v, ok := dMap["pos"]; ok {
			watermark.Pos = helper.String(v.(string))
		}
		if v, ok := dMap["loc_mode"]; ok {
			watermark.LocMode = helper.String(v.(string))
		}
		if v, ok := dMap["dx"]; ok {
			watermark.Dx = helper.String(v.(string))
		}
		if v, ok := dMap["dy"]; ok {
			watermark.Dy = helper.String(v.(string))
		}
		if v, ok := dMap["start_time"]; ok {
			watermark.StartTime = helper.String(v.(string))
		}
		if v, ok := dMap["end_time"]; ok {
			watermark.EndTime = helper.String(v.(string))
		}
		if slideConfigMap, ok := helper.InterfaceToMap(dMap, "slide_config"); ok {
			slideConfig := ci.SlideConfig{}
			if v, ok := slideConfigMap["slide_mode"]; ok {
				slideConfig.SlideMode = helper.String(v.(string))
			}
			if v, ok := slideConfigMap["x_slide_speed"]; ok {
				slideConfig.XSlideSpeed = helper.String(v.(string))
			}
			if v, ok := slideConfigMap["x_slide_speed"]; ok {
				slideConfig.XSlideSpeed = helper.String(v.(string))
			}
			watermark.SlideConfig = &slideConfig
		}
		if imageMap, ok := helper.InterfaceToMap(dMap, "image"); ok {
			image := ci.Image{}
			if v, ok := imageMap["url"]; ok {
				image.Url = helper.String(v.(string))
			}
			if v, ok := imageMap["mode"]; ok {
				image.Mode = helper.String(v.(string))
			}
			if v, ok := imageMap["width"]; ok {
				image.Width = helper.String(v.(string))
			}
			if v, ok := imageMap["height"]; ok {
				image.Height = helper.String(v.(string))
			}
			if v, ok := imageMap["transparency"]; ok {
				image.Transparency = helper.String(v.(string))
			}
			if v, ok := imageMap["background"]; ok {
				image.Background = helper.String(v.(string))
			}
			watermark.Image = &image
		}
		if textMap, ok := helper.InterfaceToMap(dMap, "text"); ok {
			text := ci.Text{}
			if v, ok := textMap["font_size"]; ok {
				text.FontSize = helper.String(v.(string))
			}
			if v, ok := textMap["font_type"]; ok {
				text.FontType = helper.String(v.(string))
			}
			if v, ok := textMap["font_color"]; ok {
				text.FontColor = helper.String(v.(string))
			}
			if v, ok := textMap["transparency"]; ok {
				text.Transparency = helper.String(v.(string))
			}
			if v, ok := textMap["text"]; ok {
				text.Text = helper.String(v.(string))
			}
			watermark.Text = &text
		}
		request.Watermark = &watermark
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaWatermarkTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaPicProcessTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaPicProcessTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaPicProcessTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_pic_process_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaPicProcessTemplateId := d.Id()

	mediaPicProcessTemplate, err := service.DescribeCiMediaPicProcessTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaPicProcessTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaPicProcessTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaPicProcessTemplate.Name != nil {
		_ = d.Set("name", mediaPicProcessTemplate.Name)
	}

	if mediaPicProcessTemplate.Watermark != nil {
		watermarkMap := map[string]interface{}{}

		if mediaPicProcessTemplate.Watermark.Type != nil {
			watermarkMap["type"] = mediaPicProcessTemplate.Watermark.Type
		}

		if mediaPicProcessTemplate.Watermark.Pos != nil {
			watermarkMap["pos"] = mediaPicProcessTemplate.Watermark.Pos
		}

		if mediaPicProcessTemplate.Watermark.LocMode != nil {
			watermarkMap["loc_mode"] = mediaPicProcessTemplate.Watermark.LocMode
		}

		if mediaPicProcessTemplate.Watermark.Dx != nil {
			watermarkMap["dx"] = mediaPicProcessTemplate.Watermark.Dx
		}

		if mediaPicProcessTemplate.Watermark.Dy != nil {
			watermarkMap["dy"] = mediaPicProcessTemplate.Watermark.Dy
		}

		if mediaPicProcessTemplate.Watermark.StartTime != nil {
			watermarkMap["start_time"] = mediaPicProcessTemplate.Watermark.StartTime
		}

		if mediaPicProcessTemplate.Watermark.EndTime != nil {
			watermarkMap["end_time"] = mediaPicProcessTemplate.Watermark.EndTime
		}

		if mediaPicProcessTemplate.Watermark.SlideConfig != nil {
			slideConfigMap := map[string]interface{}{}

			if mediaPicProcessTemplate.Watermark.SlideConfig.SlideMode != nil {
				slideConfigMap["slide_mode"] = mediaPicProcessTemplate.Watermark.SlideConfig.SlideMode
			}

			if mediaPicProcessTemplate.Watermark.SlideConfig.XSlideSpeed != nil {
				slideConfigMap["x_slide_speed"] = mediaPicProcessTemplate.Watermark.SlideConfig.XSlideSpeed
			}

			if mediaPicProcessTemplate.Watermark.SlideConfig.XSlideSpeed != nil {
				slideConfigMap["x_slide_speed"] = mediaPicProcessTemplate.Watermark.SlideConfig.XSlideSpeed
			}

			watermarkMap["slide_config"] = []interface{}{slideConfigMap}
		}

		if mediaPicProcessTemplate.Watermark.Image != nil {
			imageMap := map[string]interface{}{}

			if mediaPicProcessTemplate.Watermark.Image.Url != nil {
				imageMap["url"] = mediaPicProcessTemplate.Watermark.Image.Url
			}

			if mediaPicProcessTemplate.Watermark.Image.Mode != nil {
				imageMap["mode"] = mediaPicProcessTemplate.Watermark.Image.Mode
			}

			if mediaPicProcessTemplate.Watermark.Image.Width != nil {
				imageMap["width"] = mediaPicProcessTemplate.Watermark.Image.Width
			}

			if mediaPicProcessTemplate.Watermark.Image.Height != nil {
				imageMap["height"] = mediaPicProcessTemplate.Watermark.Image.Height
			}

			if mediaPicProcessTemplate.Watermark.Image.Transparency != nil {
				imageMap["transparency"] = mediaPicProcessTemplate.Watermark.Image.Transparency
			}

			if mediaPicProcessTemplate.Watermark.Image.Background != nil {
				imageMap["background"] = mediaPicProcessTemplate.Watermark.Image.Background
			}

			watermarkMap["image"] = []interface{}{imageMap}
		}

		if mediaPicProcessTemplate.Watermark.Text != nil {
			textMap := map[string]interface{}{}

			if mediaPicProcessTemplate.Watermark.Text.FontSize != nil {
				textMap["font_size"] = mediaPicProcessTemplate.Watermark.Text.FontSize
			}

			if mediaPicProcessTemplate.Watermark.Text.FontType != nil {
				textMap["font_type"] = mediaPicProcessTemplate.Watermark.Text.FontType
			}

			if mediaPicProcessTemplate.Watermark.Text.FontColor != nil {
				textMap["font_color"] = mediaPicProcessTemplate.Watermark.Text.FontColor
			}

			if mediaPicProcessTemplate.Watermark.Text.Transparency != nil {
				textMap["transparency"] = mediaPicProcessTemplate.Watermark.Text.Transparency
			}

			if mediaPicProcessTemplate.Watermark.Text.Text != nil {
				textMap["text"] = mediaPicProcessTemplate.Watermark.Text.Text
			}

			watermarkMap["text"] = []interface{}{textMap}
		}

		_ = d.Set("watermark", []interface{}{watermarkMap})
	}

	return nil
}

func resourceTencentCloudCiMediaPicProcessTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_pic_process_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaWatermarkTemplateRequest()

	mediaPicProcessTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "watermark"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("watermark") {
		if dMap, ok := helper.InterfacesHeadMap(d, "watermark"); ok {
			watermark := ci.Watermark{}
			if v, ok := dMap["type"]; ok {
				watermark.Type = helper.String(v.(string))
			}
			if v, ok := dMap["pos"]; ok {
				watermark.Pos = helper.String(v.(string))
			}
			if v, ok := dMap["loc_mode"]; ok {
				watermark.LocMode = helper.String(v.(string))
			}
			if v, ok := dMap["dx"]; ok {
				watermark.Dx = helper.String(v.(string))
			}
			if v, ok := dMap["dy"]; ok {
				watermark.Dy = helper.String(v.(string))
			}
			if v, ok := dMap["start_time"]; ok {
				watermark.StartTime = helper.String(v.(string))
			}
			if v, ok := dMap["end_time"]; ok {
				watermark.EndTime = helper.String(v.(string))
			}
			if slideConfigMap, ok := helper.InterfaceToMap(dMap, "slide_config"); ok {
				slideConfig := ci.SlideConfig{}
				if v, ok := slideConfigMap["slide_mode"]; ok {
					slideConfig.SlideMode = helper.String(v.(string))
				}
				if v, ok := slideConfigMap["x_slide_speed"]; ok {
					slideConfig.XSlideSpeed = helper.String(v.(string))
				}
				if v, ok := slideConfigMap["x_slide_speed"]; ok {
					slideConfig.XSlideSpeed = helper.String(v.(string))
				}
				watermark.SlideConfig = &slideConfig
			}
			if imageMap, ok := helper.InterfaceToMap(dMap, "image"); ok {
				image := ci.Image{}
				if v, ok := imageMap["url"]; ok {
					image.Url = helper.String(v.(string))
				}
				if v, ok := imageMap["mode"]; ok {
					image.Mode = helper.String(v.(string))
				}
				if v, ok := imageMap["width"]; ok {
					image.Width = helper.String(v.(string))
				}
				if v, ok := imageMap["height"]; ok {
					image.Height = helper.String(v.(string))
				}
				if v, ok := imageMap["transparency"]; ok {
					image.Transparency = helper.String(v.(string))
				}
				if v, ok := imageMap["background"]; ok {
					image.Background = helper.String(v.(string))
				}
				watermark.Image = &image
			}
			if textMap, ok := helper.InterfaceToMap(dMap, "text"); ok {
				text := ci.Text{}
				if v, ok := textMap["font_size"]; ok {
					text.FontSize = helper.String(v.(string))
				}
				if v, ok := textMap["font_type"]; ok {
					text.FontType = helper.String(v.(string))
				}
				if v, ok := textMap["font_color"]; ok {
					text.FontColor = helper.String(v.(string))
				}
				if v, ok := textMap["transparency"]; ok {
					text.Transparency = helper.String(v.(string))
				}
				if v, ok := textMap["text"]; ok {
					text.Text = helper.String(v.(string))
				}
				watermark.Text = &text
			}
			request.Watermark = &watermark
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaWatermarkTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaPicProcessTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaPicProcessTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaPicProcessTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_pic_process_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	mediaPicProcessTemplateId := d.Id()

	if err := service.DeleteCiMediaPicProcessTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
