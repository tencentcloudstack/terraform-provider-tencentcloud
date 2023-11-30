package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCiMediaWatermarkTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaWatermarkTemplateCreate,
		Read:   resourceTencentCloudCiMediaWatermarkTemplateRead,
		Update: resourceTencentCloudCiMediaWatermarkTemplateUpdate,
		Delete: resourceTencentCloudCiMediaWatermarkTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "bucket name.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"watermark": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "container format.",
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
							Description: "Horizontal offset, 1: In the picture watermark, if Background is true, when locMode is Relativity, it is %, value range: [-300 0]; when locMode is Absolute, it is px, value range: [-4096 0] ], 2: In the picture watermark, if Background is false, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096], 3: In text watermark, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096], 4: When Pos is Top, Bottom and Center, the parameter is invalid.",
						},
						"dy": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Vertical offset, 1: In the picture watermark, if Background is true, when locMode is Relativity, it is %, value range: [-300 0]; when locMode is Absolute, it is px, value range: [-4096 0] ],2: In the picture watermark, if Background is false, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096],3: In text watermark, when locMode is Relativity, it is %, value range: [0 100]; when locMode is Absolute, it is px, value range: [0 4096], 4: When Pos is Left, Right and Center, the parameter is invalid.",
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
										Description: "font type.",
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

func resourceTencentCloudCiMediaWatermarkTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_watermark_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = cos.CreateMediaWatermarkTemplateOptions{
			Tag: "Watermark",
		}
		bucket     string
		templateId string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "watermark"); ok {
		watermark := cos.Watermark{}
		if v, ok := dMap["type"]; ok {
			watermark.Type = v.(string)
		}
		if v, ok := dMap["pos"]; ok {
			watermark.Pos = v.(string)
		}
		if v, ok := dMap["loc_mode"]; ok {
			watermark.LocMode = v.(string)
		}
		if v, ok := dMap["dx"]; ok {
			watermark.Dx = v.(string)
		}
		if v, ok := dMap["dy"]; ok {
			watermark.Dy = v.(string)
		}
		if v, ok := dMap["start_time"]; ok {
			watermark.StartTime = v.(string)
		}
		if v, ok := dMap["end_time"]; ok {
			watermark.EndTime = v.(string)
		}
		if imageMap, ok := helper.InterfaceToMap(dMap, "image"); ok {
			image := cos.Image{}
			if v, ok := imageMap["url"]; ok {
				image.Url = v.(string)
			}
			if v, ok := imageMap["mode"]; ok {
				image.Mode = v.(string)
			}
			if v, ok := imageMap["width"]; ok {
				image.Width = v.(string)
			}
			if v, ok := imageMap["height"]; ok {
				image.Height = v.(string)
			}
			if v, ok := imageMap["transparency"]; ok {
				image.Transparency = v.(string)
			}
			if v, ok := imageMap["background"]; ok {
				image.Background = v.(string)
			}
			watermark.Image = &image
		}
		if textMap, ok := helper.InterfaceToMap(dMap, "text"); ok {
			text := cos.Text{}
			if v, ok := textMap["font_size"]; ok {
				text.FontSize = v.(string)
			}
			if v, ok := textMap["font_type"]; ok {
				text.FontType = v.(string)
			}
			if v, ok := textMap["font_color"]; ok {
				text.FontColor = v.(string)
			}
			if v, ok := textMap["transparency"]; ok {
				text.Transparency = v.(string)
			}
			if v, ok := textMap["text"]; ok {
				text.Text = v.(string)
			}
			watermark.Text = &text
		}
		request.Watermark = &watermark
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.CreateMediaWatermarkTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaWatermarkTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaWatermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaWatermarkTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_watermark_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	mediaWatermarkTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaWatermarkTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)

	if mediaWatermarkTemplate.Name != "" {
		_ = d.Set("name", mediaWatermarkTemplate.Name)
	}

	if mediaWatermarkTemplate.Watermark != nil {
		watermarkMap := map[string]interface{}{}

		if mediaWatermarkTemplate.Watermark.Type != "" {
			watermarkMap["type"] = mediaWatermarkTemplate.Watermark.Type
		}

		if mediaWatermarkTemplate.Watermark.Pos != "" {
			watermarkMap["pos"] = mediaWatermarkTemplate.Watermark.Pos
		}

		if mediaWatermarkTemplate.Watermark.LocMode != "" {
			watermarkMap["loc_mode"] = mediaWatermarkTemplate.Watermark.LocMode
		}

		if mediaWatermarkTemplate.Watermark.Dx != "" {
			watermarkMap["dx"] = mediaWatermarkTemplate.Watermark.Dx
		}

		if mediaWatermarkTemplate.Watermark.Dy != "" {
			watermarkMap["dy"] = mediaWatermarkTemplate.Watermark.Dy
		}

		if mediaWatermarkTemplate.Watermark.StartTime != "" {
			watermarkMap["start_time"] = mediaWatermarkTemplate.Watermark.StartTime
		}

		if mediaWatermarkTemplate.Watermark.EndTime != "" {
			watermarkMap["end_time"] = mediaWatermarkTemplate.Watermark.EndTime
		}

		if mediaWatermarkTemplate.Watermark.Image != nil {
			imageMap := map[string]interface{}{}

			if mediaWatermarkTemplate.Watermark.Image.Url != "" {
				imageMap["url"] = mediaWatermarkTemplate.Watermark.Image.Url
			}

			if mediaWatermarkTemplate.Watermark.Image.Mode != "" {
				imageMap["mode"] = mediaWatermarkTemplate.Watermark.Image.Mode
			}

			if mediaWatermarkTemplate.Watermark.Image.Width != "" {
				imageMap["width"] = mediaWatermarkTemplate.Watermark.Image.Width
			}

			if mediaWatermarkTemplate.Watermark.Image.Height != "" {
				imageMap["height"] = mediaWatermarkTemplate.Watermark.Image.Height
			}

			if mediaWatermarkTemplate.Watermark.Image.Transparency != "" {
				imageMap["transparency"] = mediaWatermarkTemplate.Watermark.Image.Transparency
			}

			if mediaWatermarkTemplate.Watermark.Image.Background != "" {
				imageMap["background"] = mediaWatermarkTemplate.Watermark.Image.Background
			}

			watermarkMap["image"] = []interface{}{imageMap}
		}

		if mediaWatermarkTemplate.Watermark.Text != nil {
			textMap := map[string]interface{}{}

			if mediaWatermarkTemplate.Watermark.Text.FontSize != "" {
				textMap["font_size"] = mediaWatermarkTemplate.Watermark.Text.FontSize
			}

			if mediaWatermarkTemplate.Watermark.Text.FontType != "" {
				textMap["font_type"] = mediaWatermarkTemplate.Watermark.Text.FontType
			}

			if mediaWatermarkTemplate.Watermark.Text.FontColor != "" {
				textMap["font_color"] = mediaWatermarkTemplate.Watermark.Text.FontColor
			}

			if mediaWatermarkTemplate.Watermark.Text.Transparency != "" {
				textMap["transparency"] = mediaWatermarkTemplate.Watermark.Text.Transparency
			}

			if mediaWatermarkTemplate.Watermark.Text.Text != "" {
				textMap["text"] = mediaWatermarkTemplate.Watermark.Text.Text
			}

			watermarkMap["text"] = []interface{}{textMap}
		}

		_ = d.Set("watermark", []interface{}{watermarkMap})
	}

	return nil
}

func resourceTencentCloudCiMediaWatermarkTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_watermark_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cos.CreateMediaWatermarkTemplateOptions{
		Tag: "Watermark",
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if d.HasChange("watermark") {
		if dMap, ok := helper.InterfacesHeadMap(d, "watermark"); ok {
			watermark := cos.Watermark{}
			if v, ok := dMap["type"]; ok {
				watermark.Type = v.(string)
			}
			if v, ok := dMap["pos"]; ok {
				watermark.Pos = v.(string)
			}
			if v, ok := dMap["loc_mode"]; ok {
				watermark.LocMode = v.(string)
			}
			if v, ok := dMap["dx"]; ok {
				watermark.Dx = v.(string)
			}
			if v, ok := dMap["dy"]; ok {
				watermark.Dy = v.(string)
			}
			if v, ok := dMap["start_time"]; ok {
				watermark.StartTime = v.(string)
			}
			if v, ok := dMap["end_time"]; ok {
				watermark.EndTime = v.(string)
			}

			if imageMap, ok := helper.InterfaceToMap(dMap, "image"); ok {
				image := cos.Image{}
				if v, ok := imageMap["url"]; ok {
					image.Url = v.(string)
				}
				if v, ok := imageMap["mode"]; ok {
					image.Mode = v.(string)
				}
				if v, ok := imageMap["width"]; ok {
					image.Width = v.(string)
				}
				if v, ok := imageMap["height"]; ok {
					image.Height = v.(string)
				}
				if v, ok := imageMap["transparency"]; ok {
					image.Transparency = v.(string)
				}
				if v, ok := imageMap["background"]; ok {
					image.Background = v.(string)
				}
				watermark.Image = &image
			}
			if textMap, ok := helper.InterfaceToMap(dMap, "text"); ok {
				text := cos.Text{}
				if v, ok := textMap["font_size"]; ok {
					text.FontSize = v.(string)
				}
				if v, ok := textMap["font_type"]; ok {
					text.FontType = v.(string)
				}
				if v, ok := textMap["font_color"]; ok {
					text.FontColor = v.(string)
				}
				if v, ok := textMap["transparency"]; ok {
					text.Transparency = v.(string)
				}
				if v, ok := textMap["text"]; ok {
					text.Text = v.(string)
				}
				watermark.Text = &text
			}
			request.Watermark = &watermark
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaWatermarkTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaWatermarkTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaWatermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaWatermarkTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_watermark_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if err := service.DeleteCiMediaTemplateById(ctx, bucket, templateId); err != nil {
		return err
	}

	return nil
}
