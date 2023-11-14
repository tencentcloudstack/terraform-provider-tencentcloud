/*
Provides a resource to create a live timeshift_template

Example Usage

```hcl
resource "tencentcloud_live_timeshift_template" "timeshift_template" {
  template_name = ""
  duration =
  description = ""
  area = ""
  item_duration =
  remove_watermark =
  transcode_template_ids =
}
```

Import

live timeshift_template can be imported using the id, e.g.

```
terraform import tencentcloud_live_timeshift_template.timeshift_template timeshift_template_id
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

func resourceTencentCloudLiveTimeshiftTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLiveTimeshiftTemplateCreate,
		Read:   resourceTencentCloudLiveTimeshiftTemplateRead,
		Update: resourceTencentCloudLiveTimeshiftTemplateUpdate,
		Delete: resourceTencentCloudLiveTimeshiftTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"template_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name.Maximum length: 255 bytes.Only letters, numbers, underscores, and hyphens are supported.",
			},

			"duration": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The time shifting duration.Unit: Second.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The template description.Only letters, numbers, underscores, and hyphens are supported.",
			},

			"area": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The region.`Mainland`: The Chinese mainland.`Overseas`: Outside the Chinese mainland.Default value: `Mainland`.",
			},

			"item_duration": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "The segment size.Value range: 3-10.Unit: Second.Default value: 5.",
			},

			"remove_watermark": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to remove watermarks.If you pass in `true`, the original stream will be recorded.Default value: `false`.",
			},

			"transcode_template_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "The transcoding template IDs.This API works only if `RemoveWatermark` is `false`.",
			},
		},
	}
}

func resourceTencentCloudLiveTimeshiftTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_timeshift_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = live.NewCreateLiveTimeShiftTemplateRequest()
		response   = live.NewCreateLiveTimeShiftTemplateResponse()
		templateId int
	)
	if v, ok := d.GetOk("template_name"); ok {
		request.TemplateName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("duration"); ok {
		request.Duration = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("item_duration"); ok {
		request.ItemDuration = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("remove_watermark"); ok {
		request.RemoveWatermark = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("transcode_template_ids"); ok {
		transcodeTemplateIdsSet := v.(*schema.Set).List()
		for i := range transcodeTemplateIdsSet {
			transcodeTemplateIds := transcodeTemplateIdsSet[i].(int)
			request.TranscodeTemplateIds = append(request.TranscodeTemplateIds, helper.IntInt64(transcodeTemplateIds))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().CreateLiveTimeShiftTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create live timeshiftTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudLiveTimeshiftTemplateRead(d, meta)
}

func resourceTencentCloudLiveTimeshiftTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_timeshift_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	timeshiftTemplateId := d.Id()

	timeshiftTemplate, err := service.DescribeLiveTimeshiftTemplateById(ctx, templateId)
	if err != nil {
		return err
	}

	if timeshiftTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LiveTimeshiftTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if timeshiftTemplate.TemplateName != nil {
		_ = d.Set("template_name", timeshiftTemplate.TemplateName)
	}

	if timeshiftTemplate.Duration != nil {
		_ = d.Set("duration", timeshiftTemplate.Duration)
	}

	if timeshiftTemplate.Description != nil {
		_ = d.Set("description", timeshiftTemplate.Description)
	}

	if timeshiftTemplate.Area != nil {
		_ = d.Set("area", timeshiftTemplate.Area)
	}

	if timeshiftTemplate.ItemDuration != nil {
		_ = d.Set("item_duration", timeshiftTemplate.ItemDuration)
	}

	if timeshiftTemplate.RemoveWatermark != nil {
		_ = d.Set("remove_watermark", timeshiftTemplate.RemoveWatermark)
	}

	if timeshiftTemplate.TranscodeTemplateIds != nil {
		_ = d.Set("transcode_template_ids", timeshiftTemplate.TranscodeTemplateIds)
	}

	return nil
}

func resourceTencentCloudLiveTimeshiftTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_timeshift_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := live.NewModifyLiveTimeShiftTemplateRequest()

	timeshiftTemplateId := d.Id()

	request.TemplateId = &templateId

	immutableArgs := []string{"template_name", "duration", "description", "area", "item_duration", "remove_watermark", "transcode_template_ids"}

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

	if d.HasChange("duration") {
		if v, ok := d.GetOkExists("duration"); ok {
			request.Duration = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	if d.HasChange("area") {
		if v, ok := d.GetOk("area"); ok {
			request.Area = helper.String(v.(string))
		}
	}

	if d.HasChange("item_duration") {
		if v, ok := d.GetOkExists("item_duration"); ok {
			request.ItemDuration = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("remove_watermark") {
		if v, ok := d.GetOkExists("remove_watermark"); ok {
			request.RemoveWatermark = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("transcode_template_ids") {
		if v, ok := d.GetOk("transcode_template_ids"); ok {
			transcodeTemplateIdsSet := v.(*schema.Set).List()
			for i := range transcodeTemplateIdsSet {
				transcodeTemplateIds := transcodeTemplateIdsSet[i].(int)
				request.TranscodeTemplateIds = append(request.TranscodeTemplateIds, helper.IntInt64(transcodeTemplateIds))
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLiveClient().ModifyLiveTimeShiftTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update live timeshiftTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudLiveTimeshiftTemplateRead(d, meta)
}

func resourceTencentCloudLiveTimeshiftTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_live_timeshift_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}
	timeshiftTemplateId := d.Id()

	if err := service.DeleteLiveTimeshiftTemplateById(ctx, templateId); err != nil {
		return err
	}

	return nil
}
