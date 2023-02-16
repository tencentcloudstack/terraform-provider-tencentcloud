/*
Provides a resource to create a mps watermark_template

Example Usage

```hcl
resource "tencentcloud_mps_watermark_template" "watermark_template" {
  coordinate_origin = "TopLeft"
  name              = "xZxasd"
  type              = "image"
  x_pos             = "12%"
  y_pos             = "21%"

  image_template {
    height        = "17px"
    image_content = filebase64("./logo.png")
    repeat_type   = "repeat"
    width         = "12px"
  }
}
```

Import

mps watermark_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_watermark_template.watermark_template watermark_template_id
```
*/
package tencentcloud

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsWatermarkTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsWatermarkTemplateCreate,
		Read:   resourceTencentCloudMpsWatermarkTemplateRead,
		Update: resourceTencentCloudMpsWatermarkTemplateUpdate,
		Delete: resourceTencentCloudMpsWatermarkTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Watermark type, optional value:image, text, svg.",
			},

			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Watermark template name, length limit: 64 characters.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Template description information, length limit: 256 characters.",
			},

			"coordinate_origin": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Origin position, optional value:TopLeft: Indicates that the origin of the coordinates is at the upper left corner of the video image, and the origin of the watermark is the upper left corner of the picture or text.TopRight: Indicates that the origin of the coordinates is at the upper right corner of the video image, and the origin of the watermark is at the upper right corner of the picture or text.BottomLeft: Indicates that the origin of the coordinates is at the lower left corner of the video image, and the origin of the watermark is the lower left corner of the picture or text.BottomRight: Indicates that the origin of the coordinates is at the lower right corner of the video image, and the origin of the watermark is at the lower right corner of the picture or text.Default value: TopLeft.",
			},

			"x_pos": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The horizontal position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats.When the string ends with %, it means that the watermark XPos specifies a percentage for the video width, such as 10% means that XPos is 10% of the video width.When the string ends with px, it means that the watermark XPos is the specified pixel, such as 100px means that the XPos is 100 pixels.Default value: 0px.",
			},

			"y_pos": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The vertical position of the origin of the watermark from the origin of the coordinates of the video image. Support %, px two formats.When the string ends with %, it means that the watermark YPos specifies a percentage for the video height, such as 10% means that YPos is 10% of the video height.When the string ends with px, it means that the watermark YPos is the specified pixel, such as 100px means that the YPos is 100 pixels.Default value: 0px.",
			},

			"image_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Image watermark template, only when Type is image, this field is required and valid.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Watermark image[Base64](https://tools.ietf.org/html/rfc4648) encoded string. Support jpeg, png image format.",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The width of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Width is a percentage of the video width, such as 10% means that the Width is 10% of the video width.When the string ends with px, it means that the watermark Width unit is pixel, such as 100px means that the Width is 100 pixels. The value range is [8, 4096].Default value: 10%.",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The height of the watermark. Support %, px two formats:When the string ends with %, it means that the watermark Height is the percentage size of the video height, such as 10% means that the Height is 10% of the video height.When the string ends with px, it means that the watermark Height unit is pixel, such as 100px means that the Height is 100 pixels. The value range is 0 or [8, 4096].Default value: 0px. Indicates that Height is scaled according to the aspect ratio of the original watermark image.",
						},
						"repeat_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Watermark repeat type. Usage scenario: The watermark is a dynamic image. Ranges:once: After the dynamic watermark is played, it will no longer appear.repeat_last_frame: After the watermark is played, stay on the last frame.repeat: the watermark loops until the end of the video (default).",
						},
					},
				},
			},

			"text_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Text watermark template, only when Type is text, this field is required and valid.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"font_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font type, currently supports two:simkai.ttf: can support Chinese and English.arial.ttf: English only.",
						},
						"font_size": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font size, format: Npx, N is a number.",
						},
						"font_color": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font color, format: 0xRRGGBB, default value: 0xFFFFFF (white).",
						},
						"font_alpha": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "Text transparency, value range: (0, 1].0: fully transparent.1: fully opaque.Default value: 1.",
						},
					},
				},
			},

			"svg_template": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "SVG watermark template, only when Type is svg, this field is required and valid.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The width of the watermark, supports px, %, W%, H%, S%, L% six formats.When the string ends with px, it means that the watermark Width unit is pixels, such as 100px means that the Width is 100 pixels; when filling 0px and the Height is not 0px, it means that the width of the watermark is proportionally scaled according to the original SVG image; when both Width and Height are filled When 0px, it means that the width of the watermark takes the width of the original SVG image.When the string ends with W%, it means that the watermark Width is a percentage of the video width, such as 10W% means that the Width is 10% of the video width.When the string ends with H%, it means that the watermark Width is a percentage of the video height, such as 10H% means that the Width is 10% of the video height.When the string ends with S%, it means that the watermark Width is the percentage size of the short side of the video, such as 10S% means that the Width is 10% of the short side of the video.When the string ends with L%, it means that the watermark Width is the percentage size of the long side of the video, such as 10L% means that the Width is 10% of the long side of the video.When the string ends with %, it has the same meaning as W%.Default value: 10W%.",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The height of the watermark, supports px, W%, H%, S%, L% six formats:When the string ends with px, it means that the watermark Height unit is pixels, such as 100px means that the Height is 100 pixels; when filling 0px and Width is not 0px, it means that the height of the watermark is proportionally scaled according to the original SVG image; when both Width and Height are filled When 0px, it means that the height of the watermark takes the height of the original SVG image.When the string ends with W%, it means that the watermark Height is a percentage of the video width, such as 10W% means that the Height is 10% of the video width.When the string ends with H%, it means that the watermark Height is the percentage size of the video height, such as 10H% means that the Height is 10% of the video height.When the string ends with S%, it means that the watermark Height is the percentage size of the short side of the video, such as 10S% means that the Height is 10% of the short side of the video.When the string ends with L%, it means that the watermark Height is the percentage size of the long side of the video, such as 10L% means that the Height is 10% of the long side of the video.When the string ends with %, the meaning is the same as H%.Default value: 0px.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsWatermarkTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_watermark_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mps.NewCreateWatermarkTemplateRequest()
		response   = mps.NewCreateWatermarkTemplateResponse()
		definition int64
	)
	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if v, ok := d.GetOk("coordinate_origin"); ok {
		request.CoordinateOrigin = helper.String(v.(string))
	}

	if v, ok := d.GetOk("x_pos"); ok {
		request.XPos = helper.String(v.(string))
	}

	if v, ok := d.GetOk("y_pos"); ok {
		request.YPos = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "image_template"); ok {
		imageWatermarkInput := mps.ImageWatermarkInput{}
		if v, ok := dMap["image_content"]; ok {
			imageWatermarkInput.ImageContent = helper.String(v.(string))
		}
		if v, ok := dMap["width"]; ok {
			imageWatermarkInput.Width = helper.String(v.(string))
		}
		if v, ok := dMap["height"]; ok {
			imageWatermarkInput.Height = helper.String(v.(string))
		}
		if v, ok := dMap["repeat_type"]; ok {
			imageWatermarkInput.RepeatType = helper.String(v.(string))
		}
		request.ImageTemplate = &imageWatermarkInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "text_template"); ok {
		textWatermarkTemplateInput := mps.TextWatermarkTemplateInput{}
		if v, ok := dMap["font_type"]; ok {
			textWatermarkTemplateInput.FontType = helper.String(v.(string))
		}
		if v, ok := dMap["font_size"]; ok {
			textWatermarkTemplateInput.FontSize = helper.String(v.(string))
		}
		if v, ok := dMap["font_color"]; ok {
			textWatermarkTemplateInput.FontColor = helper.String(v.(string))
		}
		if v, ok := dMap["font_alpha"]; ok {
			textWatermarkTemplateInput.FontAlpha = helper.Float64(v.(float64))
		}
		request.TextTemplate = &textWatermarkTemplateInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "svg_template"); ok {
		svgWatermarkInput := mps.SvgWatermarkInput{}
		if v, ok := dMap["width"]; ok {
			svgWatermarkInput.Width = helper.String(v.(string))
		}
		if v, ok := dMap["height"]; ok {
			svgWatermarkInput.Height = helper.String(v.(string))
		}
		request.SvgTemplate = &svgWatermarkInput
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateWatermarkTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps watermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.Int64ToStr(definition))

	return resourceTencentCloudMpsWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudMpsWatermarkTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_watermark_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	definition := d.Id()

	watermarkTemplate, err := service.DescribeMpsWatermarkTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if watermarkTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsWatermarkTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if watermarkTemplate.Type != nil {
		_ = d.Set("type", watermarkTemplate.Type)
	}

	if watermarkTemplate.Name != nil {
		_ = d.Set("name", watermarkTemplate.Name)
	}

	if watermarkTemplate.Comment != nil {
		_ = d.Set("comment", watermarkTemplate.Comment)
	}

	if watermarkTemplate.CoordinateOrigin != nil {
		_ = d.Set("coordinate_origin", watermarkTemplate.CoordinateOrigin)
	}

	if watermarkTemplate.XPos != nil {
		_ = d.Set("x_pos", watermarkTemplate.XPos)
	}

	if watermarkTemplate.YPos != nil {
		_ = d.Set("y_pos", watermarkTemplate.YPos)
	}

	if watermarkTemplate.ImageTemplate != nil {
		imageTemplateMap := map[string]interface{}{}

		if watermarkTemplate.ImageTemplate.ImageUrl != nil {
			url := watermarkTemplate.ImageTemplate.ImageUrl
			res, err := http.Get(*url)
			if err != nil {
				return err
			}
			content, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			base64Encode := base64.StdEncoding.EncodeToString(content)
			imageTemplateMap["image_content"] = base64Encode
		}

		if watermarkTemplate.ImageTemplate.Width != nil {
			imageTemplateMap["width"] = watermarkTemplate.ImageTemplate.Width
		}

		if watermarkTemplate.ImageTemplate.Height != nil {
			imageTemplateMap["height"] = watermarkTemplate.ImageTemplate.Height
		}

		if watermarkTemplate.ImageTemplate.RepeatType != nil {
			imageTemplateMap["repeat_type"] = watermarkTemplate.ImageTemplate.RepeatType
		}

		_ = d.Set("image_template", []interface{}{imageTemplateMap})
	}

	if watermarkTemplate.TextTemplate != nil {
		textTemplateMap := map[string]interface{}{}

		if watermarkTemplate.TextTemplate.FontType != nil {
			textTemplateMap["font_type"] = watermarkTemplate.TextTemplate.FontType
		}

		if watermarkTemplate.TextTemplate.FontSize != nil {
			textTemplateMap["font_size"] = watermarkTemplate.TextTemplate.FontSize
		}

		if watermarkTemplate.TextTemplate.FontColor != nil {
			textTemplateMap["font_color"] = watermarkTemplate.TextTemplate.FontColor
		}

		if watermarkTemplate.TextTemplate.FontAlpha != nil {
			textTemplateMap["font_alpha"] = watermarkTemplate.TextTemplate.FontAlpha
		}

		_ = d.Set("text_template", []interface{}{textTemplateMap})
	}

	if watermarkTemplate.SvgTemplate != nil {
		svgTemplateMap := map[string]interface{}{}

		if watermarkTemplate.SvgTemplate.Width != nil {
			svgTemplateMap["width"] = watermarkTemplate.SvgTemplate.Width
		}

		if watermarkTemplate.SvgTemplate.Height != nil {
			svgTemplateMap["height"] = watermarkTemplate.SvgTemplate.Height
		}

		_ = d.Set("svg_template", []interface{}{svgTemplateMap})
	}

	return nil
}

func resourceTencentCloudMpsWatermarkTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_watermark_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyWatermarkTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToInt64Point(definition)

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("comment") {
		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}
	}

	if d.HasChange("coordinate_origin") {
		if v, ok := d.GetOk("coordinate_origin"); ok {
			request.CoordinateOrigin = helper.String(v.(string))
		}
	}

	if d.HasChange("x_pos") {
		if v, ok := d.GetOk("x_pos"); ok {
			request.XPos = helper.String(v.(string))
		}
	}

	if d.HasChange("y_pos") {
		if v, ok := d.GetOk("y_pos"); ok {
			request.YPos = helper.String(v.(string))
		}
	}

	if d.HasChange("image_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "image_template"); ok {
			imageWatermarkInput := mps.ImageWatermarkInputForUpdate{}
			if v, ok := dMap["image_content"]; ok {
				imageWatermarkInput.ImageContent = helper.String(v.(string))
			}
			if v, ok := dMap["width"]; ok {
				imageWatermarkInput.Width = helper.String(v.(string))
			}
			if v, ok := dMap["height"]; ok {
				imageWatermarkInput.Height = helper.String(v.(string))
			}
			if v, ok := dMap["repeat_type"]; ok {
				imageWatermarkInput.RepeatType = helper.String(v.(string))
			}
			request.ImageTemplate = &imageWatermarkInput
		}
	}

	if d.HasChange("text_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "text_template"); ok {
			textWatermarkTemplateInput := mps.TextWatermarkTemplateInputForUpdate{}
			if v, ok := dMap["font_type"]; ok {
				textWatermarkTemplateInput.FontType = helper.String(v.(string))
			}
			if v, ok := dMap["font_size"]; ok {
				textWatermarkTemplateInput.FontSize = helper.String(v.(string))
			}
			if v, ok := dMap["font_color"]; ok {
				textWatermarkTemplateInput.FontColor = helper.String(v.(string))
			}
			if v, ok := dMap["font_alpha"]; ok {
				textWatermarkTemplateInput.FontAlpha = helper.Float64(v.(float64))
			}
			request.TextTemplate = &textWatermarkTemplateInput
		}
	}

	if d.HasChange("svg_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "svg_template"); ok {
			svgWatermarkInput := mps.SvgWatermarkInputForUpdate{}
			if v, ok := dMap["width"]; ok {
				svgWatermarkInput.Width = helper.String(v.(string))
			}
			if v, ok := dMap["height"]; ok {
				svgWatermarkInput.Height = helper.String(v.(string))
			}
			request.SvgTemplate = &svgWatermarkInput
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyWatermarkTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps watermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudMpsWatermarkTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_watermark_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	definition := d.Id()

	if err := service.DeleteMpsWatermarkTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
