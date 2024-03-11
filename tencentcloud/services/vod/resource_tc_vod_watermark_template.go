package vod

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVodWatermarkTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVodWatermarkTemplateCreate,
		Read:   resourceTencentCloudVodWatermarkTemplateRead,
		Update: resourceTencentCloudVodWatermarkTemplateUpdate,
		Delete: resourceTencentCloudVodWatermarkTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Watermarking type. Valid values: image: image watermark; text: text watermark; svg: SVG watermark.",
			},

			"sub_app_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.",
			},

			"name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Watermarking template name. Length limit: 64 characters.",
			},

			"comment": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Template description. Length limit: 256 characters.",
			},

			"coordinate_origin": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Origin position. Valid values: TopLeft: the origin of coordinates is in the top-left corner of the video, and the origin of the watermark is in the top-left corner of the image or text; TopRight: the origin of coordinates is in the top-right corner of the video, and the origin of the watermark is in the top-right corner of the image or text; BottomLeft: the origin of coordinates is in the bottom-left corner of the video, and the origin of the watermark is in the bottom-left corner of the image or text; BottomRight: the origin of coordinates is in the bottom-right corner of the video, and the origin of the watermark is in the bottom-right corner of the image or text.Default value: TopLeft.",
			},

			"x_pos": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The horizontal position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported: If the string ends in %, the `XPos` of the watermark will be the specified percentage of the video width; for example, `10%` means that `XPos` is 10% of the video width; If the string ends in px, the `XPos` of the watermark will be the specified px; for example, `100px` means that `XPos` is 100 px.Default value: 0 px.",
			},

			"y_pos": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The vertical position of the origin of the watermark relative to the origin of coordinates of the video. % and px formats are supported: If the string ends in %, the `YPos` of the watermark will be the specified percentage of the video height; for example, `10%` means that `YPos` is 10% of the video height; If the string ends in px, the `YPos` of the watermark will be the specified px; for example, `100px` means that `YPos` is 100 px.Default value: 0 px.",
			},

			"image_template": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Image watermarking template. This field is required when `Type` is `image` and is invalid when `Type` is `text`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image_content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The [Base64](https://tools.ietf.org/html/rfc4648) encoded string of a watermark image. Only JPEG, PNG, and GIF images are supported.",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Watermark width. % and px formats are supported: If the string ends in %, the `Width` of the watermark will be the specified percentage of the video width. For example, `10%` means that `Width` is 10% of the video width;  If the string ends in px, the `Width` of the watermark will be in pixels. For example, `100px` means that `Width` is 100 pixels. Value range: [8, 4096]. Default value: 10%.",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Watermark height. % and px formats are supported: If the string ends in %, the `Height` of the watermark will be the specified percentage of the video height; for example, `10%` means that `Height` is 10% of the video height;  If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px. Valid values: 0 or [8,4096]. Default value: 0 px, which means that `Height` will be proportionally scaled according to the aspect ratio of the original watermark image.",
						},
						"repeat_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Repeat type of an animated watermark. Valid values: once: no longer appears after watermark playback ends.  repeat_last_frame: stays on the last frame after watermark playback ends.  repeat (default): repeats the playback until the video ends.",
						},
						"transparency": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Image watermark transparency: 0: completely opaque  100: completely transparent Default value: 0.",
						},
					},
				},
			},

			"text_template": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Text watermarking template. This field is required when `Type` is `text` and is invalid when `Type` is `image`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"font_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font type. Currently, two types are supported: simkai.ttf: both Chinese and English are supported;  arial.ttf: only English is supported.",
						},
						"font_size": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font size in Npx format where N is a numeric value.",
						},
						"font_color": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Font color in 0xRRGGBB format. Default value: 0xFFFFFF (white).",
						},
						"font_alpha": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "Text transparency. Value range: (0, 1] 0: completely transparent  1: completely opaque Default value: 1.",
						},
					},
				},
			},

			"svg_template": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "SVG watermarking template. This field is required when `Type` is `svg` and is invalid when `Type` is `image` or `text`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Watermark width, which supports six formats of px, %, W%, H%, S%, and L%: If the string ends in px, the `Width` of the watermark will be in px; for example, `100px` means that `Width` is 100 px; if `0px` is entered and `Height` is not `0px`, the watermark width will be proportionally scaled based on the source SVG image; if `0px` is entered for both `Width` and `Height`, the watermark width will be the width of the source SVG image;  If the string ends in `W%`, the `Width` of the watermark will be the specified percentage of the video width; for example, `10W%` means that `Width` is 10% of the video width;  If the string ends in `H%`, the `Width` of the watermark will be the specified percentage of the video height; for example, `10H%` means that `Width` is 10% of the video height;  If the string ends in `S%`, the `Width` of the watermark will be the specified percentage of the short side of the video; for example, `10S%` means that `Width` is 10% of the short side of the video;  If the string ends in `L%`, the `Width` of the watermark will be the specified percentage of the long side of the video; for example, `10L%` means that `Width` is 10% of the long side of the video;  If the string ends in %, the meaning is the same as `W%`. Default value: 10W%.",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Watermark height, which supports six formats of px, %, W%, H%, S%, and L%: If the string ends in px, the `Height` of the watermark will be in px; for example, `100px` means that `Height` is 100 px; if `0px` is entered and `Width` is not `0px`, the watermark height will be proportionally scaled based on the source SVG image; if `0px` is entered for both `Width` and `Height`, the watermark height will be the height of the source SVG image;  If the string ends in `W%`, the `Height` of the watermark will be the specified percentage of the video width; for example, `10W%` means that `Height` is 10% of the video width;  If the string ends in `H%`, the `Height` of the watermark will be the specified percentage of the video height; for example, `10H%` means that `Height` is 10% of the video height;  If the string ends in `S%`, the `Height` of the watermark will be the specified percentage of the short side of the video; for example, `10S%` means that `Height` is 10% of the short side of the video;  If the string ends in `L%`, the `Height` of the watermark will be the specified percentage of the long side of the video; for example, `10L%` means that `Height` is 10% of the long side of the video;  If the string ends in %, the meaning is the same as `H%`. Default value: 0 px.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVodWatermarkTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_watermark_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = vod.NewCreateWatermarkTemplateRequest()
		response = vod.NewCreateWatermarkTemplateResponse()
		subAppId string
	)
	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("sub_app_id"); ok {
		subAppId = helper.IntToStr(v.(int))
		request.SubAppId = helper.IntUint64(v.(int))
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
		imageWatermarkInput := vod.ImageWatermarkInput{}
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
		if v, ok := dMap["transparency"]; ok {
			imageWatermarkInput.Transparency = helper.IntInt64(v.(int))
		}
		request.ImageTemplate = &imageWatermarkInput
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "text_template"); ok {
		textWatermarkTemplateInput := vod.TextWatermarkTemplateInput{}
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
		svgWatermarkInput := vod.SvgWatermarkInput{}
		if v, ok := dMap["width"]; ok {
			svgWatermarkInput.Width = helper.String(v.(string))
		}
		if v, ok := dMap["height"]; ok {
			svgWatermarkInput.Height = helper.String(v.(string))
		}
		request.SvgTemplate = &svgWatermarkInput
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().CreateWatermarkTemplate(request)
		if e != nil {
			if sdkError, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkError.Code == "FailedOperation" && sdkError.Message == "invalid vod user" {
					return resource.RetryableError(e)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return resource.NonRetryableError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create vod watermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition := *response.Response.Definition
	d.SetId(subAppId + tccommon.FILED_SP + helper.Int64ToStr(definition))

	return resourceTencentCloudVodWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudVodWatermarkTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_watermark_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("watermark template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	watermarkTemplate, err := service.DescribeVodWatermarkTemplateById(ctx, *helper.StrToUint64Point(subAppId), helper.StrToInt64(definition))
	if err != nil {
		return err
	}

	if watermarkTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `VodWatermarkTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("sub_app_id", helper.StrToInt(subAppId))

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
			imageContentResp, err := http.Get(*watermarkTemplate.ImageTemplate.ImageUrl)
			if err != nil {
				return err
			}
			content, err := ioutil.ReadAll(imageContentResp.Body)
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

		if watermarkTemplate.ImageTemplate.Transparency != nil {
			imageTemplateMap["transparency"] = watermarkTemplate.ImageTemplate.Transparency
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

func resourceTencentCloudVodWatermarkTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_watermark_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewModifyWatermarkTemplateRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("watermark template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	request.SubAppId = helper.StrToUint64Point(subAppId)
	request.Definition = helper.StrToInt64Point(definition)

	immutableArgs := []string{"type", "sub_app_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

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
			imageWatermarkInput := vod.ImageWatermarkInputForUpdate{}
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
			if v, ok := dMap["transparency"]; ok {
				imageWatermarkInput.Transparency = helper.IntInt64(v.(int))
			}
			request.ImageTemplate = &imageWatermarkInput
		}
	}

	if d.HasChange("text_template") {
		if dMap, ok := helper.InterfacesHeadMap(d, "text_template"); ok {
			textWatermarkTemplateInput := vod.TextWatermarkTemplateInputForUpdate{}
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
			svgWatermarkInput := vod.SvgWatermarkInputForUpdate{}
			if v, ok := dMap["width"]; ok {
				svgWatermarkInput.Width = helper.String(v.(string))
			}
			if v, ok := dMap["height"]; ok {
				svgWatermarkInput.Height = helper.String(v.(string))
			}
			request.SvgTemplate = &svgWatermarkInput
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVodClient().ModifyWatermarkTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update vod watermarkTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudVodWatermarkTemplateRead(d, meta)
}

func resourceTencentCloudVodWatermarkTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vod_watermark_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("watermark template id is borken, id is %s", d.Id())
	}
	subAppId := idSplit[0]
	definition := idSplit[1]

	if err := service.DeleteVodWatermarkTemplateById(ctx, helper.StrToUInt64(subAppId), helper.StrToInt64(definition)); err != nil {
		return err
	}

	return nil
}
