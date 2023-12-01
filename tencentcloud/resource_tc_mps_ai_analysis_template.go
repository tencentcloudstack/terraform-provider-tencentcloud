/*
Provides a resource to create a mps ai_analysis_template

Example Usage

```hcl
resource "tencentcloud_mps_ai_analysis_template" "ai_analysis_template" {
  name = "terraform-test"

  classification_configure {
    switch = "OFF"
  }

  cover_configure {
    switch = "ON"
  }

  frame_tag_configure {
    switch = "ON"
  }

  tag_configure {
    switch = "ON"
  }
}

```

Import

mps ai_analysis_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_ai_analysis_template.ai_analysis_template ai_analysis_template_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMpsAiAnalysisTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsAiAnalysisTemplateCreate,
		Read:   resourceTencentCloudMpsAiAnalysisTemplateRead,
		Update: resourceTencentCloudMpsAiAnalysisTemplateUpdate,
		Delete: resourceTencentCloudMpsAiAnalysisTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Ai analysis template name, length limit: 64 characters.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Ai analysis template description information, length limit: 256 characters.",
			},

			"classification_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ai classification task control parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ai classification task switch, optional value:ON/OFF.",
						},
					},
				},
			},

			"tag_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ai tag task control parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ai tag task switch, optional value:ON/OFF.",
						},
					},
				},
			},

			"cover_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ai cover task control parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ai cover task switch, optional value:ON/OFF.",
						},
					},
				},
			},

			"frame_tag_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ai frame tag task control parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ai frame tag task switch, optional value:ON/OFF.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsAiAnalysisTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_ai_analysis_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mps.NewCreateAIAnalysisTemplateRequest()
		response   = mps.NewCreateAIAnalysisTemplateResponse()
		definition int64
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "classification_configure"); ok {
		classificationConfigureInfo := mps.ClassificationConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			classificationConfigureInfo.Switch = helper.String(v.(string))
		}
		request.ClassificationConfigure = &classificationConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "tag_configure"); ok {
		tagConfigureInfo := mps.TagConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			tagConfigureInfo.Switch = helper.String(v.(string))
		}
		request.TagConfigure = &tagConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "cover_configure"); ok {
		coverConfigureInfo := mps.CoverConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			coverConfigureInfo.Switch = helper.String(v.(string))
		}
		request.CoverConfigure = &coverConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "frame_tag_configure"); ok {
		frameTagConfigureInfo := mps.FrameTagConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			frameTagConfigureInfo.Switch = helper.String(v.(string))
		}
		request.FrameTagConfigure = &frameTagConfigureInfo
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreateAIAnalysisTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps aiAnalysisTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.Int64ToStr(definition))

	return resourceTencentCloudMpsAiAnalysisTemplateRead(d, meta)
}

func resourceTencentCloudMpsAiAnalysisTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_ai_analysis_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	definition := d.Id()

	aiAnalysisTemplate, err := service.DescribeMpsAiAnalysisTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if aiAnalysisTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsAiAnalysisTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if aiAnalysisTemplate.Name != nil {
		_ = d.Set("name", aiAnalysisTemplate.Name)
	}

	if aiAnalysisTemplate.Comment != nil {
		_ = d.Set("comment", aiAnalysisTemplate.Comment)
	}

	if aiAnalysisTemplate.ClassificationConfigure != nil {
		classificationConfigureMap := map[string]interface{}{}

		if aiAnalysisTemplate.ClassificationConfigure.Switch != nil {
			classificationConfigureMap["switch"] = aiAnalysisTemplate.ClassificationConfigure.Switch
		}

		_ = d.Set("classification_configure", []interface{}{classificationConfigureMap})
	}

	if aiAnalysisTemplate.TagConfigure != nil {
		tagConfigureMap := map[string]interface{}{}

		if aiAnalysisTemplate.TagConfigure.Switch != nil {
			tagConfigureMap["switch"] = aiAnalysisTemplate.TagConfigure.Switch
		}

		_ = d.Set("tag_configure", []interface{}{tagConfigureMap})
	}

	if aiAnalysisTemplate.CoverConfigure != nil {
		coverConfigureMap := map[string]interface{}{}

		if aiAnalysisTemplate.CoverConfigure.Switch != nil {
			coverConfigureMap["switch"] = aiAnalysisTemplate.CoverConfigure.Switch
		}

		_ = d.Set("cover_configure", []interface{}{coverConfigureMap})
	}

	if aiAnalysisTemplate.FrameTagConfigure != nil {
		frameTagConfigureMap := map[string]interface{}{}

		if aiAnalysisTemplate.FrameTagConfigure.Switch != nil {
			frameTagConfigureMap["switch"] = aiAnalysisTemplate.FrameTagConfigure.Switch
		}

		_ = d.Set("frame_tag_configure", []interface{}{frameTagConfigureMap})
	}

	return nil
}

func resourceTencentCloudMpsAiAnalysisTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_ai_analysis_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyAIAnalysisTemplateRequest()

	definition := d.Id()
	needChange := false

	request.Definition = helper.StrToInt64Point(definition)

	mutableArgs := []string{"name", "comment", "classification_configure", "tag_configure", "cover_configure", "frame_tag_configure"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "classification_configure"); ok {
			classificationConfigureInfo := mps.ClassificationConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				classificationConfigureInfo.Switch = helper.String(v.(string))
			}
			request.ClassificationConfigure = &classificationConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "tag_configure"); ok {
			tagConfigureInfo := mps.TagConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				tagConfigureInfo.Switch = helper.String(v.(string))
			}
			request.TagConfigure = &tagConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "cover_configure"); ok {
			coverConfigureInfo := mps.CoverConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				coverConfigureInfo.Switch = helper.String(v.(string))
			}
			request.CoverConfigure = &coverConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "frame_tag_configure"); ok {
			frameTagConfigureInfo := mps.FrameTagConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				frameTagConfigureInfo.Switch = helper.String(v.(string))
			}
			request.FrameTagConfigure = &frameTagConfigureInfo
		}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyAIAnalysisTemplate(request)
			if e != nil {
				return retryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps aiAnalysisTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsAiAnalysisTemplateRead(d, meta)
}

func resourceTencentCloudMpsAiAnalysisTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_ai_analysis_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	definition := d.Id()

	if err := service.DeleteMpsAiAnalysisTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
