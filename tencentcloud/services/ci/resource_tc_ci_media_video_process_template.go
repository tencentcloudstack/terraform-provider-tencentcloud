package ci

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCiMediaVideoProcessTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaVideoProcessTemplateCreate,
		Read:   resourceTencentCloudCiMediaVideoProcessTemplateRead,
		Update: resourceTencentCloudCiMediaVideoProcessTemplateUpdate,
		Delete: resourceTencentCloudCiMediaVideoProcessTemplateDelete,
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

			"color_enhance": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "color enhancement.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether color enhancement is turned on.",
						},
						"contrast": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Contrast, value range: [0, 100], empty string (indicates automatic analysis).",
						},
						"correction": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "colorcorrection, value range: [0, 100], empty string (indicating automatic analysis).",
						},
						"saturation": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Saturation, value range: [0, 100], empty string (indicating automatic analysis).",
						},
					},
				},
			},

			"ms_sharpen": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "detail enhancement, ColorEnhance and MsSharpen cannot both be empty.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether detail enhancement is enabled.",
						},
						"sharpen_level": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Enhancement level, value range: [0, 10], empty string (indicates automatic analysis).",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaVideoProcessTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_video_process_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaVideoProcessTemplateOptions{
			Tag: "VideoProcess",
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

	if dMap, ok := helper.InterfacesHeadMap(d, "color_enhance"); ok {
		colorEnhance := cos.ColorEnhance{}
		if v, ok := dMap["enable"]; ok {
			colorEnhance.Enable = v.(string)
		}
		if v, ok := dMap["contrast"]; ok {
			colorEnhance.Contrast = v.(string)
		}
		if v, ok := dMap["correction"]; ok {
			colorEnhance.Correction = v.(string)
		}
		if v, ok := dMap["saturation"]; ok {
			colorEnhance.Saturation = v.(string)
		}
		request.ColorEnhance = &colorEnhance
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ms_sharpen"); ok {
		msSharpen := cos.MsSharpen{}
		if v, ok := dMap["enable"]; ok {
			msSharpen.Enable = v.(string)
		}
		if v, ok := dMap["sharpen_level"]; ok {
			msSharpen.SharpenLevel = v.(string)
		}
		request.MsSharpen = &msSharpen
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaVideoProcessTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaVideoProcessTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaVideoProcessTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaVideoProcessTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVideoProcessTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_video_process_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	template, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if template == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)

	if template.Name != "" {
		_ = d.Set("name", template.Name)
	}

	if template.VideoProcess != nil {
		mediaVideoProcessTemplate := template.VideoProcess
		if mediaVideoProcessTemplate.ColorEnhance != nil {
			colorEnhanceMap := map[string]interface{}{}

			if mediaVideoProcessTemplate.ColorEnhance.Enable != "" {
				colorEnhanceMap["enable"] = mediaVideoProcessTemplate.ColorEnhance.Enable
			}

			if mediaVideoProcessTemplate.ColorEnhance.Contrast != "" {
				colorEnhanceMap["contrast"] = mediaVideoProcessTemplate.ColorEnhance.Contrast
			}

			if mediaVideoProcessTemplate.ColorEnhance.Correction != "" {
				colorEnhanceMap["correction"] = mediaVideoProcessTemplate.ColorEnhance.Correction
			}

			if mediaVideoProcessTemplate.ColorEnhance.Saturation != "" {
				colorEnhanceMap["saturation"] = mediaVideoProcessTemplate.ColorEnhance.Saturation
			}

			_ = d.Set("color_enhance", []interface{}{colorEnhanceMap})
		}

		if mediaVideoProcessTemplate.MsSharpen != nil {
			msSharpenMap := map[string]interface{}{}

			if mediaVideoProcessTemplate.MsSharpen.Enable != "" {
				msSharpenMap["enable"] = mediaVideoProcessTemplate.MsSharpen.Enable
			}

			if mediaVideoProcessTemplate.MsSharpen.SharpenLevel != "" {
				msSharpenMap["sharpen_level"] = mediaVideoProcessTemplate.MsSharpen.SharpenLevel
			}

			_ = d.Set("ms_sharpen", []interface{}{msSharpenMap})
		}
	}

	return nil
}

func resourceTencentCloudCiMediaVideoProcessTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_video_process_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaVideoProcessTemplateOptions{
		Tag: "VideoProcess",
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "color_enhance"); ok {
		colorEnhance := cos.ColorEnhance{}
		if v, ok := dMap["enable"]; ok {
			colorEnhance.Enable = v.(string)
		}
		if v, ok := dMap["contrast"]; ok {
			colorEnhance.Contrast = v.(string)
		}
		if v, ok := dMap["correction"]; ok {
			colorEnhance.Correction = v.(string)
		}
		if v, ok := dMap["saturation"]; ok {
			colorEnhance.Saturation = v.(string)
		}
		request.ColorEnhance = &colorEnhance
	}
	if dMap, ok := helper.InterfacesHeadMap(d, "ms_sharpen"); ok {
		msSharpen := cos.MsSharpen{}
		if v, ok := dMap["enable"]; ok {
			msSharpen.Enable = v.(string)
		}
		if v, ok := dMap["sharpen_level"]; ok {
			msSharpen.SharpenLevel = v.(string)
		}
		request.MsSharpen = &msSharpen
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaVideoProcessTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaVideoProcessTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaVideoProcessTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaVideoProcessTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVideoProcessTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_video_process_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
