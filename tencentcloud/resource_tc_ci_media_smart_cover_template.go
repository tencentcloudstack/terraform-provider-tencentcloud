/*
Provides a resource to create a ci media_smart_cover_template

Example Usage

```hcl
resource "tencentcloud_ci_media_smart_cover_template" "media_smart_cover_template" {
  bucket = "terraform-ci-xxxxxx"
  name = "smart_cover_template"
  smart_cover {
		format = "jpg"
		width = "1280"
		height = "960"
		count = "10"
		delete_duplicates = "true"
  }
}
```

Import

ci media_smart_cover_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_smart_cover_template.media_smart_cover_template terraform-ci-xxxxxx#t1ede83acc305e423799d638044d859fb7
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentyun/cos-go-sdk-v5"
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
							Description: "Image Format, value jpg, png, webp.",
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
							Description: "cover deduplication, true/false.",
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request = cos.CreateMediaSmartCoverTemplateOptions{
			Tag: "SmartCover",
		}
		templateId string
		bucket     string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "smart_cover"); ok {
		smartCover := cos.NodeSmartCover{}
		if v, ok := dMap["format"]; ok {
			smartCover.Format = v.(string)
		}
		if v, ok := dMap["width"]; ok {
			smartCover.Width = v.(string)
		}
		if v, ok := dMap["height"]; ok {
			smartCover.Height = v.(string)
		}
		if v, ok := dMap["count"]; ok {
			smartCover.Count = v.(string)
		}
		if v, ok := dMap["delete_duplicates"]; ok {
			smartCover.DeleteDuplicates = v.(string)
		}
		request.SmartCover = &smartCover
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.CreateMediaSmartCoverTemplate(ctx, &request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "CreateMediaSmartCoverTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSmartCoverTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + FILED_SP + templateId)

	return resourceTencentCloudCiMediaSmartCoverTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSmartCoverTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_smart_cover_template.read")()
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

	mediaSmartCoverTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaSmartCoverTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	if mediaSmartCoverTemplate.Name != "" {
		_ = d.Set("name", mediaSmartCoverTemplate.Name)
	}

	if mediaSmartCoverTemplate.SmartCover != nil {
		smartCoverMap := map[string]interface{}{}

		if mediaSmartCoverTemplate.SmartCover.Format != "" {
			smartCoverMap["format"] = mediaSmartCoverTemplate.SmartCover.Format
		}

		if mediaSmartCoverTemplate.SmartCover.Width != "" {
			smartCoverMap["width"] = mediaSmartCoverTemplate.SmartCover.Width
		}

		if mediaSmartCoverTemplate.SmartCover.Height != "" {
			smartCoverMap["height"] = mediaSmartCoverTemplate.SmartCover.Height
		}

		if mediaSmartCoverTemplate.SmartCover.Count != "" {
			smartCoverMap["count"] = mediaSmartCoverTemplate.SmartCover.Count
		}

		if mediaSmartCoverTemplate.SmartCover.DeleteDuplicates != "" {
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
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := cos.CreateMediaSmartCoverTemplateOptions{
		Tag: "SmartCover",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if d.HasChange("smart_cover") {
		if dMap, ok := helper.InterfacesHeadMap(d, "smart_cover"); ok {
			smartCover := cos.NodeSmartCover{}
			if v, ok := dMap["format"]; ok {
				smartCover.Format = v.(string)
			}
			if v, ok := dMap["width"]; ok {
				smartCover.Width = v.(string)
			}
			if v, ok := dMap["height"]; ok {
				smartCover.Height = v.(string)
			}
			if v, ok := dMap["count"]; ok {
				smartCover.Count = v.(string)
			}
			if v, ok := dMap["delete_duplicates"]; ok {
				smartCover.DeleteDuplicates = v.(string)
			}
			request.SmartCover = &smartCover
		}
	}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient(bucket).CI.UpdateMediaSmartCoverTemplate(ctx, &request, templateId)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaSmartCoverTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSmartCoverTemplate failed, reason:%+v", logId, err)
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
