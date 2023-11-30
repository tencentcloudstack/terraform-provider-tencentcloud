package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudCssTimeshiftTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssTimeshiftTemplateCreate,
		Read:   resourceTencentCloudCssTimeshiftTemplateRead,
		Update: resourceTencentCloudCssTimeshiftTemplateUpdate,
		Delete: resourceTencentCloudCssTimeshiftTemplateDelete,
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

func resourceTencentCloudCssTimeshiftTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_timeshift_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = css.NewCreateLiveTimeShiftTemplateRequest()
		response   = css.NewCreateLiveTimeShiftTemplateResponse()
		templateId int64
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().CreateLiveTimeShiftTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css timeshiftTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = *response.Response.TemplateId
	d.SetId(helper.Int64ToStr(templateId))

	return resourceTencentCloudCssTimeshiftTemplateRead(d, meta)
}

func resourceTencentCloudCssTimeshiftTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_timeshift_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	templateId := d.Id()
	templateIdInt64, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		return fmt.Errorf("TemplateId format type error: %s", err.Error())
	}

	timeshiftTemplate, err := service.DescribeCssTimeshiftTemplateById(ctx, templateIdInt64)
	if err != nil {
		return err
	}

	if timeshiftTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssTimeshiftTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
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

func resourceTencentCloudCssTimeshiftTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_timeshift_template.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := css.NewModifyLiveTimeShiftTemplateRequest()

	templateId := d.Id()
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)
	templateIdUint := uint64(templateIdInt64)

	request.TemplateId = &templateIdUint

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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCssClient().ModifyLiveTimeShiftTemplate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update css timeshiftTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCssTimeshiftTemplateRead(d, meta)
}

func resourceTencentCloudCssTimeshiftTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_css_timeshift_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}
	templateId := d.Id()
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)

	if err := service.DeleteCssTimeshiftTemplateById(ctx, templateIdInt64); err != nil {
		return err
	}

	return nil
}
