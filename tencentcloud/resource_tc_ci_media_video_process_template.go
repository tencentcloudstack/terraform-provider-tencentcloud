/*
Provides a resource to create a ci media_video_process_template

Example Usage

```hcl
resource "tencentcloud_ci_media_video_process_template" "media_video_process_template" {
  name = &lt;nil&gt;
  color_enhance {
		enable = &lt;nil&gt;
		contrast = &lt;nil&gt;
		correction = &lt;nil&gt;
		saturation = &lt;nil&gt;

  }
  ms_sharpen {
		enable = &lt;nil&gt;
		sharpen_level = &lt;nil&gt;

  }
}
```

Import

ci media_video_process_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_video_process_template.media_video_process_template media_video_process_template_id
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

func resourceTencentCloudCiMediaVideoProcessTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaVideoProcessTemplateCreate,
		Read:   resourceTencentCloudCiMediaVideoProcessTemplateRead,
		Update: resourceTencentCloudCiMediaVideoProcessTemplateUpdate,
		Delete: resourceTencentCloudCiMediaVideoProcessTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"color_enhance": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Color enhancement.",
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
							Description: "Colorcorrection, value range: [0, 100], empty string (indicating automatic analysis).",
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
				Description: "Detail enhancement, ColorEnhance and MsSharpen cannot both be empty.",
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
	defer logElapsed("resource.tencentcloud_ci_media_video_process_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ci.NewCreateMediaVideoProcessTemplateRequest()
		response   = ci.NewCreateMediaVideoProcessTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "color_enhance"); ok {
		colorEnhance := ci.ColorEnhance{}
		if v, ok := dMap["enable"]; ok {
			colorEnhance.Enable = helper.String(v.(string))
		}
		if v, ok := dMap["contrast"]; ok {
			colorEnhance.Contrast = helper.String(v.(string))
		}
		if v, ok := dMap["correction"]; ok {
			colorEnhance.Correction = helper.String(v.(string))
		}
		if v, ok := dMap["saturation"]; ok {
			colorEnhance.Saturation = helper.String(v.(string))
		}
		request.ColorEnhance = &colorEnhance
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ms_sharpen"); ok {
		colorEnhance := ci.ColorEnhance{}
		if v, ok := dMap["enable"]; ok {
			colorEnhance.Enable = helper.String(v.(string))
		}
		if v, ok := dMap["sharpen_level"]; ok {
			colorEnhance.SharpenLevel = helper.String(v.(string))
		}
		request.MsSharpen = &colorEnhance
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaVideoProcessTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaVideoProcessTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaVideoProcessTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVideoProcessTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_video_process_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaVideoProcessTemplateId := d.Id()

	mediaVideoProcessTemplate, err := service.DescribeCiMediaVideoProcessTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaVideoProcessTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaVideoProcessTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaVideoProcessTemplate.Name != nil {
		_ = d.Set("name", mediaVideoProcessTemplate.Name)
	}

	if mediaVideoProcessTemplate.ColorEnhance != nil {
		colorEnhanceMap := map[string]interface{}{}

		if mediaVideoProcessTemplate.ColorEnhance.Enable != nil {
			colorEnhanceMap["enable"] = mediaVideoProcessTemplate.ColorEnhance.Enable
		}

		if mediaVideoProcessTemplate.ColorEnhance.Contrast != nil {
			colorEnhanceMap["contrast"] = mediaVideoProcessTemplate.ColorEnhance.Contrast
		}

		if mediaVideoProcessTemplate.ColorEnhance.Correction != nil {
			colorEnhanceMap["correction"] = mediaVideoProcessTemplate.ColorEnhance.Correction
		}

		if mediaVideoProcessTemplate.ColorEnhance.Saturation != nil {
			colorEnhanceMap["saturation"] = mediaVideoProcessTemplate.ColorEnhance.Saturation
		}

		_ = d.Set("color_enhance", []interface{}{colorEnhanceMap})
	}

	if mediaVideoProcessTemplate.MsSharpen != nil {
		msSharpenMap := map[string]interface{}{}

		if mediaVideoProcessTemplate.MsSharpen.Enable != nil {
			msSharpenMap["enable"] = mediaVideoProcessTemplate.MsSharpen.Enable
		}

		if mediaVideoProcessTemplate.MsSharpen.SharpenLevel != nil {
			msSharpenMap["sharpen_level"] = mediaVideoProcessTemplate.MsSharpen.SharpenLevel
		}

		_ = d.Set("ms_sharpen", []interface{}{msSharpenMap})
	}

	return nil
}

func resourceTencentCloudCiMediaVideoProcessTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_video_process_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaVideoProcessTemplateRequest()

	mediaVideoProcessTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "color_enhance", "ms_sharpen"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaVideoProcessTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaVideoProcessTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaVideoProcessTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaVideoProcessTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_video_process_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	mediaVideoProcessTemplateId := d.Id()

	if err := service.DeleteCiMediaVideoProcessTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
