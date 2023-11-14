/*
Provides a resource to create a ci media_smart_cover_template

Example Usage

```hcl
resource "tencentcloud_ci_media_smart_cover_template" "media_smart_cover_template" {
  name = &lt;nil&gt;
  smart_cover {
		format = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		count = &lt;nil&gt;
		delete_duplicates = &lt;nil&gt;

  }
}
```

Import

ci media_smart_cover_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_smart_cover_template.media_smart_cover_template media_smart_cover_template_id
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

func resourceTencentCloudCiMediaSmartCoverTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaSmartCoverTemplateCreate,
		Read:   resourceTencentCloudCiMediaSmartCoverTemplateRead,
		Update: resourceTencentCloudCiMediaSmartCoverTemplateUpdate,
		Delete: resourceTencentCloudCiMediaSmartCoverTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"smart_cover": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Smart Cover Parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Image Format, value jpg、png 、webp.",
						},
						"width": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Width, value range: [128, 4096], unit: px, if only Width is set, Height is calculated according to the original ratio of the video.",
						},
						"height": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Height, value range: [128, 4096], unit: px, if only Height is set, Width is calculated according to the original video ratio.",
						},
						"count": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Number of screenshots, [1,10].",
						},
						"delete_duplicates": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cover deduplication, true/false.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaSmartCoverTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_smart_cover_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ci.NewCreateMediaSmartCoverTemplateRequest()
		response   = ci.NewCreateMediaSmartCoverTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "smart_cover"); ok {
		smartCover := ci.SmartCover{}
		if v, ok := dMap["format"]; ok {
			smartCover.Format = helper.String(v.(string))
		}
		if v, ok := dMap["width"]; ok {
			smartCover.Width = helper.String(v.(string))
		}
		if v, ok := dMap["height"]; ok {
			smartCover.Height = helper.String(v.(string))
		}
		if v, ok := dMap["count"]; ok {
			smartCover.Count = helper.String(v.(string))
		}
		if v, ok := dMap["delete_duplicates"]; ok {
			smartCover.DeleteDuplicates = helper.String(v.(string))
		}
		request.SmartCover = &smartCover
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaSmartCoverTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSmartCoverTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaSmartCoverTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSmartCoverTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_smart_cover_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaSmartCoverTemplateId := d.Id()

	mediaSmartCoverTemplate, err := service.DescribeCiMediaSmartCoverTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaSmartCoverTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaSmartCoverTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaSmartCoverTemplate.Name != nil {
		_ = d.Set("name", mediaSmartCoverTemplate.Name)
	}

	if mediaSmartCoverTemplate.SmartCover != nil {
		smartCoverMap := map[string]interface{}{}

		if mediaSmartCoverTemplate.SmartCover.Format != nil {
			smartCoverMap["format"] = mediaSmartCoverTemplate.SmartCover.Format
		}

		if mediaSmartCoverTemplate.SmartCover.Width != nil {
			smartCoverMap["width"] = mediaSmartCoverTemplate.SmartCover.Width
		}

		if mediaSmartCoverTemplate.SmartCover.Height != nil {
			smartCoverMap["height"] = mediaSmartCoverTemplate.SmartCover.Height
		}

		if mediaSmartCoverTemplate.SmartCover.Count != nil {
			smartCoverMap["count"] = mediaSmartCoverTemplate.SmartCover.Count
		}

		if mediaSmartCoverTemplate.SmartCover.DeleteDuplicates != nil {
			smartCoverMap["delete_duplicates"] = mediaSmartCoverTemplate.SmartCover.DeleteDuplicates
		}

		_ = d.Set("smart_cover", []interface{}{smartCoverMap})
	}

	return nil
}

func resourceTencentCloudCiMediaSmartCoverTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_smart_cover_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaSmartCoverTemplateRequest()

	mediaSmartCoverTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "smart_cover"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaSmartCoverTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaSmartCoverTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaSmartCoverTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSmartCoverTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_smart_cover_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	mediaSmartCoverTemplateId := d.Id()

	if err := service.DeleteCiMediaSmartCoverTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
