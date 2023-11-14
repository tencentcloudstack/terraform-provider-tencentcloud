/*
Provides a resource to create a live pad_template

Example Usage

```hcl
resource "tencentcloud_live_pad_template" "pad_template" {
  template_name = ""
  url = ""
  description = ""
  wait_duration =
  max_duration =
  type =
}
```

Import

live pad_template can be imported using the id, e.g.

```
terraform import tencentcloud_live_pad_template.pad_template pad_template_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudLivePadTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLivePadTemplateCreate,
		Read:   resourceTencentCloudLivePadTemplateRead,
		Update: resourceTencentCloudLivePadTemplateUpdate,
		Delete: resourceTencentCloudLivePadTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Template namelimit 255 bytes.",
			},

			"url": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Pad content.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Description content.limit length 1024 bytes.",
			},

			"wait_duration": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Stop stream wait time.limit: 0 - 30000 ms.",
			},

			"max_duration": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Max pad duration.limit: 0 - 9999999 ms.",
			},

			"type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Pad content type.1: picture.2: video.default: 1.",
			},
		},
	}
}

func resourceTencentCloudLivePadTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_pad_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewCreateLivePadTemplateRequest()
		response   = live.NewCreateLivePadTemplateResponse()
		templateId int
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("url"); ok {
		request.Url = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("wait_duration"); ok {
		request.WaitDuration = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("max_duration"); ok {
		request.MaxDuration = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("type"); ok {
		request.Type = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateLivePadTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live padTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudLivePadTemplateRead(d, meta)
}

func resourceTencentCloudLivePadTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_pad_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	padTemplateId := d.Id()

	padTemplate, err := service.DescribeLivePadTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if padTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LivePadTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if padTemplate.TemplateName != nil {
		_ = d.Set("template_name", padTemplate.TemplateName)
	}

	if padTemplate.Url != nil {
		_ = d.Set("url", padTemplate.Url)
	}

	if padTemplate.Description != nil {
		_ = d.Set("description", padTemplate.Description)
	}

	if padTemplate.WaitDuration != nil {
		_ = d.Set("wait_duration", padTemplate.WaitDuration)
	}

	if padTemplate.MaxDuration != nil {
		_ = d.Set("max_duration", padTemplate.MaxDuration)
	}

	if padTemplate.Type != nil {
		_ = d.Set("type", padTemplate.Type)
	}

	return nil
}

func resourceTencentCloudLivePadTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_pad_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLivePadTemplateRequest()

	padTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"template_name", "url", "description", "wait_duration", "max_duration", "type"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("template_name") {
		if v, ok := d.GetOk("template_name"); ok {
			request.TemplateName = helper.String(v.(string))
		}
	}

	if d.HasChange("url") {
		if v, ok := d.GetOk("url"); ok {
			request.Url = helper.String(v.(string))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("wait_duration") {
		if v, ok := d.GetOkExists("wait_duration"); ok {
			request.WaitDuration = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_duration") {
		if v, ok := d.GetOkExists("max_duration"); ok {
			request.MaxDuration = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("type") {
		if v, ok := d.GetOkExists("type"); ok {
			request.Type = helper.IntUint64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLivePadTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live padTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLivePadTemplateRead(d, meta)
}

func resourceTencentCloudLivePadTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_pad_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	padTemplateId := d.Id()

	if err := service.DeleteLivePadTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
