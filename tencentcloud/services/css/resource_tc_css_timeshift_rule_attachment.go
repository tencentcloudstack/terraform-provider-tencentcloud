package css

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCssTimeshiftRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCssTimeshiftRuleAttachmentCreate,
		Read:   resourceTencentCloudCssTimeshiftRuleAttachmentRead,
		Delete: resourceTencentCloudCssTimeshiftRuleAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The push domain.",
			},

			"app_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The push path, which should be the same as `AppName` in the push and playback URLs. The default value is `live`.",
			},

			"stream_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The stream name.Note: If you pass in a non-empty string, the rule will only be applied to the specified stream.",
			},

			"template_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The template ID.",
			},
		},
	}
}

func resourceTencentCloudCssTimeshiftRuleAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_timeshift_rule_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = css.NewCreateLiveTimeShiftRuleRequest()
		templateId int
		domainName string
	)
	if v, ok := d.GetOk("domain_name"); ok {
		domainName = v.(string)
		request.DomainName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("app_name"); ok {
		request.AppName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("stream_name"); ok {
		request.StreamName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("template_id"); ok {
		templateId = v.(int)
		request.TemplateId = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCssClient().CreateLiveTimeShiftRule(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create css timeshiftRuleAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{helper.IntToStr(templateId), domainName}, tccommon.FILED_SP))

	return resourceTencentCloudCssTimeshiftRuleAttachmentRead(d, meta)
}

func resourceTencentCloudCssTimeshiftRuleAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_timeshift_rule_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	templateId := idSplit[0]
	templateIdInt64, err := strconv.ParseInt(templateId, 10, 64)
	if err != nil {
		return fmt.Errorf("TemplateId format type error: %s", err.Error())
	}
	domainName := idSplit[1]

	timeshiftRuleAttachment, err := service.DescribeCssTimeshiftRuleAttachmentById(ctx, templateIdInt64, domainName)
	if err != nil {
		return err
	}

	if timeshiftRuleAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CssTimeshiftRuleAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if timeshiftRuleAttachment.DomainName != nil {
		_ = d.Set("domain_name", timeshiftRuleAttachment.DomainName)
	}

	if timeshiftRuleAttachment.AppName != nil {
		_ = d.Set("app_name", timeshiftRuleAttachment.AppName)
	}

	if timeshiftRuleAttachment.StreamName != nil {
		_ = d.Set("stream_name", timeshiftRuleAttachment.StreamName)
	}

	if timeshiftRuleAttachment.TemplateId != nil {
		_ = d.Set("template_id", timeshiftRuleAttachment.TemplateId)
	}

	return nil
}

func resourceTencentCloudCssTimeshiftRuleAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_css_timeshift_rule_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	templateId := idSplit[0]
	templateIdInt64, _ := strconv.ParseInt(templateId, 10, 64)
	domainName := idSplit[1]

	snapshotRule, err := service.DescribeCssTimeshiftRuleAttachmentById(ctx, templateIdInt64, domainName)
	if err != nil {
		return err
	}

	appName := ""
	if snapshotRule.AppName != nil {
		appName = *snapshotRule.AppName
	}

	streamName := ""
	if snapshotRule.StreamName != nil {
		streamName = *snapshotRule.StreamName
	}

	if err := service.DeleteCssTimeshiftRuleAttachmentById(ctx, domainName, appName, streamName); err != nil {
		return err
	}

	return nil
}
