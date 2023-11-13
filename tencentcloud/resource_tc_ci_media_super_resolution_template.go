/*
Provides a resource to create a ci media_super_resolution_template

Example Usage

```hcl
resource "tencentcloud_ci_media_super_resolution_template" "media_super_resolution_template" {
  name = &lt;nil&gt;
  resolution = &lt;nil&gt;
  enable_scale_up = &lt;nil&gt;
  version = &lt;nil&gt;
}
```

Import

ci media_super_resolution_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_super_resolution_template.media_super_resolution_template media_super_resolution_template_id
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

func resourceTencentCloudCiMediaSuperResolutionTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaSuperResolutionTemplateCreate,
		Read:   resourceTencentCloudCiMediaSuperResolutionTemplateRead,
		Update: resourceTencentCloudCiMediaSuperResolutionTemplateUpdate,
		Delete: resourceTencentCloudCiMediaSuperResolutionTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"resolution": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Resolution Options sdtohd: Standard Definition to Ultra Definition, hdto4k: HD to 4K.",
			},

			"enable_scale_up": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Auto scaling switch, off by default.",
			},

			"version": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Version, default value Base, Base: basic version, Enhance: enhanced version.",
			},
		},
	}
}

func resourceTencentCloudCiMediaSuperResolutionTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_super_resolution_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = ci.NewCreateMediaSuperResolutionTemplateRequest()
		response   = ci.NewCreateMediaSuperResolutionTemplateResponse()
		templateId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resolution"); ok {
		request.Resolution = helper.String(v.(string))
	}

	if v, ok := d.GetOk("enable_scale_up"); ok {
		request.EnableScaleUp = helper.String(v.(string))
	}

	if v, ok := d.GetOk("version"); ok {
		request.Version = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().CreateMediaSuperResolutionTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaSuperResolutionTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(templateId)

	return resourceTencentCloudCiMediaSuperResolutionTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSuperResolutionTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_super_resolution_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	mediaSuperResolutionTemplateId := d.Id()

	mediaSuperResolutionTemplate, err := service.DescribeCiMediaSuperResolutionTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if mediaSuperResolutionTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiMediaSuperResolutionTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if mediaSuperResolutionTemplate.Name != nil {
		_ = d.Set("name", mediaSuperResolutionTemplate.Name)
	}

	if mediaSuperResolutionTemplate.Resolution != nil {
		_ = d.Set("resolution", mediaSuperResolutionTemplate.Resolution)
	}

	if mediaSuperResolutionTemplate.EnableScaleUp != nil {
		_ = d.Set("enable_scale_up", mediaSuperResolutionTemplate.EnableScaleUp)
	}

	if mediaSuperResolutionTemplate.Version != nil {
		_ = d.Set("version", mediaSuperResolutionTemplate.Version)
	}

	return nil
}

func resourceTencentCloudCiMediaSuperResolutionTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_super_resolution_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewUpdateMediaSuperResolutionTemplateRequest()

	mediaSuperResolutionTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"name", "resolution", "enable_scale_up", "version"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().UpdateMediaSuperResolutionTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci mediaSuperResolutionTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaSuperResolutionTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaSuperResolutionTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_media_super_resolution_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	mediaSuperResolutionTemplateId := d.Id()

	if err := service.DeleteCiMediaSuperResolutionTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
